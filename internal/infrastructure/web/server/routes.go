package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth"
	"github.com/rosset7i/product_crud/config"
	"github.com/rosset7i/product_crud/internal/infrastructure/database"
	"github.com/rosset7i/product_crud/internal/infrastructure/web"
	httpSwagger "github.com/swaggo/http-swagger"
)

func (s *Server) MapHandlers(c *config.Conf) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)

	r.Route("/v1", func(r chi.Router) {
		userHandler := web.NewUserHandler(database.NewUserRepository(s.db), c)
		r.Route("/users", func(r chi.Router) {
			r.Post("/register", userHandler.Register)
			r.Post("/login", userHandler.Login)
		})

		productHandler := web.NewProductHandler(database.NewProductRepository(s.db))
		r.Route("/products", func(r chi.Router) {
			r.Use(jwtauth.Verifier(c.Auth.JwtAuth))
			r.Use(jwtauth.Authenticator)
			r.Get("/", productHandler.FetchPaged)
			r.Get("/{id}", productHandler.FetchById)
			r.Post("/", productHandler.Create)
			r.Put("/", productHandler.Update)
			r.Delete("/", productHandler.Delete)
		})
	})
	r.Get("/docs/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:7000/docs/doc.json"),
	))

	return r
}
