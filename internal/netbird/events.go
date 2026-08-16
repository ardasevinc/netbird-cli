package netbird

import (
	"context"
	"fmt"
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
