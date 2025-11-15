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
	dsn := "here goes db info"

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

	prHandler := handler.NewPullRequestHandler(prService)
	userHandler := handler.NewUserHandler(userService, prService)

	server := app.NewServer(prHandler, userHandler)
	server.RunServer(ctx)
}
