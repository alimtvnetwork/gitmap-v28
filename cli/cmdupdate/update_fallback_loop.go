package cmdupdate

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
}

// MatchPlatformAsset checks if the release asset corresponds to the target OS and architecture.
func MatchPlatformAsset(assetName, tag, platform, arch string) bool {
	lowName := strings.ToLower(assetName)
	lowPlatform := strings.ToLower(platform)
	lowArch := strings.ToLower(arch)

	// Target patterns
	// Windows: gitmap-<tag>-windows-<arch>.zip or gitmap-v<tag>-windows-<arch>.zip
	// Linux: gitmap-<tag>-linux-<arch>.tar.gz
	// Darwin: gitmap-<tag>-darwin-<arch>.tar.gz
	if lowPlatform == "windows" {
		if !strings.HasSuffix(lowName, ".zip") {
			return false
		}

		if !strings.Contains(lowName, "windows") && !strings.Contains(lowName, "win") {
			return false
		}

		return strings.Contains(lowName, lowArch)
	}

	if lowPlatform == "linux" {
		if !strings.HasSuffix(lowName, ".tar.gz") {
			return false
		}

		if !strings.Contains(lowName, "linux") {
			return false
		}

		return strings.Contains(lowName, lowArch)
	}

	if lowPlatform == "darwin" {
		if !strings.HasSuffix(lowName, ".tar.gz") {
			return false
		}

		if !strings.Contains(lowName, "darwin") && !strings.Contains(lowName, "macos") {
			return false
		}

		return strings.Contains(lowName, lowArch)
	}

	return false
}

