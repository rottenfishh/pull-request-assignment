package model

type ReassignmentResult struct {
	PullRequest   PullRequest `json:"pull_request"`
	NewReviewerId string      `json:"new_reviewer_id"`
}
