// Package cmd — migrate.go handles automatic migration of legacy directories.
package cmd

import "github.com/alimtvnetwork/gitmap-v28/gitmap/localdirs"

// migrateLegacyDirs moves old directories into .gitmap/ if found.
func migrateLegacyDirs() {
	localdirs.MigrateLegacyDirs()
}
