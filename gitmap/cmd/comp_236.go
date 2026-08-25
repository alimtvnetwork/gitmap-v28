package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp236BoundID    = "9a049b03f6fc"
	comp236Uniqueness = "b6cb293891dd"
	comp236ErrFail    = "E_COMP_236_FAIL"
	comp236OpHandle   = "HandleComp236"
)

// Input236 represents the input contract for component 236.
type Input236 struct {
	ID string
}

// Output236 represents the output contract for component 236.
type Output236 struct {
	Result bool
}

// HandleComp236 handles component 236 execution.
func HandleComp236(in Input236) (Output236, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output236{Result: false}, apperror.New(comp236OpHandle, comp236ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp236BoundID,
			"uniqueness": comp236Uniqueness,
		})
	}

	return Output236{Result: true}, nil
}
