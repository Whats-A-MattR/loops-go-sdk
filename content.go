package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// GetEmailMessage retrieves an email message by ID (GET /email-messages/{emailMessageId}).
func (c *Client) GetEmailMessage(ctx context.Context, emailMessageID string) (*EmailMessageResponse, error) {
	if emailMessageID == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailMessageId is required"}
	}
	var out EmailMessageResponse
	if err := c.do(ctx, http.MethodGet, "/email-messages/"+url.PathEscape(emailMessageID), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateEmailMessage updates an email message (POST /email-messages/{emailMessageId}).
func (c *Client) UpdateEmailMessage(ctx context.Context, emailMessageID string, req *UpdateEmailMessageRequest) (*EmailMessageResponse, error) {
	if emailMessageID == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailMessageId is required"}
	}
	if req == nil {
		return nil, &APIError{StatusCode: 400, Message: "request is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out EmailMessageResponse
	if err := c.do(ctx, http.MethodPost, "/email-messages/"+url.PathEscape(emailMessageID), body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListThemes returns themes (GET /themes). perPage 10-50, default 20; cursor optional per OpenAPI.
func (c *Client) ListThemes(ctx context.Context, perPage int, cursor string) (*ListThemesResponse, error) {
	q := url.Values{}
	if perPage > 0 {
		q.Set("perPage", strconv.Itoa(perPage))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var out ListThemesResponse
	if err := c.doWithQuery(ctx, http.MethodGet, "/themes", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetTheme retrieves a theme by ID (GET /themes/{themeId}).
func (c *Client) GetTheme(ctx context.Context, themeID string) (*ThemeResponse, error) {
	if themeID == "" {
		return nil, &APIError{StatusCode: 400, Message: "themeId is required"}
	}
	var out ThemeResponse
	if err := c.do(ctx, http.MethodGet, "/themes/"+url.PathEscape(themeID), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListComponents returns components (GET /components). perPage 10-50, default 20; cursor optional per OpenAPI.
func (c *Client) ListComponents(ctx context.Context, perPage int, cursor string) (*ListComponentsResponse, error) {
	q := url.Values{}
	if perPage > 0 {
		q.Set("perPage", strconv.Itoa(perPage))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var out ListComponentsResponse
	if err := c.doWithQuery(ctx, http.MethodGet, "/components", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetComponent retrieves a component by ID (GET /components/{componentId}).
func (c *Client) GetComponent(ctx context.Context, componentID string) (*ComponentResponse, error) {
	if componentID == "" {
		return nil, &APIError{StatusCode: 400, Message: "componentId is required"}
	}
	var out ComponentResponse
	if err := c.do(ctx, http.MethodGet, "/components/"+url.PathEscape(componentID), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateUpload requests a pre-signed URL for an upload (POST /uploads).
func (c *Client) CreateUpload(ctx context.Context, req *CreateUploadRequest) (*CreateUploadResponse, error) {
	if req == nil {
		return nil, &APIError{StatusCode: 400, Message: "request is required"}
	}
	if req.EmailMessageID == "" {
		return nil, &APIError{StatusCode: 400, Message: "emailMessageId is required"}
	}
	if req.ContentType == "" {
		return nil, &APIError{StatusCode: 400, Message: "contentType is required"}
	}
	if req.ContentLength <= 0 {
		return nil, &APIError{StatusCode: 400, Message: "contentLength must be greater than 0"}
	}
	if req.ContentLength > 4_000_000 {
		return nil, &APIError{StatusCode: 400, Message: "contentLength must not exceed 4000000 bytes"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out CreateUploadResponse
	if err := c.do(ctx, http.MethodPost, "/uploads", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// CompleteUpload finalizes an upload after file transfer (POST /uploads/{id}/complete).
func (c *Client) CompleteUpload(ctx context.Context, id string) (*CompleteUploadResponse, error) {
	if id == "" {
		return nil, &APIError{StatusCode: 400, Message: "id is required"}
	}
	var out CompleteUploadResponse
	if err := c.do(ctx, http.MethodPost, "/uploads/"+url.PathEscape(id)+"/complete", nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
