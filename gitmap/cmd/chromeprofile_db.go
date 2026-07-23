// Package cmd — chromeprofile_db.go: thin glue between the chrome-profile
// command runners and the SQLite store helpers. Failures are non-fatal
// for the CLI (we still want the on-disk JSON/CSV to land) but always
// surface as a stderr warning so users see drift.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// chromeExportRecord groups the two snapshot artifacts a single
// export emits. Either path may be empty if that format was skipped.
type chromeExportRecord struct {
	JSONPath string
	JSONSize int
	CSVPath  string
	CSVSize  int
}

// persistChromeProfile upserts the profile row and inserts one
// ChromeProfileExport row per non-empty artifact. Soft-fails: any
// error is logged to stderr without aborting the CLI.
func persistChromeProfile(name, sourcePath string, rec chromeExportRecord) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgChromeProfileDBWarn, err)
		return
	}
	defer db.Close()
	id, err := db.UpsertChromeProfile(name, sourcePath, true)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgChromeProfileDBWarn, err)
		return
	}
	if rec.JSONPath != "" {
		if err := db.InsertChromeProfileExport(id, constants.OutputJSON, rec.JSONPath, rec.JSONSize); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgChromeProfileDBWarn, err)
			return
		}
	}
	if rec.CSVPath != "" {
		if err := db.InsertChromeProfileExport(id, constants.OutputCSV, rec.CSVPath, rec.CSVSize); err != nil {
			fmt.Fprintf(os.Stderr, constants.MsgChromeProfileDBWarn, err)
			return
		}
	}
	fmt.Printf(constants.MsgChromeProfileDBSynced, name)
}

// listChromeProfilesFromDB prints the tracked profiles section after
// the on-disk listing. Soft-fails when the DB is unreachable.
func listChromeProfilesFromDB() {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()
	rows, err := db.ListChromeProfilesDB()
	if err != nil || len(rows) == 0 {
		return
	}
	fmt.Print(constants.MsgChromeProfileListDBHdr)
	for _, r := range rows {
		fmt.Printf(constants.MsgChromeProfileListDBRow, r.Name, r.ExportCount, r.LastSeen)
	}
}
