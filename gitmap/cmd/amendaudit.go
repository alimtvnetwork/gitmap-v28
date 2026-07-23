// Package cmd — amendaudit.go handles audit JSON writing and DB persistence.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// writeAmendAudit writes the audit JSON file to .gitmap/amendments/.
func writeAmendAudit(f amendFlags, commits []model.CommitEntry, branch, mode, prevName, prevEmail string) string {
	dir := constants.AmendAuditDir
	ensureAmendDir(dir)

	ts := time.Now().UTC()
	fileName := constants.AmendAuditFilePrefix + formatAuditTimestamp(ts) + ".json"
	path := filepath.Join(dir, fileName)

	record := buildAuditRecord(f, commits, branch, mode, prevName, prevEmail, ts)
	var buf bytes.Buffer
	if err := encodeAmendAuditJSON(&buf, record); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendAuditWrite, err)

		return path
	}

	if err := os.WriteFile(path, buf.Bytes(), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAmendAuditWrite, err)
	}

	return path
}

// buildAuditRecord constructs the AmendmentRecord for JSON output.
func buildAuditRecord(f amendFlags, commits []model.CommitEntry, branch, mode, prevName, prevEmail string, ts time.Time) model.AmendmentRecord {
	fromCommit := ""
	toCommit := ""

	if len(commits) > 0 {
		fromCommit = commits[0].SHA
		toCommit = commits[len(commits)-1].SHA
	}

	return model.AmendmentRecord{
		Timestamp:    ts.Format(time.RFC3339),
		Branch:       branch,
		FromCommit:   fromCommit,
		ToCommit:     toCommit,
		TotalCommits: len(commits),
		PreviousAuthor: model.AmendAuthor{
			Name:  prevName,
			Email: prevEmail,
		},
		NewAuthor: model.AmendAuthor{
			Name:  f.name,
			Email: f.email,
		},
		Mode:        mode,
		ForcePushed: f.forcePush,
		Commits:     commits,
	}
}

// saveAmendToDB persists the amendment record to the SQLite database.
func saveAmendToDB(f amendFlags, commits []model.CommitEntry, branch, mode, prevName, prevEmail string) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not save amendment to database: %v\n", err)

		return
	}
	defer db.Close()

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: database migration failed: %v\n", err)

		return
	}

	fromCommit := ""
	toCommit := ""

	if len(commits) > 0 {
		fromCommit = commits[0].SHA
		toCommit = commits[len(commits)-1].SHA
	}

	if err := db.InsertAmendment(branch, fromCommit, toCommit, len(commits),
		prevName, prevEmail, f.name, f.email, mode, f.forcePush); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not save amendment to DB: %v\n", err)
	}
}

// ensureAmendDir creates the audit directory if it doesn't exist.
func ensureAmendDir(dir string) {
	if err := os.MkdirAll(dir, constants.DirPermission); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not create audit directory %s: %v\n", dir, err)
	}
}

// formatAuditTimestamp formats a time for use in filenames.
func formatAuditTimestamp(t time.Time) string {
	s := t.Format("2006-01-02T15-04-05")

	return strings.ReplaceAll(s, ":", "-")
}
