package repositories

import (
	"context"

	"github.com/edulustosa/galleria/internal/database/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UsersRepository interface {
	Create(ctx context.Context, user *models.User) (uuid.UUID, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
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
