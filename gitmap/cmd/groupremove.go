package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runGroupRemove handles "group remove <group> <slug...>".
func runGroupRemove(args []string) {
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
		removeOneSlugFromGroup(db, groupName, slug)
	}
}

// removeOneSlugFromGroup resolves a slug and removes matching repos.
func removeOneSlugFromGroup(db *store.DB, groupName, slug string) {
	repos, err := db.FindBySlug(slug)
	if err != nil || len(repos) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrDBNoMatch, slug)

		return
	}
	for _, r := range repos {
		err := db.RemoveRepoFromGroup(groupName, r.ID)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)

			return
		}
		fmt.Printf(constants.MsgGroupRemoved, r.Slug, groupName)
	}
}
