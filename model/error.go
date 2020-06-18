package model

import (
	"fmt"
)

type ErrorResponse struct {
	Error ApiError `json:"error"`
}

type ApiError struct {
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Status  string    `json:"status"`
	Details []interface{} `json:"details,omitempty"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("API error: %s (%d)", a.Message, a.Code)
}

