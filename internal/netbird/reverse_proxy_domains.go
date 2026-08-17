package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListReverseProxyDomainsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/reverse-proxies/domains", nil, &result); err != nil {
		return nil, fmt.Errorf("list reverse proxy domains: %w", err)
	}
	return result, nil
}

func (c *Client) CreateReverseProxyDomain(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/reverse-proxies/domains", request, &result); err != nil {
		return nil, fmt.Errorf("create reverse proxy domain: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteReverseProxyDomain(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("reverse proxy domain id is required")
	}
	path := "/api/reverse-proxies/domains/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete reverse proxy domain %q: %w", id, err)
	}
	return result, nil
}
