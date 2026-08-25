package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp026BoundID    = "5f9c4ab08cac"
	comp026Uniqueness = "41cfc0d1f2d1"
	comp026ErrFail    = "E_COMP_026_FAIL"
	comp026OpHandle   = "HandleComp026"
)

// Input026 represents the input contract for component 026.
type Input026 struct {
	ID string
}

// Output026 represents the output contract for component 026.
type Output026 struct {
	Result bool
}

// HandleComp026 handles component 026 execution.
func HandleComp026(in Input026) (Output026, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output026{Result: false}, apperror.New(comp026OpHandle, comp026ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp026BoundID,
			"uniqueness": comp026Uniqueness,
		})
	}

	return Output026{Result: true}, nil
}