package service

import (
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
	"time"
)

type PullRequestService struct {
	prRepository   *repository.PullRequestRepository
	teamRepository *repository.TeamRepository
	userRepository *repository.UserRepository
}

func NewPullRequestService(prRepository *repository.PullRequestRepository) *PullRequestService {
	prRepository.Init()
	return &PullRequestService{prRepository: prRepository}
}

func (s *PullRequestService) Create(pullRequest model.PullRequestShort) (*model.PullRequest, error) {
	pullRequest.Status = "CREATED"
	pr := model.PullRequest{
		PullRequestShort:  pullRequest,
		AssignedReviewers: make([]string, 0),
		CreatedAt:         time.Now(),
		MergedAt:          time.Time{},
	}
	createdPR, err := s.prRepository.CreatePR(pr)
	if err != nil {
		return nil, err
	}

	err = s.assignReviewers(createdPR)
	if err != nil {
		return nil, err
	}

	return createdPR, nil
}

func (s *PullRequestService) assignReviewers(pr *model.PullRequest) error {
	teamId, err := s.userRepository.GetTeamNameByUserId(pr.AuthorId)
	if err != nil {
		return err
	}

	userIds, err := s.teamRepository.GetUsersIdByTeam(teamId)
	if err != nil {
		return err
	}

	for _, userId := range userIds {
		if userId != pr.AuthorId {
			pr.AssignedReviewers = append(pr.AssignedReviewers, userId)
		}
		if len(pr.AssignedReviewers) == 2 {
			break
		}
	}
	return nil
}

func (s *PullRequestService) mergePR(pullRequestId string) (*model.PullRequest, error) {
	createdPR, err := s.prRepository.MergePR(pullRequestId, "MERGED", time.Now())

	if err != nil {
		return nil, err
	}

	return createdPR, nil
}
