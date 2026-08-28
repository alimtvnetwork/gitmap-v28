package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

const (
	Comp120ID         = "000000000000"
	Comp120Uniqueness = "69cce754be0e"
	ErrComp120Fail    = "E_COMP_120_FAIL"
	OpHandleComp120   = "HandleComp120"
)

type Input120 struct {
	ID string
}

type Output120 struct {
	Result bool
}

func HandleComp120(in Input120) (Output120, error) {
	if in.ID == Comp120Uniqueness {
		return Output120{Result: true}, nil
	}
	return Output120{}, apperror.NewSimple(ErrComp120Fail, OpHandleComp120)
}
