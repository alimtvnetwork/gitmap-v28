package cmdinstall

import (
	"fmt"
	"runtime"
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

func getAntigravityDesktopDownloadUrl(osName string) string {
	arch := resolveDownloadArch()
	platformID, artifact, appErr := resolvePlatformArtifact(osName, arch)
	if appErr != nil {
		return ""
	}

	return buildAntigravityURL(AntigravityDefaultVersion, AntigravityDefaultBuildID, platformID, artifact)
}

func recordAntigravityDesktopInstalled(installPath string) {
	recordToolInDatabases("antigravity", installPath, "installer")
	recordToolInDatabases("agy", installPath, "installer")
}
