package cmdpurge

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/ghtoken"
)

// ReleasePruneOptions configures GitHub release cleanup operations.
type ReleasePruneOptions struct {
	RepoDir   string   `json:"repoDir"`
	Targets   []string `json:"targets"`
	IsDryRun  bool     `json:"isDryRun"`
	IsVerbose bool     `json:"isVerbose"`
}

// ReleasePruneSummary records the counts of modified release artifacts.
type ReleasePruneSummary struct {
	AssetsDeletedCount int  `json:"assetsDeletedCount"`
	NotesUpdatedCount  int  `json:"notesUpdatedCount"`
	ReleasesScanned    int  `json:"releasesScanned"`
	IsSuccess          bool `json:"isSuccess"`
}

type releaseApiContext struct {
	client   *http.Client
	token    string
	owner    string
	repo     string
	isDryRun bool
}

type ghAssetItem struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ghReleaseItem struct {
	ID     int64         `json:"id"`
	Name   string        `json:"name"`
	Body   string        `json:"body"`
	Assets []ghAssetItem `json:"assets"`
}

// PruneReleaseAssets discovers credentials and removes release assets and mentions matching targets.
func PruneReleaseAssets(repoDir string, targets []string) (*ReleasePruneSummary, error) {
	opts := ReleasePruneOptions{
		RepoDir:  repoDir,
		Targets:  targets,
		IsDryRun: false,
	}

	return ExecuteReleasePrune(opts)
}

// ExecuteReleasePrune runs the end-to-end GitHub release purge workflow.
func ExecuteReleasePrune(opts ReleasePruneOptions) (*ReleasePruneSummary, error) {
	if len(opts.Targets) == 0 {
		return &ReleasePruneSummary{IsSuccess: true}, nil
	}

	token, err := discoverReleaseToken()
	if err != nil {
		return nil, err
	}

	owner, repoName, err := extractRepoOwnerAndName(opts.RepoDir)
	if err != nil {
		return nil, err
	}

	ctx := releaseApiContext{
		client:   &http.Client{Timeout: 30 * time.Second},
		token:    token,
		owner:    owner,
		repo:     repoName,
		isDryRun: opts.IsDryRun,
	}

	return pruneAllReleases(ctx, opts.Targets)
}

func discoverReleaseToken() (string, error) {
	tok, _, err := ghtoken.Resolve()
	if err != nil || len(strings.TrimSpace(tok)) == 0 {
		return "", apperror.NewSimple("ERR_PURGE_RELEASE_NO_TOKEN", "no github token available for release pruning")
	}

	return strings.TrimSpace(tok), nil
}

func extractRepoOwnerAndName(repoDir string) (string, string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	if repoDir != "" {
		cmd.Dir = repoDir
	}

	out, err := cmd.Output()
	if err != nil {
		return "", "", apperror.WrapSimple(err, "read remote origin url")
	}

	return parseGitHubSlug(string(out))
}

func parseGitHubSlug(rawUrl string) (string, string, error) {
	trimmed := strings.TrimSpace(rawUrl)
	trimmed = strings.TrimSuffix(trimmed, ".git")

	for _, prefix := range []string{"https://github.com/", "http://github.com/", "git@github.com:"} {
		if strings.HasPrefix(trimmed, prefix) {
			parts := strings.Split(strings.TrimPrefix(trimmed, prefix), "/")
			if len(parts) >= 2 {
				return parts[0], parts[1], nil
			}
		}
	}

	return "", "", apperror.NewSimple("ERR_PURGE_INVALID_REMOTE", "remote origin is not a github repository")
}

func pruneAllReleases(ctx releaseApiContext, targets []string) (*ReleasePruneSummary, error) {
	releases, err := fetchReleases(ctx)
	if err != nil {
		return nil, err
	}

	summary := &ReleasePruneSummary{
		ReleasesScanned: len(releases),
		IsSuccess:       true,
	}

	for _, rel := range releases {
		if err := processReleaseItem(ctx, rel, targets, summary); err != nil {
			return nil, err
		}
	}

	return summary, nil
}

