package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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

func (c *Client) DeleteAgentNetworkSettings(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, "/api/agent-network/settings", nil, &result); err != nil {
		return nil, fmt.Errorf("delete agent-network settings: %w", err)
	}
	return result, nil
}

func (c *Client) ListAgentNetworkBudgetRulesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/agent-network/budget-rules", nil, &result); err != nil {
		return nil, fmt.Errorf("list agent-network budget rules: %w", err)
	}
	return result, nil
}

func (c *Client) GetAgentNetworkBudgetRuleRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network budget rule id is required")
	}
	path := "/api/agent-network/budget-rules/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get agent-network budget rule %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateAgentNetworkBudgetRule(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/agent-network/budget-rules", request, &result); err != nil {
		return nil, fmt.Errorf("create agent-network budget rule: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateAgentNetworkBudgetRule(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network budget rule id is required")
	}
	path := "/api/agent-network/budget-rules/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update agent-network budget rule %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteAgentNetworkBudgetRule(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network budget rule id is required")
	}
	path := "/api/agent-network/budget-rules/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete agent-network budget rule %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListAgentNetworkGuardrailsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/agent-network/guardrails", nil, &result); err != nil {
		return nil, fmt.Errorf("list agent-network guardrails: %w", err)
	}
	return result, nil
}

func (c *Client) GetAgentNetworkGuardrailRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network guardrail id is required")
	}
	path := "/api/agent-network/guardrails/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get agent-network guardrail %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateAgentNetworkGuardrail(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/agent-network/guardrails", request, &result); err != nil {
		return nil, fmt.Errorf("create agent-network guardrail: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateAgentNetworkGuardrail(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network guardrail id is required")
	}
	path := "/api/agent-network/guardrails/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update agent-network guardrail %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteAgentNetworkGuardrail(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network guardrail id is required")
	}
	path := "/api/agent-network/guardrails/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete agent-network guardrail %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListAgentNetworkPoliciesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/agent-network/policies", nil, &result); err != nil {
		return nil, fmt.Errorf("list agent-network policies: %w", err)
	}
	return result, nil
}

func (c *Client) GetAgentNetworkPolicyRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network policy id is required")
	}
	path := "/api/agent-network/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get agent-network policy %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateAgentNetworkPolicy(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/agent-network/policies", request, &result); err != nil {
		return nil, fmt.Errorf("create agent-network policy: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateAgentNetworkPolicy(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network policy id is required")
	}
	path := "/api/agent-network/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update agent-network policy %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteAgentNetworkPolicy(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("agent-network policy id is required")
	}
	path := "/api/agent-network/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete agent-network policy %q: %w", id, err)
	}
	return result, nil
}
