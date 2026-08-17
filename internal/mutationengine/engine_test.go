package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"testing"

	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/transport"
)

type fakeRemote struct {
	identity               string
	account                string
	accountBefore          json.RawMessage
	accountAfter           json.RawMessage
	postureCollection      json.RawMessage
	postureBefore          json.RawMessage
	postureAfter           json.RawMessage
	ingressCollection      json.RawMessage
	ingressBefore          json.RawMessage
	ingressAfter           json.RawMessage
	agentSettingsBefore    json.RawMessage
	agentSettingsAfter     json.RawMessage
	budgetCollection       json.RawMessage
	budgetBefore           json.RawMessage
	budgetAfter            json.RawMessage
	guardrailCollection    json.RawMessage
	guardrailBefore        json.RawMessage
	guardrailAfter         json.RawMessage
	agentPolicyCollection  json.RawMessage
	agentPolicyBefore      json.RawMessage
	agentPolicyAfter       json.RawMessage
	providerCollection     json.RawMessage
	providerBefore         json.RawMessage
	providerAfter          json.RawMessage
	providerSecretSeen     bool
	userCollection         json.RawMessage
	userBefore             json.RawMessage
	userAfter              json.RawMessage
	passwordSeen           bool
	passwordBody           json.RawMessage
	inviteResendSeen       bool
	tokenBefore            json.RawMessage
	tokenCollection        json.RawMessage
	tokenAfter             json.RawMessage
	setupKeyBefore         json.RawMessage
	setupKeyCollection     json.RawMessage
	setupKeyAfter          json.RawMessage
	setupKeyBody           json.RawMessage
	inviteBefore           json.RawMessage
	inviteCollection       json.RawMessage
	inviteAfter            json.RawMessage
	publicInviteBefore     json.RawMessage
	publicInviteToken      string
	publicInviteAccepted   bool
	publicInviteBody       json.RawMessage
	ingressPortCollection  json.RawMessage
	ingressPortBefore      json.RawMessage
	ingressPortAfter       json.RawMessage
	ingressPortBody        json.RawMessage
	before                 json.RawMessage
	after                  json.RawMessage
	groupCollection        json.RawMessage
	policyBefore           json.RawMessage
	policyAfter            json.RawMessage
	policyCollection       json.RawMessage
	routeBefore            json.RawMessage
	routeAfter             json.RawMessage
	routeCollection        json.RawMessage
	peerBefore             json.RawMessage
	peerAfter              json.RawMessage
	eventStreamingBefore   json.RawMessage
	eventStreamingAfter    json.RawMessage
	eventStreamingList     json.RawMessage
	eventStreamingBody     json.RawMessage
	identityProviderBefore json.RawMessage
	identityProviderAfter  json.RawMessage
	identityProviderList   json.RawMessage
	identityProviderBody   json.RawMessage
	proxyTokenBefore       json.RawMessage
	proxyTokenCollection   json.RawMessage
	proxyTokenAfter        json.RawMessage
	proxyTokenBody         json.RawMessage
	proxyDomainBefore      json.RawMessage
	proxyDomainCollection  json.RawMessage
	proxyDomainAfter       json.RawMessage
	proxyDomainBody        json.RawMessage
	proxyClusterCollection json.RawMessage
	proxyServiceBefore     json.RawMessage
	proxyServiceCollection json.RawMessage
	proxyServiceAfter      json.RawMessage
	proxyServiceBody       json.RawMessage
	temporaryAccessBody    json.RawMessage
	temporaryAccessResult  json.RawMessage
	edrBypassed            json.RawMessage
	networkBefore          json.RawMessage
	networkAfter           json.RawMessage
	networkCollection      json.RawMessage
	dnsZoneCollection      json.RawMessage
	nameserverCollection   json.RawMessage
	nameserverBefore       json.RawMessage
	nameserverAfter        json.RawMessage
	dnsSettingsBefore      json.RawMessage
	dnsSettingsAfter       json.RawMessage
	resourceBefore         json.RawMessage
	resourceAfter          json.RawMessage
	resourceCollection     json.RawMessage
	routerBefore           json.RawMessage
	routerAfter            json.RawMessage
	routerCollection       json.RawMessage
	updateErr              error
	updates                int
	notification           notificationChannelState
	azureIDP               azureIDPState
	googleIDP              googleIDPState
	edr                    map[string]edrIntegrationState
	scim                   map[string]scimIntegrationState
}

type notificationChannelState struct {
	before     json.RawMessage
	collection json.RawMessage
	after      json.RawMessage
	body       json.RawMessage
}

type azureIDPState struct {
	before     json.RawMessage
	collection json.RawMessage
	after      json.RawMessage
	body       json.RawMessage
}

type googleIDPState struct {
	before     json.RawMessage
	collection json.RawMessage
	after      json.RawMessage
	body       json.RawMessage
}

type edrIntegrationState struct {
	before json.RawMessage
	after  json.RawMessage
	body   json.RawMessage
}

type scimIntegrationState struct {
	before     json.RawMessage
	collection json.RawMessage
	after      json.RawMessage
	token      string
}

func (f *fakeRemote) ServerIdentity() string { return f.identity }

func (f *fakeRemote) AccountScope(_ context.Context, account string) error {
	if account != f.account {
		return errors.New("wrong account")
	}
	return nil
}

func (f *fakeRemote) GetAccountRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.accountBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.accountBefore...), nil
}

func (f *fakeRemote) UpdateAccount(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.accountBefore = append(json.RawMessage(nil), f.accountAfter...)
	return append(json.RawMessage(nil), f.accountAfter...), nil
}

func (f *fakeRemote) DeleteAccount(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.accountBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListPostureChecksRaw(_ context.Context) (json.RawMessage, error) {
	if f.postureCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.postureCollection...), nil
}

func (f *fakeRemote) CreatePostureCheck(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created posture check")
	}
	f.postureCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetPostureCheckRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.postureBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.postureBefore...), nil
}

func (f *fakeRemote) UpdatePostureCheck(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.postureBefore = append(json.RawMessage(nil), f.postureAfter...)
	return append(json.RawMessage(nil), f.postureAfter...), nil
}

func (f *fakeRemote) DeletePostureCheck(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.postureBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListIngressPeersRaw(_ context.Context) (json.RawMessage, error) {
	if f.ingressCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.ingressCollection...), nil
}

func (f *fakeRemote) CreateIngressPeer(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created ingress peer")
	}
	f.ingressCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetIngressPeerRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.ingressBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.ingressBefore...), nil
}

func (f *fakeRemote) UpdateIngressPeer(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.ingressBefore = append(json.RawMessage(nil), f.ingressAfter...)
	return append(json.RawMessage(nil), f.ingressAfter...), nil
}

func (f *fakeRemote) DeleteIngressPeer(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.ingressBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListIngressPortAllocationsRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.ingressPortCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.ingressPortCollection...), nil
}

func (f *fakeRemote) GetIngressPortAllocationRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.ingressPortBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.ingressPortBefore...), nil
}

func (f *fakeRemote) CreateIngressPortAllocation(_ context.Context, _ string, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.ingressPortAfter == nil {
		return nil, errors.New("missing created ingress port allocation")
	}
	f.ingressPortBody = append(json.RawMessage(nil), request...)
	f.ingressPortCollection = json.RawMessage("[" + string(f.ingressPortAfter) + "]")
	return append(json.RawMessage(nil), f.ingressPortAfter...), nil
}

func (f *fakeRemote) UpdateIngressPortAllocation(_ context.Context, _, _ string, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.ingressPortBefore = append(json.RawMessage(nil), f.ingressPortAfter...)
	f.ingressPortBody = append(json.RawMessage(nil), request...)
	return append(json.RawMessage(nil), f.ingressPortAfter...), nil
}

func (f *fakeRemote) DeleteIngressPortAllocation(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.ingressPortBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetAgentNetworkSettingsRaw(_ context.Context) (json.RawMessage, error) {
	if f.agentSettingsBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.agentSettingsBefore...), nil
}

func (f *fakeRemote) UpdateAgentNetworkSettings(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.agentSettingsBefore = append(json.RawMessage(nil), f.agentSettingsAfter...)
	return append(json.RawMessage(nil), f.agentSettingsAfter...), nil
}

func (f *fakeRemote) CreateAgentNetworkSettings(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.agentSettingsBefore = append(json.RawMessage(nil), f.agentSettingsAfter...)
	return append(json.RawMessage(nil), f.agentSettingsAfter...), nil
}

func (f *fakeRemote) DeleteAgentNetworkSettings(_ context.Context) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.agentSettingsBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListAgentNetworkBudgetRulesRaw(_ context.Context) (json.RawMessage, error) {
	if f.budgetCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.budgetCollection...), nil
}

func (f *fakeRemote) GetAgentNetworkBudgetRuleRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.budgetBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.budgetBefore...), nil
}

func (f *fakeRemote) CreateAgentNetworkBudgetRule(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.budgetAfter == nil {
		return nil, errors.New("missing created agent-network budget rule")
	}
	f.budgetCollection = json.RawMessage("[" + string(f.budgetAfter) + "]")
	return append(json.RawMessage(nil), f.budgetAfter...), nil
}

func (f *fakeRemote) UpdateAgentNetworkBudgetRule(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.budgetBefore = append(json.RawMessage(nil), f.budgetAfter...)
	return append(json.RawMessage(nil), f.budgetAfter...), nil
}

func (f *fakeRemote) DeleteAgentNetworkBudgetRule(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.budgetBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListAgentNetworkGuardrailsRaw(_ context.Context) (json.RawMessage, error) {
	if f.guardrailCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.guardrailCollection...), nil
}

func (f *fakeRemote) GetAgentNetworkGuardrailRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.guardrailBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.guardrailBefore...), nil
}

func (f *fakeRemote) CreateAgentNetworkGuardrail(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.guardrailAfter == nil {
		return nil, errors.New("missing created agent-network guardrail")
	}
	f.guardrailCollection = json.RawMessage("[" + string(f.guardrailAfter) + "]")
	return append(json.RawMessage(nil), f.guardrailAfter...), nil
}

func (f *fakeRemote) UpdateAgentNetworkGuardrail(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.guardrailBefore = append(json.RawMessage(nil), f.guardrailAfter...)
	return append(json.RawMessage(nil), f.guardrailAfter...), nil
}

func (f *fakeRemote) DeleteAgentNetworkGuardrail(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.guardrailBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListAgentNetworkPoliciesRaw(_ context.Context) (json.RawMessage, error) {
	if f.agentPolicyCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.agentPolicyCollection...), nil
}

func (f *fakeRemote) GetAgentNetworkPolicyRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.agentPolicyBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.agentPolicyBefore...), nil
}

func (f *fakeRemote) CreateAgentNetworkPolicy(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.agentPolicyAfter == nil {
		return nil, errors.New("missing created agent-network policy")
	}
	f.agentPolicyCollection = json.RawMessage("[" + string(f.agentPolicyAfter) + "]")
	return append(json.RawMessage(nil), f.agentPolicyAfter...), nil
}

func (f *fakeRemote) UpdateAgentNetworkPolicy(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.agentPolicyBefore = append(json.RawMessage(nil), f.agentPolicyAfter...)
	return append(json.RawMessage(nil), f.agentPolicyAfter...), nil
}

func (f *fakeRemote) DeleteAgentNetworkPolicy(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.agentPolicyBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListAgentNetworkProvidersRaw(_ context.Context) (json.RawMessage, error) {
	if f.providerCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.providerCollection...), nil
}

func (f *fakeRemote) GetAgentNetworkProviderRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.providerBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.providerBefore...), nil
}

func (f *fakeRemote) CreateAgentNetworkProvider(_ context.Context, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	var object map[string]any
	if err := json.Unmarshal(request, &object); err == nil {
		if secret, ok := object["api_key"].(string); ok && secret != "" {
			f.providerSecretSeen = true
		}
	}
	if f.providerAfter == nil {
		return nil, errors.New("missing created agent-network provider")
	}
	f.providerCollection = json.RawMessage("[" + string(f.providerAfter) + "]")
	return append(json.RawMessage(nil), f.providerAfter...), nil
}

