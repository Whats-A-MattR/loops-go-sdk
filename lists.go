package loops

import (
	"context"
	"net/http"
)

// GetLists returns mailing lists (GET /lists per OpenAPI).
func (c *Client) GetLists(ctx context.Context) ([]MailingList, error) {
	var out []MailingList
	if err := c.do(ctx, http.MethodGet, "/lists", nil, &out, nil); err != nil {
		return nil, err
	}
	return out, nil
}
