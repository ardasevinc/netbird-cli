package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// ManagedUser is the stable, non-secret user inventory view. Passwords and
// invite tokens are intentionally not represented, even if an upstream
// response includes them.
type ManagedUser struct {
	ID              string           `json:"id"`
	Email           string           `json:"email"`
	Name            string           `json:"name"`
	Role            string           `json:"role"`
	Status          string           `json:"status"`
	AutoGroups      []string         `json:"auto_groups"`
	IsCurrent       *bool            `json:"is_current,omitempty"`
	IsServiceUser   *bool            `json:"is_service_user,omitempty"`
	IsBlocked       bool             `json:"is_blocked"`
	PendingApproval bool             `json:"pending_approval"`
	Issued          *string          `json:"issued,omitempty"`
	IDPID           *string          `json:"idp_id,omitempty"`
	LastLogin       *string          `json:"last_login,omitempty"`
	Permissions     *UserPermissions `json:"permissions,omitempty"`
}

type UserPermissions struct {
	IsRestricted bool                       `json:"is_restricted"`
	Modules      map[string]map[string]bool `json:"modules"`
}

type UserInvite struct {
	ID         string   `json:"id"`
	Email      string   `json:"email"`
	Name       string   `json:"name"`
	Role       string   `json:"role"`
	AutoGroups []string `json:"auto_groups"`
	ExpiresAt  string   `json:"expires_at"`
	CreatedAt  string   `json:"created_at"`
	Expired    bool     `json:"expired"`
}

func (c *Client) ListAccounts(ctx context.Context) ([]Account, error) {
	var result []Account
	if err := c.transport.GetJSON(ctx, "/api/accounts", &result); err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return result, nil
}

func (c *Client) GetAccountRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if id == "" {
		return nil, fmt.Errorf("account id is required")
	}
	path := "/api/accounts/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get account %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) UpdateAccount(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if id == "" {
		return nil, fmt.Errorf("account id is required")
	}
	path := "/api/accounts/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update account %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ListUsers(ctx context.Context, serviceUser *bool) ([]ManagedUser, error) {
	path := "/api/users"
	if serviceUser != nil {
		query := url.Values{}
		query.Set("service_user", fmt.Sprintf("%t", *serviceUser))
		path += "?" + query.Encode()
	}
	var result []ManagedUser
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return result, nil
}

func (c *Client) ListInvites(ctx context.Context) ([]UserInvite, error) {
	var result []UserInvite
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/users/invites", nil, &result); err != nil {
		return nil, fmt.Errorf("list user invites: %w", err)
	}
	return result, nil
}
