package service

import (
	"context"
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
	"time"
)

type PullRequestService struct {
	prRepository          *repository.PullRequestRepository
	prReviewersRepository *repository.PrReviewersRepository
	teamRepository        *repository.TeamRepository
	userRepository        *repository.UserRepository
}

func NewPullRequestService(prRepo *repository.PullRequestRepository, prReviewsRepo *repository.PrReviewersRepository,
	teamRepo *repository.TeamRepository, userRepo *repository.UserRepository) *PullRequestService {
	return &PullRequestService{prRepo, prReviewsRepo, teamRepo, userRepo}
}

func (s *PullRequestService) createPR(ctx context.Context, pullRequest model.PullRequestShort) (*model.PullRequest, error) {
	pullRequest.Status = model.CREATED
	pr := model.PullRequest{
		PullRequestShort:  pullRequest,
		AssignedReviewers: make([]string, 0),
		CreatedAt:         time.Now(),
		MergedAt:          time.Time{},
	}

	createdPR, err := s.prRepository.CreatePR(ctx, pr)
	if err != nil {
		return nil, err
	}

	err = s.assignReviewers(ctx, createdPR)
	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PullRequestService) getActiveTeammatesByUser(ctx context.Context, userId string) ([]string, error) {
	teamId, err := s.userRepository.GetTeamNameByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	userIds, err := s.teamRepository.GetActiveUsersByTeam(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return userIds, nil
}

func (s *PullRequestService) assignReviewers(ctx context.Context, pr *model.PullRequest) error {
	teammates, err := s.getActiveTeammatesByUser(ctx, pr.AuthorId)

	if err != nil {
		return err
	}

	for _, userId := range teammates {
		if userId != pr.AuthorId {
			pr.AssignedReviewers = append(pr.AssignedReviewers, userId)
		}
		if len(pr.AssignedReviewers) == 2 {
			break
		}
	}
	return nil
}

// return who replaced
func (s *PullRequestService) changeReviewer(ctx context.Context, prId string, oldReviewerId string) (*model.PullRequest, error) {
	authorId, err := s.prReviewersRepository.GetAuthor(ctx, prId)
	if err != nil {
		return nil, err
	}

	teammates, err := s.teamRepository.GetActiveUsersByTeam(ctx, authorId)
	if err != nil {
		return nil, err
	}

	var newReviewerId string
	for _, userId := range teammates {
		if userId != oldReviewerId && userId != authorId {
			newReviewerId = userId
			break
		}
	}

	err = s.prReviewersRepository.ChangeReviewer(ctx, prId, oldReviewerId, newReviewerId)
	if err != nil {
		return nil, err
	}

	pr, err := s.prRepository.GetPR(ctx, prId)
	if err != nil {
		return nil, err
	}

	reviewers, err := s.prReviewersRepository.GetReviewers(ctx, prId)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers
	return pr, nil
}

func (s *PullRequestService) mergePR(ctx context.Context, pullRequestId string) (*model.PullRequest, error) {
	createdPR, err := s.prRepository.MergePR(ctx, pullRequestId, model.MERGED, time.Now())

	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PullRequestService) getPRsByUser(ctx context.Context, userId string) ([]*model.PullRequest, error) {
	pullRequestsIds, err := s.prReviewersRepository.GetPRsByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	var prs []*model.PullRequest
	for _, prId := range pullRequestsIds {
		pr, err := s.prRepository.GetPR(ctx, prId)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}
