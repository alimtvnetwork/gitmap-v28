package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp113ID         = "6c658ee83fb7"
	Comp113Uniqueness = "8f1f64db81c4"
	ErrComp113Fail    = "E_COMP_113_FAIL"
	OpHandleComp113   = "HandleComp113"
)

type Input113 struct {
	ID string
}

type Output113 struct {
	Result bool
}

func HandleComp113(in Input113) (Output113, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output113{Result: false}, apperror.New(OpHandleComp113, ErrComp113Fail, map[string]any{"id": in.ID})
	}

	return Output113{Result: true}, nil
}
