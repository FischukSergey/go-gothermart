package storage

import (
	"context"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"log/slog"
	"strconv"
)

func (db *PostgresqlDB) GetUserOrders(ctx context.Context, id int) ([]models.GetUserOrders, error) {
	const op = "storage.GetUserOrders"
	log := db.logger.With(
		slog.String("op", op),
		slog.String("user", strconv.Itoa(id)),
	)

	var orders []models.GetUserOrders

	query := `SELECT order_num, order_status, accrual, created_at FROM orders
		WHERE user_id = $1 ORDER BY created_at DESC;` //todo select another

	rows, err := db.DB.Query(ctx, query, id)
	if err != nil {
		log.Error("unable to execute query")
		return orders, fmt.Errorf("unable to execute query: %w", err)
	}
	if err = rows.Err(); err != nil {
		log.Error("rows error")
		return orders, fmt.Errorf("rows error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var row models.GetUserOrders
		err = rows.Scan(&row.Number, &row.Status, &row.Accrual, &row.UploadedAt)
		if err != nil {
			log.Error("unable to read row of query")
			return orders, fmt.Errorf("unable to read row of query: %w", err)
		}
		orders = append(orders, row)
	}
	if err := rows.Err(); err != nil {
		return orders, fmt.Errorf("scan query error: %w", err)
	}

	log.Info("selected orders successfully", orders)
	return orders, nil
}
