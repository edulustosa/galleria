package handler

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/profile"
	"github.com/edulustosa/galleria/internal/views/components"
	"github.com/edulustosa/galleria/internal/views/pages"
	"github.com/google/uuid"
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
	if err := pages.Register().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func renderErrAlert(w http.ResponseWriter, r *http.Request, errMsg string) {
	if err := components.ErrAlert(errMsg).Render(r.Context(), w); err != nil {
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
				renderErrAlert(w, r, problem)
				return
			}
		}

		_, err = authService.Register(r.Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				renderErrAlert(w, r, err.Error())
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
	if err := pages.Login().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleLogin(pool *pgxpool.Pool, store *sessions.CookieStore) http.HandlerFunc {
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
				renderErrAlert(w, r, problem)
				return
			}
		}

		userID, err := authService.Login(r.Context(), req)
		if err != nil {
			if errors.Is(err, auth.ErrInvalidCredentials) {
				renderErrAlert(w, r, err.Error())
				return
			}

			log.Printf("failed to login user: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, _ := store.Get(r, "session")
		opts := &sessions.Options{
			MaxAge:   86400 * 7,
			SameSite: http.SameSiteDefaultMode,
			Path:     "/",
			HttpOnly: true,
		}
		session.Options = opts

		session.Values["authenticated"] = true
		session.Values["user_id"] = userID.String()

		if err := session.Save(r, w); err != nil {
			log.Printf("failed to save session: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func HandleLogout(store *sessions.CookieStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = false
		session.Values["user_id"] = ""
		session.Options.MaxAge = -1

		if err := session.Save(r, w); err != nil {
			log.Printf("failed to save session: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("HX-Redirect", "/")
		w.WriteHeader(http.StatusOK)
	}
}

func getUserIdFromSession(r *http.Request, store *sessions.CookieStore) (uuid.UUID, error) {
	session, err := store.Get(r, "session")
	if err != nil {
		return uuid.Nil, err
	}

	userIdStr, ok := session.Values["user_id"].(string)
	if !ok {
		return uuid.Nil, errors.New("failed to get user_id from session")
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func HandleProfilePage(pool *pgxpool.Pool, store *sessions.CookieStore) http.HandlerFunc {
	usersRepository := repo.NewPGXUsersRepository(pool)
	imagesRepository := repo.NewPGXImagesRepo(pool)
	profileService := profile.New(usersRepository, imagesRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		userId, err := getUserIdFromSession(r, store)
		if err != nil {
			log.Printf("failed to get user_id from session: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, err := profileService.GetProfile(r.Context(), userId)
		if err != nil {
			w.Header().Set("HX-Redirect", "/login")
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := pages.Profile(user).Render(r.Context(), w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
