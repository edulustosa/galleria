package repo

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ImagesRepository interface {
	Create(ctx context.Context, image *models.Image) (uuid.UUID, error)
	GetImageByID(ctx context.Context, imageID uuid.UUID) (*models.Image, error)
	GetImagesByUserID(ctx context.Context, userID uuid.UUID) ([]models.Image, error)

	FindMany(ctx context.Context, page uint64) ([]models.Image, error)
}

type PGXImagesRepository struct {
	db *pgxpool.Pool
}

func NewPGXImagesRepository(db *pgxpool.Pool) ImagesRepository {
	return &PGXImagesRepository{
		db,
	}
}

const getImagesByUserIDQuery = "SELECT * FROM images WHERE user_id = $1"

func (r *PGXImagesRepository) GetImagesByUserID(
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

const createImageQuery = `
	INSERT INTO images (
		"user_id",
		"title",
		"author",
		"description",
		"url"
	) VALUES ($1, $2, $3, $4, $5)
	RETURNING "id";
`

func (r *PGXImagesRepository) Create(
	ctx context.Context,
	image *models.Image,
) (uuid.UUID, error) {
	row := r.db.QueryRow(
		ctx,
		createImageQuery,
		image.UserID,
		image.Title,
		image.Author,
		image.Description,
		image.URL,
	)

	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getImageByIDQuery = "SELECT * FROM images WHERE id = $1;"

func (r *PGXImagesRepository) GetImageByID(
	ctx context.Context,
	imageID uuid.UUID,
) (*models.Image, error) {
	row := r.db.QueryRow(ctx, getImageByIDQuery, imageID)

	var image models.Image
	err := row.Scan(
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

	return &image, nil
}

const findManyQuery = `
	SELECT * FROM images
	ORDER BY "likes"
	LIMIT $1
	OFFSET $2;
`

const ITEMS_PER_PAGE = 20

func (r *PGXImagesRepository) FindMany(ctx context.Context, page uint64) ([]models.Image, error) {
	skip := (page - 1) * ITEMS_PER_PAGE

	rows, err := r.db.Query(ctx, findManyQuery, ITEMS_PER_PAGE, skip)
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
