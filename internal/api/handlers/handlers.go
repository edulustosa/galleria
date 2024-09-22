package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func encode[T any](w http.ResponseWriter, status int, data T) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}

	return nil
}

type Validator interface {
	Valid() (problems map[string]string)
}

func decodeValid[T Validator](r *http.Request) (T, map[string]string, error) {
	var v T
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		return v, nil, fmt.Errorf("decode json: %w", err)
	}

	if problems := v.Valid(); len(problems) > 0 {
		return v, problems, fmt.Errorf("invalid %T: %d problems", v, len(problems))
	}

	return v, nil, nil
}

type ErrorList struct {
	Errors []Error `json:"errors"`
}

type Error struct {
	Message string `json:"message"`
	Details string `json:"details"`
}

func handleInvalidRequest(w http.ResponseWriter, problems map[string]string) {
	var errors []Error

	if len(problems) > 0 {
		errors = make([]Error, 0, len(problems))
		for field, problem := range problems {
			err := Error{
				Message: fmt.Sprintf("invalid %s", field),
				Details: problem,
			}
			errors = append(errors, err)
		}
	} else {
		errors = make([]Error, 0, 1)
		errors = append(errors, Error{Message: "invalid input"})
	}

	HandleError(w, http.StatusBadRequest, errors...)
}

func HandleError(w http.ResponseWriter, status int, err ...Error) {
	errList := ErrorList{Errors: err}
	if err := encode(w, status, errList); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type AuthResponse struct {
	UserID string `json:"userId"`
}

func HandleRegister(pool *pgxpool.Pool) http.HandlerFunc {
	usersRepository := repo.NewPGXUsersRepository(pool)
	authService := auth.New(usersRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[auth.RegisterRequest](r)
		if err != nil {
			handleInvalidRequest(w, problems)
			return
		}

		userId, err := authService.Register(r.Context(), &req)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				HandleError(w, http.StatusConflict, Error{Message: err.Error()})
				return
			}

			log.Printf("failed to register user: %v", err)
			HandleError(
				w,
				http.StatusInternalServerError,
				Error{Message: "something went wrong, please try again"},
			)
			return
		}

		resp := AuthResponse{UserID: userId.String()}
		if err = encode(w, http.StatusCreated, resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func createJWT(userId, jwtKey string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"sub": userId,
			"exp": time.Now().Add(72 * time.Hour).Unix(),
			"iat": time.Now().Unix(),
			"nbf": time.Now().Unix(),
		},
	)

	return token.SignedString([]byte(jwtKey))
}

type LoginResponse struct {
	UserID string `json:"userId"`
	Token  string `json:"token"`
}

func HandleLogin(pool *pgxpool.Pool, jwtKey string) http.HandlerFunc {
	usersRepository := repo.NewPGXUsersRepository(pool)
	authService := auth.New(usersRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := decodeValid[auth.LoginRequest](r)
		if err != nil {
			handleInvalidRequest(w, problems)
			return
		}

		userId, err := authService.Login(r.Context(), &req)
		if err != nil {
			// The only error that can be returned is ErrInvalidCredentials
			HandleError(w, http.StatusUnauthorized, Error{Message: err.Error()})
			return
		}

		token, err := createJWT(userId.String(), jwtKey)
		if err != nil {
			log.Printf("failed to create token: %v", err)
			HandleError(
				w,
				http.StatusInternalServerError,
				Error{Message: "failed to create token"},
			)
			return
		}

		resp := LoginResponse{
			UserID: userId.String(),
			Token:  token,
		}
		if err = encode(w, http.StatusOK, resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
