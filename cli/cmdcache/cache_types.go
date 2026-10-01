package cmdcache

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const maxCacheFileSize = 200 * 1024 // 200 KB
const maxJsonFileSize = 150 * 1024  // 150 KB

// SearchMatch represents a single search match in a cached file.
type SearchMatch struct {
	RelativePath string                   `json:"RelativePath"`
	LineNumber   int                      `json:"LineNumber"`
	Content      string                   `json:"Content"`
	ContextLines []store.CacheContextLine `json:"ContextLines,omitempty"`
}

// CacheSearchOptions holds filtering and formatting flags for search operations.
type CacheSearchOptions struct {
	Patterns    []string `json:"Patterns"`
	FileGlobs   []string `json:"FileGlobs"`
	LinesToShow int      `json:"LinesToShow"`
	ResultLimit int      `json:"ResultLimit"`
	IsRegex     bool     `json:"IsRegex"`
}

// CacheCreateOptions holds options for cache creation.
type CacheCreateOptions struct {
	Targets []string `json:"Targets"`
	IsKeep  bool     `json:"IsKeep"`
}

// FolderStats tracks indexed file count and total size per top-level folder.
type FolderStats struct {
	Count int
	Bytes int64
}

func parseCreateOptions(args []string) CacheCreateOptions {
	opts := CacheCreateOptions{IsKeep: false}
	for _, a := range args {
		trimmed := strings.TrimSpace(a)
		if trimmed == "--keep" || trimmed == "-k" {
			opts.IsKeep = true
		} else if !strings.HasPrefix(trimmed, "-") {
			parts := extractTargetParts(trimmed)
			opts.Targets = append(opts.Targets, parts...)
		}
	}
	if len(opts.Targets) == 0 {
		opts.Targets = []string{"."}
	}
	return opts
}

func extractTargetParts(raw string) []string {
	parts := strings.Split(raw, ",")
	var cleaned []string
	for _, p := range parts {
		trimmed := strings.Trim(strings.TrimSpace(p), `"'`)
		if trimmed != "" {
			cleaned = append(cleaned, trimmed)
		}
	}
	return cleaned
}

func parseSearchFlags(args []string) CacheSearchOptions {
	opts := CacheSearchOptions{LinesToShow: 10, ResultLimit: 20}
	for i := 0; i < len(args); i++ {
		a := args[i]
		i = handleFlagStep(args, i, a, &opts)
	}
	return opts
}

func handleFlagStep(args []string, i int, a string, opts *CacheSearchOptions) int {
	if (a == "--lines" || a == "-lines") && i+1 < len(args) {
		opts.LinesToShow, _ = strconv.Atoi(args[i+1])
		return i + 1
	}
	if (a == "--limit" || a == "-limit") && i+1 < len(args) {
		opts.ResultLimit, _ = strconv.Atoi(args[i+1])
		return i + 1
	}
	if isFilePatternFlag(a) && i+1 < len(args) {
		return handleFilePatternFlag(args, i+1, opts)
	}
	processPositionalSearchArg(a, opts)
	return i
}

func handleFilePatternFlag(args []string, nextIdx int, opts *CacheSearchOptions) int {
	val := args[nextIdx]
	if (val == "(fp)" || val == "fp") && nextIdx+1 < len(args) {
		nextIdx++
		val = args[nextIdx]
	}
	terms := extractCleanTerms(val)
	opts.FileGlobs = append(opts.FileGlobs, terms...)
	return nextIdx
}

func isFilePatternFlag(a string) bool {
	lower := strings.ToLower(strings.TrimSpace(a))
	return lower == "-file-pattern" || lower == "--file-pattern" ||
		lower == "-fp" || lower == "--fp" || lower == "(fp)" || lower == "fp"
}

func extractCleanTerms(raw string) []string {
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.Trim(strings.TrimSpace(p), `"', `)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func processPositionalSearchArg(a string, opts *CacheSearchOptions) {
	if strings.HasPrefix(a, "-") || a == "(fp)" {
		return
	}
	terms := extractCleanTerms(a)
	if len(opts.Patterns) == 0 {
		opts.Patterns = append(opts.Patterns, terms...)
		return
	}
	if len(opts.FileGlobs) == 0 && (strings.Contains(a, "*") || strings.Contains(a, "?")) {
		opts.FileGlobs = append(opts.FileGlobs, terms...)
		return
	}
	opts.Patterns = append(opts.Patterns, terms...)
}

func normalizeSearchOpts(opts CacheSearchOptions) CacheSearchOptions {
	if opts.LinesToShow <= 0 {
		opts.LinesToShow = 10
	}
	if opts.ResultLimit <= 0 {
		opts.ResultLimit = 20
	}
	return opts
}

func isFileEligible(info fs.FileInfo, absPath string, isKeep bool) bool {
	if info.IsDir() {
		return false
	}
	if isKeep {
		return isTextFile(absPath)
	}
	return isCacheSizeAcceptable(info, absPath) && isTextFile(absPath)
}

func isTextFile(absPath string) bool {
	isBinary := isBinaryFile(absPath)
	return !isBinary
}

func isCacheSizeAcceptable(info fs.FileInfo, absPath string) bool {
	if info.Size() > maxCacheFileSize {
		return false
	}
	return !isLargeJson(absPath, info.Size())
}

func isLargeJson(absPath string, size int64) bool {
	hasJsonSuffix := strings.HasSuffix(strings.ToLower(absPath), ".json")
	if !hasJsonSuffix {
		return false
	}
	return size > maxJsonFileSize
}

func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	hasErr := err != nil
	if hasErr {
		return true
	}
	defer f.Close()

	buf := make([]byte, 512)
	n, readErr := f.Read(buf)
	hasReadErr := readErr != nil
	if hasReadErr && n == 0 {
		return false
	}
	return bytes.Contains(buf[:n], []byte{0})
}

func isSkippedDir(name string) bool {
	lower := strings.ToLower(name)
	return lower == ".git" || lower == "node_modules" || lower == ".vscode" ||
		lower == ".idea" || lower == ".gitmap" || lower == "vendor"
}

func resolveFileSlug(relPath string) string {
	clean := filepath.ToSlash(relPath)
	clean = strings.TrimPrefix(clean, "./")
	parts := strings.Split(clean, "/")
	if len(parts) > 1 && parts[0] != "" {
		return store.SanitizeSlug(parts[0])
	}
	return "root"
}

func findRepoRoot() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "."
	}
	return cwd
}
