package model

type PrReviewersCount struct {
	PullRequest PullRequestShort `json:"pull_request"`
	Count       int              `json:"reviewers_count"`
}
