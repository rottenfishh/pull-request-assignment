package repository

import (
	"context"
	"pr-assignment/internal/model"

	"github.com/google/uuid"
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
           SELECT * FROM teams
           WHERE team_name = $1`
	commandTag, err := r.pool.Exec(ctx, sql, teamName)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() > 0 {
		return true, model.NewError(model.TEAM_EXISTS, "TEAM %s EXISTS", teamName)
	}
	return false, nil
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
