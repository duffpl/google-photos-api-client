package internal

import (
	"errors"
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Error ApiError `json:"error"`
}

type ApiError struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Status  string        `json:"status"`
	Details []interface{} `json:"details,omitempty"`
}

func (a ApiError) Error() string {
	return fmt.Sprintf("API common: %s (%d)", a.Message, a.Code)
}

type APIStatus struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Details []interface{} `json:"details"`
}

func GetErrorFromResponse(res *http.Response) error {
	if res.StatusCode < 400 {
		return nil
	}
	if res.StatusCode == 404 {
		return errors.New("url not found")
	}
	responseModel := &ErrorResponse{}
	err := UnmarshalResponse(res, responseModel)
	if err != nil {
		return fmt.Errorf("cannot unmarshal error response: %w", err)
	}
	return responseModel.Error
}
