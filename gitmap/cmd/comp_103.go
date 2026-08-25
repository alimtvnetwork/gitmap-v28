package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp103ID         = "454f63ac30c8"
	Comp103Uniqueness = "5cf4e26bd3d8"
	ErrComp103Fail    = "E_COMP_103_FAIL"
	OpHandleComp103   = "HandleComp103"
)

type Input103 struct {
	ID string
}

type Output103 struct {
	Result bool
}

func HandleComp103(in Input103) (Output103, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output103{Result: false}, apperror.New(OpHandleComp103, ErrComp103Fail, map[string]any{"id": in.ID})
	}

	return Output103{Result: true}, nil
}
