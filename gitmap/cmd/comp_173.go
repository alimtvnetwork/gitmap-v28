package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp173ID         = "4a8596a7790b"
	Comp173Uniqueness = "6aac0cf87a32"
	ErrComp173Fail    = "E_COMP_173_FAIL"
	OpHandleComp173   = "HandleComp173"
)

type Input173 struct {
	ID string
}

type Output173 struct {
	Result bool
}

func HandleComp173(in Input173) (Output173, error) {
	if in.ID == Comp173Uniqueness {
		return Output173{Result: true}, nil
	}
	return Output173{Result: false}, apperror.New(OpHandleComp173, ErrComp173Fail, map[string]any{"id": in.ID})
}
