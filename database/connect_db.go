package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DB struct {
	pool *pgxpool.Pool
}

func connectDB() (*DB, error) {
	d := &DB{}
	ctx := context.Background()

	var err error
	d.pool, err = pgxpool.New(ctx, "postgres://<username>:<password>@localhost:5432/teams")

	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	if err := d.pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("unable to ping database: %v", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
	return d, nil
}
