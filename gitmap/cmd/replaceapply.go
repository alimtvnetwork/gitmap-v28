package cmd

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// replacePair is one literal find/replace combo. Version mode generates
// two pairs per target version (the `-vN` and `/vN` forms).
type replacePair struct {
	old string
	new string
}

// replaceHit captures the per-file outcome of a scan pass.
type replaceHit struct {
	path    string
	count   int
	updated []byte
}

// scanReplacements visits every file once and applies all pairs in
// memory. Files with zero matches are dropped from the result.
func scanReplacements(files []string, pairs []replacePair) ([]replaceHit, int) {
	hits := make([]replaceHit, 0, 32)
	total := 0
	for _, f := range files {
		hit, ok := scanOneFile(f, pairs)
		if !ok {
			continue
		}
		hits = append(hits, hit)
		total += hit.count
	}
	return hits, total
}

// scanOneFile reads one file and applies every pair. Returns ok=false
// when the file has no matches or cannot be read.
func scanOneFile(path string, pairs []replacePair) (replaceHit, bool) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrReplaceWrite, path, err)
		return replaceHit{}, false
	}
	count := 0
	for _, p := range pairs {
		c := bytes.Count(data, []byte(p.old))
		if c == 0 {
			continue
		}
		count += c
		data = bytes.ReplaceAll(data, []byte(p.old), []byte(p.new))
	}
	if count == 0 {
		return replaceHit{}, false
	}
	return replaceHit{path: path, count: count, updated: data}, true
}

// applyHits writes every hit atomically: tempfile + rename.
func applyHits(hits []replaceHit) (int, int) {
	files, total := 0, 0
	for _, h := range hits {
		if err := atomicWrite(h.path, h.updated); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrReplaceWrite, h.path, err)
			continue
		}
		files++
		total += h.count
	}
	return files, total
}

// atomicWrite implements the temp+rename contract from §7 of the spec.
func atomicWrite(path string, data []byte) error {
	dir := filepath.Dir(path)
	tmp, err := os.CreateTemp(dir, filepath.Base(path)+".gitmap-replace-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		os.Remove(tmpName)
		return err
	}
	tmp.Close()
	return os.Rename(tmpName, path)
}

// printHits prints per-file summaries for one pair (or a generic header
// for literal mode where all pairs are user-supplied).
func printHits(hits []replaceHit, pair replacePair, quiet bool) {
	if quiet {
		return
	}
	for _, h := range hits {
		fmt.Printf(constants.MsgReplaceFileMatch, h.path, h.count, pair.old, pair.new)
	}
}

// confirmYes reads y/Y from stdin and reports whether the user agreed.
func confirmYes() bool {
	var ans string
	if _, err := fmt.Scanln(&ans); err != nil {
		return false
	}
	ans = strings.TrimSpace(ans)
	return ans == "y" || ans == "Y"
}
