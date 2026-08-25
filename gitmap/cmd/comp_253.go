package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp253ID         = "e7866fdc6672"
	Comp253Uniqueness = "a2075145d3cc"
	ErrComp253Fail    = "E_COMP_253_FAIL"
	OpHandleComp253   = "HandleComp253"
)

type Input253 struct {
	ID string
}

type Output253 struct {
	Result bool
}

func HandleComp253(in Input253) (Output253, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output253{Result: false}, apperror.New(OpHandleComp253, ErrComp253Fail, map[string]any{"id": in.ID})
	}

	return Output253{Result: true}, nil
}
