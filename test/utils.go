package test

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
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
