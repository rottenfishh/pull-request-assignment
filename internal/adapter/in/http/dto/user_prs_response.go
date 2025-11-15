package dto

import "pr-assignment/internal/model"

type UserPrsResponse struct {
	UserId       string                   `json:"user_id"`
	PullRequests []model.PullRequestShort `json:"pull_requests"`
}
