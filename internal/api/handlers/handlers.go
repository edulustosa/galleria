package handlers

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/edulustosa/galleria/internal/api"
	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/factories"
	"github.com/edulustosa/galleria/internal/profile"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthResponse struct {
	UserID string `json:"userId"`
}

type GenericResponse map[string]interface{}

func HandleRegister(pool *pgxpool.Pool) http.HandlerFunc {
	usersRepository := repo.NewPGXUsersRepository(pool)
	authService := auth.New(usersRepository)

	return func(w http.ResponseWriter, r *http.Request) {
		req, problems, err := api.DecodeValid[auth.RegisterRequest](r)
		if err != nil {
			api.HandleInvalidRequest(w, problems)
			return
		}

		userId, err := authService.Register(r.Context(), &req)
		if err != nil {
			if errors.Is(err, auth.ErrUserAlreadyExists) {
				api.HandleError(w, http.StatusConflict, api.Error{Message: err.Error()})
				return
			}

			log.Printf("failed to register user: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "something went wrong, please try again"},
			)
			return
		}

		resp := AuthResponse{UserID: userId.String()}
		if err = api.Encode(w, http.StatusCreated, resp); err != nil {
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
		req, problems, err := api.DecodeValid[auth.LoginRequest](r)
		if err != nil {
			api.HandleInvalidRequest(w, problems)
			return
		}

		userId, err := authService.Login(r.Context(), &req)
		if err != nil {
			// The only error that can be returned is ErrInvalidCredentials
			api.HandleError(w, http.StatusUnauthorized, api.Error{Message: err.Error()})
			return
		}

		token, err := createJWT(userId.String(), jwtKey)
		if err != nil {
			log.Printf("failed to create token: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "failed to create token"},
			)
			return
		}

		resp := LoginResponse{
			UserID: userId.String(),
			Token:  token,
		}
		if err = api.Encode(w, http.StatusOK, resp); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleGetUserProfile(pool *pgxpool.Pool) http.HandlerFunc {
	profile := factories.MakeProfileService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.UserIDKey).(uuid.UUID)
		user, err := profile.GetProfile(r.Context(), userId)
		if err != nil {
			// The only error that can be returned is InvalidCredentials
			api.HandleError(w, http.StatusUnauthorized, api.Error{Message: err.Error()})
			return
		}

		if err = api.Encode(w, http.StatusOK, user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleGetUserImages(pool *pgxpool.Pool) http.HandlerFunc {
	profile := factories.MakeProfileService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.UserIDKey).(uuid.UUID)
		images, err := profile.GetProfileImages(r.Context(), userId)
		if err != nil {
			log.Printf("failed to get images: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "failed to get images"},
			)
			return
		}

		if err = api.Encode(w, http.StatusOK, GenericResponse{"images": images}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandleUpdateProfile(pool *pgxpool.Pool) http.HandlerFunc {
	profileService := factories.MakeProfileService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		userId := r.Context().Value(api.UserIDKey).(uuid.UUID)
		req, problems, err := api.DecodeValid[profile.UpdateProfileRequest](r)
		if err != nil {
			api.HandleInvalidRequest(w, problems)
			return
		}

		err = profileService.Update(r.Context(), userId, &req)
		if err != nil {
			log.Printf("failed to update profile: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "failed to update profile"},
			)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
