package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/product_crud/config"
)

func New(ctx context.Context, c *config.ConfDB) *pgxpool.Pool {
	config, err := pgxpool.ParseConfig(buildConnectionString(c))
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		log.Fatal(err)
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		pool.Close()
		log.Fatal(err)
	}

	return pool
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
