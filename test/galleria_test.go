package test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/galleria"
)

func TestGalleria(t *testing.T) {
	pool, err := LoadDatabase()
	if err != nil {
		t.Fatalf("failed to connect with database: %v", err)
	}

	usersRepository := repo.NewPGXUsersRepository(pool)
	imagesRepository := repo.NewPGXImagesRepository(pool)
	sut := galleria.New(usersRepository, imagesRepository)

	testCtx := context.Background()

	t.Run("users should be able to send a image", func(t *testing.T) {
		if err = TruncateTables(pool); err != nil {
			t.Fatalf("failed to truncate tables: %v", err)
		}

		userId, err := SignUpUser(usersRepository)
		if err != nil {
			t.Fatalf("failed to sign up user: %v", err)
		}

		req := &galleria.SendImageRequest{
			Title: "image title",
			URL:   "http://image.com",
		}

		imageId, err := sut.SendImage(testCtx, userId, req)
		if err != nil {
			t.Fatalf("failed to send image: %v", err)
		}

		image, err := imagesRepository.GetImageByID(testCtx, imageId)
		if err != nil {
			t.Fatalf("failed to get image by id: %v", err)
		}

		if image.Title != req.Title || image.URL != req.URL {
			t.Errorf("unexpected image: %v", image)
		}

		PrettyPrint(image)
	})
}
