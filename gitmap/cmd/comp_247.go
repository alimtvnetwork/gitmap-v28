package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	Comp247ID         = "396f80444382"
	Comp247Uniqueness = "d18b29d80a8b"
	ErrComp247Fail    = "E_COMP_247_FAIL"
	OpHandleComp247   = "HandleComp247"
)

// Input247 represents the input contract for component 247.
type Input247 struct {
	ID string
}

// Output247 represents the output contract for component 247.
type Output247 struct {
	Result bool
}

// HandleComp247 handles component 247 execution.
func HandleComp247(in Input247) (Output247, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output247{Result: false}, apperror.New(OpHandleComp247, ErrComp247Fail, map[string]any{
			"id":         in.ID,
			"bound_id":   Comp247ID,
			"uniqueness": Comp247Uniqueness,
		})
	}

	return Output247{Result: true}, nil
}
