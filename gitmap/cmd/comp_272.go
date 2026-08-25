package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp272ID         = "1c6c0bb2c7ec"
	Comp272Uniqueness = "d359f8b537f1"
	ErrComp272Fail    = "E_COMP_272_FAIL"
	OpHandleComp272   = "HandleComp272"
)

type Input272 struct {
	ID string
}

type Output272 struct {
	Result bool
}

func HandleComp272(in Input272) (Output272, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output272{Result: false}, apperror.New(OpHandleComp272, ErrComp272Fail, map[string]any{"id": in.ID})
	}

	return Output272{Result: true}, nil
}
