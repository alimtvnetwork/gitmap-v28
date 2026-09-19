package cmdautomation

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

var (
	reConstantsVersion = regexp.MustCompile(`var\s+Version\s*=\s*"([^"]+)"`)
	rePackageVersion   = regexp.MustCompile(`"version"\s*:\s*"([^"]+)"`)
	reCargoVersion     = regexp.MustCompile(`(?m)^version\s*=\s*"([^"]+)"`)
	reSpecVersion      = regexp.MustCompile(`(?m)^\*\*Version:\*\*\s*([^\s\r\n]+)`)
	reChangelogVersion = regexp.MustCompile(`(?m)^##\s*\[?v?([0-9]+\.[0-9]+\.[0-9]+)\]?`)
)

type versionTarget struct {
	name    string
	relPath string
	re      *regexp.Regexp
	isOpt   bool
}

var versionTargets = []versionTarget{
	{"Go Constants", "cli/constants/constants.go", reConstantsVersion, false},
	{"Package Manifest", "package.json", rePackageVersion, false},
	{"Changelog Header", "changelog.md", reChangelogVersion, false},
	{"Cargo Manifest", "Cargo.toml", reCargoVersion, true},
	{"Spec Overview", "02-spec/00-overview.md", reSpecVersion, true},
}

// RunVersionSync audits and optionally synchronizes repository version manifests.
func RunVersionSync(opts VersionSyncOptions) VersionSyncResultMonad {
	start := time.Now()
	root := resolveAuditDir(opts.Dir)
	canonical := resolveCanonicalVersion(root, opts.TargetVersion)
	isCanonicalMissing := canonical == ""
	if isCanonicalMissing {
		err := apperror.NewValidationError("no canonical version found in version.json or package.json")
		return result.Fail[VersionSyncResult](err)
	}
	items := auditVersionFiles(root, canonical)
	fixedCount := performFixesIfRequested(root, items, canonical, opts.IsFixMode)
	mismatchCount := countMismatches(items)
	res := buildVersionSyncResult(canonical, items, mismatchCount, fixedCount, time.Since(start))
	return result.Ok(res)
}

func performFixesIfRequested(root string, items []VersionCheckItem, canonical string, isFixMode bool) int {
	if isFixMode {
		return applyVersionFixes(root, items, canonical)
	}
	return 0
}

func countMismatches(items []VersionCheckItem) int {
	count := 0
	for _, item := range items {
		isMismatch := item.IsMatch == false
		if isMismatch {
			count++
		}
	}
	return count
}

func buildVersionSyncResult(canonical string, items []VersionCheckItem, mismatchCount, fixedCount int, dur time.Duration) VersionSyncResult {
	isClean := mismatchCount == 0 || (fixedCount > 0 && fixedCount == mismatchCount)
	return VersionSyncResult{
		CanonicalVersion: canonical, TotalChecked: len(items),
		MismatchCount: mismatchCount, FixedCount: fixedCount,
		Items: items, Duration: dur, IsClean: isClean,
	}
}

func resolveCanonicalVersion(root, target string) string {
	hasTarget := target != ""
	if hasTarget {
		return strings.TrimPrefix(target, "v")
	}
	return readCanonicalFromFiles(root)
}

func readCanonicalFromFiles(root string) string {
	vPath := filepath.Join(root, "version.json")
	ver := readJsonVersion(vPath)
	hasVer := ver != ""
	if hasVer {
		return ver
	}
	pPath := filepath.Join(root, "package.json")
	return extractVersionFromJson(pPath, "version")
}

func readJsonVersion(vPath string) string {
	ver := extractVersionFromJson(vPath, "Version")
	hasVer := ver != ""
	if hasVer {
		return ver
	}
	return extractVersionFromJson(vPath, "version")
}

func extractVersionFromJson(path, key string) string {
	data, err := os.ReadFile(path)
	hasReadErr := err != nil
	if hasReadErr {
		return ""
	}
	var obj map[string]any
	hasJsonErr := json.Unmarshal(data, &obj) != nil
	if hasJsonErr {
		return ""
	}
	raw, hasKey := obj[key].(string)
	if hasKey {
		return strings.TrimPrefix(raw, "v")
	}
	return ""
}

