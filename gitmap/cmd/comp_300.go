package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp300ID         = "983bd614bb5a"
	Comp300Uniqueness = "284b7e6d788f"
	ErrComp300Fail    = "E_COMP_300_FAIL"
	OpHandleComp300   = "HandleComp300"
)

type Input300 struct {
	ID string
}

type Output300 struct {
	Result bool
}

func HandleComp300(in Input300) (Output300, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output300{Result: false}, apperror.New(OpHandleComp300, ErrComp300Fail, map[string]any{"id": in.ID})
	}

	return Output300{Result: true}, nil
}
