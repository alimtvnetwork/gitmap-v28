package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp111ID         = "f6e0a1e2ac41"
	Comp111Uniqueness = "9b871512327c"
	ErrComp111Fail    = "E_COMP_111_FAIL"
	OpHandleComp111   = "HandleComp111"
)

type Input111 struct {
	ID string
}

type Output111 struct {
	Result bool
}

func HandleComp111(in Input111) (Output111, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output111{Result: false}, apperror.New(OpHandleComp111, ErrComp111Fail, map[string]any{"id": in.ID})
	}

	return Output111{Result: true}, nil
}
