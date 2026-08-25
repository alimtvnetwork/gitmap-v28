package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp004BoundID         = "4b227777d4dd"
	comp004UniquenessToken = "2c624232cdd2"
	comp004ErrFail         = "E_COMP_004_FAIL"
)

// Input004 defines the input contract for component 004.
type Input004 struct {
	ID string
}

// Output004 defines the output contract for component 004.
type Output004 struct {
	Result bool
}

// HandleComp004 executes unit component 004.
func HandleComp004(in Input004) (Output004, error) {
	if in.ID == "" {
		return Output004{Result: false}, apperror.New("HandleComp004", comp004ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp004BoundID,
			"uniqueness": comp004UniquenessToken,
		})
	}

	return Output004{Result: true}, nil
}
