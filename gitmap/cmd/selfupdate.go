// Package cmd — `gitmap self-update` command.
//
// Probes the GitHub releases API for the newest tag, compares against
// constants.Version, and re-runs `gitmap self-install` non-interactively
// when a newer release is available. Honors --dry-run and --force.
package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// SelfUpdateOptions controls a self-update run.
type SelfUpdateOptions struct {
	DryRun bool
	Force  bool
	Stdout io.Writer
	Client *http.Client
}

// RunSelfUpdate is the entry point for `gitmap self-update`.
func RunSelfUpdate(opts SelfUpdateOptions) int {
	if opts.Client == nil {
		opts.Client = &http.Client{Timeout: 10 * time.Second}
	}
	latest, err := fetchLatestReleaseTag(opts.Client)
	if err != nil {
		fmt.Fprintf(opts.Stdout, "self-update: could not reach release API: %v\n", err)
		return 2
	}
	current := "v" + constants.Version
	fmt.Fprintf(opts.Stdout, "current: %s\nlatest:  %s\n", current, latest)
	if !opts.Force && !isNewer(latest, current) {
		fmt.Fprintln(opts.Stdout, "already on the latest release.")
		return 0
	}
	if opts.DryRun {
		fmt.Fprintf(opts.Stdout, "[dry-run] would install %s via `gitmap self-install --version %s`\n", latest, latest)
		return 0
	}
	return execSelfInstallTag(opts.Stdout, latest)
}

func fetchLatestReleaseTag(c *http.Client) (string, error) {
	const url = "https://api.github.com/repos/alimtvnetwork/gitmap-v28/releases/latest"
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("github status %d", resp.StatusCode)
	}
	var body struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", err
	}
	if body.TagName == "" {
		return "", fmt.Errorf("empty tag_name in release payload")
	}
	return body.TagName, nil
}

func execSelfInstallTag(w io.Writer, tag string) int {
	cmd := exec.Command("gitmap", "self-install", "--version", tag, "-y")
	cmd.Stdout = w
	cmd.Stderr = w
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(w, "self-install failed: %v\n", err)
		return 1
	}
	fmt.Fprintf(w, "self-update: installed %s\n", tag)
	return 0
}

// isNewer reports whether tagA is a strictly later semver tag than tagB.
// Tags must be of the form `vX.Y.Z`. Falls back to string compare for
// malformed tags so the command degrades gracefully.
func isNewer(tagA, tagB string) bool {
	a, b := strings.TrimPrefix(tagA, "v"), strings.TrimPrefix(tagB, "v")
	pa, pb := splitSemverTriple(a), splitSemverTriple(b)
	for i := 0; i < 3; i++ {
		if pa[i] != pb[i] {
			return pa[i] > pb[i]
		}
	}
	return false
}

func splitSemverTriple(s string) [3]int {
	var out [3]int
	parts := strings.SplitN(s, ".", 3)
	for i, p := range parts {
		if i >= 3 {
			break
		}
		n := 0
		for _, r := range p {
			if r < '0' || r > '9' {
				break
			}
			n = n*10 + int(r-'0')
		}
		out[i] = n
	}
	return out
}
