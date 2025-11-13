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

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) UpdateUserStatus(ctx context.Context, userId string, newStatus bool) error {
	sql := `
        UPDATE users
        SET status = $1, updated_at = NOW()
        WHERE id = $2
    `

	commandTag, err := r.pool.Exec(ctx, sql, newStatus, userId)

	if err != nil {
		return fmt.Errorf("error updating user status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return model.NewError(model.NOT_FOUND, "user not found", userId)
	}

	return nil
}

func (r *UserRepository) AddTeam(ctx context.Context, newTeam model.Team, teamId uuid.UUID) error {
	sql := `
        INSERT INTO users(user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE SET team_name = $3
        `

	for _, member := range newTeam.Members {
		_, err := r.pool.Exec(ctx, sql, member.UserId, member.Username,
			teamId, member.IsActive)

		if err != nil {
			return fmt.Errorf("error adding team on user: %s %w", member.UserId, err)
		}

	}
	return nil
}

func (r *UserRepository) GetTeam(ctx context.Context, teamName string) (*model.Team, error) {
	sql := `
        SELECT * FROM users WHERE team_name = $1`

	rows, err := r.pool.Query(ctx, sql, teamName)

	if err != nil {
		return nil, model.NewError(model.NOT_FOUND, "team not found", teamName)
	}

	defer rows.Close()

	team := model.Team{
		TeamName: teamName,
		Members:  make([]model.TeamMember, 0),
	}

	for rows.Next() {
		teamMember := model.TeamMember{}
		err = rows.Scan(
			&teamMember.UserId,
			&teamMember.Username,
			&teamMember.IsActive)
		if err != nil {
			return nil, fmt.Errorf("error scanning row: %w", err)
		}
		team.Members = append(team.Members, teamMember)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return &team, nil
}

func (r *UserRepository) GetTeamNameByUserId(ctx context.Context, userId string) (string, error) {
	sql := `
        SELECT team_name FROM users WHERE id = $1`

	row := r.pool.QueryRow(ctx, sql, userId)

	var teamName string
	err := row.Scan(&teamName)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.NewError(model.NOT_FOUND, "team with userId %s not found", userId)
	}

	if err != nil {
		return "", fmt.Errorf("error scanning row: %w", err)
	}
	return teamName, nil
}
