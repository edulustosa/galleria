package auth_test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/repositories"
	"github.com/edulustosa/galleria/internal/use-cases/auth"
	testutils "github.com/edulustosa/galleria/test/utils"
)

func TestAuth(t *testing.T) {
	dbpool, err := testutils.LoadDatabase()
	if err != nil {
		t.Fatal("Failed to connect with database", err.Error())
	}

	usersRepository := repositories.NewPGXUsersRepository(dbpool)
	authUseCase := auth.New(usersRepository)

	t.Run("user should be able to register", func(t *testing.T) {
		req := &auth.RegisterRequest{
			Username: "john doe",
			Email:    "johndoe@email.com",
			Password: "12345678",
		}

		userId, err := authUseCase.Register(context.Background(), req)
		if err != nil {
			t.Fatal("Failed to create user", err.Error())
		}

		t.Logf("User id %v", userId)
	})
}
