package repository

import (
	"context"
	"errors"
	"fmt"
	"pr-assignment/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// pr id user id

type PrReviewersRepository struct {
	pool *pgxpool.Pool
}

func NewPrReviewersRepository(pool *pgxpool.Pool) *PrReviewersRepository {
	return &PrReviewersRepository{pool: pool}
}

func (r *PrReviewersRepository) addReviewer(ctx context.Context, pullRequestId string, authorId string, reviewerId string) error {
	sql := `
         INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2);`

	_, err := r.pool.Exec(ctx, sql, pullRequestId, authorId, reviewerId)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrReviewersRepository) ChangeReviewer(ctx context.Context, pullRequestId string, oldReviewerId string, newReviewerId string) error {
	sql := `
        UPDATE pr_reviewers
        SET reviewer_id = $2
        WHERE pull_request_id = $1 AND reviewer_id = $3`

	_, err := r.pool.Exec(ctx, sql, pullRequestId, newReviewerId, oldReviewerId)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrReviewersRepository) GetReviewers(ctx context.Context, pullRequestId string) ([]string, error) {
	sql := `
        SELECT reviewer_id FROM pr_reviewers
        WHERE pull_request_id = $1`

	rows, err := r.pool.Query(ctx, sql, pullRequestId)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviewersIds []string
	var reviewerId string

	for rows.Next() {
		err = rows.Scan(&reviewerId)
		if err != nil {
			return nil, err
		}
		reviewersIds = append(reviewersIds, reviewerId)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NOT_FOUND, "PR reviewers not found")
	}

	return reviewersIds, nil
}

func (r *PrReviewersRepository) GetPRsByUser(ctx context.Context, userId string) ([]string, error) {
	sql := `
        SELECT pull_request_id FROM pr_reviewers
        WHERE reviewer_id = $1`
	rows, err := r.pool.Query(ctx, sql, userId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prIds []string
	var prId string

	for rows.Next() {
		err = rows.Scan(&prId)

		if err != nil {
			return nil, err
		}

		prIds = append(prIds, prId)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return prIds, nil
}
