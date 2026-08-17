package netbird

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const billingSubscriptionPath = "/api/integrations/billing/subscription"

// GetBillingSubscriptionRaw returns the account's billing subscription. The
// read is also the capability probe used before billing mutations: a server
// that does not expose billing must fail closed rather than look like an
// empty subscription.
func (c *Client) GetBillingSubscriptionRaw(ctx context.Context) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, billingSubscriptionPath, nil, &result); err != nil {
		return nil, fmt.Errorf("get billing subscription: %w", err)
	}
	return result, nil
}

func (c *Client) ActivateAWSMarketplaceBilling(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	return c.billingMutation(ctx, http.MethodPost, "/api/integrations/billing/aws/marketplace/activate", "activate AWS Marketplace subscription", request)
}

func (c *Client) EnrichAWSMarketplaceBilling(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	return c.billingMutation(ctx, http.MethodPost, "/api/integrations/billing/aws/marketplace/enrich", "enrich AWS Marketplace subscription", request)
}

func (c *Client) CreateBillingCheckout(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	return c.billingMutation(ctx, http.MethodPost, "/api/integrations/billing/checkout", "create billing checkout", request)
}

func (c *Client) UpdateBillingSubscription(ctx context.Context, request json.RawMessage) (json.RawMessage, error) {
	return c.billingMutation(ctx, http.MethodPut, billingSubscriptionPath, "update billing subscription", request)
}

func (c *Client) billingMutation(ctx context.Context, method, path, operation string, request json.RawMessage) (json.RawMessage, error) {
	var result json.RawMessage
	if _, err := c.transport.DoJSON(ctx, method, path, request, &result); err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	return result, nil
}
