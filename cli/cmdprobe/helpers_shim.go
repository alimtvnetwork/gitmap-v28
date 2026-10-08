package cmdprobe

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"fmt"
	"os"
)

// openSfDB opens the default profile DB and runs migrations. Mirrors
// the upsertToDB / scan helpers so `gitmap sf` shares the exact same
// resolution rules used by `gitmap scan`.
func openSfDB() *store.DB {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		cliexit.HandleError(err, 1)
	}

	if err := db.Migrate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		db.Close()
		cliexit.HandleError(err, 1)
	}

	return db
}
