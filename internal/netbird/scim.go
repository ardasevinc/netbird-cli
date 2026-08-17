package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

var scimIntegrationProviders = map[string]string{
	"scim":      "scim-idp",
	"okta_scim": "okta-scim-idp",
}

func scimIntegrationPath(provider, id string) (string, error) {
	slug, ok := scimIntegrationProviders[strings.ToLower(strings.TrimSpace(provider))]
	if !ok {
		return "", fmt.Errorf("unsupported SCIM integration provider %q", provider)
	}
	path := "/api/integrations/" + slug
	if id != "" {
		path += "/" + url.PathEscape(id)
	}
	return path, nil
}

func (c *Client) ListSCIMIntegrationsRaw(ctx context.Context, provider string) (json.RawMessage, error) {
	path, err := scimIntegrationPath(provider, "")
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("list %s integrations: %w", provider, err)
	}
	return result, nil
}

func (c *Client) GetSCIMIntegrationRaw(ctx context.Context, provider, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("%s integration id is required", provider)
	}
	path, err := scimIntegrationPath(provider, id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get %s integration %q: %w", provider, id, err)
	}
	return result, nil
}

func (c *Client) CreateSCIMIntegration(ctx context.Context, provider string, request json.RawMessage) (json.RawMessage, error) {
	path, err := scimIntegrationPath(provider, "")
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("create %s integration: %w", provider, err)
	}
	return result, nil
}

func (c *Client) UpdateSCIMIntegration(ctx context.Context, provider, id string, request json.RawMessage) (json.RawMessage, error) {
	path, err := scimIntegrationPath(provider, id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update %s integration %q: %w", provider, id, err)
	}
	return result, nil
}

func (c *Client) DeleteSCIMIntegration(ctx context.Context, provider, id string) (json.RawMessage, error) {
	path, err := scimIntegrationPath(provider, id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete %s integration %q: %w", provider, id, err)
	}
	return result, nil
}

func (c *Client) RegenerateSCIMToken(ctx context.Context, provider, id string) (json.RawMessage, error) {
	path, err := scimIntegrationPath(provider, id)
	if err != nil {
		return nil, err
	}
	path += "/token"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, nil, &result); err != nil {
		return nil, fmt.Errorf("regenerate %s token %q: %w", provider, id, err)
	}
	return result, nil
}
