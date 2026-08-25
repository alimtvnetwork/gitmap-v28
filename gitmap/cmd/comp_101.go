package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp101ID         = "16dc368a89b4"
	Comp101Uniqueness = "c17edaae86e4"
	ErrComp101Fail    = "E_COMP_101_FAIL"
	OpHandleComp101   = "HandleComp101"
)

type Input101 struct {
	ID string
}

type Output101 struct {
	Result bool
}

func HandleComp101(in Input101) (Output101, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output101{Result: false}, apperror.New(OpHandleComp101, ErrComp101Fail, map[string]any{"id": in.ID})
	}

	return Output101{Result: true}, nil
}
