package analysis

import (
	"context"
	"encoding/json"
	"sort"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/netbird"
)

const (
	DirectionOutbound = "outbound"
	DirectionInbound  = "inbound"
)

type AccessFlow struct {
	PeerID                    string                    `json:"peer_id"`
	Direction                 string                    `json:"direction"`
	PolicyID                  string                    `json:"policy_id"`
	PolicyName                string                    `json:"policy_name"`
	RuleID                    string                    `json:"rule_id"`
	RuleName                  string                    `json:"rule_name"`
	Action                    string                    `json:"action"`
	Protocol                  string                    `json:"protocol"`
	Ports                     []string                  `json:"ports"`
	PortRanges                []netbird.PolicyPortRange `json:"port_ranges"`
	Bidirectional             bool                      `json:"bidirectional"`
	SourceMatch               string                    `json:"source_match"`
	DestinationMatch          string                    `json:"destination_match"`
	PolicySourcePostureChecks []string                  `json:"policy_source_posture_checks"`
}

type AccessRelation struct {
	Peer               netbird.Peer `json:"peer"`
	NetworkMapAdjacent bool         `json:"network_map_adjacent"`
	OutboundFlows      []AccessFlow `json:"outbound_flows"`
	InboundFlows       []AccessFlow `json:"inbound_flows"`
	HasOutboundFlows   bool         `json:"has_outbound_flows"`
	HasInboundFlows    bool         `json:"has_inbound_flows"`
	Directions         []string     `json:"directions"`
}

type ProjectionFinding struct {
	PolicyID   string `json:"policy_id"`
	PolicyName string `json:"policy_name"`
	RuleID     string `json:"rule_id"`
	RuleName   string `json:"rule_name"`
	Code       string `json:"code"`
	Detail     string `json:"detail"`
}

// AccessObservation is deliberately scoped to one concrete probe. A future
// observation collector can append records without turning one successful
// packet exchange into a claim about an entire peer relation.
type AccessObservation struct {
	SourcePeerID      string `json:"source_peer_id"`
	DestinationPeerID string `json:"destination_peer_id"`
	Protocol          string `json:"protocol"`
	Port              *int   `json:"port"`
	Timestamp         string `json:"timestamp"`
	Result            string `json:"result"`
	Origin            string `json:"origin"`
	Method            string `json:"method"`
}

type AccessCompleteness struct {
	NetworkMap       Completeness `json:"network_map"`
	PolicyProjection Completeness `json:"policy_projection"`
	Observations     Completeness `json:"observations"`
}

type Completeness struct {
	State  string  `json:"state"`
	Reason *string `json:"reason"`
}

type AccessSummary struct {
	NetworkMapPeerCount                      int `json:"network_map_peer_count"`
	RelationCount                            int `json:"relation_count"`
	OutboundFlowCount                        int `json:"outbound_flow_count"`
	InboundFlowCount                         int `json:"inbound_flow_count"`
	MapOnlyPeerCount                         int `json:"map_only_peer_count"`
	ConfiguredPeerMissingFromNetworkMapCount int `json:"configured_peer_missing_from_network_map_count"`
	ObservationCount                         int `json:"observation_count"`
}

type AccessReport struct {
	SubjectPeer                            netbird.Peer        `json:"subject_peer"`
	NetworkMapPeers                        []netbird.Peer      `json:"network_map_peers"`
	Relations                              []AccessRelation    `json:"relations"`
	MapOnlyPeerIDs                         []string            `json:"map_only_peer_ids"`
	ConfiguredPeerMissingFromNetworkMapIDs []string            `json:"configured_peer_missing_from_network_map_ids"`
	ProjectionFindings                     []ProjectionFinding `json:"projection_findings"`
	Observations                           []AccessObservation `json:"observations"`
	Summary                                AccessSummary       `json:"summary"`
	Completeness                           AccessCompleteness  `json:"completeness"`
}

