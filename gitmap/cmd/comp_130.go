package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp130ID         = "38d66d9692ac"
	Comp130Uniqueness = "39bb88f40d3a"
	ErrComp130Fail    = "E_COMP_130_FAIL"
	OpHandleComp130   = "HandleComp130"
)

type Input130 struct {
	ID string
}

type Output130 struct {
	Result bool
}

func HandleComp130(in Input130) (Output130, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output130{Result: false}, apperror.New(OpHandleComp130, ErrComp130Fail, map[string]any{"id": in.ID})
	}

	return Output130{Result: true}, nil
}
