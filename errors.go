package loops

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// APIError represents an error response from the Loops API (success: false with message).
// It implements error and preserves the HTTP status code and raw body when available.
type APIError struct {
	StatusCode int
	Body       []byte
	Success    bool
	Message    string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("loops API error (status %d): %s", e.StatusCode, e.Message)
	}
	return fmt.Sprintf("loops API error (status %d): %s", e.StatusCode, string(e.Body))
}

// parseErrorBody attempts to parse a failure response body (ContactFailureResponse, EventFailureResponse, etc.).
func parseErrorBody(resp *http.Response, body []byte) *APIError {
	apiErr := &APIError{StatusCode: resp.StatusCode, Body: body}
	var generic struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &generic); err == nil {
		apiErr.Success = generic.Success
		apiErr.Message = generic.Message
	}
	return apiErr
}
