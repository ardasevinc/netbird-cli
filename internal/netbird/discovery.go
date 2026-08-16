package netbird

import (
	"context"
	"fmt"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

type Version struct {
	ManagementCurrentVersion  string `json:"management_current_version"`
	ManagementUpdateAvailable bool   `json:"management_update_available"`
}

type Instance struct {
	SetupRequired bool `json:"setup_required"`
}

type User struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	IsServiceUser *bool  `json:"is_service_user,omitempty"`
}

type Discovery struct {
	Version  Version  `json:"version"`
	Instance Instance `json:"instance"`
	User     *User    `json:"user,omitempty"`
}

type Client struct {
	transport *transport.Client
}

func NewClient(transportClient *transport.Client) *Client {
	return &Client{transport: transportClient}
}

func (c *Client) Discover(ctx context.Context, authenticated bool) (Discovery, error) {
	var result Discovery
	if err := c.transport.GetJSON(ctx, "/api/instance/version", &result.Version); err != nil {
		return Discovery{}, fmt.Errorf("discover server version: %w", err)
	}
	if err := c.transport.GetJSON(ctx, "/api/instance", &result.Instance); err != nil {
		return Discovery{}, fmt.Errorf("discover instance status: %w", err)
	}
	if authenticated {
		var user User
		if err := c.transport.GetJSON(ctx, "/api/users/current", &user); err != nil {
			return Discovery{}, fmt.Errorf("discover current user: %w", err)
		}
		result.User = &user
	}
	return result, nil
}
