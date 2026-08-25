package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp107ID         = "3346f2bbf6c3"
	Comp107Uniqueness = "802b906a1859"
	ErrComp107Fail    = "E_COMP_107_FAIL"
	OpHandleComp107   = "HandleComp107"
)

type Input107 struct {
	ID string
}

type Output107 struct {
	Result bool
}

func HandleComp107(in Input107) (Output107, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output107{Result: false}, apperror.New(OpHandleComp107, ErrComp107Fail, map[string]any{"id": in.ID})
	}

	return Output107{Result: true}, nil
}
