package cmdpe

import (
	"bytes"
	"encoding/json"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// triageRepo is a local checkout that maps to a GitHub repository.
type triageRepo struct {
	name string
	path string
	slug string
}

// hasGh reports whether the GitHub CLI is available on PATH.
func hasGh() bool {
	_, lookErr := exec.LookPath("gh")
	return lookErr == nil
}

// ghJSON runs gh with --json output and decodes stdout as JSON into out.
func ghJSON(out any, ghArgs ...string) *appfault.AppError {
	cmd := exec.Command("gh", ghArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if runErr := cmd.Run(); runErr != nil {
		return apperror.WrapSimple(runErr, "pe.gh: "+strings.TrimSpace(stderr.String()))
	}
	if decodeErr := json.Unmarshal(stdout.Bytes(), out); decodeErr != nil {
		return apperror.WrapSimple(decodeErr, "pe.gh.decode")
	}

	return nil
}

// resolveWorkdir returns the repo workdir from flag/env/default.
func resolveWorkdir(flagVal string) string {
	if strings.TrimSpace(flagVal) != "" {
		return flagVal
	}
	if envVal := strings.TrimSpace(os.Getenv("GITMAP_WORKDIR")); envVal != "" {
		return envVal
	}
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return "."
	}

	return filepath.Join(home, "workspace", "repos")
}

// listGitRepos enumerates directories containing a .git entry.
func listGitRepos(workdir string) []triageRepo {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return nil
	}
	var repos []triageRepo
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		repoPath := filepath.Join(workdir, e.Name())
		if _, statErr := os.Stat(filepath.Join(repoPath, ".git")); statErr != nil {
			continue
		}
		slug, hasSlug := githubSlugForDir(repoPath)
		repos = append(repos, triageRepo{name: e.Name(), path: repoPath, slug: slug})
		_ = hasSlug
	}
	sort.Slice(repos, func(i, j int) bool { return repos[i].name < repos[j].name })

	return repos
}

// githubSlugForDir resolves owner/repo from the origin remote URL.
func githubSlugForDir(dir string) (string, bool) {
	cmd := exec.Command("git", "-C", dir, "config", "--get", "remote.origin.url")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		return "", false
	}

	return parseGitHubSlug(strings.TrimSpace(stdout.String()))
}

// parseGitHubSlug extracts owner/repo from common remote URL shapes.
func parseGitHubSlug(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", false
	}
	lower := strings.ToLower(trimmed)
	var path string
	switch {
	case strings.HasPrefix(lower, "https://github.com/"):
		path = trimmed[len("https://github.com/"):]
	case strings.HasPrefix(lower, "http://github.com/"):
		path = trimmed[len("http://github.com/"):]
	case strings.HasPrefix(lower, "git@github.com:"):
		path = trimmed[len("git@github.com:"):]
	case strings.HasPrefix(lower, "ssh://git@github.com/"):
		path = trimmed[len("ssh://git@github.com/"):]
	default:
		return "", false
	}
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0] + "/" + parts[1], true
}

// gitLastCommitAge returns a human relative age of HEAD, or "unknown".
func gitLastCommitAge(dir string) string {
	cmd := exec.Command("git", "-C", dir, "log", "-1", "--format=%ar")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		return "unknown"
	}
	age := strings.TrimSpace(stdout.String())
	if age == "" {
		return "unknown"
	}

	return age
}
