package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp200ID         = "27badc983df1"
	Comp200Uniqueness = "26d228663f13"
	ErrComp200Fail    = "E_COMP_200_FAIL"
	OpHandleComp200   = "HandleComp200"
)

type Input200 struct {
	ID string
}

type Output200 struct {
	Result bool
}

func HandleComp200(in Input200) (Output200, error) {
	if in.ID == Comp200Uniqueness {
		return Output200{Result: true}, nil
	}
	return Output200{Result: false}, apperror.New(OpHandleComp200, ErrComp200Fail, map[string]any{"id": in.ID})
}
