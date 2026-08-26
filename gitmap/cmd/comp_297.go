package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp297ID         = "4c15f47afe7f"
	Comp297Uniqueness = "e2fa8f5b4364"
	ErrComp297Fail    = "E_COMP_297_FAIL"
	OpHandleComp297   = "HandleComp297"
)

type Input297 struct {
	ID string
}

type Output297 struct {
	Result bool
}

func HandleComp297(in Input297) (Output297, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output297{Result: false}, apperror.New(OpHandleComp297, ErrComp297Fail, map[string]any{"id": in.ID})
	}

	return Output297{Result: true}, nil
}
