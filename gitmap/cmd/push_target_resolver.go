// Package cmd — push_target_resolver.go resolves repositories for push with recursive top-level discovery.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ResolvePushDirectoryTargets resolves scan records from a work directory for push.
func ResolvePushDirectoryTargets(dirPath string) []model.ScanRecord {
	return ResolvePullDirectoryTargets(dirPath)
}
