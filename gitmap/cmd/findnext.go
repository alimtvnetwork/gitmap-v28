package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// findNextUsageExitCode is the conventional CLI usage-error exit code.
// Distinct from exit-1 (used for I/O / DB failures) so scripts can
// branch on the cause: "you typed something wrong" vs "the system
// failed". Mirrors getopt(3) and most Unix tools.
const findNextUsageExitCode = 2

// runFindNext dispatches `gitmap find-next [--scan-folder <id>] [--json]`.
//
// As of v3.122.0 the parser rejects unknown flags, malformed values,
// and `--json=...` boolean misuse with a clear stderr message + the
// usage header, then exits 2. Previously these were silently ignored.
func runFindNext(args []string) {
	checkHelp("find-next", args)

	scanFolderID, jsonOut, err := parseFindNextFlags(args)
	if err != nil {
		fmt.Fprint(os.Stderr, err.Error())
		fmt.Fprintln(os.Stderr, constants.MsgFindNextUsageHeader)
		os.Exit(findNextUsageExitCode)
	}

	db := openSfDB()
	defer db.Close()

	rows, err := db.FindNext(scanFolderID)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrFindNextQueryFmt, err)
		os.Exit(1)
	}

	emitFindNext(rows, jsonOut)
}

// emitFindNext writes either JSON or the human-readable summary.
func emitFindNext(rows []model.FindNextRow, jsonOut bool) {
	if jsonOut {
		emitFindNextJSON(rows)

		return
	}
	emitFindNextText(rows)
}

// emitFindNextJSON dumps the result array as stablejson to stdout.
// Encoder lives in findnextrender.go; CONTRACT for top-level key
// order is pinned by gitmap/cmd/findnextjson_contract_test.go AND
// gitmap/cmd/findnext_jsonschema_contract_test.go (the latter
// cross-checks against spec/08-json-schemas/find-next.schema.json).
func emitFindNextJSON(rows []model.FindNextRow) {
	if err := encodeFindNextJSON(os.Stdout, rows); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrFindNextJSONEncodeFmt, err)
		os.Exit(1)
	}
}

// emitFindNextText prints the human summary (header + per-repo rows + hint).
func emitFindNextText(rows []model.FindNextRow) {
	if len(rows) == 0 {
		fmt.Print(constants.MsgFindNextEmpty)

		return
	}

	fmt.Printf(constants.MsgFindNextHeaderFmt, len(rows))
	for _, r := range rows {
		fmt.Printf(constants.MsgFindNextRowFmt,
			r.Repo.Slug, r.NextVersionTag, r.Method, r.ProbedAt, r.Repo.AbsolutePath)
	}
	fmt.Print(constants.MsgFindNextDoneFmt)
}
