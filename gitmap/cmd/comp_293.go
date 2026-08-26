package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp293ID         = "7cb676d57114"
	Comp293Uniqueness = "219de1387a67"
	ErrComp293Fail    = "E_COMP_293_FAIL"
	OpHandleComp293   = "HandleComp293"
)

type Input293 struct {
	ID string
}

type Output293 struct {
	Result bool
}

func HandleComp293(in Input293) (Output293, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output293{Result: false}, apperror.New(OpHandleComp293, ErrComp293Fail, map[string]any{"id": in.ID})
	}

	return Output293{Result: true}, nil
}
