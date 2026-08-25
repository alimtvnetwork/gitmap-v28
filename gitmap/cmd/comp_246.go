package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp246ID         = "37c20f19f327"
	Comp246Uniqueness = "23e8b0175874"
	ErrComp246Fail    = "E_COMP_246_FAIL"
	OpHandleComp246   = "HandleComp246"
)

// Input246 represents the input contract for component 246.
type Input246 struct {
	ID string
}

// Output246 represents the output contract for component 246.
type Output246 struct {
	Result bool
}

// HandleComp246 handles component 246 execution.
func HandleComp246(in Input246) (Output246, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output246{Result: false}, apperror.New(OpHandleComp246, ErrComp246Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp246ID,
			"uniqueness": Comp246Uniqueness,
		})
	}

	return Output246{Result: true}, nil
}
