package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp080Uniqueness = "a512db2741cd"
	ErrComp080Fail    = "E_COMP_080_FAIL"
	OpHandleComp080   = "HandleComp080"
)

type Input080 struct {
	ID string
}

type Output080 struct {
	Result bool
}

func HandleComp080(in Input080) (Output080, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output080{Result: false}, apperror.New(OpHandleComp080, ErrComp080Fail, map[string]any{"id": in.ID})
	}

	return Output080{Result: true}, nil
}
