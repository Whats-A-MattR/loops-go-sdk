package loops

import (
	"context"
	"net/http"
)

// GetAPIKey tests the API key (GET /api-key per OpenAPI). Returns team name on success.
func (c *Client) GetAPIKey(ctx context.Context) (*APIKeyResponse, error) {
	var out APIKeyResponse
	err := c.do(ctx, http.MethodGet, "/api-key", nil, &out, nil)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
