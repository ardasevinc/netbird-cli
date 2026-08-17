package netbird

import (
	"context"
	"fmt"
	"net/http"
)

// SetupRequest is the one-shot bootstrap payload accepted while an instance
// still requires setup. Password is held in memory only and is never part of
// a stage, receipt, or owned model.
type SetupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	CreatePAT   bool   `json:"create_pat"`
	PATExpireIn int    `json:"pat_expire_in,omitempty"`
}

type SetupResponse struct {
	PersonalAccessToken string `json:"personal_access_token,omitempty"`
}

func (c *Client) GetInstance(ctx context.Context) (Instance, error) {
	var result Instance
	if err := c.transport.GetJSON(ctx, "/api/instance", &result); err != nil {
		return Instance{}, fmt.Errorf("get instance setup status: %w", err)
	}
	return result, nil
}

func (c *Client) Bootstrap(ctx context.Context, request SetupRequest) (SetupResponse, error) {
	var result SetupResponse
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, "/api/setup", request, &result); err != nil {
		return SetupResponse{}, fmt.Errorf("bootstrap instance: %w", err)
	}
	return result, nil
}
