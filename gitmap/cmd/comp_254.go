package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp254ID         = "9512d95d00d6"
	Comp254Uniqueness = "ecac903ea62d"
	ErrComp254Fail    = "E_COMP_254_FAIL"
	OpHandleComp254   = "HandleComp254"
)

// Input254 represents the input contract for component 254.
type Input254 struct {
	ID string
}

// Output254 represents the output contract for component 254.
type Output254 struct {
	Result bool
}

// HandleComp254 handles component 254 execution.
func HandleComp254(in Input254) (Output254, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output254{Result: false}, apperror.New(OpHandleComp254, ErrComp254Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp254ID,
			"uniqueness": Comp254Uniqueness,
		})
	}

	return Output254{Result: true}, nil
}
