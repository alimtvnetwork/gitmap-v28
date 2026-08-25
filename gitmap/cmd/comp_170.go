package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp170ID         = "734d0759cdb4"
	Comp170Uniqueness = "9644294ac4ff"
	ErrComp170Fail    = "E_COMP_170_FAIL"
	OpHandleComp170   = "HandleComp170"
)

type Input170 struct {
	ID string
}

type Output170 struct {
	Result bool
}

func HandleComp170(in Input170) (Output170, error) {
	if in.ID == Comp170Uniqueness {
		return Output170{Result: true}, nil
	}
	return Output170{Result: false}, apperror.New(OpHandleComp170, ErrComp170Fail, map[string]any{"id": in.ID})
}
