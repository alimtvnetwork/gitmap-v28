package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp165ID         = "bc52dd634277"
	Comp165Uniqueness = "5426d2ca50f2"
	ErrComp165Fail    = "E_COMP_165_FAIL"
	OpHandleComp165   = "HandleComp165"
)

type Input165 struct {
	ID string
}

type Output165 struct {
	Result bool
}

func HandleComp165(in Input165) (Output165, error) {
	if in.ID == Comp165Uniqueness {
		return Output165{Result: true}, nil
	}
	return Output165{Result: false}, apperror.New(OpHandleComp165, ErrComp165Fail, map[string]any{"id": in.ID})
}
