package cmdimport

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func openDb() (*store.DB, error) {
	return store.OpenDefault()
}
