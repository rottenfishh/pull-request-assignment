package dto

import "pr-assignment/internal/model"

type UserResponse struct {
	model.User `json:"user"`
}
