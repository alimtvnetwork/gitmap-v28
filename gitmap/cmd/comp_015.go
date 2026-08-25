package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp015BoundID    = "e629fa6598d7"
	comp015Uniqueness = "624b60c58c9d"
	comp015ErrFail    = "E_COMP_015_FAIL"
	comp015OpHandle   = "HandleComp015"
)

// Input015 represents the input contract for component 015.
type Input015 struct {
	ID string
}

// Output015 represents the output contract for component 015.
type Output015 struct {
	Result bool
}

// HandleComp015 handles component 015 execution.
func HandleComp015(in Input015) (Output015, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output015{Result: false}, apperror.New(comp015OpHandle, comp015ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp015BoundID,
			"uniqueness": comp015Uniqueness,
		})
	}

	return Output015{Result: true}, nil
}
