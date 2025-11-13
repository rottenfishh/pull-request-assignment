package cmd

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, "postgres://<username>:<password>@localhost:5432/pr-assignment")
}
