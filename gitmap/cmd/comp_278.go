package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp278ID         = "ee62de25ccc2"
	Comp278Uniqueness = "9d6aa3d89c01"
	ErrComp278Fail    = "E_COMP_278_FAIL"
	OpHandleComp278   = "HandleComp278"
)

type Input278 struct {
	ID string
}

type Output278 struct {
	Result bool
}

func HandleComp278(in Input278) (Output278, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output278{Result: false}, apperror.New(OpHandleComp278, ErrComp278Fail, map[string]any{"id": in.ID})
	}

	return Output278{Result: true}, nil
}
