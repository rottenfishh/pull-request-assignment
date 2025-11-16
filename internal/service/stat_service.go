package service

import (
	"context"
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
)

type StatService struct {
	prReviewsRepo *repository.PrReviewersRepository
	userRepo      *repository.UserRepository
	prRepo        *repository.PullRequestRepository
}

func NewStatService(prReviewsRepo *repository.PrReviewersRepository, userRepo *repository.UserRepository,
	prRepo *repository.PullRequestRepository) *StatService {
	return &StatService{prReviewsRepo: prReviewsRepo, userRepo: userRepo, prRepo: prRepo}
}

func (s *StatService) GetReviewsCountedByUser(ctx context.Context) ([]model.UserReviewsCount, error) {
	userReviews := []model.UserReviewsCount{}

	userMap, err := s.prReviewsRepo.GetNumberOfReviewsByUser(ctx)
	if err != nil {
		return nil, err
	}

	for userID, count := range userMap {
		user, err := s.userRepo.GetUserByID(ctx, userID)
		if err != nil {
			return nil, err
		}

		userReviews = append(userReviews, model.UserReviewsCount{User: *user, ReviewsCount: count})
	}

	return userReviews, nil
}

func (s *StatService) GetReviewsCountedByPR(ctx context.Context) ([]model.PrReviewersCount, error) {
	prReviewers := []model.PrReviewersCount{}

	prMap, err := s.prReviewsRepo.GetPrsWithReviewer(ctx)
	if err != nil {
		return nil, err
	}

	for prID, count := range prMap {
		pr, err := s.prRepo.GetPR(ctx, prID)
		if err != nil {
			return nil, err
		}

		prReviewers = append(prReviewers, model.PrReviewersCount{PullRequest: pr.PullRequestShort, Count: count})
	}

	return prReviewers, nil
}