func Access(ctx context.Context, reader Reader, peerID string) (AccessReport, error) {
	subject, err := reader.GetPeer(ctx, peerID)
	if err != nil {
		return AccessReport{}, err
	}
	networkMapPeers, err := reader.ListAccessiblePeers(ctx, peerID)
	if err != nil {
		return AccessReport{}, err
	}
	inventory, err := reader.ListPeers(ctx, "", "")
	if err != nil {
		return AccessReport{}, err
	}
	policies, err := reader.ListPolicies(ctx)
	if err != nil {
		return AccessReport{}, err
	}

	networkMapPeers = joinPeerGroups(networkMapPeers, inventory)
	sort.Slice(networkMapPeers, func(i, j int) bool { return networkMapPeers[i].ID < networkMapPeers[j].ID })
	flows, findings := configuredPeerFlows(subject, inventory, policies)
	relations, mapOnly, missing := accessRelations(networkMapPeers, inventory, flows)

	projectionState := "complete"
	var projectionReason *string
	if len(findings) > 0 {
		projectionState = "partial"
		reason := "one or more enabled peer-relevant rules use unsupported action semantics"
		projectionReason = &reason
	}
	observationReason := "no packet observations were supplied"
	report := AccessReport{
		SubjectPeer:                            subject,
		NetworkMapPeers:                        networkMapPeers,
		Relations:                              relations,
		MapOnlyPeerIDs:                         mapOnly,
		ConfiguredPeerMissingFromNetworkMapIDs: missing,
		ProjectionFindings:                     findings,
		Observations:                           []AccessObservation{},
		Completeness: AccessCompleteness{
			NetworkMap:       Completeness{State: "complete"},
			PolicyProjection: Completeness{State: projectionState, Reason: projectionReason},
			Observations:     Completeness{State: "unknown", Reason: &observationReason},
		},
	}
	report.Summary = summarizeAccess(report)
	return report, nil
}

func configuredPeerFlows(subject netbird.Peer, inventory []netbird.Peer, policies []netbird.Policy) ([]AccessFlow, []ProjectionFinding) {
	flows := make([]AccessFlow, 0)
	findings := make([]ProjectionFinding, 0)
	for _, policy := range policies {
		if !policy.Enabled {
			continue
		}
		for _, rule := range policy.Rules {
			if !rule.Enabled {
				continue
			}
			subjectSource, subjectSourceMatch, sourceSupported := peerMatchesPolicySide(subject, rule.Sources, rule.SourceResource)
			subjectDestination, subjectDestinationMatch, destinationSupported := peerMatchesPolicySide(subject, rule.Destinations, rule.DestinationResource)
			if !subjectSource && !subjectDestination {
				continue
			}
			if !strings.EqualFold(rule.Action, "accept") {
				findings = append(findings, ProjectionFinding{
					PolicyID: pointerValue(policy.ID), PolicyName: policy.Name,
					RuleID: pointerValue(rule.ID), RuleName: rule.Name,
					Code: "unsupported_action", Detail: "enabled action " + rule.Action + " is not projected as deny or allow",
				})
				continue
			}
			if !sourceSupported || !destinationSupported {
				continue
			}
			for _, peer := range inventory {
				if peer.ID == "" || peer.ID == subject.ID {
					continue
				}
				peerSource, peerSourceMatch, peerSourceSupported := peerMatchesPolicySide(peer, rule.Sources, rule.SourceResource)
				peerDestination, peerDestinationMatch, peerDestinationSupported := peerMatchesPolicySide(peer, rule.Destinations, rule.DestinationResource)
				if !peerSourceSupported || !peerDestinationSupported {
					continue
				}
				if subjectSource && peerDestination {
					flows = append(flows, accessFlow(policy, rule, peer.ID, DirectionOutbound, subjectSourceMatch, peerDestinationMatch))
					if rule.Bidirectional {
						flows = append(flows, accessFlow(policy, rule, peer.ID, DirectionInbound, peerDestinationMatch, subjectSourceMatch))
					}
				}
				if subjectDestination && peerSource {
					flows = append(flows, accessFlow(policy, rule, peer.ID, DirectionInbound, peerSourceMatch, subjectDestinationMatch))
					if rule.Bidirectional {
						flows = append(flows, accessFlow(policy, rule, peer.ID, DirectionOutbound, subjectDestinationMatch, peerSourceMatch))
					}
				}
			}
		}
	}
	flows = deduplicateAccessFlows(flows)
	sortAccessFlows(flows)
	sort.Slice(findings, func(i, j int) bool {
		return findings[i].PolicyID+findings[i].RuleID+findings[i].Code < findings[j].PolicyID+findings[j].RuleID+findings[j].Code
	})
	return flows, findings
}

func peerMatchesPolicySide(peer netbird.Peer, groups []netbird.PolicyGroup, resource *netbird.PolicyResource) (bool, string, bool) {
	if resource != nil && resource.ID != "" {
		if !strings.EqualFold(resource.Type, "peer") {
			return false, "", false
		}
		return peer.ID == resource.ID, "peer_resource", true
	}
	return groupsIntersect(peer.Groups, groups), "group", true
}

