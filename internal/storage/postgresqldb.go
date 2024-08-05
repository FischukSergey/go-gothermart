package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/FischukSergey/go-gothermart.git/internal/logger"
	"github.com/FischukSergey/go-gothermart.git/internal/models"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
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
	query := `
	CREATE TABLE IF NOT EXISTS users
	  (id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    email VARCHAR NOT NULL UNIQUE,
    password VARCHAR NOT NULL,
    role VARCHAR,
    time TIME);
	`
	_, err = db.Exec(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%w, unable to execute query", err)
	}
	_, err = db.Exec(ctx, "CREATE UNIQUE INDEX IF NOT EXISTS email_idx ON users (email);") //создаем уникальный индекс по оригинальному url
	if err != nil {
		return nil, fmt.Errorf("%w, unable to create index", err)
	}
	return &PostgresqlDB{
		DB:     db,
		logger: log,
	}, nil
}

// Register() метод принимает логин и пароль, проверяет на уникальность логин,
// сохранят в таблице users, и возвращает ошибку.
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

	//TODO вставить извлечение ID select....
	return 0, nil
}

// Login() метод принимает логин и пароль, проверяет на наличие и возвращает ошибку.
func (db *PostgresqlDB) Login(ctx context.Context, email string) (*models.User, error) {
	const op = "postgresql.Register"
	log := db.logger.With(
		slog.String("op", op),
		slog.String("email", email),
	)

	query := `SELECT email, password, role, id FROM users WHERE email=$1;`
	user := models.User{}
	err := db.DB.QueryRow(ctx, query, email).Scan(&user.Email, &user.EncryptedPassword, &user.Role, &user.ID)
	if errors.Is(err, pgx.ErrNoRows) {
		log.Error("row not found", "email", email, logger.Err(err))
		return nil, err
	}
	if err != nil {
		log.Error("unable to execute query", logger.Err(err))
		return nil, err
	}

	return &user, nil
}
