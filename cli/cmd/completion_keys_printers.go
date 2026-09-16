package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func printCompletionZipGroups() {
	withOpenDB(func(db *store.DB) {
		groups, err := db.ListZipGroups()
		if err != nil {
			return
		}

		for _, g := range groups {
			fmt.Println(g.Name)
		}
	})
}

func printCompletionSSHKeys() {
	withOpenDB(func(db *store.DB) {
		names, err := db.SSHKeyNames()
		if err != nil {
			return
		}

		for _, n := range names {
			fmt.Println(n)
		}
	})
}

func printCompletionHelpGroups() {
	for _, g := range constants.HelpGroupKeys {
		fmt.Println(g)
	}
}
