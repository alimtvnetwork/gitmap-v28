package setup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

const (
	scannerInitBufSize = 64 * 1024
	scannerMaxBufSize  = 1024 * 1024
)

// PathSnippetWriteResult describes the outcome of WritePathSnippet.
type PathSnippetWriteResult struct {
	Profile string // resolved profile path actually touched
	Action  string // "appended", "rewritten", or "noop"
	Snippet string // rendered snippet bytes (without trailing newline)
}

// WritePathSnippet renders the canonical snippet for the given shell
// and writes it to the user's profile. If the marker block already
// exists, it is rewritten in place (idempotent). If absent, it is
// appended after a blank line.
//
// shell: bash | zsh | fish | pwsh
// dir:   resolved deploy directory to inject into the snippet
// manager: header label, e.g. "gitmap setup" (default), "run.sh",
//
//	"installer". Determines the marker line so two managers can
//	coexist without overwriting each other's blocks.
//
// profile: explicit rc-file path. Pass "" to auto-resolve from $HOME +
//
//	shell.
//
// Spec: spec/04-generic-cli/21-post-install-shell-activation/02-snippets.md
func WritePathSnippet(shell, dir, manager, profile string) (PathSnippetWriteResult, error) {
	body, profile, err := resolveSnippetTarget(shell, dir, manager, profile)
	if err != nil {
		return PathSnippetWriteResult{}, err
	}
	existing, _ := os.ReadFile(profile)
	open := MarkerOpenFor(manager)
	if !strings.Contains(string(existing), open) {
		return appendSnippet(profile, body)
	}
	return rewriteProfileFile(profile, string(existing), open, MarkerClose(), body)
}

func resolveSnippetTarget(shell, dir, manager, profile string) (string, string, error) {
	body, err := RenderPathSnippet(shell, dir, manager)
	if err != nil {
		return "", "", err
	}
	resolvedPath, err := resolveProfilePath(shell, profile)
	if err != nil {
		return "", "", err
	}
	return body, resolvedPath, nil
}

func resolveProfilePath(shell, profile string) (string, error) {
	resolved, err := ensureProfilePath(shell, profile)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(resolved), constants.DirPermission); err != nil {
		return "", fmt.Errorf("create profile dir %s: %w", filepath.Dir(resolved), err)
	}
	return resolved, nil
}

func ensureProfilePath(shell, profile string) (string, error) {
	if len(profile) > 0 {
		return profile, nil
	}
	return defaultProfilePath(shell)
}

func rewriteProfileFile(
	profile,
	existing,
	open,
	close,
	body string,
) (PathSnippetWriteResult, error) {
	rewritten := rewriteSnippetBlock(existing, open, close, body)
	if rewritten == existing {
		return PathSnippetWriteResult{Profile: profile, Action: "noop", Snippet: body}, nil
	}
	wrErr := os.WriteFile(profile, []byte(rewritten), constants.FilePermission)
	if wrErr != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("rewrite profile %s: %w", profile, wrErr)
	}
	return PathSnippetWriteResult{Profile: profile, Action: "rewritten", Snippet: body}, nil
}

// appendSnippet adds the snippet (with leading blank line) to profile.
func appendSnippet(profile, body string) (PathSnippetWriteResult, error) {
	f, err := os.OpenFile(profile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, constants.FilePermission)
	if err != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("open profile %s: %w", profile, err)
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "\n%s\n", body); err != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("append snippet: %w", err)
	}
	return PathSnippetWriteResult{Profile: profile, Action: "appended", Snippet: body}, nil
}

// snippetScanState tracks block replacement progress during scanning.
type snippetScanState struct {
	skipping bool
	wrote    bool
}

func processSnippetLine(
	line,
	open,
	close,
	body string,
	state *snippetScanState,
	out *strings.Builder,
) {
	switch {
	case !state.skipping && line == open:
		state.skipping = true
		out.WriteString(body + "\n")
		state.wrote = true
	case state.skipping && line == close:
		state.skipping = false
	case !state.skipping:
		out.WriteString(line + "\n")
	}
}

func scanSnippetLines(content, open, close, body string) (string, bool) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, scannerInitBufSize), scannerMaxBufSize)
	var out strings.Builder
	var state snippetScanState
	for scanner.Scan() {
		processSnippetLine(scanner.Text(), open, close, body, &state, &out)
	}
	return out.String(), state.wrote
}

// rewriteSnippetBlock replaces the existing marker block with body.
// Lines outside the block (including order) are preserved exactly.
func rewriteSnippetBlock(content, open, close, body string) string {
	res, wrote := scanSnippetLines(content, open, close, body)
	if !wrote {
		return content
	}
	if strings.HasSuffix(content, "\n") {
		return res
	}
	return strings.TrimRight(res, "\n")
}

// defaultProfilePath picks the conventional rc file for the shell.
func defaultProfilePath(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	rel, err := profileRelPath(shell)
	if err != nil {
		return "", err
	}
	return filepath.Join(home, rel), nil
}

func profileRelPath(shell string) (string, error) {
	switch shell {
	case constants.PathSnippetShellBash:
		return ".bashrc", nil
	case constants.PathSnippetShellZsh:
		return ".zshrc", nil
	case constants.PathSnippetShellFish:
		return filepath.Join(".config", "fish", "config.fish"), nil
	case constants.PathSnippetShellPwsh:
		// PowerShell profile resolution is OS-specific; callers should
		// pass an explicit path on Windows. Fallback for cross-shell use.
		return filepath.Join(".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
	default:
		return "", fmt.Errorf("unknown shell %q", shell)
	}
}
