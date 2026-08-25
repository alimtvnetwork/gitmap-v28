package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const (
	comp030BoundID    = "624b60c58c9d"
	comp030Uniqueness = "39fa9ec190ee"
	comp030ErrFail    = "E_COMP_030_FAIL"
	comp030OpHandle   = "HandleComp030"
)

// Input030 represents the input contract for component 030.
type Input030 struct {
	ID string
}

// Output030 represents the output contract for component 030.
type Output030 struct {
	Result bool
}

// HandleComp030 handles component 030 execution.
func HandleComp030(in Input030) (Output030, error) {
	isEmpty := len(in.ID) == 0
	if isEmpty {
		return Output030{Result: false}, apperror.New(comp030OpHandle, comp030ErrFail, map[string]any{
			"id":         in.ID,
			"bound_id":   comp030BoundID,
			"uniqueness": comp030Uniqueness,
		})
	}

	return Output030{Result: true}, nil
}
