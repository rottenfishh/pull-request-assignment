package repository

import (
	"context"
	"fmt"
	"log"
	"pr-assignment/internal/model"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (r *TeamRepository) Init() {
	var err error
	r.pool, err = pgxpool.New(r.ctx, "postgres://<username>:<password>@localhost:5432/teams")

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := r.pool.Ping(r.ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
}

func (r *TeamRepository) Exists(teamName string) (bool, error) {
	sql := `
           SELECT team_name FROM teams
           WHERE team_name = $1`
	commandTag, err := r.pool.Exec(r.ctx, sql, teamName)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() > 0 {
		return true, model.NewError(model.TEAM_EXISTS, "TEAM EXISTS", teamName)
	}
	return false, nil
}
func (r *TeamRepository) AddTeam(newTeam string) error {
	sql := `
           INSERT INTO teams (team_name)`

	res, err := r.Exists(newTeam)
	if err != nil {
		return err
	}

	if !res {
		_, err = r.pool.Exec(r.ctx, sql, newTeam)
		if err != nil {
			return err
		}
	}
	return nil
}
