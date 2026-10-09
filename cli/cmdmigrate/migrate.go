// Package cmd — migrate.go handles automatic migration of legacy directories.
package cmdmigrate

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/fspath"
)

// migrateLegacyDirs moves old directories into .gitmap/ if found.
func MigrateLegacyDirs() {
	fspath.MigrateLegacyDirs()
	CleanCorruptedInstallDirsSilent()
}

func CleanCorruptedInstallDirsSilent() {
	if runtime.GOOS == "windows" {
		return
	}

	_, _ = cmddoctor.CleanCorruptedDirs(cmddoctor.CleanOptions{IsDryRun: false, IsForce: true})
}
