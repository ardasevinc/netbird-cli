package netbird

import (
	"context"
	"encoding/json"
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

func (c *Client) ListSetupKeysRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/setup-keys", nil, &result); err != nil {
		return nil, fmt.Errorf("list setup keys: %w", err)
	}
	return result, nil
}

func (c *Client) CreateSetupKey(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/setup-keys", request, &result); err != nil {
		return nil, fmt.Errorf("create setup key: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateSetupKey(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("setup key id is required")
	}
	path := "/api/setup-keys/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update setup key %q: %w", id, err)
	}
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

func (c *Client) GetSetupKeyRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("setup key id is required")
	}
	path := "/api/setup-keys/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get setup key %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteSetupKey(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("setup key id is required")
	}
	path := "/api/setup-keys/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete setup key %q: %w", id, err)
	}
	return result, nil
}

func normalizeSetupKeys(keys []SetupKey) {
	for i := range keys {
		if keys[i].AutoGroups == nil {
			keys[i].AutoGroups = []string{}
		}
	}
}
