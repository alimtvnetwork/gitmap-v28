package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp220ID         = "36790ecd55c2"
	Comp220Uniqueness = "e3f6959781c3"
	ErrComp220Fail    = "E_COMP_220_FAIL"
	OpHandleComp220   = "HandleComp220"
)

type Input220 struct {
	ID string
}

type Output220 struct {
	Result bool
}

func HandleComp220(in Input220) (Output220, error) {
	if in.ID == Comp220Uniqueness {
		return Output220{Result: true}, nil
	}
	return Output220{Result: false}, apperror.New(OpHandleComp220, ErrComp220Fail, map[string]any{"id": in.ID})
}
