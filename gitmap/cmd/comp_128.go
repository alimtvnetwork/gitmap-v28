package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp128ID         = "2747b7c71856"
	Comp128Uniqueness = "51e8ea280b44"
	ErrComp128Fail    = "E_COMP_128_FAIL"
	OpHandleComp128   = "HandleComp128"
)

type Input128 struct {
	ID string
}

type Output128 struct {
	Result bool
}

func HandleComp128(in Input128) (Output128, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output128{Result: false}, apperror.New(OpHandleComp128, ErrComp128Fail, map[string]any{"id": in.ID})
	}

	return Output128{Result: true}, nil
}
