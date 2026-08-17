package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type EDRBypassedPeer struct {
	PeerID string `json:"peer_id"`
}

var edrIntegrationProviders = map[string]struct{}{
	"intune":      {},
	"sentinelone": {},
	"falcon":      {},
	"huntress":    {},
	"fleetdm":     {},
}

func edrIntegrationPath(provider string) (string, error) {
	provider = strings.ToLower(strings.TrimSpace(provider))
	if _, ok := edrIntegrationProviders[provider]; !ok {
		return "", fmt.Errorf("unsupported EDR integration provider %q", provider)
	}
	return "/api/integrations/edr/" + url.PathEscape(provider), nil
}

func (c *Client) GetEDRIntegrationRaw(ctx context.Context, provider string) (json.RawMessage, error) {
	path, err := edrIntegrationPath(provider)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get EDR %s integration: %w", provider, err)
	}
	return result, nil
}

func (c *Client) CreateEDRIntegration(ctx context.Context, provider string, request json.RawMessage) (json.RawMessage, error) {
	path, err := edrIntegrationPath(provider)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("create EDR %s integration: %w", provider, err)
	}
	return result, nil
}

func (c *Client) UpdateEDRIntegration(ctx context.Context, provider string, request json.RawMessage) (json.RawMessage, error) {
	path, err := edrIntegrationPath(provider)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update EDR %s integration: %w", provider, err)
	}
	return result, nil
}

func (c *Client) DeleteEDRIntegration(ctx context.Context, provider string) (json.RawMessage, error) {
	path, err := edrIntegrationPath(provider)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete EDR %s integration: %w", provider, err)
	}
	return result, nil
}

func (c *Client) ListEDRBypassedPeers(ctx context.Context) ([]EDRBypassedPeer, error) {
	var result []EDRBypassedPeer
	if err := c.transport.GetJSON(ctx, "/api/peers/edr/bypassed", &result); err != nil {
		return nil, fmt.Errorf("list EDR-bypassed peers: %w", err)
	}
	if result == nil {
		result = []EDRBypassedPeer{}
	}
	return result, nil
}

func (c *Client) ListEDRBypassedPeersRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/peers/edr/bypassed", nil, &result); err != nil {
		return nil, fmt.Errorf("list EDR-bypassed peers: %w", err)
	}
	return result, nil
}

func (c *Client) BypassPeerEDR(ctx context.Context, peerID string) (json.RawMessage, error) {
	if strings.TrimSpace(peerID) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	var result json.RawMessage
	path := "/api/peers/" + url.PathEscape(peerID) + "/edr/bypass"
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, nil, &result); err != nil {
		return nil, fmt.Errorf("bypass EDR for peer %q: %w", peerID, err)
	}
	return result, nil
}

func (c *Client) RevokePeerEDRBypass(ctx context.Context, peerID string) (json.RawMessage, error) {
	if strings.TrimSpace(peerID) == "" {
		return nil, fmt.Errorf("peer id is required")
	}
	var result json.RawMessage
	path := "/api/peers/" + url.PathEscape(peerID) + "/edr/bypass"
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("revoke EDR bypass for peer %q: %w", peerID, err)
	}
	return result, nil
}
