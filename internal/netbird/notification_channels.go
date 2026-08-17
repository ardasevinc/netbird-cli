package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *Client) ListNotificationChannelsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/integrations/notifications/channels", nil, &result); err != nil {
		return nil, fmt.Errorf("list notification channels: %w", err)
	}
	return result, nil
}

func (c *Client) GetNotificationChannelRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("notification channel id is required")
	}
	path := "/api/integrations/notifications/channels/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get notification channel %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateNotificationChannel(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/integrations/notifications/channels", request, &result); err != nil {
		return nil, fmt.Errorf("create notification channel: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateNotificationChannel(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("notification channel id is required")
	}
	path := "/api/integrations/notifications/channels/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update notification channel %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteNotificationChannel(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("notification channel id is required")
	}
	path := "/api/integrations/notifications/channels/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete notification channel %q: %w", id, err)
	}
	return result, nil
}
