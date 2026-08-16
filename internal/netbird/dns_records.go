package netbird

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListDNSRecords(ctx context.Context, zoneID string) ([]DNSRecord, error) {
	path, err := dnsRecordPath(zoneID, "")
	if err != nil {
		return nil, err
	}
	var result []DNSRecord
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list dns records for zone %q: %w", zoneID, err)
	}
	if result == nil {
		result = []DNSRecord{}
	}
	return result, nil
}

func (c *Client) GetDNSRecord(ctx context.Context, zoneID, recordID string) (DNSRecord, error) {
	path, err := dnsRecordPath(zoneID, recordID)
	if err != nil {
		return DNSRecord{}, err
	}
	var result DNSRecord
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return DNSRecord{}, fmt.Errorf("get dns record %q for zone %q: %w", recordID, zoneID, err)
	}
	return result, nil
}

func dnsRecordPath(zoneID, recordID string) (string, error) {
	if strings.TrimSpace(zoneID) == "" {
		return "", fmt.Errorf("zone id is required")
	}
	path := "/api/dns/zones/" + url.PathEscape(zoneID) + "/records"
	if recordID != "" {
		if strings.TrimSpace(recordID) == "" {
			return "", fmt.Errorf("record id is required")
		}
		path += "/" + url.PathEscape(recordID)
	}
	return path, nil
}
