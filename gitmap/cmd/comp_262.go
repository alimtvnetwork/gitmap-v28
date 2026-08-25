package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp262ID         = "9e6a72557ada"
	Comp262Uniqueness = "388c2eafe5af"
	ErrComp262Fail    = "E_COMP_262_FAIL"
	OpHandleComp262   = "HandleComp262"
)

type Input262 struct {
	ID string
}

type Output262 struct {
	Result bool
}

func HandleComp262(in Input262) (Output262, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output262{Result: false}, apperror.New(OpHandleComp262, ErrComp262Fail, map[string]any{"id": in.ID})
	}

	return Output262{Result: true}, nil
}
