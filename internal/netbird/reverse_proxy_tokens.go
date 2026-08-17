package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListReverseProxyTokensRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/reverse-proxies/proxy-tokens", nil, &result); err != nil {
		return nil, fmt.Errorf("list reverse proxy tokens: %w", err)
	}
	return result, nil
}

func (c *Client) CreateReverseProxyToken(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/reverse-proxies/proxy-tokens", request, &result); err != nil {
		return nil, fmt.Errorf("create reverse proxy token: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteReverseProxyToken(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("reverse proxy token id is required")
	}
	path := "/api/reverse-proxies/proxy-tokens/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete reverse proxy token %q: %w", id, err)
	}
	return result, nil
}
