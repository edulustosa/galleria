package test

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func LoadDatabase() (*pgxpool.Pool, error) {
	if err := godotenv.Load("../.env"); err != nil {
		return nil, err
	}

	return pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
}

func TruncateTables(db *pgxpool.Pool) error {
	tables := []string{"users", "images", "comments"}

	for _, table := range tables {
		_, err := db.Exec(context.Background(), fmt.Sprintf("TRUNCATE %s CASCADE", table))
		if err != nil {
			return err
		}
	}

	return nil
}

func PrettyPrint(data ...any) {
	for _, d := range data {
		b, _ := json.MarshalIndent(d, "", "  ")
		fmt.Println(string(b))
	}
}

func SignUpUser(usersRepository repo.UsersRepository) (uuid.UUID, error) {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.DefaultCost)

	return usersRepository.Create(context.Background(), &models.User{
		Username:     "john doe",
		Email:        "johndoe@email.com",
		PasswordHash: string(hashedPassword),
	})
}

func CreateImage(imagesRepository repo.ImagesRepository, userID uuid.UUID) (uuid.UUID, error) {
	return imagesRepository.Create(context.Background(), &models.Image{
		UserID: userID,
		Title:  "image title",
		URL:    "https://example.com/image.jpg",
	})
}
