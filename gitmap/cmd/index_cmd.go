package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/indexer"
)

func runIndex(args []string) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		return apperror.WrapSimple(err, "error")
	}
	defer mainDB.Close()
	defer db.Close()

	cwd, _ := os.Getwd()
	w := indexer.NewWalker(cwd, db, false)
	fmt.Println("Indexing starting...")
	if err := w.Walk(ctx, 4); err != nil {
		return apperror.WrapSimple(err, "error")
	}
	fmt.Println("Indexing complete.")
	return nil
}
