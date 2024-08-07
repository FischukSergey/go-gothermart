package main

import (
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/login"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/orders"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/register"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	mwlogger "github.com/FischukSergey/go-gothermart.git/internal/app/middleware/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/storage"
	stdlog "log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/jackc/pgconn"
)

const ( //уровни логирования
	envLocal = "local" //уровень по умолчанию
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	ParseFlags()                        //инициализируем флаги/переменные окружения конфигурации сервера
	log := setupLogger(FlagLevelLogger) //инициализируем логер с заданным уровнем

	r := chi.NewRouter() //инициализируем роутер и middleware

	var DatabaseDSN *pgconn.Config //инициализируем базу данных
	DatabaseDSN, err := pgconn.ParseConfig(FlagDatabaseDSN)
	if err != nil {
		stdlog.Fatal("Ошибка парсинга строки инициализации БД Postgres")
	}

	storageDB, err := storage.NewDB(DatabaseDSN, log)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer storageDB.DB.Close()
	log.Info("database connection", slog.String("database", DatabaseDSN.Database))

	//инициализируем middleware
	r.Use(mwlogger.NewMwLogger(log)) //маршрут в middleware за логированием
	//r.Use(gzipper.NewMwGzipper(log)) //работа со сжатыми запросами/сжатие ответов
	r.Use(auth.AuthToken(log)) //ID session аутентификация пользователя/JWToken в  cookie
	//инициализируем хендлеры
	r.Post("/api/user/register", register.Register(log, storageDB))
	r.Post("/api/user/login", login.LoginAuth(log, storageDB))
	r.Post("/api/user/orders", orders.OrderSave(log, storageDB))

	srv := &http.Server{ //запускаем сервер
		Addr:         FlagServerPort,
		Handler:      r,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Info("Initializing server", slog.String("address", srv.Addr))

	if err := srv.ListenAndServe(); err != nil {
		fmt.Printf("Ошибка при запуске сервера: %s", err.Error())
		return
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

/*
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
*/
