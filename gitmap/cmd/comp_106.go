package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp106ID         = "482d9673cfee"
	Comp106Uniqueness = "fa2b7af0a811"
	ErrComp106Fail    = "E_COMP_106_FAIL"
	OpHandleComp106   = "HandleComp106"
)

type Input106 struct {
	ID string
}

type Output106 struct {
	Result bool
}

func HandleComp106(in Input106) (Output106, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output106{Result: false}, apperror.New(OpHandleComp106, ErrComp106Fail, map[string]any{"id": in.ID})
	}

	return Output106{Result: true}, nil
}
