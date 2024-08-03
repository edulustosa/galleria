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

type RegisterResponse struct {
	ID uuid.UUID `json:"id"`
}

func (a *Auth) Register(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	_, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHashBytes),
	}

	userId, err := a.usersRepository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{userId}, nil
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID uuid.UUID `json:"id"`
}

func (a *Auth) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	user, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	return &LoginResponse{user.ID}, nil
}
