package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp017BoundID    = "4523540f1504"
	comp017Uniqueness = "86e501496586"
	comp017ErrFail    = "E_COMP_017_FAIL"
	comp017OpHandle   = "HandleComp017"
)

// Input017 represents the input contract for component 017.
type Input017 struct {
	ID string
}

// Output017 represents the output contract for component 017.
type Output017 struct {
	Result bool
}

// HandleComp017 handles component 017 execution.
func HandleComp017(in Input017) (Output017, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output017{Result: false}, apperror.New(comp017OpHandle, comp017ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp017BoundID,
			"uniqueness": comp017Uniqueness,
		})
	}

	return Output017{Result: true}, nil
}
