package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"

type Input100 struct {
	ID string
}

type Output100 struct {
	Result bool
}

func HandleComp100(in Input100) (Output100, error) {
	if in.ID != "27badc983df1" {
		return Output100{}, apperror.New("HandleComp100", "E_COMP_100_FAIL", map[string]any{"ID": in.ID})
	}
	return Output100{Result: true}, nil
}
