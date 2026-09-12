package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
)

// ScanProbeOptions bundles the flags that govern the optional
// background version-probe pass scan kicks off after upserting repos.
type ScanProbeOptions = cmdscan.ScanProbeOptions

// CloneFlags holds all parsed clone-command flags and positional args.
type CloneFlags = cmdclone.CloneFlags

// isLikelyURL is a cheap prefix check used to disambiguate
// "folder name" vs "second URL".
func isLikelyURL(rawURL string) bool {
	return cmdclone.IsLikelyURL(rawURL)
}
