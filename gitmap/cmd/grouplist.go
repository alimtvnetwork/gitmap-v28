package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runGroupList handles "group list".
func runGroupList() error {
	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}
	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil && isLegacyDataError(err) {
		fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
		return apperror.NewSimple("fatal error", "E9000")
	}
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	printGroupList(db, groups)
	printHints(groupListHints())
	return nil
}

// printGroupList renders the group table to stdout.
func printGroupList(db *store.DB, groups []model.Group) {
	if len(groups) == 0 {
		fmt.Println(constants.MsgGroupEmpty)

		return
	}
	fmt.Println(constants.MsgGroupHeader)
	fmt.Println(constants.MsgListSeparator)
	for _, g := range groups {
		count, countErr := db.CountGroupRepos(g.Name)
		if countErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not count repos for group %s: %v\n", g.Name, countErr)
		}
		fmt.Printf(constants.MsgGroupRowFmt, g.Name, count, g.Description)
	}
}
