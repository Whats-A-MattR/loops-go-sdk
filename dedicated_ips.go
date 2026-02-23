package loops

import (
	"context"
	"net/http"
)

// GetDedicatedSendingIPs returns dedicated sending IP addresses (GET /dedicated-sending-ips per OpenAPI).
func (c *Client) GetDedicatedSendingIPs(ctx context.Context) ([]string, error) {
	var out []string
	if err := c.do(ctx, http.MethodGet, "/dedicated-sending-ips", nil, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}
