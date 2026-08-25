package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp013BoundID    = "3fdba35f04dc"
	comp013Uniqueness = "5f9c4ab08cac"
	comp013ErrFail    = "E_COMP_013_FAIL"
	comp013OpHandle   = "HandleComp013"
)

// Input013 represents the input contract for component 013.
type Input013 struct {
	ID string
}

// Output013 represents the output contract for component 013.
type Output013 struct {
	Result bool
}

// HandleComp013 handles component 013 execution.
func HandleComp013(in Input013) (Output013, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output013{Result: false}, apperror.New(comp013OpHandle, comp013ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp013BoundID,
			"uniqueness": comp013Uniqueness,
		})
	}

	return Output013{Result: true}, nil
}
