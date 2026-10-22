// Package cmdtemprelease implements CLI command handlers for gitmap.
package cmdtemprelease

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// runTempReleaseList lists all temp-release branches.
func runTempReleaseList(args []string) error {
	isJSONOutput := hasTRListFlag(args, "--json")

	db, err := openDB()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrListDBFailed)
	}

	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ DB migration failed: %v\n", err)
	}

	releases, err := db.ListTempReleases()
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBareFmt)
	}

	return renderTRListOutput(releases, isJSONOutput)
}

func renderTRListOutput(releases []model.TempRelease, isJSONOutput bool) error {
	if !isJSONOutput {
		printTRList(releases)

		return nil
	}

	if errEncode := encodeTempReleaseListJSON(os.Stdout, releases); errEncode != nil {
		return apperror.WrapSimple(errEncode, "encode temp release json")
	}

	return nil
}

// printTRList prints temp-release records in terminal format.
func printTRList(releases []model.TempRelease) {
	if len(releases) == 0 {
		fmt.Print(constants.MsgTRListEmpty)

		return
	}

	fmt.Printf(constants.MsgTRListHeader, len(releases))

	for _, r := range releases {
		printTRListRow(r)
	}
}

func printTRListRow(r model.TempRelease) {
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

// hasTRListFlag checks if a flag is present in the args.
func hasTRListFlag(args []string, flag string) bool {
	for _, a := range args {
		if a == flag {
			return true
		}
	}

	return false
}
