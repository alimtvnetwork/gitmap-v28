package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp203ID         = "4621c1d55fa4"
	Comp203Uniqueness = "f64f410744d9"
	ErrComp203Fail    = "E_COMP_203_FAIL"
	OpHandleComp203   = "HandleComp203"
)

type Input203 struct {
	ID string
}

type Output203 struct {
	Result bool
}

func HandleComp203(in Input203) (Output203, error) {
	if in.ID == Comp203Uniqueness {
		return Output203{Result: true}, nil
	}
	return Output203{Result: false}, apperror.New(OpHandleComp203, ErrComp203Fail, map[string]any{"id": in.ID})
}
