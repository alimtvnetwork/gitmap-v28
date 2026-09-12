package cmdinstall

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

const agManagerGitURL = "https://github.com/lbjlaq/Antigravity-Manager.git"

func fetchAgManagerFromGitTags() (string, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "git", "ls-remote", "--tags", "--refs", agManagerGitURL)
	out, err := cmd.Output()
	if err != nil {
		return "", "", fmt.Errorf("git ls-remote failed: %w", err)
	}

	return resolveHighestTagFromOutput(string(out))
}

func resolveHighestTagFromOutput(raw string) (string, string, error) {
	ver, isFound := findHighestSemverInOutput(raw)
	if !isFound {
		return "", "", fmt.Errorf("no valid semver release tags found in git repository")
	}

	tag := "v" + ver.CoreString()
	url := constructAgManagerAssetURL(tag, ver.CoreString())
	if url == "" {
		return "", "", fmt.Errorf("unsupported platform or architecture for %s", tag)
	}

	return url, ver.CoreString(), nil
}

func findHighestSemverInOutput(raw string) (release.Version, bool) {
	var highest release.Version
	isFound := false
	for _, line := range strings.Split(raw, "\n") {
		v, ok := parseTagFromGitLine(line)
		if ok && (!isFound || v.GreaterThan(highest)) {
			highest = v
			isFound = true
		}
	}

	return highest, isFound
}

func parseTagFromGitLine(line string) (release.Version, bool) {
	trimmed := strings.TrimSpace(line)
	idx := strings.Index(trimmed, "refs/tags/")
	if idx < 0 {
		return release.Version{}, false
	}

	tagStr := trimmed[idx+len("refs/tags/"):]
	v, err := release.Parse(tagStr)
	if err != nil {
		return release.Version{}, false
	}

	return v, true
}

func resolveAgManagerAssetURL(reqVer string) (string, string, error) {
	clean := strings.TrimPrefix(strings.TrimSpace(reqVer), "v")
	if clean == "" {
		return getAgManagerAssetURL()
	}

	tag := "v" + clean
	url := constructAgManagerAssetURL(tag, clean)
	if url == "" {
		return "", "", fmt.Errorf("unsupported platform for version %s", clean)
	}

	return url, clean, nil
}
