package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp271ID         = "3635a91e3da8"
	Comp271Uniqueness = "2d86377d4cc3"
	ErrComp271Fail    = "E_COMP_271_FAIL"
	OpHandleComp271   = "HandleComp271"
)

type Input271 struct {
	ID string
}

type Output271 struct {
	Result bool
}

func HandleComp271(in Input271) (Output271, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output271{Result: false}, apperror.New(OpHandleComp271, ErrComp271Fail, map[string]any{"id": in.ID})
	}

	return Output271{Result: true}, nil
}
