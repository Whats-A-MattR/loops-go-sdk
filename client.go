package loops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Client is the Loops API client. All methods are safe for concurrent use.
type Client struct {
	apiKey  string
	baseURL string
	client  *http.Client
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithBaseURL sets the API base URL (default: DefaultBaseURL from OpenAPI spec).
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = strings.TrimSuffix(baseURL, "/")
	}
}

// WithHTTPClient sets the *http.Client used for requests (default: http.DefaultClient).
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.client = client
	}
}

// NewClient returns a new Loops API client. API key is required (Bearer auth per OpenAPI securitySchemes).
func NewClient(apiKey string, opts ...ClientOption) *Client {
	c := &Client{
		apiKey:  apiKey,
		baseURL: DefaultBaseURL,
		client:  http.DefaultClient,
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// mergeBody marshals v to JSON, then merges extra keys (for OpenAPI additionalProperties) and returns the final JSON.
func mergeBody(v interface{}, extra map[string]interface{}) ([]byte, error) {
	if len(extra) == 0 {
		return json.Marshal(v)
	}
	data, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var m map[string]interface{}
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, err
	}
	for k, val := range extra {
		m[k] = val
	}
	return json.Marshal(m)
}

func (c *Client) do(ctx context.Context, method, path string, body []byte, result interface{}, opts *doOpts) error {
	var bodyReader io.Reader
	if len(body) > 0 {
		bodyReader = bytes.NewReader(body)
	}
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")
	if opts != nil {
		for k, v := range opts.headers {
			req.Header.Set(k, v)
		}
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	slurp, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return parseErrorBody(resp, slurp)
	}
	if result != nil && len(slurp) > 0 {
		if err := json.Unmarshal(slurp, result); err != nil {
			return fmt.Errorf("decode response: %w", err)
		}
	}
	return nil
}

type doOpts struct {
	headers map[string]string
	query   url.Values
}

func (c *Client) doWithQuery(ctx context.Context, method, path string, query url.Values, body []byte, result interface{}) error {
	if len(query) > 0 {
		path = path + "?" + query.Encode()
	}
	return c.do(ctx, method, path, body, result, nil)
}

func (c *Client) doWithHeaders(ctx context.Context, method, path string, headers map[string]string, body []byte, result interface{}) error {
	return c.do(ctx, method, path, body, result, &doOpts{headers: headers})
}
