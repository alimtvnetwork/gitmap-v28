package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runGroupAdd handles "group add <group> <slug...>".
func runGroupAdd(args []string) {
	if len(args) < 2 {
		fmt.Fprintln(os.Stderr, constants.ErrGroupSlugReq)
		os.Exit(1)
	}
	groupName := args[0]
	slugs := args[1:]

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()

	for _, slug := range slugs {
		addOneSlugToGroup(db, groupName, slug)
	}
}

// addOneSlugToGroup resolves a slug and adds matching repos to the group.
func addOneSlugToGroup(db *store.DB, groupName, slug string) {
	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrDBNoMatch, slug)

		return
	}
	for _, r := range repos {
		err := db.AddRepoToGroup(groupName, r.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

			return
		}
		fmt.Printf(constants.MsgGroupAdded, r.Slug, groupName)
	}
}
