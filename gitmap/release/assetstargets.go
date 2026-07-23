// Package release — assetstargets.go defines the default cross-compile target matrix.
package release

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// DefaultTargets returns the standard 6-target cross-compilation matrix.
func DefaultTargets() []BuildTarget {
	return []BuildTarget{
		{GOOS: "windows", GOARCH: "amd64"},
		{GOOS: "windows", GOARCH: "arm64"},
		{GOOS: "linux", GOARCH: "amd64"},
		{GOOS: "linux", GOARCH: "arm64"},
		{GOOS: "darwin", GOARCH: "amd64"},
		{GOOS: "darwin", GOARCH: "arm64"},
	}
}

// ResolveTargets determines the final target list using three-layer config:
// CLI flag (highest) → config.json release.targets (middle) → defaults (lowest).
func ResolveTargets(flagTargets string, configTargets []model.ReleaseTarget) ([]BuildTarget, error) {
	if len(flagTargets) > 0 {
		return ParseTargets(flagTargets)
	}

	if len(configTargets) > 0 {
		return convertConfigTargets(configTargets), nil
	}

	return DefaultTargets(), nil
}

// convertConfigTargets converts model.ReleaseTarget slice to BuildTarget slice.
func convertConfigTargets(targets []model.ReleaseTarget) []BuildTarget {
	result := make([]BuildTarget, len(targets))

	for i, t := range targets {
		result[i] = BuildTarget{GOOS: t.GOOS, GOARCH: t.GOARCH}
	}

	return result
}

// ParseTargets parses a comma-separated "os/arch" string into BuildTargets.
// Example: "windows/amd64,linux/arm64"
func ParseTargets(input string) ([]BuildTarget, error) {
	if len(input) == 0 {
		return DefaultTargets(), nil
	}

	parts := strings.Split(input, ",")
	targets := make([]BuildTarget, 0, len(parts))

	for _, p := range parts {
		t, err := parseOneTarget(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}

		targets = append(targets, t)
	}

	return targets, nil
}

// parseOneTarget parses "os/arch" into a BuildTarget.
func parseOneTarget(s string) (BuildTarget, error) {
	parts := strings.SplitN(s, "/", 2)
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return BuildTarget{}, fmt.Errorf("invalid target format %q — expected os/arch", s)
	}

	return BuildTarget{GOOS: parts[0], GOARCH: parts[1]}, nil
}

// DescribeTargets returns human-readable names for dry-run output.
func DescribeTargets(binName, version string, targets []BuildTarget) []string {
	names := make([]string, 0, len(targets))

	for _, t := range targets {
		names = append(names, formatOutputName(binName, version, t))
	}

	return names
}
