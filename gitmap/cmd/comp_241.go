package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp241ID         = "749fc650cacb"
	Comp241Uniqueness = "d4679c618f1a"
	ErrComp241Fail    = "E_COMP_241_FAIL"
	OpHandleComp241   = "HandleComp241"
)

// Input241 represents the input contract for component 241.
type Input241 struct {
	ID string
}

// Output241 represents the output contract for component 241.
type Output241 struct {
	Result bool
}

// HandleComp241 handles component 241 execution.
func HandleComp241(in Input241) (Output241, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty || in.ID != Comp241Uniqueness {
		return Output241{Result: false}, apperror.New(OpHandleComp241, ErrComp241Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp241ID,
			"uniqueness": Comp241Uniqueness,
		})
	}

	return Output241{Result: true}, nil
}
