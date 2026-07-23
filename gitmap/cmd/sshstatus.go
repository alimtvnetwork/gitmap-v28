// Package cmd — `gitmap ssh status` subcommand (v6.57.0).
//
// Single-screen SSH health summary: ssh-agent reachability, loaded
// identities (best-effort via `ssh-add -l`), and per-host probe of
// `ssh -T -o BatchMode=yes git@<host>`. Side-effect free; safe to
// run repeatedly. Exits 0 always — diagnostic, not gating.
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// sshStatusHosts is the default reachability probe set. Kept short
// to keep `gitmap ssh status` snappy on cold-cache networks.
var sshStatusHosts = []string{"github.com", "gitlab.com", "bitbucket.org"}

// runSSHStatus prints the SSH health summary.
func runSSHStatus(_ []string) {
	fmt.Fprint(os.Stdout, constants.MsgSSHStatusHeader)
	printSSHAgentStatus()
	printSSHLoadedKeys()
	printSSHProbe()
	fmt.Fprint(os.Stdout, constants.MsgSSHStatusFooter)
}

// printSSHAgentStatus reports whether ssh-agent is reachable.
func printSSHAgentStatus() {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		fmt.Fprint(os.Stdout, constants.MsgSSHStatusAgentMissing)

		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgSSHStatusAgentRunning, sock)
}

// printSSHLoadedKeys lists identities reported by `ssh-add -l`.
// Soft-fails on missing binary or no-identities exit codes.
func printSSHLoadedKeys() {
	out, err := exec.Command("ssh-add", "-l").CombinedOutput()
	text := strings.TrimSpace(string(out))
	if err != nil || strings.Contains(text, "no identities") {
		fmt.Fprintf(os.Stdout, constants.MsgSSHStatusKeysHeader, 0)
		fmt.Fprint(os.Stdout, constants.MsgSSHStatusKeysNone)

		return
	}
	lines := splitNonEmptyLines(text)
	fmt.Fprintf(os.Stdout, constants.MsgSSHStatusKeysHeader, len(lines))
	for _, ln := range lines {
		fmt.Fprintf(os.Stdout, constants.MsgSSHStatusKeyLine, ln)
	}
}

// printSSHProbe runs a batch-mode `ssh -T` against each host.
func printSSHProbe() {
	fmt.Fprint(os.Stdout, constants.MsgSSHStatusProbeHeader)
	for _, host := range sshStatusHosts {
		ok, detail := probeSSHHost(host)
		if ok {
			fmt.Fprintf(os.Stdout, constants.MsgSSHStatusProbeOK, host, detail)

			continue
		}
		fmt.Fprintf(os.Stdout, constants.MsgSSHStatusProbeFail, host, detail)
	}
}

// probeSSHHost runs one batch-mode probe and classifies the result.
// Returns (ok, oneLineDetail). The "ok" bucket includes GitHub's
// well-known successful-auth-with-no-shell exit code 1 message.
func probeSSHHost(host string) (bool, string) {
	cmd := exec.Command("ssh", "-T", "-o", "BatchMode=yes",
		"-o", "StrictHostKeyChecking=accept-new",
		"-o", "ConnectTimeout=5", "git@"+host)
	out, _ := cmd.CombinedOutput()
	text := strings.ReplaceAll(strings.TrimSpace(string(out)), "\n", " | ")
	lower := strings.ToLower(text)
	if strings.Contains(lower, "successfully authenticated") ||
		strings.Contains(lower, "does not provide shell access") ||
		strings.Contains(lower, "logged in as") {
		return true, firstSegment(text, 80)
	}
	if text == "" {
		return false, "no response (timeout or network blocked)"
	}

	return false, firstSegment(text, 120)
}

// splitNonEmptyLines splits on newlines and drops blank entries.
func splitNonEmptyLines(s string) []string {
	parts := strings.Split(s, "\n")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}

	return out
}

// firstSegment truncates s to max runes with an ellipsis.
func firstSegment(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max] + "…"
}
