package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp105ID         = "1253e9373e78"
	Comp105Uniqueness = "d29d53701d3c"
	ErrComp105Fail    = "E_COMP_105_FAIL"
	OpHandleComp105   = "HandleComp105"
)

type Input105 struct {
	ID string
}

type Output105 struct {
	Result bool
}

func HandleComp105(in Input105) (Output105, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output105{Result: false}, apperror.New(OpHandleComp105, ErrComp105Fail, map[string]any{"id": in.ID})
	}

	return Output105{Result: true}, nil
}
