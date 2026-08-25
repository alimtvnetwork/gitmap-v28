package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp093ID         = "6e4001871c0c"
	Comp093Uniqueness = "2811745d7b8d"
	ErrComp093Fail    = "E_COMP_093_FAIL"
	OpHandleComp093   = "HandleComp093"
)

type Input093 struct {
	ID string
}

type Output093 struct {
	Result bool
}

func HandleComp093(in Input093) (Output093, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output093{Result: false}, apperror.New(OpHandleComp093, ErrComp093Fail, map[string]any{"id": in.ID})
	}

	return Output093{Result: true}, nil
}
