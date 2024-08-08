package orders

import (
	"context"
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/lib/luhn"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type OrderSaver interface {
	CreateOrder(ctx context.Context, order models.Order) error
}

// OrderSave сохранение номера заказа. Пользователь должен быть авторизован.
func OrderSave(log *slog.Logger, storage OrderSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Saving order")

		userID := r.Context().Value(auth.CtxKeyUser).(int)
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Request bad", logger.Err(err))
			return
		}
		if string(body) == "" {
			http.Error(w, "Request is empty", http.StatusBadRequest)
			log.Error("Request is empty")
			return
		}

		//проверяем номер заказа на контрольную сумму
		if !luhn.Valid(string(body)) {
			http.Error(w, "Request is invalid", http.StatusUnprocessableEntity)
			log.Error("Request is invalid")
			return
		}
		order := models.Order{
			UserID:    userID,
			OrderID:   string(body),
			Status:    "NEW",
			CreatedAt: time.Now().Format(time.RFC3339),
		}

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		//пишем заказ в базу и обрабатываем ошибку, если есть
		err = storage.CreateOrder(ctx, order)
		if err != nil {
			if errors.Is(err, models.ErrOrderUploadedSameUser) {
				http.Error(w, err.Error(), http.StatusOK)
				return
			}
			if errors.Is(err, models.ErrOrderUploadedAnotherUser) {
				http.Error(w, err.Error(), http.StatusConflict)
				return
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error creating order", logger.Err(err))
			return
		}

		w.WriteHeader(http.StatusAccepted)
	}
}
