package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type NetworkResource struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Address     string      `json:"address"`
	Type        string      `json:"type"`
	Description *string     `json:"description,omitempty"`
	Enabled     bool        `json:"enabled"`
	Groups      []PeerGroup `json:"groups"`
}

type NetworkRouter struct {
	ID         string   `json:"id"`
	Enabled    bool     `json:"enabled"`
	Masquerade bool     `json:"masquerade"`
	Metric     int      `json:"metric"`
	Peer       *string  `json:"peer,omitempty"`
	PeerGroups []string `json:"peer_groups"`
}

func (c *Client) ListNetworkResources(ctx context.Context, networkID string) ([]NetworkResource, error) {
	path, err := networkChildPath(networkID, "resources", "")
	if err != nil {
		return nil, err
	}
	var result []NetworkResource
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list resources for network %q: %w", networkID, err)
	}
	if result == nil {
		result = []NetworkResource{}
	}
	return result, nil
}

func (c *Client) GetNetworkResource(ctx context.Context, networkID, resourceID string) (NetworkResource, error) {
	path, err := networkChildPath(networkID, "resources", resourceID)
	if err != nil {
		return NetworkResource{}, err
	}
	var result NetworkResource
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return NetworkResource{}, fmt.Errorf("get resource %q for network %q: %w", resourceID, networkID, err)
	}
	return result, nil
}

func (c *Client) GetNetworkResourceRaw(ctx context.Context, networkID, resourceID string) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "resources", resourceID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get resource %q for network %q: %w", resourceID, networkID, err)
	}
	return result, nil
}

func (c *Client) UpdateNetworkResource(ctx context.Context, networkID, resourceID string, request json.RawMessage) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "resources", resourceID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update resource %q for network %q: %w", resourceID, networkID, err)
	}
	return result, nil
}

func (c *Client) DeleteNetworkResource(ctx context.Context, networkID, resourceID string) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "resources", resourceID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete resource %q for network %q: %w", resourceID, networkID, err)
	}
	return result, nil
}

func (c *Client) ListNetworkRouters(ctx context.Context, networkID string) ([]NetworkRouter, error) {
	path, err := networkChildPath(networkID, "routers", "")
	if err != nil {
		return nil, err
	}
	var result []NetworkRouter
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list routers for network %q: %w", networkID, err)
	}
	if result == nil {
		result = []NetworkRouter{}
	}
	for i := range result {
		if result[i].PeerGroups == nil {
			result[i].PeerGroups = []string{}
		}
	}
	return result, nil
}

func (c *Client) GetNetworkRouter(ctx context.Context, networkID, routerID string) (NetworkRouter, error) {
	path, err := networkChildPath(networkID, "routers", routerID)
	if err != nil {
		return NetworkRouter{}, err
	}
	var result NetworkRouter
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return NetworkRouter{}, fmt.Errorf("get router %q for network %q: %w", routerID, networkID, err)
	}
	if result.PeerGroups == nil {
		result.PeerGroups = []string{}
	}
	return result, nil
}

func (c *Client) GetNetworkRouterRaw(ctx context.Context, networkID, routerID string) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "routers", routerID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get router %q for network %q: %w", routerID, networkID, err)
	}
	return result, nil
}

func (c *Client) UpdateNetworkRouter(ctx context.Context, networkID, routerID string, request json.RawMessage) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "routers", routerID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update router %q for network %q: %w", routerID, networkID, err)
	}
	return result, nil
}

func (c *Client) DeleteNetworkRouter(ctx context.Context, networkID, routerID string) (json.RawMessage, error) {
	path, err := networkChildPath(networkID, "routers", routerID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete router %q for network %q: %w", routerID, networkID, err)
	}
	return result, nil
}

func (c *Client) ListAllNetworkRouters(ctx context.Context) ([]NetworkRouter, error) {
	var result []NetworkRouter
	if err := c.transport.GetJSON(ctx, "/api/networks/routers", &result); err != nil {
		return nil, fmt.Errorf("list all network routers: %w", err)
	}
	if result == nil {
		result = []NetworkRouter{}
	}
	for i := range result {
		if result[i].PeerGroups == nil {
			result[i].PeerGroups = []string{}
		}
	}
	return result, nil
}

func networkChildPath(networkID, resource, childID string) (string, error) {
	if strings.TrimSpace(networkID) == "" {
		return "", fmt.Errorf("network id is required")
	}
	path := "/api/networks/" + url.PathEscape(networkID) + "/" + resource
	if childID != "" {
		path += "/" + url.PathEscape(childID)
	}
	return path, nil
}
