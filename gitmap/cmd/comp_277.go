package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp277ID         = "27d719c754aa"
	Comp277Uniqueness = "833cd8c0e698"
	ErrComp277Fail    = "E_COMP_277_FAIL"
	OpHandleComp277   = "HandleComp277"
)

type Input277 struct {
	ID string
}

type Output277 struct {
	Result bool
}

func HandleComp277(in Input277) (Output277, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output277{Result: false}, apperror.New(OpHandleComp277, ErrComp277Fail, map[string]any{"id": in.ID})
	}

	return Output277{Result: true}, nil
}
