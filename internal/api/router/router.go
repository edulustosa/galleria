package router

import (
	"net/http"

	"github.com/edulustosa/galleria/internal/api/handlers"
	"github.com/edulustosa/galleria/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool, jwtKey string) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	addRoutes(r, pool, jwtKey)

	return r
}

func addRoutes(r chi.Router, pool *pgxpool.Pool, jwtKey string) {
	r.Post("/register", handlers.HandleRegister(pool))
	r.Post("/login", handlers.HandleLogin(pool, jwtKey))

	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware([]byte(jwtKey)))

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("pong"))
		})
	})
}
