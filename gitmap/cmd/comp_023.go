package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp023BoundID    = "535fa30d7e25"
	comp023Uniqueness = "25fc0e7096fc"
	comp023ErrFail    = "E_COMP_023_FAIL"
	comp023OpHandle   = "HandleComp023"
)

// Input023 represents the input contract for component 023.
type Input023 struct {
	ID string
}

// Output023 represents the output contract for component 023.
type Output023 struct {
	Result bool
}

// HandleComp023 handles component 023 execution.
func HandleComp023(in Input023) (Output023, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output023{Result: false}, apperror.New(comp023OpHandle, comp023ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp023BoundID,
			"uniqueness": comp023Uniqueness,
		})
	}

	return Output023{Result: true}, nil
}