func accessFlow(policy netbird.Policy, rule netbird.PolicyRule, peerID, direction, sourceMatch, destinationMatch string) AccessFlow {
	ports := append([]string(nil), rule.Ports...)
	if ports == nil {
		ports = []string{}
	}
	portRanges := append([]netbird.PolicyPortRange(nil), rule.PortRanges...)
	if portRanges == nil {
		portRanges = []netbird.PolicyPortRange{}
	}
	posture := append([]string(nil), policy.SourcePostureChecks...)
	if posture == nil {
		posture = []string{}
	}
	return AccessFlow{
		PeerID: peerID, Direction: direction,
		PolicyID: pointerValue(policy.ID), PolicyName: policy.Name,
		RuleID: pointerValue(rule.ID), RuleName: rule.Name,
		Action: strings.ToLower(rule.Action), Protocol: rule.Protocol, Ports: ports, PortRanges: portRanges,
		Bidirectional: rule.Bidirectional, SourceMatch: sourceMatch, DestinationMatch: destinationMatch,
		PolicySourcePostureChecks: posture,
	}
}

func deduplicateAccessFlows(flows []AccessFlow) []AccessFlow {
	unique := make([]AccessFlow, 0, len(flows))
	seen := make(map[string]struct{}, len(flows))
	for _, flow := range flows {
		encoded, _ := json.Marshal(flow)
		key := string(encoded)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		unique = append(unique, flow)
	}
	return unique
}

func accessRelations(networkMapPeers, inventory []netbird.Peer, flows []AccessFlow) ([]AccessRelation, []string, []string) {
	peers := make(map[string]netbird.Peer, len(inventory)+len(networkMapPeers))
	for _, peer := range inventory {
		peers[peer.ID] = peer
	}
	adjacent := make(map[string]bool, len(networkMapPeers))
	for _, peer := range networkMapPeers {
		peers[peer.ID] = peer
		adjacent[peer.ID] = true
	}
	flowsByPeer := make(map[string][]AccessFlow)
	for _, flow := range flows {
		flowsByPeer[flow.PeerID] = append(flowsByPeer[flow.PeerID], flow)
	}
	ids := make([]string, 0, len(adjacent)+len(flowsByPeer))
	seen := make(map[string]struct{}, len(adjacent)+len(flowsByPeer))
	for id := range adjacent {
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	for id := range flowsByPeer {
		if _, ok := seen[id]; !ok {
			ids = append(ids, id)
		}
	}
	sort.Strings(ids)
	relations := make([]AccessRelation, 0, len(ids))
	mapOnly := make([]string, 0)
	missing := make([]string, 0)
	for _, id := range ids {
		outbound := make([]AccessFlow, 0)
		inbound := make([]AccessFlow, 0)
		for _, flow := range flowsByPeer[id] {
			if flow.Direction == DirectionOutbound {
				outbound = append(outbound, flow)
			} else {
				inbound = append(inbound, flow)
			}
		}
		directions := make([]string, 0, 2)
		if len(outbound) > 0 {
			directions = append(directions, DirectionOutbound)
		}
		if len(inbound) > 0 {
			directions = append(directions, DirectionInbound)
		}
		relation := AccessRelation{
			Peer: peers[id], NetworkMapAdjacent: adjacent[id],
			OutboundFlows: outbound, InboundFlows: inbound,
			HasOutboundFlows: len(outbound) > 0, HasInboundFlows: len(inbound) > 0,
			Directions: directions,
		}
		relations = append(relations, relation)
		if adjacent[id] && len(outbound)+len(inbound) == 0 {
			mapOnly = append(mapOnly, id)
		}
		if !adjacent[id] && len(outbound)+len(inbound) > 0 {
			missing = append(missing, id)
		}
	}
	return relations, mapOnly, missing
}

func sortAccessFlows(flows []AccessFlow) {
	sort.Slice(flows, func(i, j int) bool {
		left, right := flows[i], flows[j]
		for _, pair := range [][2]string{{left.PeerID, right.PeerID}, {left.Direction, right.Direction}, {left.PolicyID, right.PolicyID}, {left.RuleID, right.RuleID}} {
			if pair[0] != pair[1] {
				return pair[0] < pair[1]
			}
		}
		return left.Protocol+strings.Join(left.Ports, ",") < right.Protocol+strings.Join(right.Ports, ",")
	})
}

func summarizeAccess(report AccessReport) AccessSummary {
	summary := AccessSummary{
		NetworkMapPeerCount: len(report.NetworkMapPeers), RelationCount: len(report.Relations),
		MapOnlyPeerCount:                         len(report.MapOnlyPeerIDs),
		ConfiguredPeerMissingFromNetworkMapCount: len(report.ConfiguredPeerMissingFromNetworkMapIDs),
		ObservationCount:                         len(report.Observations),
	}
	for _, relation := range report.Relations {
		summary.OutboundFlowCount += len(relation.OutboundFlows)
		summary.InboundFlowCount += len(relation.InboundFlows)
	}
	return summary
}
