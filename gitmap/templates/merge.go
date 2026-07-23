// Package templates — Merge primitive.
//
// Merge inserts (or updates in place) a gitmap-managed marker block inside
// a target file (.gitattributes or .gitignore). The block is bracketed by
// well-known sentinel comments so re-running the merge is idempotent:
//
//	# >>> gitmap:<kind>/<lang> >>>
//	... template body, verbatim ...
//	# <<< gitmap:<kind>/<lang> <<<
//
// On the second and later runs, the existing marker block is located by
// regex, replaced with the new body, and the rest of the file is kept
// byte-for-byte. Hand edits OUTSIDE the block survive untouched; hand
// edits INSIDE the block are intentionally overwritten — the user is
// expected to fork the template to ~/.gitmap/templates/ if they want
// custom content.
package templates

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

// MergeOutcome describes what Merge did to the target file.
type MergeOutcome int

const (
	// MergeCreated means the target file did not exist and was created.
	MergeCreated MergeOutcome = iota
	// MergeInserted means the file existed but had no prior gitmap block;
	// the new block was appended at the end.
	MergeInserted
	// MergeUpdated means a prior gitmap block was found and its body was
	// replaced (or kept identical — see Changed).
	MergeUpdated
)

// MergeResult is the structured return value of Merge.
type MergeResult struct {
	Path     string       // absolute path of the target file
	Outcome  MergeOutcome // what happened on disk
	Changed  bool         // true when bytes on disk differ from before
	BlockTag string       // e.g. "lfs/common" — the marker tag used
}

// Merge writes (or refreshes) a gitmap-managed marker block in targetPath
// containing body, identified by tag (e.g. "lfs/common"). It is safe to
// call repeatedly: the second call is a no-op when body has not changed.
func Merge(targetPath, tag string, body []byte) (MergeResult, error) {
	abs, err := filepath.Abs(targetPath)
	if err != nil {
		return MergeResult{}, fmt.Errorf("resolve %q: %w", targetPath, err)
	}

	prior, existed, readErr := readIfExists(abs)
	if readErr != nil {
		return MergeResult{}, readErr
	}

	block := buildBlock(tag, body)

	next, outcome := composeNext(prior, existed, tag, block)

	if bytes.Equal(prior, next) {
		return MergeResult{Path: abs, Outcome: outcome, Changed: false, BlockTag: tag}, nil
	}

	if mkErr := os.MkdirAll(filepath.Dir(abs), 0o755); mkErr != nil {
		return MergeResult{}, fmt.Errorf("mkdir %q: %w", filepath.Dir(abs), mkErr)
	}
	if wErr := os.WriteFile(abs, next, 0o644); wErr != nil {
		return MergeResult{}, fmt.Errorf("write %q: %w", abs, wErr)
	}

	return MergeResult{Path: abs, Outcome: outcome, Changed: true, BlockTag: tag}, nil
}

// readIfExists returns the file contents, or nil + existed=false when the
// file is missing. Other errors propagate.
func readIfExists(path string) (data []byte, existed bool, err error) {
	data, err = os.ReadFile(path)
	if err == nil {
		return data, true, nil
	}
	if os.IsNotExist(err) {
		return nil, false, nil
	}

	return nil, false, fmt.Errorf("read %q: %w", path, err)
}

// composeNext returns the new file bytes and the structural outcome
// (created / inserted / updated). It does NOT decide whether bytes
// actually changed — Merge handles that with bytes.Equal.
func composeNext(prior []byte, existed bool, tag string, block []byte) ([]byte, MergeOutcome) {
	if !existed {
		return block, MergeCreated
	}

	re := blockRegex(tag)
	if re.Match(prior) {
		return re.ReplaceAll(prior, block), MergeUpdated
	}

	return appendBlock(prior, block), MergeInserted
}

// buildBlock wraps body with the begin/end marker comments and a single
// trailing newline so subsequent appends start on a fresh line.
func buildBlock(tag string, body []byte) []byte {
	var b bytes.Buffer
	fmt.Fprintf(&b, "# >>> gitmap:%s >>>\n", tag)
	b.Write(trimTrailingNewlines(body))
	b.WriteByte('\n')
	fmt.Fprintf(&b, "# <<< gitmap:%s <<<\n", tag)

	return b.Bytes()
}

// blockRegex returns a regex that matches an existing marker block for
// tag, including its surrounding marker lines. (?s) lets `.` match
// newlines so the body in between can be any size.
func blockRegex(tag string) *regexp.Regexp {
	q := regexp.QuoteMeta(tag)
	pattern := `(?s)# >>> gitmap:` + q + ` >>>\n.*?# <<< gitmap:` + q + ` <<<\n?`

	return regexp.MustCompile(pattern)
}

// appendBlock places block at the end of prior, ensuring exactly one
// blank line of separation. Files without a trailing newline get one
// before the block so the marker comment starts cleanly on its own line.
func appendBlock(prior, block []byte) []byte {
	var b bytes.Buffer
	b.Write(prior)
	if !bytes.HasSuffix(prior, []byte("\n")) {
		b.WriteByte('\n')
	}
	if !bytes.HasSuffix(prior, []byte("\n\n")) && len(prior) > 0 {
		b.WriteByte('\n')
	}
	b.Write(block)

	return b.Bytes()
}

// trimTrailingNewlines collapses any trailing newlines on body so the
// block has a single, predictable terminator regardless of how the
// embedded template happens to end.
func trimTrailingNewlines(body []byte) []byte {
	out := body
	for len(out) > 0 && out[len(out)-1] == '\n' {
		out = out[:len(out)-1]
	}

	return out
}
