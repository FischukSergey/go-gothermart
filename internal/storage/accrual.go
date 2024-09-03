package storage

import (
	"context"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/app/services"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// GetAccrualOrders выдает все еще не обработанные заказы
func (db *PostgresqlDB) GetAccrualOrders(ctx context.Context) ([]services.OrderUpdate, error) {
	const op = "storage.GetAccrualOrders"
	log := db.logger.With(
		slog.String("op", op),
	)
	var orders []services.OrderUpdate

	query := `SELECT order_num, user_id FROM orders WHERE order_status = $1;`

	rows, err := db.DB.Query(ctx, query, models.StatusOrderNew)
	if err != nil {
		log.Error("unable to execute query")
		return orders, fmt.Errorf("unable to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order services.OrderUpdate
		err = rows.Scan(&order.OrderID, &order.UserID)
		if err != nil {
			log.Error("unable to read row of query")
			return orders, fmt.Errorf("unable to read row of query: %w", err)
		}
		orders = append(orders, order)
	}
	if err := rows.Err(); err != nil {
		return orders, fmt.Errorf("scan query error: %w", err)
	}

	log.Info("NEW orders selected successfully")
	return orders, nil
}

// UpdateAccrualOrder обновляет данные о статусе заказа в соответствии с данными, полученными
// от микросервиса начисления баллов лояльности и пересчитывает баланс баллов пользователя id
func (db *PostgresqlDB) UpdateAccrualOrder(ctx context.Context, order models.AccrualOrder, userID int) error {
	const op = "storage.UpdateAccrualOrder"
	log := db.logger.With(
		slog.String("op", op))

	tx, err := db.DB.Begin(ctx) //открываем транзакцию изменения таблиц
	if err != nil {
		log.Error("unable to begin transaction")
		return fmt.Errorf("%w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		_ = tx.Rollback(ctx)
	}(tx, ctx)

	//обновляем статус заказа
	query := "UPDATE orders SET order_status=$1, accrual=$2 WHERE order_num=$3;"
	_, err = tx.Exec(ctx, query, order.Status, order.Accrual, order.Order)
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error("unable to rollback transaction to update accrual orders")
			return fmt.Errorf("%w", err)
		}
		log.Error("unable to execute query")
		return fmt.Errorf("unable to execute query: %w", err)
	}

	//обновляем баланс
	queryBalance := "UPDATE users SET balance=balance+$1 WHERE id=$2;"
	_, err = tx.Exec(ctx, queryBalance, order.Accrual, userID) //todo id user
	if err != nil {
		err = tx.Rollback(ctx)
		if err != nil {
			log.Error("unable to rollback transaction to update balance")
			return fmt.Errorf("%w", err)
		}
		log.Error("unable to execute query")
		return fmt.Errorf("unable to execute query: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		log.Error("unable to commit transaction to update accrual orders")
		return fmt.Errorf("%w", err)
	}

	return nil
}
