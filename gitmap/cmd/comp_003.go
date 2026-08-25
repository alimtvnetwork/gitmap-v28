package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const comp003UniquenessToken = "e7f6c011776e"

// Input003 defines the input contract for component 003.
type Input003 struct {
	ID string
}

// Output003 defines the output contract for component 003.
type Output003 struct {
	Result bool
}

// HandleComp003 executes unit component 003.
func HandleComp003(in Input003) (Output003, error) {
	if in.ID == "" {
		return Output003{Result: false}, apperror.New("HandleComp003", "E_COMP_003_FAIL", map[string]any{
			"id":    in.ID,
			"token": comp003UniquenessToken,
		})
	}

	return Output003{Result: true}, nil
}
