package service

import (
	"pr-assignment/internal/adapter/out/repository"
	"pr-assignment/internal/model"
	"time"
)

type PullRequestService struct {
	prRepository *repository.PullRequestRepository
}

func NewPullRequestService(prRepository *repository.PullRequestRepository) *PullRequestService {
	prRepository.Init()
	return &PullRequestService{prRepository: prRepository}
}

func (s *PullRequestService) Create(pullRequest model.PullRequestShort) error {
	pullRequest.Status = "CREATED"
	pr := model.PullRequest{
		PullRequestShort:  pullRequest,
		AssignedReviewers: make([]string, 0),
		CreatedAt:         time.Now(),
		MergedAt:          time.Time{},
	}
	err := s.prRepository.CreatePR(pr)
	if err != nil {
		return err
	}
	return nil
}

func (s *PullRequestService) mergePR(pullRequestId string) error {
	err := s.prRepository.MergePR(pullRequestId, "MERGED", time.Now())
	if err != nil {
		return err
	}
	return nil
}
