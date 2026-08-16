package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// Route is the stable topology view used by reachability analysis. Provider
// fields are retained only when they have a direct management-plane meaning.
type Route struct {
	ID                  string   `json:"id"`
	Description         string   `json:"description"`
	Domains             []string `json:"domains,omitempty"`
	Enabled             bool     `json:"enabled"`
	Groups              []string `json:"groups"`
	KeepRoute           bool     `json:"keep_route"`
	Masquerade          bool     `json:"masquerade"`
	Metric              int      `json:"metric"`
	Network             *string  `json:"network,omitempty"`
	NetworkID           string   `json:"network_id"`
	NetworkType         string   `json:"network_type"`
	Peer                *string  `json:"peer,omitempty"`
	PeerGroups          []string `json:"peer_groups,omitempty"`
	AccessControlGroups []string `json:"access_control_groups,omitempty"`
	SkipAutoApply       *bool    `json:"skip_auto_apply,omitempty"`
}

type routeWire struct {
	ID                  string    `json:"id"`
	Description         string    `json:"description"`
	Domains             *[]string `json:"domains,omitempty"`
	Enabled             bool      `json:"enabled"`
	Groups              []string  `json:"groups"`
	KeepRoute           bool      `json:"keep_route"`
	Masquerade          bool      `json:"masquerade"`
	Metric              int       `json:"metric"`
	Network             *string   `json:"network,omitempty"`
	NetworkID           string    `json:"network_id"`
	NetworkType         string    `json:"network_type"`
	Peer                *string   `json:"peer,omitempty"`
	PeerGroups          *[]string `json:"peer_groups,omitempty"`
	AccessControlGroups *[]string `json:"access_control_groups,omitempty"`
	SkipAutoApply       *bool     `json:"skip_auto_apply,omitempty"`
}

func (r routeWire) normalized() Route {
	return Route{
		ID:                  r.ID,
		Description:         r.Description,
		Domains:             stringsOrEmpty(r.Domains),
		Enabled:             r.Enabled,
		Groups:              stringsOrEmpty(&r.Groups),
		KeepRoute:           r.KeepRoute,
		Masquerade:          r.Masquerade,
		Metric:              r.Metric,
		Network:             r.Network,
		NetworkID:           r.NetworkID,
		NetworkType:         r.NetworkType,
		Peer:                r.Peer,
		PeerGroups:          stringsOrEmpty(r.PeerGroups),
		AccessControlGroups: stringsOrEmpty(r.AccessControlGroups),
		SkipAutoApply:       r.SkipAutoApply,
	}
}

type Network struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Description       *string  `json:"description,omitempty"`
	Policies          []string `json:"policies"`
	Resources         []string `json:"resources"`
	Routers           []string `json:"routers"`
	RoutingPeersCount int      `json:"routing_peers_count"`
}

type networkWire struct {
	ID                string    `json:"id"`
	Name              string    `json:"name"`
	Description       *string   `json:"description,omitempty"`
	Policies          *[]string `json:"policies"`
	Resources         *[]string `json:"resources"`
	Routers           *[]string `json:"routers"`
	RoutingPeersCount int       `json:"routing_peers_count"`
}

func (n networkWire) normalized() Network {
	return Network{
		ID:                n.ID,
		Name:              n.Name,
		Description:       n.Description,
		Policies:          stringsOrEmpty(n.Policies),
		Resources:         stringsOrEmpty(n.Resources),
		Routers:           stringsOrEmpty(n.Routers),
		RoutingPeersCount: n.RoutingPeersCount,
	}
}

func (c *Client) ListRoutes(ctx context.Context) ([]Route, error) {
	var result []routeWire
	if err := c.transport.GetJSON(ctx, "/api/routes", &result); err != nil {
		return nil, fmt.Errorf("list routes: %w", err)
	}
	return normalizeRoutes(result), nil
}

func (c *Client) GetRoute(ctx context.Context, id string) (Route, error) {
	path, err := topologyPath("routes", id)
	if err != nil {
		return Route{}, err
	}
	var result routeWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Route{}, fmt.Errorf("get route %q: %w", id, err)
	}
	return result.normalized(), nil
}

func (c *Client) GetRouteRaw(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := topologyPath("routes", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get route %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdateRoute(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	path, err := topologyPath("routes", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update route %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListNetworks(ctx context.Context) ([]Network, error) {
	var result []networkWire
	if err := c.transport.GetJSON(ctx, "/api/networks", &result); err != nil {
		return nil, fmt.Errorf("list networks: %w", err)
	}
	normalized := make([]Network, len(result))
	for i, network := range result {
		normalized[i] = network.normalized()
	}
	return normalized, nil
}

func (c *Client) GetNetwork(ctx context.Context, id string) (Network, error) {
	path, err := topologyPath("networks", id)
	if err != nil {
		return Network{}, err
	}
	var result networkWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Network{}, fmt.Errorf("get network %q: %w", id, err)
	}
	return result.normalized(), nil
}

func (c *Client) GetNetworkRaw(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := topologyPath("networks", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get network %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdateNetwork(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	path, err := topologyPath("networks", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update network %q: %w", id, err)
	}
	return result, nil
}

func normalizeRoutes(routes []routeWire) []Route {
	result := make([]Route, len(routes))
	for i, route := range routes {
		result[i] = route.normalized()
	}
	return result
}

func stringsOrEmpty(value *[]string) []string {
	if value == nil || *value == nil {
		return []string{}
	}
	return append([]string(nil), (*value)...)
}

func topologyPath(resource, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("%s id is required", strings.TrimSuffix(resource, "s"))
	}
	return "/api/" + resource + "/" + url.PathEscape(id), nil
}
