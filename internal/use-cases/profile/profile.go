package profile

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
)

type Profile struct {
	usersRepository  repo.UsersRepository
	imagesRepository repo.ImagesRepository
}

func New(usersRepository repo.UsersRepository, imagesRepository repo.ImagesRepository) *Profile {
	return &Profile{
		usersRepository,
		imagesRepository,
	}
}

type UpdateProfileRequest struct {
	ID                uuid.UUID `json:"id"`
	Username          *string   `json:"username,omitempty"`
	Bio               *string   `json:"bio,omitempty"`
	ProfilePictureURL *string   `json:"profilePictureURL,omitempty"`
}

var ErrInvalidCredentials error

func (p *Profile) Update(ctx context.Context, req *UpdateProfileRequest) error {
	user, err := p.usersRepository.FindByID(ctx, req.ID)
	if err != nil {
		return ErrInvalidCredentials
	}

	if req.Username != nil {
		user.Username = *req.Username
	}

	if req.Bio != nil {
		user.Bio = req.Bio
	}

	if req.ProfilePictureURL != nil {
		user.ProfilePictureURL = req.ProfilePictureURL
	}

	return p.usersRepository.Update(ctx, user)
}

func (p *Profile) GetProfile(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := p.usersRepository.FindByID(ctx, id)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (p *Profile) GetProfileImages(ctx context.Context, id uuid.UUID) ([]models.Image, error) {
	images, err := p.imagesRepository.GetImagesByUserID(ctx, id)
	if err != nil {
		return nil, err
	}

	return images, nil
}
