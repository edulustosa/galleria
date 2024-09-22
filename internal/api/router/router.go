package router

import (
	"net/http"

	"github.com/edulustosa/galleria/internal/api/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	addRoutes(r, pool)

	return r
}

func addRoutes(r chi.Router, pool *pgxpool.Pool) {
	r.Post("/register", handlers.HandleRegister(pool))
	r.Post("/login", handlers.HandleLogin(pool))
}
