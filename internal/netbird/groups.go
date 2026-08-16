package netbird

import (
	"context"
)

type Group struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	PeersCount     int    `json:"peers_count"`
	ResourcesCount int    `json:"resources_count"`
}

func (c *Client) ListGroups(ctx context.Context) ([]Group, error) {
	var groups []Group
	if err := c.transport.GetJSON(ctx, "/api/groups", &groups); err != nil {
		return nil, err
	}
	return groups, nil
}
