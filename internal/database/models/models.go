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

type Image struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Title       string
	Author      *string
	Description *string
	URL         string
	Likes       int
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}
