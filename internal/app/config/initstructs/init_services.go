package initstructs

import (
	"pr-assignment/internal/service"
)

type Services struct {
	userService        *service.UserService
	pullRequestService *service.PullRequestService
	statService        *service.StatService
}

func InitServices(repos Repositories) Services {
	userService := service.NewUserService(repos.userRepo, repos.teamRepo)
	prService := service.NewPullRequestService(repos.prRepo, repos.prReviewsRepo, repos.teamRepo, repos.userRepo, userService)
	statService := service.NewStatService(repos.prReviewsRepo, repos.userRepo, repos.prRepo)

	return Services{
		userService:        userService,
		pullRequestService: prService,
		statService:        statService,
	}
}
