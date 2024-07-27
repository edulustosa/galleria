package auth

import (
	"context"
	"errors"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repositories"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	usersRepository repositories.UsersRepository
}

func New(usersRepository repositories.UsersRepository) *Auth {
	return &Auth{usersRepository}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (a *Auth) Register(ctx context.Context, req *RegisterRequest) (uuid.UUID, error) {
	_, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err == nil {
		return uuid.UUID{}, errors.New("user already exists")
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return uuid.UUID{}, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHashBytes),
	}

	return a.usersRepository.Create(ctx, user)
}
