package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp008BoundID    = "2c624232cdd2"
	comp008Uniqueness = "b17ef6d19c7a"
	comp008ErrFail    = "E_COMP_008_FAIL"
	comp008OpHandle   = "HandleComp008"
)

// Input008 represents the input contract for component 008.
type Input008 struct {
	ID string
}

// Output008 represents the output contract for component 008.
type Output008 struct {
	Result bool
}

// HandleComp008 handles component 008 execution.
func HandleComp008(input Input008) (Output008, error) {
	isEmpty := len(input.ID) == 0
	if isEmpty {
		return Output008{Result: false}, apperror.New(comp008OpHandle, comp008ErrFail, map[string]any{
			"id":         input.ID,
			"bound_id":   comp008BoundID,
			"uniqueness": comp008Uniqueness,
		})
	}

	return Output008{Result: true}, nil
}
