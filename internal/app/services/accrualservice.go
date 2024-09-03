package services

import (
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type OrderUpdate struct { //структура для массива новых записей с id пользователя
	OrderID string
	UserID  int
}

type AccrualServices interface {
	GetAccrualOrders(ctx context.Context) ([]OrderUpdate, error)
	UpdateAccrualOrder(ctx context.Context, accrualOrder models.AccrualOrder, userID int) error
}

// AccrualService в фоновом режиме проводит запросы к микросервису начисления баллов лояльности
// и обновляет данные в основном сервисе заказов
func AccrualService(ctx context.Context, accrual models.Accrual, storage AccrualServices, log *slog.Logger, wg *sync.WaitGroup) {
	workPool := make(chan struct{}, accrual.MaxWorker)
	defer close(workPool)

	//var wg sync.WaitGroup

	log = log.With(slog.String("Service", "Accrual"))

	ticker := time.NewTicker(time.Second * time.Duration(accrual.TimeTicker))
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var orders []OrderUpdate
			orders, err := storage.GetAccrualOrders(ctx) //запрос в БД за слайсом необработанных заказов
			if err != nil {
				log.Error("Failed to get accrual orders", "error", err)
			}
			if len(orders) == 0 {
				continue
			}
			for _, order := range orders {
				workPool <- struct{}{}
				wg.Add(1)
				log.Info("Getting orders status and accrual for order " + order.OrderID)
				go processedAccrual(ctx, workPool, order, storage, log, accrual, wg)
				//time.Sleep(1 * time.Second)
			}
		case <-ctx.Done():
			log.Info("Accrual service shutting down")
			return
		}
	}
}

func processedAccrual(ctx context.Context, workPool chan struct{}, order OrderUpdate,
	storage AccrualServices, log *slog.Logger, accrual models.Accrual, wg *sync.WaitGroup) {
	defer func() {
		<-workPool //освобождаем буфер для запуска следующей горутины по завершению работы функции
		wg.Done()  //освобождаем счетчик работающих горутин
	}()
	//time.Sleep(3 * time.Second) //проверка работы wg.Wait
	//TODO 1.Проработать замечания ментора. 2.Timeout на ответ сервиса Accrual?

	uri := accrual.AccrualServerAddress + "/api/orders/" + order.OrderID

	for retries := accrual.MaxRetries; retries > 0; retries-- { //делаем несколько попыток получить данные о заказе
		res, err := http.Get(uri)
		if err != nil {
			log.Error("accrual service http get error", "err", err)
			return
		}

		switch res.StatusCode {
		case http.StatusInternalServerError:
			log.Error("accrual service http error", "err", res.Body)
			time.Sleep(time.Second * 1) //попытаемся еще раз
			closeBody(res.Body, log)
			continue

		case http.StatusNoContent:
			log.Info("accrual service http no content", slog.String("order", order.OrderID))
			closeBody(res.Body, log)
			return

		case http.StatusTooManyRequests:
			timeout := res.Header.Get("Retry-After")
			log.Info("accrual service http retry timeout", slog.String("timeout", timeout))
			if timeout == "" {
				timeout = "10"
			}
			t, err := strconv.Atoi(timeout)
			if err != nil {
				log.Error("timeout accrual service error", "err", err)
				closeBody(res.Body, log)
				break
			}
			time.Sleep(time.Duration(t) * time.Second) //todo таймаут для всех горутин????
			closeBody(res.Body, log)
			continue

		case http.StatusOK:
			var updatedOrderStatus models.AccrualOrder
			if err := render.DecodeJSON(res.Body, &updatedOrderStatus); err != nil {
				log.Error("error unmarshal accrual order status", "err", err)
				closeBody(res.Body, log)
				break
			}
			if updatedOrderStatus.Status == models.StatusOrderProcessed ||
				updatedOrderStatus.Status == models.StatusOrderInvalid { //если обработка заказа окончена
				err := storage.UpdateAccrualOrder(ctx, updatedOrderStatus, order.UserID) //пишем в БД
				if err != nil {
					log.Error("error update accrual order status", "err", err)
					closeBody(res.Body, log)
					break
				}
				log.Info("accrual order updated", "order", updatedOrderStatus.Order)
			} else {
				log.Info("accrual order processing", "order", updatedOrderStatus.Order)
			}

			closeBody(res.Body, log)
			return

		default:
			log.Error("unexpected accrual service http code")
			closeBody(res.Body, log)
			return
		}

		err = res.Body.Close()
		if err != nil {
			log.Error("error close body GET", "err", err)
		}

	}
}

func closeBody(body io.ReadCloser, log *slog.Logger) {
	err := body.Close()
	if err != nil {
		log.Error("error close body GET", "err", err)
	}
}
