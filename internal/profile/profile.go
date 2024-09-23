package profile

import (
	"context"
	"errors"
	"net/url"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
)

type Profile struct {
	usersRepository  repo.UsersRepository
	imagesRepository repo.ImagesRepository
}

func New(
	usersRepository repo.UsersRepository,
	imagesRepository repo.ImagesRepository,
) *Profile {
	return &Profile{
		usersRepository,
		imagesRepository,
	}
}

type UpdateProfileRequest struct {
	Username          *string `json:"username"`
	Bio               *string `json:"bio"`
	ProfilePictureURL *string `json:"avatar"`
}

func (r UpdateProfileRequest) Valid() (problems map[string]string) {
	problems = make(map[string]string)

	if r.Username != nil && (len(*r.Username) < 3 || len(*r.Username) > 32) {
		problems["username"] = "must be between 3 and 32 characters long"
	}

	if r.Bio != nil && len(*r.Bio) > 500 {
		problems["bio"] = "must be between 8 and 500 characters long"
	}

	if r.ProfilePictureURL != nil {
		parsedURL, err := url.Parse(*r.ProfilePictureURL)
		if err != nil {
			problems["profilePictureURL"] = "invalid URL"
		} else if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
			problems["profilePictureURL"] = "must be a valid HTTP or HTTPS url"
		}
	}

	return problems
}

var ErrInvalidCredentials = errors.New("invalid credentials")

func (p *Profile) Update(
	ctx context.Context,
	id uuid.UUID,
	req *UpdateProfileRequest,
) error {
	user, err := p.usersRepository.FindByID(ctx, id)
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
