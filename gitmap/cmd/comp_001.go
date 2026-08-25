package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp001ID         = "6b86b273ff34"
	Comp001Uniqueness = "d4735e3a265e"
	ErrComp001Fail    = "E_COMP_001_FAIL"
	OpHandleComp001   = "HandleComp001"
)

type Input001 struct {
	ID string
}

type Output001 struct {
	Result bool
}

func HandleComp001(in Input001) (Output001, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output001{Result: false}, apperror.New(OpHandleComp001, ErrComp001Fail, map[string]any{"id": in.ID})
	}

	return Output001{Result: true}, nil
}
