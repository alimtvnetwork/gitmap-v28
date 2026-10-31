// Package cmdmacro — cat_reader.go: file inspection, line extraction, and TokyoNight rendering for gitmap cat.
package cmdmacro

import (
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
)

// ExecuteCatWithOptions inspects and streams target files with range, head, tail, and query filters.
func ExecuteCatWithOptions(opts CatOptions) error {
	if opts.TargetFile == "" {
		printCatUsage()

		return apperror.NewValidationError("missing filepath for cat")
	}

	data, err := os.ReadFile(opts.TargetFile)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read file %q", opts.TargetFile))
	}

	if opts.IsInfo {
		printCatFileInfo(opts.TargetFile, data)

		return nil
	}

	if shouldShowCatImageSummary(opts) {
		printCatImageSummary(opts.TargetFile, data)

		return nil
	}

	return renderCatPayload(opts, data)
}

func printCatFileInfo(path string, data []byte) {
	lines := extractCatLines(data)
	ui := output.UI()

	fmt.Fprintf(ui, "\n  %sFile Info: %s%s\n", constants.ColorCyan, path, constants.ColorReset)
	fmt.Fprintf(ui, "  Lines:    %d\n", len(lines))
	fmt.Fprintf(ui, "  Bytes:    %d\n", len(data))
	fmt.Fprintf(ui, "  Size:     %s\n\n", formatImageFileSize(len(data)))
}

func shouldShowCatImageSummary(opts CatOptions) bool {
	if !isImageExtension(opts.TargetFile) {
		return false
	}

	hasModifiers := opts.Query != "" || opts.HasLineRange || opts.HeadCount > 0 || opts.TailCount > 0

	return !hasModifiers
}

func renderCatPayload(opts CatOptions, data []byte) error {
	lines := extractCatLines(data)
	matches := resolveCatMatches(lines, opts)

	if opts.IsRaw {
		emitRawCatMatches(matches, opts.HasLineNumbers)

		return nil
	}

	emitTokyoCatMatches(opts, matches, len(lines), len(data))

	return nil
}

func extractCatLines(data []byte) []string {
	content := string(data)
	if content == "" {
		return nil
	}

	lines := strings.Split(content, "\n")
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	return lines
}

func resolveCatMatches(lines []string, opts CatOptions) []LineMatch {
	if opts.HasLineRange {
		window := extractLinesByWindow(lines, opts)
		if opts.Query != "" {
			return filterLineMatchesByQuery(window, opts.Query, opts.Limit)
		}

		return window
	}

	if opts.Query != "" {
		return filterLinesWithWindow(lines, opts)
	}

	return extractLinesByWindow(lines, opts)
}

func filterLineMatchesByQuery(window []LineMatch, query string, limit int) []LineMatch {
	var matches []LineMatch
	tokens := extractQueryTokens(query)

	for _, m := range window {
		if lineMatchesAnyToken(m.Text, tokens) {
			matches = append(matches, m)
			if limit > 0 && len(matches) >= limit {
				break
			}
		}
	}

	return matches
}

func filterLinesWithWindow(lines []string, opts CatOptions) []LineMatch {
	matches := filterLinesByQuery(lines, opts.Query, 0)
	if opts.TailCount > 0 && len(matches) > opts.TailCount {
		return matches[len(matches)-opts.TailCount:]
	}

	if opts.HeadCount > 0 && len(matches) > opts.HeadCount {
		return matches[:opts.HeadCount]
	}

	if opts.Limit > 0 && len(matches) > opts.Limit {
		return matches[:opts.Limit]
	}

	return matches
}

func filterLinesByQuery(lines []string, query string, limit int) []LineMatch {
	var matches []LineMatch
	tokens := extractQueryTokens(query)

	for idx, line := range lines {
		if lineMatchesAnyToken(line, tokens) {
			matches = append(matches, LineMatch{Number: idx + 1, Text: line})
			if limit > 0 && len(matches) >= limit {
				break
			}
		}
	}

	return matches
}

func extractQueryTokens(query string) []string {
	if strings.Contains(query, "|") {
		parts := strings.Split(query, "|")
		var tokens []string
		for _, p := range parts {
			trimmed := strings.ToLower(strings.TrimSpace(p))
			if trimmed != "" {
				tokens = append(tokens, trimmed)
			}
		}

		return tokens
	}

	return []string{strings.ToLower(query)}
}

func lineMatchesAnyToken(line string, tokens []string) bool {
	lLower := strings.ToLower(line)
	for _, token := range tokens {
		if strings.Contains(lLower, token) {
			return true
		}
	}

	return false
}

func extractLinesByWindow(lines []string, opts CatOptions) []LineMatch {
	total := len(lines)
	start, end := resolveWindowBounds(opts, total)
	var matches []LineMatch

	for idx := start; idx <= end && idx <= total; idx++ {
		matches = append(matches, LineMatch{Number: idx, Text: lines[idx-1]})
	}

	return matches
}

func resolveWindowBounds(opts CatOptions, total int) (int, int) {
	if opts.TailCount > 0 {
		return resolveTailBounds(opts.TailCount, total)
	}

	start := opts.StartLine
	if start <= 0 {
		start = 1
	}

	end := total
	if opts.EndLine > 0 && opts.EndLine < total {
		end = opts.EndLine
	}

	if opts.HeadCount > 0 && opts.HeadCount < (end-start+1) {
		end = start + opts.HeadCount - 1
	}

	return start, end
}

func resolveTailBounds(tailCount, total int) (int, int) {
	start := total - tailCount + 1
	if start < 1 {
		start = 1
	}

	return start, total
}

func emitRawCatMatches(matches []LineMatch, hasLineNumbers bool) {
	out := output.Raw()
	for _, m := range matches {
		if hasLineNumbers {
			fmt.Fprintf(out, "%6d | %s\n", m.Number, m.Text)
			continue
		}

		fmt.Fprintln(out, m.Text)
	}
}

func emitTokyoCatMatches(opts CatOptions, matches []LineMatch, totalLines, fileSize int) {
	ui := output.UI()
	banner := buildCatHeaderBanner(opts, matches, totalLines, fileSize)

	fmt.Fprintf(ui, "\n  %s%s%s\n", constants.ColorCyan, banner, constants.ColorReset)
	emitRawCatMatches(matches, opts.HasLineNumbers)
	fmt.Fprintf(ui, "  %s--------------------------------%s\n\n", constants.ColorCyan, constants.ColorReset)
}

func buildCatHeaderBanner(opts CatOptions, matches []LineMatch, totalLines, fileSize int) string {
	if opts.Query != "" && opts.TailCount > 0 {
		return fmt.Sprintf("--- %s (last %d matches for %q across %d lines, %d bytes) ---",
			opts.TargetFile, len(matches), opts.Query, totalLines, fileSize)
	}

	if opts.Query != "" {
		return fmt.Sprintf("--- %s (%d matches for %q across %d lines, %d bytes) ---",
			opts.TargetFile, len(matches), opts.Query, totalLines, fileSize)
	}

	if opts.HasLineRange || opts.HeadCount > 0 || opts.TailCount > 0 {
		start := 1
		end := totalLines
		if len(matches) > 0 {
			start = matches[0].Number
			end = matches[len(matches)-1].Number
		}

		return fmt.Sprintf("--- %s (lines %d-%d of %d, %d bytes) ---",
			opts.TargetFile, start, end, totalLines, fileSize)
	}

	return fmt.Sprintf("--- %s (%d lines, %d bytes) ---",
		opts.TargetFile, totalLines, fileSize)
}
