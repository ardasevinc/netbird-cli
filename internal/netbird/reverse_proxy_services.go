package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListReverseProxyClustersRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/reverse-proxies/clusters", nil, &result); err != nil {
		return nil, fmt.Errorf("list reverse proxy clusters: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteReverseProxyCluster(ctx context.Context, address string) (json.RawMessage, error) {
	if strings.TrimSpace(address) == "" {
		return nil, fmt.Errorf("reverse proxy cluster address is required")
	}
	path := "/api/reverse-proxies/clusters/" + url.PathEscape(address)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete reverse proxy cluster %q: %w", address, err)
	}
	return result, nil
}

func (c *Client) ListReverseProxyServicesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/reverse-proxies/services", nil, &result); err != nil {
		return nil, fmt.Errorf("list reverse proxy services: %w", err)
	}
	return result, nil
}

func (c *Client) GetReverseProxyServiceRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("reverse proxy service id is required")
	}
	path := "/api/reverse-proxies/services/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get reverse proxy service %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateReverseProxyService(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/reverse-proxies/services", request, &result); err != nil {
		return nil, fmt.Errorf("create reverse proxy service: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateReverseProxyService(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("reverse proxy service id is required")
	}
	path := "/api/reverse-proxies/services/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update reverse proxy service %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteReverseProxyService(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("reverse proxy service id is required")
	}
	path := "/api/reverse-proxies/services/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete reverse proxy service %q: %w", id, err)
	}
	return result, nil
}
