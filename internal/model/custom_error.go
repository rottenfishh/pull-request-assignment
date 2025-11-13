package model

import "fmt"

type ErrCode string

const (
	TEAM_EXISTS  ErrCode = "TEAM_EXISTS"
	PR_EXISTS    ErrCode = "PR_EXISTS"
	PR_MERGED    ErrCode = "PR_MERGED"
	NOT_ASSIGNED ErrCode = "NOT_ASSIGNED"
	NO_CANDIDATE ErrCode = "NO_CANDIDATE"
	NOT_FOUND    ErrCode = "NOT_FOUND"
)

type CustomError struct {
	message string
	code    ErrCode
	err     error
}

func (e *CustomError) Error() string {
	return fmt.Sprintf("%s : %s", e.code, e.message)
}

func NewError(code ErrCode, format string, a ...any) *CustomError {
	return &CustomError{code: code, message: fmt.Sprintf(format, a...)}
}
