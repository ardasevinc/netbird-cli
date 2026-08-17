package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func (c *Client) GetAgentNetworkSettingsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/agent-network/settings", nil, &result); err != nil {
		return nil, fmt.Errorf("get agent-network settings: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateAgentNetworkSettings(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, "/api/agent-network/settings", request, &result); err != nil {
		return nil, fmt.Errorf("update agent-network settings: %w", err)
	}
	return result, nil
}

func (c *Client) CreateAgentNetworkSettings(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/agent-network/settings", request, &result); err != nil {
		return nil, fmt.Errorf("create agent-network settings: %w", err)
	}
	return result, nil
}
