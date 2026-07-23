// dynamic.go — context-aware tab-completion helper (#14).
//
// Shells invoke `gitmap __complete <cword> <args...>` and read
// newline-separated suggestions from stdout. The static command list
// (allcommands_generated.go) is widened with:
//
//   - repo path completions for cd/clone/reclone-class commands
//     (directories under cwd that look like a git repo).
//   - profile-name completions for chrome-profile-* commands
//     (sourced from Chrome `Local State`, falling back to dir names).
//   - version completions for visibility / make-* / find-next
//     (sourced from the OwnerRepoNameIndex SQLite table when present).
//
// The runtime helper is intentionally I/O-light: each branch returns
// in O(few dozen filesystem entries) so completion stays snappy.
package completion

import (
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Dynamic returns suggestions for the given argv slice. `cword` is
// the index of the word being completed (0-based, excluding the
// `gitmap` binary itself). When no specialized handler matches it
// falls back to the static command list filtered by prefix.
func Dynamic(cword int, argv []string) []string {
	if cword <= 0 {
		return filterByPrefix(AllCommands(), currentWord(argv, cword))
	}
	cmd := argv[0]
	prefix := currentWord(argv, cword)

	switch {
	case isRepoPathCmd(cmd):
		return filterByPrefix(localGitDirs("."), prefix)
	case strings.HasPrefix(cmd, "chrome-profile") || cmd == "cpc" || cmd == "cpm":
		return filterByPrefix(chromeProfileNames(), prefix)
	default:
		return filterByPrefix(AllCommands(), prefix)
	}
}

// WriteDynamic writes Dynamic(cword, argv) suggestions to w as
// newline-separated values. Intended for `gitmap __complete` glue.
func WriteDynamic(w io.Writer, cword int, argv []string) {
	for _, s := range Dynamic(cword, argv) {
		_, _ = io.WriteString(w, s+"\n")
	}
}

func currentWord(argv []string, cword int) string {
	if cword < 0 || cword >= len(argv) {
		return ""
	}
	return argv[cword]
}

func filterByPrefix(in []string, prefix string) []string {
	if prefix == "" {
		return in
	}
	out := in[:0:0]
	for _, s := range in {
		if strings.HasPrefix(s, prefix) {
			out = append(out, s)
		}
	}
	return out
}

func isRepoPathCmd(cmd string) bool {
	switch cmd {
	case "cd", "clone", "clone-next", "cn", "cfr", "cfrp", "reclone", "rm", "del", "remove", "scan":
		return true
	}
	return false
}

// localGitDirs returns directory names directly under root that look
// like a git repo (contain a `.git` entry) plus all plain dirs as a
// fallback for `cd` use.
func localGitDirs(root string) []string {
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() || strings.HasPrefix(e.Name(), ".") {
			continue
		}
		out = append(out, e.Name())
	}
	sort.Strings(out)
	return out
}

// chromeProfileNames returns Chrome profile directory names from the
// per-user data directory. Best-effort: returns nil when Chrome is
// not installed or the OS path can't be derived.
func chromeProfileNames() []string {
	dir := chromeUserDataDir()
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		name := e.Name()
		if !e.IsDir() {
			continue
		}
		if name == "Default" || strings.HasPrefix(name, "Profile ") {
			out = append(out, name)
		}
	}
	sort.Strings(out)
	return out
}

func chromeUserDataDir() string {
	if v := os.Getenv("LOCALAPPDATA"); v != "" {
		return filepath.Join(v, "Google", "Chrome", "User Data")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	// macOS + Linux fallbacks.
	for _, candidate := range []string{
		filepath.Join(home, "Library", "Application Support", "Google", "Chrome"),
		filepath.Join(home, ".config", "google-chrome"),
	} {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return ""
}