// FindExecutableAsset finds an asset matching the platform and arch with non-zero size.
func FindExecutableAsset(assets []GitHubReleaseAsset, tag, platform, arch string) (GitHubReleaseAsset, bool) {
	for _, a := range assets {
		if MatchPlatformAsset(a.Name, tag, platform, arch) && a.Size > 0 {
			return a, true
		}
	}

	// Secondary check: if size is not reported by API, match name anyway
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
		low := strings.ToLower(a.Name)
		if strings.Contains(low, "checksum") || strings.Contains(low, "sha256") {
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
	client := &http.Client{Timeout: 8 * time.Second}

	req, errReq := http.NewRequest(http.MethodGet, url, nil)
	if errReq != nil {
		return nil, errReq
	}

	req.Header.Set("User-Agent", constants.UpdaterBin)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, errResp := client.Do(req)
	if errResp != nil {
		return nil, errResp
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var releases []GitHubRelease
	errDec := json.NewDecoder(resp.Body).Decode(&releases)
	if errDec != nil {
		return nil, errDec
	}

	return releases, nil
}

// ResolveUpdateTargetWithFallback resolves the target version, probing up to maxFallback tags.
func ResolveUpdateTargetWithFallback(slug string, requestedVersion string, maxFallback int, isForce bool) (*ReleaseCandidate, error) {
	platform := runtime.GOOS
	arch := runtime.GOARCH
	if maxFallback <= 0 {
		maxFallback = 5
	}

	dbConn, _ := store.OpenReleaseCacheDB()
	if dbConn != nil {
		defer dbConn.Close()
	}

	// 1. Explicit requested version
	if requestedVersion != "" {
		tag := FormatVersionTag(requestedVersion)

		// Check SQLite cache first if not forced
		if !isForce && dbConn != nil {
			cached, errCache := store.GetCachedRelease(dbConn, tag, platform, arch)
			if errCache == nil && cached != nil {
				if cached.HasExecutables {
					return &ReleaseCandidate{
						Tag:           cached.Tag,
						Version:       cached.Version,
						AssetURL:      cached.AssetURL,
						AssetSize:     cached.AssetSize,
						ChecksumURL:   cached.ChecksumURL,
						FallbackDepth: 0,
						IsCached:      true,
					}, nil
				}

				return nil, apperror.NewWithDetails(
					"update.explicit.no_exec",
					"E1205",
					fmt.Sprintf("requested release %s has no executable assets for %s/%s", tag, platform, arch),
					"cmdupdate",
					apperror.ErrorTypeValidation,
					apperror.SeverityError,
					nil,
				)
			}
		}

		// Query GitHub API for specific release tag
		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", constants.UpdateRepoOwner, slug, tag)
		client := &http.Client{Timeout: 8 * time.Second}
		req, _ := http.NewRequest(http.MethodGet, url, nil)
		if req != nil {
			req.Header.Set("User-Agent", constants.UpdaterBin)
			resp, errResp := client.Do(req)
			if errResp == nil && resp.StatusCode == http.StatusOK {
				defer resp.Body.Close()

				var ghRel GitHubRelease
				if errDec := json.NewDecoder(resp.Body).Decode(&ghRel); errDec == nil {
					asset, hasAsset := FindExecutableAsset(ghRel.Assets, tag, platform, arch)
					checksumURL := FindChecksumAsset(ghRel.Assets)

					if dbConn != nil {
						_ = store.UpsertReleaseCache(dbConn, store.ReleaseCacheRecord{
							Tag:            tag,
							Version:        NormalizeVersion(tag),
							HasExecutables: hasAsset,
							AssetURL:       asset.BrowserDownloadURL,
							AssetSize:      asset.Size,
							Platform:       platform,
							Arch:           arch,
							ChecksumURL:    checksumURL,
						})
					}

					if hasAsset {
						return &ReleaseCandidate{
							Tag:           tag,
							Version:       NormalizeVersion(tag),
							AssetURL:      asset.BrowserDownloadURL,
							AssetSize:      asset.Size,
							ChecksumURL:   checksumURL,
							FallbackDepth: 0,
							IsCached:      false,
						}, nil
					}
				}
			}
		}

		return nil, apperror.NewWithDetails(
			"update.explicit.not_found",
			"E1204",
			fmt.Sprintf("requested release %s was not found or lacks executable assets for %s/%s", tag, platform, arch),
			"cmdupdate",
			apperror.ErrorTypeValidation,
			apperror.SeverityError,
			nil,
		)
	}

	// 2. Dynamic Fallback Loop (Latest / Unspecified)
	releases, errReleases := FetchRecentReleases(slug, 10)
	if errReleases == nil && len(releases) > 0 {
		probeDepth := 0

		for i, rel := range releases {
			if rel.Draft {
				continue
			}

			tag := rel.TagName
			if tag == "" {
				tag = rel.Name
			}

			asset, hasAsset := FindExecutableAsset(rel.Assets, tag, platform, arch)
			checksumURL := FindChecksumAsset(rel.Assets)

			if dbConn != nil {
				_ = store.UpsertReleaseCache(dbConn, store.ReleaseCacheRecord{
					Tag:            tag,
					Version:        NormalizeVersion(tag),
					HasExecutables: hasAsset,
					AssetURL:       asset.BrowserDownloadURL,
					AssetSize:      asset.Size,
					Platform:       platform,
					Arch:           arch,
					ChecksumURL:    checksumURL,
				})
			}

			if hasAsset {
				return &ReleaseCandidate{
					Tag:           tag,
					Version:       NormalizeVersion(tag),
					AssetURL:      asset.BrowserDownloadURL,
					AssetSize:      asset.Size,
					ChecksumURL:   checksumURL,
					FallbackDepth: i,
					IsCached:      false,
				}, nil
			}

			fmt.Fprintf(os.Stderr, "  [warn] Release %s has no executable assets for %s/%s. Probing prior releases...\n", tag, platform, arch)
			probeDepth++
			if probeDepth >= maxFallback {
				break
			}
		}

		return nil, apperror.NewWithDetails(
			"update.fallback.exhausted",
			"E1205",
			fmt.Sprintf("no valid releases with executable assets found in the last %d releases", maxFallback),
			"cmdupdate",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
	}

	// 3. Fallback to cached releases in SQLite if network query failed
	if dbConn != nil {
		cachedList, errList := store.ListCachedReleases(dbConn, maxFallback)
		if errList == nil && len(cachedList) > 0 {
			for i, rec := range cachedList {
				if rec.HasExecutables && rec.Platform == platform && rec.Arch == arch {
					return &ReleaseCandidate{
						Tag:           rec.Tag,
						Version:       rec.Version,
						AssetURL:      rec.AssetURL,
						AssetSize:     rec.AssetSize,
						ChecksumURL:   rec.ChecksumURL,
						FallbackDepth: i,
						IsCached:      true,
					}, nil
				}
			}
		}
	}

	return nil, apperror.NewWithDetails(
		"update.resolution.failed",
		"E1205",
		"failed to fetch releases and no valid cached release found",
		"cmdupdate",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
}
