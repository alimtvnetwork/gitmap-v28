// Package cmdmacro — cat_reader.go: file inspection, line extraction, and TokyoNight rendering for gitmap cat.
package cmdmacro

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/output"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ExecuteCatWithOptions inspects and streams target files with range, head, tail, and query filters.
func ExecuteCatWithOptions(opts CatOptions) error {
	if opts.TargetFile == "" {
		printCatUsage()

		return apperror.NewValidationError("missing filepath for cat")
	}

	data, err := readCatTargetData(opts.TargetFile)
	if err != nil {
		return err
	}

	recordHeatmapCatTouch(opts)

	return dispatchCatContent(opts, data)
}

func recordHeatmapCatTouch(opts CatOptions) {
	cleanPath := filepath.ToSlash(filepath.Clean(opts.TargetFile))
	if cleanPath == "" || cleanPath == "." {
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	if rel, errRel := filepath.Rel(cwd, cleanPath); errRel == nil && !strings.HasPrefix(rel, "..") {
		cleanPath = filepath.ToSlash(rel)
	}

	insertCatHeatmapAction(cwd, cleanPath, opts)
}

func insertCatHeatmapAction(cwd, cleanPath string, opts CatOptions) {
	hdb, err := store.OpenHeatmapRootDB(cwd)
	if err != nil {
		return
	}

	defer hdb.Close()

	_, _ = hdb.LogAction(store.HeatmapActionLogRecord{
		ActionType:     "READ",
		FilePath:       cleanPath,
		StartLine:      opts.StartLine,
		EndLine:        opts.EndLine,
		ToolName:       "cat",
		ReasoningState: "EXPLORATION",
		ReasoningNotes: "inspected file via gitmap cat",
	})
}

func readCatTargetData(target string) ([]byte, error) {
	stat, statErr := os.Stat(target)
	if statErr == nil && stat.IsDir() {
		return nil, buildCatDirectoryError(target)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		return nil, handleCatReadError(target, err)
	}

	return data, nil
}

func handleCatReadError(target string, err error) error {
	if os.IsNotExist(err) {
		return buildCatFileNotFoundError(target)
	}

	return apperror.WrapSimple(err, fmt.Sprintf("read file %q", target))
}

func dispatchCatContent(opts CatOptions, data []byte) error {
	if opts.IsInfo {
		printCatFileInfo(opts.TargetFile, data)

		return nil
	}

	if isCatImageSummaryShown(opts) {
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

// isCatImageSummaryShown replaces shouldShowCatImageSummary (31-Oct-26):
// the boolean-guidelines linter bans the should* prefix for bool funcs.
func isCatImageSummaryShown(opts CatOptions) bool {
	if !isImageExtension(opts.TargetFile) {
		return false
	}

	hasModifiers := opts.Query != "" || opts.HasLineRange || opts.HeadCount > 0 || opts.TailCount > 0

	return !hasModifiers
}

func renderCatPayload(opts CatOptions, data []byte) error {
	lines := extractCatLines(data)
	matches := resolveCatMatches(lines, opts)

	if opts.JsonField != "" {
		emitFieldExtractedMatches(matches, opts.JsonField)

		return nil
	}

	if opts.IsRaw {
		emitRawCatMatches(matches, opts.HasLineNumbers)

		return nil
	}

	emitTokyoCatMatches(opts, matches, len(lines), len(data))

	return nil
}

func emitFieldExtractedMatches(matches []LineMatch, field string) {
	out := output.Raw()
	for _, m := range matches {
		var obj map[string]interface{}
		// WHY guard-continue: the nested unmarshal + field-lookup
		// ifs tripped the nested-if linter; flat guards do not.
		if err := json.Unmarshal([]byte(m.Text), &obj); err != nil {
			fmt.Fprintln(out, m.Text)
			continue
		}
		val, ok := obj[field]
		if !ok {
			fmt.Fprintln(out, m.Text)
			continue
		}
		switch v := val.(type) {
		case string:
			fmt.Fprintln(out, v)
		default:
			b, _ := json.MarshalIndent(v, "", "  ")
			fmt.Fprintln(out, string(b))
		}
	}
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

// catBannerWindow returns the first/last match numbers, or the full
// range when there are no matches. Single-level guard only.
func catBannerWindow(matches []LineMatch, totalLines int) (int, int) {
	if len(matches) == 0 {
		return 1, totalLines
	}

	return matches[0].Number, matches[len(matches)-1].Number
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
		// WHY helper: the inner len(matches) check inside this
		// branch tripped the nested-if linter.
		start, end := catBannerWindow(matches, totalLines)

		return fmt.Sprintf("--- %s (lines %d-%d of %d, %d bytes) ---",
			opts.TargetFile, start, end, totalLines, fileSize)
	}

	return fmt.Sprintf("--- %s (%d lines, %d bytes) ---",
		opts.TargetFile, totalLines, fileSize)
}
