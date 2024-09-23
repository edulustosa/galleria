package test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/profile"
	"golang.org/x/crypto/bcrypt"
)

func TestAuth_Profile(t *testing.T) {
	dbpool, err := LoadDatabase()
	if err != nil {
		t.Fatal("Failed to connect with database", err.Error())
	}

	usersRepository := repo.NewPGXUsersRepository(dbpool)
	imagesRepository := repo.NewPGXImagesRepository(dbpool)
	profileService := profile.New(usersRepository, imagesRepository)

	t.Run("user should be able to update profile", func(t *testing.T) {
		if err = TruncateTables(dbpool); err != nil {
			t.Fatal("Failed to truncate tables", err.Error())
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

		userID, _ := usersRepository.Create(context.Background(), &models.User{
			Username:     "john doe",
			Email:        "johndoe@email.com",
			PasswordHash: string(hashedPassword),
		})

		newUsername := "Robert"
		bio := "This is the new bio"
		profilePictureURL := "http://my_url.com"

		update := &profile.UpdateProfileRequest{
			Username:          &newUsername,
			Bio:               &bio,
			ProfilePictureURL: &profilePictureURL,
		}

		err := profileService.Update(context.Background(), userID, update)
		if err != nil {
			t.Fatal("Failed to update user:", err.Error())
		}

		user, _ := usersRepository.FindByID(context.Background(), userID)
		PrettyPrint(user)
	})

	t.Run("fields not specified should keep unchanged", func(t *testing.T) {
		if err = TruncateTables(dbpool); err != nil {
			t.Fatal("Failed to truncate tables", err.Error())
		}

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

		userID, _ := usersRepository.Create(context.Background(), &models.User{
			Username:     "john doe",
			Email:        "johndoe@email.com",
			PasswordHash: string(hashedPassword),
		})

		bio := "This is the new bio"

		update := &profile.UpdateProfileRequest{
			Bio: &bio,
		}

		err := profileService.Update(context.Background(), userID, update)
		if err != nil {
			t.Fatal("Failed to update user:", err.Error())
		}

		user, _ := usersRepository.FindByID(context.Background(), userID)
		PrettyPrint(user)
	})
}
