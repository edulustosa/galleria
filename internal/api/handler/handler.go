package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/views"
	"github.com/edulustosa/galleria/internal/views/components"
	"github.com/gorilla/sessions"
	"github.com/jackc/pgx/v5/pgxpool"
)

type validator interface {
	Valid() (problems map[string]string)
}

func validate(v validator) (problems map[string]string, err error) {
	if problems = v.Valid(); len(problems) > 0 {
		return problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return nil, nil
}

func HandleRegisterPage(w http.ResponseWriter, r *http.Request) {
	if err := views.Register().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRegister(pool *pgxpool.Pool) http.HandlerFunc {
	usersStore := repo.NewPGXUsersRepository(pool)
	authService := auth.New(usersStore)

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		req := &auth.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		problems, err := validate(req)
		if err != nil {
			for _, problem := range problems {
				err := components.ErrAlert(problem).Render(r.Context(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		_, err = authService.Register(r.Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				err := components.ErrAlert(err.Error()).Render(r.Context(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			log.Printf("failed to register user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusCreated)
	}
}

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := views.Login().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLogin(
	pool *pgxpool.Pool,
	store *sessions.CookieStore,
) http.HandlerFunc {
	usersStore := repo.NewPGXUsersRepository(pool)
	authService := auth.New(usersStore)

	return func(w http.ResponseWriter, r *http.Request) {
		email := r.FormValue("email")
		password := r.FormValue("password")

		req := &auth.LoginRequest{
			Email:    email,
			Password: password,
		}

		problems, err := validate(req)
		if err != nil {
			for _, problem := range problems {
				err := components.ErrAlert(problem).Render(r.Context(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}

		userID, err := authService.Login(r.Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				err := components.ErrAlert(err.Error()).Render(r.Context(), w)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}

			log.Printf("failed to login user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := store.New(r, "session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session.Values["userId"] = userID.ID
		if err := session.Save(r, w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}
