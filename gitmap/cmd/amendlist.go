// Package cmd — amendlist.go handles the amend-list command.
package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runAmendList handles the "amend-list" command.
func runAmendList(args []string) {
	checkHelp("amend-list", args)
	asJSON := hasAmendListJSONFlag(args)
	limit := parseAmendListLimit(args)
	branch := parseAmendListBranch(args)

	amendments := loadAmendments(branch)
	amendments = applyAmendmentLimit(amendments, limit)

	if asJSON {
		printAmendmentsJSON(amendments)

		return
	}

	printAmendmentsTerminal(amendments)
}

// hasAmendListJSONFlag checks if --json is present in args.
func hasAmendListJSONFlag(args []string) bool {
	for _, arg := range args {
		if arg == constants.FlagJSON {
			return true
		}
	}

	return false
}

// parseAmendListLimit extracts the --limit N value from args.
func parseAmendListLimit(args []string) int {
	for i, arg := range args {
		if arg == constants.FlagLimit && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil && n > 0 {
				return n
			}
		}
	}

	return 0
}

// parseAmendListBranch extracts the --branch value from args.
func parseAmendListBranch(args []string) string {
	for i, arg := range args {
		if arg == constants.FlagAmendListBranch && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

// loadAmendments opens the DB and fetches amendments, optionally filtered by branch.
func loadAmendments(branch string) []store.AmendmentRow {
	db, err := openDB()
	if err != nil {
		fmt.Fprintln(os.Stderr, constants.ErrNoDatabase)
		os.Exit(1)
	}
	defer db.Close()

	var amendments []store.AmendmentRow

	if branch != "" {
		amendments, err = db.ListAmendmentsByBranch(branch)
	} else {
		amendments, err = db.ListAmendments()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendListFailed, err)
		os.Exit(1)
	}

	return amendments
}

// applyAmendmentLimit trims amendments to at most n items (0 means no limit).
func applyAmendmentLimit(amendments []store.AmendmentRow, n int) []store.AmendmentRow {
	if n <= 0 || n >= len(amendments) {
		return amendments
	}

	return amendments[:n]
}

// printAmendmentsTerminal renders amendments as a table to stdout.
func printAmendmentsTerminal(amendments []store.AmendmentRow) {
	if len(amendments) == 0 {
		fmt.Println(constants.MsgAmendListEmpty)

		return
	}

	fmt.Printf(constants.MsgAmendListHeader, len(amendments))
	fmt.Println(constants.MsgAmendListSeparator)
	fmt.Println(constants.MsgAmendListColumns)

	for _, a := range amendments {
		printAmendmentRow(a)
	}
}

// printAmendmentRow prints a single amendment row.
func printAmendmentRow(a store.AmendmentRow) {
	forcePushed := constants.MsgNo
	if a.ForcePushed == 1 {
		forcePushed = constants.MsgYes
	}

	fmt.Printf(constants.MsgAmendListRowFmt,
		a.Branch, a.Mode, a.TotalCommits,
		a.PreviousName, a.PreviousEmail,
		a.NewName, a.NewEmail,
		forcePushed, a.CreatedAt)
}

// printAmendmentsJSON renders amendments as stable JSON to stdout.
func printAmendmentsJSON(amendments []store.AmendmentRow) {
	if err := encodeAmendListJSON(os.Stdout, amendments); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode amendments to JSON: %v\n", err)
	}
}
