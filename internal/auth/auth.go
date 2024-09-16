package auth

import (
	"context"
	"errors"
	"net/mail"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	usersRepository repo.UsersRepository
}

func New(usersRepository repo.UsersRepository) *Auth {
	return &Auth{
		usersRepository,
	}
}

type RegisterRequest struct {
	Username string
	Email    string
	Password string
}

func (r *RegisterRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		problems["email"] = "invalid email"
	}

	if len(r.Password) < 8 || len(r.Password) > 128 {
		problems["password"] = "password must be between 8 and 128 characters"
	}

	if len(r.Username) < 3 || len(r.Username) > 32 {
		problems["username"] = "username must be between 3 and 32 characters"
	}

	return problems
}

type UserIDResponse struct {
	ID uuid.UUID `json:"id"`
}

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) Register(
	ctx context.Context,
	req *RegisterRequest,
) (*UserIDResponse, error) {
	_, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
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

	return &UserIDResponse{userId}, nil
}

type LoginRequest struct {
	Email    string
	Password string
}

func (r *LoginRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		problems["email"] = "invalid email"
	}

	if len(r.Password) < 8 || len(r.Password) > 128 {
		problems["password"] = "password must be between 8 and 128 characters"
	}

	return problems
}

func (a *Auth) Login(
	ctx context.Context,
	req *LoginRequest,
) (*UserIDResponse, error) {
	user, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return &UserIDResponse{user.ID}, nil
}
