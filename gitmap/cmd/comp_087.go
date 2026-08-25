package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp087ID         = "bdd2d3af3a5a"
	Comp087Uniqueness = "41e521adf8ae"
	ErrComp087Fail    = "E_COMP_087_FAIL"
	OpHandleComp087   = "HandleComp087"
)

type Input087 struct {
	ID string
}

type Output087 struct {
	Result bool
}

func HandleComp087(in Input087) (Output087, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output087{Result: false}, apperror.New(OpHandleComp087, ErrComp087Fail, map[string]any{"id": in.ID})
	}

	return Output087{Result: true}, nil
}
