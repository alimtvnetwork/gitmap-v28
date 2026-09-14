package cmdinstall

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func detectHostOS() string {
	return runtime.GOOS
}

func is32BitArch(arch string) bool {
	return arch == "386" || arch == "arm" || arch == "x86"
}

func is64BitX86(arch string) bool {
	return arch == "amd64" || arch == "x64" || arch == "x86_64"
}

func is64BitArm(arch string) bool {
	return arch == "arm64" || arch == "aarch64"
}

func createUnsupported32BitError(arch string) *apperror.AppError {
	return apperror.NewWithDetails(
		"detectHostArch",
		ErrUnsupportedPlatform,
		"32-bit architectures are not supported by Antigravity",
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"arch": arch},
	)
}

func createUnknownArchError(arch string) *apperror.AppError {
	return apperror.NewWithDetails(
		"detectHostArch",
		ErrUnsupportedPlatform,
		fmt.Sprintf("unsupported CPU architecture: %s", arch),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"arch": arch},
	)
}

func normalizeArch(arch string) (string, *apperror.AppError) {
	if is32BitArch(arch) {
		return "", createUnsupported32BitError(arch)
	}
	if is64BitX86(arch) {
		return "x64", nil
	}
	if is64BitArm(arch) {
		return "arm64", nil
	}
	return "", createUnknownArchError(arch)
}

func detectHostArch() (string, *apperror.AppError) {
	return normalizeArch(runtime.GOARCH)
}

func readOSReleaseFile() (string, bool) {
	paths := []string{"/etc/os-release", "/usr/lib/os-release"}
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err == nil {
			return string(data), true
		}
	}
	return "", false
}

func parseDistroPrettyName(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "PRETTY_NAME=") {
			val := strings.TrimPrefix(trimmed, "PRETTY_NAME=")
			return strings.Trim(val, "\"")
		}
	}
	return ""
}

func parseDistroFallbackName(content string) string {
	for _, line := range strings.Split(content, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "NAME=") {
			val := strings.TrimPrefix(trimmed, "NAME=")
			return strings.Trim(val, "\"")
		}
	}
	return "Linux"
}

func detectAntigravityLinuxDistro() string {
	content, hasRelease := readOSReleaseFile()
	if !hasRelease {
		return "Linux"
	}
	pretty := parseDistroPrettyName(content)
	if pretty != "" {
		return pretty
	}
	return parseDistroFallbackName(content)
}

func resolveDistroOrDetails(osName string) string {
	if osName == "linux" {
		return detectAntigravityLinuxDistro()
	}
	if osName == "windows" {
		return "Windows"
	}
	return "macOS"
}

func formatOSDisplayName(osName string) string {
	if osName == "darwin" || osName == "macos" {
		return "macOS"
	}
	if osName == "windows" {
		return "Windows"
	}
	return "Linux"
}

func formatCategoricalAnnouncement(info AntigravityPlatformInfo, version string) string {
	return fmt.Sprintf(
		"[STEP ] Antigravity IDE installer - version %s\n"+
			"[INFO ] Detected Platform: %s (%s) [arch: %s]\n"+
			"[INFO ] Target Platform: %s | Artifact: %s\n"+
			"[INFO ] Download URL: %s\n",
		version,
		formatOSDisplayName(info.OS),
		info.DistroName,
		info.Arch,
		info.PlatformID,
		info.ArtifactName,
		info.DownloadURL,
	)
}

func announcePlatform(info AntigravityPlatformInfo, version string) {
	fmt.Print(formatCategoricalAnnouncement(info, version))
}

func matchWindowsArtifact(arch string) (string, string, bool) {
	if arch == "arm64" {
		return "windows-arm", "Antigravity-arm64.exe", true
	}
	if arch == "x64" {
		return "windows-x64", "Antigravity-x64.exe", true
	}
	return "", "", false
}

func matchDarwinArtifact(arch string) (string, string, bool) {
	if arch == "arm64" {
		return "darwin-arm", "Antigravity.dmg", true
	}
	if arch == "x64" {
		return "darwin-x64", "Antigravity.dmg", true
	}
	return "", "", false
}

func matchLinuxArtifact(arch string) (string, string, bool) {
	if arch == "arm64" {
		return "linux-arm", "Antigravity.tar.gz", true
	}
	if arch == "x64" {
		return "linux-x64", "Antigravity.tar.gz", true
	}
	return "", "", false
}

func matchPlatformArtifact(osName, arch string) (string, string, bool) {
	if osName == "windows" {
		return matchWindowsArtifact(arch)
	}
	if osName == "darwin" || osName == "macos" {
		return matchDarwinArtifact(arch)
	}
	if osName == "linux" {
		return matchLinuxArtifact(arch)
	}
	return "", "", false
}

func createUnsupportedPlatformError(osName, arch string) *apperror.AppError {
	return apperror.NewWithDetails(
		"resolvePlatformArtifact",
		ErrUnsupportedPlatform,
		fmt.Sprintf("no build artifact available for %s/%s", osName, arch),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"os": osName, "arch": arch},
	)
}

func resolvePlatformArtifact(osName, arch string) (string, string, *apperror.AppError) {
	platformID, artifact, isResolved := matchPlatformArtifact(osName, arch)
	if !isResolved {
		return "", "", createUnsupportedPlatformError(osName, arch)
	}
	return platformID, artifact, nil
}

func resolveWindowsAppDataDir() string {
	local := os.Getenv("LOCALAPPDATA")
	if local != "" {
		return local
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, "AppData", "Local")
}

func resolveDefaultInstallDir(osName, prefix string) string {
	if prefix != "" {
		return prefix
	}
	if osName == "windows" {
		return filepath.Join(resolveWindowsAppDataDir(), "Programs", "Antigravity")
	}
	if osName == "darwin" || osName == "macos" {
		return "/Applications"
	}
	return "/opt/antigravity"
}

func resolveDefaultBinDir(osName string) string {
	if osName == "windows" {
		return filepath.Join(resolveWindowsAppDataDir(), "agy", "bin")
	}
	return "/usr/local/bin"
}

func assemblePlatformInfo(osName, arch, platformID, artifact, version, buildID, prefix string) AntigravityPlatformInfo {
	return AntigravityPlatformInfo{
		OS:           osName,
		Arch:         arch,
		DistroName:   resolveDistroOrDetails(osName),
		PlatformID:   platformID,
		ArtifactName: artifact,
		DownloadURL:  buildAntigravityURL(version, buildID, platformID, artifact),
		InstallDir:   resolveDefaultInstallDir(osName, prefix),
		BinDir:       resolveDefaultBinDir(osName),
	}
}

func detectAntigravityPlatform(version, buildID, prefix string) (AntigravityPlatformInfo, *apperror.AppError) {
	hostOS := detectHostOS()
	arch, errArch := detectHostArch()
	if errArch != nil {
		return AntigravityPlatformInfo{}, errArch
	}
	platformID, artifact, errArtifact := resolvePlatformArtifact(hostOS, arch)
	if errArtifact != nil {
		return AntigravityPlatformInfo{}, errArtifact
	}
	return assemblePlatformInfo(hostOS, arch, platformID, artifact, version, buildID, prefix), nil
}
