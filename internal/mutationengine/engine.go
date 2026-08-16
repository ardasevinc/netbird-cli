package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/ardasevinc/netbird-cli/internal/analysis"
	"github.com/ardasevinc/netbird-cli/internal/ledger"
	"github.com/ardasevinc/netbird-cli/internal/mutation"
	"github.com/ardasevinc/netbird-cli/internal/operations"
)

// Remote is the smallest adapter surface needed by the first consequential
// operation. It deliberately exposes no arbitrary method-and-URL dispatcher.
type Remote interface {
	ServerIdentity() string
	AccountScope(context.Context, string) error
	GetGroup(context.Context, string) (json.RawMessage, error)
	UpdateGroup(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteGroup(context.Context, string) (json.RawMessage, error)
	GetPolicyRaw(context.Context, string) (json.RawMessage, error)
	UpdatePolicy(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePolicy(context.Context, string) (json.RawMessage, error)
	GetRouteRaw(context.Context, string) (json.RawMessage, error)
	UpdateRoute(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteRoute(context.Context, string) (json.RawMessage, error)
	GetPeerRaw(context.Context, string) (json.RawMessage, error)
	UpdatePeer(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeletePeer(context.Context, string) (json.RawMessage, error)
	GetNetworkRaw(context.Context, string) (json.RawMessage, error)
	UpdateNetwork(context.Context, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetwork(context.Context, string) (json.RawMessage, error)
	GetNetworkResourceRaw(context.Context, string, string) (json.RawMessage, error)
	UpdateNetworkResource(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetworkResource(context.Context, string, string) (json.RawMessage, error)
	GetNetworkRouterRaw(context.Context, string, string) (json.RawMessage, error)
	UpdateNetworkRouter(context.Context, string, string, json.RawMessage) (json.RawMessage, error)
	DeleteNetworkRouter(context.Context, string, string) (json.RawMessage, error)
}

type Ledger interface {
	Get(context.Context, string, int) (ledger.Stage, error)
	BeginAttempt(context.Context, string, int) (ledger.Attempt, error)
	SetAttemptState(context.Context, string, string) error
	RecordReceipt(context.Context, ledger.Receipt) error
}

type ApplyInput struct {
	StageID          string
	Revision         int
	Profile          string
	ServerIdentity   string
	AccountID        string
	Acknowledgements []string
	AckAllBlocking   bool
}

type requestTarget struct {
	ID        string `json:"id"`
	NetworkID string `json:"network_id"`
}

type Result struct {
	StageID   string                 `json:"stage_id"`
	Revision  int                    `json:"revision"`
	AttemptID string                 `json:"attempt_id,omitempty"`
	State     mutation.DispatchState `json:"state"`
	Reason    string                 `json:"reason,omitempty"`
}

type ApplyError struct {
	Result Result
	Err    error
}

func (e *ApplyError) Error() string { return e.Err.Error() }

func (e *ApplyError) Unwrap() error { return e.Err }

func Apply(ctx context.Context, store Ledger, remote Remote, input ApplyInput) (Result, error) {
	if input.StageID == "" || input.Revision < 1 {
		return Result{}, &ApplyError{Err: errors.New("apply requires an exact stage id and positive revision")}
	}
	stage, err := store.Get(ctx, input.StageID, input.Revision)
	if err != nil {
		return Result{}, &ApplyError{Err: err}
	}
	result := Result{StageID: stage.ID, Revision: stage.Revision}
	if stage.Cancelled {
		return result, &ApplyError{Result: result, Err: errors.New("stage revision is cancelled")}
	}
	definition, err := operations.Lookup(stage.Operation)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	if definition.Mutation != operations.Consequential || !definition.DispatcherAdmitted {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("operation %q is not admitted for dispatch", stage.Operation)}
	}
	if input.Profile == "" || stage.Profile != input.Profile {
		return result, &ApplyError{Result: result, Err: errors.New("stage profile does not match the selected profile")}
	}
	if input.ServerIdentity == "" || stage.ServerIdentity == "" || stage.ServerIdentity != input.ServerIdentity || remote.ServerIdentity() != input.ServerIdentity {
		return result, &ApplyError{Result: result, Err: errors.New("server identity does not match the staged identity")}
	}
	if input.AccountID == "" || stage.AccountID == "" || stage.AccountID != input.AccountID {
		return result, &ApplyError{Result: result, Err: errors.New("account scope does not match the staged account")}
	}
	if err := remote.AccountScope(ctx, input.AccountID); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	var request requestTarget
	if err := json.Unmarshal(stage.Request, &request); err != nil || request.ID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires a target id", stage.Operation)}
	}
	if (stage.Operation == "networks.resources.update" || stage.Operation == "networks.resources.delete" || stage.Operation == "networks.routers.update" || stage.Operation == "networks.routers.delete") && request.NetworkID == "" {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("%s stage request requires network_id", stage.Operation)}
	}
	findings := make([]mutation.Finding, 0, len(stage.Findings))
	for _, finding := range stage.Findings {
		findings = append(findings, mutation.Finding{Code: finding.Code, Severity: mutation.Severity(finding.Severity), Message: finding.Message})
	}
	if err := mutation.ValidateAcknowledgements(findings, input.Acknowledgements, input.AckAllBlocking); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	liveBefore, err := readPreimage(ctx, remote, stage.Operation, request)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("re-read %s preimage: %w", stage.Operation, err)}
	}
	preimage, err := mutation.ClassifyPreimage(stage.Before, liveBefore, stage.IntendedAfter)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("classify %s preimage: %w", stage.Operation, err)}
	}
	if preimage == mutation.PreimageDrifted {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("staged %s preimage drifted; create a new revision", stage.Operation)}
	}
	impact, err := mutationImpact(stage.Operation, liveBefore, stage.IntendedAfter)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("recompute %s mutation impact: %w", stage.Operation, err)}
	}
	if len(stage.Impact) != 0 && string(stage.Impact) != "{}" {
		liveImpact, err := json.Marshal(impact)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("encode %s mutation impact: %w", stage.Operation, err)}
		}
		equal, err := mutation.Equivalent(stage.Impact, liveImpact)
		if err != nil || !equal {
			return result, &ApplyError{Result: result, Err: errors.New("staged mutation impact changed; create a new revision")}
		}
	}
	attempt, err := store.BeginAttempt(ctx, stage.ID, stage.Revision)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	result.AttemptID = attempt.ID
	if preimage == mutation.PreimageAlreadySatisfied {
		return finish(ctx, store, result, mutation.AlreadySatisfied, "remote state already equals intended state")
	}
	if _, err := dispatch(ctx, remote, stage.Operation, request, stage.Request); err != nil {
		state := classifyDispatchError(err)
		return finish(ctx, store, result, state, stage.Operation+" did not produce a confirmed success")
	}
	if isDeleteOperation(stage.Operation) {
		if err := confirmDeleted(ctx, remote, stage.Operation, request); err != nil && !isNotFound(err) {
			return finish(ctx, store, result, mutation.Unknown, "delete may have applied, but absence could not be confirmed")
		}
		return finish(ctx, store, result, mutation.ConfirmedSuccess, "remote "+strings.TrimSuffix(stage.Operation, ".delete")+" is absent after delete")
	}
	liveAfter, err := readPreimage(ctx, remote, stage.Operation, request)
	if err != nil {
		return finish(ctx, store, result, mutation.Unknown, "update may have applied, but read-back was inconclusive")
	}
	equal, err := mutation.Equivalent(liveAfter, stage.IntendedAfter)
	if err != nil {
		return finish(ctx, store, result, mutation.Unknown, "update response could not be compared with intended state")
	}
	if !equal {
		return finish(ctx, store, result, mutation.Partial, "remote state differs from intended state after update")
	}
	return finish(ctx, store, result, mutation.ConfirmedSuccess, "remote state matches intended state")
}

