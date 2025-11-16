package model

import (
	"errors"
)

type ErrorResponse struct {
	Error CustomError `json:"error"`
}

func ParseErrorResponse(err error) ErrorResponse {
	var response ErrorResponse

	var customErr *CustomError
	if errors.As(err, &customErr) {
		response.Error.Code = customErr.Code
		response.Error.Message = customErr.Message
	} else {
		response.Error.Code = InternalError
		response.Error.Message = err.Error()
	}
	return response
}
