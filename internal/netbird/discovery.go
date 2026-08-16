package netbird

import (
	"context"
	"errors"
	"fmt"
	"net/http"

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
	ID             string            `json:"id"`
	Settings       AccountSettings   `json:"settings"`
	Domain         string            `json:"domain"`
	DomainCategory string            `json:"domain_category"`
	CreatedAt      string            `json:"created_at"`
	CreatedBy      string            `json:"created_by"`
	Onboarding     AccountOnboarding `json:"onboarding"`
}

type AccountSettings struct {
	PeerLoginExpirationEnabled      bool                     `json:"peer_login_expiration_enabled"`
	PeerLoginExpiration             int                      `json:"peer_login_expiration"`
	PeerInactivityExpirationEnabled bool                     `json:"peer_inactivity_expiration_enabled"`
	PeerInactivityExpiration        int                      `json:"peer_inactivity_expiration"`
	RegularUsersViewBlocked         bool                     `json:"regular_users_view_blocked"`
	GroupsPropagationEnabled        bool                     `json:"groups_propagation_enabled"`
	JWTGroupsEnabled                bool                     `json:"jwt_groups_enabled"`
	JWTGroupsClaimName              string                   `json:"jwt_groups_claim_name"`
	JWTAllowGroups                  []string                 `json:"jwt_allow_groups"`
	RoutingPeerDNSResolutionEnabled bool                     `json:"routing_peer_dns_resolution_enabled"`
	DNSDomain                       string                   `json:"dns_domain"`
	NetworkRange                    string                   `json:"network_range"`
	NetworkRangeV6                  string                   `json:"network_range_v6"`
	PeerExposeEnabled               bool                     `json:"peer_expose_enabled"`
	PeerExposeGroups                []string                 `json:"peer_expose_groups"`
	Extra                           AccountExtraSettings     `json:"extra"`
	DashboardFeatures               AccountDashboardFeatures `json:"dashboard_features"`
	LazyConnectionEnabled           bool                     `json:"lazy_connection_enabled"`
	AutoUpdateVersion               string                   `json:"auto_update_version"`
	AutoUpdateAlways                bool                     `json:"auto_update_always"`
	MetricsPushEnabled              bool                     `json:"metrics_push_enabled"`
	AgentNetworkOnly                bool                     `json:"agent_network_only"`
	EmbeddedIDPEnabled              bool                     `json:"embedded_idp_enabled"`
	LocalAuthDisabled               bool                     `json:"local_auth_disabled"`
	LocalMFAEnabled                 bool                     `json:"local_mfa_enabled"`
	IPv6EnabledGroups               []string                 `json:"ipv6_enabled_groups"`
}

type AccountDashboardFeatures struct {
	AgentNetwork *bool `json:"agent_network,omitempty"`
}

type AccountExtraSettings struct {
	PeerApprovalEnabled                bool     `json:"peer_approval_enabled"`
	UserApprovalRequired               bool     `json:"user_approval_required"`
	NetworkTrafficLogsEnabled          bool     `json:"network_traffic_logs_enabled"`
	NetworkTrafficLogsGroups           []string `json:"network_traffic_logs_groups"`
	NetworkTrafficPacketCounterEnabled bool     `json:"network_traffic_packet_counter_enabled"`
}

type AccountOnboarding struct {
	SignupFormPending     bool `json:"signup_form_pending"`
	OnboardingFlowPending bool `json:"onboarding_flow_pending"`
}

type Discovery struct {
	Version        Version  `json:"version"`
	Instance       Instance `json:"instance"`
	User           *User    `json:"user,omitempty"`
	IdentityStatus string   `json:"identity_status"`
}

type Client struct {
	transport *transport.Client
}

func NewClient(transportClient *transport.Client) *Client {
	return &Client{transport: transportClient}
}

func (c *Client) ServerIdentity() string { return c.transport.Origin() }

func (c *Client) Discover(ctx context.Context, authenticated bool) (Discovery, error) {
	result := Discovery{IdentityStatus: "not_requested"}
	if err := c.transport.GetJSON(ctx, "/api/instance/version", &result.Version); err != nil {
		return Discovery{}, fmt.Errorf("discover server version: %w", err)
	}
	if err := c.transport.GetJSON(ctx, "/api/instance", &result.Instance); err != nil {
		return Discovery{}, fmt.Errorf("discover instance status: %w", err)
	}
	if authenticated {
		result.IdentityStatus = "available"
		var user User
		if _, err := c.transport.DoJSON(ctx, http.MethodGet, "/api/users/current", nil, &user); err != nil {
			var requestErr *transport.RequestError
			if errors.As(err, &requestErr) && requestErr.StatusCode == http.StatusForbidden {
				result.IdentityStatus = "unavailable"
				return result, nil
			}
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
