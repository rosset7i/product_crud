package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/rosset7i/product_crud/config"
	_ "github.com/rosset7i/product_crud/docs"
	"github.com/rosset7i/product_crud/internal/infra/database"
	"github.com/rosset7i/product_crud/internal/infra/webserver/handlers"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title           product_crud API
// @version         1.0
// @description     product_crud is an API for managing users and products.
// @description     It provides authentication endpoints and a product catalog with CRUD operations.
//
// @termsOfService  http://swagger.io/terms/
//
// @contact.name    product_crud API Support
// @contact.url     http://www.swagger.io/support
// @contact.email   support@swagger.io
//
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
//
// @host            localhost:7000
// @BasePath        /
//
// @securityDefinitions.apiKey  Bearer
// @in                        header
// @name                      Authorization
func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPool(context.Background(), cfg.ConnectionString)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	server := &http.Server{
		Addr:    cfg.WebServerAddress,
		Handler: router(cfg, db),
	}

	shutdown := make(chan bool)

	go gracefulShutdown(server, shutdown)

	log.Printf("Starting server at %s", cfg.WebServerAddress)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-shutdown
	log.Printf("Server exited.")
}

func gracefulShutdown(server *http.Server, shutdown chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	<-ctx.Done()
	stop()

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Graceful shutdown complete.")
	shutdown <- true
}

func router(cfg *config.Config, db *pgxpool.Pool) http.Handler {
	userHandler := handlers.NewUserHandler(database.NewUserRepository(db), cfg)
	productHandler := handlers.NewProductHandler(database.NewProductRepository(db))

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/users", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)
	})

	r.Route("/products", func(r chi.Router) {
		r.Use(jwtauth.Verifier(cfg.TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Get("/", productHandler.FetchPaged)
		r.Get("/{id}", productHandler.FetchById)
		r.Post("/", productHandler.Create)
		r.Put("/", productHandler.Update)
		r.Delete("/", productHandler.Delete)
	})

	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7000/docs/doc.json"),
	))

	return r
}
