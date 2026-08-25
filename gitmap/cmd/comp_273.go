package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp273ID         = "303c8bd55875"
	Comp273Uniqueness = "6fc8f95bc646"
	ErrComp273Fail    = "E_COMP_273_FAIL"
	OpHandleComp273   = "HandleComp273"
)

type Input273 struct {
	ID string
}

type Output273 struct {
	Result bool
}

func HandleComp273(in Input273) (Output273, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output273{Result: false}, apperror.New(OpHandleComp273, ErrComp273Fail, map[string]any{"id": in.ID})
	}

	return Output273{Result: true}, nil
}
