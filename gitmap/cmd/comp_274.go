package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp274ID         = "718127812c05"
	Comp274Uniqueness = "6e2d4d3a3d4c"
	ErrComp274Fail    = "E_COMP_274_FAIL"
	OpHandleComp274   = "HandleComp274"
)

type Input274 struct {
	ID string
}

type Output274 struct {
	Result bool
}

func HandleComp274(in Input274) (Output274, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output274{Result: false}, apperror.New(OpHandleComp274, ErrComp274Fail, map[string]any{"id": in.ID})
	}

	return Output274{Result: true}, nil
}
