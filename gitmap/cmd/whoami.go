package cmd

// `gitmap whoami` — diagnose git identity + push-auth mismatch in cwd.
// Answers: "who am I committing as, and which credential/SSH key will
// `git push` actually use?" Common failure mode: HTTPS remote where
// Windows Credential Manager caches a different GitHub account than
// the local user.email, causing `Permission denied to <other-user>`.

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// runWhoAmI prints the effective git identity + transport + probable
// auth principal for the current repo, and offers copy-paste fixes.
// Exits 0 always; this is a diagnostic, never a mutator.
func runWhoAmI(_ []string) {
	if !isGitRepoCWD() {
		fmt.Fprintln(os.Stderr, "✗ not a git repository (run `gitmap whoami` inside a repo)")
		exitWith(1)

		return
	}

	printWhoAmIIdentity()
	url := printWhoAmITransport()
	printWhoAmIAuth(url)
	printWhoAmISSHKeys()
	printWhoAmIFixHints(url)
}

// printWhoAmISSHKeys lists private keys in ~/.ssh so the user can
// pick the right one to bind via `gitmap ssh-bind <key>`. Public
// keys (.pub) and known_hosts/config files are excluded.
func printWhoAmISSHKeys() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}
	dir := home + string(os.PathSeparator) + ".ssh"
	entries, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("\n── SSH keys (~/.ssh) ──\n  (none — directory missing)")

		return
	}
	fmt.Println("\n── SSH keys (~/.ssh) ──")
	found := false
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if strings.HasSuffix(name, ".pub") || name == "known_hosts" || name == "config" || strings.HasPrefix(name, "known_hosts") {
			continue
		}
		fmt.Printf("  %s\n", name)
		found = true
	}
	if !found {
		fmt.Println("  (no private keys found)")
		fmt.Println("  generate one: ssh-keygen -t ed25519 -C \"you@example.com\" -f ~/.ssh/id_ed25519_<label>")
	}
}


// printWhoAmIIdentity prints local vs global user.name/email so the
// user can see exactly which identity gets stamped on new commits.
func printWhoAmIIdentity() {
	fmt.Println("── Git identity (commit authorship) ──")
	fmt.Printf("  local  name : %s\n", gitCfg("--local", "user.name"))
	fmt.Printf("  local  email: %s\n", gitCfg("--local", "user.email"))
	fmt.Printf("  global name : %s\n", gitCfg("--global", "user.name"))
	fmt.Printf("  global email: %s\n", gitCfg("--global", "user.email"))
	fmt.Printf("  effective   : %s <%s>\n", gitCfg("user.name"), gitCfg("user.email"))
}

// printWhoAmITransport prints the origin URL and detected transport
// (HTTPS vs SSH). Returns the raw URL for downstream auth analysis.
func printWhoAmITransport() string {
	url := gitCfg("--get", "remote.origin.url")
	fmt.Println("\n── Remote transport ──")
	fmt.Printf("  origin: %s\n", url)
	fmt.Printf("  kind  : %s\n", classifyTransport(url))

	return url
}

// classifyTransport returns "SSH", "HTTPS", or "other/none".
func classifyTransport(url string) string {
	if url == "" {
		return "none"
	}
	if strings.HasPrefix(url, "git@") || strings.HasPrefix(url, "ssh://") {
		return "SSH"
	}
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return "HTTPS"
	}

	return "other"
}

// printWhoAmIAuth reports which principal `git push` will actually
// authenticate as: for HTTPS the cached credential-helper entry, for
// SSH the resolved identity key.
func printWhoAmIAuth(url string) {
	fmt.Println("\n── Push authentication principal ──")
	switch classifyTransport(url) {
	case "HTTPS":
		fmt.Printf("  credential.helper : %s\n", gitCfg("credential.helper"))
		fmt.Printf("  cached user (host): %s\n", probeHTTPSCachedUser(url))
		fmt.Println("  note              : commit email is IGNORED for auth;")
		fmt.Println("                      the OS credential store decides.")
	case "SSH":
		fmt.Printf("  core.sshCommand   : %s\n", gitCfg("core.sshCommand"))
		fmt.Printf("  ssh -T github.com : %s\n", probeSSHIdentity(url))
	default:
		fmt.Println("  (no HTTPS/SSH origin — nothing to authenticate)")
	}
}

