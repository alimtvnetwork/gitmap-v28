package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp014BoundID    = "8527a891e224"
	comp014Uniqueness = "59e19706d51d"
	comp014ErrFail    = "E_COMP_014_FAIL"
	comp014OpHandle   = "HandleComp014"
)

// Input014 represents the input contract for component 014.
type Input014 struct {
	ID string
}

// Output014 represents the output contract for component 014.
type Output014 struct {
	Result bool
}

// HandleComp014 handles component 014 execution.
func HandleComp014(in Input014) (Output014, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output014{Result: false}, apperror.New(comp014OpHandle, comp014ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp014BoundID,
			"uniqueness": comp014Uniqueness,
		})
	}

	return Output014{Result: true}, nil
}
