// Package cmd — `gitmap pull-release-cd` (alias `prc`).
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runPullReleaseCD implements `gitmap pull-release-cd` (alias `prc`).
//
// Syntax: `gitmap prc <name-or-url> <version>[, <name-or-url> <version> ...]`.
// For each parsed entry it resolves the repo path (cloning URLs first),
// then spawns `gitmap pull-release <version> -y` with cwd set to that
// path. See spec/01-app/112-pull-release-cd.md.
func runPullReleaseCD(args []string) {
	checkHelp(constants.CmdPullReleaseCD, args)
	printCanonicalCmdBanner(constants.CmdPullReleaseCD, constants.CmdPullReleaseCDAlias)

	entries, err := parsePRCEntries(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ %v\n", err)
		os.Exit(2)
	}

	self, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  ✗ resolve gitmap binary: %v\n", err)
		os.Exit(1)
	}

	results := executePRCEntries(self, entries)
	printPRCSummary(results)

	for _, r := range results {
		if !r.ok {
			os.Exit(1)
		}
	}
}

// prcEntry is one `<name-or-url> <version>` pair.
type prcEntry struct {
	token   string // raw slug or URL
	version string
}

// prcResult is the per-entry execution status.
type prcResult struct {
	entry  prcEntry
	slug   string
	path   string
	ok     bool
	reason string
}

// parsePRCEntries normalises os.Args[2:] into a list of entries by
// joining on spaces, splitting on commas, and pulling out two tokens
// per segment.
func parsePRCEntries(args []string) ([]prcEntry, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("usage: gitmap prc <name-or-url> <version>[, ...]")
	}

	joined := strings.Join(args, " ")
	segments := strings.Split(joined, ",")

	var out []prcEntry
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if seg == "" {
			continue
		}
		parts := strings.Fields(seg)
		if len(parts) < 2 {
			return nil, fmt.Errorf("entry %q missing version (expected `<name-or-url> <version>`)", seg)
		}
		if len(parts) > 2 {
			return nil, fmt.Errorf("entry %q has extra tokens; only `<name-or-url> <version>` allowed", seg)
		}
		out = append(out, prcEntry{token: parts[0], version: parts[1]})
	}

	if len(out) == 0 {
		return nil, fmt.Errorf("no entries parsed")
	}

	return out, nil
}

// executePRCEntries runs each entry sequentially and collects results.
func executePRCEntries(self string, entries []prcEntry) []prcResult {
	results := make([]prcResult, 0, len(entries))
	for i, e := range entries {
		fmt.Fprintf(os.Stderr, "\n──── [%d/%d] %s @ %s ────\n", i+1, len(entries), e.token, e.version)
		results = append(results, executeOnePRCEntry(self, e))
	}

	return results
}

// executeOnePRCEntry resolves one entry's path and runs pull-release on it.
func executeOnePRCEntry(self string, e prcEntry) prcResult {
	res := prcResult{entry: e}

	slug, path, err := resolvePRCTarget(self, e.token)
	if err != nil {
		res.reason = err.Error()

		return res
	}
	res.slug = slug
	res.path = path

	if err := runPRCSubprocess(self, path, e.version); err != nil {
		res.reason = err.Error()

		return res
	}

	res.ok = true

	return res
}

// resolvePRCTarget returns (slug, absolutePath) for token. URL tokens are
// cloned first via a subprocess `gitmap clone <url>`, then re-looked-up
// by the slug derived from the URL's last path segment.
func resolvePRCTarget(self, token string) (string, string, error) {
	if isPRCURL(token) {
		if err := runPRCClone(self, token); err != nil {
			return "", "", fmt.Errorf("clone %s: %w", token, err)
		}
		slug := slugFromURL(token)

		path, err := lookupPRCPath(slug)
		if err != nil {
			return slug, "", err
		}

		return slug, path, nil
	}

	path, err := lookupPRCPath(token)
	if err != nil {
		return token, "", err
	}

	return token, path, nil
}

// isPRCURL recognises an HTTPS or SSH git URL.
func isPRCURL(token string) bool {
	return strings.Contains(token, "://") || strings.HasPrefix(token, constants.PrefixSSH)
}

// slugFromURL strips the URL down to its last path segment minus `.git`.
func slugFromURL(url string) string {
	url = strings.TrimSuffix(url, constants.ExtGit)
	if i := strings.LastIndexAny(url, "/:"); i >= 0 {
		url = url[i+1:]
	}

	return strings.ToLower(url)
}

// lookupPRCPath finds a registered repo by slug and returns its abs path.
func lookupPRCPath(slug string) (string, error) {
	db, err := openDB()
	if err != nil {
		return "", fmt.Errorf("open db: %w", err)
	}
	defer db.Close()

	repos, err := db.FindBySlug(strings.ToLower(slug))
	if err != nil {
		return "", fmt.Errorf("lookup %s: %w", slug, err)
	}
	if len(repos) == 0 {
		return "", fmt.Errorf("repo %q not found — run `gitmap scan` first", slug)
	}

	return repos[0].AbsolutePath, nil
}

// runPRCClone shells out to `gitmap clone <url>` so the cloner's
// existing scan-and-register pipeline runs unchanged.
func runPRCClone(self, url string) error {
	cmd := exec.Command(self, constants.CmdClone, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// runPRCSubprocess invokes `gitmap pull-release <version> -y` with
// cwd set to dir. Isolated so a failure cannot corrupt our own cwd
// or os.Exit the whole batch.
func runPRCSubprocess(self, dir, version string) error {
	cmd := exec.Command(self, constants.CmdReleasePull, version, "-y")
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// printPRCSummary writes a one-line-per-entry status table to stderr.
func printPRCSummary(results []prcResult) {
	fmt.Fprintln(os.Stderr, "\n──── pull-release-cd summary ────")
	for _, r := range results {
		status := "OK  "
		detail := r.path
		if !r.ok {
			status = "FAIL"
			detail = r.reason
		}
		fmt.Fprintf(os.Stderr, "  [%s] %s %s — %s\n", status, r.entry.token, r.entry.version, detail)
	}
}
