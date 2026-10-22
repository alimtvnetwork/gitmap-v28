package cmdpe

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
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

func appendTriageRepo(repos []triageRepo, workdir string, e os.DirEntry) []triageRepo {
	if !e.IsDir() {
		return repos
	}

	repoPath := filepath.Join(workdir, e.Name())
	if _, statErr := os.Stat(filepath.Join(repoPath, constants.DotGit)); statErr != nil {
		return repos
	}

	slug, _ := githubSlugForDir(repoPath)

	return append(repos, triageRepo{name: e.Name(), path: repoPath, slug: slug})
}

// listGitRepos enumerates directories containing a .git entry.
func listGitRepos(workdir string) []triageRepo {
	entries, readErr := os.ReadDir(workdir)
	if readErr != nil {
		return nil
	}

	var repos []triageRepo
	for _, e := range entries {
		repos = appendTriageRepo(repos, workdir, e)
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].name < repos[j].name })

	return repos
}

// githubSlugForDir resolves owner/repo from the origin remote URL.
func githubSlugForDir(dir string) (string, bool) {
	origin := gitutil.GetRemoteOriginURL(dir)
	if origin == "" {
		return "", false
	}

	return parseGitHubSlug(origin)
}

// parseGitHubSlugFromURL extracts owner/repo from a remote URL.
func parseGitHubSlugFromURL(url string) string {
	slug, isGitHub := parseGitHubSlug(url)
	if !isGitHub {
		return ""
	}

	return slug
}

// parseGitHubSlug extracts owner/repo from common remote URL shapes.
func parseGitHubSlug(raw string) (string, bool) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "", false
	}

	path := extractSlugPath(trimmed)
	if path == "" {
		return "", false
	}

	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", false
	}

	return parts[0] + "/" + parts[1], true
}

func extractSlugPath(trimmed string) string {
	lower := strings.ToLower(trimmed)
	var path string
	switch {
	case strings.HasPrefix(lower, string(constants.GitHubHTTPSPrefix)):
		path = trimmed[len(constants.GitHubHTTPSPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubHTTPPrefix)):
		path = trimmed[len(constants.GitHubHTTPPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHPrefix)):
		path = trimmed[len(constants.GitHubSSHPrefix):]
	case strings.HasPrefix(lower, string(constants.GitHubSSHSchemePrefix)):
		path = trimmed[len(constants.GitHubSSHSchemePrefix):]
	default:
		return ""
	}

	return strings.TrimSuffix(path, constants.GitDotGitSuffix)
}

// gitLastCommitAge returns a human relative age of HEAD, or "unknown".
func gitLastCommitAge(dir string) string {
	out, runErr := gitutil.RunGitInDir(dir, constants.GitLog, "-1", "--format=%ar")
	if runErr != nil {
		return "unknown"
	}

	age := strings.TrimSpace(string(out))
	if age == "" {
		return "unknown"
	}

	return age
}
