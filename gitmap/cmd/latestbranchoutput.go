// Package cmd — latest-branch output formatters.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/gitutil"
)

// dispatchLatestOutput routes to the correct output formatter.
func dispatchLatestOutput(result latestBranchResult, items []gitutil.RemoteBranchInfo, cfg latestBranchConfig) {
	if cfg.format == constants.OutputJSON {
		printLatestJSON(result, items, cfg.top)

		return
	}
	if cfg.format == constants.OutputCSV {
		printLatestCSV(items, result.selectedRemote, cfg.top)

		return
	}
	printLatestTerminal(result, items, cfg.top)
}

// printLatestJSON outputs the latest branch result as JSON.
func printLatestJSON(result latestBranchResult, items []gitutil.RemoteBranchInfo, top int) {
	if err := encodeLatestBranchJSON(os.Stdout, result, items, top); err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ Failed to encode latest branch JSON: %v\n", err)
	}
}

// CSV output lives in latestbranchcsv.go (file-size budget split).

// printLatestTerminal outputs the latest branch result as text.
func printLatestTerminal(result latestBranchResult, items []gitutil.RemoteBranchInfo, top int) {
	fmt.Println()
	printTerminalHeader(result)
	if top > 0 {
		printTerminalTopTable(items, result.selectedRemote, top)
	}
	fmt.Println()
}

// printTerminalHeader prints the main latest-branch info block.
func printTerminalHeader(result latestBranchResult) {
	fmt.Printf(constants.LBTermLatestFmt, strings.Join(result.branchNames, ", "))
	fmt.Printf(constants.LBTermRemoteFmt, result.selectedRemote)
	fmt.Printf(constants.LBTermSHAFmt, result.shortSha)
	fmt.Printf(constants.LBTermDateFmt, result.commitDate)
	fmt.Printf(constants.LBTermSubjectFmt, result.latest.Subject)
	fmt.Printf(constants.LBTermRefFmt, result.latest.RemoteRef)
}

// printTerminalTopTable prints the top-N branches table.
func printTerminalTopTable(items []gitutil.RemoteBranchInfo, remote string, top int) {
	count := resolveTopCount(top, len(items))
	fmt.Println()
	fmt.Printf(constants.LBTermTopHdrFmt, count, remote)
	printTerminalTopHeader()
	for _, item := range items[:count] {
		printTerminalTopRow(item)
	}
}

// printTerminalTopHeader prints the table column headers.
func printTerminalTopHeader() {
	fmt.Printf(constants.LBTermRowFmt,
		constants.LatestBranchTableColumns[0], constants.LatestBranchTableColumns[1],
		constants.LatestBranchTableColumns[2], constants.LatestBranchTableColumns[3])
}

// printTerminalTopRow prints a single branch row.
func printTerminalTopRow(item gitutil.RemoteBranchInfo) {
	fmt.Printf(constants.LBTermRowFmt,
		gitutil.FormatDisplayDate(item.CommitDate),
		gitutil.StripRemotePrefix(item.RemoteRef),
		gitutil.TruncSha(item.Sha),
		item.Subject)
}
