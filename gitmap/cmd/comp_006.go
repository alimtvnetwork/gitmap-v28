package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp006BoundID    = "e7f6c011776e"
	comp006Uniqueness = "6b51d431df5d"
	comp006ErrFail    = "E_COMP_006_FAIL"
	comp006OpHandle   = "HandleComp006"
)

// Input006 represents the input contract for component 006.
type Input006 struct {
	ID string
}

// Output006 represents the output contract for component 006.
type Output006 struct {
	Result bool
}

// HandleComp006 handles component 006 execution.
func HandleComp006(input Input006) (Output006, error) {
	isEmpty := len(input.ID) == 0
	if isEmpty {
		return Output006{Result: false}, apperror.New(comp006OpHandle, comp006ErrFail, map[string]any{
			"id":         input.ID,
			"bound_id":   comp006BoundID,
			"uniqueness": comp006Uniqueness,
		})
	}

	return Output006{Result: true}, nil
}
