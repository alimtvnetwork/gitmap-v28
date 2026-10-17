package cmdci

import (
	"bytes"
	"encoding/json"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os/exec"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

// ciRun is the subset of `gh run list --json` fields ci open needs.
type ciRun struct {
	DatabaseID int64  `json:"databaseId"`
	Conclusion string `json:"conclusion"`
	Title      string `json:"displayTitle"`
	URL        string `json:"url"`
}

// ciJob is the subset of `gh run view --json jobs` fields ci open needs.
type ciJob struct {
	DatabaseID int64  `json:"databaseId"`
	Name       string `json:"name"`
	Conclusion string `json:"conclusion"`
}

// hasCiGh reports whether the GitHub CLI is available on PATH.
func hasCiGh() bool {
	_, lookErr := exec.LookPath("gh")
	return lookErr == nil
}

// ciGhJSON runs gh and decodes stdout JSON into out.
func ciGhJSON(out any, ghArgs ...string) *appfault.AppError {
	cmd := exec.Command("gh", ghArgs...)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if runErr := cmd.Run(); runErr != nil {
		return apperror.WrapSimple(runErr, "ci.gh: "+strings.TrimSpace(stderr.String()))
	}
	if decodeErr := json.Unmarshal(stdout.Bytes(), out); decodeErr != nil {
		return apperror.WrapSimple(decodeErr, "ci.gh.decode")
	}

	return nil
}

// ciRepoSlug resolves owner/repo from an argument or the cwd origin remote.
func ciRepoSlug(arg string) (string, *appfault.AppError) {
	if strings.TrimSpace(arg) != "" {
		return resolveCiSlugToken(strings.TrimSpace(arg))
	}
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		return "", appfault.NewValidationError("ci open: pass a repo argument or run inside a git checkout with an origin remote")
	}

	return resolveCiSlugToken(strings.TrimSpace(stdout.String()))
}

// resolveCiSlugToken accepts owner/repo or a github remote URL.
func resolveCiSlugToken(token string) (string, *appfault.AppError) {
	if strings.Count(token, "/") == 1 && !strings.Contains(token, ":") && !strings.Contains(token, " ") {
		return token, nil
	}
	lower := strings.ToLower(token)
	var path string
	switch {
	case strings.HasPrefix(lower, "https://github.com/"):
		path = token[len("https://github.com/"):]
	case strings.HasPrefix(lower, "git@github.com:"):
		path = token[len("git@github.com:"):]
	default:
		return "", appfault.NewValidationError("ci open: cannot resolve a github repo from " + token)
	}
	path = strings.TrimSuffix(path, ".git")
	parts := strings.Split(path, "/")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", appfault.NewValidationError("ci open: cannot resolve a github repo from " + token)
	}

	return parts[0] + "/" + parts[1], nil
}
