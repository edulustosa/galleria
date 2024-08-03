package profile

import (
	"context"
	"errors"

	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
)

type Profile struct {
	usersRepository repo.UsersRepository
}

func New(usersRepository repo.UsersRepository) *Profile {
	return &Profile{usersRepository}
}

type UpdateProfileRequest struct {
	ID                uuid.UUID `json:"id"`
	Username          *string   `json:"username"`
	Bio               *string   `json:"bio"`
	ProfilePictureURL *string   `json:"profilePictureURL"`
}

func (p *Profile) Update(ctx context.Context, req *UpdateProfileRequest) error {
	user, err := p.usersRepository.FindByID(ctx, req.ID)
	if err != nil {
		return errors.New("invalid credentials")
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
