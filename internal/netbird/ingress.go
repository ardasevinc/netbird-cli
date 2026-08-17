package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type IngressPeer struct {
	ID           string `json:"id"`
	PeerID       string `json:"peer_id"`
	IngressIP    string `json:"ingress_ip"`
	Region       string `json:"region"`
	Connected    bool   `json:"connected"`
	Enabled      bool   `json:"enabled"`
	Fallback     bool   `json:"fallback"`
	AvailableTCP int    `json:"available_tcp"`
	AvailableUDP int    `json:"available_udp"`
}

type ingressPeerWire struct {
	ID             string `json:"id"`
	PeerID         string `json:"peer_id"`
	IngressIP      string `json:"ingress_ip"`
	Region         string `json:"region"`
	Connected      bool   `json:"connected"`
	Enabled        bool   `json:"enabled"`
	Fallback       bool   `json:"fallback"`
	AvailablePorts struct {
		TCP int `json:"tcp"`
		UDP int `json:"udp"`
	} `json:"available_ports"`
}

func (p ingressPeerWire) normalized() IngressPeer {
	return IngressPeer{ID: p.ID, PeerID: p.PeerID, IngressIP: p.IngressIP, Region: p.Region, Connected: p.Connected, Enabled: p.Enabled, Fallback: p.Fallback, AvailableTCP: p.AvailablePorts.TCP, AvailableUDP: p.AvailablePorts.UDP}
}

type IngressPortMapping struct {
	IngressStart    int    `json:"ingress_start"`
	IngressEnd      int    `json:"ingress_end"`
	TranslatedStart int    `json:"translated_start"`
	TranslatedEnd   int    `json:"translated_end"`
	Protocol        string `json:"protocol"`
}

type IngressPortAllocation struct {
	ID                string               `json:"id"`
	Name              string               `json:"name"`
	IngressIP         string               `json:"ingress_ip"`
	IngressPeerID     string               `json:"ingress_peer_id"`
	Region            string               `json:"region"`
	Enabled           bool                 `json:"enabled"`
	PortRangeMappings []IngressPortMapping `json:"port_range_mappings"`
}

func (c *Client) ListIngressPeers(ctx context.Context) ([]IngressPeer, error) {
	var wire []ingressPeerWire
	if err := c.transport.GetJSON(ctx, "/api/ingress/peers", &wire); err != nil {
		return nil, fmt.Errorf("list ingress peers: %w", err)
	}
	result := make([]IngressPeer, len(wire))
	for i, item := range wire {
		result[i] = item.normalized()
	}
	return result, nil
}

func (c *Client) ListIngressPeersRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/ingress/peers", nil, &result); err != nil {
		return nil, fmt.Errorf("list ingress peers: %w", err)
	}
	return result, nil
}

func (c *Client) CreateIngressPeer(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/ingress/peers", request, &result); err != nil {
		return nil, fmt.Errorf("create ingress peer: %w", err)
	}
	return result, nil
}

func (c *Client) GetIngressPeer(ctx context.Context, id string) (IngressPeer, error) {
	if strings.TrimSpace(id) == "" {
		return IngressPeer{}, fmt.Errorf("ingress peer id is required")
	}
	var wire ingressPeerWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/ingress/peers/"+url.PathEscape(id), nil, &wire); err != nil {
		return IngressPeer{}, fmt.Errorf("get ingress peer %q: %w", id, err)
	}
	return wire.normalized(), nil
}

func (c *Client) GetIngressPeerRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ingress peer id is required")
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/ingress/peers/"+url.PathEscape(id), nil, &result); err != nil {
		return nil, fmt.Errorf("get ingress peer %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdateIngressPeer(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ingress peer id is required")
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, "/api/ingress/peers/"+url.PathEscape(id), request, &result); err != nil {
		return nil, fmt.Errorf("update ingress peer %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteIngressPeer(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("ingress peer id is required")
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, "/api/ingress/peers/"+url.PathEscape(id), nil, &result); err != nil {
		return nil, fmt.Errorf("delete ingress peer %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListIngressPortAllocations(ctx context.Context, peerID string) ([]IngressPortAllocation, error) {
	path, err := ingressPortsPath(peerID, "")
	if err != nil {
		return nil, err
	}
	var result []IngressPortAllocation
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list ingress ports for peer %q: %w", peerID, err)
	}
	if result == nil {
		result = []IngressPortAllocation{}
	}
	for i := range result {
		if result[i].PortRangeMappings == nil {
			result[i].PortRangeMappings = []IngressPortMapping{}
		}
	}
	return result, nil
}

func (c *Client) GetIngressPortAllocation(ctx context.Context, peerID, allocationID string) (IngressPortAllocation, error) {
	path, err := ingressPortsPath(peerID, allocationID)
	if err != nil {
		return IngressPortAllocation{}, err
	}
	var result IngressPortAllocation
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return IngressPortAllocation{}, fmt.Errorf("get ingress port allocation %q for peer %q: %w", allocationID, peerID, err)
	}
	if result.PortRangeMappings == nil {
		result.PortRangeMappings = []IngressPortMapping{}
	}
	return result, nil
}

func (c *Client) ListIngressPortAllocationsRaw(ctx context.Context, peerID string) (json.RawMessage, error) {
	path, err := ingressPortsPath(peerID, "")
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("list ingress ports for peer %q: %w", peerID, err)
	}
	return result, nil
}

func (c *Client) GetIngressPortAllocationRaw(ctx context.Context, peerID, allocationID string) (json.RawMessage, error) {
	path, err := ingressPortsPath(peerID, allocationID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get ingress port allocation %q for peer %q: %w", allocationID, peerID, err)
	}
	return result, nil
}

func (c *Client) CreateIngressPortAllocation(ctx context.Context, peerID string, request json.RawMessage) (json.RawMessage, error) {
	path, err := ingressPortsPath(peerID, "")
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("create ingress port allocation for peer %q: %w", peerID, err)
	}
	return result, nil
}

func (c *Client) UpdateIngressPortAllocation(ctx context.Context, peerID, allocationID string, request json.RawMessage) (json.RawMessage, error) {
	path, err := ingressPortsPath(peerID, allocationID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update ingress port allocation %q for peer %q: %w", allocationID, peerID, err)
	}
	return result, nil
}

func (c *Client) DeleteIngressPortAllocation(ctx context.Context, peerID, allocationID string) (json.RawMessage, error) {
	path, err := ingressPortsPath(peerID, allocationID)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete ingress port allocation %q for peer %q: %w", allocationID, peerID, err)
	}
	return result, nil
}

func ingressPortsPath(peerID, allocationID string) (string, error) {
	if strings.TrimSpace(peerID) == "" {
		return "", fmt.Errorf("peer id is required")
	}
	path := "/api/peers/" + url.PathEscape(peerID) + "/ingress/ports"
	if allocationID != "" {
		if strings.TrimSpace(allocationID) == "" {
			return "", fmt.Errorf("allocation id is required")
		}
		path += "/" + url.PathEscape(allocationID)
	}
	return path, nil
}