func readPreimage(ctx context.Context, remote Remote, operation string, target requestTarget) (json.RawMessage, error) {
	switch operation {
	case "groups.update":
		return remote.GetGroup(ctx, target.ID)
	case "groups.delete":
		return remote.GetGroup(ctx, target.ID)
	case "policies.update":
		return remote.GetPolicyRaw(ctx, target.ID)
	case "policies.delete":
		return remote.GetPolicyRaw(ctx, target.ID)
	case "routes.update":
		return remote.GetRouteRaw(ctx, target.ID)
	case "routes.delete":
		return remote.GetRouteRaw(ctx, target.ID)
	case "peers.update":
		return remote.GetPeerRaw(ctx, target.ID)
	case "peers.delete":
		return remote.GetPeerRaw(ctx, target.ID)
	case "networks.update":
		return remote.GetNetworkRaw(ctx, target.ID)
	case "networks.delete":
		return remote.GetNetworkRaw(ctx, target.ID)
	case "networks.resources.delete":
		return remote.GetNetworkResourceRaw(ctx, target.NetworkID, target.ID)
	case "networks.resources.update":
		return remote.GetNetworkResourceRaw(ctx, target.NetworkID, target.ID)
	case "networks.routers.delete":
		return remote.GetNetworkRouterRaw(ctx, target.NetworkID, target.ID)
	case "networks.routers.update":
		return remote.GetNetworkRouterRaw(ctx, target.NetworkID, target.ID)
	default:
		return nil, fmt.Errorf("operation %q has no preimage reader", operation)
	}
}

