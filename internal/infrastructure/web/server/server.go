package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rosset7i/product_crud/config"
	"github.com/rosset7i/product_crud/internal/infrastructure/database"
)

type Server struct {
	c         *config.Conf
	db        *pgxpool.Pool
	container *Container
}

func NewServer(c *config.Conf) *Server {
	return &Server{
		c:  c,
		db: database.New(context.Background(), &c.DB),
	}
}

func (s *Server) Run() {
	s.init()

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.c.Server.Port),
		Handler:      s.MapHandlers(s.c),
		ReadTimeout:  s.c.Server.TimeoutRead,
		WriteTimeout: s.c.Server.TimeoutWrite,
		IdleTimeout:  s.c.Server.TimeoutIdle,
	}

	go func() {
		log.Printf("Starting server at :%d", s.c.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown

	log.Println("shutting down gracefully, press Ctrl+C again to force")

	s.db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Graceful shutdown complete.")
}
