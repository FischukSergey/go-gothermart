package storage

import (
	"context"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	_ "github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PostgresqlDB struct {
	DB     *pgxpool.Pool
	logger *slog.Logger
}

// NewDB инициализация пула соединений с базой postgresql
func NewDB(dbConfig *pgconn.Config, log *slog.Logger) (*PostgresqlDB, error) {

	connect := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
		dbConfig.User, dbConfig.Password, dbConfig.Host, strconv.Itoa(int(dbConfig.Port)), dbConfig.Database)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	db, err := pgxpool.New(ctx, connect) //config.FlagDatabaseDSN)
	if err != nil {
		return nil, fmt.Errorf("%w, unable to create connection db:%s", err, dbConfig.Database)
	}
	return &PostgresqlDB{
		DB:     db,
		logger: log,
	}, nil
}

func (db *PostgresqlDB) Register(ctx context.Context, u *models.User) (int, error) {
	const op = "postgresql.Register"
	log := db.logger.With(
		slog.String("op", op),
		slog.String("email", u.Email),
	)

	//готовим запрос на вставку
	query := `INSERT INTO users (email, password, role) VALUES($1,$2,$3);`
	_, err := db.DB.Exec(ctx, query, u.Email, u.EncryptedPassword, "customer") //TODO вставить роль пользователя
	//обработка ошибки сохранения нового пользователя
	if err != nil {
		//если login неуникальный
		if strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
			//БД возвращает ошибку на русском языке. Из-за этого не обрабатывается ошибка. Как исправить не нашел.
			//var pgErr *pgconn.PgError
			//if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("%s: %w", u.Email, models.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: не удалось выполнить запись в базу %w", op, err)
	}
	log.Info("Success create user", "email", u.Email)
	return 0, nil
}
