package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp238ID         = "8ae4c23b80d1"
	Comp238Uniqueness = "e73cb135243c"
	ErrComp238Fail    = "E_COMP_238_FAIL"
	OpHandleComp238   = "HandleComp238"
)

type Input238 struct {
	ID string
}

type Output238 struct {
	Result bool
}

func HandleComp238(in Input238) (Output238, error) {
	if in.ID == Comp238Uniqueness {
		return Output238{Result: true}, nil
	}
	return Output238{Result: false}, apperror.New(OpHandleComp238, ErrComp238Fail, map[string]any{"id": in.ID})
}
