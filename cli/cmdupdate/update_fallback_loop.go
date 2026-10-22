package cmdupdate

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// GitHubReleaseAsset represents a downloadable file attached to a release.
type GitHubReleaseAsset struct {
	Name               string `json:"name"`
	Size               int64  `json:"size"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

// GitHubRelease represents a GitHub release payload.
type GitHubRelease struct {
	TagName    string               `json:"tag_name"`
	Name       string               `json:"name"`
	Body       string               `json:"body"`
	Draft      bool                 `json:"draft"`
	Prerelease bool                 `json:"prerelease"`
	Assets     []GitHubReleaseAsset `json:"assets"`
}

// ReleaseCandidate encapsulates the chosen target release for update.
type ReleaseCandidate struct {
	Tag           string
	Version       string
	AssetURL      string
	AssetSize     int64
	ChecksumURL   string
	FallbackDepth int
	IsCached      bool
	ReleaseNotes  string
}

// MatchPlatformAsset checks if the release asset corresponds to the target OS and architecture.
func MatchPlatformAsset(assetName, tag, platform, arch string) bool {
	lowName, lowPlatform, lowArch := strings.ToLower(assetName), strings.ToLower(platform), strings.ToLower(arch)
	switch lowPlatform {
	case "windows":
		return matchWindowsAsset(lowName, lowArch)
	case "linux":
		return matchLinuxAsset(lowName, lowArch)
	case "darwin":
		return matchDarwinAsset(lowName, lowArch)
	default:
		return false
	}
}

func matchWindowsAsset(name, arch string) bool {
	return strings.HasSuffix(name, ".zip") && (strings.Contains(name, "windows") || strings.Contains(name, "win")) && strings.Contains(name, arch)
}

func matchLinuxAsset(name, arch string) bool {
	return strings.HasSuffix(name, ".tar.gz") && strings.Contains(name, "linux") && strings.Contains(name, arch)
}

func matchDarwinAsset(name, arch string) bool {
	return strings.HasSuffix(name, ".tar.gz") && (strings.Contains(name, "darwin") || strings.Contains(name, "macos")) && strings.Contains(name, arch)
}

// FindExecutableAsset finds an asset matching the platform and arch with non-zero size.
func FindExecutableAsset(assets []GitHubReleaseAsset, tag, platform, arch string) (GitHubReleaseAsset, bool) {
	for _, a := range assets {
		if MatchPlatformAsset(a.Name, tag, platform, arch) && a.Size > 0 {
			return a, true
		}
	}
	for _, a := range assets {
		if MatchPlatformAsset(a.Name, tag, platform, arch) {
			return a, true
		}
	}
	return GitHubReleaseAsset{}, false
}

// FindChecksumAsset finds the checksums.txt asset if present.
func FindChecksumAsset(assets []GitHubReleaseAsset) string {
	for _, a := range assets {
		if strings.Contains(strings.ToLower(a.Name), "checksum") || strings.Contains(strings.ToLower(a.Name), "sha256") {
			return a.BrowserDownloadURL
		}
	}
	return ""
}

// FetchRecentReleases queries the GitHub releases API for the latest releases.
func FetchRecentReleases(slug string, limit int) ([]GitHubRelease, error) {
	if limit <= 0 {
		limit = 10
	}
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases?per_page=%d", constants.UpdateRepoOwner, slug, limit)
	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}
	req.Header.Set("User-Agent", constants.UpdaterBin)
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, errResp := (&http.Client{Timeout: 8 * time.Second}).Do(req)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, apperror.NewExecutionError(fmt.Sprintf("GitHub API returned status %d", resp.StatusCode))
	}
	var releases []GitHubRelease
	if errDec := json.NewDecoder(resp.Body).Decode(&releases); errDec != nil {
		return nil, errDec
	}
	return releases, nil
}

// ResolveUpdateTargetWithFallback resolves the target version, probing up to maxFallback tags.
func ResolveUpdateTargetWithFallback(slug string, requestedVersion string, maxFallback int, isForce bool) (*ReleaseCandidate, error) {
	platform, arch := runtime.GOOS, runtime.GOARCH
	if maxFallback <= 0 {
		maxFallback = 5
	}
	dbConn, _ := store.OpenReleaseCacheDB()
	if dbConn != nil {
		defer dbConn.Close()
	}
	if requestedVersion != "" {
		return resolveExplicitVersion(dbConn, slug, requestedVersion, platform, arch, isForce)
	}
	if cand, err := resolveFallbackFromReleases(dbConn, slug, platform, arch, maxFallback); err == nil {
		return cand, nil
	}
	return resolveFallbackFromCache(dbConn, platform, arch, maxFallback)
}

func checkExplicitCached(dbConn *sql.DB, tag, platform, arch string, isForce bool) (*ReleaseCandidate, error, bool) {
	if isForce || dbConn == nil {
		return nil, nil, false
	}
	cand, err := resolveExplicitCached(dbConn, tag, platform, arch)
	return cand, err, cand != nil || err != nil
}

func resolveExplicitVersion(dbConn *sql.DB, slug, requestedVersion, platform, arch string, isForce bool) (*ReleaseCandidate, error) {
	tag := FormatVersionTag(requestedVersion)
	if cand, err, hasCached := checkExplicitCached(dbConn, tag, platform, arch, isForce); hasCached {
		return cand, err
	}
	cand, err := fetchExplicitReleaseCandidate(slug, tag, platform, arch, dbConn)
	if err == nil && cand != nil {
		return cand, nil
	}
	msg := fmt.Sprintf("requested release %s was not found or lacks executable assets for %s/%s", tag, platform, arch)
	return nil, apperror.NewWithDetails("update.explicit.not_found", "E1204", msg, "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
}

func resolveExplicitCached(dbConn *sql.DB, tag, platform, arch string) (*ReleaseCandidate, error) {
	cached, errCache := store.GetCachedRelease(dbConn, tag, platform, arch)
	if errCache != nil || cached == nil {
		return nil, nil
	}
	if !cached.HasExecutables {
		msg := fmt.Sprintf("requested release %s has no executable assets for %s/%s", tag, platform, arch)
		return nil, apperror.NewWithDetails("update.explicit.no_exec", "E1205", msg, "cmdupdate", apperror.ErrorTypeValidation, apperror.SeverityError, nil)
	}
	return &ReleaseCandidate{
		Tag: cached.Tag, Version: cached.Version, AssetURL: cached.AssetURL,
		AssetSize: cached.AssetSize, ChecksumURL: cached.ChecksumURL, FallbackDepth: 0, IsCached: true,
	}, nil
}

func fetchExplicitReleaseCandidate(slug, tag, platform, arch string, dbConn *sql.DB) (*ReleaseCandidate, error) {
	ghRel, err := fetchTagRelease(slug, tag)
	if err != nil {
		return nil, err
	}
	asset, hasAsset := FindExecutableAsset(ghRel.Assets, tag, platform, arch)
	checksumURL := FindChecksumAsset(ghRel.Assets)
	cacheReleaseRecord(dbConn, tag, platform, arch, asset, checksumURL, hasAsset)
	if !hasAsset {
		return nil, apperror.NewNotFoundError("no executable asset")
	}
	return &ReleaseCandidate{
		Tag: tag, Version: NormalizeVersion(tag), AssetURL: asset.BrowserDownloadURL,
		AssetSize: asset.Size, ChecksumURL: checksumURL, FallbackDepth: 0, IsCached: false, ReleaseNotes: ghRel.Body,
	}, nil
}

func fetchTagRelease(slug, tag string) (*GitHubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", constants.UpdateRepoOwner, slug, tag)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", constants.UpdaterBin)
	resp, errResp := (&http.Client{Timeout: 8 * time.Second}).Do(req)
	if errResp != nil {
		return nil, errResp
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, apperror.NewExecutionError(fmt.Sprintf("bad status: %d", resp.StatusCode))
	}
	var ghRel GitHubRelease
	if errDec := json.NewDecoder(resp.Body).Decode(&ghRel); errDec != nil {
		return nil, errDec
	}
	return &ghRel, nil
}

func resolveFallbackFromReleases(dbConn *sql.DB, slug, platform, arch string, maxFallback int) (*ReleaseCandidate, error) {
	releases, errReleases := FetchRecentReleases(slug, 10)
	if errReleases != nil || len(releases) == 0 {
		return nil, apperror.NewNotFoundError("no releases found")
	}
	probeDepth := 0
	for i, rel := range releases {
		if cand, isFound := probeReleaseCandidate(dbConn, rel, platform, arch, i); isFound {
			return cand, nil
		}
		if rel.Draft {
			continue
		}
		probeDepth++
		if probeDepth >= maxFallback {
			break
		}
	}
	msg := fmt.Sprintf("no valid releases with executable assets found in the last %d releases", maxFallback)
	return nil, apperror.NewWithDetails("update.fallback.exhausted", "E1205", msg, "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
}

func probeReleaseCandidate(dbConn *sql.DB, rel GitHubRelease, platform, arch string, depth int) (*ReleaseCandidate, bool) {
	if rel.Draft {
		return nil, false
	}
	tag := resolveReleaseTag(rel)
	asset, hasAsset := FindExecutableAsset(rel.Assets, tag, platform, arch)
	checksumURL := FindChecksumAsset(rel.Assets)
	cacheReleaseRecord(dbConn, tag, platform, arch, asset, checksumURL, hasAsset)
	if !hasAsset {
		return nil, false
	}
	return &ReleaseCandidate{
		Tag: tag, Version: NormalizeVersion(tag), AssetURL: asset.BrowserDownloadURL,
		AssetSize: asset.Size, ChecksumURL: checksumURL, FallbackDepth: depth, IsCached: false, ReleaseNotes: rel.Body,
	}, true
}

func resolveReleaseTag(rel GitHubRelease) string {
	if len(rel.TagName) > 0 {
		return rel.TagName
	}
	return rel.Name
}

func cacheReleaseRecord(dbConn *sql.DB, tag, platform, arch string, asset GitHubReleaseAsset, checksumURL string, hasAsset bool) {
	if dbConn == nil {
		return
	}
	_ = store.UpsertReleaseCache(dbConn, store.ReleaseCacheRecord{
		Tag: tag, Version: NormalizeVersion(tag), HasExecutables: hasAsset,
		AssetURL: asset.BrowserDownloadURL, AssetSize: asset.Size, Platform: platform, Arch: arch, ChecksumURL: checksumURL,
	})
}

func resolveFallbackFromCache(dbConn *sql.DB, platform, arch string, maxFallback int) (*ReleaseCandidate, error) {
	if dbConn == nil {
		return nil, errResolutionFailed()
	}
	cachedList, errList := store.ListCachedReleases(dbConn, maxFallback)
	if errList != nil || len(cachedList) == 0 {
		return nil, errResolutionFailed()
	}
	for i, rec := range cachedList {
		if rec.HasExecutables && rec.Platform == platform && rec.Arch == arch {
			return &ReleaseCandidate{
				Tag: rec.Tag, Version: rec.Version, AssetURL: rec.AssetURL, AssetSize: rec.AssetSize, ChecksumURL: rec.ChecksumURL, FallbackDepth: i, IsCached: true,
			}, nil
		}
	}
	return nil, errResolutionFailed()
}

func errResolutionFailed() error {
	return apperror.NewWithDetails("update.resolution.failed", "E1205", "failed to fetch releases and no valid cached release found", "cmdupdate", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
}
