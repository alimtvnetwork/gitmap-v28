package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp201ID         = "43974ed74066"
	Comp201Uniqueness = "b7c7470e59e2"
	ErrComp201Fail    = "E_COMP_201_FAIL"
	OpHandleComp201   = "HandleComp201"
)

type Input201 struct {
	ID string
}

type Output201 struct {
	Result bool
}

func HandleComp201(in Input201) (Output201, error) {
	if in.ID == Comp201Uniqueness {
		return Output201{Result: true}, nil
	}
	return Output201{Result: false}, apperror.New(OpHandleComp201, ErrComp201Fail, map[string]any{"id": in.ID})
}
