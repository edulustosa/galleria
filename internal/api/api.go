package api

import (
	"net/http"

	"github.com/edulustosa/galleria/internal/api/handler"
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

	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	addRoutes(r, pool)

	return r
}

func addRoutes(r chi.Router, pool *pgxpool.Pool) {
	r.Get("/register", handler.HandleRegisterPage)
	r.Post("/register", handler.HandleRegister(pool))

	r.Get("/login", handler.HandleLoginPage)
}
