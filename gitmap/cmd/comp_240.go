package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp240ID         = "6af1f692e949"
	Comp240Uniqueness = "ddfe0e8d462a"
	ErrComp240Fail    = "E_COMP_240_FAIL"
	OpHandleComp240   = "HandleComp240"
)

// Input240 represents the input contract for component 240.
type Input240 struct {
	ID string
}

// Output240 represents the output contract for component 240.
type Output240 struct {
	Result bool
}

// HandleComp240 handles component 240 execution.
func HandleComp240(in Input240) (Output240, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output240{Result: false}, apperror.New(OpHandleComp240, ErrComp240Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp240ID,
			"uniqueness": Comp240Uniqueness,
		})
	}

	return Output240{Result: true}, nil
}
