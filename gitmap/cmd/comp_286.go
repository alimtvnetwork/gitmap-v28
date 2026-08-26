package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp286ID         = "00328ce57bbc"
	Comp286Uniqueness = "5e74cb2ad4e2"
	ErrComp286Fail    = "E_COMP_286_FAIL"
	OpHandleComp286   = "HandleComp286"
)

type Input286 struct {
	ID string
}

type Output286 struct {
	Result bool
}

func HandleComp286(in Input286) (Output286, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output286{Result: false}, apperror.New(OpHandleComp286, ErrComp286Fail, map[string]any{"id": in.ID})
	}

	return Output286{Result: true}, nil
}
