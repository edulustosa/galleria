package test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/galleria"
	"github.com/google/uuid"
)

func TestGalleria_SendImage(t *testing.T) {
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

	t.Run(
		"users should not be able to send a image with a wrong id",
		func(t *testing.T) {
			if err = TruncateTables(pool); err != nil {
				t.Fatalf("failed to truncate tables: %v", err)
			}

			req := &galleria.SendImageRequest{
				Title: "image title",
				URL:   "http://image.com",
			}

			_, err := sut.SendImage(testCtx, uuid.New(), req)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}

			if err.Error() != galleria.ErrUserNotFound.Error() {
				t.Errorf("unexpected error: %v", err)
			}
		},
	)
}

func TestGalleria_Display(t *testing.T) {
	pool, err := LoadDatabase()
	if err != nil {
		t.Fatalf("failed to connect with database: %v", err)
	}

	usersRepository := repo.NewPGXUsersRepository(pool)
	imagesRepository := repo.NewPGXImagesRepository(pool)
	sut := galleria.New(usersRepository, imagesRepository)

	ctx := context.Background()

	t.Run("users should be able to view images on the galleria", func(t *testing.T) {
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

		_, err = sut.SendImage(ctx, userId, req)
		if err != nil {
			t.Fatalf("failed to send image: %v", err)
		}

		posts, err := sut.Display(ctx, 1)
		if err != nil {
			t.Fatalf("failed to display images: %v", err)
		}

		if len(posts) != 1 {
			t.Errorf("unexpected number of posts: %v", len(posts))
		}

		PrettyPrint(posts)
	})

	t.Run("users should be able to view images paginated", func(t *testing.T) {
		if err = TruncateTables(pool); err != nil {
			t.Fatalf("failed to truncate tables: %v", err)
		}

		userId, err := SignUpUser(usersRepository)
		if err != nil {
			t.Fatalf("failed to sign up user: %v", err)
		}

		for i := 0; i < 22; i++ {
			req := &galleria.SendImageRequest{
				Title: "image title",
				URL:   "http://image.com",
			}

			_, err = sut.SendImage(ctx, userId, req)
			if err != nil {
				t.Fatalf("failed to send image: %v", err)
			}
		}

		posts, err := sut.Display(ctx, 2)
		if err != nil {
			t.Fatalf("failed to display images: %v", err)
		}

		if len(posts) != 2 {
			t.Errorf("unexpected number of posts: %v", len(posts))
		}

		PrettyPrint(posts)
	})
}
