package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp060ID         = "39fa9ec190ee"
	Comp060Uniqueness = "2abaca4911e6"
	ErrComp060Fail    = "E_COMP_060_FAIL"
	OpHandleComp060   = "HandleComp060"
)

type Input060 struct {
	ID string
}

type Output060 struct {
	Result bool
}

func HandleComp060(in Input060) (Output060, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output060{Result: false}, apperror.New(OpHandleComp060, ErrComp060Fail, map[string]any{"id": in.ID})
	}

	return Output060{Result: true}, nil
}
