package test

import (
	"context"
	"testing"

	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/galleria"
)

func TestComment(t *testing.T) {
	pool, err := LoadDatabase()
	if err != nil {
		t.Fatalf("failed to connect with database: %v", err)
	}

	usersRepository := repo.NewPGXUsersRepository(pool)
	imagesRepository := repo.NewPGXImagesRepository(pool)
	commentsRepository := repo.NewPGXCommentsRepository(pool)
	sut := galleria.New(usersRepository, imagesRepository, commentsRepository)

	ctx := context.Background()

	t.Run("users should be able to comment on a post", func(t *testing.T) {
		if err := TruncateTables(pool); err != nil {
			t.Fatalf("failed to truncate tables: %v", err)
		}

		userID, err := SignUpUser(usersRepository)
		if err != nil {
			t.Fatalf("failed to sign up user: %v", err)
		}

		imageID, err := CreateImage(imagesRepository, userID)
		if err != nil {
			t.Fatalf("failed to create image: %v", err)
		}

		comm := "this is a comment"

		commentID, err := sut.AddComment(ctx, userID, imageID, comm)
		if err != nil {
			t.Errorf("failed to add comment: %v", err)
		}

		comment, err := commentsRepository.FindByID(ctx, commentID)
		if err != nil {
			t.Fatalf("failed to find comment: %v", err)
		}

		if comment.Content != comm {
			t.Errorf("expected comment content to be %s, got %s", comm, comment.Content)
		}

		PrettyPrint(comment)
	})
}
