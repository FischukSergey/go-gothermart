package storage

import (
	"context"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"log/slog"
)

func (db *PostgresqlDB) GetAccrualOrders(ctx context.Context) ([]string, error) {
	const op = "storage.GetAccrualOrders"
	log := db.logger.With(
		slog.String("op", op),
	)
	var orders []string

	query := `SELECT order_num FROM orders WHERE order_status = $1;`

	rows, err := db.DB.Query(ctx, query, models.StatusOrderNew)
	if err != nil {
		log.Error("unable to execute query")
		return orders, fmt.Errorf("unable to execute query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order string
		err = rows.Scan(&order)
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

func (db *PostgresqlDB) UpdateAccrualOrder(ctx context.Context, order models.AccrualOrder) error {
	const op = "storage.UpdateAccrualOrder"
	log := db.logger.With(
		slog.String("op", op))

	query := "UPDATE orders SET order_status=$1, accrual=$2 WHERE order_num=$3;"
	_, err := db.DB.Exec(ctx, query, order.Status, order.Accrual, order.Order)
	if err != nil {
		log.Error("unable to execute query")
		return fmt.Errorf("unable to execute query: %w", err)
	}
	return nil
}
