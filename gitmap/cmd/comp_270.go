package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp270ID         = "d8d1790737d5"
	Comp270Uniqueness = "84f01dd97c68"
	ErrComp270Fail    = "E_COMP_270_FAIL"
	OpHandleComp270   = "HandleComp270"
)

type Input270 struct {
	ID string
}

type Output270 struct {
	Result bool
}

func HandleComp270(in Input270) (Output270, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output270{Result: false}, apperror.New(OpHandleComp270, ErrComp270Fail, map[string]any{"id": in.ID})
	}

	return Output270{Result: true}, nil
}
