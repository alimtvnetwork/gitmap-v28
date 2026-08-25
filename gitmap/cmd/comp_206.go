package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp206ID         = "5cf4e26bd3d8"
	Comp206Uniqueness = "fabf5b7fedb3"
	ErrComp206Fail    = "E_COMP_206_FAIL"
	OpHandleComp206   = "HandleComp206"
)

type Input206 struct {
	ID string
}

type Output206 struct {
	Result bool
}

func HandleComp206(in Input206) (Output206, error) {
	if in.ID == Comp206Uniqueness {
		return Output206{Result: true}, nil
	}
	return Output206{Result: false}, apperror.New(OpHandleComp206, ErrComp206Fail, map[string]any{"id": in.ID})
}
