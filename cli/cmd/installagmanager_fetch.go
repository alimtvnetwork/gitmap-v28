package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"
)

type agManagerAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type agManagerRelease struct {
	TagName string           `json:"tag_name"`
	Assets  []agManagerAsset `json:"assets"`
}

func getAgManagerAssetURL() (string, string, error) {
	url, ver, err := fetchAgManagerFromAPI()
	if err == nil && url != "" {
		return url, ver, nil
	}

	url, ver, err = fetchAgManagerFromRedirect()
	if err == nil && url != "" {
		return url, ver, nil
	}

	return fetchAgManagerFromGitTags()
}

func fetchAgManagerFromAPI() (string, string, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.github.com/repos/lbjlaq/Antigravity-Manager/releases/latest", nil)
	if err != nil {
		return "", "", err
	}

	req.Header.Set("User-Agent", "Gitmap-Installer")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("github api status: %d", resp.StatusCode)
	}

	return parseAgManagerReleaseResponse(resp)
}

func parseAgManagerReleaseResponse(resp *http.Response) (string, string, error) {
	var release agManagerRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return "", "", err
	}

	ver := strings.TrimPrefix(release.TagName, "v")
	url, err := matchAgManagerAsset(release.Assets)
	if err != nil {
		return "", "", err
	}

	return url, ver, nil
}

func fetchAgManagerFromRedirect() (string, string, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Timeout: 8 * time.Second,
	}

	resp, err := client.Get("https://github.com/lbjlaq/Antigravity-Manager/releases/latest")
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close()
	loc := resp.Header.Get("Location")
	if loc == "" {
		return "", "", fmt.Errorf("no redirect location found for latest release")
	}

	return parseTagFromLocation(loc)
}

func parseTagFromLocation(loc string) (string, string, error) {
	idx := strings.LastIndex(loc, "/")
	if idx == -1 || idx+1 >= len(loc) {
		return "", "", fmt.Errorf("invalid release tag location: %s", loc)
	}

	tag := loc[idx+1:]
	ver := strings.TrimPrefix(tag, "v")
	url := constructAgManagerAssetURL(tag, ver)

	return url, ver, nil
}

func constructAgManagerAssetURL(tag, ver string) string {
	base := fmt.Sprintf("https://github.com/lbjlaq/Antigravity-Manager/releases/download/%s", tag)
	switch runtime.GOOS {
	case "windows":
		return fmt.Sprintf("%s/Antigravity.Tools_%s_x64-setup.exe", base, ver)
	case "darwin":
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/Antigravity.Tools_%s_aarch64.dmg", base, ver)
		}

		return fmt.Sprintf("%s/Antigravity.Tools_%s_x64.dmg", base, ver)
	case "linux":
		if runtime.GOARCH == "arm64" {
			return fmt.Sprintf("%s/Antigravity.Tools_%s_arm64.deb", base, ver)
		}

		return fmt.Sprintf("%s/Antigravity.Tools_%s_amd64.deb", base, ver)
	}

	return ""
}

func matchAgManagerAsset(assets []agManagerAsset) (string, error) {
	osStr, archStr := runtime.GOOS, runtime.GOARCH
	for _, asset := range assets {
		if isIgnoredAsset(asset.Name) {
			continue
		}

		if matchAssetOS(asset.Name, osStr, archStr) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no suitable asset found for %s %s", osStr, archStr)
}

func isIgnoredAsset(name string) bool {
	return strings.HasSuffix(name, ".sig") || strings.HasSuffix(name, "updater.json")
}

func matchAssetOS(name, osStr, archStr string) bool {
	n := strings.ToLower(name)
	switch osStr {
	case "windows":
		return (strings.HasSuffix(n, ".exe") || strings.HasSuffix(n, ".msi")) && matchArch(n, archStr)
	case "darwin":
		return strings.HasSuffix(n, ".dmg") && matchArch(n, archStr)
	case "linux":
		return (strings.HasSuffix(n, ".deb") || strings.HasSuffix(n, ".appimage")) && matchArch(n, archStr)
	}

	return false
}

func matchArch(n, archStr string) bool {
	if archStr == "arm64" {
		return strings.Contains(n, "aarch64") || strings.Contains(n, "arm64")
	}

	if archStr == "amd64" {
		return strings.Contains(n, "x64") || strings.Contains(n, "amd64") || strings.Contains(n, "x86_64")
	}

	return false
}
