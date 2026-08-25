package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp036BoundID    = "76a50887d8f1"
	comp036Uniqueness = "872261620421"
	comp036ErrFail    = "E_COMP_036_FAIL"
	comp036OpHandle   = "HandleComp036"
)

// Input036 represents the input contract for component 036.
type Input036 struct {
	ID string
}

// Output036 represents the output contract for component 036.
type Output036 struct {
	Result bool
}

// HandleComp036 handles component 036 execution.
func HandleComp036(in Input036) (Output036, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output036{Result: false}, apperror.New(comp036OpHandle, comp036ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp036BoundID,
			"uniqueness": comp036Uniqueness,
		})
	}

	return Output036{Result: true}, nil
}
