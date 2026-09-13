package cmdinstall

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func parseVersionPart(part string) int {
	var val int
	for _, ch := range part {
		if ch >= '0' && ch <= '9' {
			val = val*10 + int(ch-'0')
		}
	}
	return val
}

func compareParts(p1, p2 int) int {
	if p1 > p2 {
		return 1
	}
	if p1 < p2 {
		return -1
	}
	return 0
}

func compareVersionIndices(parts1, parts2 []string, idx int) int {
	p1, p2 := 0, 0
	if idx < len(parts1) {
		p1 = parseVersionPart(parts1[idx])
	}
	if idx < len(parts2) {
		p2 = parseVersionPart(parts2[idx])
	}
	return compareParts(p1, p2)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func compareVersionNumbers(v1, v2 string) int {
	p1 := strings.Split(v1, ".")
	p2 := strings.Split(v2, ".")
	maxLen := maxInt(len(p1), len(p2))
	for i := 0; i < maxLen; i++ {
		cmp := compareVersionIndices(p1, p2, i)
		if cmp != 0 {
			return cmp
		}
	}
	return 0
}

func isVersionAtLeast(actual, required string) bool {
	return compareVersionNumbers(actual, required) >= 0
}

func checkToolAvailable(tool string) *apperror.AppError {
	_, err := exec.LookPath(tool)
	if err == nil {
		return nil
	}
	return apperror.NewWithDetails(
		"checkToolAvailable",
		ErrPrerequisiteFailed,
		fmt.Sprintf("required tool not found in PATH: %s", tool),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"tool": tool},
	)
}

func checkRequiredTools(tools []string) *apperror.AppError {
	for _, t := range tools {
		if appErr := checkToolAvailable(t); appErr != nil {
			return appErr
		}
	}
	return nil
}

func extractGlibcVersion(output string) string {
	for _, word := range strings.Fields(output) {
		trimmed := strings.Trim(word, "(),;")
		parts := strings.Split(trimmed, ".")
		if len(parts) >= 2 && parseVersionPart(parts[0]) > 0 {
			return trimmed
		}
	}
	return ""
}

func readGlibcVersion() string {
	out, err := exec.Command("ldd", "--version").Output()
	if err != nil {
		return ""
	}
	firstLine := strings.Split(string(out), "\n")[0]
	return extractGlibcVersion(firstLine)
}

func verifyLinuxGlibc() *apperror.AppError {
	version := readGlibcVersion()
	if version == "" {
		return nil
	}
	if isVersionAtLeast(version, "2.28") {
		return nil
	}
	return apperror.NewWithDetails(
		"verifyLinuxGlibc",
		ErrPrerequisiteFailed,
		fmt.Sprintf("glibc >= 2.28 is required (found %s)", version),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"found_glibc": version, "min_glibc": "2.28"},
	)
}

func verifyLinuxPrereqs() *apperror.AppError {
	if appErr := checkRequiredTools([]string{"tar", "curl"}); appErr != nil {
		return appErr
	}
	return verifyLinuxGlibc()
}

func extractWindowsBuildNumber(verOutput string) int {
	low := strings.ToLower(verOutput)
	idx := strings.Index(low, "version ")
	if idx == -1 {
		return 0
	}
	sub := strings.TrimRight(verOutput[idx+8:], "]\r\n ")
	parts := strings.Split(sub, ".")
	if len(parts) < 3 {
		return 0
	}
	return parseVersionPart(parts[2])
}

func readWindowsBuildNumber() int {
	cmd := exec.Command("cmd", "/c", "ver")
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	return extractWindowsBuildNumber(string(out))
}

func verifyWindowsPrereqs() *apperror.AppError {
	build := readWindowsBuildNumber()
	if build == 0 {
		return nil
	}
	if build >= 17763 {
		return nil
	}
	return apperror.NewWithDetails(
		"verifyWindowsPrereqs",
		ErrPrerequisiteFailed,
		fmt.Sprintf("Windows 10 1809 (build 17763) or newer required (found build %d)", build),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"build": build, "min_build": 17763},
	)
}

func readDarwinProductVersion() string {
	cmd := exec.Command("sw_vers", "-productVersion")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func verifyDarwinVersion() *apperror.AppError {
	prod := readDarwinProductVersion()
	if prod == "" {
		return nil
	}
	if isVersionAtLeast(prod, "12.0") {
		return nil
	}
	return apperror.NewWithDetails(
		"verifyDarwinVersion",
		ErrPrerequisiteFailed,
		fmt.Sprintf("macOS 12.0 (Monterey) or later is required (found %s)", prod),
		"installer",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		map[string]any{"version": prod, "min_version": "12.0"},
	)
}

func verifyDarwinPrereqs() *apperror.AppError {
	if appErr := checkRequiredTools([]string{"curl", "hdiutil", "ditto", "xattr"}); appErr != nil {
		return appErr
	}
	return verifyDarwinVersion()
}

func verifyAntigravityPrerequisites(info AntigravityPlatformInfo) *apperror.AppError {
	if info.IsWindows() {
		return verifyWindowsPrereqs()
	}
	if info.IsDarwin() {
		return verifyDarwinPrereqs()
	}
	return verifyLinuxPrereqs()
}
