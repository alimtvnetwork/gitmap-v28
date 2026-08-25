package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp126ID         = "65a699905c02"
	Comp126Uniqueness = "d6e5a20b30f8"
	ErrComp126Fail    = "E_COMP_126_FAIL"
	OpHandleComp126   = "HandleComp126"
)

type Input126 struct {
	ID string
}

type Output126 struct {
	Result bool
}

func HandleComp126(in Input126) (Output126, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output126{Result: false}, apperror.New(OpHandleComp126, ErrComp126Fail, map[string]any{"id": in.ID})
	}

	return Output126{Result: true}, nil
}
