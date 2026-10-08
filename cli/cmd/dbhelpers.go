package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDB() (*store.DB, error) {
	return store.OpenDefault()
}

type statusSummary struct {
	Total   int
	Clean   int
	Dirty   int
	Ahead   int
	Behind  int
	Stashed int
	Missing int
}
