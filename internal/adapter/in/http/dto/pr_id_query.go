package dto

type PullRequestIDQuery struct {
	PrID string `form:"pull_request_id" json:"pull_request_id"`
}
