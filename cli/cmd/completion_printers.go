package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/completion"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func withOpenDB(fn func(db *store.DB)) {
	db, err := openDB()
	if err != nil {
		return
	}

	defer db.Close()
	fn(db)
}

func printCompletionRepos() {
	withOpenDB(func(db *store.DB) {
		repos, err := db.ListRepos()
		if err != nil {
			return
		}

		for _, r := range repos {
			fmt.Println(r.Slug)
		}
	})
}

func printCompletionGroups() {
	withOpenDB(func(db *store.DB) {
		groups, err := db.ListGroups()
		if err != nil {
			return
		}

		for _, g := range groups {
			fmt.Println(g.Name)
		}
	})
}

func printCompletionCommands() {
	for _, cmd := range completion.AllCommands() {
		fmt.Println(cmd)
	}
}

func printCompletionAliases() {
	withOpenDB(func(db *store.DB) {
		aliases, err := db.ListAliases()
		if err != nil {
			return
		}

		for _, a := range aliases {
			fmt.Println(a.Alias)
		}
	})
}
