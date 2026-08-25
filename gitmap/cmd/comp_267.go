package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp267ID         = "8acc23987b89"
	Comp267Uniqueness = "5ef6514ed330"
	ErrComp267Fail    = "E_COMP_267_FAIL"
	OpHandleComp267   = "HandleComp267"
)

type Input267 struct {
	ID string
}

type Output267 struct {
	Result bool
}

func HandleComp267(in Input267) (Output267, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output267{Result: false}, apperror.New(OpHandleComp267, ErrComp267Fail, map[string]any{"id": in.ID})
	}

	return Output267{Result: true}, nil
}
