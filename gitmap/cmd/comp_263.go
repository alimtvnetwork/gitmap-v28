package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp263ID         = "4be84111a613"
	Comp263Uniqueness = "f7c2599681e9"
	ErrComp263Fail    = "E_COMP_263_FAIL"
	OpHandleComp263   = "HandleComp263"
)

type Input263 struct {
	ID string
}

type Output263 struct {
	Result bool
}

func HandleComp263(in Input263) (Output263, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output263{Result: false}, apperror.New(OpHandleComp263, ErrComp263Fail, map[string]any{"id": in.ID})
	}

	return Output263{Result: true}, nil
}
