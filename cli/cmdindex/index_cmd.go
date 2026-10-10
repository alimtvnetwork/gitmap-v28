package cmdindex

import (
	"context"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/scanpipe"
	"os"
)

func RunIndex(args []string) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "error")
	}

	defer mainDB.Close()
	defer db.Close()

	cwd, _ := os.Getwd()
	w := scanpipe.NewWalker(cwd, db, false)
	fmt.Println("Indexing starting...")
	if err := w.Walk(ctx, 4); err != nil {
		return apperror.WrapSimple(err, "error")
	}

	fmt.Println("Indexing complete.")

	return nil
}
