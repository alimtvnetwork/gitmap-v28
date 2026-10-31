// Package cmdmacro — cat_suggest.go: typo resolution, directory detection, and remediation suggestions for gitmap cat.
package cmdmacro

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/diag"
)

func buildCatDirectoryError(target string) error {
	appErr := apperror.NewValidation("cat", "E1002", fmt.Sprintf("%q is a directory, not a file", target))
	targetSlash := filepath.ToSlash(target)
	base := filepath.Base(target)

	appErr.WithSuggestions(
		diag.Suggestion{
			Command:     "gitmap tree " + targetSlash,
			Description: "List files in directory " + base,
			Category:    diag.CategoryRemediation,
			ActionType:  diag.ActionManual,
		},
		diag.Suggestion{
			Command:     "gitmap find " + targetSlash,
			Description: "Search files inside " + base,
			Category:    diag.CategoryRemediation,
			ActionType:  diag.ActionManual,
		},
	)

	return appErr
}

func buildCatFileNotFoundError(target string) error {
	msg := fmt.Sprintf("file %q not found", target)
	appErr := apperror.NewNotFound("cat", "E4004", msg)
	candidates := findCatCandidates(target)

	for _, cand := range candidates {
		appErr.WithSuggestions(diag.Suggestion{
			Command:     "gitmap cat " + cand,
			Description: "View " + filepath.Base(cand),
			Category:    diag.CategoryTypo,
			ActionType:  diag.ActionManual,
		})
	}

	return appErr
}

func findCatCandidates(target string) []string {
	clean := filepath.ToSlash(target)
	dir := filepath.Dir(clean)
	base := filepath.Base(clean)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	return filterCandidateEntries(dir, base, entries)
}

func filterCandidateEntries(dir, base string, entries []os.DirEntry) []string {
	var candidates []string
	stem := strings.TrimSuffix(base, filepath.Ext(base))

	for _, e := range entries {
		isMatch := evaluateCandidateEntry(e, base, stem)
		if isMatch {
			candidates = append(candidates, formatCandidatePath(dir, e.Name()))
		}

		if len(candidates) >= 3 {
			break
		}
	}

	return candidates
}

func evaluateCandidateEntry(e os.DirEntry, base, stem string) bool {
	if e.IsDir() {
		return false
	}

	return isCatCandidateMatch(base, stem, e.Name())
}

func formatCandidatePath(dir, name string) string {
	if dir == "." || dir == "" {
		return name
	}

	return filepath.ToSlash(filepath.Join(dir, name))
}

func isCatCandidateMatch(base, stem, name string) bool {
	if strings.EqualFold(base, name) {
		return false
	}

	nameLower := strings.ToLower(name)
	prefix := resolveCandidatePrefix(stem)

	if prefix != "" && strings.HasPrefix(nameLower, prefix) {
		return true
	}

	if cmdclone.Levenshtein(strings.ToLower(base), nameLower) <= 4 {
		return true
	}

	return false
}

func resolveCandidatePrefix(stem string) string {
	lower := strings.ToLower(stem)
	idx := strings.IndexAny(lower, "_-.")

	if idx >= 2 {
		return lower[:idx]
	}

	if len(lower) >= 3 {
		return lower[:3]
	}

	return lower
}
