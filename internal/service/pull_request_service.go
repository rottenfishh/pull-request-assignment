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
	_, err := s.userRepository.GetUserById(ctx, prBody.AuthorId)
	if err != nil {
		return nil, err
	}

	pullRequest := model.PullRequestShort{
		PullRequestId:   prBody.PullRequestId,
		PullRequestName: prBody.PullRequestName,
		AuthorId:        prBody.AuthorId,
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

func (s *PullRequestService) ChangeReviewer(ctx context.Context, prId string, oldReviewerId string) (*model.ReassignmentResult, error) {
	pullRequest, err := s.prRepository.GetPR(ctx, prId)
	if err != nil {
		return nil, err
	}

	if pullRequest.Status == model.MERGED {
		return nil, model.NewError(model.PR_MERGED, "PR already merged")
	}
	teamName, err := s.userRepository.GetTeamNameByUserId(ctx, pullRequest.AuthorId)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	teammates, err := s.userRepository.GetActiveUsersByTeam(ctx, teamName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	reviewers, err := s.prReviewersRepository.GetReviewers(ctx, prId)
	fmt.Println(reviewers)

	if !s.inReviewers(reviewers, oldReviewerId) {
		return nil, model.NewError(model.NOT_ASSIGNED, "Old reviewer was not assigned to PR")
	}

	var newReviewerId string
	newReviewerId = oldReviewerId
	for _, userId := range teammates {
		res, _ := s.checkAllowedToReview(reviewers, pullRequest.AuthorId, userId)

		if res {
			newReviewerId = userId
			break
		}
	}

	err = s.prReviewersRepository.ChangeReviewer(ctx, prId, oldReviewerId, newReviewerId)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	pr, err := s.prRepository.GetPR(ctx, prId)
	if err != nil {
		return nil, err
	}

	reviewers, err = s.prReviewersRepository.GetReviewers(ctx, prId)
	if err != nil {
		return nil, err
	}

	pr.AssignedReviewers = reviewers

	response := model.ReassignmentResult{
		PullRequest:   *pr,
		NewReviewerId: newReviewerId,
	}
	return &response, nil
}

func (s *PullRequestService) MergePR(ctx context.Context, pullRequestId string) (*model.PullRequest, error) {
	createdPR, err := s.prRepository.MergePR(ctx, pullRequestId, model.MERGED, time.Now())

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	reviewers, err := s.prReviewersRepository.GetReviewers(ctx, pullRequestId)
	fmt.Println(reviewers)
	if err != nil {
		return nil, err
	}
	createdPR.AssignedReviewers = reviewers
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

func (s *PullRequestService) ReassignReviewsAfterDeath(ctx context.Context, deadReviewerId string) error {
	pullRequestsIds, err := s.prReviewersRepository.GetPRsByUser(ctx, deadReviewerId)
	if err != nil {
		return err
	}
	for _, prId := range pullRequestsIds {
		_, err := s.ChangeReviewer(ctx, prId, deadReviewerId)
		if err != nil {
			return err
		}
	}
	return nil
}
