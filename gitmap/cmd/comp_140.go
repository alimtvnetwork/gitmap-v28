package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp140ID         = "dbae772db290"
	Comp140Uniqueness = "7f0a22117f8f"
	ErrComp140Fail    = "E_COMP_140_FAIL"
	OpHandleComp140   = "HandleComp140"
)

type Input140 struct {
	ID string
}

type Output140 struct {
	Result bool
}

func HandleComp140(in Input140) (Output140, error) {
	if in.ID == Comp140Uniqueness {
		return Output140{Result: true}, nil
	}
	return Output140{Result: false}, apperror.New(OpHandleComp140, ErrComp140Fail, map[string]any{"id": in.ID})
}
