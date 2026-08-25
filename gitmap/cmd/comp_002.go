package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp002BoundID   = "d4735e3a265e"
	comp002DataToken = "4b227777d4dd"
	comp002ErrFail   = "E_COMP_002_FAIL"
)

// Input002 represents the input contract for component 002.
type Input002 struct {
	ID string
}

// Output002 represents the output contract for component 002.
type Output002 struct {
	Result bool
}

// HandleComp002 handles component 002 execution.
func HandleComp002(in Input002) (Output002, error) {
	if in.ID == "" {
		return Output002{Result: false}, apperror.New("HandleComp002", comp002ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp002BoundID,
			"data_token": comp002DataToken,
		})
	}

	return Output002{Result: true}, nil
}
