package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp025BoundID    = "b7a56873cd77"
	comp025Uniqueness = "1a6562590ef1"
	comp025ErrFail    = "E_COMP_025_FAIL"
	comp025OpHandle   = "HandleComp025"
)

// Input025 represents the input contract for component 025.
type Input025 struct {
	ID string
}

// Output025 represents the output contract for component 025.
type Output025 struct {
	Result bool
}

// HandleComp025 handles component 025 execution.
func HandleComp025(in Input025) (Output025, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output025{Result: false}, apperror.New(comp025OpHandle, comp025ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp025BoundID,
			"uniqueness": comp025Uniqueness,
		})
	}

	return Output025{Result: true}, nil
}
