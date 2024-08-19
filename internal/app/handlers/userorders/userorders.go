package userorders

import (
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"time"
)

type UserOrdersGetter interface {
	GetUserOrders(ctx context.Context, id int) ([]models.GetUserOrders, error)
}

func UserOrders(log *slog.Logger, storage UserOrdersGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("UserOrders started")
		defer log.Debug("UserOrders finished")

		w.Header().Set("Content-Type", "application/json")
		userID := r.Context().Value(auth.CtxKeyUser).(int)

		var orders []models.GetUserOrders

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		//получаем ордера из базы и обрабатываем ошибку, если есть
		orders, err := storage.GetUserOrders(ctx, userID)
		if err != nil { //ошибка работы БД
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error getting user orders: ", err)
			//render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}
		if len(orders) == 0 { //нет данных для ответа
			w.WriteHeader(http.StatusNoContent)
			log.Info("not found data for user")
			//render.JSON(w, r, "not found data for user")
			return
		}

		render.JSON(w, r, orders) //пишем JSON ответ
		w.WriteHeader(http.StatusOK)
		log.Info("UserOrders finished successfully", "user", userID)
	}
}
