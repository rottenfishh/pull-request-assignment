package model

import "fmt"

type ErrCode string

const (
	TEAM_EXISTS    ErrCode = "TEAM_EXISTS"
	PR_EXISTS      ErrCode = "PR_EXISTS"
	PR_MERGED      ErrCode = "PR_MERGED"
	NOT_ASSIGNED   ErrCode = "NOT_ASSIGNED"
	NO_CANDIDATE   ErrCode = "NO_CANDIDATE"
	NOT_FOUND      ErrCode = "NOT_FOUND"
	INTERNAL_ERROR ErrCode = "INTERNAL_ERROR"
)

type CustomError struct {
	Message string  `json:"code"`
	Code    ErrCode `json:"message"`
	err     error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s : %s", e.Code, e.Message)
}

func NewError(code ErrCode, format string, a ...any) *CustomError {
	return &CustomError{Code: code, Message: fmt.Sprintf(format, a...)}
}
