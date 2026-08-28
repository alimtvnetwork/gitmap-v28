package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/indexer"
	"github.com/pterm/pterm"
)

func runIndex(args []string) error {
	ctx := context.Background()
	mainDB, db, err := getRepoDB(ctx)
	if err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	defer mainDB.Close()
	defer db.Close()

	cwd, _ := os.Getwd()
	w := indexer.NewWalker(cwd, db, false)
	fmt.Println("Indexing starting...")
	if err := w.Walk(ctx, 4); err != nil {
		pterm.Error.Println(err)
		os.Exit(1)
	}
	fmt.Println("Indexing complete.")
	return nil
}
