package netbird

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

// PersonalAccessToken is metadata only. The management API does not return a
// token value from these read endpoints, and the owned model has no secret
// field by design.
type PersonalAccessToken struct {
	ID             string  `json:"id"`
	Name           string  `json:"name"`
	CreatedBy      string  `json:"created_by"`
	CreatedAt      string  `json:"created_at"`
	ExpirationDate string  `json:"expiration_date"`
	LastUsed       *string `json:"last_used,omitempty"`
}

func (c *Client) ListPersonalAccessTokens(ctx context.Context, userID string) ([]PersonalAccessToken, error) {
	if strings.TrimSpace(userID) == "" {
		return nil, fmt.Errorf("user id is required")
	}
	var result []PersonalAccessToken
	path := "/api/users/" + url.PathEscape(userID) + "/tokens"
	if err := c.transport.GetJSON(ctx, path, &result); err != nil {
		return nil, fmt.Errorf("list tokens for user %q: %w", userID, err)
	}
	if result == nil {
		result = []PersonalAccessToken{}
	}
	return result, nil
}

func (c *Client) GetPersonalAccessToken(ctx context.Context, userID, tokenID string) (PersonalAccessToken, error) {
	if strings.TrimSpace(userID) == "" {
		return PersonalAccessToken{}, fmt.Errorf("user id is required")
	}
	if strings.TrimSpace(tokenID) == "" {
		return PersonalAccessToken{}, fmt.Errorf("token id is required")
	}
	var result PersonalAccessToken
	path := "/api/users/" + url.PathEscape(userID) + "/tokens/" + url.PathEscape(tokenID)
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return PersonalAccessToken{}, fmt.Errorf("get token %q for user %q: %w", tokenID, userID, err)
	}
	return result, nil
}
