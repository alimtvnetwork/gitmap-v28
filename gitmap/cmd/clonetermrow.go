package cmd

// clonetermrow.go — single-row helpers that build the standardized
// RepoTermBlock input for one clone-now / clone-from row at a time.
// Used by the streaming BeforeRow hooks so each block prints
// IMMEDIATELY before its row's git clone shells out (vs. the
// previous behavior of dumping every block up-front before
// Execute started).
//
// Design rationale: clonetermplan.go's batch helpers iterate the
// whole Plan and print N blocks back-to-back. That made the user
// wait through every ls-remote round-trip before any clone progress
// appeared. Streaming flips that — one ls-remote → one block → one
// clone → repeat — which matches what `gitmap clone <url1> <url2>`
// already does and gives the user a live, scannable transcript.

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonefrom"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenow"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// printCloneNowTermBlockRow emits one RepoTermBlock for one
// clone-now row. Mirrors the per-row branch-source logic from
// printCloneNowTermBlocks (the batch variant) so streaming and
// upfront output stay byte-for-byte identical per row — only the
// timing changes.
//
// Faithfulness: clonenow's executor (clonenow/execute.go
// buildGitArgs) only passes `-b` when row.Branch is non-empty —
// the ls-remote-detected fallback we show on the `branch:` line
// is informational only. CmdBranch is therefore pinned to
// row.Branch (NOT the detected fallback) so the printed cmd
// matches the real argv exactly. URL/dest are passed through
// from the executor's resolveRowDisplay.
func printCloneNowTermBlockRow(index, total int, row clonenow.Row,
	url, dest string) {
	_ = total // total is reserved for a future "[i/N]" prefix; unused today
	branch := row.Branch
	source := "manifest"
	if len(branch) == 0 {
		branch = detectRemoteHEAD(url)
		source = remoteBranchSource(branch)
	}
	in := CloneTermBlockInput{
		Index:        index,
		Name:         pickCloneNowName(row, dest),
		Branch:       branch,
		BranchSource: source,
		OriginalURL:  url,
		TargetURL:    url,
		Dest:         dest,
		CmdBranch:    row.Branch, // executor uses row.Branch, NOT detected
		// Non-nil empty Pre opts into pickCmdBranch's explicit
		// "no -b" sentinel when row.Branch is empty. Without this,
		// pickCmdBranch's legacy fallback renders `-b <detected>` in
		// the cmd: line even though clonenow.BuildGitArgs emits no
		// -b flag — the documented v5.x verify-cmd-faithful drift.
		CmdExtraArgsPre: cmdExtraArgsPreForRowBranch(row.Branch),
	}
	maybePrintCloneTermBlock(constants.OutputTerminal, in)
	// Cross-check the printed cmd: against the executor's real argv
	// (no-op unless --verify-cmd-faithful is on). Single source of
	// truth: clonenow.BuildGitArgs is the same builder Execute calls.
	runCmdFaithfulCheck(in, clonenow.BuildGitArgs(row, url, dest))
	runCmdPrintArgv(clonenow.BuildGitArgs(row, url, dest))
}

// cmdExtraArgsPreForRowBranch returns the explicit "no -b" sentinel
// (non-nil empty slice) when rowBranch is empty so pickCmdBranch
// does NOT fall back to in.Branch — keeping the displayed cmd:
// faithful to the executor argv, which itself emits -b only when
// rowBranch is non-empty. When rowBranch is set we return nil so
// downstream Post-flag appenders keep working unchanged.
func cmdExtraArgsPreForRowBranch(rowBranch string) []string {
	if len(rowBranch) == 0 {
		return []string{}
	}

	return nil
}

// printCloneFromTermBlockRow emits one RepoTermBlock for one
// clone-from row. clone-from never rewrites URLs, so OriginalURL
// and TargetURL are both the row's URL. Branch source mirrors
// clonefrom.branchSourceForRow ("manifest" if pinned, else the
// ls-remote-discovered default — falling through to "(unknown)"
// inside the renderer if detection fails).
//
// Faithfulness: the executor (clonefrom/execute.go buildGitArgs)
// only passes `-b` when row.Branch is non-empty, and adds
// `--depth=N` AFTER `-b` when row.Depth > 0. The printed cmd
// mirrors that exactly via CmdBranch (= row.Branch, NOT the
// ls-remote fallback) and CmdExtraArgsPost.
func printCloneFromTermBlockRow(index, total int, row clonefrom.Row,
	dest string) {
	_ = total // reserved for future "[i/N]" prefix
	branch := row.Branch
	source := "manifest"
	if len(branch) == 0 {
		branch = detectRemoteHEAD(row.URL)
		source = remoteBranchSource(branch)
	}
	var post []string
	if row.Depth > 0 {
		post = []string{fmt.Sprintf(constants.CloneFromDepthFlagFmt, row.Depth)}
	}
	in := CloneTermBlockInput{
		Index:            index,
		Name:             dest,
		Branch:           branch,
		BranchSource:     source,
		OriginalURL:      row.URL,
		TargetURL:        row.URL,
		Dest:             dest,
		CmdBranch:        row.Branch, // executor uses row.Branch, NOT detected
		CmdExtraArgsPre:  cmdExtraArgsPreForRowBranch(row.Branch),
		CmdExtraArgsPost: post,
	}
	maybePrintCloneTermBlock(constants.OutputTerminal, in)
	// Cross-check displayed cmd: vs the executor's real argv. The
	// argv builder is the SAME function execute.go calls, so any
	// future drift between the two paths surfaces here loudly.
	runCmdFaithfulCheck(in, clonefrom.BuildGitArgs(row, dest))
	runCmdPrintArgv(clonefrom.BuildGitArgs(row, dest))
}
