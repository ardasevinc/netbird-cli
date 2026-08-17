package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListEventStreamingIntegrationsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/event-streaming", nil, &result); err != nil {
		return nil, fmt.Errorf("list event-streaming integrations: %w", err)
	}
	return result, nil
}

func (c *Client) GetEventStreamingIntegrationRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("event-streaming integration id is required")
	}
	path := "/api/event-streaming/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get event-streaming integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateEventStreamingIntegration(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/event-streaming", request, &result); err != nil {
		return nil, fmt.Errorf("create event-streaming integration: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateEventStreamingIntegration(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("event-streaming integration id is required")
	}
	path := "/api/event-streaming/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update event-streaming integration %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteEventStreamingIntegration(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("event-streaming integration id is required")
	}
	path := "/api/event-streaming/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete event-streaming integration %q: %w", id, err)
	}
	return result, nil
}
