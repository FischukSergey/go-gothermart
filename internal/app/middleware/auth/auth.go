package auth

import (
	"context"
	jwtoken "github.com/FischukSergey/go-gothermart.git/internal/lib/jwt"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"log/slog"
	"net/http"
	"strconv"
)

type ctxKey int

const (
	CtxKeyUser ctxKey = iota + 1
)

func AuthToken(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {

		log.Debug("middleware authorize started")

		Authorize := func(w http.ResponseWriter, r *http.Request) {
			//если api для регистрации или авторизации, то ничего не делаем и передаем дальше
			if r.RequestURI == "/api/user/login" || r.RequestURI == "/api/user/register" {
				next.ServeHTTP(w, r)
			} else { //если любой другой api-проверяем на валидность, извлекаем id  и помещаем в контекст

				var userID int

				tokenCookie, err := r.Cookie("token")
				if err != nil { //токен не нашелся
					log.Error("could not get token cookie", logger.Err(err))
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				//если токен есть, проверим на валидность и получим ID
				userID = jwtoken.GetJWTokenUserID(tokenCookie.Value)
				if userID == -1 { //токен есть, но не валиден
					log.Error("invalid token or absent id user")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				//если все успешно - пишем в контекст ID пользователя
				log.Info("token validate, user authorized", slog.String("user ID:", strconv.Itoa(userID)))
				next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), CtxKeyUser, userID)))
			}
		}
		return http.HandlerFunc(Authorize)
	}
}
