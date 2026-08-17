package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListAzureIDPsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/integrations/azure-idp", nil, &result); err != nil {
		return nil, fmt.Errorf("list Azure IDP integrations: %w", err)
	}
	return result, nil
}

func (c *Client) GetAzureIDPRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("azure IDP integration id is required")
	}
	path := "/api/integrations/azure-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get azure IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateAzureIDP(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/integrations/azure-idp", request, &result); err != nil {
		return nil, fmt.Errorf("create azure IDP integration: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateAzureIDP(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("azure IDP integration id is required")
	}
	path := "/api/integrations/azure-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update azure IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteAzureIDP(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("azure IDP integration id is required")
	}
	path := "/api/integrations/azure-idp/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete azure IDP integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) SyncAzureIDP(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("azure IDP integration id is required")
	}
	path := "/api/integrations/azure-idp/" + url.PathEscape(id) + "/sync"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, nil, &result); err != nil {
		return nil, fmt.Errorf("sync azure IDP integration %q: %w", id, err)
	}
	return result, nil
}
