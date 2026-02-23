package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

// CreateContactProperty creates a contact property (POST /contacts/properties). Name and type required per OpenAPI.
func (c *Client) CreateContactProperty(ctx context.Context, req *ContactPropertyCreateRequest) (*ContactPropertySuccessResponse, error) {
	if req == nil || req.Name == "" || req.Type == "" {
		return nil, &APIError{StatusCode: 400, Message: "name and type are required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out ContactPropertySuccessResponse
	if err := c.do(ctx, http.MethodPost, "/contacts/properties", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// ListContactProperties returns contact properties (GET /contacts/properties). List param: "all" or "custom" per OpenAPI.
func (c *Client) ListContactProperties(ctx context.Context, list string) ([]ContactProperty, error) {
	q := url.Values{}
	if list != "" {
		q.Set("list", list)
	}
	var out []ContactProperty
	if err := c.doWithQuery(ctx, http.MethodGet, "/contacts/properties", q, nil, &out); err != nil {
		return nil, err
	}
	return out, nil
}
