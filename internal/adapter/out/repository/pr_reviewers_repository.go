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

func (r *PrReviewersRepository) AddReviewer(ctx context.Context, pullRequestID string, reviewerID string) error {
	sql := `
         INSERT INTO pr_reviewers (pull_request_id, reviewer_id) VALUES ($1, $2);`

	_, err := r.pool.Exec(ctx, sql, pullRequestID, reviewerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrReviewersRepository) ChangeReviewer(ctx context.Context, pullRequestID string, oldReviewerID string,
	newReviewerID string) error {
	sql := `
        UPDATE pr_reviewers
        SET reviewer_id = $2
        WHERE pull_request_id = $1 AND reviewer_id = $3`

	_, err := r.pool.Exec(ctx, sql, pullRequestID, newReviewerID, oldReviewerID)
	if err != nil {
		return err
	}

	return nil
}

func (r *PrReviewersRepository) GetReviewers(ctx context.Context, pullRequestID string) ([]string, error) {
	sql := `
        SELECT reviewer_id FROM pr_reviewers
        WHERE pull_request_id = $1`

	rows, err := r.pool.Query(ctx, sql, pullRequestID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var reviewersIDs []string
	var reviewerID string

	for rows.Next() {
		err = rows.Scan(&reviewerID)
		if err != nil {
			return nil, err
		}
		reviewersIDs = append(reviewersIDs, reviewerID)
	}

	if len(reviewersIDs) == 0 {
		return nil, model.NewError(model.NotFound, "PR reviewers not found")
	}

	return reviewersIDs, nil
}

func (r *PrReviewersRepository) GetPRsByUser(ctx context.Context, userID string) ([]string, error) {
	sql := `
        SELECT pull_request_id FROM pr_reviewers
        WHERE reviewer_id = $1`
	rows, err := r.pool.Query(ctx, sql, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prIDs []string
	var prID string

	for rows.Next() {
		err = rows.Scan(&prID)

		if err != nil {
			return nil, err
		}

		prIDs = append(prIDs, prID)
	}

	if len(prIDs) == 0 {
		return nil, model.NewError(model.NotFound, "PR reviewers not found")
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return prIDs, nil
}

// user - number of prs where they are reviewer
func (r *PrReviewersRepository) GetNumberOfReviewsByUser(ctx context.Context) (map[string]int, error) {
	sql := `
          SELECT reviewer_id, COUNT(*) FROM pr_reviewers
          GROUP BY reviewer_id`
	rows, err := r.pool.Query(ctx, sql)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	usersWithReviewCount := make(map[string]int)

	var reviewerID string
	var count int
	for rows.Next() {
		err = rows.Scan(&reviewerID, &count)
		if err != nil {
			return nil, err
		}
		usersWithReviewCount[reviewerID] = count
	}
	return usersWithReviewCount, nil
}

func (r *PrReviewersRepository) GetPrsWithReviewer(ctx context.Context) (map[string]int, error) {
	sql := `
          SELECT pull_request_id, COUNT(*) FROM pr_reviewers
          GROUP BY pull_request_id`

	rows, err := r.pool.Query(ctx, sql)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, model.NewError(model.NotFound, "PR reviewers not found")
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var prID string
	var count int

	prsMap := make(map[string]int)

	for rows.Next() {
		err = rows.Scan(&prID, &count)
		if err != nil {
			return nil, err
		}
		prsMap[prID] = count
	}

	return prsMap, nil
}
