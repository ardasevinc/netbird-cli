package analysis

import (
	"context"
	"sort"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/netbird"
)

// Reader is the read-only management surface required to explain a peer's
// current reachability. The accessible-peer endpoint remains authoritative;
// policy matches are explanatory evidence only.
type Reader interface {
	GetPeer(context.Context, string) (netbird.Peer, error)
	ListPeers(context.Context, string, string) ([]netbird.Peer, error)
	ListAccessiblePeers(context.Context, string) ([]netbird.Peer, error)
	ListPolicies(context.Context) ([]netbird.Policy, error)
}

type PolicyEvidence struct {
	ReachablePeerID string `json:"reachable_peer_id"`
	PolicyID        string `json:"policy_id"`
	PolicyName      string `json:"policy_name"`
	RuleID          string `json:"rule_id"`
	RuleName        string `json:"rule_name"`
	Action          string `json:"action"`
	Protocol        string `json:"protocol"`
	Direction       string `json:"direction"`
}

type Summary struct {
	ReachablePeerCount            int `json:"reachable_peer_count"`
	PolicyEvidenceCount           int `json:"policy_evidence_count"`
	UnexplainedReachablePeerCount int `json:"unexplained_reachable_peer_count"`
}

type Report struct {
	SourcePeer                  netbird.Peer     `json:"source_peer"`
	ReachablePeers              []netbird.Peer   `json:"reachable_peers"`
	PolicyEvidence              []PolicyEvidence `json:"policy_evidence"`
	UnexplainedReachablePeerIDs []string         `json:"unexplained_reachable_peer_ids"`
	Summary                     Summary          `json:"summary"`
	Completeness                map[string]any   `json:"completeness"`
}

func Reachability(ctx context.Context, reader Reader, peerID string) (Report, error) {
	source, err := reader.GetPeer(ctx, peerID)
	if err != nil {
		return Report{}, err
	}
	reachable, err := reader.ListAccessiblePeers(ctx, peerID)
	if err != nil {
		return Report{}, err
	}
	inventory, err := reader.ListPeers(ctx, "", "")
	if err != nil {
		return Report{}, err
	}
	policies, err := reader.ListPolicies(ctx)
	if err != nil {
		return Report{}, err
	}

	reachable = joinPeerGroups(reachable, inventory)
	sort.Slice(reachable, func(i, j int) bool { return reachable[i].ID < reachable[j].ID })
	evidence := policyEvidence(source, reachable, policies)
	unexplained := unexplainedPeers(reachable, evidence)
	return Report{
		SourcePeer:                  source,
		ReachablePeers:              reachable,
		PolicyEvidence:              evidence,
		UnexplainedReachablePeerIDs: unexplained,
		Summary: Summary{
			ReachablePeerCount:            len(reachable),
			PolicyEvidenceCount:           len(evidence),
			UnexplainedReachablePeerCount: len(unexplained),
		},
		Completeness: map[string]any{"state": "complete", "reason": nil},
	}, nil
}

func joinPeerGroups(reachable, inventory []netbird.Peer) []netbird.Peer {
	groupsByID := make(map[string][]netbird.PeerGroup, len(inventory))
	for _, peer := range inventory {
		groupsByID[peer.ID] = peer.Groups
	}
	joined := make([]netbird.Peer, len(reachable))
	copy(joined, reachable)
	for i, peer := range joined {
		if groups, ok := groupsByID[peer.ID]; ok {
			peer.Groups = groups
			joined[i] = peer
		}
	}
	return joined
}

func policyEvidence(source netbird.Peer, reachable []netbird.Peer, policies []netbird.Policy) []PolicyEvidence {
	result := make([]PolicyEvidence, 0)
	for _, peer := range reachable {
		for _, policy := range policies {
			for _, rule := range policy.Rules {
				if !policy.Enabled || !rule.Enabled {
					continue
				}
				if groupsIntersect(source.Groups, rule.Sources) && groupsIntersect(peer.Groups, rule.Destinations) {
					result = append(result, evidence(policy, rule, peer.ID, "source_to_destination"))
				}
				if rule.Bidirectional && groupsIntersect(source.Groups, rule.Destinations) && groupsIntersect(peer.Groups, rule.Sources) {
					result = append(result, evidence(policy, rule, peer.ID, "destination_to_source"))
				}
			}
		}
	}
	sort.Slice(result, func(i, j int) bool {
		left := result[i]
		right := result[j]
		for _, pair := range [][2]string{{left.ReachablePeerID, right.ReachablePeerID}, {left.PolicyID, right.PolicyID}, {left.RuleID, right.RuleID}, {left.Direction, right.Direction}} {
			if pair[0] != pair[1] {
				return pair[0] < pair[1]
			}
		}
		return left.PolicyName+left.RuleName < right.PolicyName+right.RuleName
	})
	return result
}

func evidence(policy netbird.Policy, rule netbird.PolicyRule, peerID, direction string) PolicyEvidence {
	return PolicyEvidence{
		ReachablePeerID: peerID,
		PolicyID:        pointerValue(policy.ID),
		PolicyName:      policy.Name,
		RuleID:          pointerValue(rule.ID),
		RuleName:        rule.Name,
		Action:          rule.Action,
		Protocol:        rule.Protocol,
		Direction:       direction,
	}
}

func groupsIntersect(peers []netbird.PeerGroup, policyGroups []netbird.PolicyGroup) bool {
	for _, peerGroup := range peers {
		for _, policyGroup := range policyGroups {
			if peerGroup.ID != "" && peerGroup.ID == policyGroup.ID {
				return true
			}
		}
	}
	return false
}

func unexplainedPeers(reachable []netbird.Peer, evidence []PolicyEvidence) []string {
	explained := make(map[string]struct{}, len(evidence))
	for _, item := range evidence {
		if strings.EqualFold(item.Action, "accept") {
			explained[item.ReachablePeerID] = struct{}{}
		}
	}
	result := make([]string, 0)
	for _, peer := range reachable {
		if _, ok := explained[peer.ID]; !ok {
			result = append(result, peer.ID)
		}
	}
	return result
}

func pointerValue(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
