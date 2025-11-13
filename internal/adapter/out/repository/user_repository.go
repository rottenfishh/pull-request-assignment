package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pr-assignment/internal/model"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (r *UserRepository) Init() {
	var err error
	r.pool, err = pgxpool.New(r.ctx, "postgres://<username>:<password>@localhost:5432/users")

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := r.pool.Ping(r.ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
}

func (r *UserRepository) UpdateUserStatus(userId string, newStatus bool) error {
	sql := `
        UPDATE users
        SET status = $1, updated_at = NOW()
        WHERE id = $2
    `

	commandTag, err := r.pool.Exec(r.ctx, sql, newStatus, userId)

	if err != nil {
		return fmt.Errorf("error updating user status: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return model.NewError(model.NOT_FOUND, "user not found", userId)
	}

	return nil
}

func (r *UserRepository) AddTeam(newTeam model.Team, teamId uuid.UUID) error {
	sql := `
        INSERT INTO users(user_id, username, team_name, is_active)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (user_id) DO UPDATE SET team_name = $3
        `

	for _, member := range newTeam.Members {
		_, err := r.pool.Exec(r.ctx, sql, member.UserId, member.Username,
			teamId, member.IsActive)

		if err != nil {
			return fmt.Errorf("error adding team on user: %s %w", member.UserId, err)
		}

	}
	return nil
}

func (r *UserRepository) GetTeam(teamName string) (*model.Team, error) {
	sql := `
        SELECT * FROM users WHERE team_name = $1`

	rows, err := r.pool.Query(r.ctx, sql, teamName)

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

func (r *UserRepository) GetTeamNameByUserId(userId string) (string, error) {
	sql := `
        SELECT team_name FROM users WHERE id = $1`

	row := r.pool.QueryRow(r.ctx, sql, userId)
	if errors.Is(row, pgx.ErrNoRows) {
		return "", model.NewError(model.NOT_FOUND, "team with userId %s not found", userId)
	}

	var teamName string
	err := row.Scan(&teamName)
	if err != nil {
		return "", fmt.Errorf("error scanning row: %w", err)
	}
	return teamName, nil
}
