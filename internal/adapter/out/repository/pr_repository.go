package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"pr-assignment/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func (r *PullRequestRepository) Init() {
	var err error
	r.pool, err = pgxpool.New(r.ctx, "postgres://<username>:<password>@localhost:5432/pull_requests")

	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := r.pool.Ping(r.ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
}

func (r *PullRequestRepository) CreatePR(pr model.PullRequest) (*model.PullRequest, error) {
	sql := `
        INSERT INTO pull_requests(pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (pull_request_id) DO NOTHING
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt
        `

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(r.ctx, sql, pr.PullRequestId, pr.PullRequestName,
		pr.AuthorId, pr.Status, pr.CreatedAt, pr.MergedAt).Scan(
		&pullRequest.PullRequestId,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorId,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.PR_EXISTS, "%s already exists", pr.PullRequestId)
	}
	if err != nil {
		return nil, err
	}
	return &pullRequest, nil
}

func (r *PullRequestRepository) MergePR(pullRequestId string, status string, time time.Time) (*model.PullRequest, error) {
	sql := `
        UPDATE pull-requests
        SET status = $3, mergedAt = $2
        WHERE pull_request_id = $1
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt`

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(r.ctx, sql, pullRequestId, time, status).Scan(
		&pullRequest.PullRequestId,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorId,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NOT_FOUND, "%s not found", pullRequestId)
	}
	if err != nil {
		return nil, err
	}
	return &pullRequest, nil
}
