package cmdfix

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/usercontext"
)

// ResolvedGitIdentity holds the resolved author name and email.
type ResolvedGitIdentity struct {
	Name         string
	Email        string
	IsConfigured bool
}

// resolveSavedGitIdentity searches for saved credentials in:
// repo/global git config, bound profile, stored profiles, or env.
func resolveSavedGitIdentity(repoPath string) ResolvedGitIdentity {
	if ident, hasIdent := resolveFromGitConfig(repoPath); hasIdent {
		return ident
	}
	if ident, hasIdent := resolveFromBoundProfile(repoPath); hasIdent {
		return ident
	}
	if ident, hasIdent := resolveFromStoredProfiles(); hasIdent {
		return ident
	}
	if ident, hasIdent := resolveFromEnv(); hasIdent {
		return ident
	}

	return ResolvedGitIdentity{IsConfigured: false}
}

func resolveFromGitConfig(repoPath string) (ResolvedGitIdentity, bool) {
	name := queryRepoGitConfig(repoPath, "user.name")
	email := queryRepoGitConfig(repoPath, "user.email")
	if name == "" {
		name = queryGlobalGitConfig("user.name")
	}
	if email == "" {
		email = queryGlobalGitConfig("user.email")
	}
	if name != "" && email != "" {
		return ResolvedGitIdentity{Name: name, Email: email, IsConfigured: true}, true
	}

	return ResolvedGitIdentity{}, false
}

func queryRepoGitConfig(repoPath, key string) string {
	out, err := exec.Command("git", "-C", repoPath, "config", key).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func queryGlobalGitConfig(key string) string {
	out, err := exec.Command("git", "config", "--global", key).Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

func resolveFromBoundProfile(repoPath string) (ResolvedGitIdentity, bool) {
	prof, isBound := usercontext.ResolveProjectUser(repoPath)
	if !isBound || prof == nil || prof.Name == "" || prof.Email == "" {
		return ResolvedGitIdentity{}, false
	}

	return ResolvedGitIdentity{Name: prof.Name, Email: prof.Email, IsConfigured: true}, true
}

func resolveFromStoredProfiles() (ResolvedGitIdentity, bool) {
	cfg, err := store.LoadGitProfiles()
	if err != nil {
		return ResolvedGitIdentity{}, false
	}
	for _, p := range cfg.Profiles {
		if p.Email != "" && p.Name != "" {
			return ResolvedGitIdentity{Name: p.Name, Email: p.Email, IsConfigured: true}, true
		}
	}

	return ResolvedGitIdentity{}, false
}

func resolveFromEnv() (ResolvedGitIdentity, bool) {
	name := os.Getenv("GIT_AUTHOR_NAME")
	email := os.Getenv("GIT_AUTHOR_EMAIL")
	if name != "" && email != "" {
		return ResolvedGitIdentity{Name: name, Email: email, IsConfigured: true}, true
	}

	return ResolvedGitIdentity{}, false
}

// buildCommitArgsWithIdentity prepares commit args using Git's saved credentials.
// If repo git config already has user.name/email, standard args are used.
// If saved credentials exist elsewhere, they are passed as -c overrides.
func buildCommitArgsWithIdentity(repoPath, message string) ([]string, bool) {
	repoEmail := queryRepoGitConfig(repoPath, "user.email")
	repoName := queryRepoGitConfig(repoPath, "user.name")
	if repoEmail != "" && repoName != "" {
		return []string{"-C", repoPath, "commit", "-m", message}, true
	}
	ident := resolveSavedGitIdentity(repoPath)
	if !ident.IsConfigured {
		printMissingIdentityWarning()

		return nil, false
	}

	return []string{
		"-C", repoPath,
		"-c", "user.name=" + ident.Name,
		"-c", "user.email=" + ident.Email,
		"commit", "-m", message,
	}, true
}

// HasConfiguredGitIdentity reports whether git resolves an author identity.
func HasConfiguredGitIdentity(repoPath string) bool {
	ident := resolveSavedGitIdentity(repoPath)

	return ident.IsConfigured
}

func printMissingIdentityWarning() {
	fmt.Fprintln(os.Stderr, "⚠ Git author identity is not configured (user.name and user.email are missing).")
	fmt.Fprintln(os.Stderr, "  Please configure your Git identity using:")
	fmt.Fprintln(os.Stderr, "    gitmap user config --global --name \"<Name>\" --email \"<Email>\"")
	fmt.Fprintln(os.Stderr, "  Or change it in Git directly:")
	fmt.Fprintln(os.Stderr, "    git config --global user.name \"<Name>\" && git config --global user.email \"<Email>\"")
}
