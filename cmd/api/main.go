package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/zippy/config"
	"github.com/rosset7i/zippy/internal/infra/database"
	"github.com/rosset7i/zippy/internal/infra/webserver/handlers"
)

func main() {
	config := config.LoadConfig()

	connectionString := fmt.Sprintf(
		"dbname=%v user=%v password=%v host=%v port=%v sslmode=disable client_encoding=UTF8",
		config.DBName,
		config.DBUser,
		config.DBPassword,
		config.DBHost,
		config.DBPort,
	)

	db, err := sql.Open(config.DBDriver, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	userHandler := handlers.NewUserHandler(database.NewUser(db))
	productHandler := handlers.NewProductHandler(database.NewProduct(db))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Route("/users", func(r chi.Router) {
		r.Get("/", userHandler.FetchByEmail)
		r.Post("/", userHandler.Create)
	})
	r.Route("/products", func(r chi.Router) {
		r.Get("/", productHandler.FetchPaged)
		r.Get("/fetch-by-id", productHandler.FetchById)
		r.Post("/", productHandler.Create)
		r.Put("/", productHandler.Update)
		r.Delete("/", productHandler.Delete)
	})
	if err := http.ListenAndServe(config.WebServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
