// Package cmd implements CLI command handlers for gitmap.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runTempReleaseList lists all temp-release branches.
func runTempReleaseList(args []string) {
	jsonOutput := hasTRListFlag(args, "--json")

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ DB migration failed: %v\n", err)
	}

	releases, err := db.ListTempReleases()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}

	if jsonOutput {
		if err := encodeTempReleaseListJSON(os.Stdout, releases); err != nil {
			fmt.Fprintf(os.Stderr, "  ✗ Failed to encode temp releases to JSON: %v\n", err)
		}

		return
	}

	printTRList(releases)
}

// printTRList prints temp-release records in terminal format.
func printTRList(releases []model.TempRelease) {
	if len(releases) == 0 {
		fmt.Print(constants.MsgTRListEmpty)

		return
	}

	fmt.Printf(constants.MsgTRListHeader, len(releases))

	for _, r := range releases {
		short := r.CommitSha
		if len(short) > constants.ShaDisplayLength {
			short = short[:constants.ShaDisplayLength]
		}

		msg := r.CommitMessage
		if len(msg) > 50 {
			msg = msg[:50]
		}

		fmt.Printf(constants.MsgTRListRow, r.Branch, short, msg, r.CreatedAt)
	}
}

// hasTRListFlag checks if a flag is present in the args.
func hasTRListFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}

	return false
}
