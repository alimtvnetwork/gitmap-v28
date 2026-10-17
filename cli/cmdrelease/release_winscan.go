// Package cmdrelease — release_winscan.go: `gitmap release winscan <tag>`.
//
// Pre-publish scan of the Windows release asset (zip) against the 7 Windows
// Defender heuristic trigger patterns from RCA 273 (commit 385aba38e).
// Downloads the asset to a temp dir via `gh release download`, extracts the
// exe with archive/zip, and byte-scans it. Read-only.
//
// Wiring (lead): same pattern as RunReleaseVerify — add `winscan` to
// RunRelease's sub routing with an explicit nil check.
package cmdrelease

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
)

// winscanPattern is one Defender heuristic trigger pattern (RCA 273 §3).
// A pattern matches when ALL needles hit the binary bytes; the reported
// offset is the first needle's first hit.
type winscanPattern struct {
	id      string
	name    string
	needles []*regexp.Regexp
}

var (
	winscanWordReg    = regexp.MustCompile(`(?i)(^|[^A-Za-z])reg(\.exe)?([^A-Za-z]|$)`)
	winscanWordAdd    = regexp.MustCompile(`(?i)(^|[^A-Za-z])add([^A-Za-z]|$)`)
	winscanWordDelete = regexp.MustCompile(`(?i)(^|[^A-Za-z])delete([^A-Za-z]|$)`)
	winscanWordCmd    = regexp.MustCompile(`(?i)(^|[^A-Za-z])cmd(\.exe)?([^A-Za-z]|$)`)
	winscanRunKeyPath = regexp.MustCompile(`Software\\Microsoft\\Windows\\CurrentVersion\\Run`)
)

// winscanPatterns — the 7 trigger patterns confirmed in the 385aba38e diff
// (.ai-memory/plans/273-windows-virus-rca.md §3).
var winscanPatterns = []winscanPattern{
	{id: "P1", name: "reg.exe Run-key persistence write", needles: []*regexp.Regexp{
		winscanWordReg, winscanWordAdd, winscanRunKeyPath}},
	{id: "P2", name: "vssadmin shadow-copy create/delete", needles: []*regexp.Regexp{
		regexp.MustCompile(`(?i)vssadmin`), regexp.MustCompile(`(?i)shadow`)}},
	{id: "P3", name: "Winlogon auto-login credential write", needles: []*regexp.Regexp{
		regexp.MustCompile(`Winlogon`), regexp.MustCompile(`DefaultPassword`)}},
	{id: "P4", name: "hidden PowerShell WMI process enumeration", needles: []*regexp.Regexp{
		regexp.MustCompile(`(?i)powershell`), regexp.MustCompile(`Win32_Process`)}},
	{id: "P5", name: "spawned taskkill process kill", needles: []*regexp.Regexp{
		regexp.MustCompile(`(?i)taskkill`)}},
	{id: "P6", name: "delayed self-deletion via cmd.exe", needles: []*regexp.Regexp{
		winscanWordCmd, regexp.MustCompile(`(?i)del /F`)}},
	{id: "P7", name: "reg.exe Run-key deletion loop", needles: []*regexp.Regexp{
		winscanWordReg, winscanWordDelete, regexp.MustCompile(`\\Run`)}},
}

// winscanMatch is one pattern hit with its byte offset.
type winscanMatch struct {
	pattern winscanPattern
	offset  int
}

// RunReleaseWinScan executes `gitmap release winscan <tag>` — read-only.
// args is the tail after the `winscan` token: <tag>.
func RunReleaseWinScan(args []string) *appfault.AppError {
	if isReleaseNewCmdHelpRequested(args) {
		fmt.Print(releaseWinScanHelpText)
		cliexit.Exit(0)
	}

	if len(args) == 0 || strings.TrimSpace(args[0]) == "" {
		return appfault.NewValidationError("missing required argument: <tag>\n" + releaseWinScanUsageLine +
			"example: gitmap release winscan v6.558.0")
	}
	tag := strings.TrimSpace(args[0])

	if appErr := requireGhCLI(); appErr != nil {
		return appErr
	}

	tmpDir, err := os.MkdirTemp("", "gitmap-winscan-")
	if err != nil {
		return appfault.NewValidationError("could not create temp dir: " + err.Error())
	}
	defer os.RemoveAll(tmpDir)

	zipPath, appErr := downloadWindowsAsset(tag, tmpDir)
	if appErr != nil {
		return appErr
	}
	exePath, appErr := extractLargestExe(zipPath, tmpDir)
	if appErr != nil {
		return appErr
	}

	data, err := os.ReadFile(exePath)
	if err != nil {
		return appfault.NewValidationError("could not read extracted exe: " + err.Error())
	}

	return reportWinScan(filepath.Base(zipPath), filepath.Base(exePath), data)
}

