// Package scripts embeds the canonical install / uninstall shell scripts
// so that subcommands like `gitmap self-install` can run them offline.
//
// Spec: spec/01-app/90-self-install-uninstall.md
package scripts

import (
	"embed"
	"io/fs"
)

//go:embed install.ps1 install.sh uninstall.ps1 release-version.ps1 release-version.sh
var files embed.FS

// FS returns the embedded read-only filesystem rooted at gitmap/scripts/.
// Callers should use fs.ReadFile to extract a script's bytes.
func FS() fs.FS {
	return files
}

// ReadFile is a convenience wrapper around fs.ReadFile against the
// embedded scripts directory.
func ReadFile(name string) ([]byte, error) {
	return fs.ReadFile(files, name)
}
