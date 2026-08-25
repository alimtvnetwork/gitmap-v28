package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp040ID         = "d59eced1ded0"
	Comp040Uniqueness = "48449a14a4ff"
	ErrComp040Fail    = "E_COMP_040_FAIL"
	OpHandleComp040   = "HandleComp040"
)

type Input040 struct {
	ID string
}

type Output040 struct {
	Result bool
}

func HandleComp040(in Input040) (Output040, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output040{Result: false}, apperror.New(OpHandleComp040, ErrComp040Fail, map[string]any{"id": in.ID})
	}

	return Output040{Result: true}, nil
}
