package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp010BoundID    = "4a44dc153642"
	comp010Uniqueness = "f5ca38f748a1"
	comp010ErrFail    = "E_COMP_010_FAIL"
	comp010OpHandle   = "HandleComp010"
)

// Input010 represents the input contract for component 010.
type Input010 struct {
	ID string
}

// Output010 represents the output contract for component 010.
type Output010 struct {
	Result bool
}

// HandleComp010 handles component 010 execution.
func HandleComp010(in Input010) (Output010, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output010{Result: false}, apperror.New(comp010OpHandle, comp010ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp010BoundID,
			"uniqueness": comp010Uniqueness,
		})
	}

	return Output010{Result: true}, nil
}
