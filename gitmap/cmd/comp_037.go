package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp037ID         = "7a61b53701be"
	Comp037Uniqueness = "eb624dbe56eb"
	ErrComp037Fail    = "E_COMP_037_FAIL"
	OpHandleComp037   = "HandleComp037"
)

type Input037 struct {
	ID string
}

type Output037 struct {
	Result bool
}

func HandleComp037(in Input037) (Output037, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output037{Result: false}, apperror.New(OpHandleComp037, ErrComp037Fail, map[string]any{"id": in.ID})
	}

	return Output037{Result: true}, nil
}
