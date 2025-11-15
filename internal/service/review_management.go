package service

import (
	"context"
	"pr-assignment/internal/model"
)

func (s *PullRequestService) checkAllowedToReview(reviewers []string, authorId string, newReviewerId string) (bool, error) {

	if authorId == newReviewerId {
		return false, nil
	}

	for id := range reviewers {
		if reviewers[id] == newReviewerId {
			return false, nil
		}
	}

	return true, nil
}

func (s *PullRequestService) AssignReviewers(ctx context.Context, pr *model.PullRequest) error {
	teammates, err := s.userService.GetActiveTeammatesByUser(ctx, pr.AuthorId)
	if err != nil {
		return err
	}

	for _, userId := range teammates {
		if userId != pr.AuthorId {
			pr.AssignedReviewers = append(pr.AssignedReviewers, userId)
			err = s.prReviewersRepository.AddReviewer(ctx, pr.PullRequestId, userId)
			if err != nil {
				return err
			}
			if len(pr.AssignedReviewers) == 2 {
				break
			}
		}
	}
	return nil
}

func (s *PullRequestService) inReviewers(reviewers []string, oldReviewerId string) bool {
	flag := false
	for id := range reviewers {
		if reviewers[id] == oldReviewerId {
			flag = true
			break
		}
	}
	return flag
}