// downloadWindowsAsset downloads the windows amd64 zip for the tag into dir.
// The repo resolves from the current directory's git remote via gh.
func downloadWindowsAsset(tag, dir string) (string, *appfault.AppError) {
	if _, appErr := runGhCLI("release", "download", tag,
		"--pattern", "*windows*amd64*.zip", "--dir", dir); appErr != nil {
		return "", appfault.NewValidationError(
			fmt.Sprintf("no windows amd64 zip asset found for tag %s: %s", tag, appErr.Error()))
	}

	zips, err := filepath.Glob(filepath.Join(dir, "*.zip"))
	if err != nil || len(zips) == 0 {
		return "", appfault.NewValidationError(
			fmt.Sprintf("no windows amd64 zip asset downloaded for tag %s", tag))
	}
	sort.Strings(zips)

	return zips[0], nil
}

// extractLargestExe unzips the asset and returns the largest .exe inside.
func extractLargestExe(zipPath, dir string) (string, *appfault.AppError) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", appfault.NewValidationError("could not open zip: " + err.Error())
	}
	defer reader.Close()

	best := ""
	var bestSize int64 = -1
	for _, file := range reader.File {
		if file.FileInfo().IsDir() || !strings.HasSuffix(strings.ToLower(file.Name), ".exe") {
			continue
		}
		if file.UncompressedSize64 <= uint64(bestSize) {
			continue
		}
		outPath := filepath.Join(dir, "exe", filepath.Base(file.Name))
		if appErr := extractZipFile(file, outPath); appErr != nil {
			return "", appErr
		}
		best = outPath
		bestSize = int64(file.UncompressedSize64)
	}

	if best == "" {
		return "", appfault.NewValidationError("no .exe found inside " + filepath.Base(zipPath))
	}

	return best, nil
}

// extractZipFile writes one zip entry to outPath.
func extractZipFile(file *zip.File, outPath string) *appfault.AppError {
	src, err := file.Open()
	if err != nil {
		return appfault.NewValidationError("could not read zip entry: " + err.Error())
	}
	defer src.Close()

	if err := os.MkdirAll(filepath.Dir(outPath), 0o755); err != nil {
		return appfault.NewValidationError("could not create temp dir: " + err.Error())
	}
	dst, err := os.Create(outPath)
	if err != nil {
		return appfault.NewValidationError("could not write temp file: " + err.Error())
	}
	defer dst.Close()

	buf := make([]byte, 32*1024)
	for {
		n, readErr := src.Read(buf)
		if n > 0 {
			if _, writeErr := dst.Write(buf[:n]); writeErr != nil {
				return appfault.NewValidationError("could not write temp file: " + writeErr.Error())
			}
		}
		if readErr != nil {
			break
		}
	}

	return nil
}

// scanBinary runs the 7 trigger patterns over the exe bytes.
func scanBinary(data []byte) []winscanMatch {
	matches := []winscanMatch{}
	for _, pattern := range winscanPatterns {
		offset, isMatched := matchWinscanPattern(data, pattern)
		if isMatched {
			matches = append(matches, winscanMatch{pattern: pattern, offset: offset})
		}
	}

	return matches
}

// matchWinscanPattern reports whether every needle hits; offset is the
// earliest needle hit.
func matchWinscanPattern(data []byte, pattern winscanPattern) (int, bool) {
	earliest := -1
	for _, needle := range pattern.needles {
		loc := needle.FindIndex(data)
		if loc == nil {
			return 0, false
		}
		if earliest < 0 || loc[0] < earliest {
			earliest = loc[0]
		}
	}

	return earliest, true
}

// reportWinScan prints the scan verdict; any pattern hit fails the command.
func reportWinScan(zipName, exeName string, data []byte) *appfault.AppError {
	matches := scanBinary(data)
	fmt.Printf("%s (%s, %.1f MB):\n", zipName, exeName, float64(len(data))/(1024*1024))

	if len(matches) == 0 {
		fmt.Printf("  0/%d trigger patterns ✓\n", len(winscanPatterns))

		return nil
	}

	for _, match := range matches {
		fmt.Printf("  ✗ %s (%s) matched at offset 0x%x\n", match.pattern.id, match.pattern.name, match.offset)
	}

	return appfault.NewValidationError(
		fmt.Sprintf("winscan: %d/%d Defender trigger patterns matched — review before publishing",
			len(matches), len(winscanPatterns)))
}
