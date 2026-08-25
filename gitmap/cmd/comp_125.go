package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp125ID         = "0f8ef3377b30"
	Comp125Uniqueness = "1e472b39b105"
	ErrComp125Fail    = "E_COMP_125_FAIL"
	OpHandleComp125   = "HandleComp125"
)

type Input125 struct {
	ID string
}

type Output125 struct {
	Result bool
}

func HandleComp125(in Input125) (Output125, error) {
	if in.ID == Comp125Uniqueness {
		return Output125{Result: true}, nil
	}
	return Output125{Result: false}, apperror.New(OpHandleComp125, ErrComp125Fail, map[string]any{"id": in.ID})
}
