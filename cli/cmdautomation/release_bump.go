package cmdautomation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunReleaseBump computes the next SemVer version and updates repository manifests.
func RunReleaseBump(opts ReleaseBumpOptions) ReleaseBumpResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	curVer := readCanonicalFromFiles(root)
	isCurVerEmpty := curVer == ""
	if isCurVerEmpty {
		curVer = "1.0.0"
	}
	tier := resolveBumpTier(opts.BumpType, root)
	newVer := resolveTargetVersion(opts.TargetVersion, curVer, tier)
	updatedFiles := planOrApplyRelease(root, curVer, newVer, opts)
	tag := fmt.Sprintf("v%s", newVer)
	isTagged := executeTagIfRequested(root, tag, opts.IsTag, opts.IsDryRun)
	isPushed := executePushIfRequested(root, tag, opts.IsPush, opts.IsDryRun)
	res := buildReleaseBumpResult(curVer, newVer, tier, tag, updatedFiles, isTagged, isPushed, time.Since(start))
	return result.Ok(res)
}

func resolveBumpTier(tier, root string) string {
	hasTier := tier != "" && tier != "auto"
	if hasTier {
		return strings.ToLower(tier)
	}
	return detectBumpTierFromCommits(root)
}

func detectBumpTierFromCommits(root string) string {
	cmd := exec.Command("git", "-C", root, "log", "-n", "25", "--oneline")
	out, err := cmd.Output()
	hasErr := err != nil
	if hasErr {
		return "patch"
	}
	text := strings.ToLower(string(out))
	isMajor := strings.Contains(text, "breaking") || strings.Contains(text, "!:")
	if isMajor {
		return "major"
	}
	isMinor := strings.Contains(text, "feat:") || strings.Contains(text, "feat(")
	if isMinor {
		return "minor"
	}
	return "patch"
}

func resolveTargetVersion(target, current, tier string) string {
	hasTarget := target != ""
	if hasTarget {
		return strings.TrimPrefix(target, "v")
	}
	return computeNextSemVer(current, tier)
}

func computeNextSemVer(current, tier string) string {
	cleaned := strings.TrimPrefix(current, "v")
	parts := strings.Split(cleaned, ".")
	isInvalid := len(parts) < 3
	if isInvalid {
		return "1.0.1"
	}
	major, _ := strconv.Atoi(parts[0])
	minor, _ := strconv.Atoi(parts[1])
	patch, _ := strconv.Atoi(parts[2])
	switch tier {
	case "major":
		return fmt.Sprintf("%d.0.0", major+1)
	case "minor":
		return fmt.Sprintf("%d.%d.0", major, minor+1)
	default:
		return fmt.Sprintf("%d.%d.%d", major, minor, patch+1)
	}
}

func planOrApplyRelease(root, curVer, newVer string, opts ReleaseBumpOptions) []string {
	manifests := []string{"version.json", "package.json", resolveConstantsPath(root), "changelog.md"}
	if opts.IsDryRun {
		return manifests
	}
	return updateReleaseManifests(root, newVer, opts.Bullets)
}

func updateReleaseManifests(root, newVer string, bullets []string) []string {
	var updated []string
	updateVersionJson(filepath.Join(root, "version.json"), newVer)
	updated = append(updated, "version.json")
	updatePackageJson(filepath.Join(root, "package.json"), newVer)
	updated = append(updated, "package.json")
	cPath := resolveConstantsPath(root)
	updateConstantsVersion(filepath.Join(root, cPath), newVer)
	updated = append(updated, cPath)
	updateChangelogFile(filepath.Join(root, "changelog.md"), newVer, bullets)
	updated = append(updated, "changelog.md")
	writeReleaseNotesFile(root, newVer, bullets)
	rnRel := fmt.Sprintf(".ai-memory/release/release-notes-v%s.md", newVer)
	updated = append(updated, rnRel)
	updateUserPreferencesFile(filepath.Join(root, ".ai-memory/user-preferences"), newVer)
	return updated
}

func updateVersionJson(fullPath, newVer string) bool {
	content := readFileContent(fullPath)
	isMissing := content == ""
	if isMissing {
		return false
	}
	reV := regexp.MustCompile(`("Version"\s*:\s*)"[^"]+"`)
	reLower := regexp.MustCompile(`("version"\s*:\s*)"[^"]+"`)
	updated := reV.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newVer))
	updated = reLower.ReplaceAllString(updated, fmt.Sprintf(`${1}"%s"`, newVer))
	writeErr := os.WriteFile(fullPath, []byte(updated), 0644)
	return writeErr == nil
}

