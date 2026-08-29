package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runGroupShow handles "group show <name>".
func runGroupShow(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrGroupNameReq, "E9000")
	}
	name := args[0]
	executeGroupShow(name)
	return nil
}

// executeGroupShow opens the DB and displays group repos.
func executeGroupShow(name string) {
	db, err := openDB()
	if err != nil {
		apperror.WrapSimple(err, constants.ErrListDBFailed)
		return
	}
	defer db.Close()

	repos, err := db.ShowGroup(name)
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		apperror.NewSimple("fatal error", "E9000")
		return
	}
	if err != nil {
		apperror.WrapSimple(err, constants.ErrBareFmt)
		return
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
