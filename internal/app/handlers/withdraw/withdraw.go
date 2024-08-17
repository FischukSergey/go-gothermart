package withdraw

import (
	"context"
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/lib/luhn"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type OrderBalanceWithdraw interface {
	CreateOrderWithdraw(ctx context.Context, order models.Order) error
}

// OrderWithdraw сохранение номера заказа со списанием баллов.
// Пользователь должен быть авторизован.
func OrderWithdraw(log *slog.Logger, storage OrderBalanceWithdraw) http.HandlerFunc {
	type request struct {
		Order string  `json:"order"`
		Sum   float32 `json:"sum"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Saving withdraw order")
		defer log.Debug("Saving withdraw order finished")

		userID := r.Context().Value(auth.CtxKeyUser).(int)

		req := &request{}
		if err := render.Decode(r, req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error("Bad request", "error", err)
			return
		}
		//проверяем номер заказа на контрольную сумму
		if !luhn.Valid(req.Order) {
			http.Error(w, "Number order is invalid", http.StatusUnprocessableEntity)
			log.Error("Number order is invalid")
			return
		}

		orderWithdraw := models.Order{
			UserID:    userID,
			OrderID:   req.Order,
			Withdraw:  req.Sum,
			Status:    "",
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		//пишем заказ в базу и обрабатываем ошибку, если есть
		err := storage.CreateOrderWithdraw(ctx, orderWithdraw)
		if err != nil {
			if errors.Is(err, models.ErrInsufficientFunds) {
				http.Error(w, err.Error(), http.StatusPaymentRequired)
				return
			}
			if errors.Is(err, models.ErrOrderExists) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error creating order withdraw", logger.Err(err))
			return
		}

		log.Info("Order withdraw successfully created", slog.String("order_id", req.Order))
		w.WriteHeader(http.StatusOK)
	}
}
