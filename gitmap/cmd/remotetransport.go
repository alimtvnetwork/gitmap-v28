package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ApplyTransportFlag rewrites the `remote.origin.url` of dir to the
// requested transport and persists it via `git remote set-url`.
// Returns (changed, oldURL, newURL, err).
//
// When neither flag is set, returns immediately with changed=false.
// When both flags are set, --ssh wins and a one-line warning is
// printed to stderr (mirrors `gitmap clone` semantics).
//
// Unrecognised origin URLs fail-open: a warning is printed but no
// error is returned, so the caller can still run git push/pull.
//
// Spec: spec/01-app/111-push-pull-transport-flags.md
func ApplyTransportFlag(dir string, useSSH, useHTTPS bool) (bool, string, string, error) {
	if !useSSH && !useHTTPS {
		return false, "", "", nil
	}
	if useSSH && useHTTPS {
		fmt.Fprintln(os.Stderr, "⚠ both --ssh and --https set; --ssh wins")
		useHTTPS = false
	}

	old, err := currentOriginURL(dir)
	if err != nil {
		return false, "", "", err
	}

	var converted string
	var ok bool
	if useSSH {
		converted, ok = ConvertURLToSSH(old)
	} else {
		converted, ok = ConvertURLToHTTPS(old)
	}
	if !ok {
		fmt.Fprintf(os.Stderr, "⚠ remote.origin.url %q is not a recognised Git URL; skipping transport rewrite\n", old)

		return false, old, old, nil
	}
	if converted == old {
		return false, old, old, nil
	}

	if err := setOriginURL(dir, converted); err != nil {
		return false, old, converted, err
	}
	fmt.Printf("→ remote.origin.url: %s → %s\n", old, converted)

	return true, old, converted, nil
}

// currentOriginURL reads `git -C dir config --get remote.origin.url`.
func currentOriginURL(dir string) (string, error) {
	out, err := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		return "", fmt.Errorf("read remote.origin.url in %s: %w (is `git remote add origin <url>` set?)", dir, err)
	}

	return strings.TrimSpace(string(out)), nil
}

// setOriginURL runs `git -C dir remote set-url origin <url>`.
func setOriginURL(dir, url string) error {
	cmd := exec.Command("git", "-C", dir, "remote", "set-url", "origin", url)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("git remote set-url origin %s: %w (%s)", url, err, strings.TrimSpace(string(out)))
	}

	return nil
}
