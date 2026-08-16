package netbird

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type AuditEvent struct {
	ID             string            `json:"id"`
	Activity       string            `json:"activity"`
	ActivityCode   string            `json:"activity_code"`
	InitiatorEmail string            `json:"initiator_email"`
	InitiatorID    string            `json:"initiator_id"`
	InitiatorName  string            `json:"initiator_name"`
	Meta           map[string]string `json:"meta"`
	TargetID       string            `json:"target_id"`
	Timestamp      string            `json:"timestamp"`
}

func (c *Client) ListAuditEvents(ctx context.Context) ([]AuditEvent, error) {
	var result []AuditEvent
	if err := c.transport.GetJSON(ctx, "/api/events/audit", &result); err != nil {
		return nil, fmt.Errorf("list audit events: %w", err)
	}
	return result, nil
}

type NetworkTrafficEndpoint struct {
	Address     string                 `json:"address"`
	DNSLabel    *string                `json:"dns_label"`
	GeoLocation NetworkTrafficLocation `json:"geo_location"`
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	OS          *string                `json:"os"`
	Type        string                 `json:"type"`
}

type NetworkTrafficLocation struct {
	CityName    string `json:"city_name"`
	CountryCode string `json:"country_code"`
}

type NetworkTrafficSubEvent struct {
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"`
}

type NetworkTrafficPolicy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type NetworkTrafficUser struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type NetworkTrafficICMP struct {
	Code int `json:"code"`
	Type int `json:"type"`
}

type NetworkTrafficEvent struct {
	Destination NetworkTrafficEndpoint   `json:"destination"`
	Direction   string                   `json:"direction"`
	Events      []NetworkTrafficSubEvent `json:"events"`
	FlowID      string                   `json:"flow_id"`
	ICMP        NetworkTrafficICMP       `json:"icmp"`
	NumOfDrops  int                      `json:"num_of_drops"`
	NumOfEnds   int                      `json:"num_of_ends"`
	NumOfStarts int                      `json:"num_of_starts"`
	Policy      NetworkTrafficPolicy     `json:"policy"`
	Protocol    int                      `json:"protocol"`
	ReporterID  string                   `json:"reporter_id"`
	RxBytes     int                      `json:"rx_bytes"`
	RxPackets   int                      `json:"rx_packets"`
	Source      NetworkTrafficEndpoint   `json:"source"`
	TxBytes     int                      `json:"tx_bytes"`
	TxPackets   int                      `json:"tx_packets"`
	User        NetworkTrafficUser       `json:"user"`
	WindowEnd   time.Time                `json:"window_end"`
	WindowStart time.Time                `json:"window_start"`
}

type NetworkTrafficPage struct {
	Data         []NetworkTrafficEvent `json:"data"`
	Page         int                   `json:"page"`
	PageSize     int                   `json:"page_size"`
	TotalPages   int                   `json:"total_pages"`
	TotalRecords int                   `json:"total_records"`
}

type ProxyAccessLog struct {
	AuthMethodUsed  *string            `json:"auth_method_used,omitempty"`
	BytesDownload   int64              `json:"bytes_download"`
	BytesUpload     int64              `json:"bytes_upload"`
	CityName        *string            `json:"city_name,omitempty"`
	CountryCode     *string            `json:"country_code,omitempty"`
	DurationMS      int                `json:"duration_ms"`
	Host            string             `json:"host"`
	ID              string             `json:"id"`
	Metadata        *map[string]string `json:"metadata,omitempty"`
	Method          string             `json:"method"`
	Path            string             `json:"path"`
	Protocol        *string            `json:"protocol,omitempty"`
	Reason          *string            `json:"reason,omitempty"`
	ServiceID       string             `json:"service_id"`
	SourceIP        *string            `json:"source_ip,omitempty"`
	StatusCode      int                `json:"status_code"`
	SubdivisionCode *string            `json:"subdivision_code,omitempty"`
	Timestamp       time.Time          `json:"timestamp"`
	UserID          *string            `json:"user_id,omitempty"`
}

type ProxyAccessLogsPage struct {
	Data         []ProxyAccessLog `json:"data"`
	Page         int              `json:"page"`
	PageSize     int              `json:"page_size"`
	TotalPages   int              `json:"total_pages"`
	TotalRecords int              `json:"total_records"`
}

type EventPageOptions struct {
	Page     int
	PageSize int
}

func (o EventPageOptions) query() string {
	values := url.Values{}
	if o.Page > 0 {
		values.Set("page", strconv.Itoa(o.Page))
	}
	if o.PageSize > 0 {
		values.Set("page_size", strconv.Itoa(o.PageSize))
	}
	return values.Encode()
}

func (c *Client) ListNetworkTrafficEvents(ctx context.Context, options EventPageOptions) (NetworkTrafficPage, error) {
	path := "/api/events/network-traffic"
	if query := options.query(); query != "" {
		path += "?" + query
	}
	var result NetworkTrafficPage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return NetworkTrafficPage{}, fmt.Errorf("list network traffic events: %w", err)
	}
	return result, nil
}

func (c *Client) ListProxyAccessLogs(ctx context.Context, options EventPageOptions) (ProxyAccessLogsPage, error) {
	path := "/api/events/proxy"
	if query := options.query(); query != "" {
		path += "?" + query
	}
	var result ProxyAccessLogsPage
	if _, err := c.transport.DoJSON(ctx, http.MethodGet, path, nil, &result); err != nil {
		return ProxyAccessLogsPage{}, fmt.Errorf("list proxy access logs: %w", err)
	}
	return result, nil
}
