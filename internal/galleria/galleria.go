package galleria

import (
	"context"

	"github.com/edulustosa/galleria/helpers"
	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
)

type Galleria struct {
	usersRepository  repo.UsersRepository
	imagesRepository repo.ImagesRepository
	// commentsRepository repo.CommentsRepository
}

func New(
	usersRepository repo.UsersRepository,
	imagesRepository repo.ImagesRepository,
	// commentsRepository repo.CommentsRepository,
) *Galleria {
	return &Galleria{
		usersRepository:  usersRepository,
		imagesRepository: imagesRepository,
		// commentsRepository: commentsRepository,
	}
}

// Unimplemented method
func (g *Galleria) Display(ctx context.Context) ([]models.Image, error) {
	return nil, nil
}

type SendImageRequest struct {
	Title       string  `json:"title"`
	Author      *string `json:"author"`
	Description *string `json:"description"`
	URL         string  `json:"url"`
}

func (r SendImageRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if r.Title == "" || len(r.Title) > 255 {
		problems["title"] = "title must be between 1 and 255 characters"
	}

	if r.Author != nil && len(*r.Author) > 50 {
		problems["author"] = "author must be less than 255 characters"
	}

	if r.Description != nil && len(*r.Description) > 500 {
		problems["description"] = "description must be less than 255 characters"
	}

	if err := helpers.ValidateURL(r.URL); err != nil {
		problems["url"] = "invalid url scheme"
	}

	return problems
}

func (g *Galleria) SendImage(
	ctx context.Context,
	userId uuid.UUID,
	req *SendImageRequest,
) (imageId uuid.UUID, err error) {
	image := &models.Image{
		Title:       req.Title,
		UserID:      userId,
		Author:      req.Author,
		Description: req.Description,
		URL:         req.URL,
	}

	return g.imagesRepository.Create(ctx, image)
}
