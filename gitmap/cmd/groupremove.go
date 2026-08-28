package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// runGroupRemove handles "group remove <group> <slug...>".
func runGroupRemove(args []string) error {
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
		removeOneSlugFromGroup(db, groupName, slug)
	}
	return nil
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
