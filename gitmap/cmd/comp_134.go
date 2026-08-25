package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp134ID         = "5d389f5e2e34"
	Comp134Uniqueness = "8b496bf96bbc"
	ErrComp134Fail    = "E_COMP_134_FAIL"
	OpHandleComp134   = "HandleComp134"
)

type Input134 struct {
	ID string
}

type Output134 struct {
	Result bool
}

func HandleComp134(in Input134) (Output134, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output134{Result: false}, apperror.New(OpHandleComp134, ErrComp134Fail, map[string]any{"id": in.ID})
	}

	return Output134{Result: true}, nil
}
