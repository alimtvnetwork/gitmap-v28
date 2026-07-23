package setup

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	body, err := RenderPathSnippet(shell, dir, manager)
	if err != nil {
		return PathSnippetWriteResult{}, err
	}
	if len(profile) == 0 {
		profile, err = defaultProfilePath(shell)
		if err != nil {
			return PathSnippetWriteResult{}, err
		}
	}

	if mkErr := os.MkdirAll(filepath.Dir(profile), 0o755); mkErr != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("create profile dir %s: %w", filepath.Dir(profile), mkErr)
	}

	existing, _ := os.ReadFile(profile)
	open := MarkerOpenFor(manager)
	close := MarkerClose()

	if strings.Contains(string(existing), open) {
		rewritten := rewriteSnippetBlock(string(existing), open, close, body)
		if rewritten == string(existing) {
			return PathSnippetWriteResult{Profile: profile, Action: "noop", Snippet: body}, nil
		}
		if wrErr := os.WriteFile(profile, []byte(rewritten), 0o644); wrErr != nil {
			return PathSnippetWriteResult{}, fmt.Errorf("rewrite profile %s: %w", profile, wrErr)
		}
		return PathSnippetWriteResult{Profile: profile, Action: "rewritten", Snippet: body}, nil
	}

	return appendSnippet(profile, body)
}

// appendSnippet adds the snippet (with leading blank line) to profile.
func appendSnippet(profile, body string) (PathSnippetWriteResult, error) {
	f, err := os.OpenFile(profile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("open profile %s: %w", profile, err)
	}
	defer f.Close()
	if _, err = fmt.Fprintf(f, "\n%s\n", body); err != nil {
		return PathSnippetWriteResult{}, fmt.Errorf("append snippet: %w", err)
	}

	return PathSnippetWriteResult{Profile: profile, Action: "appended", Snippet: body}, nil
}

// rewriteSnippetBlock replaces the existing marker block with body.
// Lines outside the block (including order) are preserved exactly.
func rewriteSnippetBlock(content, open, close, body string) string {
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var out strings.Builder
	skip := false
	wrote := false
	for scanner.Scan() {
		line := scanner.Text()
		switch {
		case !skip && line == open:
			skip = true
			out.WriteString(body)
			out.WriteString("\n")
			wrote = true
		case skip && line == close:
			skip = false
		case !skip:
			out.WriteString(line)
			out.WriteString("\n")
		}
	}
	if !wrote {
		return content
	}
	// Preserve trailing-newline state: original ends with newline iff result should.
	if !strings.HasSuffix(content, "\n") {
		return strings.TrimRight(out.String(), "\n")
	}

	return out.String()
}

// defaultProfilePath picks the conventional rc file for the shell.
func defaultProfilePath(shell string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home dir: %w", err)
	}
	switch shell {
	case "bash":
		return filepath.Join(home, ".bashrc"), nil
	case "zsh":
		return filepath.Join(home, ".zshrc"), nil
	case "fish":
		return filepath.Join(home, ".config", "fish", "config.fish"), nil
	case "pwsh":
		// PowerShell profile resolution is OS-specific; callers should
		// pass an explicit path on Windows. Fallback for cross-shell use.
		return filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"), nil
	}

	return "", fmt.Errorf("unknown shell %q", shell)
}
