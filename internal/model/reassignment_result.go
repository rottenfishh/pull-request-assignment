package model

type ReassignmentResult struct {
	PullRequest   PullRequest `json:"pull_request"`
	NewReviewerID string      `json:"new_reviewer_id"`
}
