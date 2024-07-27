package testutils

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func LoadDatabase() (*pgxpool.Pool, error) {
	if err := godotenv.Load("../../../.env"); err != nil {
		return nil, err
	}

	return pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
}
