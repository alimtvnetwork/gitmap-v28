// Package cmdzsh provides .zshrc content deduplication and appending.
package cmdzsh

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// AppendZshrcEntries deduplicates and appends custom lines to target ~/.zshrc.
func AppendZshrcEntries(targetHome, customContent string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}

	home := ResolveTargetHome(targetHome)
	zshrcPath := filepath.Join(home, ".zshrc")

	return appendDeduplicatedContent(zshrcPath, customContent)
}

func appendDeduplicatedContent(zshrcPath, customContent string) *apperror.AppError {
	existingBytes, _ := os.ReadFile(zshrcPath)
	existingLines := strings.Split(string(existingBytes), constants.NewLineUnix)
	newLines := strings.Split(customContent, constants.NewLineUnix)
	missing := filterMissingLines(existingLines, newLines)
	if len(missing) == 0 {
		return nil
	}

	return writeAppendedLines(zshrcPath, missing)
}

func filterMissingLines(existing, candidates []string) []string {
	lookup := buildLineLookup(existing)

	return collectMissingLines(candidates, lookup)
}

func buildLineLookup(lines []string) map[string]bool {
	lookup := make(map[string]bool, len(lines))
	for _, l := range lines {
		lookup[strings.TrimSpace(l)] = true
	}

	return lookup
}

func collectMissingLines(candidates []string, lookup map[string]bool) []string {
	var missing []string
	for _, c := range candidates {
		trimmed := strings.TrimSpace(c)
		if isCandidateLineMissing(trimmed, lookup) {
			missing = append(missing, c)
		}
	}

	return missing
}

func isCandidateLineMissing(trimmed string, lookup map[string]bool) bool {
	if trimmed == "" {
		return false
	}

	_, isFound := lookup[trimmed]

	return !isFound
}

func writeAppendedLines(zshrcPath string, missing []string) *apperror.AppError {
	f, err := os.OpenFile(zshrcPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.writeAppendedLines.open")
	}

	defer f.Close()

	payload := strings.Join(missing, constants.NewLineUnix) + constants.NewLineUnix
	_, err = f.WriteString(payload)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.writeAppendedLines.write")
	}

	return nil
}

// AppendZshrcFromFile reads a source file and appends deduplicated lines to ~/.zshrc.
func AppendZshrcFromFile(targetHome, srcFilePath string, isDryRun bool) *apperror.AppError {
	data, err := os.ReadFile(srcFilePath)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.AppendZshrcFromFile")
	}

	return AppendZshrcEntries(targetHome, string(data), isDryRun)
}

func writeZshrcBytes(path string, data []byte) *apperror.AppError {
	err := os.WriteFile(path, data, 0644)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.writeZshrcBytes")
	}

	return nil
}
