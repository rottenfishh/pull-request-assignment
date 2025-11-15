package init_structs

import (
	"pr-assignment/internal/adapter/in/http/handler"
)

type Handlers struct {
	UserHandler        *handler.UserHandler
	PullRequestHandler *handler.PullRequestHandler
	StatHandler        *handler.StatHandler
}

func InitHandlers(services Services) Handlers {
	userHandler := handler.NewUserHandler(services.userService, services.pullRequestService)
	prHandler := handler.NewPullRequestHandler(services.pullRequestService)
	statHandler := handler.NewStatHandler(services.statService)

	return Handlers{
		UserHandler:        userHandler,
		PullRequestHandler: prHandler,
		StatHandler:        statHandler,
	}
}
