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

func (r *PullRequestRepository) GetPR(ctx context.Context, pullRequestId string) (*model.PullRequest, error) {
	sql := `
        SELECT * FROM pr_repository
        WHERE pr_id = $1`

	row := r.pool.QueryRow(ctx, sql, pullRequestId)

	pullRequest := model.PullRequest{}
	err := row.Scan(
		&pullRequest.PullRequestId,
		&pullRequest.PullRequestName,
		&pullRequest.AuthorId,
		&pullRequest.Status,
		&pullRequest.CreatedAt,
		&pullRequest.MergedAt)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NOT_FOUND, "NO SUCH RESOURCE")
	}

	if err != nil {
		return nil, err
	}

	return &pullRequest, nil
}

func (r *PullRequestRepository) CreatePR(ctx context.Context, pr model.PullRequest) (*model.PullRequest, error) {
	sql := `
        INSERT INTO pull_requests(pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt)
        VALUES ($1, $2, $3, $4, $5, $6)
        ON CONFLICT (pull_request_id) DO NOTHING
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt
        `

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(ctx, sql, pr.PullRequestId, pr.PullRequestName,
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

func (r *PullRequestRepository) MergePR(ctx context.Context, pullRequestId string, status string, time time.Time) (*model.PullRequest, error) {
	sql := `
        UPDATE pull-requests
        SET status = $3, mergedAt = $2
        WHERE pull_request_id = $1
        RETURNING pull_request_id, pull_request_name, 
                                  author_id, status, createdAt, mergedAt`

	pullRequest := model.PullRequest{}
	err := r.pool.QueryRow(ctx, sql, pullRequestId, time, status).Scan(
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
