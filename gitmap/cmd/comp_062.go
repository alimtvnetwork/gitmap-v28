package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp062ID         = "81b8a03f97e8"
	Comp062Uniqueness = "6affdae3b3c1"
	ErrComp062Fail    = "E_COMP_062_FAIL"
	OpHandleComp062   = "HandleComp062"
)

type Input062 struct {
	ID string
}

type Output062 struct {
	Result bool
}

func HandleComp062(in Input062) (Output062, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output062{Result: false}, apperror.New(OpHandleComp062, ErrComp062Fail, map[string]any{"id": in.ID})
	}

	return Output062{Result: true}, nil
}
