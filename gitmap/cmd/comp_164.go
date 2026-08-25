package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp164ID         = "3f9807cb9ae9"
	Comp164Uniqueness = "2452984f72ef"
	ErrComp164Fail    = "E_COMP_164_FAIL"
	OpHandleComp164   = "HandleComp164"
)

type Input164 struct {
	ID string
}

type Output164 struct {
	Result bool
}

func HandleComp164(in Input164) (Output164, error) {
	if in.ID == Comp164Uniqueness {
		return Output164{Result: true}, nil
	}
	return Output164{Result: false}, apperror.New(OpHandleComp164, ErrComp164Fail, map[string]any{"id": in.ID})
}