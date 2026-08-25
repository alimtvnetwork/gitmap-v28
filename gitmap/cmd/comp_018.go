package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp018BoundID    = "4ec9599fc203"
	comp018Uniqueness = "76a50887d8f1"
	comp018ErrFail    = "E_COMP_018_FAIL"
	comp018OpHandle   = "HandleComp018"
)

// Input018 represents the input contract for component 018.
type Input018 struct {
	ID string
}

// Output018 represents the output contract for component 018.
type Output018 struct {
	Result bool
}

// HandleComp018 handles component 018 execution.
func HandleComp018(in Input018) (Output018, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output018{Result: false}, apperror.New(comp018OpHandle, comp018ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp018BoundID,
			"uniqueness": comp018Uniqueness,
		})
	}

	return Output018{Result: true}, nil
}
