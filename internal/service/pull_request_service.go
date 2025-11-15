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

func (s *PullRequestService) CreatePR(ctx context.Context, pullRequest model.PullRequestShort) (*model.PullRequest, error) {
	_, err := s.userRepository.GetUserById(ctx, pullRequest.AuthorId)
	if err != nil {
		return nil, err
	}

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

	err = s.AssignReviewers(ctx, createdPR)
	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PullRequestService) GetActiveTeammatesByUser(ctx context.Context, userId string) ([]string, error) {
	teamId, err := s.userRepository.GetTeamNameByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}

	userIds, err := s.userRepository.GetActiveUsersByTeam(ctx, teamId)
	if err != nil {
		return nil, err
	}

	return userIds, nil
}

func (s *PullRequestService) AssignReviewers(ctx context.Context, pr *model.PullRequest) error {
	teammates, err := s.GetActiveTeammatesByUser(ctx, pr.AuthorId)

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
func (s *PullRequestService) ChangeReviewer(ctx context.Context, prId string, oldReviewerId string) (*model.ReassignmentResult, error) {
	pullRequest, err := s.prRepository.GetPR(ctx, prId)
	if err != nil {
		return nil, err
	}

	if pullRequest.Status == model.MERGED {
		return nil, model.NewError(model.PR_MERGED, "PR already merged")
	}

	teammates, err := s.userRepository.GetActiveUsersByTeam(ctx, pullRequest.AuthorId)
	if err != nil {
		return nil, err
	}

	var newReviewerId string
	for _, userId := range teammates {
		if userId != oldReviewerId && userId != pullRequest.AuthorId {
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

	response := model.ReassignmentResult{
		PullRequest:   *pr,
		OldReviewerId: oldReviewerId,
	}
	return &response, nil
}

func (s *PullRequestService) MergePR(ctx context.Context, pullRequestId string) (*model.PullRequest, error) {
	createdPR, err := s.prRepository.MergePR(ctx, pullRequestId, model.MERGED, time.Now())

	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PullRequestService) GetPRsByUser(ctx context.Context, userId string) ([]*model.PullRequest, error) {
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
