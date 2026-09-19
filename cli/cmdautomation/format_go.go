package cmdautomation

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type normalizedContent struct {
	data    []byte
	hasBom  bool
	hasCrlf bool
}

type formatTotals struct {
	formattedCount int
	violationCount int
}

// RunFormatGo formats Go source files, organizing imports and normalizing encoding.
func RunFormatGo(opts FormatGoOptions) FormatGoResultMonad {
	start := time.Now()
	baseDir := resolveArtifactBaseDir(opts.Dir)
	files := collectGoFiles(opts, baseDir)
	items := make([]FormatGoItem, 0, len(files))
	for _, f := range files {
		items = append(items, formatSingleGoFile(f, opts))
	}
	res := aggregateFormatGo(items, start)
	return result.Ok(res)
}

func stripBomIfPresent(raw []byte) []byte {
	if bytes.HasPrefix(raw, []byte("\xef\xbb\xbf")) {
		return bytes.TrimPrefix(raw, []byte("\xef\xbb\xbf"))
	}
	return raw
}

func normalizeCrlfIfPresent(data []byte) []byte {
	if bytes.Contains(data, []byte("\r\n")) {
		return bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
	}
	return data
}

func normalizeFileBytes(raw []byte) normalizedContent {
	hasBom := bytes.HasPrefix(raw, []byte("\xef\xbb\xbf"))
	stripped := stripBomIfPresent(raw)
	hasCrlf := bytes.Contains(stripped, []byte("\r\n"))
	cleaned := normalizeCrlfIfPresent(stripped)
	return normalizedContent{
		data:    cleaned,
		hasBom:  hasBom,
		hasCrlf: hasCrlf,
	}
}

func parseAndSortAst(content []byte) []byte {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", content, parser.ParseComments)
	if err != nil {
		return content
	}
	ast.SortImports(fset, file)
	var buf bytes.Buffer
	if err := format.Node(&buf, fset, file); err != nil {
		return content
	}
	return buf.Bytes()
}

func formatGoSource(content []byte) []byte {
	sorted := parseAndSortAst(content)
	formatted, err := format.Source(sorted)
	if err != nil {
		return sorted
	}
	return formatted
}

func buildFormatItem(path string, norm normalizedContent, hasReorder bool, isChanged bool) FormatGoItem {
	return FormatGoItem{
		Path:               path,
		IsFormatted:        isChanged,
		HasBom:             norm.hasBom,
		HasCrlf:            norm.hasCrlf,
		HasReorderImports: hasReorder,
	}
}

func processGoFileContent(raw []byte, path string, opts FormatGoOptions) FormatGoItem {
	norm := normalizeFileBytes(raw)
	formatted := formatGoSource(norm.data)
	isIdentical := bytes.Equal(raw, formatted)
	isChanged := (isIdentical == false)
	isSameAst := bytes.Equal(norm.data, formatted)
	hasReorder := (isSameAst == false)
	item := buildFormatItem(path, norm, hasReorder, isChanged)
	if opts.IsWrite && isChanged {
		_ = os.WriteFile(path, formatted, 0644)
	}
	return item
}

func formatSingleGoFile(path string, opts FormatGoOptions) FormatGoItem {
	raw, err := os.ReadFile(path)
	if err != nil {
		return FormatGoItem{
			Path:         path,
			ErrorMessage: err.Error(),
		}
	}
	return processGoFileContent(raw, path, opts)
}

func fetchStagedGoFiles(baseDir string) []string {
	out := fetchGitCommandOutput(baseDir, "diff", "--cached", "--name-only", "--diff-filter=ACM")
	lines := strings.Split(out, "\n")
	var staged []string
	for _, l := range lines {
		p := filepath.ToSlash(strings.TrimSpace(l))
		if strings.HasSuffix(p, ".go") && len(p) > 0 {
			staged = append(staged, p)
		}
	}
	return staged
}

func scanGoDirectory(dir string) []string {
	var files []string
	_ = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() && isSkipDirectory(d.Name()) {
			return filepath.SkipDir
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".go") {
			files = append(files, filepath.ToSlash(p))
		}
		return nil
	})
	return files
}

func resolveExplicitGoPaths(paths []string, baseDir string) []string {
	var files []string
	for _, p := range paths {
		target := filepath.Join(baseDir, p)
		info, err := os.Stat(target)
		if err == nil && !info.IsDir() && strings.HasSuffix(p, ".go") {
			files = append(files, filepath.ToSlash(p))
		}
		if err == nil && info.IsDir() {
			files = append(files, scanGoDirectory(target)...)
		}
	}
	return files
}

func collectGoFiles(opts FormatGoOptions, baseDir string) []string {
	if opts.IsStaged {
		return fetchStagedGoFiles(baseDir)
	}
	if len(opts.Paths) > 0 {
		return resolveExplicitGoPaths(opts.Paths, baseDir)
	}
	return scanGoDirectory(baseDir)
}

func computeFormatTotals(items []FormatGoItem) formatTotals {
	var tot formatTotals
	for _, it := range items {
		if it.IsFormatted {
			tot.formattedCount++
			tot.violationCount++
		}
	}
	return tot
}

func aggregateFormatGo(items []FormatGoItem, start time.Time) FormatGoResult {
	tot := computeFormatTotals(items)
	isClean := (tot.violationCount == 0)
	return FormatGoResult{
		TotalFiles:     len(items),
		FormattedCount: tot.formattedCount,
		ViolationCount: tot.violationCount,
		Items:          items,
		Duration:       time.Since(start),
		IsClean:        isClean,
	}
}
