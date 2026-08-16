package netbird

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// SetupKey is the safe inventory view. Key is intentionally excluded because
// the management API may return one-time secret material in this response.
type SetupKey struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Type          string   `json:"type"`
	State         string   `json:"state"`
	Valid         bool     `json:"valid"`
	Revoked       bool     `json:"revoked"`
	Ephemeral     bool     `json:"ephemeral"`
	AllowExtraDNS bool     `json:"allow_extra_dns_labels"`
	AutoGroups    []string `json:"auto_groups"`
	UsageLimit    int      `json:"usage_limit"`
	UsedTimes     int      `json:"used_times"`
	Expires       string   `json:"expires"`
	LastUsed      string   `json:"last_used"`
	UpdatedAt     string   `json:"updated_at"`
}

func (c *Client) ListSetupKeys(ctx context.Context) ([]SetupKey, error) {
	var result []SetupKey
	if err := c.transport.GetJSON(ctx, "/api/setup-keys", &result); err != nil {
		return nil, fmt.Errorf("list setup keys: %w", err)
	}
	normalizeSetupKeys(result)
	return result, nil
}

func (c *Client) GetSetupKey(ctx context.Context, id string) (SetupKey, error) {
	if strings.TrimSpace(id) == "" {
		return SetupKey{}, fmt.Errorf("setup key id is required")
	}
	var result SetupKey
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/setup-keys/"+url.PathEscape(id), nil, &result); err != nil {
		return SetupKey{}, fmt.Errorf("get setup key %q: %w", id, err)
	}
	normalizeSetupKeys([]SetupKey{result})
	return result, nil
}

func normalizeSetupKeys(keys []SetupKey) {
	for i := range keys {
		if keys[i].AutoGroups == nil {
			keys[i].AutoGroups = []string{}
		}
	}
}
