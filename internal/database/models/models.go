package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                uuid.UUID
	Username          string
	Email             string
	PasswordHash      string
	Bio               *string
	ProfilePictureURL *string
	CreatedAt         pgtype.Timestamp
	UpdatedAt         pgtype.Timestamp
}
