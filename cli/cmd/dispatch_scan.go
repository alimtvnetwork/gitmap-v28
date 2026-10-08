// Package cmd — dispatch_scan.go: thin scan dispatch with pre-check.
//
// Kept from cli/cmd/clihelpers.go during the spec-243 cmd split (Wave D).
// The special-repos pre-check is dispatch-level; the scan itself lives in
// cmdscan.
package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
)

func runScan(args []string) error {
	_, _ = CheckSpecialReposOnScan(resolveSpecialWorkBaseDir(), false)
	return cmdscan.RunScan(args)
}
