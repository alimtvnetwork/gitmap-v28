// Package cmdrelease — release_verify.go: `gitmap release verify <repo> <tag>`.
//
// Post-release integrity audit via `gh release view`: asset count, draft
// flag, checksums presence, and tag→commit == repo HEAD. Read-only.
//
// Wiring (lead): inside RunRelease's sub routing (`sub := strings.ToLower(args[0])`):
//  1. add `if sub == "verify" { if appErr := RunReleaseVerify(args[1:]); appErr != nil { return appErr }; return nil }`
//  2. the explicit nil check is required — *appfault.AppError is a typed nil
//     on success, and a bare `return RunReleaseVerify(...)` would surface it
//     as a non-nil error.
package cmdrelease

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// ghReleaseView is the parsed `gh release view --json` payload.
type ghReleaseView struct {
	TagName string `json:"tagName"`
	IsDraft bool   `json:"isDraft"`
	Assets  []struct {
		Name string `json:"name"`
	} `json:"assets"`
}

// RunReleaseVerify executes `gitmap release verify <repo> <tag>` — read-only.
// args is the tail after the `verify` token: <repo> <tag>.
func RunReleaseVerify(args []string) *appfault.AppError {
	if isReleaseNewCmdHelpRequested(args) {
		fmt.Print(releaseVerifyHelpText)
		cliexit.Exit(0)
	}

	if len(args) < 2 {
		return appfault.NewValidationError("missing required arguments: <repo> <tag>\n" +
			releaseVerifyUsageLine + "example: gitmap release verify gitmap-v28 v6.558.0")
	}
	repo := normalizeReleaseRepo(args[0])
	tag := strings.TrimSpace(args[1])

	if appErr := requireGhCLI(); appErr != nil {
		return appErr
	}

	view, appErr := fetchGhReleaseView(repo, tag)
	if appErr != nil {
		return appErr
	}
	tagCommit, appErr := resolveTagCommit(repo, tag)
	if appErr != nil {
		return appErr
	}
	headCommit, appErr := resolveRepoHead(repo)
	if appErr != nil {
		return appErr
	}

	return reportReleaseVerify(repo, view, tagCommit, headCommit)
}

// fetchGhReleaseView loads the release object for the tag.
func fetchGhReleaseView(repo, tag string) (*ghReleaseView, *appfault.AppError) {
	out, appErr := runGhCLI("release", "view", tag, "--repo", repo, "--json", "tagName,isDraft,assets")
	if appErr != nil {
		return nil, appErr
	}

	var view ghReleaseView
	if err := json.Unmarshal([]byte(out), &view); err != nil {
		return nil, appfault.NewValidationError("could not parse gh release view output: " + err.Error())
	}

	return &view, nil
}

// resolveTagCommit peels the tag ref to its commit SHA (annotated tags peel
// one extra level through the tag object).
func resolveTagCommit(repo, tag string) (string, *appfault.AppError) {
	type refObject struct {
		SHA  string `json:"sha"`
		Type string `json:"type"`
	}
	type tagRef struct {
		Object refObject `json:"object"`
	}

	out, appErr := runGhCLI("api", "repos/"+repo+"/git/refs/tags/"+tag, "--jq", "{sha: .object.sha, type: .object.type}")
	if appErr != nil {
		return "", appErr
	}

	var ref tagRef
	if err := json.Unmarshal([]byte(out), &ref); err != nil {
		return "", appfault.NewValidationError("could not parse tag ref: " + err.Error())
	}
	if ref.Object.Type == "tag" {
		out, appErr = runGhCLI("api", "repos/"+repo+"/git/tags/"+ref.Object.SHA, "--jq", ".object.sha")
		if appErr != nil {
			return "", appErr
		}

		var commitSHA string
		if err := json.Unmarshal([]byte(out), &commitSHA); err != nil {
			return "", appfault.NewValidationError("could not parse tag object: " + err.Error())
		}

		return commitSHA, nil
	}

	return ref.Object.SHA, nil
}

// resolveRepoHead returns the default branch HEAD commit SHA.
func resolveRepoHead(repo string) (string, *appfault.AppError) {
	branch, appErr := runGhCLI("api", "repos/"+repo, "--jq", ".default_branch")
	if appErr != nil {
		return "", appErr
	}
	head, appErr := runGhCLI("api", "repos/"+repo+"/commits/"+branch, "--jq", ".sha")
	if appErr != nil {
		return "", appErr
	}

	var sha string
	if err := json.Unmarshal([]byte(head), &sha); err != nil {
		return "", appfault.NewValidationError("could not parse HEAD commit: " + err.Error())
	}

	return sha, nil
}

// hasChecksumAsset reports whether any asset looks like a checksum file.
func hasChecksumAsset(view *ghReleaseView) bool {
	for _, asset := range view.Assets {
		lower := strings.ToLower(asset.Name)
		if strings.Contains(lower, "checksum") || strings.Contains(lower, "sha256") {
			return true
		}
	}

	return false
}

// shortSHA renders the first 8 chars of a commit SHA.
func shortSHA(sha string) string {
	if len(sha) > 8 {
		return sha[:8]
	}

	return sha
}

// reportReleaseVerify prints the audit rows; a tag/HEAD mismatch fails the
// command with a non-zero exit.
func reportReleaseVerify(repo string, view *ghReleaseView, tagCommit, headCommit string) *appfault.AppError {
	fmt.Printf("%s (%s):\n", view.TagName, repo)
	fmt.Printf("  assets:    %d\n", len(view.Assets))
	fmt.Printf("  draft:     %v\n", view.IsDraft)

	checksums := "missing"
	if hasChecksumAsset(view) {
		checksums = "present ✓"
	}
	fmt.Printf("  checksums: %s\n", checksums)

	tagShort := shortSHA(tagCommit)
	headShort := shortSHA(headCommit)
	if tagCommit != "" && tagCommit == headCommit {
		fmt.Printf("  tag→%s = HEAD ✓\n", tagShort)

		return nil
	}

	fmt.Printf("  tag→%s ≠ HEAD %s ✗\n", tagShort, headShort)

	return appfault.NewValidationError(
		fmt.Sprintf("release verify failed: tag %s points at %s, repo HEAD is %s",
			view.TagName, tagShort, headShort))
}
