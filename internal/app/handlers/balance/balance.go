package balance

import (
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Balancer interface {
	GetUserBalance(ctx context.Context, userID int) (current, withdrawn float32, err error)
}

// GetBalance получение текущего баланса и суммы всех списаний пользователя.
// Пользователь должен быть авторизован.
func GetBalance(log *slog.Logger, storage Balancer) http.HandlerFunc {
	type response struct {
		Current   float32 `json:"current"`
		Withdrawn float32 `json:"withdrawn"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting balance and withdrawn")
		defer log.Debug("Getting balance and withdrawn finished")

		w.Header().Set("Content-Type", "application/json")
		userID := r.Context().Value(auth.CtxKeyUser).(int)

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		//пишем заказ в базу и обрабатываем ошибку, если есть
		var resp response
		var err error
		resp.Current, resp.Withdrawn, err = storage.GetUserBalance(ctx, userID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error("Error getting balance and withdrawn: ", err)
			render.JSON(w, r, map[string]string{"error": err.Error()})
			return
		}

		render.JSON(w, r, resp) //пишем JSON ответ
		w.WriteHeader(http.StatusOK)
		log.Info("Balance and withdrawn successful", slog.String("user_id", strconv.Itoa(userID)))
	}
}
