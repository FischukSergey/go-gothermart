package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// CreateOrderWithdraw добавление новой записи(ордера) заказа на списание баллов
// проверяет на достаточность средств на списание (balance in users)
func (db *PostgresqlDB) CreateOrderWithdraw(ctx context.Context, order models.Order) error {
	const op = "postgresql.CreateOrderWithdraw"
	log := db.logger.With(
		slog.String("op", op),
		slog.String("order", order.OrderID),
	)

	//проверка на уже существующий в базе номер заказа
	row := db.DB.QueryRow(ctx, "SELECT user_id, order_num FROM orders WHERE order_num=$1", order.OrderID)
	var useOrder models.Order
	err := row.Scan(
		&useOrder.UserID,
		&useOrder.OrderID,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) { //если заказа не существует, все хорошо и будем делать запись
			log.Info("order is new and will be created", slog.String("номер заказа:", order.OrderID))
		} else {
			log.Info("unrecognizable error", slog.String("номер заказа:", order.OrderID))
			return err
		}
	}
	if useOrder.OrderID != "" {
		log.Info("order has already been created")
		return models.ErrOrderExists
	}

	//получаем текущий баланс
	row = db.DB.QueryRow(ctx, "SELECT balance FROM users WHERE id=$1", order.UserID)
	var balance float32
	err = row.Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("no balance for user", order.UserID)
			return errors.New("no balance for user")
		}
		log.Error("error scanning row", slog.String("error", err.Error()))
		return err
	}

	//проверяем достаточность средств и делаем запись
	switch balance >= order.Withdraw {

	case true: //если средств достаточно
		tx, err := db.DB.Begin(ctx) //открываем транзакцию изменения таблиц
		if err != nil {
			log.Error("unable to begin transaction")
			return fmt.Errorf("%w", err)
		}
		query := `INSERT INTO orders (user_id, order_num, order_status, created_at, withdraw)
			VALUES ($1, $2, $3, $4, $5);`
		_, err = tx.Exec(ctx, query,
			order.UserID,
			order.OrderID,
			order.Status,
			order.CreatedAt,
			order.Withdraw)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Error("unable to rollback transaction to create withdraw")
				return fmt.Errorf("%w", err)
			}
			log.Error("unable to create order withdraw", err)
			return fmt.Errorf("unable to create order withdraw: %w", err)
		}
		queryBalance := "UPDATE users SET balance=balance-$1 WHERE id=$2;"
		_, err = tx.Exec(ctx, queryBalance, order.Withdraw, order.UserID)
		if err != nil {
			err = tx.Rollback(ctx)
			if err != nil {
				log.Error("unable to rollback transaction to update balance")
				return fmt.Errorf("%w", err)
			}
			log.Error("unable to update balance")
			return fmt.Errorf("unable to update balance: %w", err)
		}

		err = tx.Commit(ctx) //закрываем транзакцию
		if err != nil {
			log.Error("unable to commit transaction to create order withdraw")
			return fmt.Errorf("%w", err)
		}

	case false:
		return models.ErrInsufficientFunds //возвращаем ошибку о недостаточности средств на счете
	}

	log.Info("order withdraw has been created successfully")
	return nil
}
