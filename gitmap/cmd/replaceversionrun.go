package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runReplaceLiteral implements `gitmap replace "<old>" "<new>"`.
func runReplaceLiteral(oldS, newS string, opts replaceOpts) {
	if oldS == "" {
		fmt.Fprint(os.Stderr, constants.ErrReplaceEmptyOld)
		os.Exit(1)
	}
	root := repoRoot()
	files := loadRepoFiles(root, opts.exts, opts.extCaseIns)

	pair := replacePair{old: oldS, new: newS}
	hits, total := scanReplacements(files, []replacePair{pair})
	printHits(hits, pair, opts.quiet)
	fmt.Printf(constants.MsgReplaceSummary, len(hits), total)

	if total == 0 {
		fmt.Print(constants.MsgReplaceNoMatches)
		return
	}
	if opts.dryRun || !confirmLiteral(hits, total, opts) {
		return
	}
	commitHits(hits)
}

// runReplaceVersion implements `-N` and `all` modes. n==0 means all.
func runReplaceVersion(n int, opts replaceOpts, isAll bool) {
	base, k := detectVersion()
	targets := versionTargets(k, n)
	if len(targets) == 0 {
		fmt.Print(constants.MsgReplaceAlreadyAtV1)
		return
	}
	root := repoRoot()
	files := loadRepoFiles(root, opts.exts, opts.extCaseIns)

	hits, total := scanVersionTargets(files, base, k, targets, opts.quiet)
	fmt.Printf(constants.MsgReplaceSummary, len(hits), total)
	if total == 0 {
		fmt.Print(constants.MsgReplaceNoMatches)
		return
	}
	if opts.dryRun || !confirmVersion(targets, k, opts, isAll) {
		return
	}
	commitHits(hits)
}

// scanVersionTargets runs one scan per target version, accumulating
// hits across all passes. Per-target output is printed inline.
func scanVersionTargets(
	files []string, base string, k int, targets []int, quiet bool,
) ([]replaceHit, int) {
	all := make([]replaceHit, 0, 16)
	total := 0
	for _, t := range targets {
		pairs := pairsForTarget(base, t, k)
		hits, sum := scanReplacements(files, pairs)
		printHits(hits, pairs[0], quiet)
		all = append(all, hits...)
		total += sum
	}
	return all, total
}

// loadRepoFiles wraps walkRepoFiles with the standardized error exit.
// The exts allow-list comes from --ext (nil = include every text file)
// and caseInsensitive comes from --ext-case (default insensitive).
func loadRepoFiles(root string, exts []string, caseInsensitive bool) []string {
	files, err := walkRepoFiles(root, exts, caseInsensitive)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceWalk, err)
		os.Exit(2)
	}
	fmt.Printf(constants.MsgReplaceScanning, len(files), root)
	return files
}

// confirmLiteral asks for y/N (or honors --yes) for literal mode.
func confirmLiteral(hits []replaceHit, total int, opts replaceOpts) bool {
	if opts.yes {
		return true
	}
	fmt.Printf(constants.MsgReplaceConfirmLit, total, len(hits))
	if confirmYes() {
		return true
	}
	fmt.Print(constants.MsgReplaceAborted)
	return false
}

// confirmVersion asks for y/N (or honors --yes) for version mode.
func confirmVersion(targets []int, k int, opts replaceOpts, _ bool) bool {
	if opts.yes {
		return true
	}
	fmt.Printf(constants.MsgReplaceConfirmVer, targets[0], targets[len(targets)-1], k)
	if confirmYes() {
		return true
	}
	fmt.Print(constants.MsgReplaceAborted)
	return false
}

// commitHits flushes the in-memory rewrites to disk and prints summary.
func commitHits(hits []replaceHit) {
	files, total := applyHits(hits)
	fmt.Printf(constants.MsgReplaceApplied, total, files)
}