func updatePackageJson(fullPath, newVer string) bool {
	content := readFileContent(fullPath)
	isMissing := content == ""
	if isMissing {
		return false
	}
	rePkg := regexp.MustCompile(`("version"\s*:\s*)"[^"]+"`)
	updated := rePkg.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newVer))
	writeErr := os.WriteFile(fullPath, []byte(updated), 0644)
	return writeErr == nil
}

func updateConstantsVersion(fullPath, newVer string) bool {
	content := readFileContent(fullPath)
	isMissing := content == ""
	if isMissing {
		return false
	}
	reConst := regexp.MustCompile(`(var\s+Version\s*=\s*)"[^"]+"`)
	updated := reConst.ReplaceAllString(content, fmt.Sprintf(`${1}"%s"`, newVer))
	writeErr := os.WriteFile(fullPath, []byte(updated), 0644)
	return writeErr == nil
}

func updateChangelogFile(fullPath, newVer string, bullets []string) bool {
	content := readFileContent(fullPath)
	today := time.Now().UTC().Format("2006-01-02")
	bulletLines := formatBulletPoints(bullets, newVer)
	entry := fmt.Sprintf("## [v%s] %s Release v%s\n\n### Added / Changed / Fixed / Removed\n\n%s\n\n",
		newVer, today, newVer, bulletLines)
	writeErr := os.WriteFile(fullPath, []byte(entry+content), 0644)
	return writeErr == nil
}

func writeReleaseNotesFile(root, newVer string, bullets []string) bool {
	relDir := filepath.Join(root, ".ai-memory", "release")
	_ = os.MkdirAll(relDir, 0755)
	rnPath := filepath.Join(relDir, fmt.Sprintf("release-notes-v%s.md", newVer))
	bulletLines := formatBulletPoints(bullets, newVer)
	body := fmt.Sprintf("## Release Notes v%s\n\n%s\n", newVer, bulletLines)
	writeErr := os.WriteFile(rnPath, []byte(body), 0644)
	return writeErr == nil
}

func formatBulletPoints(bullets []string, newVer string) string {
	hasBullets := len(bullets) > 0
	if hasBullets {
		var lines []string
		for _, b := range bullets {
			lines = append(lines, fmt.Sprintf("- %s", b))
		}
		return strings.Join(lines, "\n")
	}
	return fmt.Sprintf("- Bumped version to v%s across manifests\n- Automated release synchronization", newVer)
}

func updateUserPreferencesFile(fullPath, newVer string) bool {
	content := readFileContent(fullPath)
	isMissing := content == ""
	if isMissing {
		return false
	}
	rePref := regexp.MustCompile(`Release Mode Active \(v[0-9.]+\)`)
	updated := rePref.ReplaceAllString(content, fmt.Sprintf("Release Mode Active (v%s)", newVer))
	writeErr := os.WriteFile(fullPath, []byte(updated), 0644)
	return writeErr == nil
}

func executeTagIfRequested(root, tag string, isTag, isDryRun bool) bool {
	isAllowed := isTag && isDryRun == false
	if isAllowed {
		cmd := exec.Command("git", "-C", root, "tag", "-a", tag, "-m", fmt.Sprintf("Release %s", tag))
		return cmd.Run() == nil
	}
	return isTag
}

func executePushIfRequested(root, tag string, isPush, isDryRun bool) bool {
	isAllowed := isPush && isDryRun == false
	if isAllowed {
		cmd := exec.Command("git", "-C", root, "push", "origin", tag)
		return cmd.Run() == nil
	}
	return isPush
}

func buildReleaseBumpResult(cur, newVer, tier, tag string, files []string, isTagged, isPushed bool, dur time.Duration) ReleaseBumpResult {
	branch := fmt.Sprintf("release/v%s", newVer)
	return ReleaseBumpResult{
		PreviousVersion: cur, NewVersion: newVer, BumpType: tier,
		UpdatedFiles: files, ReleaseBranch: branch, TagName: tag,
		IsTagged: isTagged, IsPushed: isPushed, Duration: dur, IsSuccess: true,
	}
}
