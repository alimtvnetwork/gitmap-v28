package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runGroupList handles "group list".
func runGroupList() {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	groups, err := db.ListGroups()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}

	printGroupList(db, groups)
	printHints(groupListHints())
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
