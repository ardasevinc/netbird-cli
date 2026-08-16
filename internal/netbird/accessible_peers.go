package netbird

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) ListAccessiblePeers(ctx context.Context, peerID string) ([]Peer, error) {
	if strings.TrimSpace(peerID) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	var result []peerWire
	path := "/api/peers/" + url.PathEscape(peerID) + "/accessible-peers"
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list accessible peers for %q: %w", peerID, err)
	}
	normalized := make([]Peer, len(result))
	for i, peer := range result {
		normalized[i] = peer.normalized()
	}
	return normalized, nil
}
