package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runGroupShow handles "group show <name>".
func runGroupShow(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple(constants.ErrGroupNameReq, "E9000")
	}

	name := args[0]
	if appErr := executeGroupShow(name); appErr != nil {
		return appErr
	}

	return nil
}

func fetchGroupRepos(db *store.DB, name string) ([]model.ScanRecord, *apperror.AppError) {
	repos, err := db.ShowGroup(name)
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)

		return nil, apperror.NewSimple("legacy data error", "E9000")
	}

	if err != nil {
		return nil, apperror.WrapSimple(err, constants.ErrBareFmt)
	}

	return repos, nil
}

// executeGroupShow opens the DB and displays group repos.
func executeGroupShow(name string) *apperror.AppError {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	repos, appErr := fetchGroupRepos(db, name)
	if appErr != nil {
		return appErr
	}

	printGroupShowOutput(name, repos)

	return nil
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
