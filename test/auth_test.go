package test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/use-cases/auth"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth_Register(t *testing.T) {
	dbpool, err := LoadDatabase()
	if err != nil {
		t.Fatal("Failed to connect with database", err.Error())
	}

	usersRepository := repo.NewPGXUsersRepository(dbpool)
	authUseCase := auth.New(usersRepository)

	t.Run("user should be able to register", func(t *testing.T) {
		if err = TruncateTables(dbpool); err != nil {
			t.Fatal("Failed to truncate tables", err.Error())
		}

		req := &auth.RegisterRequest{
			Username: "john doe",
			Email:    "johndoe@email.com",
			Password: "12345678",
		}

		resp, err := authUseCase.Register(context.Background(), req)
		if err != nil {
			t.Fatal("Failed to create user", err.Error())
		}

		t.Logf("User id: %v", resp.ID)
	})

	t.Run("user should not be able to register with same email", func(t *testing.T) {
		if err = TruncateTables(dbpool); err != nil {
			t.Fatal("Failed to truncate tables", err.Error())
		}

		user1 := &auth.RegisterRequest{
			Username: "john doe",
			Email:    "johndoe@email.com",
			Password: "12345678",
		}

		user2 := &auth.RegisterRequest{
			Username: "john doe",
			Email:    "johndoe@email.com",
			Password: "12345678",
		}

		_, err := authUseCase.Register(context.Background(), user1)
		if err != nil {
			t.Fatal("Failed to create user", err.Error())
		}

		_, err = authUseCase.Register(context.Background(), user2)
		if err == nil {
			t.Fatal("User 2 created")
		}

		t.Log("User 2 not created:", err.Error())
	})
}

func TestAuth_Login(t *testing.T) {
	dbpool, err := LoadDatabase()
	if err != nil {
		t.Fatal("Failed to connect with database", err.Error())
	}

	usersRepository := repo.NewPGXUsersRepository(dbpool)
	authUseCase := auth.New(usersRepository)

	t.Run("user should be able to authenticate", func(t *testing.T) {
		if err = TruncateTables(dbpool); err != nil {
			t.Fatal("Failed to truncate tables", err.Error())
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

		usersRepository.Create(context.Background(), &models.User{
			Username:     "john doe",
			Email:        "johndoe@email.com",
			PasswordHash: string(hashedPassword),
		})

		req := &auth.LoginRequest{
			Email:    "johndoe@email.com",
			Password: "12345678",
		}

		resp, err := authUseCase.Login(context.Background(), req)
		if err != nil {
			t.Fatal("Failed to authenticate", err.Error())
		}

		t.Log("Authenticated successfully, id:", resp.ID)
	})
}
