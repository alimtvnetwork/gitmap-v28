package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runSSHList displays all stored SSH keys as an aligned table or JSON.
func runSSHList(args ...string) {
	jsonOut := hasFlagInArgs(args, constants.FlagSSHJSON)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	keys, err := db.ListSSHKeys()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}

	if jsonOut {
		printSSHListJSON(keys)

		return
	}

	if len(keys) == 0 {
		fmt.Println("  No SSH keys stored. Run 'gitmap ssh' to generate one.")

		return
	}

	fmt.Fprintf(os.Stdout, constants.MsgSSHListHeader, len(keys))
	fmt.Fprintf(os.Stdout, constants.MsgSSHListColumns, "Name", "Path", "Fingerprint", "Created")
	fmt.Fprintf(os.Stdout, constants.MsgSSHListColumns,
		"───────────────", "──────────────────────────────",
		"─────────────────────────", "──────────")

	for _, k := range keys {
		created := k.CreatedAt
		if len(created) > 10 {
			created = created[:10]
		}

		fmt.Fprintf(os.Stdout, constants.MsgSSHListRow, k.Name, k.PrivatePath, k.Fingerprint, created)
	}
}

// printSSHListJSON outputs SSH keys as JSON.
func printSSHListJSON(keys []model.SSHKey) {
	if err := encodeSSHListJSON(os.Stdout, keys); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
	}
}

// hasFlagInArgs checks if a flag is present in the given args slice.
func hasFlagInArgs(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}

	return false
}
