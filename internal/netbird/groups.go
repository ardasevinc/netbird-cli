package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Group struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
}

func (c *Client) ListGroups(ctx context.Context) ([]Group, error) {
	var groups []Group
	if err := c.transport.GetJSON(ctx, "/api/groups", &groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func (c *Client) ListGroupsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/groups", nil, &result); err != nil {
		return nil, fmt.Errorf("list groups: %w", err)
	}
	return result, nil
}

func (c *Client) CreateGroup(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/groups", request, &result); err != nil {
		return nil, fmt.Errorf("create group: %w", err)
	}
	return result, nil
}

// GetGroup returns the normalized JSON document so the mutation engine can
// compare the complete remote preimage without losing fields at the adapter
// boundary.
func (c *Client) GetGroup(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := groupPath(id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get group %q: %w", id, err)
	}
	return result, nil
}

// UpdateGroup sends the exact staged request once and returns the endpoint's
// response. Replay decisions belong to the mutation engine.
func (c *Client) UpdateGroup(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	path, err := groupPath(id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update group %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteGroup(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := groupPath(id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete group %q: %w", id, err)
	}
	return result, nil
}

func groupPath(id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("group id is required")
	}
	return "/api/groups/" + url.PathEscape(id), nil
}
