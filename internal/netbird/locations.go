package netbird

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

func (c *Client) ListCountries(ctx context.Context) ([]string, error) {
	var result []string
	if err := c.transport.GetJSON(ctx, "/api/locations/countries", &result); err != nil {
		return nil, fmt.Errorf("list countries: %w", err)
	}
	if result == nil {
		result = []string{}
	}
	return result, nil
}

func (c *Client) ListCountryCities(ctx context.Context, country string) ([]string, error) {
	if strings.TrimSpace(country) == "" {
		return nil, fmt.Errorf("country is required")
	}
	var result []string
	path := "/api/locations/countries/" + url.PathEscape(country) + "/cities"
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list cities for country %q: %w", country, err)
	}
	if result == nil {
		result = []string{}
	}
	return result, nil
}
