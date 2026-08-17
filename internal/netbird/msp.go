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

const mspTenantsPath = "/api/integrations/msp/tenants"

func (c *Client) ListMSPTenantsRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, mspTenantsPath, nil, &result); err != nil {
		return nil, fmt.Errorf("list MSP tenants: %w", err)
	}
	return result, nil
}

func (c *Client) GetMSPTenantRaw(ctx context.Context, id string) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("MSP tenant id is required")
	}
	collection, err := c.ListMSPTenantsRaw(ctx)
	if err != nil {
		return nil, err
	}
	var tenants []json.RawMessage
	if err := json.Unmarshal(collection, &tenants); err != nil {
		return nil, fmt.Errorf("decode MSP tenant collection: %w", err)
	}
	for _, tenant := range tenants {
		var candidate struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(tenant, &candidate); err != nil {
			return nil, fmt.Errorf("decode MSP tenant: %w", err)
		}
		if candidate.ID == id {
			return tenant, nil
		}
	}
	return nil, &transport.RequestError{Dispatched: true, StatusCode: http.StatusNotFound, Description: fmt.Sprintf("MSP tenant %q was not found", id)}
}

func (c *Client) CreateMSPTenant(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPost, mspTenantsPath, request, &result); err != nil {
		return nil, fmt.Errorf("create MSP tenant: %w", err)
	}
	return result, nil
}

func (c *Client) UpdateMSPTenant(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("MSP tenant id is required")
	}
	path := mspTenantsPath + "/" + url.PathEscape(id)
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodPut, path, request, &result); err != nil {
		return nil, fmt.Errorf("update MSP tenant %q: %w", id, err)
	}
	return result, nil
}

func (c *Client) VerifyMSPTenantDNS(ctx context.Context, id string) (json.RawMessage, error) {
	return c.mspTenantAction(ctx, id, "dns", "verify MSP tenant DNS", http.MethodPost, nil)
}

func (c *Client) InviteMSPTenant(ctx context.Context, id string) (json.RawMessage, error) {
	return c.mspTenantAction(ctx, id, "invite", "invite MSP tenant", http.MethodPost, nil)
}

func (c *Client) RespondMSPTenantInvite(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	return c.mspTenantAction(ctx, id, "invite", "respond to MSP tenant invite", http.MethodPut, request)
}

func (c *Client) CreateMSPTenantSubscription(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	return c.mspTenantAction(ctx, id, "subscription", "create MSP tenant subscription", http.MethodPost, request)
}

func (c *Client) UnlinkMSPTenant(ctx context.Context, id string, request json.RawMessage) (json.RawMessage, error) {
	return c.mspTenantAction(ctx, id, "unlink", "unlink MSP tenant", http.MethodPost, request)
}

func (c *Client) mspTenantAction(ctx context.Context, id, suffix, operation string, method string, request json.RawMessage) (json.RawMessage, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("MSP tenant id is required")
	}
	path := mspTenantsPath + "/" + url.PathEscape(id) + "/" + suffix
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, method, path, request, &result); err != nil {
		return nil, fmt.Errorf("%s %q: %w", operation, id, err)
	}
	return result, nil
}
