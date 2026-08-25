package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp124ID         = "6affdae3b3c1"
	Comp124Uniqueness = "766cb53c753b"
	ErrComp124Fail    = "E_COMP_124_FAIL"
	OpHandleComp124   = "HandleComp124"
)

type Input124 struct {
	ID string
}

type Output124 struct {
	Result bool
}

func HandleComp124(in Input124) (Output124, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output124{Result: false}, apperror.New(OpHandleComp124, ErrComp124Fail, map[string]any{"id": in.ID})
	}

	return Output124{Result: true}, nil
}
