package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp251ID         = "c75d3f1f5bcd"
	Comp251Uniqueness = "5344c4110f48"
	ErrComp251Fail    = "E_COMP_251_FAIL"
	OpHandleComp251   = "HandleComp251"
)

type Input251 struct {
	ID string
}

type Output251 struct {
	Result bool
}

func HandleComp251(in Input251) (Output251, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output251{Result: false}, apperror.New(OpHandleComp251, ErrComp251Fail, map[string]any{"id": in.ID})
	}

	return Output251{Result: true}, nil
}
