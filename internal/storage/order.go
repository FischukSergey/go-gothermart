package storage

import (
	"context"
	"errors"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/jackc/pgx/v5"
	"log/slog"
)

// CreateOrder() добавление новой записи заказа (ордера)
func (db *PostgresqlDB) CreateOrder(ctx context.Context, order models.Order) error {
	const op = "postgresql.Register"
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

	//если заказ уже существует проверяем кто его записал
	if useOrder.OrderID != "" {
		if useOrder.UserID == order.UserID {
			log.Info("order already loaded same user", slog.String("номер заказа:", order.OrderID))
			return models.ErrOrderUploadedSameUser
		}
		log.Info("order already loaded another user", slog.String("номер заказа:", order.OrderID))
		return models.ErrOrderUploadedAnotherUser
	}

	//так как заказа не существует, делаем запись
	row = db.DB.QueryRow(ctx,
		`INSERT INTO orders (user_id, order_num, order_status, created_at)
		VALUES ($1, $2, $3, $4) RETURNING id;`,
		order.UserID,
		order.OrderID,
		order.Status,
		order.CreatedAt)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Error("failed to create order", logger.Err(err))
		return err
	}

	return nil
}
