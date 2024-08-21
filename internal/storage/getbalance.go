package storage

import (
	"context"
	"errors"
	"fmt"
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
	tx, err := db.DB.Begin(ctx) //открываем транзакцию для синхронизации данных из таблиц
	if err != nil {
		log.Error("unable to begin transaction")
		return 0, 0, fmt.Errorf("%w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Error("unable to rollback transaction to get current balance")
		}
	}(tx, ctx)
	//получаем текущий баланс
	row := db.DB.QueryRow(ctx, "SELECT balance FROM users WHERE id=$1 FOR UPDATE;", userID)

	err = row.Scan(&current)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("no balance for user", slog.String("userID", strconv.Itoa(userID)))
			return 0, 0, err
		}
		log.Error("error scanning row", slog.String("error", err.Error()))
		return 0, 0, err
	}

	//считаем списания
	row = db.DB.QueryRow(ctx, `SELECT sum(withdraw) FROM "orders" WHERE user_id=$1 FOR UPDATE;`, userID)
	err = row.Scan(&withdrawn)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("no withdrawn for user", slog.String("userID", strconv.Itoa(userID)))
			return 0, 0, err
		}
		log.Error("error scanning row", slog.String("error", err.Error()))
		return 0, 0, err
	}
	err = tx.Commit(ctx) //закрываем транзакцию
	if err != nil {
		log.Error("unable to commit transaction to get current balance")
		return 0, 0, fmt.Errorf("%w", err)
	}
	return current, withdrawn, nil
}
