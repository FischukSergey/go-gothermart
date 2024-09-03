package withdrawals

import (
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
)

type GetUserWithdrawAll interface {
	GetAllWithdraw(ctx context.Context, userID int) ([]models.GetAllWithdraw, error)
}

// OrderWithdrawAll выдача всех заказов со списанием баллов.
// Пользователь должен быть авторизован.
func OrderWithdrawAll(log *slog.Logger, storage GetUserWithdrawAll) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Get all withdraw orders")
		defer log.Debug("Get all withdraw orders finished")

		w.Header().Set("Content-Type", "application/json")
		userID := r.Context().Value(auth.CtxKeyUser).(int)

		//ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		//defer cancel()
		//пишем заказ в базу и обрабатываем ошибку, если есть
		var err error
		var resp []models.GetAllWithdraw
		resp, err = storage.GetAllWithdraw(r.Context(), userID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("Error found order withdraw", logger.Err(err))
			return
		}

		if len(resp) == 0 { //нет данных для ответа
			w.WriteHeader(http.StatusNoContent)
			log.Info("not found withdraw orders for user")
			render.JSON(w, r, "not found withdraw orders for user")
			return
		}

		render.JSON(w, r, resp)
		log.Info("get orders withdraw successfully", slog.String("user_id", strconv.Itoa(userID)))
		w.WriteHeader(http.StatusOK)
	}
}
