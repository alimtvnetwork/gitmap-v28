package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp204ID         = "fc56dbc6d465"
	Comp204Uniqueness = "e6f47e008cc5"
	ErrComp204Fail    = "E_COMP_204_FAIL"
	OpHandleComp204   = "HandleComp204"
)

type Input204 struct {
	ID string
}

type Output204 struct {
	Result bool
}

func HandleComp204(in Input204) (Output204, error) {
	if in.ID == Comp204Uniqueness {
		return Output204{Result: true}, nil
	}
	return Output204{Result: false}, apperror.New(OpHandleComp204, ErrComp204Fail, map[string]any{"id": in.ID})
}
