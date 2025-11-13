package model

import "time"

type PullRequest struct {
	PullRequestShort
	AssignedReviewers []string  `json:"assigned_reviewers"`
	CreatedAt         time.Time `json:"createdAt"`
	MergedAt          time.Time `json:"mergedAt"`
}
