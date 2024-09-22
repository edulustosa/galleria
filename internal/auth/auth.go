package auth

import (
	"context"
	"errors"
	"fmt"
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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r RegisterRequest) Valid() (problems map[string]string) {
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

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

func (a *Auth) Register(
	ctx context.Context,
	req *RegisterRequest,
) (uuid.UUID, error) {
	_, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err == nil {
		return uuid.Nil, ErrUserAlreadyExists
	}

	passwordHashBytes, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return uuid.Nil, err
	}

	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(passwordHashBytes),
	}

	return a.usersRepository.Create(ctx, user)
}

type LoginRequest struct {
	Email    string
	Password string
}

func (r LoginRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	_, err := mail.ParseAddress(r.Email)
	if err != nil {
		problems["email"] = fmt.Sprintf("%s is not a valid email", r.Email)
	}

	if len(r.Password) < 8 || len(r.Password) > 128 {
		problems["password"] = "password must be between 8 and 128 characters"
	}

	return problems
}

func (a *Auth) Login(ctx context.Context, req *LoginRequest) (uuid.UUID, error) {
	user, err := a.usersRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(req.Password),
	)
	if err != nil {
		return uuid.Nil, ErrInvalidCredentials
	}

	return user.ID, nil
}