// probeHTTPSCachedUser asks git's credential helper what username it
// would hand to GitHub for this URL. Empty = no cached entry.
func probeHTTPSCachedUser(url string) string {
	in := fmt.Sprintf("url=%s\n\n", url)
	cmd := exec.Command("git", "credential", "fill")
	cmd.Stdin = strings.NewReader(in)
	out, err := cmd.Output()
	if err != nil {
		return "(none cached / helper unavailable)"
	}
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "username=") {
			return strings.TrimPrefix(line, "username=")
		}
	}

	return "(no username returned)"
}

// probeSSHIdentity runs `ssh -T git@<host>` in probe mode; GitHub
// replies with the authenticated login in the greeting stderr line.
func probeSSHIdentity(url string) string {
	host := extractSSHHost(url)
	if host == "" {
		return "(cannot parse host)"
	}
	cmd := exec.Command("ssh", "-o", "BatchMode=yes", "-T", "git@"+host)
	var buf strings.Builder
	cmd.Stderr = &buf
	_ = cmd.Run() // GitHub always exits non-zero for -T; we want stderr.
	msg := strings.TrimSpace(buf.String())
	if msg == "" {
		return "(no response — ssh key rejected or ssh missing)"
	}

	return whoamiFirstLine(msg)
}

// extractSSHHost pulls the host out of `git@host:owner/repo` or
// `ssh://git@host/owner/repo`.
func extractSSHHost(url string) string {
	if strings.HasPrefix(url, "ssh://") {
		rest := strings.TrimPrefix(url, "ssh://")
		if at := strings.Index(rest, "@"); at >= 0 {
			rest = rest[at+1:]
		}
		if slash := strings.IndexAny(rest, "/:"); slash >= 0 {
			return rest[:slash]
		}

		return rest
	}
	if at := strings.Index(url, "@"); at >= 0 {
		rest := url[at+1:]
		if colon := strings.Index(rest, ":"); colon >= 0 {
			return rest[:colon]
		}
	}

	return ""
}

// whoamiFirstLine returns everything up to the first newline in s.
func whoamiFirstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}

	return s
}

// printWhoAmIFixHints prints copy-paste PowerShell/bash fixes tuned
// to the detected transport so the user does not need to remember
// the credential-helper or per-repo sshCommand incantations.
func printWhoAmIFixHints(url string) {
	fmt.Println("\n── Fix hints ──")
	switch classifyTransport(url) {
	case "HTTPS":
		fmt.Println("  Wrong GitHub account cached? Clear it, then push again:")
		if runtime.GOOS == "windows" {
			fmt.Println("    cmdkey /delete:LegacyGeneric:target=git:https://github.com")
		} else {
			fmt.Println("    printf 'protocol=https\\nhost=github.com\\n\\n' | git credential-manager erase")
		}
		fmt.Println("  Or pin the correct user into the remote URL:")
		fmt.Println("    git remote set-url origin https://<correct-user>@github.com/<owner>/<repo>.git")
	case "SSH":
		fmt.Println("  Force this repo to use a specific key (no global change):")
		fmt.Println("    gitmap ssh-bind <key-filename>     # e.g. id_ed25519_aukgit")
		fmt.Println("    (equivalent to: git config core.sshCommand \"ssh -i ~/.ssh/<key> -F /dev/null\")")

		fmt.Println("  Or switch to HTTPS: gitmap push --https")
	}
	fmt.Println()
}

// gitCfg wraps `git config <args>` and returns trimmed stdout, or ""
// if git exits non-zero (unset key, no git, etc.).
func gitCfg(args ...string) string {
	cmd := exec.Command("git", append([]string{"config"}, args...)...)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}
