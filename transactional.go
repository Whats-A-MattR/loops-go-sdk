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

// ListTransactionalResources returns a paginated list of transactional resources (GET /transactionals). perPage 10–50, default 20; cursor optional per OpenAPI.
func (c *Client) ListTransactionalResources(ctx context.Context, perPage int, cursor string) (*ListTransactionalsResourceResponse, error) {
	q := url.Values{}
	if perPage > 0 {
		q.Set("perPage", strconv.Itoa(perPage))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var out ListTransactionalsResourceResponse
	if err := c.doWithQuery(ctx, http.MethodGet, "/transactionals", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateTransactional creates a new transactional (POST /transactionals). name is required per OpenAPI.
func (c *Client) CreateTransactional(ctx context.Context, req *CreateTransactionalRequest) (*TransactionalResourceResponse, error) {
	if req == nil || req.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "name is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out TransactionalResourceResponse
	if err := c.do(ctx, http.MethodPost, "/transactionals", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTransactional retrieves a transactional by ID (GET /transactionals/{transactionalId}).
func (c *Client) GetTransactional(ctx context.Context, transactionalID string) (*TransactionalResourceResponse, error) {
	if transactionalID == "" {
		return nil, &APIError{StatusCode: 400, Message: "transactionalId is required"}
	}
	var out TransactionalResourceResponse
	if err := c.do(ctx, http.MethodGet, "/transactionals/"+url.PathEscape(transactionalID), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateTransactionalResource updates a transactional's name (POST /transactionals/{transactionalId}).
func (c *Client) UpdateTransactionalResource(ctx context.Context, transactionalID string, req *UpdateTransactionalRequest) (*TransactionalResourceResponse, error) {
	if transactionalID == "" {
		return nil, &APIError{StatusCode: 400, Message: "transactionalId is required"}
	}
	if req == nil || req.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "name is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out TransactionalResourceResponse
	if err := c.do(ctx, http.MethodPost, "/transactionals/"+url.PathEscape(transactionalID), body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// EnsureTransactionalDraft ensures a draft email message exists for the transactional (POST /transactionals/{transactionalId}/draft).
func (c *Client) EnsureTransactionalDraft(ctx context.Context, transactionalID string) (*TransactionalDraftResponse, error) {
	if transactionalID == "" {
		return nil, &APIError{StatusCode: 400, Message: "transactionalId is required"}
	}
	var out TransactionalDraftResponse
	if err := c.do(ctx, http.MethodPost, "/transactionals/"+url.PathEscape(transactionalID)+"/draft", nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// PublishTransactional publishes a transactional's current draft (POST /transactionals/{transactionalId}/publish).
func (c *Client) PublishTransactional(ctx context.Context, transactionalID string) (*TransactionalResourceResponse, error) {
	if transactionalID == "" {
		return nil, &APIError{StatusCode: 400, Message: "transactionalId is required"}
	}
	var out TransactionalResourceResponse
	if err := c.do(ctx, http.MethodPost, "/transactionals/"+url.PathEscape(transactionalID)+"/publish", nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
