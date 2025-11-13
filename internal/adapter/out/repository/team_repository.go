package repository

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"pr-assignment/internal/model"
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

func (r *TeamRepository) AddTeam(newTeam model.Team, teamId uuid.UUID) error {
	sql := `
           INSERT INTO teams (team_id, team_name, user_id)`

	teamExists, err := r.Exists(newTeam.TeamName)

	if err != nil {
		return err
	}

	if teamExists {
		return model.NewError(model.TEAM_EXISTS, "TEAM EXISTS", "Team %s already exists", newTeam.TeamName)
	}

	for _, member := range newTeam.Members {
		_, err = r.pool.Exec(r.ctx, sql, teamId, newTeam.TeamName, member.UserId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *TeamRepository) GetUsersIdByTeam(teamId string) ([]string, error) {
	sql := `
        SELECT user_id FROM teams 
        WHERE team_id = $1`
	userIds := make([]string, 0)
	rows, err := r.pool.Query(r.ctx, sql, teamId)
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
