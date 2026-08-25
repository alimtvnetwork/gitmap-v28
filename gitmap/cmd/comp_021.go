package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp021BoundID    = "6f4b6612125f"
	comp021Uniqueness = "73475cb40a56"
	comp021ErrFail    = "E_COMP_021_FAIL"
	comp021OpHandle   = "HandleComp021"
)

// Input021 represents the input contract for component 021.
type Input021 struct {
	ID string
}

// Output021 represents the output contract for component 021.
type Output021 struct {
	Result bool
}

// HandleComp021 handles component 021 execution.
func HandleComp021(in Input021) (Output021, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output021{Result: false}, apperror.New(comp021OpHandle, comp021ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp021BoundID,
			"uniqueness": comp021Uniqueness,
		})
	}

	return Output021{Result: true}, nil
}
