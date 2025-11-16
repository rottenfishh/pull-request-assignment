package service

import (
	"context"
	"fmt"
	"pr-assignment/internal/adapter/in/http/dto"
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
	"time"
)

type PullRequestService struct {
	prRepository          *repository.PullRequestRepository
	prReviewersRepository *repository.PrReviewersRepository
	teamRepository        *repository.TeamRepository
	userRepository        *repository.UserRepository
	userService           *UserService
}

func NewPullRequestService(prRepo *repository.PullRequestRepository, prReviewsRepo *repository.PrReviewersRepository,
	teamRepo *repository.TeamRepository, userRepo *repository.UserRepository, userService *UserService) *PullRequestService {

	return &PullRequestService{prRepo, prReviewsRepo,
		teamRepo, userRepo, userService}
}

func (s *PullRequestService) CreatePR(ctx context.Context, prBody dto.PullRequestQuery) (*model.PullRequest, error) {
	_, err := s.userRepository.GetUserByID(ctx, prBody.AuthorID)
	if err != nil {
		return nil, err
	}

	pullRequest := model.PullRequestShort{
		PullRequestID:   prBody.PullRequestID,
		PullRequestName: prBody.PullRequestName,
		AuthorID:        prBody.AuthorID,
		Status:          model.CREATED,
	}
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

func (s *PullRequestService) ChangeReviewer(ctx context.Context, prID string, oldReviewerID string) (*model.ReassignmentResult, error) {
	pullRequest, err := s.prRepository.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	if pullRequest.Status == model.MERGED {
		return nil, model.NewError(model.PrMerged, "PR already merged")
	}
	teamName, err := s.userRepository.GetTeamNameByUserID(ctx, pullRequest.AuthorID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	teammates, err := s.userRepository.GetActiveUsersByTeam(ctx, teamName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	reviewers, err := s.prReviewersRepository.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	if !s.inReviewers(reviewers, oldReviewerID) {
		return nil, model.NewError(model.NotAssigned, "Old reviewer was not assigned to PR")
	}

	var newReviewerID string
	newReviewerID = oldReviewerID
	for _, userID := range teammates {
		res, _ := s.checkAllowedToReview(reviewers, pullRequest.AuthorID, userID)

		if res {
			newReviewerID = userID
			break
		}
	}

	err = s.prReviewersRepository.ChangeReviewer(ctx, prID, oldReviewerID, newReviewerID)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pr, err := s.prRepository.GetPR(ctx, prID)
	if err != nil {
		return nil, err
	}

	reviewers, err = s.prReviewersRepository.GetReviewers(ctx, prID)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers

	response := model.ReassignmentResult{
		PullRequest:   *pr,
		NewReviewerID: newReviewerID,
	}
	return &response, nil
}

func (s *PullRequestService) MergePR(ctx context.Context, pullRequestID string) (*model.PullRequest, error) {
	createdPR, err := s.prRepository.MergePR(ctx, pullRequestID, model.MERGED, time.Now())

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	reviewers, err := s.prReviewersRepository.GetReviewers(ctx, pullRequestID)
	if err != nil {
		return nil, err
	}
	createdPR.AssignedReviewers = reviewers
	return createdPR, nil
}

func (s *PullRequestService) GetPRsByUser(ctx context.Context, userID string) ([]*model.PullRequest, error) {
	pullRequestsIDs, err := s.prReviewersRepository.GetPRsByUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	var prs []*model.PullRequest
	for _, prID := range pullRequestsIDs {
		pr, err := s.prRepository.GetPR(ctx, prID)
		if err != nil {
			return nil, err
		}
		prs = append(prs, pr)
	}
	return prs, nil
}

func (s *PullRequestService) ReassignReviewsAfterDeath(ctx context.Context, deadReviewerID string) error {
	pullRequestsIDs, err := s.prReviewersRepository.GetPRsByUser(ctx, deadReviewerID)
	if err != nil {
		return err
	}

	for _, prID := range pullRequestsIDs {
		_, err := s.ChangeReviewer(ctx, prID, deadReviewerID)
		if err != nil {
			return err
		}
	}
	return nil
}