func fetchReleases(ctx releaseApiContext) ([]ghReleaseItem, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=100", ctx.owner, ctx.repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, apperror.WrapSimple(err, "create releases request")
	}

	applyGitHubHeaders(req, ctx.token)
	resp, err := ctx.client.Do(req)
	if err != nil {
		return nil, apperror.WrapSimple(err, "execute releases request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, apperror.NewSimple("ERR_PURGE_RELEASE_FAILED", fmt.Sprintf("github releases request failed with status: %d", resp.StatusCode))
	}

	var items []ghReleaseItem
	if err := json.NewDecoder(resp.Body).Decode(&items); err != nil {
		return nil, apperror.WrapSimple(err, "decode releases json")
	}

	return items, nil
}

func applyGitHubHeaders(req *http.Request, token string) {
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Set("User-Agent", "GitMap-History-Purge")
}

func processReleaseItem(ctx releaseApiContext, rel ghReleaseItem, targets []string, summary *ReleasePruneSummary) error {
	deleted, err := pruneMatchingAssets(ctx, rel.Assets, targets)
	if err != nil {
		return err
	}
	summary.AssetsDeletedCount += deleted

	isUpdated, err := redactAndPatchBody(ctx, rel, targets)
	if err != nil {
		return err
	}
	if isUpdated {
		summary.NotesUpdatedCount++
	}

	return nil
}

func pruneMatchingAssets(ctx releaseApiContext, assets []ghAssetItem, targets []string) (int, error) {
	count := 0
	for _, a := range assets {
		if isAssetTargetMatch(a.Name, targets) {
			if err := deleteAssetById(ctx, a.ID); err != nil {
				return count, err
			}
			count++
		}
	}

	return count, nil
}

func isAssetTargetMatch(assetName string, targets []string) bool {
	lowerName := strings.ToLower(assetName)
	for _, t := range targets {
		base := strings.ToLower(filepath.Base(t))
		if lowerName == base || strings.Contains(lowerName, base) {
			return true
		}
	}

	return false
}

func deleteAssetById(ctx releaseApiContext, assetId int64) error {
	if ctx.isDryRun {
		return nil
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/assets/%d", ctx.owner, ctx.repo, assetId)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return apperror.WrapSimple(err, "create delete asset request")
	}

	applyGitHubHeaders(req, ctx.token)
	resp, err := ctx.client.Do(req)
	if err != nil {
		return apperror.WrapSimple(err, "execute delete asset request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 && resp.StatusCode != http.StatusNotFound {
		return apperror.NewSimple("ERR_PURGE_RELEASE_FAILED", fmt.Sprintf("delete asset %d failed with status %d", assetId, resp.StatusCode))
	}

	return nil
}

func redactAndPatchBody(ctx releaseApiContext, rel ghReleaseItem, targets []string) (bool, error) {
	newBody, hasChanged := redactReleaseNotes(rel.Body, targets)
	if !hasChanged {
		return false, nil
	}

	if ctx.isDryRun {
		return true, nil
	}

	return true, patchReleaseBody(ctx, rel.ID, newBody)
}

func redactReleaseNotes(body string, targets []string) (string, bool) {
	result := body
	hasChanged := false

	for _, t := range targets {
		cleanTarget := strings.TrimSpace(t)
		if len(cleanTarget) > 0 && strings.Contains(result, cleanTarget) {
			result = strings.ReplaceAll(result, cleanTarget, "[REDACTED_BY_GITMAP_HISTORY_PURGE]")
			hasChanged = true
		}
		base := filepath.Base(cleanTarget)
		if len(base) > 0 && strings.Contains(result, base) {
			result = strings.ReplaceAll(result, base, "[REDACTED_BY_GITMAP_HISTORY_PURGE]")
			hasChanged = true
		}
	}

	return result, hasChanged
}

func patchReleaseBody(ctx releaseApiContext, releaseId int64, newBody string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/%d", ctx.owner, ctx.repo, releaseId)
	payload, _ := json.Marshal(map[string]string{"body": newBody})

	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewReader(payload))
	if err != nil {
		return apperror.WrapSimple(err, "create patch release request")
	}

	applyGitHubHeaders(req, ctx.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := ctx.client.Do(req)
	if err != nil {
		return apperror.WrapSimple(err, "execute patch release request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return apperror.NewSimple("ERR_PURGE_RELEASE_FAILED", fmt.Sprintf("patch release %d failed with status %d", releaseId, resp.StatusCode))
	}

	return nil
}
