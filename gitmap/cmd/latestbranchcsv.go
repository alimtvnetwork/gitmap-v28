// Package cmd — latest-branch CSV formatter.
//
// Split out from latestbranchoutput.go to keep that file under the
// 200-line code-style budget. CRLF line endings are forced on every
// row (header + data) so output is byte-identical across platforms
// (RFC 4180); pinned by gitmap/cmd/csvcrlf_contract_test.go.
package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// printLatestCSV outputs the latest branch result as CSV to stdout.
func printLatestCSV(items []gitutil.RemoteBranchInfo, remote string, top int) {
	if err := encodeLatestBranchCSV(os.Stdout, items, remote, top); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to write CSV: %v\n", err)
	}
}

// encodeLatestBranchCSV writes the CSV to w. Takes an io.Writer so
// contract tests can capture bytes into a buffer for byte-exact
// CRLF / comma assertions.
func encodeLatestBranchCSV(
	w io.Writer, items []gitutil.RemoteBranchInfo, remote string, top int,
) error {
	count := resolveTopCount(top, len(items))
	cw := csv.NewWriter(w)
	cw.UseCRLF = true
	if err := cw.Write(constants.LatestBranchCSVHeaders); err != nil {

		return err
	}
	for _, item := range items[:count] {
		writeCSVRow(cw, item, remote)
	}
	cw.Flush()

	return cw.Error()
}

// resolveTopCount determines how many items to display.
func resolveTopCount(top, total int) int {
	count := 1
	if top > 0 {
		count = top
	}
	if count > total {
		count = total
	}

	return count
}

// writeCSVRow writes a single CSV row for a branch item.
func writeCSVRow(w *csv.Writer, item gitutil.RemoteBranchInfo, remote string) {
	if err := w.Write([]string{
		gitutil.StripRemotePrefix(item.RemoteRef),
		remote,
		gitutil.TruncSha(item.Sha),
		gitutil.FormatDisplayDate(item.CommitDate),
		item.Subject,
		item.RemoteRef,
	}); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to write CSV row: %v\n", err)
	}
}
