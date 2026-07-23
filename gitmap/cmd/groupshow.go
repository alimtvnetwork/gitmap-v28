package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runGroupShow handles "group show <name>".
func runGroupShow(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupNameReq)
		os.Exit(1)
	}
	name := args[0]
	executeGroupShow(name)
}

// executeGroupShow opens the DB and displays group repos.
func executeGroupShow(name string) {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	repos, err := db.ShowGroup(name)
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
	printGroupShowOutput(name, repos)
}

// printGroupShowOutput renders repos in a group with header and rows.
func printGroupShowOutput(name string, repos []model.ScanRecord) {
	fmt.Printf(constants.MsgGroupShowHeader, name, len(repos))
	fmt.Println(constants.MsgListSeparator)
	for _, r := range repos {
		fmt.Printf(constants.MsgGroupShowRowFmt, r.Slug, r.AbsolutePath)
	}
	fmt.Println()
}
