package handler

import (
	"fmt"
	"net/http"

	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/views"
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

func HandleLoginPage(w http.ResponseWriter, r *http.Request) {
	if err := views.Login().Render(r.Context(), w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func HandleRegister(pool *pgxpool.Pool) http.HandlerFunc {
	// usersStore := repo.NewPGXUsersRepository(pool)
	// authService := auth.New(usersStore)

	return func(w http.ResponseWriter, r *http.Request) {
		username := r.FormValue("username")
		email := r.FormValue("email")
		password := r.FormValue("password")

		req := &auth.RegisterRequest{
			Username: username,
			Email:    email,
			Password: password,
		}

		fmt.Println(req)

		// problems, err := validate(req)
		// if err != nil {

		// }

		// _, err = authService.Register(r.Context(), req)
	}
}
