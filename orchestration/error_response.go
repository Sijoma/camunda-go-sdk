package orchestration

import (
	"encoding/json"
	"fmt"
	"io"
)

type ErrorResponse struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Status   int    `json:"status"`
	Detail   string `json:"detail"`
	Instance string `json:"instance"`
}

func (e ErrorResponse) Error() string {
	return fmt.Sprintf("Type: %s (Status: %d) - Detail: %s, Instance: %s", e.Type, e.Status, e.Detail, e.Instance)
}

// handleErrorResponse decodes and returns an error from an ErrorResponse
// The prefix parameter allows the caller to provide context for the error.
func handleErrorResponse(body io.ReadCloser, prefix string) error {
	var errorResponse ErrorResponse
	err := json.NewDecoder(body).Decode(&errorResponse)
	if err != nil {
		return fmt.Errorf("%s: failed to decode error response: %w", prefix, err)
	}

	return fmt.Errorf("%s: %w", prefix, errorResponse)
}
