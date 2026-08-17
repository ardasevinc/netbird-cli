package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IdentityProvider struct {
	ID       *string `json:"id,omitempty"`
	Name     string  `json:"name"`
	Type     string  `json:"type"`
	Issuer   string  `json:"issuer"`
	ClientID string  `json:"client_id"`
}

type PostureCheck struct {
	ID          string                     `json:"id"`
	Name        string                     `json:"name"`
	Description *string                    `json:"description,omitempty"`
	Checks      map[string]json.RawMessage `json:"checks"`
}

func (c *Client) ListIdentityProviders(ctx context.Context) ([]IdentityProvider, error) {
	var result []IdentityProvider
	if err := c.transport.GetJSON(ctx, "/api/identity-providers", &result); err != nil {
		return nil, fmt.Errorf("list identity providers: %w", err)
	}
	return result, nil
}

func (c *Client) GetIdentityProvider(ctx context.Context, id string) (IdentityProvider, error) {
	if strings.TrimSpace(id) == "" {
		return IdentityProvider{}, fmt.Errorf("identity provider id is required")
	}
	var result IdentityProvider
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/identity-providers/"+url.PathEscape(id), nil, &result); err != nil {
		return IdentityProvider{}, fmt.Errorf("get identity provider %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListPostureChecks(ctx context.Context) ([]PostureCheck, error) {
	var result []PostureCheck
	if err := c.transport.GetJSON(ctx, "/api/posture-checks", &result); err != nil {
		return nil, fmt.Errorf("list posture checks: %w", err)
	}
	return result, nil
}

func (c *Client) ListPostureChecksRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/posture-checks", nil, &result); err != nil {
		return nil, fmt.Errorf("list posture checks: %w", err)
	}
	return result, nil
}

func (c *Client) CreatePostureCheck(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/posture-checks", request, &result); err != nil {
		return nil, fmt.Errorf("create posture check: %w", err)
	}
	return result, nil
}

func (c *Client) GetPostureCheck(ctx context.Context, id string) (PostureCheck, error) {
	if strings.TrimSpace(id) == "" {
		return PostureCheck{}, fmt.Errorf("posture check id is required")
	}
	var result PostureCheck
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/posture-checks/"+url.PathEscape(id), nil, &result); err != nil {
		return PostureCheck{}, fmt.Errorf("get posture check %q: %w", id, err)
	}
	return result, nil
}
