package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strconv"
)

func (db *PostgresqlDB) GetUserBalance(ctx context.Context, userID int) (current, withdrawn float32, err error) {
	const op = "postgresql.GetUserBalance"
	log := db.logger.With(
		slog.String("op", op),
		slog.String("user", strconv.Itoa(userID)),
	)

	//получаем текущий баланс
	row := db.DB.QueryRow(ctx, "SELECT balance FROM users WHERE id=$1;", userID)

	err = row.Scan(&current)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("no balance for user", userID)
			return 0, 0, err
		}
		log.Error("error scanning row", slog.String("error", err.Error()))
		return 0, 0, err
	}

	//считаем списания
	row = db.DB.QueryRow(ctx, `SELECT sum(withdraw) FROM "orders" WHERE user_id=$1;`, userID)
	err = row.Scan(&withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("no withdrawn for user", userID)
			return 0, 0, err
		}
		log.Error("error scanning row", slog.String("error", err.Error()))
		return 0, 0, err
	}

	return current, withdrawn, nil
}
