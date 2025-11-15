package main

import (
	"context"
	_ "pr-assignment/docs"
	"pr-assignment/internal/adapter/in/http/handler"
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/app"
	"pr-assignment/internal/app/config/db"
	"pr-assignment/internal/service"
)

func main() {
	ctx := context.Background()
	dsn := "dsn"

	database, err := db.NewDb(ctx, dsn)
	if err != nil {
		panic(err)
	}
	defer database.Pool.Close()

	err = database.RunMigrations()
	if err != nil {
		panic(err)
	}

	teamRepo := repository.NewTeamRepository(database.Pool)
	prRepo := repository.NewPullRequestRepository(database.Pool)
	userRepo := repository.NewUserRepository(database.Pool)
	prReviewersRepo := repository.NewPrReviewersRepository(database.Pool)

	prService := service.NewPullRequestService(prRepo, prReviewersRepo, teamRepo, userRepo)
	userService := service.NewUserService(userRepo, teamRepo)
	statService := service.NewStatService(prReviewersRepo, userRepo, prRepo)

	prHandler := handler.NewPullRequestHandler(prService)
	userHandler := handler.NewUserHandler(userService, prService)
	statHandler := handler.NewStatHandler(statService)

	server := app.NewServer(prHandler, userHandler, statHandler)

	err = server.RunServer(ctx)
	if err != nil {
		panic(err)
	}
}
