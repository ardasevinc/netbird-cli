package netbird

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/transport"
)

func TestBillingIntegrationRoutesAndBodies(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("content-type", "application/json")
		var body map[string]any
		if r.Body != nil {
			_ = json.NewDecoder(r.Body).Decode(&body)
		}
		w.WriteHeader(http.StatusOK)
		switch {
		case r.URL.Path == billingSubscriptionPath && r.Method == http.MethodGet:
			_, _ = w.Write([]byte(`{"active":true,"plan_tier":"basic","price_id":"price_basic"}`))
		case r.URL.Path == "/api/integrations/billing/aws/marketplace/activate" && r.Method == http.MethodPost:
			assertBillingBody(t, body, map[string]any{"plan_tier": "business"})
			_, _ = w.Write([]byte(`{"active":true,"plan_tier":"business","price_id":"price_business"}`))
		case r.URL.Path == "/api/integrations/billing/aws/marketplace/enrich" && r.Method == http.MethodPost:
			assertBillingBody(t, body, map[string]any{"aws_user_id": "aws-user"})
			_, _ = w.Write([]byte(`{"ok":true}`))
		case r.URL.Path == "/api/integrations/billing/checkout" && r.Method == http.MethodPost:
			assertBillingBody(t, body, map[string]any{"baseURL": "https://app.example/success", "priceID": "price_business", "enableTrial": true})
			_, _ = w.Write([]byte(`{"session_id":"cs_test_123","url":"https://checkout.example/cs_test_123"}`))
		case r.URL.Path == billingSubscriptionPath && r.Method == http.MethodPut:
			assertBillingBody(t, body, map[string]any{"priceID": "price_business", "plan_tier": "business"})
			_, _ = w.Write([]byte(`{"active":true,"plan_tier":"business","price_id":"price_business"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client, err := transport.New(transport.Config{BaseURL: server.URL})
	if err != nil {
		t.Fatal(err)
	}
	api := NewClient(client)
	ctx := context.Background()
	if got, err := api.GetBillingSubscriptionRaw(ctx); err != nil || len(got) == 0 {
		t.Fatalf("get subscription: got=%s err=%v", got, err)
	}
	if _, err := api.ActivateAWSMarketplaceBilling(ctx, json.RawMessage(`{"plan_tier":"business"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.EnrichAWSMarketplaceBilling(ctx, json.RawMessage(`{"aws_user_id":"aws-user"}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.CreateBillingCheckout(ctx, json.RawMessage(`{"baseURL":"https://app.example/success","priceID":"price_business","enableTrial":true}`)); err != nil {
		t.Fatal(err)
	}
	if _, err := api.UpdateBillingSubscription(ctx, json.RawMessage(`{"priceID":"price_business","plan_tier":"business"}`)); err != nil {
		t.Fatal(err)
	}
}

func assertBillingBody(t *testing.T, got, want map[string]any) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("request body = %#v, want %#v", got, want)
	}
}
