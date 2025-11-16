package main

import (
	"context"
	"log"
	_ "pr-assignment/docs"
	"pr-assignment/internal/app"
	"pr-assignment/internal/app/config/db"
	"pr-assignment/internal/app/config/env"
	"pr-assignment/internal/app/config/initstructs"
)

func main() {
	ctx := context.Background()

	configDb, err := env.LoadConfigEnv()
	if err != nil {
		log.Fatalf("unable to load config: %e", err)
	}

	database, err := db.InitDatabase(ctx, *configDb)
	if err != nil {
		log.Fatalf("unable to init database: %e", err)
	}

	defer database.Pool.Close()

	repos := initstructs.InitRepositories(database.Pool)
	services := initstructs.InitServices(repos)
	handlers := initstructs.InitHandlers(services)

	server := app.NewServer(handlers.PullRequestHandler, handlers.UserHandler, handlers.StatHandler)

	err = server.RunServer()
	if err != nil {
		log.Fatalf("unable to run server: %e", err)
	}
}