func (f *fakeRemote) UpdateAgentNetworkProvider(_ context.Context, _ string, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	var object map[string]any
	if err := json.Unmarshal(request, &object); err == nil {
		if secret, ok := object["api_key"].(string); ok && secret != "" {
			f.providerSecretSeen = true
		}
	}
	f.providerBefore = append(json.RawMessage(nil), f.providerAfter...)
	return append(json.RawMessage(nil), f.providerAfter...), nil
}

func (f *fakeRemote) DeleteAgentNetworkProvider(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.providerBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListUsersRaw(_ context.Context) (json.RawMessage, error) {
	if f.userCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.userCollection...), nil
}

func (f *fakeRemote) GetUserRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.userBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.userBefore...), nil
}

func (f *fakeRemote) CreateUser(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.userAfter == nil {
		return nil, errors.New("missing created user")
	}
	f.userCollection = json.RawMessage("[" + string(f.userAfter) + "]")
	return append(json.RawMessage(nil), f.userAfter...), nil
}

func (f *fakeRemote) UpdateUser(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.userBefore = append(json.RawMessage(nil), f.userAfter...)
	return append(json.RawMessage(nil), f.userAfter...), nil
}

func (f *fakeRemote) DeleteUser(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.userBefore = nil
	return nil, nil
}

func (f *fakeRemote) ApproveUser(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.userBefore = append(json.RawMessage(nil), f.userAfter...)
	return append(json.RawMessage(nil), f.userAfter...), nil
}

func (f *fakeRemote) RejectUser(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.userBefore = nil
	return nil, nil
}

func (f *fakeRemote) ChangeUserPassword(_ context.Context, _ string, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.passwordSeen = true
	f.passwordBody = append(json.RawMessage(nil), request...)
	return json.RawMessage(`{"success":true}`), nil
}

func (f *fakeRemote) ResendUserInvite(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.inviteResendSeen = true
	return nil, nil
}

func (f *fakeRemote) GetPersonalAccessTokenRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.tokenBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.tokenBefore...), nil
}

func (f *fakeRemote) DeletePersonalAccessToken(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.tokenBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListPersonalAccessTokensRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.tokenCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.tokenCollection...), nil
}

func (f *fakeRemote) CreatePersonalAccessToken(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.tokenAfter == nil {
		return nil, errors.New("missing created token")
	}
	f.tokenCollection = json.RawMessage("[" + string(f.tokenAfter) + "]")
	return json.RawMessage(`{"plain_token":"one-time-secret","personal_access_token":` + string(f.tokenAfter) + `}`), nil
}

func (f *fakeRemote) GetSetupKeyRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.setupKeyBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.setupKeyBefore...), nil
}

func (f *fakeRemote) DeleteSetupKey(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.setupKeyBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListSetupKeysRaw(_ context.Context) (json.RawMessage, error) {
	if f.setupKeyCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.setupKeyCollection...), nil
}

func (f *fakeRemote) CreateSetupKey(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.setupKeyAfter == nil {
		return nil, errors.New("missing created setup key")
	}
	f.setupKeyCollection = json.RawMessage("[" + string(f.setupKeyAfter) + "]")
	return json.RawMessage(`{"id":"key-1","name":"bootstrap","key":"one-time-key"}`), nil
}

func (f *fakeRemote) UpdateSetupKey(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.setupKeyBody = append(json.RawMessage(nil), body...)
	if f.setupKeyAfter == nil {
		return nil, errors.New("missing updated setup key")
	}
	f.setupKeyBefore = append(json.RawMessage(nil), f.setupKeyAfter...)
	return append(json.RawMessage(nil), f.setupKeyAfter...), nil
}

func (f *fakeRemote) ListInvitesRaw(_ context.Context) (json.RawMessage, error) {
	if f.inviteCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.inviteCollection...), nil
}

func (f *fakeRemote) GetInviteRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.inviteBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.inviteBefore...), nil
}

func (f *fakeRemote) CreateInvite(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.inviteAfter == nil {
		return nil, errors.New("missing created invite")
	}
	f.inviteCollection = json.RawMessage("[" + string(f.inviteAfter) + "]")
	return json.RawMessage(`{"id":"invite-1","invite_token":"one-time-invite"}`), nil
}

func (f *fakeRemote) DeleteInvite(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.inviteBefore = nil
	return nil, nil
}

func (f *fakeRemote) RegenerateInvite(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.inviteAfter != nil {
		f.inviteBefore = append(json.RawMessage(nil), f.inviteAfter...)
	}
	return json.RawMessage(`{"invite_token":"replacement-invite"}`), nil
}

func (f *fakeRemote) GetPublicInviteRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.publicInviteBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.publicInviteBefore...), nil
}

func (f *fakeRemote) AcceptInvite(_ context.Context, token string, request json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.publicInviteAccepted = true
	f.publicInviteToken = token
	f.publicInviteBody = append(json.RawMessage(nil), request...)
	return json.RawMessage(`{"success":true}`), nil
}

func (f *fakeRemote) ListDNSZonesRaw(_ context.Context) (json.RawMessage, error) {
	if f.dnsZoneCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.dnsZoneCollection...), nil
}

func (f *fakeRemote) ListNameserverGroupsRaw(_ context.Context) (json.RawMessage, error) {
	if f.nameserverCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.nameserverCollection...), nil
}

func (f *fakeRemote) CreateNameserverGroup(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created nameserver group")
	}
	f.nameserverCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetNameserverGroupRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.nameserverBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.nameserverBefore...), nil
}

func (f *fakeRemote) UpdateNameserverGroup(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.nameserverBefore = append(json.RawMessage(nil), f.nameserverAfter...)
	return append(json.RawMessage(nil), f.nameserverAfter...), nil
}

func (f *fakeRemote) DeleteNameserverGroup(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.nameserverBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetDNSSettingsRaw(_ context.Context) (json.RawMessage, error) {
	if f.dnsSettingsBefore == nil {
		return json.RawMessage(`{"disabled_management_groups":[]}`), nil
	}
	return append(json.RawMessage(nil), f.dnsSettingsBefore...), nil
}

func (f *fakeRemote) UpdateDNSSettings(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.dnsSettingsBefore = append(json.RawMessage(nil), f.dnsSettingsAfter...)
	return append(json.RawMessage(nil), f.dnsSettingsAfter...), nil
}

func (f *fakeRemote) CreateDNSZone(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created dns zone")
	}
	f.dnsZoneCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetDNSZoneRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) UpdateDNSZone(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) DeleteDNSZone(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) ListDNSRecordsRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.dnsZoneCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.dnsZoneCollection...), nil
}

func (f *fakeRemote) CreateDNSRecord(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created dns record")
	}
	f.dnsZoneCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) GetDNSRecordRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) UpdateDNSRecord(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) DeleteDNSRecord(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) GetGroup(_ context.Context, _ string) (json.RawMessage, error) {
	if f.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.before...), nil
}

func (f *fakeRemote) ListGroupsRaw(_ context.Context) (json.RawMessage, error) {
	if f.groupCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.groupCollection...), nil
}

func (f *fakeRemote) CreateGroup(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.after == nil {
		return nil, errors.New("missing created group")
	}
	f.groupCollection = json.RawMessage("[" + string(f.after) + "]")
	return append(json.RawMessage(nil), f.after...), nil
}

func (f *fakeRemote) UpdateGroup(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = append(json.RawMessage(nil), f.after...)
	return f.before, nil
}

func (f *fakeRemote) DeleteGroup(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.before = nil
	return nil, nil
}

func (f *fakeRemote) GetPolicyRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.policyBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.policyBefore...), nil
}

func (f *fakeRemote) ListPoliciesRaw(_ context.Context) (json.RawMessage, error) {
	if f.policyCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.policyCollection...), nil
}

func (f *fakeRemote) CreatePolicy(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.policyAfter == nil {
		return nil, errors.New("missing created policy")
	}
	f.policyCollection = json.RawMessage("[" + string(f.policyAfter) + "]")
	return append(json.RawMessage(nil), f.policyAfter...), nil
}

func (f *fakeRemote) UpdatePolicy(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.policyBefore = append(json.RawMessage(nil), f.policyAfter...)
	return f.policyBefore, nil
}

func (f *fakeRemote) DeletePolicy(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.policyBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetRouteRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.routeBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.routeBefore...), nil
}

func (f *fakeRemote) ListRoutesRaw(_ context.Context) (json.RawMessage, error) {
	if f.routeCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.routeCollection...), nil
}

func (f *fakeRemote) CreateRoute(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.routeAfter == nil {
		return nil, errors.New("missing created route")
	}
	f.routeCollection = json.RawMessage("[" + string(f.routeAfter) + "]")
	return append(json.RawMessage(nil), f.routeAfter...), nil
}

func (f *fakeRemote) UpdateRoute(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routeBefore = append(json.RawMessage(nil), f.routeAfter...)
	return f.routeBefore, nil
}

func (f *fakeRemote) DeleteRoute(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routeBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetPeerRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.peerBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.peerBefore...), nil
}

func (f *fakeRemote) UpdatePeer(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.peerBefore = append(json.RawMessage(nil), f.peerAfter...)
	return f.peerBefore, nil
}

func (f *fakeRemote) CreateTemporaryAccessPeer(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.temporaryAccessBody = append(json.RawMessage(nil), body...)
	if f.temporaryAccessResult == nil {
		return nil, errors.New("missing temporary access result")
	}
	return append(json.RawMessage(nil), f.temporaryAccessResult...), nil
}

func (f *fakeRemote) CreatePeerJob(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.temporaryAccessBody = append(json.RawMessage(nil), body...)
	if f.temporaryAccessResult == nil {
		return nil, errors.New("missing peer job result")
	}
	return append(json.RawMessage(nil), f.temporaryAccessResult...), nil
}

func (f *fakeRemote) ListEventStreamingIntegrationsRaw(_ context.Context) (json.RawMessage, error) {
	if f.eventStreamingList == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.eventStreamingList...), nil
}

func (f *fakeRemote) GetEventStreamingIntegrationRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.eventStreamingBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.eventStreamingBefore...), nil
}

func (f *fakeRemote) CreateEventStreamingIntegration(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.eventStreamingBody = append(json.RawMessage(nil), body...)
	if f.eventStreamingAfter == nil {
		return nil, errors.New("missing event-streaming integration")
	}
	f.eventStreamingList = json.RawMessage("[" + string(f.eventStreamingAfter) + "]")
	return append(json.RawMessage(nil), f.eventStreamingAfter...), nil
}

func (f *fakeRemote) UpdateEventStreamingIntegration(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.eventStreamingBody = append(json.RawMessage(nil), body...)
	if f.eventStreamingAfter == nil {
		return nil, errors.New("missing event-streaming integration")
	}
	f.eventStreamingBefore = append(json.RawMessage(nil), f.eventStreamingAfter...)
	return append(json.RawMessage(nil), f.eventStreamingAfter...), nil
}

func (f *fakeRemote) DeleteEventStreamingIntegration(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.eventStreamingBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListIdentityProvidersRaw(_ context.Context) (json.RawMessage, error) {
	if f.identityProviderList == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.identityProviderList...), nil
}

func (f *fakeRemote) GetIdentityProviderRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.identityProviderBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.identityProviderBefore...), nil
}

func (f *fakeRemote) CreateIdentityProvider(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.identityProviderBody = append(json.RawMessage(nil), body...)
	if f.identityProviderAfter == nil {
		return nil, errors.New("missing identity provider")
	}
	f.identityProviderList = json.RawMessage("[" + string(f.identityProviderAfter) + "]")
	return append(json.RawMessage(nil), f.identityProviderAfter...), nil
}

func (f *fakeRemote) UpdateIdentityProvider(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.identityProviderBody = append(json.RawMessage(nil), body...)
	if f.identityProviderAfter == nil {
		return nil, errors.New("missing identity provider")
	}
	f.identityProviderBefore = append(json.RawMessage(nil), f.identityProviderAfter...)
	return append(json.RawMessage(nil), f.identityProviderAfter...), nil
}