func dispatch(ctx context.Context, remote Remote, operation string, target requestTarget, request json.RawMessage) (json.RawMessage, error) {
	switch operation {
	case "groups.update":
		return remote.UpdateGroup(ctx, target.ID, request)
	case "groups.delete":
		return remote.DeleteGroup(ctx, target.ID)
	case "policies.update":
		return remote.UpdatePolicy(ctx, target.ID, request)
	case "policies.delete":
		return remote.DeletePolicy(ctx, target.ID)
	case "routes.update":
		return remote.UpdateRoute(ctx, target.ID, request)
	case "routes.delete":
		return remote.DeleteRoute(ctx, target.ID)
	case "peers.update":
		return remote.UpdatePeer(ctx, target.ID, request)
	case "peers.delete":
		return remote.DeletePeer(ctx, target.ID)
	case "networks.update":
		return remote.UpdateNetwork(ctx, target.ID, request)
	case "networks.delete":
		return remote.DeleteNetwork(ctx, target.ID)
	case "networks.resources.delete":
		return remote.DeleteNetworkResource(ctx, target.NetworkID, target.ID)
	case "networks.resources.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNetworkResource(ctx, target.NetworkID, target.ID, body)
	case "networks.routers.delete":
		return remote.DeleteNetworkRouter(ctx, target.NetworkID, target.ID)
	case "networks.routers.update":
		body, err := stripTargetFields(request)
		if err != nil {
			return nil, fmt.Errorf("prepare %s request: %w", operation, err)
		}
		return remote.UpdateNetworkRouter(ctx, target.NetworkID, target.ID, body)
	default:
		return nil, fmt.Errorf("operation %q has no dispatcher", operation)
	}
}

func stripTargetFields(request json.RawMessage) (json.RawMessage, error) {
	var object map[string]any
	if err := json.Unmarshal(request, &object); err != nil {
		return nil, err
	}
	delete(object, "id")
	delete(object, "network_id")
	return json.Marshal(object)
}

func mutationImpact(operation string, before, intendedAfter json.RawMessage) (analysis.ImpactReport, error) {
	switch operation {
	case "groups.update":
		return analysis.GroupUpdateImpact(before, intendedAfter)
	case "groups.delete":
		return analysis.GroupDeleteImpact(before)
	case "policies.update":
		return analysis.PolicyUpdateImpact(before, intendedAfter)
	case "policies.delete":
		return analysis.PolicyDeleteImpact(before)
	case "routes.update":
		return analysis.RouteUpdateImpact(before, intendedAfter)
	case "routes.delete":
		return analysis.RouteDeleteImpact(before)
	case "peers.update":
		return analysis.PeerUpdateImpact(before, intendedAfter)
	case "peers.delete":
		return analysis.PeerDeleteImpact(before)
	case "networks.update":
		return analysis.NetworkUpdateImpact(before, intendedAfter)
	case "networks.delete":
		return analysis.NetworkDeleteImpact(before)
	case "networks.resources.delete":
		return analysis.NetworkResourceDeleteImpact(before)
	case "networks.resources.update":
		return analysis.NetworkResourceUpdateImpact(before, intendedAfter)
	case "networks.routers.delete":
		return analysis.NetworkRouterDeleteImpact(before)
	case "networks.routers.update":
		return analysis.NetworkRouterUpdateImpact(before, intendedAfter)
	default:
		return analysis.ImpactReport{}, fmt.Errorf("operation %q has no impact analyzer", operation)
	}
}

func confirmDeleted(ctx context.Context, remote Remote, operation string, target requestTarget) error {
	_, err := readPreimage(ctx, remote, operation, target)
	if err == nil {
		return errors.New("resource still exists after delete")
	}
	return err
}

func isNotFound(err error) bool {
	var status interface{ StatusCodeState() int }
	return errors.As(err, &status) && status.StatusCodeState() == 404
}

func isDeleteOperation(operation string) bool {
	return operation == "groups.delete" || operation == "policies.delete" || operation == "routes.delete" || operation == "peers.delete" || operation == "networks.delete" || operation == "networks.resources.delete" || operation == "networks.routers.delete"
}

func classifyDispatchError(err error) mutation.DispatchState {
	var dispatched interface{ DispatchedState() bool }
	if errors.As(err, &dispatched) {
		if !dispatched.DispatchedState() {
			return mutation.NotDispatched
		}
		var status interface{ StatusCodeState() int }
		if errors.As(err, &status) && status.StatusCodeState() >= 400 && status.StatusCodeState() < 500 {
			return mutation.DefinitivelyRejected
		}
	}
	return mutation.Unknown
}

func finish(ctx context.Context, store Ledger, result Result, state mutation.DispatchState, reason string) (Result, error) {
	result.State = state
	result.Reason = reason
	if err := store.SetAttemptState(ctx, result.AttemptID, string(state)); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	encoded, err := json.Marshal(result)
	if err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	if err := store.RecordReceipt(ctx, ledger.Receipt{AttemptID: result.AttemptID, StageID: result.StageID, Revision: result.Revision, State: string(state), Result: encoded}); err != nil {
		if state == mutation.ConfirmedSuccess || state == mutation.AlreadySatisfied {
			result.State = mutation.EffectConfirmedReceiptFail
		}
		return result, &ApplyError{Result: result, Err: fmt.Errorf("persist mutation receipt: %w", err)}
	}
	if state == mutation.ConfirmedSuccess || state == mutation.AlreadySatisfied {
		return result, nil
	}
	return result, &ApplyError{Result: result, Err: errors.New(reason)}
}
