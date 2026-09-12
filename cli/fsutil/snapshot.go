// Package fsutil — snapshot.go provides snapshot selection helpers.
package fsutil

import "strings"

// PickLatestSnapshotWithTag returns the latest snapshot matching the given tag.
func PickLatestSnapshotWithTag(snapshots []string, tag string) string {
	for i := len(snapshots) - 1; i >= 0; i-- {
		if strings.HasSuffix(snapshots[i], tag) {
			return snapshots[i]
		}
	}

	return ""
}
