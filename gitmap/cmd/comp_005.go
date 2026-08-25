package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp005BoundID         = "ef2d127de37b"
	comp005UniquenessToken = "4a44dc153642"
	comp005ErrFail         = "E_COMP_005_FAIL"
)

// Input005 defines the input contract for component 005.
type Input005 struct {
	ID string
}

// Output005 defines the output contract for component 005.
type Output005 struct {
	Result bool
}

// HandleComp005 executes unit component 005.
func HandleComp005(in Input005) (Output005, error) {
	if in.ID == "" {
		return Output005{Result: false}, apperror.New("HandleComp005", comp005ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp005BoundID,
			"uniqueness": comp005UniquenessToken,
		})
	}

	return Output005{Result: true}, nil
}
