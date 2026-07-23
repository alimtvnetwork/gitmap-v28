// Package release handles version parsing, release workflows,
// GitHub integration, and release metadata management.
package release

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// Version represents a parsed semantic version.
type Version struct {
	Major      int
	Minor      int
	Patch      int
	PreRelease string
	Raw        string
}

// Parse parses a version string into a Version struct.
// Accepts formats: v1, v1.2, v1.2.3, v1.2.3-rc.1
func Parse(input string) (Version, error) {
	raw := input
	s := strings.TrimPrefix(input, "v")
	if len(s) == 0 {
		return Version{}, fmt.Errorf("empty version string")
	}

	preRelease := extractPreRelease(s)
	core := s
	if len(preRelease) > 0 {
		core = s[:len(s)-len(preRelease)-1]
	}

	major, minor, patch, err := parseCoreSegments(core)
	if err != nil {
		return Version{}, err
	}

	return Version{
		Major: major, Minor: minor, Patch: patch,
		PreRelease: preRelease, Raw: raw,
	}, nil
}

// extractPreRelease splits off the pre-release suffix after a hyphen.
func extractPreRelease(s string) string {
	idx := strings.Index(s, "-")
	if idx < 0 {
		return ""
	}

	return s[idx+1:]
}

// parseCoreSegments parses the X.Y.Z portion, padding missing segments.
func parseCoreSegments(core string) (major, minor, patch int, err error) {
	parts := strings.Split(core, ".")
	if len(parts) > 3 {
		return 0, 0, 0, fmt.Errorf("too many version segments")
	}

	major, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, 0, fmt.Errorf("invalid major version: %s", parts[0])
	}

	minor = 0
	if len(parts) >= 2 {
		minor, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid minor version: %s", parts[1])
		}
	}

	patch = 0
	if len(parts) >= 3 {
		patch, err = strconv.Atoi(parts[2])
		if err != nil {
			return 0, 0, 0, fmt.Errorf("invalid patch version: %s", parts[2])
		}
	}

	return major, minor, patch, nil
}

// String returns the padded version with v prefix (e.g. v1.2.3).
func (v Version) String() string {
	base := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.PreRelease) > 0 {
		return base + "-" + v.PreRelease
	}

	return base
}

// CoreString returns the version without v prefix (e.g. 1.2.3).
func (v Version) CoreString() string {
	base := fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
	if len(v.PreRelease) > 0 {
		return base + "-" + v.PreRelease
	}

	return base
}

// IsPreRelease returns true if the version has a pre-release suffix.
func (v Version) IsPreRelease() bool {
	return len(v.PreRelease) > 0
}

// GreaterThan returns true if v is strictly greater than other.
func (v Version) GreaterThan(other Version) bool {
	if v.Major > other.Major {
		return true
	}
	if v.Major < other.Major {
		return false
	}
	if v.Minor > other.Minor {
		return true
	}
	if v.Minor < other.Minor {
		return false
	}
	if v.Patch > other.Patch {
		return true
	}
	if v.Patch < other.Patch {
		return false
	}

	return preReleaseGreater(v.PreRelease, other.PreRelease)
}

// preReleaseGreater compares pre-release precedence.
// Stable (empty) is greater than any pre-release.
func preReleaseGreater(a, b string) bool {
	if len(a) == 0 && len(b) > 0 {
		return true
	}

	return false
}

// Bump increments the version by the given level (major, minor, patch).
func Bump(v Version, level string) (Version, error) {
	if level == constants.BumpMajor {
		return Version{Major: v.Major + 1, Minor: 0, Patch: 0}, nil
	}
	if level == constants.BumpMinor {
		return Version{Major: v.Major, Minor: v.Minor + 1, Patch: 0}, nil
	}
	if level == constants.BumpPatch {
		return Version{Major: v.Major, Minor: v.Minor, Patch: v.Patch + 1}, nil
	}

	return Version{}, fmt.Errorf("invalid bump level: %s (use major, minor, or patch)", level)
}
