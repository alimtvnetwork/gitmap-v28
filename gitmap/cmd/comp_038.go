package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp038ID         = "aea92132c4cb"
	Comp038Uniqueness = "f74efabef12e"
	ErrComp038Fail    = "E_COMP_038_FAIL"
	OpHandleComp038   = "HandleComp038"
)

type Input038 struct {
	ID string
}

type Output038 struct {
	Result bool
}

func HandleComp038(in Input038) (Output038, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output038{Result: false}, apperror.New(OpHandleComp038, ErrComp038Fail, map[string]any{"id": in.ID})
	}

	return Output038{Result: true}, nil
}
