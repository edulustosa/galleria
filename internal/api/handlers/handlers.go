package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/edulustosa/galleria/internal/api"
	"github.com/edulustosa/galleria/internal/auth"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/factories"
	"github.com/edulustosa/galleria/internal/galleria"
	"github.com/edulustosa/galleria/internal/profile"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthResponse struct {
	UserID string `json:"userId"`
}

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

		if err = api.Encode(w, http.StatusOK, api.JSON{"images": images}); err != nil {
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

func HandleGalleria(pool *pgxpool.Pool) http.HandlerFunc {
	galleria := factories.MakeGalleriaService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		var page uint64 = 1
		pageStr := r.URL.Query().Get("page")
		if pageStr != "" {
			p, err := strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				api.HandleError(w, http.StatusBadRequest, api.Error{
					Message: "invalid page",
					Details: "page must be a positive integer",
				})
				return
			}
			page = p
		}

		posts, err := galleria.Display(r.Context(), page)
		if err != nil {
			log.Printf("failed to get images: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "something went wrong, please try again"},
			)
			return
		}

		if err = api.Encode(w, http.StatusOK, api.JSON{"posts": posts}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

type AddCommentRequest struct {
	Comment string `json:"comment"`
}

func (r AddCommentRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if len(r.Comment) > 500 {
		problems["comment"] = "comment must be less than 500 characters"
	}

	return
}

func HandleAddComment(pool *pgxpool.Pool) http.HandlerFunc {
	galleriaService := factories.MakeGalleriaService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(api.UserIDKey).(uuid.UUID)
		postId, err := uuid.Parse(chi.URLParam(r, "postId"))
		if err != nil {
			api.HandleError(w, http.StatusBadRequest, api.Error{
				Message: "invalid post id",
				Details: "post id must be a valid UUID",
			})
			return
		}

		req, problems, err := api.DecodeValid[AddCommentRequest](r)
		if err != nil {
			api.HandleInvalidRequest(w, problems)
			return
		}

		commentID, err := galleriaService.AddComment(r.Context(), userID, postId, req.Comment)
		if err != nil {
			if errors.Is(err, galleria.ErrImageNotFound) || errors.Is(err, galleria.ErrUserNotFound) {
				api.HandleError(w, http.StatusNotFound, api.Error{Message: err.Error()})
				return
			}

			log.Printf("failed to add comment: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "something went wrong, please try again"},
			)
			return
		}

		if err := api.Encode(w, http.StatusCreated, api.JSON{"commentId": commentID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func HandlePostComments(pool *pgxpool.Pool) http.HandlerFunc {
	galleriaService := factories.MakeGalleriaService(pool)

	return func(w http.ResponseWriter, r *http.Request) {
		postId, err := uuid.Parse(chi.URLParam(r, "postId"))
		if err != nil {
			api.HandleError(w, http.StatusBadRequest, api.Error{
				Message: "invalid post id",
				Details: "post id must be a valid UUID",
			})
			return
		}

		comments, err := galleriaService.GetComments(r.Context(), postId)
		if err != nil {
			if errors.Is(err, galleria.ErrImageNotFound) {
				api.HandleError(w, http.StatusNotFound, api.Error{Message: err.Error()})
				return
			}

			log.Printf("failed to get comments: %v", err)
			api.HandleError(
				w,
				http.StatusInternalServerError,
				api.Error{Message: "something went wrong, please try again"},
			)
			return
		}

		if err = api.Encode(w, http.StatusOK, api.JSON{"comments": comments}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
