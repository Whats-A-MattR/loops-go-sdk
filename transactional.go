package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// SendTransactional sends a transactional email (POST /transactional). Email and transactionalId required per OpenAPI.
// IdempotencyKey is optional (max 100 chars).
func (c *Client) SendTransactional(ctx context.Context, req *TransactionalRequest, idempotencyKey string) (*TransactionalSuccessResponse, error) {
	if req == nil || req.Email == "" || req.TransactionalID == "" {
		return nil, &APIError{StatusCode: 400, Message: "email and transactionalId are required"}
	}
	body, err := json.Marshal(req)
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
	var out TransactionalSuccessResponse
	if err := c.doWithHeaders(ctx, http.MethodPost, "/transactional", headers, body, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListTransactionals returns published transactional emails (GET /transactional). perPage 10–50, default 20; cursor optional per OpenAPI.
func (c *Client) ListTransactionals(ctx context.Context, perPage int, cursor string) (*ListTransactionalsResponse, error) {
	q := url.Values{}
	if perPage > 0 {
		q.Set("perPage", strconv.Itoa(perPage))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var out ListTransactionalsResponse
	if err := c.doWithQuery(ctx, http.MethodGet, "/transactional", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
