package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp108ID         = "9537f32ec759"
	Comp108Uniqueness = "0f4121d0ef1d"
	ErrComp108Fail    = "E_COMP_108_FAIL"
	OpHandleComp108   = "HandleComp108"
)

type Input108 struct {
	ID string
}

type Output108 struct {
	Result bool
}

func HandleComp108(in Input108) (Output108, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output108{Result: false}, apperror.New(OpHandleComp108, ErrComp108Fail, map[string]any{"id": in.ID})
	}

	return Output108{Result: true}, nil
}
