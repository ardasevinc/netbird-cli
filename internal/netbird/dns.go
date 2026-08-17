package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type Nameserver struct {
	IP   string `json:"ip"`
	Type string `json:"ns_type"`
	Port int    `json:"port"`
}

type NameserverGroup struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Domains              []string     `json:"domains"`
	Enabled              bool         `json:"enabled"`
	Groups               []string     `json:"groups"`
	Nameservers          []Nameserver `json:"nameservers"`
	Primary              bool         `json:"primary"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

type DNSSettings struct {
	DisabledManagementGroups []string `json:"disabled_management_groups"`
}

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Type    string `json:"type"`
	Content string `json:"content"`
	TTL     int    `json:"ttl"`
}

type DNSZone struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Domain             string      `json:"domain"`
	Enabled            bool        `json:"enabled"`
	EnableSearchDomain bool        `json:"enable_search_domain"`
	DistributionGroups []string    `json:"distribution_groups"`
	Records            []DNSRecord `json:"records"`
}

type nameserverGroupWire struct {
	ID                   string       `json:"id"`
	Name                 string       `json:"name"`
	Description          string       `json:"description"`
	Domains              *[]string    `json:"domains"`
	Enabled              bool         `json:"enabled"`
	Groups               *[]string    `json:"groups"`
	Nameservers          []Nameserver `json:"nameservers"`
	Primary              bool         `json:"primary"`
	SearchDomainsEnabled bool         `json:"search_domains_enabled"`
}

func (n nameserverGroupWire) normalized() NameserverGroup {
	return NameserverGroup{
		ID:                   n.ID,
		Name:                 n.Name,
		Description:          n.Description,
		Domains:              stringsOrEmpty(n.Domains),
		Enabled:              n.Enabled,
		Groups:               stringsOrEmpty(n.Groups),
		Nameservers:          append([]Nameserver(nil), n.Nameservers...),
		Primary:              n.Primary,
		SearchDomainsEnabled: n.SearchDomainsEnabled,
	}
}

type dnsZoneWire struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	Domain             string      `json:"domain"`
	Enabled            bool        `json:"enabled"`
	EnableSearchDomain bool        `json:"enable_search_domain"`
	DistributionGroups *[]string   `json:"distribution_groups"`
	Records            []DNSRecord `json:"records"`
}

func (z dnsZoneWire) normalized() DNSZone {
	return DNSZone{
		ID:                 z.ID,
		Name:               z.Name,
		Domain:             z.Domain,
		Enabled:            z.Enabled,
		EnableSearchDomain: z.EnableSearchDomain,
		DistributionGroups: stringsOrEmpty(z.DistributionGroups),
		Records:            append([]DNSRecord(nil), z.Records...),
	}
}

func (c *Client) ListNameserverGroups(ctx context.Context) ([]NameserverGroup, error) {
	var result []nameserverGroupWire
	if err := c.transport.GetJSON(ctx, "/api/dns/nameservers", &result); err != nil {
		return nil, fmt.Errorf("list nameserver groups: %w", err)
	}
	normalized := make([]NameserverGroup, len(result))
	for i, group := range result {
		normalized[i] = group.normalized()
	}
	return normalized, nil
}

func (c *Client) GetNameserverGroup(ctx context.Context, id string) (NameserverGroup, error) {
	path, err := dnsPath("nameservers", id)
	if err != nil {
		return NameserverGroup{}, err
	}
	var result nameserverGroupWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return NameserverGroup{}, fmt.Errorf("get nameserver group %q: %w", id, err)
	}
	return result.normalized(), nil
}

func (c *Client) GetDNSSettings(ctx context.Context) (DNSSettings, error) {
	var result DNSSettings
	if err := c.transport.GetJSON(ctx, "/api/dns/settings", &result); err != nil {
		return DNSSettings{}, fmt.Errorf("get dns settings: %w", err)
	}
	if result.DisabledManagementGroups == nil {
		result.DisabledManagementGroups = []string{}
	}
	return result, nil
}

func (c *Client) ListDNSZones(ctx context.Context) ([]DNSZone, error) {
	var result []dnsZoneWire
	if err := c.transport.GetJSON(ctx, "/api/dns/zones", &result); err != nil {
		return nil, fmt.Errorf("list dns zones: %w", err)
	}
	normalized := make([]DNSZone, len(result))
	for i, zone := range result {
		normalized[i] = zone.normalized()
	}
	return normalized, nil
}

func (c *Client) ListDNSZonesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/dns/zones", nil, &result); err != nil {
		return nil, fmt.Errorf("list dns zones: %w", err)
	}
	return result, nil
}

func (c *Client) CreateDNSZone(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/dns/zones", request, &result); err != nil {
		return nil, fmt.Errorf("create dns zone: %w", err)
	}
	return result, nil
}

func (c *Client) GetDNSZoneRaw(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := dnsPath("zones", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get dns zone %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdateDNSZone(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	path, err := dnsPath("zones", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update dns zone %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteDNSZone(ctx context.Context, id string) (json.RawMessage, error) {
	path, err := dnsPath("zones", id)
	if err != nil {
		return nil, err
	}
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete dns zone %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) GetDNSZone(ctx context.Context, id string) (DNSZone, error) {
	path, err := dnsPath("zones", id)
	if err != nil {
		return DNSZone{}, err
	}
	var result dnsZoneWire
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return DNSZone{}, fmt.Errorf("get dns zone %q: %w", id, err)
	}
	return result.normalized(), nil
}

func dnsPath(resource, id string) (string, error) {
	if strings.TrimSpace(id) == "" {
		return "", fmt.Errorf("%s id is required", strings.TrimSuffix(resource, "s"))
	}
	return "/api/dns/" + resource + "/" + url.PathEscape(id), nil
}
