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
	AccountID     string `json:"account_id,omitempty"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Role          string `json:"role"`
	IsServiceUser *bool  `json:"is_service_user,omitempty"`
}

type Account struct {
	ID string `json:"id"`
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

func (c *Client) ServerIdentity() string { return c.transport.Origin() }

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

// AccountScope verifies that the authenticated token can see the configured
// account. The management API returns the caller's visible account list.
func (c *Client) AccountScope(ctx context.Context, accountID string) error {
	if accountID == "" {
		return fmt.Errorf("account scope is required")
	}
	var accounts []Account
	if err := c.transport.GetJSON(ctx, "/api/accounts", &accounts); err != nil {
		return fmt.Errorf("verify account scope: %w", err)
	}
	for _, account := range accounts {
		if account.ID == accountID {
			return nil
		}
	}
	return fmt.Errorf("authenticated credential cannot access account %q", accountID)
}
