package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type PolicyGroup struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count,omitempty"`
	ResourcesCount int    `json:"resources_count,omitempty"`
}

type PolicyResource struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type PolicyPortRange struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type PolicyRule struct {
	ID                  *string             `json:"id,omitempty"`
	Name                string              `json:"name"`
	Description         *string             `json:"description,omitempty"`
	Action              string              `json:"action"`
	Protocol            string              `json:"protocol"`
	Enabled             bool                `json:"enabled"`
	Bidirectional       bool                `json:"bidirectional"`
	Sources             []PolicyGroup       `json:"sources,omitempty"`
	Destinations        []PolicyGroup       `json:"destinations,omitempty"`
	Ports               []string            `json:"ports,omitempty"`
	PortRanges          []PolicyPortRange   `json:"port_ranges,omitempty"`
	AuthorizedGroups    map[string][]string `json:"authorized_groups,omitempty"`
	SourceResource      *PolicyResource     `json:"sourceResource,omitempty"`
	DestinationResource *PolicyResource     `json:"destinationResource,omitempty"`
}

type Policy struct {
	ID                  *string      `json:"id,omitempty"`
	Name                string       `json:"name"`
	Description         *string      `json:"description,omitempty"`
	Enabled             bool         `json:"enabled"`
	Rules               []PolicyRule `json:"rules"`
	SourcePostureChecks []string     `json:"source_posture_checks"`
}

func (c *Client) ListPolicies(ctx context.Context) ([]Policy, error) {
	var result []Policy
	if err := c.transport.GetJSON(ctx, "/api/policies", &result); err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	return result, nil
}

func (c *Client) ListPoliciesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/policies", nil, &result); err != nil {
		return nil, fmt.Errorf("list policies: %w", err)
	}
	return result, nil
}

func (c *Client) CreatePolicy(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/policies", request, &result); err != nil {
		return nil, fmt.Errorf("create policy: %w", err)
	}
	return result, nil
}

func (c *Client) GetPolicy(ctx context.Context, id string) (Policy, error) {
	if strings.TrimSpace(id) == "" {
		return Policy{}, fmt.Errorf("policy id is required")
	}
	path := "/api/policies/" + url.PathEscape(id)
	var result Policy
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Policy{}, fmt.Errorf("get policy %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) GetPolicyRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("policy id is required")
	}
	path := "/api/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get policy %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdatePolicy(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("policy id is required")
	}
	path := "/api/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update policy %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeletePolicy(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("policy id is required")
	}
	path := "/api/policies/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete policy %q: %w", id, err)
	}
	return result, nil
}
