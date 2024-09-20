package api

import (
	"net/http"

	"github.com/edulustosa/galleria/internal/api/handler"
	"github.com/edulustosa/galleria/internal/api/middlewares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServer(pool *pgxpool.Pool, store *sessions.CookieStore) http.Handler {
	r := chi.NewMux()

	r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		middleware.Recoverer,
	)

	fs := http.FileServer(http.Dir("assets"))
	r.Handle("/assets/*", http.StripPrefix("/assets/", fs))

	addRoutes(r, pool, store)

	return r
}

func addRoutes(r chi.Router, pool *pgxpool.Pool, store *sessions.CookieStore) {
	r.Get("/register", handler.HandleRegisterPage)
	r.Post("/register", handler.HandleRegister(pool))

	r.Get("/login", handler.HandleLoginPage)
	r.Post("/login", handler.HandleLogin(pool, store))
	r.Get("/logout", handler.HandleLogout(store))

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middlewares.AuthMiddleware(store))

		r.Get("/profile", handler.HandleProfilePage(pool, store))
		// r.Put("/profile", handler.HandleUpdateProfile(pool))
	})
}
