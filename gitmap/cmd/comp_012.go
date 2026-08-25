package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp012BoundID    = "6b51d431df5d"
	comp012Uniqueness = "c2356069e9d1"
	comp012ErrFail    = "E_COMP_012_FAIL"
	comp012OpHandle   = "HandleComp012"
)

// Input012 represents the input contract for component 012.
type Input012 struct {
	ID string
}

// Output012 represents the output contract for component 012.
type Output012 struct {
	Result bool
}

// HandleComp012 handles component 012 execution.
func HandleComp012(in Input012) (Output012, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output012{Result: false}, apperror.New(comp012OpHandle, comp012ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp012BoundID,
			"uniqueness": comp012Uniqueness,
		})
	}

	return Output012{Result: true}, nil
}
