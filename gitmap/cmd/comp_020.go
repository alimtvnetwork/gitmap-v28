package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp020BoundID    = "f5ca38f748a1"
	comp020Uniqueness = "d59eced1ded0"
	comp020ErrFail    = "E_COMP_020_FAIL"
	comp020OpHandle   = "HandleComp020"
)

// Input020 represents the input contract for component 020.
type Input020 struct {
	ID string
}

// Output020 represents the output contract for component 020.
type Output020 struct {
	Result bool
}

// HandleComp020 handles component 020 execution.
func HandleComp020(in Input020) (Output020, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output020{Result: false}, apperror.New(comp020OpHandle, comp020ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp020BoundID,
			"uniqueness": comp020Uniqueness,
		})
	}

	return Output020{Result: true}, nil
}
