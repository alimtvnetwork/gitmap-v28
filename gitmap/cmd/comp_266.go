package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp266ID         = "ea5b27556fbb"
	Comp266Uniqueness = "68f10bf021d7"
	ErrComp266Fail    = "E_COMP_266_FAIL"
	OpHandleComp266   = "HandleComp266"
)

type Input266 struct {
	ID string
}

type Output266 struct {
	Result bool
}

func HandleComp266(in Input266) (Output266, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output266{Result: false}, apperror.New(OpHandleComp266, ErrComp266Fail, map[string]any{"id": in.ID})
	}

	return Output266{Result: true}, nil
}
