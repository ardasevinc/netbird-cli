package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/transport"
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

type PublicUserInvite struct {
	Email     string `json:"email"`
	Name      string `json:"name"`
	ExpiresAt string `json:"expires_at"`
	Valid     bool   `json:"valid"`
	InvitedBy string `json:"invited_by"`
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

func (c *Client) DeleteAccount(ctx context.Context, id string) (json.RawMessage, error) {
	if id == "" {
		return nil, fmt.Errorf("account id is required")
	}
	path := "/api/accounts/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete account %q: %w", id, err)
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

func (c *Client) ListUsersRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/users", nil, &result); err != nil {
		return nil, fmt.Errorf("list users: %w", err)
	}
	return result, nil
}

func (c *Client) GetUserRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get user %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) CreateUser(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/users", request, &result); err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateUser(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update user %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) DeleteUser(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete user %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ApproveUser(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id) + "/approve"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, nil, &result); err != nil {
		return nil, fmt.Errorf("approve user %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) RejectUser(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id) + "/reject"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("reject user %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) ChangeUserPassword(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	path := "/api/users/" + url.PathEscape(id) + "/password"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("change password for user %q: %w", id, err)
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

func (c *Client) ListInvitesRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/users/invites", nil, &result); err != nil {
		return nil, fmt.Errorf("list user invites: %w", err)
	}
	return result, nil
}

func (c *Client) GetInviteRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("invite id is required")
	}
	collection, err := c.ListInvitesRaw(ctx)
	if err != nil {
		return nil, err
	}
	var invites []json.RawMessage
	if err := json.Unmarshal(collection, &invites); err != nil {
		return nil, fmt.Errorf("decode user invites: %w", err)
	}
	for _, invite := range invites {
		var candidate struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(invite, &candidate); err != nil {
			return nil, fmt.Errorf("decode user invite: %w", err)
		}
		if candidate.ID == id {
			return invite, nil
		}
	}
	return nil, &transport.RequestError{Dispatched: true, StatusCode: http.StatusNotFound, Description: "invite not found"}
}

func (c *Client) CreateInvite(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/users/invites", request, &result); err != nil {
		return nil, fmt.Errorf("create user invite: %w", err)
	}
	return result, nil
}

func (c *Client) DeleteInvite(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("invite id is required")
	}
	path := "/api/users/invites/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodDelete, path, nil, &result); err != nil {
		return nil, fmt.Errorf("delete user invite %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) RegenerateInvite(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("invite id is required")
	}
	path := "/api/users/invites/" + url.PathEscape(id) + "/regenerate"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("regenerate user invite %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) GetPublicInviteRaw(ctx context.Context, token string) (json.RawMessage, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("invite token is required")
	}
	path := "/api/users/invites/" + url.PathEscape(token)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return nil, fmt.Errorf("get public user invite: %w", err)
	}
	return sanitizePublicInvite(result)
}

func (c *Client) AcceptInvite(ctx context.Context, token string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(token) == "" {
		return nil, fmt.Errorf("invite token is required")
	}
	path := "/api/users/invites/" + url.PathEscape(token) + "/accept"
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, path, request, &result); err != nil {
		return nil, fmt.Errorf("accept user invite: %w", err)
	}
	return result, nil
}

func sanitizePublicInvite(raw json.RawMessage) (json.RawMessage, error) {
	var object map[string]any
	if err := json.Unmarshal(raw, &object); err != nil {
		return nil, fmt.Errorf("decode public user invite: %w", err)
	}
	for _, field := range []string{"token", "invite_token", "password"} {
		delete(object, field)
	}
	return json.Marshal(object)
}
