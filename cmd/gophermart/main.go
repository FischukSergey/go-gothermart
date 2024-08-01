package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"log/slog"
	"net/http"
	"os"
	"time"
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
