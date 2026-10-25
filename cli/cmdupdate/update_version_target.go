package cmdupdate

import (
	"slices"
	"strconv"
	"strings"
)

// NormalizeVersion strips any leading 'v' or 'V' and trims whitespace.
func NormalizeVersion(v string) string {
	cleaned := strings.TrimSpace(v)
	cleaned = strings.TrimPrefix(cleaned, "v")
	cleaned = strings.TrimPrefix(cleaned, "V")

	return cleaned
}

// FormatVersionTag formats a version string with a leading 'v'.
func FormatVersionTag(v string) string {
	norm := NormalizeVersion(v)
	if norm == "" {
		return ""
	}

	return "v" + norm
}

func cleanSemVerPart(part string) int {
	clean := strings.Split(strings.Split(part, "-")[0], "+")[0]
	n, _ := strconv.Atoi(clean)

	return n
}

func parseSemVerTriple(v string) [3]int {
	var out [3]int
	parts := strings.Split(NormalizeVersion(v), ".")

	for i := 0; i < 3 && i < len(parts); i++ {
		out[i] = cleanSemVerPart(parts[i])
	}

	return out
}

func compareTriples(t1, t2 [3]int) int {
	for i := 0; i < 3; i++ {
		if t1[i] < t2[i] {
			return -1
		}

		if t1[i] > t2[i] {
			return 1
		}
	}

	return 0
}

// CompareSemVer compares two semantic version strings.
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2.
func CompareSemVer(v1, v2 string) int {
	norm1, norm2 := NormalizeVersion(v1), NormalizeVersion(v2)
	if norm1 == norm2 {
		return 0
	}

	if norm1 == "" {
		return -1
	}

	if norm2 == "" {
		return 1
	}

	return compareTriples(parseSemVerTriple(norm1), parseSemVerTriple(norm2))
}

// IsNewerVersion returns true if target version is strictly newer than current version.
func IsNewerVersion(target, current string) bool {
	return CompareSemVer(target, current) > 0
}

// IsAlreadyUpdated checks if current version is at or ahead of target version and force is not set.
func IsAlreadyUpdated(currentVersion, targetVersion string, isForce bool) bool {
	if isForce {
		return false
	}

	normCurrent := NormalizeVersion(currentVersion)
	normTarget := NormalizeVersion(targetVersion)
	if normCurrent == "" || normTarget == "" {
		return false
	}

	return CompareSemVer(normTarget, normCurrent) <= 0
}

// IsUpdateAllowed checks whether an update to targetVersion is permitted given currentVersion and options.
func IsUpdateAllowed(currentVersion, targetVersion string, isForce, isAllowDowngrade bool) bool {
	if isForce || isAllowDowngrade {
		return true
	}

	return IsNewerVersion(targetVersion, currentVersion)
}

// ParseUpdateOptions extracts UpdateOptions from CLI arguments.
func ParseUpdateOptions(args []string) UpdateOptions {
	opts := UpdateOptions{
		MaxFallbackTags: 5,
	}

	for i := 0; i < len(args); i++ {
		i = parseUpdateArg(args, i, &opts)
	}

	applyTargetPinnedVersion(&opts)

	return opts
}

func parseUpdateArg(args []string, i int, opts *UpdateOptions) int {
	arg := args[i]
	if parseUpdateFlag(arg, opts) {
		return i
	}

	if nextIdx, hasVersion := parseVersionFlag(args, i, opts); hasVersion {
		return nextIdx
	}

	parsePositionalVersion(arg, opts)

	return i
}

func parseUpdateFlag(arg string, opts *UpdateOptions) bool {
	if parseUpdateModeFlag(arg, opts) {
		return true
	}

	return parseOutputFlag(arg, opts)
}

func parseUpdateModeFlag(arg string, opts *UpdateOptions) bool {
	switch arg {
	case "-f", "--force":
		opts.IsForce = true
	case "-d", "--allow-downgrade":
		opts.IsAllowDowngrade = true
	default:
		return false
	}

	return true
}

func parseOutputFlag(arg string, opts *UpdateOptions) bool {
	switch arg {
	case "-n", "--dry-run":
		opts.IsDryRun = true
	case "-j", "--json":
		opts.IsJSON = true
	case "-q", "--quiet":
		opts.IsQuiet = true
	default:
		return false
	}

	return true
}

func parseVersionFlag(args []string, i int, opts *UpdateOptions) (int, bool) {
	arg := args[i]
	if slices.Contains([]string{"-v", "--version"}, arg) && i+1 < len(args) {
		opts.TargetVersion = args[i+1]

		return i + 1, true
	}

	if val, hasPrefix := extractVersionPrefix(arg); hasPrefix {
		opts.TargetVersion = val

		return i, true
	}

	return i, false
}

func extractVersionPrefix(arg string) (string, bool) {
	if strings.HasPrefix(arg, "--version=") {
		return strings.TrimPrefix(arg, "--version="), true
	}

	if strings.HasPrefix(arg, "-v=") {
		return strings.TrimPrefix(arg, "-v="), true
	}

	return "", false
}

func parsePositionalVersion(arg string, opts *UpdateOptions) {
	if strings.HasPrefix(arg, "-") || opts.TargetVersion != "" {
		return
	}

	if isSemverLikeVersion(arg) {
		opts.TargetVersion = arg
	}
}

func applyTargetPinnedVersion(opts *UpdateOptions) {
	if targetPinnedVersion != "" && opts.TargetVersion == "" {
		opts.TargetVersion = targetPinnedVersion
	}
}

func isDigitsOnly(s string) bool {
	if len(s) == 0 {
		return false
	}

	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}

	return true
}

func isSemverLikeVersion(s string) bool {
	parts := strings.Split(NormalizeVersion(s), ".")
	if len(parts) < 2 {
		return false
	}

	for _, p := range parts {
		if !isDigitsOnly(p) {
			return false
		}
	}

	return true
}
