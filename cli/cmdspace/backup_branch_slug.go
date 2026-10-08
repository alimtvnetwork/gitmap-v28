package cmdspace

import (
	"fmt"
	"strings"
)

// maxBackupSlugLen caps the generated slug at 60 chars (spec 243.2).
const maxBackupSlugLen = 60

// makeBackupSlug derives a branch slug from the task string: lowercase,
// spaces/underscores to hyphens, only [a-z0-9-] kept, repeats collapsed,
// trimmed, max 60 chars. Empty result is an error asking for a better
// string.
func makeBackupSlug(task string) (string, error) {
	lowered := strings.ToLower(task)
	hyphenated := strings.NewReplacer(" ", "-", "_", "-").Replace(lowered)

	var b strings.Builder
	for _, r := range hyphenated {
		if isSlugChar(r) {
			b.WriteRune(r)
		}
	}

	slug := collapseHyphens(b.String())
	slug = strings.Trim(slug, "-")
	if len(slug) > maxBackupSlugLen {
		slug = strings.TrimRight(slug[:maxBackupSlugLen], "-")
	}

	if slug == "" {
		return "", fmt.Errorf("task %q produces an empty backup slug — describe the task with letters or numbers", task)
	}

	return slug, nil
}

// isSlugChar reports whether r survives into a backup slug.
func isSlugChar(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= '0' && r <= '9' {
		return true
	}

	return r == '-'
}

// collapseHyphens reduces every run of hyphens to a single hyphen.
func collapseHyphens(s string) string {
	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	return s
}
