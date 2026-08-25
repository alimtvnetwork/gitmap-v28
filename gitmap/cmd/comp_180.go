package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp180ID         = "7b69759630f8"
	Comp180Uniqueness = "838f461c2fa6"
	ErrComp180Fail    = "E_COMP_180_FAIL"
	OpHandleComp180   = "HandleComp180"
)

type Input180 struct {
	ID string
}

type Output180 struct {
	Result bool
}

func HandleComp180(in Input180) (Output180, error) {
	if in.ID == Comp180Uniqueness {
		return Output180{Result: true}, nil
	}
	return Output180{Result: false}, apperror.New(OpHandleComp180, ErrComp180Fail, map[string]any{"id": in.ID})
}
