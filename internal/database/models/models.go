package models

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID                uuid.UUID        `json:"id"`
	Username          string           `json:"username"`
	Email             string           `json:"email"`
	PasswordHash      string           `json:"-"`
	Bio               *string          `json:"bio"`
	ProfilePictureURL *string          `json:"profilePictureURL"`
	CreatedAt         pgtype.Timestamp `json:"createdAt"`
	UpdatedAt         pgtype.Timestamp `json:"updatedAt"`
}

type Image struct {
	ID          uuid.UUID        `json:"id"`
	UserID      uuid.UUID        `json:"userId"`
	Title       string           `json:"title"`
	Author      *string          `json:"author"`
	Description *string          `json:"description"`
	URL         string           `json:"url"`
	CreatedAt   pgtype.Timestamp `json:"createdAt"`
	UpdatedAt   pgtype.Timestamp `json:"updatedAt"`
}

type Post struct {
	Image    Image   `json:"image"`
	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}

type Comment struct {
	ID        uuid.UUID        `json:"id"`
	UserID    uuid.UUID        `json:"userId"`
	ImageID   uuid.UUID        `json:"imageId"`
	Content   string           `json:"content"`
	CreatedAt pgtype.Timestamp `json:"createdAt"`
	UpdatedAt pgtype.Timestamp `json:"updatedAt"`

	Username string  `json:"username"`
	Avatar   *string `json:"avatar"`
}
