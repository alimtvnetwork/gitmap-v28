// Package gitutil — uncommitted_classifier.go formats file status badges.
package gitutil

import "fmt"

func FormatDirtyBadges(d DirtyDiagnosis) string {
	if !d.IsDirty {
		return "clean"
	}
	return fmt.Sprintf("M:%d / U:%d / D:%d", d.ModifiedCount, d.UntrackedCount, d.DeletedCount)
}
