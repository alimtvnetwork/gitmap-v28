package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp122ID         = "1be00341082e"
	Comp122Uniqueness = "82c01ce15b43"
	ErrComp122Fail    = "E_COMP_122_FAIL"
	OpHandleComp122   = "HandleComp122"
)

type Input122 struct {
	ID string
}

type Output122 struct {
	Result bool
}

func HandleComp122(in Input122) (Output122, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output122{Result: false}, apperror.New(OpHandleComp122, ErrComp122Fail, map[string]any{"id": in.ID})
	}

	return Output122{Result: true}, nil
}
