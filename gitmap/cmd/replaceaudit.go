package cmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runReplaceAudit implements `gitmap replace --audit`. It scans every
// eligible file for any older `<base>-vT` or `<base>/vT` reference and
// prints `path:line: matched-text` per occurrence. Never writes.
// Honors --ext to restrict the scan surface.
func runReplaceAudit(opts replaceOpts) {
	base, k := detectVersion()
	targets := versionTargets(k, 0)
	if len(targets) == 0 {
		fmt.Print(constants.MsgReplaceAlreadyAtV1)
		return
	}
	root := repoRoot()
	files := loadRepoFiles(root, opts.exts, opts.extCaseIns)

	needles := buildAuditNeedles(base, targets)
	totalHits := scanAudit(files, needles)
	if totalHits == 0 {
		fmt.Print(constants.MsgReplaceAuditClean)
	}
}

// buildAuditNeedles flattens every (target -> dash form, slash form)
// pair into a single search list for the audit scan.
func buildAuditNeedles(base string, targets []int) [][]byte {
	out := make([][]byte, 0, len(targets)*2)
	for _, t := range targets {
		out = append(out,
			[]byte(fmt.Sprintf("%s-v%d", base, t)),
			[]byte(fmt.Sprintf("%s/v%d", base, t)),
		)
	}
	return out
}

// scanAudit prints per-line hits across all files. Returns the total
// number of lines reported so the caller can print a clean message.
func scanAudit(files []string, needles [][]byte) int {
	total := 0
	for _, f := range files {
		total += scanAuditFile(f, needles)
	}
	return total
}

// scanAuditFile streams one file line-by-line and prints any line that
// contains at least one needle.
func scanAuditFile(path string, needles [][]byte) int {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceWrite, path, err)
		return 0
	}
	hits := 0
	lineNum := 0
	for _, line := range bytes.Split(data, []byte("\n")) {
		lineNum++
		if lineContainsAny(line, needles) {
			fmt.Fprintf(os.Stdout, constants.MsgReplaceAuditMatch, path, lineNum, string(line))
			hits++
		}
	}
	return hits
}

// lineContainsAny returns true when the line contains any needle.
func lineContainsAny(line []byte, needles [][]byte) bool {
	for _, n := range needles {
		if bytes.Contains(line, n) {
			return true
		}
	}
	return false
}
