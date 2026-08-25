package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp256ID         = "51e8ea280b44"
	Comp256Uniqueness = "94f8607915df"
	ErrComp256Fail    = "E_COMP_256_FAIL"
	OpHandleComp256   = "HandleComp256"
)

type Input256 struct {
	ID string
}

type Output256 struct {
	Result bool
}

func HandleComp256(in Input256) (Output256, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output256{Result: false}, apperror.New(OpHandleComp256, ErrComp256Fail, map[string]any{"id": in.ID})
	}

	return Output256{Result: true}, nil
}
