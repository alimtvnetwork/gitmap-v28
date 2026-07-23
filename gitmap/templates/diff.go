// Package templates — Diff primitive.
//
// Diff compares the current on-disk content of a target file against
// what `add <kind> <lang>` would produce, focused exclusively on the
// gitmap-managed marker block. Hand edits OUTSIDE the block are
// invisible to the diff (consistent with Merge's contract: outside
// content is untouched, so a diff would only add noise).
//
// The output is a small unified-style hunk list. We deliberately do not
// pull a Myers diff dependency — comparing two block bodies line-by-line
// covers every case `templates diff` cares about, in ~80 LOC.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// DiffStatus enumerates the high-level outcomes of a single Diff call.
type DiffStatus int

const (
	// DiffNoChange means the template body matches the on-disk block
	// (or both are absent). Exit code 0 territory.
	DiffNoChange DiffStatus = iota
	// DiffMissingBlock means the file exists but has no gitmap block
	// for tag — `add` would insert one.
	DiffMissingBlock
	// DiffMissingFile means the file itself is absent — `add` would
	// create it.
	DiffMissingFile
	// DiffBlockChanged means the file has a gitmap block but its body
	// differs from the template — `add` would update it.
	DiffBlockChanged
)

// DiffResult is the structured outcome of Diff. Hunks is empty when
// Status == DiffNoChange.
type DiffResult struct {
	Path     string
	Tag      string
	Status   DiffStatus
	Hunks    []string // unified-style "+/-" lines, banner-prefixed
	BlockOld []byte   // body bytes currently on disk (nil if absent)
	BlockNew []byte   // body bytes the template would write
}

// Diff loads targetPath, locates the gitmap:<tag> block (if any), and
// compares its body to body. Returns a structured result; never writes
// to disk. The caller decides how to render Hunks (raw/ANSI).
func Diff(targetPath, tag string, body []byte) (DiffResult, error) {
	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return DiffResult{}, fmt.Errorf("resolve %q: %w", targetPath, err)
	}

	prior, existed, readErr := readIfExists(abs)
	if readErr != nil {
		return DiffResult{}, readErr
	}

	wantBody := bytes.TrimRight(body, "\n")
	res := DiffResult{Path: abs, Tag: tag, BlockNew: wantBody}

	if !existed {
		res.Status = DiffMissingFile
		res.Hunks = renderAdditionHunk(tag, wantBody)

		return res, nil
	}

	gotBody, found := extractBlockBody(prior, tag)
	if !found {
		res.Status = DiffMissingBlock
		res.Hunks = renderAdditionHunk(tag, wantBody)

		return res, nil
	}

	res.BlockOld = gotBody
	if bytes.Equal(gotBody, wantBody) {
		res.Status = DiffNoChange

		return res, nil
	}

	res.Status = DiffBlockChanged
	res.Hunks = renderChangeHunk(tag, gotBody, wantBody)

	return res, nil
}

// extractBlockBody returns the body lines between the marker comments
// for tag. found=false means the block is absent.
func extractBlockBody(file []byte, tag string) (body []byte, found bool) {
	re := blockRegex(tag)
	loc := re.FindIndex(file)
	if loc == nil {
		return nil, false
	}
	matched := file[loc[0]:loc[1]]
	// Drop the leading `# >>> gitmap:tag >>>\n` and trailing
	// `# <<< gitmap:tag <<<\n?` lines to recover the raw body.
	header := []byte(fmt.Sprintf("# >>> gitmap:%s >>>\n", tag))
	footer := []byte(fmt.Sprintf("# <<< gitmap:%s <<<", tag))
	stripped := bytes.TrimPrefix(matched, header)
	if i := bytes.LastIndex(stripped, footer); i >= 0 {
		stripped = stripped[:i]
	}

	return bytes.TrimRight(stripped, "\n"), true
}

// renderAdditionHunk formats every line of body as a `+` addition under
// a single `@@ gitmap:<tag> @@` banner. Used when the block is wholly
// new (file missing, or block missing inside an existing file).
func renderAdditionHunk(tag string, body []byte) []string {
	out := []string{fmt.Sprintf("@@ gitmap:%s @@", tag)}
	for _, line := range splitDiffLines(body) {
		out = append(out, "+"+line)
	}

	return out
}

// renderChangeHunk emits a paired `-old / +new` hunk. We don't try to
// align matching lines (no LCS) — the marker block is small enough
// that a flat removal-then-addition is honest and readable.
func renderChangeHunk(tag string, oldBody, newBody []byte) []string {
	out := []string{fmt.Sprintf("@@ gitmap:%s @@", tag)}
	for _, line := range splitDiffLines(oldBody) {
		out = append(out, "-"+line)
	}
	for _, line := range splitDiffLines(newBody) {
		out = append(out, "+"+line)
	}

	return out
}

// splitDiffLines splits on "\n" and drops the trailing empty token from
// a final newline. Preserves blank lines in between so they show up as
// "+" / "-" entries rather than getting silently swallowed.
func splitDiffLines(body []byte) []string {
	if len(body) == 0 {
		return nil
	}
	s := string(body)
	parts := strings.Split(s, "\n")
	if parts[len(parts)-1] == "" {
		parts = parts[:len(parts)-1]
	}

	return parts
}

// LoadFile is a thin os.ReadFile shim re-exported for callers (cmd
// package) that want to fall back gracefully when the target is
// missing without re-implementing the os.IsNotExist check.
func LoadFile(path string) ([]byte, bool, error) {
	data, err := os.ReadFile(path)
	if err == nil {
		return data, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}

	return nil, false, fmt.Errorf("read %q: %w", path, err)
}
