package repo

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	Update(ctx context.Context, user *models.User) error
}

type PGXUsersRepository struct {
	db *pgxpool.Pool
}

func NewPGXUsersRepository(db *pgxpool.Pool) UsersRepository {
	return &PGXUsersRepository{db}
}

const createUser = `
	INSERT INTO users (
		"username",
		"email",
		"password_hash"
	) VALUES ($1, $2, $3)
	RETURNING "id";
`

func (r *PGXUsersRepository) Create(ctx context.Context, user *models.User) (uuid.UUID, error) {
	row := r.db.QueryRow(
		ctx,
		createUser,
		user.Username,
		user.Email,
		user.PasswordHash,
	)

	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const findUserByEmail = "SELECT * FROM users WHERE email = $1;"

func (r *PGXUsersRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	row := r.db.QueryRow(ctx, findUserByEmail, email)

	var user models.User
	err := row.Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ProfilePictureURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

const findByIDQuery = "SELECT * FROM users WHERE id = $1;"

func (r *PGXUsersRepository) FindByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	var user models.User
	err := r.db.QueryRow(ctx, findByIDQuery, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Bio,
		&user.ProfilePictureURL,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	return &user, err
}

const update = `
	UPDATE users
	SET "username" = $1, "bio" = $2, "profile_picture_url" = $3
	WHERE id = $4;
`

func (r *PGXUsersRepository) Update(ctx context.Context, user *models.User) error {
	_, err := r.db.Exec(ctx, update,
		user.Username,
		user.Bio,
		user.ProfilePictureURL,
		user.ID,
	)

	return err
}
