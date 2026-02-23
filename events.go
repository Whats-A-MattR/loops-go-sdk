package loops

import (
	"context"
	"net/http"
)

const idempotencyKeyHeader = "Idempotency-Key"

// SendEvent sends an event (POST /events/send). EventName required; provide email or userId per OpenAPI.
// IdempotencyKey is optional (max 100 chars per OpenAPI).
func (c *Client) SendEvent(ctx context.Context, req *EventRequest, idempotencyKey string) (*EventSuccessResponse, error) {
	if req == nil || req.EventName == "" {
		return nil, &APIError{StatusCode: 400, Message: "eventName is required"}
	}
	if req.Email == "" && req.UserID == "" {
		return nil, &APIError{StatusCode: 400, Message: "email or userId is required"}
	}
	body, err := mergeBody(req, req.Extra)
	if err != nil {
		return nil, err
	}
	headers := make(map[string]string)
	if len(idempotencyKey) > 0 {
		if len(idempotencyKey) > 100 {
			idempotencyKey = idempotencyKey[:100]
		}
		headers[idempotencyKeyHeader] = idempotencyKey
	}
	var out EventSuccessResponse
	if err := c.doWithHeaders(ctx, http.MethodPost, "/events/send", headers, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
