package model

import "time"

type User struct {
	ID                string
	Username          string
	Email             string
	PasswordHash      string
	Bio               *string
	ProfilePictureURL *string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
