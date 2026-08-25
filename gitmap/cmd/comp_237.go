package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp237ID         = "f0bc318fb896"
	Comp237Uniqueness = "98144d79af44"
	ErrComp237Fail    = "E_COMP_237_FAIL"
	OpHandleComp237   = "HandleComp237"
)

type Input237 struct {
	ID string
}

type Output237 struct {
	Result bool
}

func HandleComp237(in Input237) (Output237, error) {
	if in.ID == Comp237Uniqueness {
		return Output237{Result: true}, nil
	}
	return Output237{Result: false}, apperror.New(OpHandleComp237, ErrComp237Fail, map[string]any{"id": in.ID})
}
