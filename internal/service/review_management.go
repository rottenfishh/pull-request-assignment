package service

import (
	"context"
	"pr-assignment/internal/model"
)

func (s *PullRequestService) checkAllowedToReview(reviewers []string, authorID string, newReviewerID string) (bool, error) {

	if authorID == newReviewerID {
		return false, nil
	}

	for id := range reviewers {
		if reviewers[id] == newReviewerID {
			return false, nil
		}
	}

	return true, nil
}

func (s *PullRequestService) AssignReviewers(ctx context.Context, pr *model.PullRequest) error {
	teammates, err := s.userService.GetActiveTeammatesByUser(ctx, pr.AuthorID)
	if err != nil {
		return err
	}

	for _, userID := range teammates {
		if userID != pr.AuthorID {
			pr.AssignedReviewers = append(pr.AssignedReviewers, userID)
			err = s.prReviewersRepository.AddReviewer(ctx, pr.PullRequestID, userID)
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

func (s *PullRequestService) inReviewers(reviewers []string, oldReviewerID string) bool {
	flag := false
	for id := range reviewers {
		if reviewers[id] == oldReviewerID {
			flag = true
			break
		}
	}
	return flag
}
