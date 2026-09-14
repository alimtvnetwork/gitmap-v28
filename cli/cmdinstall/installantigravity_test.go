package cmdinstall

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func assertPlatformArtifactAndURL(t *testing.T, osName, arch, expPlatform, expArtifact string) {
	platformID, artifact, appErr := resolvePlatformArtifact(osName, arch)
	if appErr != nil {
		t.Fatalf("unexpected error resolving %s/%s: %v", osName, arch, appErr)
	}

	if platformID != expPlatform || artifact != expArtifact {
		t.Fatalf("expected (%s, %s), got (%s, %s)", expPlatform, expArtifact, platformID, artifact)
	}

	url := buildAntigravityURL("2.13.0", "6362815968182272", platformID, artifact)
	if !strings.Contains(url, expPlatform+"/"+expArtifact) {
		t.Fatalf("expected URL to contain %s/%s, got %s", expPlatform, expArtifact, url)
	}
}

func TestAntigravityWindowsUrls(t *testing.T) {
	assertPlatformArtifactAndURL(t, "windows", "x64", "windows-x64", "Antigravity-x64.exe")
	assertPlatformArtifactAndURL(t, "windows", "arm64", "windows-arm", "Antigravity-arm64.exe")
}

func TestAntigravityDarwinUrls(t *testing.T) {
	assertPlatformArtifactAndURL(t, "darwin", "arm64", "darwin-arm", "Antigravity.dmg")
	assertPlatformArtifactAndURL(t, "darwin", "x64", "darwin-x64", "Antigravity.dmg")
}

func TestAntigravityLinuxUrls(t *testing.T) {
	assertPlatformArtifactAndURL(t, "linux", "x64", "linux-x64", "Antigravity.tar.gz")
	assertPlatformArtifactAndURL(t, "linux", "arm64", "linux-arm", "Antigravity.tar.gz")
}

func TestAntigravityWindowsDownloadUrl(t *testing.T) {
	winUrl := getAntigravityDesktopDownloadUrl("windows")
	if !strings.Contains(winUrl, "Antigravity-") || !strings.HasSuffix(winUrl, ".exe") {
		t.Fatalf("expected windows URL pointing to Antigravity executable, got %s", winUrl)
	}

	if !strings.Contains(winUrl, "antigravity-public") {
		t.Fatalf("expected windows URL from antigravity-public, got %s", winUrl)
	}
}

func TestAntigravityLinuxDownloadUrl(t *testing.T) {
	linuxUrl := getAntigravityDesktopDownloadUrl("linux")
	if !strings.Contains(linuxUrl, "Antigravity.tar.gz") {
		t.Fatalf("expected linux URL pointing to Antigravity.tar.gz, got %s", linuxUrl)
	}

	if !strings.Contains(linuxUrl, "antigravity-public") {
		t.Fatalf("expected linux URL from antigravity-public, got %s", linuxUrl)
	}
}

func assertAnnouncementContains(t *testing.T, output, snippet string) {
	if !strings.Contains(output, snippet) {
		t.Errorf("announcement missing expected snippet: %q", snippet)
	}
}

func buildTestPlatformInfo() AntigravityPlatformInfo {
	return AntigravityPlatformInfo{
		OS:           "linux",
		Arch:         "x64",
		DistroName:   "Ubuntu 24.04 LTS",
		PlatformID:   "linux-x64",
		ArtifactName: "Antigravity.tar.gz",
		DownloadURL:  "https://storage.googleapis.com/antigravity-public/test.tar.gz",
		InstallDir:   "/opt/antigravity",
		BinDir:       "/usr/local/bin",
	}
}

func TestCategoricalAnnouncementFormatting(t *testing.T) {
	info := buildTestPlatformInfo()
	output := formatCategoricalAnnouncement(info, "2.13.0")

	assertAnnouncementContains(t, output, "[STEP ] Antigravity IDE installer - version 2.13.0")
	assertAnnouncementContains(t, output, "[INFO ] Detected Platform: Linux (Ubuntu 24.04 LTS) [arch: x64]")
	assertAnnouncementContains(t, output, "[INFO ] Target Platform: linux-x64 | Artifact: Antigravity.tar.gz")
	assertAnnouncementContains(t, output, "[INFO ] Download URL: https://storage.googleapis.com/antigravity-public/test.tar.gz")
}

func TestArchPrereqErrorCodeMapping(t *testing.T) {
	_, err32 := normalizeArch("386")
	if err32 == nil || err32.Code != ErrUnsupportedPlatform {
		t.Fatalf("expected ErrUnsupportedPlatform for 32-bit arch, got %v", err32)
	}

	_, errUnknown := normalizeArch("mips")
	if errUnknown == nil || errUnknown.Code != ErrUnsupportedPlatform {
		t.Fatalf("expected ErrUnsupportedPlatform for unknown arch, got %v", errUnknown)
	}
}

func TestUnsupportedPlatformErrorCodeMapping(t *testing.T) {
	_, _, errOS := resolvePlatformArtifact("solaris", "x64")
	if errOS == nil || errOS.Code != ErrUnsupportedPlatform {
		t.Fatalf("expected ErrUnsupportedPlatform for solaris, got %v", errOS)
	}
}

func TestVersionComparisonLogic(t *testing.T) {
	if isVersionAtLeast("2.27", "2.28") {
		t.Errorf("expected 2.27 to be less than 2.28")
	}

	if !isVersionAtLeast("2.28", "2.28") {
		t.Errorf("expected 2.28 to be at least 2.28")
	}

	if !isVersionAtLeast("2.35", "2.28") {
		t.Errorf("expected 2.35 to be at least 2.28")
	}
}

