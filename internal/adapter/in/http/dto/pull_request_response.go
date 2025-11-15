package dto

import (
	"pr-assignment/internal/model"
	"time"
)

type PrResponse struct {
	model.PullRequestShort
	AssignedReviewers []string `json:"assigned_reviewers"`
}

type PrMerged struct {
	PrResponse
	MergedAt time.Time `json:"merged_at"`
}

type PrMergedResponse struct {
	PrMerged PrMerged `json:"pr"`
}

type PrReassignResponse struct {
	PrResponse `json:"pr"`
	ReplacedBy string `json:"replaced_by"`
}
