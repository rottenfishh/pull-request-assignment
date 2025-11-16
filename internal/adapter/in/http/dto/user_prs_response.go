package dto

import "pr-assignment/internal/model"

type UserPrsResponse struct {
	UserID       string                   `json:"user_id"`
	PullRequests []model.PullRequestShort `json:"pull_requests"`
}
