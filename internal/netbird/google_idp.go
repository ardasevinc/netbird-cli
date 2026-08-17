package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListGoogleIDPsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/integrations/google-idp", nil, &result); err != nil {
		return nil, fmt.Errorf("list google IDP integrations: %w", err)
	}
	return result, nil
}

func (c *Client) GetGoogleIDPRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("google IDP integration id is required")
	}
	path := "/api/integrations/google-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get google IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateGoogleIDP(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/integrations/google-idp", request, &result); err != nil {
		return nil, fmt.Errorf("create google IDP integration: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateGoogleIDP(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("google IDP integration id is required")
	}
	path := "/api/integrations/google-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update google IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteGoogleIDP(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("google IDP integration id is required")
	}
	path := "/api/integrations/google-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete google IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) SyncGoogleIDP(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("google IDP integration id is required")
	}
	path := "/api/integrations/google-idp/" + url.PathEscape(id) + "/sync"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, nil, &result); err != nil {
		return nil, fmt.Errorf("sync google IDP integration %q: %w", id, err)
	}
	return result, nil
}
