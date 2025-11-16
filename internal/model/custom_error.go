package model

import "fmt"

type ErrCode string

const (
	DefaultError  ErrCode = "SOME_ERROR"
	TeamExists    ErrCode = "TEAM_EXISTS"
	PrExists      ErrCode = "PR_EXISTS"
	PrMerged      ErrCode = "PR_MERGED"
	NotAssigned   ErrCode = "NOT_ASSIGNED"
	NoCandidate   ErrCode = "NO_CANDIDATE"
	NotFound      ErrCode = "NOT_FOUND"
	InternalError ErrCode = "INTERNAL_ERROR"
)

type CustomError struct {
	Message string  `json:"message"`
	Code    ErrCode `json:"code"`
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s : %s", e.Code, e.Message)
}

func NewError(code ErrCode, format string, a ...any) *CustomError {
	return &CustomError{Code: code, Message: fmt.Sprintf(format, a...)}
}
