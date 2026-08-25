package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp255ID         = "9556b82499cc"
	Comp255Uniqueness = "5e5c743a015f"
	ErrComp255Fail    = "E_COMP_255_FAIL"
	OpHandleComp255   = "HandleComp255"
)

type Input255 struct {
	ID string
}

type Output255 struct {
	Result bool
}

func HandleComp255(in Input255) (Output255, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output255{Result: false}, apperror.New(OpHandleComp255, ErrComp255Fail, map[string]any{"id": in.ID})
	}

	return Output255{Result: true}, nil
}
