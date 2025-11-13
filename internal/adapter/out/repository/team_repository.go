package repository

import (
	"context"
	"fmt"
	"pr-assignment/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// team id team name user id is_active
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
	commandTag, err := r.pool.Exec(ctx, sql, teamName)
	if err != nil {
		return false, err
	}
	if commandTag.RowsAffected() > 0 {
		return true, model.NewError(model.TEAM_EXISTS, "TEAM EXISTS", teamName)
	}
	return false, nil
}

func (r *TeamRepository) AddTeam(ctx context.Context, newTeam model.Team, teamId uuid.UUID) error {
	sql := `
           INSERT INTO teams (team_id, team_name, user_id, is_active)`

	teamExists, err := r.Exists(ctx, newTeam.TeamName)

	if err != nil {
		return err
	}

	if teamExists {
		return model.NewError(model.TEAM_EXISTS, "TEAM EXISTS", "Team %s already exists", newTeam.TeamName)
	}

	for _, member := range newTeam.Members {
		_, err = r.pool.Exec(ctx, sql, teamId, newTeam.TeamName, member.UserId, member.IsActive)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) GetActiveUsersByTeam(ctx context.Context, teamId string) ([]string, error) {
	sql := `
        SELECT user_id FROM teams 
        WHERE team_id = $1 
        AND is_active = true`
	userIds := make([]string, 0)
	rows, err := r.pool.Query(ctx, sql, teamId)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var userId string
	for rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			return nil, err
		}
		userIds = append(userIds, userId)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	if len(userIds) == 0 {
		return nil, model.NewError(model.NOT_FOUND, "team not found or has no users", teamId)
	}

	return userIds, nil
}
