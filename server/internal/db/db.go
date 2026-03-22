package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
}

type Service struct {
	pool *pgxpool.Pool
}

func Must(config Config) *Service {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host,
		config.Port,
		config.Username,
		config.Password,
		config.Name,
	)

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	return &Service{pool: pool}
}

func (s *Service) Handle() *pgxpool.Pool {
	return s.pool
}

func (s *Service) Close() error {
	s.pool.Close()
	return nil
}
