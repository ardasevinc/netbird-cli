package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type PeerGroup struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
}

// Peer is the owned, stable inventory view. Provider-specific fields remain
// behind the adapter until they earn a versioned nb schema.
type Peer struct {
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	Hostname             string      `json:"hostname"`
	IP                   string      `json:"ip"`
	IPv6                 *string     `json:"ipv6,omitempty"`
	OS                   string      `json:"os"`
	Version              string      `json:"version"`
	Connected            bool        `json:"connected"`
	Ephemeral            bool        `json:"ephemeral"`
	LoginExpired         bool        `json:"login_expired"`
	ApprovalRequired     bool        `json:"approval_required"`
	AccessiblePeersCount int         `json:"accessible_peers_count,omitempty"`
	Groups               []PeerGroup `json:"groups"`
	CreatedAt            string      `json:"created_at"`
	LastSeen             string      `json:"last_seen"`
}

type peerWire struct {
	ID                   string      `json:"id"`
	Name                 string      `json:"name"`
	Hostname             string      `json:"hostname"`
	IP                   string      `json:"ip"`
	IPv6                 *string     `json:"ipv6,omitempty"`
	OS                   string      `json:"os"`
	Version              string      `json:"version"`
	Connected            bool        `json:"connected"`
	Ephemeral            bool        `json:"ephemeral"`
	LoginExpired         bool        `json:"login_expired"`
	ApprovalRequired     bool        `json:"approval_required"`
	AccessiblePeersCount int         `json:"accessible_peers_count"`
	Groups               []PeerGroup `json:"groups"`
	CreatedAt            string      `json:"created_at"`
	LastSeen             string      `json:"last_seen"`
}

func (p peerWire) normalized() Peer {
	groups := p.Groups
	if groups == nil {
		groups = []PeerGroup{}
	}
	p.Groups = groups
	return Peer(p)
}

func (c *Client) ListPeers(ctx context.Context, name, ip string) ([]Peer, error) {
	path := "/api/peers"
	query := url.Values{}
	if name != "" {
		query.Set("name", name)
	}
	if ip != "" {
		query.Set("ip", ip)
	}
	if encoded := query.Encode(); encoded != "" {
		path += "?" + encoded
	}
	var result []peerWire
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list peers: %w", err)
	}
	normalized := make([]Peer, len(result))
	for i, peer := range result {
		normalized[i] = peer.normalized()
	}
	return normalized, nil
}

func (c *Client) GetPeer(ctx context.Context, id string) (Peer, error) {
	if strings.TrimSpace(id) == "" {
		return Peer{}, fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(id)
	var result peerWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return Peer{}, fmt.Errorf("get peer %q: %w", id, err)
	}
	return result.normalized(), nil
}

func (c *Client) GetPeerRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get peer %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdatePeer(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update peer %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateTemporaryAccessPeer(ctx context.Context, peerID string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(peerID) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(peerID) + "/temporary-access"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("create temporary access peer for %q: %w", peerID, err)
	}
	return result, nil
}

func (c *Client) DeletePeer(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete peer %q: %w", id, err)
	}
	return result, nil
}
