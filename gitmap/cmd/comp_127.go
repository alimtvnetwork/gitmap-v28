package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp127ID         = "922c7954216c"
	Comp127Uniqueness = "9512d95d00d6"
	ErrComp127Fail    = "E_COMP_127_FAIL"
	OpHandleComp127   = "HandleComp127"
)

type Input127 struct {
	ID string
}

type Output127 struct {
	Result bool
}

func HandleComp127(in Input127) (Output127, error) {
	if in.ID == Comp127Uniqueness {
		return Output127{Result: true}, nil
	}
	return Output127{Result: false}, apperror.New(OpHandleComp127, ErrComp127Fail, map[string]any{"id": in.ID})
}
