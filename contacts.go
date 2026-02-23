package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// CreateContact adds a contact (POST /contacts/create). Email is required per OpenAPI ContactRequest.
func (c *Client) CreateContact(ctx context.Context, req *ContactRequest) (*ContactSuccessResponse, error) {
	if req == nil || req.Email == "" {
		return nil, &APIError{StatusCode: 400, Message: "email is required"}
	}
	body, err := mergeBody(req, req.Extra)
	if err != nil {
		return nil, err
	}
	var out ContactSuccessResponse
	if err := c.do(ctx, http.MethodPost, "/contacts/create", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateContact updates a contact (PUT /contacts/update). Provide either email or userId per OpenAPI.
func (c *Client) UpdateContact(ctx context.Context, req *ContactUpdateRequest) (*ContactSuccessResponse, error) {
	if req == nil || (req.Email == "" && req.UserID == "") {
		return nil, &APIError{StatusCode: 400, Message: "email or userId is required"}
	}
	body, err := mergeBody(req, req.Extra)
	if err != nil {
		return nil, err
	}
	var out ContactSuccessResponse
	if err := c.do(ctx, http.MethodPut, "/contacts/update", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// FindContact finds a contact by email or userId (GET /contacts/find). Only one parameter allowed per OpenAPI.
func (c *Client) FindContact(ctx context.Context, email, userId string) ([]Contact, error) {
	if (email != "" && userId != "") || (email == "" && userId == "") {
		return nil, &APIError{StatusCode: 400, Message: "exactly one of email or userId is required"}
	}
	q := url.Values{}
	if email != "" {
		q.Set("email", email)
	} else {
		q.Set("userId", userId)
	}
	var out []Contact
	if err := c.doWithQuery(ctx, http.MethodGet, "/contacts/find", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// DeleteContact deletes a contact (POST /contacts/delete). Include only one of email or userId per OpenAPI.
func (c *Client) DeleteContact(ctx context.Context, req *ContactDeleteRequest) (*ContactDeleteResponse, error) {
	if req == nil {
		return nil, &APIError{StatusCode: 400, Message: "request is required"}
	}
	hasEmail := req.Email != ""
	hasUserID := req.UserID != ""
	if hasEmail == hasUserID {
		return nil, &APIError{StatusCode: 400, Message: "exactly one of email or userId is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out ContactDeleteResponse
	if err := c.do(ctx, http.MethodPost, "/contacts/delete", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
