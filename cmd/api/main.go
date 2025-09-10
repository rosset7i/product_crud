package main

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/product_crud/config"
	_ "github.com/rosset7i/product_crud/docs"
	"github.com/rosset7i/product_crud/internal/infrastructure/web/server"
)

// @title           product_crud API
// @version         1.0

// @securityDefinitions.apiKey  Bearer
// @in                        header
// @name                      Authorization
func main() {
	s := server.NewServer(config.New())
	s.Run()
}
