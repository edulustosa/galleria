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

	r.Get("/galleria", handlers.HandleGalleria(pool))
	r.Get("/galleria/posts/{postId}/comments", handlers.HandlePostComments(pool))

	r.Group(func(r chi.Router) {
		r.Use(middlewares.JWTAuthMiddleware([]byte(jwtKey)))

		r.Get("/profile", handlers.HandleGetUserProfile(pool))
		r.Get("/profile/images", handlers.HandleGetUserImages(pool))
		r.Patch("/profile", handlers.HandleUpdateProfile(pool))

		r.Post("/galleria/posts/{postId}", handlers.HandleAddComment(pool))
	})
}
