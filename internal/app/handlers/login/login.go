package login

import (
	"context"
	"github.com/FischukSergey/go-gothermart.git/internal/lib/jwt"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Loginer interface {
	Login(ctx context.Context, email string) (*models.User, error)
}

func LoginAuth(log *slog.Logger, storage Loginer) http.HandlerFunc {
	type request struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		log.Debug("registering user")

		req := &request{}
		if err := render.Decode(r, req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error("Bad request", "error", err)
			return
		}
		u := &models.User{
			Email:    req.Login,
			Password: req.Password,
		}

		//проводим валидацию логина и пароля
		err := u.Validate()
		if err != nil {
			log.Error("login or password failure", logger.Err(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()
		//ищем в базе логин и пароль, в случае успеха получаем объект user
		user, err := storage.Login(ctx, u.Email)

		if err != nil {
			log.Error("login or password failure", logger.Err(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//проверяем пароль
		if err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(req.Password)); err != nil {
			log.Error("password failure", logger.Err(err))
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		//создаем токен соединения
		token, err := jwtoken.NewToken(*user)

		if err != nil {
			log.Error("can't create JWToken", logger.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//w.Header().Set("Authorization", "Bearer "+token)
		cookie := &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(jwtoken.ExpiresKey * time.Hour),
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)

		log.Info("user login successfully",
			slog.String("email", user.Email),
			slog.String("uid", strconv.Itoa(user.ID)),
		)
	}
}
