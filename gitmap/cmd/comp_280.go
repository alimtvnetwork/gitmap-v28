package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp280ID         = "7f0a22117f8f"
	Comp280Uniqueness = "6bcaea988250"
	ErrComp280Fail    = "E_COMP_280_FAIL"
	OpHandleComp280   = "HandleComp280"
)

type Input280 struct {
	ID string
}

type Output280 struct {
	Result bool
}

func HandleComp280(in Input280) (Output280, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output280{Result: false}, apperror.New(OpHandleComp280, ErrComp280Fail, map[string]any{"id": in.ID})
	}

	return Output280{Result: true}, nil
}