func (f *fakeRemote) DeleteIdentityProvider(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.identityProviderBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListReverseProxyTokensRaw(_ context.Context) (json.RawMessage, error) {
	if f.proxyTokenCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.proxyTokenCollection...), nil
}

func (f *fakeRemote) CreateReverseProxyToken(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyTokenBody = append(json.RawMessage(nil), body...)
	if f.proxyTokenAfter == nil {
		return nil, errors.New("missing proxy token")
	}
	f.proxyTokenCollection = json.RawMessage("[" + string(f.proxyTokenAfter) + "]")
	return json.RawMessage(`{"id":"token-1","name":"byop","plain_token":"one-time-proxy-token"}`), nil
}

func (f *fakeRemote) DeleteReverseProxyToken(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyTokenBefore = nil
	f.proxyTokenCollection = json.RawMessage(`[]`)
	return nil, nil
}

func (f *fakeRemote) ListReverseProxyDomainsRaw(_ context.Context) (json.RawMessage, error) {
	if f.proxyDomainCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.proxyDomainCollection...), nil
}

func (f *fakeRemote) CreateReverseProxyDomain(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyDomainBody = append(json.RawMessage(nil), body...)
	if f.proxyDomainAfter == nil {
		return nil, errors.New("missing proxy domain")
	}
	f.proxyDomainCollection = json.RawMessage("[" + string(f.proxyDomainAfter) + "]")
	return append(json.RawMessage(nil), f.proxyDomainAfter...), nil
}

func (f *fakeRemote) DeleteReverseProxyDomain(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyDomainBefore = nil
	f.proxyDomainCollection = json.RawMessage(`[]`)
	return nil, nil
}

func (f *fakeRemote) ListReverseProxyClustersRaw(_ context.Context) (json.RawMessage, error) {
	if f.proxyClusterCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.proxyClusterCollection...), nil
}

func (f *fakeRemote) DeleteReverseProxyCluster(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyClusterCollection = json.RawMessage(`[]`)
	return nil, nil
}

func (f *fakeRemote) ListReverseProxyServicesRaw(_ context.Context) (json.RawMessage, error) {
	if f.proxyServiceCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.proxyServiceCollection...), nil
}

func (f *fakeRemote) GetReverseProxyServiceRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.proxyServiceBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.proxyServiceBefore...), nil
}

func (f *fakeRemote) CreateReverseProxyService(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyServiceBody = append(json.RawMessage(nil), body...)
	if f.proxyServiceAfter == nil {
		return nil, errors.New("missing proxy service")
	}
	f.proxyServiceCollection = json.RawMessage("[" + string(f.proxyServiceAfter) + "]")
	return append(json.RawMessage(nil), f.proxyServiceAfter...), nil
}

func (f *fakeRemote) UpdateReverseProxyService(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyServiceBody = append(json.RawMessage(nil), body...)
	if f.proxyServiceAfter == nil {
		return nil, errors.New("missing proxy service")
	}
	f.proxyServiceBefore = append(json.RawMessage(nil), f.proxyServiceAfter...)
	return append(json.RawMessage(nil), f.proxyServiceAfter...), nil
}

func (f *fakeRemote) DeleteReverseProxyService(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.proxyServiceBefore = nil
	f.proxyServiceCollection = json.RawMessage(`[]`)
	return nil, nil
}

func (f *fakeRemote) ListNotificationChannelsRaw(_ context.Context) (json.RawMessage, error) {
	if f.notification.collection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.notification.collection...), nil
}

func (f *fakeRemote) GetNotificationChannelRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.notification.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.notification.before...), nil
}

func (f *fakeRemote) CreateNotificationChannel(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.notification.body = append(json.RawMessage(nil), body...)
	if f.notification.after == nil {
		return nil, errors.New("missing notification channel")
	}
	f.notification.collection = json.RawMessage("[" + string(f.notification.after) + "]")
	return append(json.RawMessage(nil), f.notification.after...), nil
}

func (f *fakeRemote) UpdateNotificationChannel(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.notification.body = append(json.RawMessage(nil), body...)
	if f.notification.after == nil {
		return nil, errors.New("missing notification channel")
	}
	f.notification.before = append(json.RawMessage(nil), f.notification.after...)
	return append(json.RawMessage(nil), f.notification.after...), nil
}

func (f *fakeRemote) DeleteNotificationChannel(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.notification.before = nil
	f.notification.collection = json.RawMessage(`[]`)
	return nil, nil
}

func (f *fakeRemote) ListAzureIDPsRaw(_ context.Context) (json.RawMessage, error) {
	if f.azureIDP.collection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.azureIDP.collection...), nil
}

func (f *fakeRemote) GetAzureIDPRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.azureIDP.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.azureIDP.before...), nil
}

func (f *fakeRemote) CreateAzureIDP(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.azureIDP.body = append(json.RawMessage(nil), body...)
	if f.azureIDP.after == nil {
		return nil, errors.New("missing Azure IDP")
	}
	f.azureIDP.collection = json.RawMessage("[" + string(f.azureIDP.after) + "]")
	return append(json.RawMessage(nil), f.azureIDP.after...), nil
}

func (f *fakeRemote) UpdateAzureIDP(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.azureIDP.body = append(json.RawMessage(nil), body...)
	if f.azureIDP.after == nil {
		return nil, errors.New("missing Azure IDP")
	}
	f.azureIDP.before = append(json.RawMessage(nil), f.azureIDP.after...)
	return append(json.RawMessage(nil), f.azureIDP.after...), nil
}

func (f *fakeRemote) DeleteAzureIDP(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.azureIDP.before = nil
	f.azureIDP.collection = json.RawMessage(`[]`)
	return json.RawMessage(`{}`), nil
}

func (f *fakeRemote) SyncAzureIDP(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return json.RawMessage(`{"result":"ok"}`), nil
}

func (f *fakeRemote) ListGoogleIDPsRaw(_ context.Context) (json.RawMessage, error) {
	if f.googleIDP.collection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.googleIDP.collection...), nil
}

func (f *fakeRemote) GetGoogleIDPRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.googleIDP.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.googleIDP.before...), nil
}

func (f *fakeRemote) CreateGoogleIDP(_ context.Context, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.googleIDP.body = append(json.RawMessage(nil), body...)
	if f.googleIDP.after == nil {
		return nil, errors.New("missing google IDP")
	}
	f.googleIDP.collection = json.RawMessage("[" + string(f.googleIDP.after) + "]")
	return append(json.RawMessage(nil), f.googleIDP.after...), nil
}

func (f *fakeRemote) UpdateGoogleIDP(_ context.Context, _ string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.googleIDP.body = append(json.RawMessage(nil), body...)
	if f.googleIDP.after == nil {
		return nil, errors.New("missing google IDP")
	}
	f.googleIDP.before = append(json.RawMessage(nil), f.googleIDP.after...)
	return append(json.RawMessage(nil), f.googleIDP.after...), nil
}

func (f *fakeRemote) DeleteGoogleIDP(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.googleIDP.before = nil
	f.googleIDP.collection = json.RawMessage(`[]`)
	return json.RawMessage(`{}`), nil
}

func (f *fakeRemote) SyncGoogleIDP(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	return json.RawMessage(`{"result":"ok"}`), nil
}

func (f *fakeRemote) edrState(provider string) *edrIntegrationState {
	if f.edr == nil {
		f.edr = make(map[string]edrIntegrationState)
	}
	state := f.edr[provider]
	f.edr[provider] = state
	return &state
}

func (f *fakeRemote) GetEDRIntegrationRaw(_ context.Context, provider string) (json.RawMessage, error) {
	state := f.edrState(provider)
	if state.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), state.before...), nil
}

func (f *fakeRemote) CreateEDRIntegration(_ context.Context, provider string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.edrState(provider)
	state.body = append(json.RawMessage(nil), body...)
	if state.after == nil {
		return nil, errors.New("missing EDR integration")
	}
	state.before = append(json.RawMessage(nil), state.after...)
	f.edr[provider] = *state
	return append(json.RawMessage(nil), state.after...), nil
}

func (f *fakeRemote) UpdateEDRIntegration(_ context.Context, provider string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.edrState(provider)
	state.body = append(json.RawMessage(nil), body...)
	if state.after == nil {
		return nil, errors.New("missing EDR integration")
	}
	state.before = append(json.RawMessage(nil), state.after...)
	f.edr[provider] = *state
	return append(json.RawMessage(nil), state.after...), nil
}

func (f *fakeRemote) DeleteEDRIntegration(_ context.Context, provider string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.edrState(provider)
	state.before = nil
	f.edr[provider] = *state
	return json.RawMessage(`{}`), nil
}

func (f *fakeRemote) scimState(provider string) *scimIntegrationState {
	if f.scim == nil {
		f.scim = make(map[string]scimIntegrationState)
	}
	state := f.scim[provider]
	f.scim[provider] = state
	return &state
}

func (f *fakeRemote) ListSCIMIntegrationsRaw(_ context.Context, provider string) (json.RawMessage, error) {
	state := f.scimState(provider)
	if state.collection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), state.collection...), nil
}

func (f *fakeRemote) GetSCIMIntegrationRaw(_ context.Context, provider, _ string) (json.RawMessage, error) {
	state := f.scimState(provider)
	if state.before == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), state.before...), nil
}

func (f *fakeRemote) CreateSCIMIntegration(_ context.Context, provider string, body json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.scimState(provider)
	if state.after == nil {
		return nil, errors.New("missing SCIM integration")
	}
	state.collection = json.RawMessage("[" + string(state.after) + "]")
	state.before = append(json.RawMessage(nil), state.after...)
	f.scim[provider] = *state
	return append(json.RawMessage(nil), state.after...), nil
}

func (f *fakeRemote) UpdateSCIMIntegration(_ context.Context, provider, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.scimState(provider)
	state.before = append(json.RawMessage(nil), state.after...)
	f.scim[provider] = *state
	return append(json.RawMessage(nil), state.after...), nil
}

func (f *fakeRemote) DeleteSCIMIntegration(_ context.Context, provider, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.scimState(provider)
	state.before = nil
	state.collection = json.RawMessage(`[]`)
	f.scim[provider] = *state
	return json.RawMessage(`{}`), nil
}

func (f *fakeRemote) RegenerateSCIMToken(_ context.Context, provider, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	state := f.scimState(provider)
	if state.token == "" {
		state.token = "nbs-one-time"
	}
	f.scim[provider] = *state
	return json.RawMessage(`{"auth_token":"nbs-one-time"}`), nil
}

func (f *fakeRemote) DeletePeer(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.peerBefore = nil
	return nil, nil
}

func (f *fakeRemote) ListEDRBypassedPeersRaw(_ context.Context) (json.RawMessage, error) {
	if f.edrBypassed == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.edrBypassed...), nil
}

func (f *fakeRemote) BypassPeerEDR(_ context.Context, peerID string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.edrBypassed = json.RawMessage(`[{"peer_id":"` + peerID + `"}]`)
	return json.RawMessage(`{"peer_id":"` + peerID + `"}`), nil
}

func (f *fakeRemote) RevokePeerEDRBypass(_ context.Context, peerID string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	var peers []map[string]string
	if err := json.Unmarshal(f.edrBypassed, &peers); err == nil {
		filtered := peers[:0]
		for _, peer := range peers {
			if peer["peer_id"] != peerID {
				filtered = append(filtered, peer)
			}
		}
		encoded, _ := json.Marshal(filtered)
		f.edrBypassed = encoded
	}
	return nil, nil
}

func (f *fakeRemote) GetNetworkRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.networkBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.networkBefore...), nil
}

func (f *fakeRemote) ListNetworksRaw(_ context.Context) (json.RawMessage, error) {
	if f.networkCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.networkCollection...), nil
}

func (f *fakeRemote) CreateNetwork(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.networkAfter == nil {
		return nil, errors.New("missing created network")
	}
	f.networkCollection = json.RawMessage("[" + string(f.networkAfter) + "]")
	return append(json.RawMessage(nil), f.networkAfter...), nil
}

