// Package cmd: audit-legacy scans the workspace for forbidden legacy
// strings (default: gitmap-v27 / gitmap-v27 / gitmap-v27) and exits 1 on // gitmap-legacy-ref-allow
// any hit. Designed as a regression guard for remixes / rename commits.
package cmd

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// auditLegacyHit is one matched line.
type auditLegacyHit struct {
	File    string `json:"file"`
	Line    int    `json:"line"`
	Pattern string `json:"pattern"`
	Text    string `json:"text"`
}

// auditLegacyOpts holds parsed CLI flags.
type auditLegacyOpts struct {
	Patterns   []*regexp.Regexp
	Raw        []string
	Root       string
	AsJSON     bool
	ReportPath string // empty = no report file written
	WriteDiffs bool   // --diffs: emit per-file unified diffs alongside the report
}

// runAuditLegacy is the dispatch entry point.
func runAuditLegacy(args []string) {
	checkHelp(constants.CmdAuditLegacy, args)
	opts, err := parseAuditLegacyArgs(args)
	if err != nil {
		cliexit.Fail(constants.CmdAuditLegacy, "parse-args", "", err, 2)
	}
	hits, n, walkErr := scanAuditLegacy(opts)
	if walkErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrAuditLegacyWalk, opts.Root, walkErr)
		os.Exit(2)
	}
	emitAuditLegacy(opts, hits, n)
	plans := writeAuditLegacyDiffs(opts, hits)
	writeAuditLegacyReport(opts, hits, n, plans)
	if len(hits) > 0 {
		os.Exit(1)
	}
}

// (parseAuditLegacyArgs and pattern helpers live in auditlegacy_parse.go)

// scanAuditLegacy walks the root and returns hits + scanned-file count.
func scanAuditLegacy(opts auditLegacyOpts) ([]auditLegacyHit, int, error) {
	state := &auditWalkState{patterns: opts.Patterns}
	walkErr := filepath.WalkDir(opts.Root, state.visit)

	return state.hits, state.fileCount, walkErr
}

// auditWalkState accumulates results during filepath.WalkDir.
type auditWalkState struct {
	patterns  []*regexp.Regexp
	hits      []auditLegacyHit
	fileCount int
}

// visit is the WalkDir callback.
func (s *auditWalkState) visit(path string, d fs.DirEntry, err error) error {
	if err != nil || d == nil {
		return nil
	}
	if d.IsDir() {
		return skipAuditDir(d.Name())
	}
	if !isAuditScannable(path) {
		return nil
	}
	s.fileCount++
	s.hits = append(s.hits, scanAuditLegacyFile(path, s.patterns)...)

	return nil
}

// skipAuditDir returns fs.SkipDir for ignored directories.
func skipAuditDir(name string) error {
	switch name {
	case ".git", "node_modules", "dist", "build", "bin", ".next",
		".gitmap", "vendor", "coverage":
		return fs.SkipDir
	}

	return nil
}

// isAuditScannable filters out binary / non-text files by extension.
func isAuditScannable(path string) bool {
	skip := []string{".png", ".jpg", ".jpeg", ".gif", ".webp", ".ico",
		".pdf", ".zip", ".gz", ".tar", ".exe", ".dll", ".so",
		".dylib", ".bin", ".db", ".sqlite", ".woff", ".woff2", ".ttf"}
	low := strings.ToLower(path)
	for _, ext := range skip {
		if strings.HasSuffix(low, ext) {
			return false
		}
	}

	return true
}

// scanAuditLegacyFile reads one file and returns all matches.
func scanAuditLegacyFile(path string, pats []*regexp.Regexp) []auditLegacyHit {
	f, err := os.Open(path)
	if err != nil {
		return nil
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	return collectAuditLineHits(path, scanner, pats)
}

// collectAuditLineHits iterates the scanner and gathers matches.
func collectAuditLineHits(path string, scanner *bufio.Scanner, pats []*regexp.Regexp) []auditLegacyHit {
	var hits []auditLegacyHit
	lineNo := 0
	for scanner.Scan() {
		lineNo++
		hits = append(hits, matchAuditLine(path, lineNo, scanner.Text(), pats)...)
	}

	return hits
}

// matchAuditLine returns one hit per matching pattern on a single line.
func matchAuditLine(path string, lineNo int, line string, pats []*regexp.Regexp) []auditLegacyHit {
	var hits []auditLegacyHit
	for _, re := range pats {
		if re.MatchString(line) {
			hits = append(hits, auditLegacyHit{
				File: path, Line: lineNo, Pattern: re.String(), Text: line,
			})
		}
	}

	return hits
}
