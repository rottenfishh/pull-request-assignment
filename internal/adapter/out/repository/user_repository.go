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

func (r *UserRepository) UpdateUserStatus(ctx context.Context, userID string, newStatus bool) (*model.User, error) {
	sql := `
        UPDATE users
        SET is_active = $1
        WHERE user_id = $2
        RETURNING user_id, username, team_name, is_active
    `

	row := r.pool.QueryRow(ctx, sql, newStatus, userID)

	user := model.User{}
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NotFound, "user not found %s", userID)
	}

	if err != nil {
		return nil, fmt.Errorf("error updating user status: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) AddTeam(ctx context.Context, newTeam model.Team, teamID uuid.UUID) error {
	sql := `
        INSERT INTO users(user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE SET team_name = $3
        `

	for _, member := range newTeam.Members {
		_, err := r.pool.Exec(ctx, sql, member.UserID, member.Username,
			teamID, member.IsActive)

		if err != nil {
			return fmt.Errorf("error adding team on user: %s %w", member.UserID, err)
		}

	}
	return nil
}

func (r *UserRepository) GetTeam(ctx context.Context, teamID string) (*model.Team, error) {
	sql := `
        SELECT * FROM users WHERE team_name = $1`

	rows, err := r.pool.Query(ctx, sql, teamID)

	if err != nil {
		fmt.Println("error getting team userrepo")
		return nil, model.NewError(model.NotFound, "team not found %s", teamID)
	}

	defer rows.Close()

	team := model.Team{
		TeamName: teamID,
		Members:  make([]model.TeamMember, 0),
	}

	for rows.Next() {
		teamMember := model.TeamMember{}
		err = rows.Scan(
			&teamMember.UserID,
			&teamMember.Username,
			&teamID,
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

func (r *UserRepository) GetTeamNameByUserID(ctx context.Context, userID string) (string, error) {
	sql := `
        SELECT team_name FROM users WHERE user_id = $1`

	row := r.pool.QueryRow(ctx, sql, userID)

	var teamName string
	err := row.Scan(&teamName)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.NewError(model.NotFound, "team with userID %s not found", userID)
	}

	if err != nil {
		return "", fmt.Errorf("error scanning row: %w", err)
	}
	return teamName, nil
}

func (r *UserRepository) GetActiveUsersByTeam(ctx context.Context, teamID string) ([]string, error) {
	sql := `
        SELECT user_id FROM users
        WHERE team_name = $1 
        AND is_active = true`

	userIDs := make([]string, 0)
	rows, err := r.pool.Query(ctx, sql, teamID)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}

	defer rows.Close()

	var userID string
	for rows.Next() {
		err = rows.Scan(&userID)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, userID)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	if len(userIDs) == 0 {
		return nil, model.NewError(model.NotFound, "team not found or has no users %s", teamID)
	}

	return userIDs, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	sql := `
           SELECT * FROM users WHERE user_id = $1`
	row := r.pool.QueryRow(ctx, sql, userID)
	user := model.User{}
	err := row.Scan(&user.UserID, &user.Username, &user.TeamName, &user.IsActive)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NotFound, "user not found %s", userID)
	}
	if err != nil {
		return nil, fmt.Errorf("error scanning row: %w", err)
	}
	return &user, nil
}