func TestGlibcAndWindowsBuildParsing(t *testing.T) {
	glibc := extractGlibcVersion("ldd (Ubuntu GLIBC 2.35-0ubuntu3.8) 2.35")
	if glibc != "2.35" {
		t.Errorf("expected glibc 2.35, got %q", glibc)
	}

	build := extractWindowsBuildNumber("Microsoft Windows [Version 10.0.19045.3803]")
	if build != 19045 {
		t.Errorf("expected windows build 19045, got %d", build)
	}
}

func createDownloadFailedError(url string) *apperror.AppError {
	return apperror.NewWithDetails(
		"downloadFileWithRetry",
		ErrDownloadFailed,
		fmt.Sprintf("failed to download from %s after 3 attempts", url),
		"installer",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		map[string]any{"url": url, "attempts": 3},
	)
}

func TestDownloadFailedErrorCodeMapping(t *testing.T) {
	appErr := createDownloadFailedError("https://storage.googleapis.com/test.zip")
	if appErr == nil || appErr.Code != ErrDownloadFailed {
		t.Fatalf("expected ErrDownloadFailed, got %v", appErr)
	}
}

func TestDefaultInstallDirResolutions(t *testing.T) {
	if dir := resolveDefaultInstallDir("linux", ""); dir != "/opt/antigravity" {
		t.Errorf("expected /opt/antigravity for linux, got %s", dir)
	}

	if dir := resolveDefaultInstallDir("darwin", ""); dir != "/Applications" {
		t.Errorf("expected /Applications for darwin, got %s", dir)
	}

	if dir := resolveDefaultInstallDir("linux", "/custom/dir"); dir != "/custom/dir" {
		t.Errorf("expected /custom/dir override, got %s", dir)
	}
}

func TestDefaultBinDirResolutions(t *testing.T) {
	if bin := resolveDefaultBinDir("linux"); bin != "/usr/local/bin" {
		t.Errorf("expected /usr/local/bin for linux, got %s", bin)
	}

	if bin := resolveDefaultBinDir("darwin"); bin != "/usr/local/bin" {
		t.Errorf("expected /usr/local/bin for darwin, got %s", bin)
	}

	winBin := resolveDefaultBinDir("windows")
	if !strings.HasSuffix(winBin, filepath.Join("agy", "bin")) {
		t.Errorf("expected windows bin dir ending with agy/bin, got %s", winBin)
	}
}

func assertToolAlias(t *testing.T, alias, expected string) {
	actual := resolveToolAlias(alias)
	if actual != expected {
		t.Errorf("expected alias %q to resolve to %q, got %q", alias, expected, actual)
	}
}

func TestAntigravityIdeAliases(t *testing.T) {
	assertToolAlias(t, "antigravity", constants.ToolAntigravity)
	assertToolAlias(t, "antigravity-ide", constants.ToolAntigravity)
	assertToolAlias(t, "antigravity-desktop", constants.ToolAntigravity)
	assertToolAlias(t, "ide", constants.ToolAntigravity)
}

func TestAntigravityCliAliases(t *testing.T) {
	assertToolAlias(t, "ag", constants.ToolAgy)
	assertToolAlias(t, "agy", constants.ToolAgy)
	assertToolAlias(t, "antigravity-cli", constants.ToolAgy)
	assertToolAlias(t, "agy-cli", constants.ToolAgy)
	assertToolAlias(t, "ag-cli", constants.ToolAgy)
}

func TestAntigravityManagerAliases(t *testing.T) {
	assertToolAlias(t, "agm", constants.ToolAgManager)
	assertToolAlias(t, "ag-tools", constants.ToolAgManager)
	assertToolAlias(t, "antigravity-tools", constants.ToolAgManager)
}

func TestAntigravityBinaryNames(t *testing.T) {
	if bin := toolBinaryName(constants.ToolAntigravity); bin != "antigravity" {
		t.Errorf("expected ToolAntigravity binary name 'antigravity', got %q", bin)
	}

	if bin := toolBinaryName(constants.ToolAgy); bin != "agy" {
		t.Errorf("expected ToolAgy binary name 'agy', got %q", bin)
	}
}

func TestAntigravityProbeMap(t *testing.T) {
	cfgAntigravity, hasAntigravity := toolProbeMap[constants.ToolAntigravity]
	if !hasAntigravity || len(cfgAntigravity.bins) == 0 {
		t.Fatalf("expected ToolAntigravity in toolProbeMap")
	}

	cfgAgy, hasAgy := toolProbeMap[constants.ToolAgy]
	if !hasAgy || len(cfgAgy.bins) == 0 {
		t.Fatalf("expected ToolAgy in toolProbeMap")
	}
}

func TestAgyInstallDefaultTarget(t *testing.T) {
	opts := installOptions{DryRun: true}
	if err := dispatchAgyInstallTarget("cli", opts); err != nil {
		t.Fatalf("expected nil error for cli dry-run, got %v", err)
	}

	agyInstallDryRun = true
	defer func() { agyInstallDryRun = false }()

	if errDefault := runAgyInstallCmd(agyInstallCmd, []string{}); errDefault != nil {
		t.Fatalf("expected nil error for default agy install, got %v", errDefault)
	}
}

func TestAgyDesktopFinderIsolation(t *testing.T) {
	candidates := getAntigravityAppPaths()
	if len(candidates) == 0 {
		t.Fatalf("expected candidate paths for antigravity desktop app")
	}

	for _, path := range candidates {
		if strings.HasSuffix(path, "agy") || strings.HasSuffix(path, "agy.exe") {
			t.Errorf("candidate path %s should not be agy binary", path)
		}
	}
}
