package main

import (
	"context"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/balance"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/login"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/orders"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/register"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/userorders"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/withdraw"
	"github.com/FischukSergey/go-gothermart.git/internal/app/handlers/withdrawals"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/auth"
	"github.com/FischukSergey/go-gothermart.git/internal/app/middleware/gzipper"
	mwlogger "github.com/FischukSergey/go-gothermart.git/internal/app/middleware/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/app/services"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/FischukSergey/go-gothermart.git/internal/storage"
	stdlog "log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
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

	rout := chi.NewRouter() //инициализируем роутер и middleware

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

	defer func() {
		storageDB.DB.Close()
		log.Info("database closed")
	}()

	log.Info("database connection", slog.String("database", DatabaseDSN.Database))

	rout.Route("/api/user", func(r chi.Router) {
		//инициализируем middleware
		r.Use(mwlogger.NewMwLogger(log)) //маршрут в middleware за логированием
		r.Use(gzipper.NewMwGzipper(log)) //работа со сжатыми запросами/сжатие ответов
		r.Use(auth.AuthToken(log))       //ID session аутентификация пользователя/JWToken в  cookie

		//инициализируем хендлеры
		r.Post("/register", register.Register(log, storageDB))
		r.Post("/login", login.LoginAuth(log, storageDB))
		r.Post("/balance/withdraw", withdraw.OrderWithdraw(log, storageDB))
		r.Post("/orders", orders.OrderSave(log, storageDB))
		r.Get("/orders", userorders.UserOrders(log, storageDB))
		r.Get("/balance", balance.GetBalance(log, storageDB))
		r.Get("/withdrawals", withdrawals.OrderWithdrawAll(log, storageDB))
	})

	//инициализируем сервис и сервер расчета баллов (accrual)
	ctx, cancel := context.WithCancel(context.Background())

	accrual := models.Accrual{
		MaxWorker:            3,
		TimeTicker:           5,
		MaxRetries:           3,
		AccrualServerAddress: FlagAccrualSystemAddress,
	}
	var wg sync.WaitGroup
	go services.AccrualService(ctx, accrual, storageDB, log, &wg)

	srv := &http.Server{ //запускаем сервер
		Addr:         FlagServerPort,
		Handler:      rout,
		ReadTimeout:  4 * time.Second,
		WriteTimeout: 4 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Info("Initializing server")
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Error("Ошибка при запуске сервера", logger.Err(err))
		}
	}()
	log.Info("Server started", slog.String("address", srv.Addr))

	//Остановка процессов
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-done //ждем сигнал прерывания
	//останавливаем сервер на прием новых запросов и дорабатываем принятые
	log.Info("Server stopping", slog.String("address", srv.Addr))
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", logger.Err(err))
		return
	}
	log.Info("api server stopped")

	cancel()  //тормозим по контексту accrual сервис
	wg.Wait() //дожидаемся отработки горутин воркера
	time.Sleep(1 * time.Second)
	log.Info("accrual server stopped")
	//последним по defer закрываем базу данных
}

// setupLogger() настройка уровня доступа к логам из переменной среды
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