func auditVersionFiles(root, canonical string) []VersionCheckItem {
	var items []VersionCheckItem
	for _, target := range versionTargets {
		relPath := target.relPath
		isConstants := target.relPath == "cli/constants/constants.go"
		if isConstants {
			relPath = resolveConstantsPath(root)
		}
		item := auditTargetFile(target.name, relPath, filepath.Join(root, relPath), canonical, target.re, target.isOpt)
		hasTargetFile := item.TargetFile != ""
		if hasTargetFile {
			items = append(items, item)
		}
	}
	return items
}

func resolveConstantsPath(root string) string {
	p1 := filepath.Join(root, "cli/constants/constants_version.go")
	isP1Exist := fileExists(p1)
	if isP1Exist {
		return "cli/constants/constants_version.go"
	}
	return "cli/constants/constants.go"
}

func auditTargetFile(name, relPath, fullPath, canonical string, re *regexp.Regexp, isOpt bool) VersionCheckItem {
	isMissing := fileExists(fullPath) == false
	if isMissing && isOpt {
		return VersionCheckItem{}
	}
	if isMissing {
		return makeStatusItem(name, relPath, "", canonical, false, fmt.Sprintf("File not found: %s", relPath))
	}
	content := readFileContent(fullPath)
	match := re.FindStringSubmatch(content)
	isUnmatched := len(match) < 2
	if isUnmatched {
		return makeStatusItem(name, relPath, "", canonical, false, fmt.Sprintf("Version pattern not found in %s", relPath))
	}
	current := strings.TrimPrefix(match[1], "v")
	isMatch := current == canonical
	msg := formatMatchMessage(name, current, canonical, isMatch)
	return makeStatusItem(name, relPath, current, canonical, isMatch, msg)
}

func formatMatchMessage(name, current, canonical string, isMatch bool) string {
	if isMatch {
		return fmt.Sprintf("%s matches canonical v%s", name, canonical)
	}
	return fmt.Sprintf("%s version 'v%s' != canonical 'v%s'", name, current, canonical)
}

func makeStatusItem(name, relPath, cur, exp string, isMatch bool, msg string) VersionCheckItem {
	return VersionCheckItem{
		Name: name, TargetFile: relPath, CurrentVersion: cur,
		ExpectedVersion: exp, IsMatch: isMatch, Message: msg,
	}
}

func applyVersionFixes(root string, items []VersionCheckItem, canonical string) int {
	fixed := 0
	for _, item := range items {
		isMatch := item.IsMatch
		if isMatch {
			continue
		}
		isFixed := fixVersionFile(filepath.Join(root, item.TargetFile), canonical)
		if isFixed {
			fixed++
		}
	}
	return fixed
}

func fixVersionFile(fullPath, canonical string) bool {
	data, err := os.ReadFile(fullPath)
	hasErr := err != nil
	if hasErr {
		return false
	}
	updated := replaceVersionInContent(fullPath, string(data), canonical)
	isUnchanged := updated == string(data)
	if isUnchanged {
		return false
	}
	writeErr := os.WriteFile(fullPath, []byte(updated), 0644)
	return writeErr == nil
}

func replaceVersionInContent(path, content, canonical string) string {
	base := filepath.Base(path)
	switch base {
	case "package.json":
		return rePackageVersion.ReplaceAllString(content, fmt.Sprintf(`"version": "%s"`, canonical))
	case "Cargo.toml":
		return reCargoVersion.ReplaceAllString(content, fmt.Sprintf(`version = "%s"`, canonical))
	case "00-overview.md":
		return reSpecVersion.ReplaceAllString(content, fmt.Sprintf(`**Version:** %s`, canonical))
	case "version.json":
		res := regexp.MustCompile(`"Version"\s*:\s*"[^"]+"`).ReplaceAllString(content, fmt.Sprintf(`"Version": "%s"`, canonical))
		return rePackageVersion.ReplaceAllString(res, fmt.Sprintf(`"version": "%s"`, canonical))
	default:
		return reConstantsVersion.ReplaceAllString(content, fmt.Sprintf(`var Version = "%s"`, canonical))
	}
}

func readFileContent(fullPath string) string {
	data, err := os.ReadFile(fullPath)
	isFail := err != nil
	if isFail {
		return ""
	}
	return string(data)
}

func fileExists(fullPath string) bool {
	info, err := os.Stat(fullPath)
	isMissing := err != nil || info.IsDir()
	if isMissing {
		return false
	}
	return true
}
