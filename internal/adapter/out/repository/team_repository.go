package repository

import (
	"context"
	"errors"
	"fmt"
	"pr-assignment/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TeamRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *TeamRepository {
	return &TeamRepository{pool: pool}
}

func (r *TeamRepository) Exists(ctx context.Context, teamName string) (bool, error) {
	sql := `
           SELECT team_name FROM teams
           WHERE team_name = $1`
	fmt.Println("team name " + teamName)
	var name string
	queryRow := r.pool.QueryRow(ctx, sql, teamName)

	err := queryRow.Scan(&name)
	fmt.Println(err)
	if errors.Is(err, pgx.ErrNoRows) {
		return false, nil
	}

	if err != nil {
		return false, nil
	}

	return true, nil
}

func (r *TeamRepository) GetTeamId(ctx context.Context, teamName string) (string, error) {
	sql := `
           SELECT team_id FROM teams
           WHERE team_name = $1`

	queryRow := r.pool.QueryRow(ctx, sql, teamName)

	var teamId string
	err := queryRow.Scan(&teamId)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.NewError(model.NOT_FOUND, "%s team table not found", teamName)
	}
	if err != nil {
		return "", err
	}
	return teamId, nil
}

func (r *TeamRepository) AddTeam(ctx context.Context, newTeam model.Team, teamId uuid.UUID) error {
	sql := `
           INSERT INTO teams (team_id, team_name) VALUES ($1, $2)`

	teamExists, err := r.Exists(ctx, newTeam.TeamName)

	if err != nil {
		return err
	}

	if teamExists {
		return model.NewError(model.TEAM_EXISTS, "Team %s already exists", newTeam.TeamName)
	}

	_, err = r.pool.Exec(ctx, sql, teamId, newTeam.TeamName)

	return nil
}
