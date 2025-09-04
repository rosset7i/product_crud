package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/zippy/config"
	_ "github.com/rosset7i/zippy/docs"
	"github.com/rosset7i/zippy/internal/infra/database"
	"github.com/rosset7i/zippy/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Zippy API
// @version 1.0
// @description IDK.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:7000
// @BasePath /
// @securityDefinitions.apiKey ApiKeyAuth
// @in header
// @name Authorization
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

	userHandler := handlers.NewUserHandler(database.NewUser(db), config)
	productHandler := handlers.NewProductHandler(database.NewProduct(db))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Route("/users", func(r chi.Router) {
		r.Post("/", userHandler.Create)
		r.Post("/", userHandler.Login)
	})
	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(config.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.FetchPaged)
		r.Get("/{id}", productHandler.FetchById)
		r.Post("/", productHandler.Create)
		r.Put("/", productHandler.Update)
		r.Delete("/", productHandler.Delete)
	})
	r.Get("/docs/*", httpSwagger.Handler(httpSwagger.URL("http://localhost:7000/docs/doc.json")))
	if err := http.ListenAndServe(config.WebServerAddress, r); err != nil {
		log.Fatal(err)
	}
}
