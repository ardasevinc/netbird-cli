package mutation

import (
	"errors"
	"fmt"
	"sort"
)

type Severity string

const (
	Info     Severity = "info"
	Warning  Severity = "warning"
	Blocking Severity = "blocking"
)

type Finding struct {
	Code     string   `json:"code"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
}

type DispatchState string

const (
	NotDispatched              DispatchState = "not_dispatched"
	DefinitivelyRejected       DispatchState = "definitively_rejected"
	ConfirmedSuccess           DispatchState = "confirmed_success"
	AlreadySatisfied           DispatchState = "already_satisfied"
	Partial                    DispatchState = "partial"
	Unknown                    DispatchState = "unknown"
	EffectConfirmedReceiptFail DispatchState = "effect_confirmed_receipt_failed"
)

func BlockingCodes(findings []Finding) []string {
	codes := make([]string, 0)
	for _, finding := range findings {
		if finding.Severity == Blocking {
			codes = append(codes, finding.Code)
		}
	}
	sort.Strings(codes)
	return codes
}

func ValidateAcknowledgements(findings []Finding, acknowledgements []string, ackAll bool) error {
	blocking := BlockingCodes(findings)
	if len(blocking) == 0 {
		return nil
	}
	if ackAll {
		return nil
	}
	provided := make(map[string]struct{}, len(acknowledgements))
	for _, code := range acknowledgements {
		provided[code] = struct{}{}
	}
	missing := make([]string, 0)
	for _, code := range blocking {
		if _, ok := provided[code]; !ok {
			missing = append(missing, code)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("blocking findings require acknowledgement: %v", missing)
	}
	return nil
}

func RetryAllowed(state DispatchState) bool {
	return state == NotDispatched
}

func ValidateDispatchTransition(from, to DispatchState) error {
	if from == "" {
		return errors.New("dispatch state is empty")
	}
	if from == Unknown || from == Partial {
		return errors.New("uncertain dispatch state cannot transition without reconciliation")
	}
	return nil
}
