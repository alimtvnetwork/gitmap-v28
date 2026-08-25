package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp250ID         = "1e472b39b105"
	Comp250Uniqueness = "0604cd3138fe"
	ErrComp250Fail    = "E_COMP_250_FAIL"
	OpHandleComp250   = "HandleComp250"
)

// Input250 represents the input contract for component 250.
type Input250 struct {
	ID string
}

// Output250 represents the output contract for component 250.
type Output250 struct {
	Result bool
}

// HandleComp250 handles component 250 execution.
func HandleComp250(in Input250) (Output250, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output250{Result: false}, apperror.New(OpHandleComp250, ErrComp250Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp250ID,
			"uniqueness": Comp250Uniqueness,
		})
	}

	return Output250{Result: true}, nil
}
