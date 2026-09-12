package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// runGroupAdd handles "group add <group> <slug...>".
func runGroupAdd(args []string) error {
	if len(args) < 2 {
		return apperror.NewSimple(constants.ErrGroupSlugReq, "E9000")
	}

	groupName := args[0]
	slugs := args[1:]

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()

	for _, slug := range slugs {
		addOneSlugToGroup(db, groupName, slug)
	}

	return nil
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
