package cmdupdate

import (
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

// IsAlreadyUpdated checks if current version matches target version and force is not set.
func IsAlreadyUpdated(currentVersion, targetVersion string, isForce bool) bool {
	if isForce {
		return false
	}

	normCurrent := NormalizeVersion(currentVersion)
	normTarget := NormalizeVersion(targetVersion)

	if normCurrent == "" || normTarget == "" {
		return false
	}

	return normCurrent == normTarget
}

// ParseUpdateOptions extracts UpdateOptions from CLI arguments.
func ParseUpdateOptions(args []string) UpdateOptions {
	opts := UpdateOptions{
		MaxFallbackTags: 5,
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]

		switch {
		case arg == "-f" || arg == "--force":
			opts.IsForce = true

		case arg == "-n" || arg == "--dry-run":
			opts.IsDryRun = true

		case arg == "-j" || arg == "--json":
			opts.IsJSON = true

		case arg == "-q" || arg == "--quiet":
			opts.IsQuiet = true

		case (arg == "-v" || arg == "--version") && i+1 < len(args):
			opts.TargetVersion = args[i+1]
			i++

		case strings.HasPrefix(arg, "--version="):
			opts.TargetVersion = strings.TrimPrefix(arg, "--version=")

		case strings.HasPrefix(arg, "-v="):
			opts.TargetVersion = strings.TrimPrefix(arg, "-v=")

		case !strings.HasPrefix(arg, "-") && opts.TargetVersion == "":
			if isSemverLikeVersion(arg) {
				opts.TargetVersion = arg
			}
		}
	}

	if targetPinnedVersion != "" && opts.TargetVersion == "" {
		opts.TargetVersion = targetPinnedVersion
	}

	return opts
}

func isSemverLikeVersion(s string) bool {
	norm := NormalizeVersion(s)
	parts := strings.Split(norm, ".")
	if len(parts) < 2 {
		return false
	}

	for _, p := range parts {
		if len(p) == 0 {
			return false
		}

		for _, r := range p {
			if r < '0' || r > '9' {
				return false
			}
		}
	}

	return true
}
