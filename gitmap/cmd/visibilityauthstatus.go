// Package cmd — visibilityauthstatus.go: preflight `gh auth status` /
// `glab auth status` gate shared by `make-all-*`, `vu`, and `vr`.
// Fails fast with a Code Red message BEFORE any provider mutation so
// an unauthenticated CLI cannot leave a half-populated audit run.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §preflight.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// mustEnsureProviderAuth runs `<cli> auth status` and exits
// ExitVisAuthFailed when the CLI is not logged in. Caller must have
// already passed mustEnsureProviderCLI (LookPath gate).
func mustEnsureProviderAuth(provider string, verbose bool) {
	cli := providerCLI(provider)
	args := []string{"auth", "status"}
	if verbose {
		fmt.Fprintf(os.Stderr, constants.MsgVisVerboseExec, cli, strings.Join(args, " "))
	}
	cmd := exec.Command(cli, args...)
	out, err := cmd.CombinedOutput()
	if err == nil {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrVisAuthStatusFailedFmt, cli, err, strings.TrimSpace(string(out)))
	os.Exit(constants.ExitVisAuthFailed)
}
