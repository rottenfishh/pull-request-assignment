package model

type ReassignmentResult struct {
	PullRequest   PullRequest `json:"pull_request"`
	OldReviewerId string      `json:"old_reviewer_id"`
}
