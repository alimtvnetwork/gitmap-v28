package cmdinstall

import (
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func resolveBuildVersion(version string) string {
	if version != "" {
		return version
	}

	return AntigravityDefaultVersion
}

func resolveBuildID(buildID string) string {
	if buildID != "" {
		return buildID
	}

	return AntigravityDefaultBuildID
}

func buildAntigravityURL(version, buildID, platform, artifact string) string {
	v := resolveBuildVersion(version)
	b := resolveBuildID(buildID)

	return fmt.Sprintf("%s/%s-%s/%s/%s", AntigravityBaseURL, v, b, platform, artifact)
}

func resolveDownloadArch() string {
	if runtime.GOARCH == "arm64" || runtime.GOARCH == "arm" {
		return "arm64"
	}

	return "x64"
}

func parseManifestVersionAndBuild(content string) (string, string) {
	var version, buildID string
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "version:") {
			version = strings.TrimSpace(strings.TrimPrefix(trimmed, "version:"))
			version = strings.Trim(version, "\"'")
		}

		if strings.Contains(trimmed, "antigravity-hub/") {
			parts := strings.Split(trimmed, "antigravity-hub/")
			if len(parts) > 1 {
				subParts := strings.Split(parts[1], "/")
				if len(subParts) > 0 {
					tag := subParts[0]
					dashParts := strings.Split(tag, "-")
					if len(dashParts) > 1 {
						buildID = dashParts[1]
					}
				}
			}
		}
	}

	return version, buildID
}

func fetchLiveReleaseMetadata(manifestURL string) (string, string) {
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get(manifestURL)
	if err != nil || resp.StatusCode != http.StatusOK {
		return AntigravityDefaultVersion, AntigravityDefaultBuildID
	}
	defer resp.Body.Close()

	body, errRead := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
	if errRead != nil {
		return AntigravityDefaultVersion, AntigravityDefaultBuildID
	}

	ver, bld := parseManifestVersionAndBuild(string(body))
	if ver == "" || bld == "" {
		return AntigravityDefaultVersion, AntigravityDefaultBuildID
	}

	return ver, bld
}

func getAntigravityDesktopDownloadUrl(osName string) string {
	arch := resolveDownloadArch()
	platformID, artifact, appErr := resolvePlatformArtifact(osName, arch)
	if appErr != nil {
		return ""
	}

	ver, bld := fetchLiveReleaseMetadata(AntigravityManifestURL)

	return buildAntigravityURL(ver, bld, platformID, artifact)
}

func recordAntigravityDesktopInstalled(installPath string) {
	recordToolInDatabases("antigravity", installPath, "installer")
	recordToolInDatabases("agy", installPath, "installer")
}
