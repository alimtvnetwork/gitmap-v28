package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp287ID         = "d7cdaa5ca058"
	Comp287Uniqueness = "8e28c5eb829e"
	ErrComp287Fail    = "E_COMP_287_FAIL"
	OpHandleComp287   = "HandleComp287"
)

type Input287 struct {
	ID string
}

type Output287 struct {
	Result bool
}

func HandleComp287(in Input287) (Output287, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output287{Result: false}, apperror.New(OpHandleComp287, ErrComp287Fail, map[string]any{"id": in.ID})
	}

	return Output287{Result: true}, nil
}
