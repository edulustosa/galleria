package factories

import (
	"github.com/edulustosa/galleria/internal/database/repo"
	"github.com/edulustosa/galleria/internal/profile"
	"github.com/jackc/pgx/v5/pgxpool"
)

func MakeProfileService(pool *pgxpool.Pool) *profile.Profile {
	usersRepository := repo.NewPGXUsersRepository(pool)
	imagesRepository := repo.NewPGXImagesRepository(pool)
	return profile.New(usersRepository, imagesRepository)
}
