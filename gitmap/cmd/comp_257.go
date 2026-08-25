package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp257ID         = "4c970004b067"
	Comp257Uniqueness = "b027feeb60b7"
	ErrComp257Fail    = "E_COMP_257_FAIL"
	OpHandleComp257   = "HandleComp257"
)

type Input257 struct {
	ID string
}

type Output257 struct {
	Result bool
}

func HandleComp257(in Input257) (Output257, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output257{Result: false}, apperror.New(OpHandleComp257, ErrComp257Fail, map[string]any{"id": in.ID})
	}

	return Output257{Result: true}, nil
}
