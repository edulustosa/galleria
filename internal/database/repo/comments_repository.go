package repo

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LikeUpdate int8

const (
	Increment LikeUpdate = 1
	Decrement LikeUpdate = -1
)

type CommentsRepository interface {
	Create(ctx context.Context, comment *models.Comment) (uuid.UUID, error)
	FindByID(ctx context.Context, commentID uuid.UUID) (*models.Comment, error)

	FindByImageID(
		ctx context.Context,
		imageID uuid.UUID,
	) ([]models.Comment, error)

	UpdateLikes(
		ctx context.Context,
		commentID uuid.UUID,
		update LikeUpdate,
	) error
}

type PGXCommentsRepository struct {
	pool *pgxpool.Pool
}

func NewPGXCommentsRepository(pool *pgxpool.Pool) CommentsRepository {
	return &PGXCommentsRepository{pool}
}

const findCommentByIDQuery = "SELECT * FROM comments WHERE id = $1;"

func (r *PGXCommentsRepository) FindByID(
	ctx context.Context,
	commentID uuid.UUID,
) (*models.Comment, error) {
	var comment models.Comment
	err := r.pool.QueryRow(ctx, findCommentByIDQuery, commentID).Scan(
		&comment.ID,
		&comment.UserID,
		&comment.ImageID,
		&comment.Content,
		&comment.Likes,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &comment, nil
}

const createCommentQuery = `
	INSERT INTO comments (user_id, image_id, content)
	VALUES ($1, $2, $3)
	RETURNING id
`

func (r *PGXCommentsRepository) Create(
	ctx context.Context,
	comment *models.Comment,
) (uuid.UUID, error) {
	var commentID uuid.UUID
	err := r.pool.QueryRow(
		ctx,
		createCommentQuery,
		comment.UserID,
		comment.ImageID,
		comment.Content,
	).Scan(&commentID)
	if err != nil {
		return uuid.Nil, err
	}

	return commentID, nil
}

const findCommentsByImageIDQuery = `
	SELECT
		comments.id,
		comments.user_id,
		comments.image_id,
		comments.content,
		comments.likes,
		users.username,
		users.profile_picture_url
	FROM comments
	JOIN users ON comments.user_id = users.id
	ORDER BY comments.likes
	WHERE comments.image_id = $1;
`

func (r *PGXCommentsRepository) FindByImageID(
	ctx context.Context,
	imageID uuid.UUID,
) ([]models.Comment, error) {
	rows, err := r.pool.Query(ctx, findCommentsByImageIDQuery, imageID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var comment models.Comment

		err := rows.Scan(
			&comment.ID,
			&comment.UserID,
			&comment.ImageID,
			&comment.Content,
			&comment.Likes,
			&comment.Username,
			&comment.Avatar,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, comment)
	}

	return comments, nil
}

const updateLikesQuery = "UPDATE comments SET likes = likes + $1 WHERE id = $2;"

func (r *PGXCommentsRepository) UpdateLikes(
	ctx context.Context,
	commentID uuid.UUID,
	update LikeUpdate,
) error {
	_, err := r.pool.Exec(ctx, updateLikesQuery, update, commentID)
	if err != nil {
		return err
	}

	return nil
}
