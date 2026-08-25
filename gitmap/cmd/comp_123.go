package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp123ID         = "a665a4592042"
	Comp123Uniqueness = "37c20f19f327"
	ErrComp123Fail    = "E_COMP_123_FAIL"
	OpHandleComp123   = "HandleComp123"
)

type Input123 struct {
	ID string
}

type Output123 struct {
	Result bool
}

func HandleComp123(in Input123) (Output123, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output123{Result: false}, apperror.New(OpHandleComp123, ErrComp123Fail, map[string]any{"id": in.ID})
	}

	return Output123{Result: true}, nil
}
