// Package cmd — status_target_resolver.go resolves repositories for status with recursive top-level discovery.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ResolveStatusDirectoryTargets resolves scan records from a work directory for status.
func ResolveStatusDirectoryTargets(dirPath string) []model.ScanRecord {
	return ResolvePullDirectoryTargets(dirPath)
}
