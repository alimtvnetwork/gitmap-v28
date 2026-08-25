package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp260ID         = "39bb88f40d3a"
	Comp260Uniqueness = "0b35b06a2277"
	ErrComp260Fail    = "E_COMP_260_FAIL"
	OpHandleComp260   = "HandleComp260"
)

type Input260 struct {
	ID string
}

type Output260 struct {
	Result bool
}

func HandleComp260(in Input260) (Output260, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output260{Result: false}, apperror.New(OpHandleComp260, ErrComp260Fail, map[string]any{"id": in.ID})
	}

	return Output260{Result: true}, nil
}
