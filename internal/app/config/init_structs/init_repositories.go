package init_structs

import (
	"pr-assignment/internal/adapter/out/repository"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	teamRepo      *repository.TeamRepository
	prRepo        *repository.PullRequestRepository
	userRepo      *repository.UserRepository
	prReviewsRepo *repository.PrReviewersRepository
}

func InitRepositories(pool *pgxpool.Pool) Repositories {
	teamRepo := repository.NewTeamRepository(pool)
	prRepo := repository.NewPullRequestRepository(pool)
	userRepo := repository.NewUserRepository(pool)
	prReviewersRepo := repository.NewPrReviewersRepository(pool)

	return Repositories{
		teamRepo:      teamRepo,
		prRepo:        prRepo,
		userRepo:      userRepo,
		prReviewsRepo: prReviewersRepo,
	}
}
