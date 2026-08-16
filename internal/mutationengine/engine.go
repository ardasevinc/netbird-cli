package mutationengine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

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
	var request struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(stage.Request, &request); err != nil || request.ID == "" {
		return result, &ApplyError{Result: result, Err: errors.New("groups.update stage request requires a group id")}
	}
	findings := make([]mutation.Finding, 0, len(stage.Findings))
	for _, finding := range stage.Findings {
		findings = append(findings, mutation.Finding{Code: finding.Code, Severity: mutation.Severity(finding.Severity), Message: finding.Message})
	}
	if err := mutation.ValidateAcknowledgements(findings, input.Acknowledgements, input.AckAllBlocking); err != nil {
		return result, &ApplyError{Result: result, Err: err}
	}
	liveBefore, err := remote.GetGroup(ctx, request.ID)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("re-read group preimage: %w", err)}
	}
	preimage, err := mutation.ClassifyPreimage(stage.Before, liveBefore, stage.IntendedAfter)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("classify group preimage: %w", err)}
	}
	if preimage == mutation.PreimageDrifted {
		return result, &ApplyError{Result: result, Err: errors.New("staged group preimage drifted; create a new revision")}
	}
	impact, err := analysis.GroupUpdateImpact(liveBefore, stage.IntendedAfter)
	if err != nil {
		return result, &ApplyError{Result: result, Err: fmt.Errorf("recompute group mutation impact: %w", err)}
	}
	if len(stage.Impact) != 0 && string(stage.Impact) != "{}" {
		liveImpact, err := json.Marshal(impact)
		if err != nil {
			return result, &ApplyError{Result: result, Err: fmt.Errorf("encode group mutation impact: %w", err)}
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
	if _, err := remote.UpdateGroup(ctx, request.ID, stage.Request); err != nil {
		state := classifyDispatchError(err)
		return finish(ctx, store, result, state, "group update did not produce a confirmed success")
	}
	liveAfter, err := remote.GetGroup(ctx, request.ID)
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
