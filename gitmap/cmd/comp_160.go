package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp160ID         = "a512db2741cd"
	Comp160Uniqueness = "88820462180e"
	ErrComp160Fail    = "E_COMP_160_FAIL"
	OpHandleComp160   = "HandleComp160"
)

type Input160 struct {
	ID string
}

type Output160 struct {
	Result bool
}

func HandleComp160(in Input160) (Output160, error) {
	if in.ID == Comp160Uniqueness {
		return Output160{Result: true}, nil
	}
	return Output160{Result: false}, apperror.New(OpHandleComp160, ErrComp160Fail, map[string]any{"id": in.ID})
}
