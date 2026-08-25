package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp024BoundID    = "c2356069e9d1"
	comp024Uniqueness = "98010bd9270f"
	comp024ErrFail    = "E_COMP_024_FAIL"
	comp024OpHandle   = "HandleComp024"
)

// Input024 represents the input contract for component 024.
type Input024 struct {
	ID string
}

// Output024 represents the output contract for component 024.
type Output024 struct {
	Result bool
}

// HandleComp024 handles component 024 execution.
func HandleComp024(in Input024) (Output024, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output024{Result: false}, apperror.New(comp024OpHandle, comp024ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp024BoundID,
			"uniqueness": comp024Uniqueness,
		})
	}

	return Output024{Result: true}, nil
}
