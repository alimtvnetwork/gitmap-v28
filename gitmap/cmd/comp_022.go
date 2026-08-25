package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp022BoundID    = "785f3ec7eb32"
	comp022Uniqueness = "71ee45a3c0db"
	comp022ErrFail    = "E_COMP_022_FAIL"
	comp022OpHandle   = "HandleComp022"
)

// Input022 represents the input contract for component 022.
type Input022 struct {
	ID string
}

// Output022 represents the output contract for component 022.
type Output022 struct {
	Result bool
}

// HandleComp022 handles component 022 execution.
func HandleComp022(in Input022) (Output022, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output022{Result: false}, apperror.New(comp022OpHandle, comp022ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp022BoundID,
			"uniqueness": comp022Uniqueness,
		})
	}

	return Output022{Result: true}, nil
}
