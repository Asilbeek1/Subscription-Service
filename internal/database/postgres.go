package database

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Asilbeek1/Subscription-Service/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
	log  *slog.Logger
}

func OpenDB(ctx context.Context, cfg config.Postgres, log *slog.Logger) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.DB_USER,
		cfg.DB_PASSWORD,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.SSLMode,
	)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	poolConfig.MaxConns = int32(cfg.MaxOpenConns)
	poolConfig.MaxConnLifetime = cfg.ConnLifeTime

	db, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &DB{pool: db, log: log}, nil
}
