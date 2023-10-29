package dtos

import "time"

type ErrorResponse struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error(), Timestamp: time.Now()}
}
