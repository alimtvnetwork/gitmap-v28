package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp275ID         = "3a1dfb05d725"
	Comp275Uniqueness = "f89f8d0e735a"
	ErrComp275Fail    = "E_COMP_275_FAIL"
	OpHandleComp275   = "HandleComp275"
)

type Input275 struct {
	ID string
}

type Output275 struct {
	Result bool
}

func HandleComp275(in Input275) (Output275, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output275{Result: false}, apperror.New(OpHandleComp275, ErrComp275Fail, map[string]any{"id": in.ID})
	}

	return Output275{Result: true}, nil
}
