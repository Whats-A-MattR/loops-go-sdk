package loops

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strconv"
)

// ListCampaigns returns campaigns (GET /campaigns). perPage 10-50, default 20; cursor optional per OpenAPI.
func (c *Client) ListCampaigns(ctx context.Context, perPage int, cursor string) (*ListCampaignsResponse, error) {
	q := url.Values{}
	if perPage > 0 {
		q.Set("perPage", strconv.Itoa(perPage))
	}
	if cursor != "" {
		q.Set("cursor", cursor)
	}
	var out ListCampaignsResponse
	if err := c.doWithQuery(ctx, http.MethodGet, "/campaigns", q, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// CreateCampaign creates a draft campaign (POST /campaigns). Name is required per OpenAPI.
func (c *Client) CreateCampaign(ctx context.Context, req *CreateCampaignRequest) (*CreateCampaignResponse, error) {
	if req == nil || req.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "name is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out CreateCampaignResponse
	if err := c.do(ctx, http.MethodPost, "/campaigns", body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// GetCampaign retrieves a campaign by ID (GET /campaigns/{campaignId}).
func (c *Client) GetCampaign(ctx context.Context, campaignID string) (*CampaignResponse, error) {
	if campaignID == "" {
		return nil, &APIError{StatusCode: 400, Message: "campaignId is required"}
	}
	var out CampaignResponse
	if err := c.do(ctx, http.MethodGet, "/campaigns/"+url.PathEscape(campaignID), nil, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}

// UpdateCampaign updates a draft campaign (POST /campaigns/{campaignId}). Campaign ID and name are required per OpenAPI.
func (c *Client) UpdateCampaign(ctx context.Context, campaignID string, req *UpdateCampaignRequest) (*CampaignResponse, error) {
	if campaignID == "" {
		return nil, &APIError{StatusCode: 400, Message: "campaignId is required"}
	}
	if req == nil || req.Name == "" {
		return nil, &APIError{StatusCode: 400, Message: "name is required"}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	var out CampaignResponse
	if err := c.do(ctx, http.MethodPost, "/campaigns/"+url.PathEscape(campaignID), body, &out, nil); err != nil {
		return nil, err
	}
	return &out, nil
}