func (f *fakeRemote) UpdateNetwork(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.networkBefore = append(json.RawMessage(nil), f.networkAfter...)
	return f.networkBefore, nil
}

func (f *fakeRemote) DeleteNetwork(_ context.Context, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.networkBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetNetworkResourceRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.resourceBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.resourceBefore...), nil
}

func (f *fakeRemote) ListNetworkResourcesRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.resourceCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.resourceCollection...), nil
}

func (f *fakeRemote) CreateNetworkResource(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.resourceAfter == nil {
		return nil, errors.New("missing created resource")
	}
	f.resourceCollection = json.RawMessage("[" + string(f.resourceAfter) + "]")
	return append(json.RawMessage(nil), f.resourceAfter...), nil
}

func (f *fakeRemote) UpdateNetworkResource(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.resourceBefore = append(json.RawMessage(nil), f.resourceAfter...)
	return f.resourceBefore, nil
}

func (f *fakeRemote) DeleteNetworkResource(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.resourceBefore = nil
	return nil, nil
}

func (f *fakeRemote) GetNetworkRouterRaw(_ context.Context, _, _ string) (json.RawMessage, error) {
	if f.routerBefore == nil {
		return nil, &transport.RequestError{Dispatched: true, StatusCode: 404, Description: "not found"}
	}
	return append(json.RawMessage(nil), f.routerBefore...), nil
}

func (f *fakeRemote) ListNetworkRoutersRaw(_ context.Context, _ string) (json.RawMessage, error) {
	if f.routerCollection == nil {
		return json.RawMessage(`[]`), nil
	}
	return append(json.RawMessage(nil), f.routerCollection...), nil
}

func (f *fakeRemote) CreateNetworkRouter(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	if f.routerAfter == nil {
		return nil, errors.New("missing created router")
	}
	f.routerCollection = json.RawMessage("[" + string(f.routerAfter) + "]")
	return append(json.RawMessage(nil), f.routerAfter...), nil
}

func (f *fakeRemote) UpdateNetworkRouter(_ context.Context, _, _ string, _ json.RawMessage) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routerBefore = append(json.RawMessage(nil), f.routerAfter...)
	return f.routerBefore, nil
}

func (f *fakeRemote) DeleteNetworkRouter(_ context.Context, _, _ string) (json.RawMessage, error) {
	f.updates++
	if f.updateErr != nil {
		return nil, f.updateErr
	}
	f.routerBefore = nil
	return nil, nil
}

func stageForTest(t *testing.T, store *ledger.Store, before, after string) ledger.Stage {
	t.Helper()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
	})
	if err != nil {
		t.Fatal(err)
	}
	return stage
}

func TestApplyJournalsAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old","peers_count":0,"resources_count":0}`
	after := `{"id":"g1","name":"new","peers_count":0,"resources_count":0}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || result.AttemptID == "" {
		t.Fatalf("unexpected result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil || receipt.State != string(mutation.ConfirmedSuccess) {
		t.Fatalf("receipt=%+v err=%v", receipt, err)
	}
}

func TestApplyDispatchesGroupCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"engineering"}`
	created := `{"id":"g1","name":"engineering","peers_count":0,"resources_count":0}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.create",
		Request:        json.RawMessage(`{"name":"engineering"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"group_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a group can create a new policy or resource membership principal; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"group_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.group_create", Severity: "blocking", Message: "creating the group may alter policy or resource membership and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", groupCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected group create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","domain":"office.internal","enabled":true}`
	created := `{"id":"zone-1","name":"office","domain":"office.internal","enabled":true,"distribution_groups":[],"records":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.create",
		Request:        json.RawMessage(intended),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_zone_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_create", Severity: "blocking", Message: "creating the DNS zone may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsZoneCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"zone-1","domain":"office.internal","distribution_groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.delete",
		Request:        json.RawMessage(`{"id":"zone-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_zone_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_delete", Severity: "blocking", Message: "deleting the DNS zone may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSZoneUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"zone-1","domain":"office.internal","enabled":true}`
	after := `{"id":"zone-1","domain":"corp.internal","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.zones.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_zone_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a DNS zone can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_zone_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_zone_change", Severity: "blocking", Message: "the proposed DNS zone change may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns zone update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	request := `{"zone_id":"zone-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	intended := `{"name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	created := `{"id":"record-2","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.create",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_record_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_create", Severity: "blocking", Message: "creating the DNS record may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsZoneCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	after := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}`
	request := `{"zone_id":"zone-1","id":"record-1","name":"db","type":"A","content":"10.0.0.6","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.update",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_record_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_change", Severity: "blocking", Message: "the proposed DNS record change may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSRecordDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"record-1","name":"db","type":"A","content":"10.0.0.5","ttl":60}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.records.delete",
		Request:        json.RawMessage(`{"zone_id":"zone-1","id":"record-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_record_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a DNS record can change name-resolution behavior for distributed peers; affected peers and records require live analysis"],"completeness":{"state":"unknown","reason":"dns_record_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_record_delete", Severity: "blocking", Message: "deleting the DNS record may alter name resolution and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns record delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	created := `{"id":"ns-1","name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.create",
		Request:        json.RawMessage(intended),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_create_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_create", Severity: "blocking", Message: "creating the nameserver group may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverCollection: []byte(before), after: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ns-1","name":"office","description":"old","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	after := `{"id":"ns-1","name":"office","description":"new","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	request := `{"id":"ns-1","description":"new"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.update",
		Request:        json.RawMessage(request),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_change", Severity: "blocking", Message: "the proposed nameserver group change may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverBefore: []byte(before), nameserverAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNameserverGroupDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ns-1","name":"office","description":"office resolvers","domains":["office.internal"],"enabled":true,"groups":["g1"],"nameservers":[{"ip":"10.0.0.53","ns_type":"udp","port":53}],"primary":false,"search_domains_enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.nameservers.delete",
		Request:        json.RawMessage(`{"id":"ns-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"dns_nameserver_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a nameserver group can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_nameserver_delete_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_nameserver_delete", Severity: "blocking", Message: "deleting the nameserver group may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", nameserverBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected nameserver delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesDNSSettingsUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"disabled_management_groups":["g1"]}`
	after := `{"disabled_management_groups":["g1","g2"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "dns.settings.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"dns_settings_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating DNS settings can change resolver behavior for distributed peers; affected peers and domains require live analysis"],"completeness":{"state":"unknown","reason":"dns_settings_update_requires_dns_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.dns_settings_change", Severity: "blocking", Message: "the proposed DNS settings change may alter resolver behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", dnsSettingsBefore: []byte(before), dnsSettingsAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected dns settings update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAccountUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"account-1","settings":{"peer_login_expiration_enabled":true}}`
	after := `{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "accounts.update",
		Request:        json.RawMessage(`{"id":"account-1","settings":{"peer_login_expiration_enabled":false}}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"account_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating account settings can change authentication, posture, DNS, or network behavior; affected peers and resources require live analysis"],"completeness":{"state":"unknown","reason":"account_update_requires_management_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.account_change", Severity: "blocking", Message: "the proposed account change may alter management-plane behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", accountBefore: []byte(before), accountAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected account update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAccountDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"account-1","domain":"example.test"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "accounts.delete",
		Request:        json.RawMessage(`{"id":"account-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"account_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an account can remove management access and account resources; affected peers and resources require live analysis"],"completeness":{"state":"unknown","reason":"account_delete_requires_management_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.account_delete", Severity: "blocking", Message: "deleting the account may remove management access and resources and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", accountBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected account delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPostureCheckCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"pc-2","name":"managed","checks":{}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "posture_checks.create",
		Request:        json.RawMessage(`{"name":"managed","checks":{}}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"posture_check_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a posture check can change policy admission for peers; affected peers and policies require live analysis"],"completeness":{"state":"unknown","reason":"posture_check_create_requires_policy_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.posture_check_create", Severity: "blocking", Message: "creating the posture check may alter policy admission and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", postureCollection: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected posture check create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPostureCheckUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"pc-1","name":"managed","checks":{}}`
	after := `{"id":"pc-1","name":"managed-v2","checks":{}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "posture_checks.update",
		Request:        json.RawMessage(`{"id":"pc-1","name":"managed-v2","checks":{}}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"posture_check_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a posture check can change policy admission for peers; affected peers and policies require live analysis"],"completeness":{"state":"unknown","reason":"posture_check_update_requires_policy_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.posture_check_change", Severity: "blocking", Message: "the proposed posture check change may alter policy admission and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", postureBefore: []byte(before), postureAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected posture check update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPostureCheckDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"pc-1","name":"managed","checks":{}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "posture_checks.delete",
		Request:        json.RawMessage(`{"id":"pc-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"posture_check_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a posture check can change policy admission for peers; affected peers and policies require live analysis"],"completeness":{"state":"unknown","reason":"posture_check_delete_requires_policy_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.posture_check_delete", Severity: "blocking", Message: "deleting the posture check may alter policy admission and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", postureBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected posture check delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesIngressPeerCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"ing-2","peer_id":"peer-1","region":"eu","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "ingress.peers.create",
		Request:        json.RawMessage(`{"peer_id":"peer-1","region":"eu","enabled":true}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"ingress_peer_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an ingress peer can add an external reachability path; affected peers and allocations require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_peer_create_requires_topology_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.ingress_peer_create", Severity: "blocking", Message: "creating the ingress peer may add external reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", ingressCollection: []byte(before), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected ingress peer create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesIngressPeerUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ing-1","peer_id":"peer-1","enabled":true}`
	after := `{"id":"ing-1","peer_id":"peer-1","enabled":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "ingress.peers.update",
		Request:        json.RawMessage(`{"id":"ing-1","enabled":false}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"ingress_peer_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating an ingress peer can change an external reachability path; affected peers and allocations require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_peer_update_requires_topology_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.ingress_peer_change", Severity: "blocking", Message: "the proposed ingress peer change may alter external reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", ingressBefore: []byte(before), ingressAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected ingress peer update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesIngressPeerDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"ing-1","peer_id":"peer-1","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "ingress.peers.delete",
		Request:        json.RawMessage(`{"id":"ing-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"ingress_peer_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an ingress peer can remove an external reachability path; affected peers and allocations require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_peer_delete_requires_topology_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.ingress_peer_delete", Severity: "blocking", Message: "deleting the ingress peer may remove external reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", ingressBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected ingress peer delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesIngressPortAllocationMutations(t *testing.T) {
	cases := []struct {
		name, operation, request, before, after, impact, finding, message string
	}{
		{
			name: "create", operation: "peers.ingress.ports.create",
			request: `{"peer_id":"peer-1","name":"web","enabled":true}`,
			before:  `[]`, after: `{"id":"alloc-1","name":"web","enabled":true}`,
			impact:  `{"classification":"ingress_port_allocation_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an ingress port allocation can expose peer services through an external ingress path; affected peers and mappings require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_port_allocation_create_requires_topology_analysis"}}`,
			finding: "impact.ingress_port_allocation_create", message: "creating the ingress port allocation may expose peer services externally and requires exact acknowledgement",
		},
		{
			name: "update", operation: "peers.ingress.ports.update",
			request: `{"peer_id":"peer-1","id":"alloc-1","enabled":false}`,
			before:  `{"id":"alloc-1","name":"web","enabled":true}`, after: `{"id":"alloc-1","name":"web","enabled":false}`,
			impact:  `{"classification":"ingress_port_allocation_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating an ingress port allocation can change externally exposed peer services and port mappings; affected peers require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_port_allocation_update_requires_topology_analysis"}}`,
			finding: "impact.ingress_port_allocation_change", message: "the proposed ingress port allocation change may alter external peer exposure and requires exact acknowledgement",
		},
		{
			name: "delete", operation: "peers.ingress.ports.delete",
			request: `{"peer_id":"peer-1","id":"alloc-1"}`,
			before:  `{"id":"alloc-1","name":"web","enabled":true}`, after: `{}`,
			impact:  `{"classification":"ingress_port_allocation_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an ingress port allocation can remove externally exposed peer services and port mappings; affected peers require live topology analysis"],"completeness":{"state":"unknown","reason":"ingress_port_allocation_delete_requires_topology_analysis"}}`,
			finding: "impact.ingress_port_allocation_delete", message: "deleting the ingress port allocation may remove external peer exposure and requires exact acknowledgement",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store, err := ledger.Open(t.TempDir() + "/ledger.db")
			if err != nil {
				t.Fatal(err)
			}
			defer store.Close()
			stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: tc.operation, Request: json.RawMessage(tc.request), Before: json.RawMessage(tc.before), IntendedAfter: json.RawMessage(tc.after), Impact: json.RawMessage(tc.impact), Findings: []ledger.Finding{{Code: tc.finding, Severity: "blocking", Message: tc.message}}})
			if err != nil {
				t.Fatal(err)
			}
			remote := &fakeRemote{identity: "https://nb.test", account: "account-1", ingressPortCollection: []byte(tc.before), ingressPortBefore: []byte(tc.before), ingressPortAfter: []byte(tc.after)}
			result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
			if err != nil {
				t.Fatal(err)
			}
			if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
				t.Fatalf("unexpected ingress port result: %+v updates=%d", result, remote.updates)
			}
			if tc.name != "delete" {
				var body map[string]any
				if err := json.Unmarshal(remote.ingressPortBody, &body); err != nil {
					t.Fatal(err)
				}
				if _, ok := body["peer_id"]; ok {
					t.Fatalf("peer_id leaked into dispatched ingress port body: %s", remote.ingressPortBody)
				}
			}
		})
	}
}

func TestApplyDispatchesEDRBypassMutations(t *testing.T) {
	cases := []struct {
		name, operation, before, after, impact, finding, message string
		bypassed                                                 json.RawMessage
	}{
		{
			name: "create", operation: "peers.edr.bypass.create", before: `[]`, after: `{"peer_id":"peer-1"}`,
			impact:  `{"classification":"edr_bypass_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["bypassing EDR compliance immediately grants the peer network access outside the normal compliance control","the bypassed-peer collection is the authoritative postcondition; the API response is not used as sole effect proof"],"completeness":{"state":"unknown","reason":"edr_bypass_create_requires_compliance_analysis"}}`,
			finding: "impact.edr_bypass_create", message: "bypassing EDR compliance immediately grants peer network access and requires exact acknowledgement",
		},
		{
			name: "delete", operation: "peers.edr.bypass.delete", before: `[{"peer_id":"peer-1"}]`, after: `{}`,
			impact:  `{"classification":"edr_bypass_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["revoking an EDR bypass restores the peer's normal compliance gate and may remove current network access","the bypassed-peer collection is the authoritative absence postcondition"],"completeness":{"state":"unknown","reason":"edr_bypass_delete_requires_compliance_analysis"}}`,
			finding: "impact.edr_bypass_delete", message: "revoking the EDR bypass restores compliance gating and may remove peer network access; exact acknowledgement is required",
			bypassed: json.RawMessage(`[{"peer_id":"peer-1"}]`),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			store, err := ledger.Open(t.TempDir() + "/ledger.db")
			if err != nil {
				t.Fatal(err)
			}
			defer store.Close()
			stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: tc.operation, Request: json.RawMessage(`{"id":"peer-1"}`), Before: json.RawMessage(tc.before), IntendedAfter: json.RawMessage(tc.after), Impact: json.RawMessage(tc.impact), Findings: []ledger.Finding{{Code: tc.finding, Severity: "blocking", Message: tc.message}}})
			if err != nil {
				t.Fatal(err)
			}
			remote := &fakeRemote{identity: "https://nb.test", account: "account-1", edrBypassed: tc.bypassed}
			result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
			if err != nil {
				t.Fatal(err)
			}
			if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
				t.Fatalf("unexpected EDR bypass result: %+v updates=%d", result, remote.updates)
			}
		})
	}
}

func TestApplyDispatchesAgentNetworkSettingsUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"enabled":true}`
	after := `{"enabled":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.settings.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"agent_network_settings_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating agent-network settings can change Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_settings_update_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_settings_change", Severity: "blocking", Message: "the proposed agent-network settings change may alter Cloud agent behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentSettingsBefore: []byte(before), agentSettingsAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent-network settings result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkSettingsCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{}`
	after := `{"enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.settings.create",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"agent_network_settings_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating agent-network settings can enable Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_settings_create_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_settings_create", Severity: "blocking", Message: "creating agent-network settings may enable Cloud agent behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentSettingsBefore: []byte(before), agentSettingsAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent-network settings create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkSettingsDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.settings.delete",
		Request:        json.RawMessage(`{}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"agent_network_settings_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting agent-network settings can disable Cloud agent routing and provider behavior; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_settings_delete_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_settings_delete", Severity: "blocking", Message: "deleting agent-network settings may disable Cloud agent behavior and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentSettingsBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent-network settings delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkBudgetRuleCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"rule-1","name":"monthly","enabled":true,"limits":{"budget_limit":{"enabled":true,"group_cap_usd":100,"user_cap_usd":0,"window_seconds":3600}}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.budget_rules.create",
		Request:        json.RawMessage(`{"name":"monthly","enabled":true}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"agent_network_budget_rule_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an agent-network budget rule can change spend and token admission for targeted callers; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_budget_rule_create_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_budget_rule_create", Severity: "blocking", Message: "creating the agent-network budget rule may change spend and token admission and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", budgetCollection: []byte(before), budgetAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected budget rule create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkBudgetRuleUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rule-1","name":"monthly","enabled":true}`
	after := `{"id":"rule-1","name":"monthly","enabled":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.budget_rules.update",
		Request:        json.RawMessage(after),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"agent_network_budget_rule_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating an agent-network budget rule can change spend and token admission for targeted callers; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_budget_rule_update_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_budget_rule_change", Severity: "blocking", Message: "the proposed agent-network budget rule change may change spend and token admission and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", budgetBefore: []byte(before), budgetAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected budget rule update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkBudgetRuleDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rule-1","name":"monthly","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "agent_network.budget_rules.delete",
		Request:        json.RawMessage(`{"id":"rule-1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"agent_network_budget_rule_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an agent-network budget rule can remove spend and token admission limits for targeted callers; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_budget_rule_delete_requires_capability_analysis"}}`),
		Findings:       []ledger.Finding{{Code: "impact.agent_network_budget_rule_delete", Severity: "blocking", Message: "deleting the agent-network budget rule may remove spend and token limits and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", budgetBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected budget rule delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkGuardrailCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"guard-1","name":"strict","checks":{"model_allowlist":{"enabled":true}}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.guardrails.create", Request: json.RawMessage(`{"name":"strict"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"agent_network_guardrail_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an agent-network guardrail can alter model and prompt admission for attached policies; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_guardrail_create_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_guardrail_create", Severity: "blocking", Message: "creating the agent-network guardrail may alter policy admission and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", guardrailCollection: []byte(before), guardrailAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected guardrail create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkGuardrailUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"guard-1","name":"strict","checks":{"model_allowlist":{"enabled":true}}}`
	after := `{"id":"guard-1","name":"strict","checks":{"model_allowlist":{"enabled":false}}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.guardrails.update", Request: json.RawMessage(after), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"agent_network_guardrail_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating an agent-network guardrail can alter model and prompt admission for attached policies; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_guardrail_update_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_guardrail_change", Severity: "blocking", Message: "the proposed agent-network guardrail change may alter policy admission and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", guardrailBefore: []byte(before), guardrailAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected guardrail update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkGuardrailDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"guard-1","name":"strict","checks":{"model_allowlist":{"enabled":true}}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.guardrails.delete", Request: json.RawMessage(`{"id":"guard-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"agent_network_guardrail_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an agent-network guardrail can remove model and prompt admission controls from attached policies; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_guardrail_delete_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_guardrail_delete", Severity: "blocking", Message: "deleting the agent-network guardrail may remove policy admission controls and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", guardrailBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected guardrail delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkPolicyCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"policy-1","name":"engineering","enabled":true,"source_groups":["group-1"],"destination_provider_ids":["provider-1"],"guardrail_ids":[],"limits":{"token_limit":{"enabled":false},"budget_limit":{"enabled":false}}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.policies.create", Request: json.RawMessage(`{"name":"engineering"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"agent_network_policy_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an agent-network policy can add provider reachability and admission limits for source groups; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_policy_create_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_policy_create", Severity: "blocking", Message: "creating the agent-network policy may add provider reachability and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentPolicyCollection: []byte(before), agentPolicyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent policy create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyResolvesAgentNetworkProviderSecretOnlyForDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"provider-1","name":"OpenAI"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.providers.create", Request: json.RawMessage(`{"name":"OpenAI","api_key_ref":"env:PROVIDER_KEY"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"agent_network_provider_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an agent-network provider can add upstream reachability and secret-bearing routing configuration; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_provider_create_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_provider_create", Severity: "blocking", Message: "creating the agent-network provider may add upstream reachability and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", providerCollection: []byte(before), providerAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "env:PROVIDER_KEY" {
			t.Fatalf("unexpected secret ref %q", ref)
		}
		return "ephemeral", nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || !remote.providerSecretSeen || remote.updates != 1 {
		t.Fatalf("unexpected provider create result: %+v updates=%d secret_seen=%v", result, remote.updates, remote.providerSecretSeen)
	}
	if strings.Contains(string(stage.Request), "ephemeral") || strings.Contains(string(stage.Request), "api_key\":\"") {
		t.Fatal("provider secret was persisted in staged request")
	}
}

func TestApplyDispatchesUserCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"user-1","email":"a@example.com","role":"user","auto_groups":[],"is_blocked":false,"pending_approval":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.create", Request: json.RawMessage(`{"email":"a@example.com","role":"user","auto_groups":[],"is_service_user":false}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"user_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a user can grant account access and assign automatic peer groups; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"user_create_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.user_create", Severity: "blocking", Message: "creating the user may grant account access and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", userCollection: []byte(before), userAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected user create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesUserRejectAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"user-1","email":"a@example.com","pending_approval":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.reject", Request: json.RawMessage(`{"id":"user-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"user_reject","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["rejecting a pending user removes a pending account-access edge; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"user_reject_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.user_reject", Severity: "blocking", Message: "rejecting the user removes a pending account-access edge and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", userBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected user reject result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyResolvesUserPasswordRefsWithoutPersistingSecrets(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"user-1","status":"active"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.password.update", Request: json.RawMessage(`{"id":"user-1","old_password_ref":"env:OLD_PASSWORD","new_password_ref":"env:NEW_PASSWORD"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(before), Impact: json.RawMessage(`{"classification":"user_password_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["changing a user's password changes an authentication credential; the password values are resolved from external references only at dispatch and are never persisted","the management API exposes no password read-back representation; success is confirmed by the endpoint response after exact user preimage validation"],"completeness":{"state":"unknown","reason":"user_password_change_has_no_readable_credential_state"}}`), Findings: []ledger.Finding{{Code: "impact.user_password_change", Severity: "blocking", Message: "changing the user password changes an authentication credential and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", userBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		switch ref {
		case "env:OLD_PASSWORD":
			return "OldPass123!", nil
		case "env:NEW_PASSWORD":
			return "NewPass123!", nil
		default:
			t.Fatalf("unexpected secret ref %q", ref)
			return "", errors.New("unexpected ref")
		}
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || !remote.passwordSeen || remote.updates != 1 {
		t.Fatalf("unexpected password result: %+v updates=%d seen=%v", result, remote.updates, remote.passwordSeen)
	}
	var body map[string]any
	if err := json.Unmarshal(remote.passwordBody, &body); err != nil || body["old_password"] != "OldPass123!" || body["new_password"] != "NewPass123!" {
		t.Fatalf("unexpected dispatched password body: %s", remote.passwordBody)
	}
	if strings.Contains(string(stage.Request), "OldPass123!") || strings.Contains(string(stage.Request), "NewPass123!") {
		t.Fatal("password leaked into staged request")
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "OldPass123!") || strings.Contains(string(receipt.Result), "NewPass123!") {
		t.Fatal("password leaked into persisted receipt")
	}
}

func TestApplyDispatchesUserInviteResendAndConfirmsMetadata(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"user-1","status":"pending"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.invite.resend", Request: json.RawMessage(`{"id":"user-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(before), Impact: json.RawMessage(`{"classification":"user_invite_resend","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["resending a user invitation creates an external enrollment communication side effect; the endpoint returns no durable token or response body","success is confirmed by the endpoint response plus unchanged user metadata after exact preimage validation"],"completeness":{"state":"unknown","reason":"user_invite_resend_has_external_delivery_side_effect"}}`), Findings: []ledger.Finding{{Code: "impact.user_invite_resend", Severity: "blocking", Message: "resending the user invitation triggers external enrollment delivery and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", userBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || !remote.inviteResendSeen || remote.updates != 1 {
		t.Fatalf("unexpected invite resend result: %+v updates=%d seen=%v", result, remote.updates, remote.inviteResendSeen)
	}
}

func TestApplyAcceptsInviteWithExternalRefsWithoutPersistingSecrets(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"email":"a@example.com","name":"New","valid":true}`
	after := `{"success":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.invites.accept", Request: json.RawMessage(`{"invite_token_ref":"env:INVITE_TOKEN","password_ref":"env:INVITE_PASSWORD"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"invite_accept","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["accepting an invite creates an account-access edge from a public token; the invite token and password are resolved from external references only at dispatch and are never persisted","the unauthenticated endpoint's success response is the only available effect proof because the accepted invite is no longer readable as pending metadata"],"completeness":{"state":"unknown","reason":"invite_accept_has_no_pending_invite_readback"}}`), Findings: []ledger.Finding{{Code: "impact.invite_accept", Severity: "blocking", Message: "accepting the invite creates account access from a public token and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", publicInviteBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		switch ref {
		case "env:INVITE_TOKEN":
			return "nbi-token-secret", nil
		case "env:INVITE_PASSWORD":
			return "NewPass123!", nil
		default:
			t.Fatalf("unexpected secret ref %q", ref)
			return "", errors.New("unexpected ref")
		}
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || !remote.publicInviteAccepted || remote.publicInviteToken != "nbi-token-secret" || remote.updates != 1 {
		t.Fatalf("unexpected invite acceptance result: %+v updates=%d token=%q accepted=%v", result, remote.updates, remote.publicInviteToken, remote.publicInviteAccepted)
	}
	var body map[string]any
	if err := json.Unmarshal(remote.publicInviteBody, &body); err != nil || body["password"] != "NewPass123!" {
		t.Fatalf("unexpected invite acceptance body: %s", remote.publicInviteBody)
	}
	if strings.Contains(string(stage.Request), "nbi-token-secret") || strings.Contains(string(stage.Request), "NewPass123!") {
		t.Fatal("invite acceptance secret leaked into staged request")
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "nbi-token-secret") || strings.Contains(string(receipt.Result), "NewPass123!") {
		t.Fatal("invite acceptance secret leaked into persisted receipt")
	}
}

func TestApplyDispatchesUserTokenDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"token-1","name":"agent","created_by":"user-1"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.tokens.delete", Request: json.RawMessage(`{"user_id":"user-1","token_id":"token-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"user_token_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["deleting a personal access token revokes the credential represented by the exact token preimage"],"completeness":{"state":"complete","reason":null}}`), Findings: []ledger.Finding{{Code: "impact.user_token_delete", Severity: "blocking", Message: "deleting the personal access token revokes a credential and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", tokenBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected token delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyReturnsPersonalAccessTokenOnceWithoutPersistingIt(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"token-1","name":"agent","created_by":"user-1"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.tokens.create", Request: json.RawMessage(`{"user_id":"user-1","name":"agent","expires_in":30}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"user_token_create","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a personal access token changes credential inventory without changing data-plane reachability; the one-time token value is returned only in the successful apply result and is never persisted"],"completeness":{"state":"complete","reason":null}}`), Findings: []ledger.Finding{{Code: "impact.user_token_create", Severity: "blocking", Message: "creating the personal access token returns a one-time secret and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", tokenCollection: []byte(before), tokenAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "one-time-secret" || remote.updates != 1 {
		t.Fatalf("unexpected token create result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "one-time-secret") || strings.Contains(string(stage.Request), "one-time-secret") {
		t.Fatal("one-time token leaked into persisted state")
	}
}

func TestApplyDispatchesSetupKeyDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"key-1","name":"bootstrap","valid":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "setup_keys.delete", Request: json.RawMessage(`{"id":"key-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"setup_key_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a setup key can prevent new peer enrollment through that credential; already-enrolled peers require separate live analysis"],"completeness":{"state":"unknown","reason":"setup_key_delete_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.setup_key_delete", Severity: "blocking", Message: "deleting the setup key may stop new peer enrollment and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", setupKeyBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected setup key delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyReturnsSetupKeyOnceWithoutPersistingIt(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"key-1","name":"bootstrap","valid":true,"auto_groups":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "setup_keys.create", Request: json.RawMessage(`{"name":"bootstrap","type":"reusable","expires_in":86400,"auto_groups":[],"usage_limit":0}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"setup_key_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a setup key can expand peer enrollment authority; the one-time key value is returned only in the successful apply result and is never persisted"],"completeness":{"state":"unknown","reason":"setup_key_create_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.setup_key_create", Severity: "blocking", Message: "creating the setup key expands peer enrollment authority and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", setupKeyCollection: []byte(before), setupKeyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "one-time-key" || remote.updates != 1 {
		t.Fatalf("unexpected setup key create result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "one-time-key") || strings.Contains(string(stage.Request), "one-time-key") {
		t.Fatal("setup key leaked into persisted state")
	}
}

func TestApplyDispatchesSetupKeyUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"key-1","name":"bootstrap","revoked":false,"auto_groups":[]}`
	after := `{"id":"key-1","name":"bootstrap","revoked":true,"auto_groups":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "setup_keys.update", Request: json.RawMessage(`{"id":"key-1","revoked":true,"auto_groups":[]}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"setup_key_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating a setup key can change revocation or auto-group enrollment authority; already-enrolled peers require separate live analysis"],"completeness":{"state":"unknown","reason":"setup_key_update_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.setup_key_change", Severity: "blocking", Message: "changing the setup key may alter peer enrollment authority and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", setupKeyBefore: []byte(before), setupKeyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(remote.setupKeyBody, &body); err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || body["revoked"] != true || body["auto_groups"] == nil {
		t.Fatalf("unexpected setup key update result: %+v updates=%d body=%s", result, remote.updates, remote.setupKeyBody)
	}
}

func TestApplyDispatchesTemporaryAccessAndConfirmsResponse(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"peer-1","name":"target","connected":true}`
	after := `{"id":"temp-1","name":"temp-host","rules":["tcp/80"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "peers.temporary_access.create", Request: json.RawMessage(`{"id":"peer-1","name":"temp-host","wg_pub_key":"pub","rules":["tcp/80"]}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"temporary_access_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a temporary access peer grants a short-lived scoped path to the target peer; automatic cleanup is controlled by the remote peer lifecycle and is not readable as durable state"],"completeness":{"state":"unknown","reason":"temporary_access_peer_lifetime_is_external"}}`), Findings: []ledger.Finding{{Code: "impact.temporary_access_create", Severity: "blocking", Message: "creating a temporary access peer grants scoped network access and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before), temporaryAccessResult: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(remote.temporaryAccessBody, &body); err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || body["id"] != nil || body["name"] != "temp-host" {
		t.Fatalf("unexpected temporary access result: %+v updates=%d body=%s", result, remote.updates, remote.temporaryAccessBody)
	}
}

func TestApplyDispatchesPeerJobAndConfirmsResponse(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"peer-1","name":"target","connected":true}`
	after := `{"id":"job-1","status":"pending","workload":{"type":"bundle"}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "peers.jobs.create", Request: json.RawMessage(`{"peer_id":"peer-1","workload":{"type":"bundle","parameters":{"anonymize":true}}}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"peer_job_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a remote job executes a diagnostic workload on the selected peer and may collect sensitive logs; the peer must be online and the response is the authoritative job-creation proof"],"completeness":{"state":"unknown","reason":"remote_job_execution_and_collection"}}`), Findings: []ledger.Finding{{Code: "impact.peer_job_create", Severity: "blocking", Message: "creating the remote job executes a workload on the peer and may collect sensitive diagnostics; exact acknowledgement is required"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before), temporaryAccessResult: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(remote.temporaryAccessBody, &body); err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || body["peer_id"] != nil || body["workload"] == nil {
		t.Fatalf("unexpected peer job result: %+v updates=%d body=%s", result, remote.updates, remote.temporaryAccessBody)
	}
}

func TestApplyEventStreamingResolvesConfigOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"stream-1","platform":"s3","enabled":true,"config":{"secret_access_key":"****"}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "event_streaming.create", Request: json.RawMessage(`{"platform":"s3","enabled":true,"config_ref":"pa:stream-config"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"event_streaming_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating an event-streaming integration exports future account activity to an external platform; the resolved configuration is dispatched in memory and is never persisted"],"completeness":{"state":"unknown","reason":"event_streaming_external_delivery"}}`), Findings: []ledger.Finding{{Code: "impact.event_streaming_create", Severity: "blocking", Message: "creating the event-streaming integration exports account activity and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", eventStreamingList: []byte(before), eventStreamingAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:stream-config" {
			t.Fatalf("unexpected secret ref: %s", ref)
		}
		return `{"bucket":"logs","secret_access_key":"real-secret"}`, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "real-secret") || !strings.Contains(string(remote.eventStreamingBody), "real-secret") {
		t.Fatalf("unexpected event-streaming result: %+v body=%s", result, remote.eventStreamingBody)
	}
}

func TestApplyIdentityProviderResolvesClientSecretOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"idp-1","name":"zitadel","type":"oidc","issuer":"https://idp.example","client_id":"client-1"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "identity_providers.create", Request: json.RawMessage(`{"type":"oidc","name":"zitadel","issuer":"https://idp.example","client_id":"client-1","client_secret_ref":"pa:idp-secret"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"identity_provider_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating an identity provider changes account authentication ingress; client secrets are resolved in memory and never persisted"],"completeness":{"state":"unknown","reason":"identity_provider_authentication_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.identity_provider_create", Severity: "blocking", Message: "creating the identity provider changes authentication ingress and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", identityProviderList: []byte(before), identityProviderAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:idp-secret" {
			t.Fatalf("unexpected secret ref: %s", ref)
		}
		return "real-client-secret", nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "real-client-secret") || !strings.Contains(string(remote.identityProviderBody), "real-client-secret") {
		t.Fatalf("unexpected identity provider result: %+v body=%s", result, remote.identityProviderBody)
	}
}

func TestApplyReturnsReverseProxyTokenOnceWithoutPersistingIt(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"token-1","name":"byop","expires_at":null}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "reverse_proxy_tokens.create", Request: json.RawMessage(`{"name":"byop","expires_in":0}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"reverse_proxy_token_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a reverse proxy token grants an external proxy credential; the clear token is returned once after dispatch and is never persisted"],"completeness":{"state":"unknown","reason":"reverse_proxy_token_external_credential"}}`), Findings: []ledger.Finding{{Code: "impact.reverse_proxy_token_create", Severity: "blocking", Message: "creating the reverse proxy token grants an external proxy credential and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", proxyTokenCollection: []byte(before), proxyTokenAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "one-time-proxy-token" || remote.updates != 1 || strings.Contains(string(stage.Request), "one-time-proxy-token") {
		t.Fatalf("unexpected reverse proxy token result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyConfirmsReverseProxyDomainDeleteFromCollection(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	object := `{"id":"domain-1","domain":"app.example.com"}`
	before := "[" + object + "]"
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "reverse_proxy_domains.delete", Request: json.RawMessage(`{"id":"domain-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"reverse_proxy_domain_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["deleting a custom reverse proxy domain removes a public DNS and ingress binding"],"completeness":{"state":"unknown","reason":"reverse_proxy_domain_public_exposure"}}`), Findings: []ledger.Finding{{Code: "impact.reverse_proxy_domain_delete", Severity: "blocking", Message: "deleting the reverse proxy domain removes public DNS and ingress exposure and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", proxyDomainCollection: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected reverse proxy domain result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyReverseProxyServiceResolvesAuthOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"service-1","name":"app","domain":"app.example.com","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "reverse_proxy_services.create", Request: json.RawMessage(`{"name":"app","domain":"app.example.com","enabled":true,"auth_ref":"pa:proxy-auth"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"reverse_proxy_service_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a reverse proxy service publishes a public ingress route to internal targets; authentication and target configuration are dispatched only after explicit acknowledgement"],"completeness":{"state":"unknown","reason":"reverse_proxy_service_public_exposure"}}`), Findings: []ledger.Finding{{Code: "impact.reverse_proxy_service_create", Severity: "blocking", Message: "creating the reverse proxy service publishes public ingress to internal targets and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", proxyServiceCollection: []byte(before), proxyServiceAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:proxy-auth" {
			t.Fatalf("unexpected auth ref: %s", ref)
		}
		return `{"password_auth":{"username":"viewer","password":"real-secret"}}`, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "real-secret") || !strings.Contains(string(remote.proxyServiceBody), "real-secret") {
		t.Fatalf("unexpected reverse proxy service result: %+v body=%s", result, remote.proxyServiceBody)
	}
}

func TestApplyNotificationChannelResolvesTargetOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"channel-1","type":"webhook","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "notification_channels.create", Request: json.RawMessage(`{"type":"webhook","enabled":true,"target_ref":"pa:notification-target"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"notification_channel_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a notification channel changes external account event delivery; target credentials are resolved in memory and never persisted"],"completeness":{"state":"unknown","reason":"notification_channel_external_delivery"}}`), Findings: []ledger.Finding{{Code: "impact.notification_channel_create", Severity: "blocking", Message: "creating the notification channel changes external account delivery and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", notification: notificationChannelState{collection: []byte(before), after: []byte(after)}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:notification-target" {
			t.Fatalf("unexpected target ref: %s", ref)
		}
		return `{"url":"https://hooks.example","headers":{"Authorization":"real-secret"}}`, nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "real-secret") || !strings.Contains(string(remote.notification.body), "real-secret") {
		t.Fatalf("unexpected notification channel result: %+v body=%s", result, remote.notification.body)
	}
}

func TestApplyAzureIDPResolvesClientSecretOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":1,"enabled":true,"client_id":"client","tenant_id":"tenant","host":"microsoft.com"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "azure_idp.create", Request: json.RawMessage(`{"client_id":"client","tenant_id":"tenant","host":"microsoft.com","client_secret_ref":"pa:azure-secret"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"azure_idp_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating an Azure identity integration changes external authentication and directory synchronization; the client secret is resolved in memory and never persisted"],"completeness":{"state":"unknown","reason":"azure_idp_authentication_and_sync_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.azure_idp_create", Severity: "blocking", Message: "creating the Azure identity integration changes external authentication and directory synchronization and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", azureIDP: azureIDPState{collection: []byte(before), after: []byte(after)}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:azure-secret" {
			t.Fatalf("unexpected secret ref: %s", ref)
		}
		return "base64-secret", nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "base64-secret") || !strings.Contains(string(remote.azureIDP.body), "base64-secret") {
		t.Fatalf("unexpected Azure IDP result: %+v body=%s", result, remote.azureIDP.body)
	}
}

func TestApplyAzureIDPSyncRequiresExactPreimageAndEndpointProof(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":1,"enabled":true,"client_id":"client","tenant_id":"tenant","host":"microsoft.com"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "azure_idp.sync", Request: json.RawMessage(`{"id":"1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(before), Impact: json.RawMessage(`{"classification":"azure_idp_sync","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["triggering an Azure identity synchronization can create or update account users and groups from the external directory; the endpoint success response is the declared proof"],"completeness":{"state":"unknown","reason":"azure_idp_external_directory_sync"}}`), Findings: []ledger.Finding{{Code: "impact.azure_idp_sync", Severity: "blocking", Message: "triggering Azure directory synchronization may create or update account users and groups and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", azureIDP: azureIDPState{before: []byte(before)}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected Azure IDP sync result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyGoogleIDPResolvesServiceAccountKeyOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":1,"enabled":true,"customer_id":"customer"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "google_idp.create", Request: json.RawMessage(`{"customer_id":"customer","service_account_key_ref":"pa:google-key"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"google_idp_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a Google identity integration changes external authentication and directory synchronization; the service-account key is resolved in memory and never persisted"],"completeness":{"state":"unknown","reason":"google_idp_authentication_and_sync_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.google_idp_create", Severity: "blocking", Message: "creating the Google identity integration changes external authentication and directory synchronization and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", googleIDP: googleIDPState{collection: []byte(before), after: []byte(after)}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:google-key" {
			t.Fatalf("unexpected secret ref: %s", ref)
		}
		return "json-service-account-key", nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "json-service-account-key") || !strings.Contains(string(remote.googleIDP.body), "json-service-account-key") {
		t.Fatalf("unexpected Google IDP result: %+v body=%s", result, remote.googleIDP.body)
	}
}

func TestApplyGoogleIDPSyncRequiresExactPreimageAndEndpointProof(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":1,"enabled":true,"customer_id":"customer"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "google_idp.sync", Request: json.RawMessage(`{"id":"1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(before), Impact: json.RawMessage(`{"classification":"google_idp_sync","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["triggering a Google identity synchronization can create or update account users and groups from the external directory; the endpoint success response is the declared proof"],"completeness":{"state":"unknown","reason":"google_idp_external_directory_sync"}}`), Findings: []ledger.Finding{{Code: "impact.google_idp_sync", Severity: "blocking", Message: "triggering Google directory synchronization may create or update account users and groups and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", googleIDP: googleIDPState{before: []byte(before)}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected Google IDP sync result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyEDRIntuneResolvesSecretOnlyAtDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `null`
	after := `{"id":1,"enabled":true,"client_id":"client","tenant_id":"tenant","groups":[],"last_synced_interval":24}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "edr.intune.create", Request: json.RawMessage(`{"client_id":"client","tenant_id":"tenant","secret_ref":"pa:intune-secret","groups":[],"last_synced_interval":24}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"edr_intune_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating an EDR integration changes device-compliance enforcement and may deny or restore peer access; credentials are resolved in memory and never persisted"],"completeness":{"state":"unknown","reason":"edr_integration_compliance_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.edr_intune_create", Severity: "blocking", Message: "creating the EDR integration changes device-compliance enforcement and peer access and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", edr: map[string]edrIntegrationState{"intune": {after: []byte(after)}}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true, SecretResolver: func(ref string) (string, error) {
		if ref != "pa:intune-secret" {
			t.Fatalf("unexpected secret ref: %s", ref)
		}
		return "real-intune-secret", nil
	}})
	if err != nil {
		t.Fatal(err)
	}
	state := remote.edr["intune"]
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 || strings.Contains(string(stage.Request), "real-intune-secret") || !strings.Contains(string(state.body), "real-intune-secret") {
		t.Fatalf("unexpected EDR result: %+v body=%s", result, state.body)
	}
}

func TestApplyEDRDeleteConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":1,"enabled":true,"api_url":"https://edr.example","groups":[],"last_synced_interval":24,"match_attributes":{}}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "edr.fleetdm.delete", Request: json.RawMessage(`{}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`null`), Impact: json.RawMessage(`{"classification":"edr_fleetdm_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["deleting an EDR integration removes a device-compliance gate and may change peer access"],"completeness":{"state":"unknown","reason":"edr_integration_compliance_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.edr_fleetdm_delete", Severity: "blocking", Message: "deleting the EDR integration removes a device-compliance gate and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", edr: map[string]edrIntegrationState{"fleetdm": {before: []byte(before)}}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected EDR delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplySCIMCreateConfirmsCollectionAndTokenIsOneTime(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":1,"enabled":true,"prefix":"acme","provider":"okta"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "scim.create", Request: json.RawMessage(`{"prefix":"acme","provider":"okta"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"scim_scim_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["creating a SCIM identity integration changes external user and group provisioning; its one-time token is returned only to the caller and never persisted"],"completeness":{"state":"unknown","reason":"scim_provisioning_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.scim_scim_create", Severity: "blocking", Message: "creating the SCIM integration changes external user and group provisioning and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", scim: map[string]scimIntegrationState{"scim": {collection: []byte(before), after: []byte(after)}}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected SCIM create result: %+v updates=%d", result, remote.updates)
	}

	tokenBefore := after
	tokenStage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "scim.token", Request: json.RawMessage(`{"id":"1"}`), Before: json.RawMessage(tokenBefore), IntendedAfter: json.RawMessage(tokenBefore), Impact: json.RawMessage(`{"classification":"scim_scim_token","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["regenerating the SCIM token revokes the prior external provisioning credential; the replacement token is returned once and never persisted"],"completeness":{"state":"unknown","reason":"scim_token_credential_boundary"}}`), Findings: []ledger.Finding{{Code: "impact.scim_scim_token", Severity: "blocking", Message: "regenerating the SCIM token revokes the prior provisioning credential and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	result, err = Apply(context.Background(), store, remote, ApplyInput{StageID: tokenStage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "nbs-one-time" || strings.Contains(string(receipt.Result), "nbs-one-time") {
		t.Fatalf("unexpected SCIM token result: %+v receipt=%s", result, receipt.Result)
	}
}

func TestApplyConfirmsReverseProxyClusterDeleteFromCollection(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[{"address":"proxy.example.com","type":"account"}]`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "reverse_proxy_clusters.delete", Request: json.RawMessage(`{"id":"proxy.example.com"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"reverse_proxy_cluster_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["deleting an account reverse proxy cluster removes its public ingress infrastructure; running proxy processes may remain connected until stopped or revoked separately"],"completeness":{"state":"unknown","reason":"reverse_proxy_cluster_public_exposure"}}`), Findings: []ledger.Finding{{Code: "impact.reverse_proxy_cluster_delete", Severity: "blocking", Message: "deleting the reverse proxy cluster removes public ingress infrastructure and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", proxyClusterCollection: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected reverse proxy cluster result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyReturnsInviteTokenOnceWithoutPersistingIt(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	after := `{"id":"invite-1","email":"a@example.com","name":"New","role":"user","auto_groups":[],"expired":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.invites.create", Request: json.RawMessage(`{"email":"a@example.com","name":"New","role":"user","auto_groups":[]}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"invite_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating an invite can expand account enrollment authority; the one-time invite token is returned only in the successful apply result and is never persisted"],"completeness":{"state":"unknown","reason":"invite_create_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.invite_create", Severity: "blocking", Message: "creating an invite expands enrollment authority and returns a one-time token; exact acknowledgement is required"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", inviteCollection: []byte(before), inviteAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "one-time-invite" || remote.updates != 1 {
		t.Fatalf("unexpected invite create result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "one-time-invite") {
		t.Fatal("invite token leaked into persisted receipt")
	}
}

func TestApplyDispatchesInviteDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"invite-1","email":"a@example.com"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.invites.delete", Request: json.RawMessage(`{"invite_id":"invite-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"invite_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an invite removes a pending enrollment edge; already-created users require separate account analysis"],"completeness":{"state":"unknown","reason":"invite_delete_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.invite_delete", Severity: "blocking", Message: "deleting the invite may remove pending enrollment and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", inviteBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected invite delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyReturnsRegeneratedInviteTokenOnce(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"invite-1","email":"a@example.com"}`
	after := `{"id":"invite-1","email":"a@example.com","expires_at":"later"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "users.invites.regenerate", Request: json.RawMessage(`{"invite_id":"invite-1","expires_in":3600}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"invite_regenerate","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["regenerating an invite invalidates the previous enrollment token and creates a new one-time token"],"completeness":{"state":"unknown","reason":"invite_regenerate_requires_enrollment_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.invite_regenerate", Severity: "blocking", Message: "regenerating the invite invalidates its prior token and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", inviteBefore: []byte(before), inviteAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || result.OneTimeSecret != "replacement-invite" || remote.updates != 1 {
		t.Fatalf("unexpected invite regenerate result: %+v updates=%d", result, remote.updates)
	}
	receipt, err := store.GetReceipt(context.Background(), result.AttemptID)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(receipt.Result), "replacement-invite") {
		t.Fatal("replacement invite token leaked into persisted receipt")
	}
}

func TestApplyDispatchesAgentNetworkPolicyUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"policy-1","name":"engineering","enabled":true}`
	after := `{"id":"policy-1","name":"engineering","enabled":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.policies.update", Request: json.RawMessage(after), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(after), Impact: json.RawMessage(`{"classification":"agent_network_policy_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["updating an agent-network policy can alter provider reachability, guardrails, or admission limits for source groups; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_policy_update_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_policy_change", Severity: "blocking", Message: "the proposed agent-network policy change may alter provider reachability and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentPolicyBefore: []byte(before), agentPolicyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent policy update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesAgentNetworkPolicyDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"policy-1","name":"engineering","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", Operation: "agent_network.policies.delete", Request: json.RawMessage(`{"id":"policy-1"}`), Before: json.RawMessage(before), IntendedAfter: json.RawMessage(`{}`), Impact: json.RawMessage(`{"classification":"agent_network_policy_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting an agent-network policy can remove provider reachability and admission limits for source groups; affected peers and account resources require capability-aware live analysis"],"completeness":{"state":"unknown","reason":"agent_network_policy_delete_requires_capability_analysis"}}`), Findings: []ledger.Finding{{Code: "impact.agent_network_policy_delete", Severity: "blocking", Message: "deleting the agent-network policy may remove provider reachability and requires exact acknowledgement"}}})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", agentPolicyBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected agent policy delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesGroupDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"group","peers_count":2,"resources_count":1}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.delete",
		Request:        json.RawMessage(`{"id":"g1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"group_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a group can change policy membership and peer access; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"group_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.group_delete", Severity: "blocking", Message: "deleting the group may alter policy membership and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected group delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyRefusesDriftBeforeJournalingOrDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage := stageForTest(t, store, `{"id":"g1","name":"old"}`, `{"id":"g1","name":"new"}`)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"someone-else"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	_, err = Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || remote.updates != 0 {
		t.Fatalf("drift was not refused: err=%v updates=%d", err, remote.updates)
	}
	if _, err := store.GetAttempt(context.Background(), "missing"); err == nil {
		t.Fatal("unexpected attempt")
	}
}

func TestApplyDoesNotReplayAfterAmbiguousDispatch(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, Description: "connection lost"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.Unknown {
		t.Fatalf("unexpected ambiguous result: %+v err=%v", result, err)
	}
	attempt, err := store.GetAttempt(context.Background(), result.AttemptID)
	if err != nil || attempt.State != string(mutation.Unknown) {
		t.Fatalf("attempt=%+v err=%v", attempt, err)
	}
}

func TestApplyReturnsAlreadySatisfiedWithoutWrite(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, after, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(after), after: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err != nil || result.State != mutation.AlreadySatisfied || remote.updates != 0 {
		t.Fatalf("unexpected no-op result: %+v updates=%d err=%v", result, remote.updates, err)
	}
}

func TestApplyClassifiesDefinitiveRemoteRejection(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"g1","name":"old"}`
	after := `{"id":"g1","name":"new"}`
	stage := stageForTest(t, store, before, after)
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(before), after: []byte(after), updateErr: &transport.RequestError{Dispatched: true, StatusCode: 403, Description: "remote rejected request"}}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.State != mutation.DefinitivelyRejected {
		t.Fatalf("unexpected rejection result: %+v err=%v", result, err)
	}
}

func TestApplyRefusesChangedImpactEvidence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "groups.update",
		Request:        json.RawMessage(`{"id":"g1","name":"new"}`),
		Before:         json.RawMessage(`{"id":"g1","name":"old"}`),
		IntendedAfter:  json.RawMessage(`{"id":"g1","name":"new"}`),
		Impact:         json.RawMessage(`{"classification":"unknown"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", before: []byte(`{"id":"g1","name":"old"}`), after: []byte(`{"id":"g1","name":"new"}`)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1"})
	if err == nil || result.AttemptID != "" || remote.updates != 0 {
		t.Fatalf("changed impact was not refused before dispatch: result=%+v err=%v updates=%d", result, err, remote.updates)
	}
}

func TestApplyDispatchesPolicyUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"old","rules":[]}`
	after := `{"id":"p1","name":"new","rules":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.update",
		Request:        json.RawMessage(`{"id":"p1","name":"new","rules":[]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["policy metadata changed without changing policy rules"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyBefore: []byte(before), policyAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPolicyCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`
	created := `{"id":"p1","name":"allow-office","enabled":true,"rules":[{"action":"accept"}],"source_posture_checks":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.create",
		Request:        json.RawMessage(`{"name":"allow-office","enabled":true,"rules":[{"action":"accept"}]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"policy_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a policy can add access edges; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"policy_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.policy_create", Severity: "blocking", Message: "creating the policy may add reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyCollection: []byte(before), policyAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPolicyDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"policy","rules":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "policies.delete",
		Request:        json.RawMessage(`{"id":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"policy_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a policy can remove access edges; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"policy_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.policy_delete", Severity: "blocking", Message: "deleting the policy may remove access and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", policyBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected policy delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","description":"old","enabled":true,"metric":10,"groups":["g1"]}`
	after := `{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.update",
		Request:        json.RawMessage(`{"id":"r1","description":"new","enabled":true,"metric":10,"groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["route description changed without changing routing behavior"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeBefore: []byte(before), routeAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}`
	created := `{"id":"r1","description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"],"metric":10}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.create",
		Request:        json.RawMessage(`{"description":"private subnet","enabled":true,"network":"10.0.0.0/24","groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"route_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a route can add a reachable destination or alter path selection; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"route_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.route_create", Severity: "blocking", Message: "creating the route may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeCollection: []byte(before), routeAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesRouteDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","description":"route","enabled":true,"network":"10.0.0.0/24"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "routes.delete",
		Request:        json.RawMessage(`{"id":"r1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"route_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a route can change network reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"route_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.route_delete", Severity: "blocking", Message: "deleting the route may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routeBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected route delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPeerUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"old","approval_required":false,"connected":true}`
	after := `{"id":"p1","name":"new","approval_required":false,"connected":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "peers.update",
		Request:        json.RawMessage(`{"id":"p1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["peer name changed without changing peer access or connectivity state"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before), peerAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected peer result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesPeerDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"p1","name":"peer","connected":true,"approval_required":false}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "peers.delete",
		Request:        json.RawMessage(`{"id":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"peer_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a peer removes an access and connectivity principal; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"peer_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.peer_delete", Severity: "blocking", Message: "deleting the peer may remove access or connectivity and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", peerBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected peer delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"n1","name":"old","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	after := `{"id":"n1","name":"new","description":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.update",
		Request:        json.RawMessage(`{"id":"n1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["network metadata changed without changing attached policies, resources, or routers"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkBefore: []byte(before), networkAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"office","description":"primary"}`
	created := `{"id":"n1","name":"office","description":"primary","policies":[],"resources":[],"routers":[],"routing_peers_count":0}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.create",
		Request:        json.RawMessage(`{"name":"office","description":"primary"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network can add a topology scope for policies, resources, and routers; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_create", Severity: "blocking", Message: "creating the network may add topology and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkCollection: []byte(before), networkAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"n1","name":"office","policies":["p1"],"resources":["r1"],"routers":["rt1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.delete",
		Request:        json.RawMessage(`{"id":"n1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network can remove attached policies, resources, and routers; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_delete", Severity: "blocking", Message: "deleting the network may remove attached topology and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", networkBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.delete",
		Request:        json.RawMessage(`{"network_id":"n1","id":"r1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_resource_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network resource can remove a reachable destination; affected peers and routes require live topology analysis"],"completeness":{"state":"unknown","reason":"network_resource_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_resource_delete", Severity: "blocking", Message: "deleting the network resource may remove reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"r1","name":"old","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	after := `{"id":"r1","name":"new","description":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.update",
		Request:        json.RawMessage(`{"network_id":"n1","id":"r1","name":"new"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"metadata_only","reachability":"unchanged","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"high","evidence":["network resource metadata changed without changing its address, enablement, or group assignment"],"completeness":{"state":"complete","reason":null}}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceBefore: []byte(before), resourceAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource update result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkResourceCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`
	created := `{"id":"r1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"],"type":"subnet"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.resources.create",
		Request:        json.RawMessage(`{"network_id":"n1","name":"db","address":"10.0.0.0/24","enabled":true,"groups":["g1"]}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_resource_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network resource can add a reachable destination; affected peers and routes require live topology analysis"],"completeness":{"state":"unknown","reason":"network_resource_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_resource_create", Severity: "blocking", Message: "creating the network resource may add reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", resourceCollection: []byte(before), resourceAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network resource create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterCreateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `[]`
	intended := `{"enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`
	created := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":[]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.create",
		Request:        json.RawMessage(`{"network_id":"n1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(intended),
		Impact:         json.RawMessage(`{"classification":"network_router_create","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["creating a network router can add route availability and peer reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_create_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_create", Severity: "blocking", Message: "creating the network router may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerCollection: []byte(before), routerAfter: []byte(created)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router create result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterDeleteAndConfirmsAbsence(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1"}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.delete",
		Request:        json.RawMessage(`{"network_id":"n1","id":"rt1"}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(`{}`),
		Impact:         json.RawMessage(`{"classification":"network_router_delete","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["deleting a network router can change route availability and peer reachability; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_delete_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_delete", Severity: "blocking", Message: "deleting the network router may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerBefore: []byte(before)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router delete result: %+v updates=%d", result, remote.updates)
	}
}

func TestApplyDispatchesNetworkRouterUpdateAndConfirmsReadBack(t *testing.T) {
	store, err := ledger.Open(t.TempDir() + "/ledger.db")
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()
	before := `{"id":"rt1","enabled":true,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`
	after := `{"id":"rt1","enabled":false,"masquerade":true,"metric":10,"peer":"p1","peer_groups":["g1"]}`
	stage, err := store.Create(context.Background(), ledger.StageInput{
		Profile:        "default",
		ServerIdentity: "https://nb.test",
		AccountID:      "account-1",
		Operation:      "networks.routers.update",
		Request:        json.RawMessage(`{"network_id":"n1","id":"rt1","enabled":false}`),
		Before:         json.RawMessage(before),
		IntendedAfter:  json.RawMessage(after),
		Impact:         json.RawMessage(`{"classification":"network_router_change","reachability":"potentially_changed","affected_peer_ids":[],"affected_resource_ids":[],"confidence":"medium","evidence":["network router update changes reachability fields: [enabled]; affected peers and resources require live topology analysis"],"completeness":{"state":"unknown","reason":"network_router_change_requires_topology"}}`),
		Findings:       []ledger.Finding{{Code: "impact.network_router_change", Severity: "blocking", Message: "the proposed network router change may alter reachability and requires exact acknowledgement"}},
	})
	if err != nil {
		t.Fatal(err)
	}
	remote := &fakeRemote{identity: "https://nb.test", account: "account-1", routerBefore: []byte(before), routerAfter: []byte(after)}
	result, err := Apply(context.Background(), store, remote, ApplyInput{
		StageID: stage.ID, Revision: 1, Profile: "default", ServerIdentity: "https://nb.test", AccountID: "account-1", AckAllBlocking: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.State != mutation.ConfirmedSuccess || remote.updates != 1 {
		t.Fatalf("unexpected network router update result: %+v updates=%d", result, remote.updates)
	}
}
