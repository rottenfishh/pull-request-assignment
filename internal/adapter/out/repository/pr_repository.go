package repository

import (
	"context"
	"errors"
	"pr-assignment/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PullRequestRepository struct {
	pool *pgxpool.Pool
}

func NewPullRequestRepository(pool *pgxpool.Pool) *PullRequestRepository {
	return &PullRequestRepository{pool: pool}
}

func (r *PullRequestRepository) GetPR(ctx context.Context, pullRequestID string) (*model.PullRequest, error) {
	sql := `
        SELECT * FROM pull_requests
        WHERE pull_request_id = $1`

	row := r.pool.QueryRow(ctx, sql, pullRequestID)

	pullRequest := model.PullRequest{}
	err := row.Scan(
		&pullRequest.PullRequestID,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NotFound, "NO SUCH RESOURCE")
	}

	if err != nil {
		return nil, err
	}

	return &pullRequest, nil
}

func (r *PullRequestRepository) CreatePR(ctx context.Context, pr model.PullRequest) (*model.PullRequest, error) {
	sql := `
        INSERT INTO pull_requests(pull_request_id, pull_request_name, 
                                  author_id, status, created_at, merged_at)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (pull_request_id) DO NOTHING
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, created_at, merged_at
        `

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(ctx, sql, pr.PullRequestID, pr.PullRequestName,
		pr.AuthorID, pr.Status, pr.CreatedAt, pr.MergedAt).Scan(
		&pullRequest.PullRequestID,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.PrExists, "%s already exists", pr.PullRequestID)
	}

	if err != nil {
		return nil, err
	}

	return &pullRequest, nil
}

func (r *PullRequestRepository) MergePR(ctx context.Context, pullRequestID string, status model.PRstatus, time time.Time) (*model.PullRequest, error) {
	sql := `
        UPDATE pull_requests
        SET status = $3, merged_at = $2
        WHERE pull_request_id = $1
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, created_at, merged_at`

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(ctx, sql, pullRequestID, time, status).Scan(
		&pullRequest.PullRequestID,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorID,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NotFound, "%s not found", pullRequestID)
	}
	if err != nil {
		return nil, err
	}
	return &pullRequest, nil
}

func (r *PullRequestRepository) GetAuthor(ctx context.Context, pullRequestID string) (string, error) {
	sql := `
        SELECT author_id
        FROM pull_requests
        WHERE pull_request_id = $1`

	row := r.pool.QueryRow(ctx, sql, pullRequestID)

	var authorID string
	err := row.Scan(&authorID)

	if err != nil {
		return "", err
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return "", model.NewError(model.NotFound, "Pull request %s not found", pullRequestID)
	}

	return authorID, nil
}
