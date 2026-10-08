package cmdchromeprofile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runSmartChromeImport(opts chromeTransferOptions) error {
	candidates, explicitTarget, err := collectImportCandidates(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chrome-profile-import: ERROR %v\n", err)
		printChromeProfileImportUsage()

		return err
	}

	if len(candidates) == 0 {
		fmt.Fprintf(os.Stderr, "chrome-profile-import: ERROR no matching profile snapshot files found\n")
		printChromeProfileImportUsage()

		return fmt.Errorf("no matching snapshot files found to import")
	}

	filtered := filterImportCandidates(candidates, opts)
	if len(filtered) == 0 {
		fmt.Println("No snapshot files remained after applying filters (--except / --email).")

		return nil
	}

	if opts.Limit > 0 && len(filtered) > opts.Limit {
		filtered = filtered[:opts.Limit]
	}

	fmt.Printf("\n\033[1;96mStarting Chrome Profile Import (%d snapshot(s))\033[0m\n\n", len(filtered))
	successCount := 0
	for _, f := range filtered {
		if err := importSingleSnapshotWithStepLogging(f, explicitTarget); err != nil {
			fmt.Fprintf(os.Stderr, "  \033[1;91m✗ Failed to import %s:\033[0m %v\n", f, err)
			continue
		}

		successCount++
	}

	fmt.Printf("\033[1;92m✓ Chrome Profile Import Complete:\033[0m %d of %d profile(s) imported successfully.\n\n",
		successCount, len(filtered))

	return nil
}

func printChromeProfileImportUsage() {
	fmt.Fprintln(os.Stderr, "  usage: gitmap chrome profile import [file.json|glob|dir|email] [flags]")
	fmt.Fprintln(os.Stderr, "  flags: --except <rule> | --limit <n> | --email <email>")
}

func collectImportCandidates(opts chromeTransferOptions) ([]string, string, error) {
	targetDir := resolveTargetDirFromOpts(opts)
	if opts.Email != "" {
		return findCandidatesByEmail(opts.Email, targetDir)
	}

	if len(opts.Positional) == 0 {
		return findCandidatesInDir(".")
	}

	if len(opts.Positional) == 1 {
		return resolveSinglePositionalCandidate(opts.Positional[0])
	}

	return resolveMultiplePositionalCandidates(opts.Positional)
}

func resolveTargetDirFromOpts(opts chromeTransferOptions) string {
	if len(opts.Positional) > 0 && isDirectoryPath(opts.Positional[0]) {
		return opts.Positional[0]
	}

	return "."
}

func findCandidatesByEmail(email, targetDir string) ([]string, string, error) {
	files := collectSearchSnapshotFiles(targetDir)
	for _, f := range files {
		meta, err := readSnapshotMetadata(f)
		if err == nil && meta.Email != "" && strings.EqualFold(meta.Email, email) {
			return []string{f}, "", nil
		}
	}

	return nil, "", fmt.Errorf("no snapshot file found containing email %q (searched %d files in %s)", email, len(files), targetDir)
}

func collectSearchSnapshotFiles(dir string) []string {
	if dir == "" {
		dir = "."
	}

	files := scanSnapshotFiles(dir)
	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if dir != gitmapChromeDir && chromeProfilePathExists(gitmapChromeDir) {
		files = append(files, scanSnapshotFiles(gitmapChromeDir)...)
	}

	return files
}

func extractCandidatePaths(candidates []DiscoveredProfileCandidate) []string {
	var paths []string
	for _, c := range candidates {
		paths = append(paths, c.SourcePath)
	}

	return paths
}

func findCandidatesInDir(dir string) ([]string, string, error) {
	if dir == "*.*" || dir == "*" {
		dir = "."
	}

	if candidates, err := DiscoverProfileCandidates(dir); err == nil && len(candidates) > 0 {
		return extractCandidatePaths(candidates), "", nil
	}

	files := fallbackCheckFiles(scanSnapshotFiles(dir), dir)
	if len(files) == 0 {
		return nil, "", fmt.Errorf("no profile snapshot files (.json, .zip, .sqlite) found in %q", dir)
	}

	return files, "", nil
}

func resolveSinglePositionalCandidate(pos0 string) ([]string, string, error) {
	if pos0 == "." || pos0 == "*.*" || pos0 == "*" || isDirectoryPath(pos0) {
		return findCandidatesInDir(pos0)
	}

	if strings.ContainsAny(pos0, "*?[") {
		return resolveGlobCandidate(pos0)
	}

	if strings.Contains(pos0, "@") && !isSnapshotFileExtension(pos0) {
		return findCandidatesByEmail(pos0, ".")
	}

	if chromeProfilePathExists(pos0) {
		return []string{pos0}, "", nil
	}

	for _, ext := range []string{constants.ExtZIP, constants.ExtJSON, constants.ExtSQLite} {
		if candExt := pos0 + ext; chromeProfilePathExists(candExt) {
			return []string{candExt}, "", nil
		}
	}

	gitmapFile := filepath.Join(constants.GitMapDir, "chrome", pos0+".json")
	if chromeProfilePathExists(gitmapFile) {
		return []string{gitmapFile}, "", nil
	}

	return nil, "", fmt.Errorf("snapshot file or directory %q not found", pos0)
}

func resolveGlobCandidate(pattern string) ([]string, string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, "", err
	}

	if len(matches) == 0 {
		return nil, "", fmt.Errorf("no files matching pattern %q found", pattern)
	}

	return matches, "", nil
}

func resolveMultiplePositionalCandidates(pos []string) ([]string, string, error) {
	if len(pos) == 2 && isImportableSnapshot(pos[0]) && !isImportableSnapshot(pos[1]) {
		return []string{pos[0]}, pos[1], nil
	}

	var out []string
	for _, p := range pos {
		appendCandidateMatches(p, &out)
	}

	return out, "", nil
}

func appendCandidateMatches(p string, out *[]string) {
	if strings.ContainsAny(p, "*?[") {
		matches, _ := filepath.Glob(p)
		*out = append(*out, matches...)

		return
	}

	if chromeProfilePathExists(p) {
		*out = append(*out, p)
	}
}

func filterImportCandidates(candidates []string, opts chromeTransferOptions) []string {
	var out []string
	for _, f := range candidates {
		meta, err := readSnapshotMetadata(f)
		if err != nil {
			out = append(out, f)
			continue
		}

		if opts.Email != "" && !strings.EqualFold(meta.Email, opts.Email) {
			continue
		}

		if isEx, rule := isExcepted(opts.Except, meta.FileName, meta.ProfileName, meta.DisplayName, meta.Email); isEx {
			fmt.Printf("  \033[1;93m↷ Skipping\033[0m %s (matched --except %q)\n", f, rule)
			continue
		}

		out = append(out, f)
	}

	return out
}
