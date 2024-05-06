package dtos

import "time"

type ErrorResponse struct {
	Message   string    `json:"Message"`
	Timestamp time.Time `json:"Timestamp"`
}

func NewErrorResponse(err error) *ErrorResponse {
	return &ErrorResponse{Message: err.Error(), Timestamp: time.Now()}
}
