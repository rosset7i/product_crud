package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/product_crud/config"
)

func New(ctx context.Context, c *config.ConfDB) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(buildConnectionString(c))
	if err != nil {
		return nil, err
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		return nil, err
	}

	return pool, nil
}

func buildConnectionString(c *config.ConfDB) string {
	return fmt.Sprintf(
		"dbname=%v user=%v password=%v host=%v port=%v sslmode=disable client_encoding=UTF8",
		c.DBName,
		c.Username,
		c.Password,
		c.Host,
		c.Port,
	)
}
