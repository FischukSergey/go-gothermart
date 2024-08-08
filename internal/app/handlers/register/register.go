package register

import (
	"context"
	"errors"
	jwtoken "github.com/FischukSergey/go-gothermart.git/internal/lib/jwt"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/go-chi/render"
)

type UserRegister interface {
	Register(ctx context.Context, u *models.User) (id int, err error)
}

func Register(log *slog.Logger, storage UserRegister) http.HandlerFunc {
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
		if u.Email == "" || u.Password == "" { //не должны быть пустыми
			w.WriteHeader(http.StatusBadRequest)
			log.Error("Bad request, login or password empty")
			return
		}
		//проводим валидацию логина и пароля
		err := u.Validate()
		if err != nil {
			log.Error("login or password failure", logger.Err(err))
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		//кодируем пароль
		passHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error("failed to generate password hash", "error", err)
			return
		}

		u.EncryptedPassword = string(passHash)

		ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
		defer cancel()

		//вызываем метод записи в БД
		id, err := storage.Register(ctx, u)
		//обработка ошибки вставки уже существующего Login (email)
		var res []string
		if errors.Is(err, models.ErrUserExists) {
			res = strings.Split(err.Error(), ":")
			http.Error(w, "request failed, login exists", http.StatusConflict)
			log.Error("Request create user failed, login exists",
				slog.String("email:", res[0]),
			)
			return
		}

		if err != nil {
			log.Error("failed to register user", "error", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		u.ID = id
		//создаем токен соединения
		token, err := jwtoken.NewToken(*u)
		if err != nil {
			log.Error("can't create JWToken", logger.Err(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		//w.Header().Set("Authorization", "Bearer "+token)
		cookie := &http.Cookie{
			Name:    "token",
			Value:   token,
			Expires: time.Now().Add(72 * time.Hour),
		}
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)

		log.Info("user register successfully",
			slog.String("email", u.Email),
			slog.String("uid", strconv.Itoa(u.ID)),
		)
	}
}
