package repo

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImagesRepository interface {
	GetImagesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Image, error)
}

type PGXImagesRepo struct {
	db *pgxpool.Pool
}

func NewPGXImagesRepository(db *pgxpool.Pool) ImagesRepository {
	return &PGXImagesRepo{
		db,
	}
}

const getImagesByUserIDQuery = "SELECT * FROM images WHERE user_id = $1"

func (r *PGXImagesRepo) GetImagesByUserID(
	ctx context.Context,
	userID uuid.UUID,
) ([]models.Image, error) {
	rows, err := r.db.Query(ctx, getImagesByUserIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []models.Image
	for rows.Next() {
		var image models.Image

		err := rows.Scan(
			&image.ID,
			&image.UserID,
			&image.Title,
			&image.Author,
			&image.Description,
			&image.URL,
			&image.Likes,
			&image.CreatedAt,
			&image.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		images = append(images, image)
	}

	return images, nil
}
