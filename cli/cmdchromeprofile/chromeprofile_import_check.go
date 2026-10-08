package cmdchromeprofile

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func collectSnapshotFileCandidates(target string) []DiscoveredProfileCandidate {
	var candidates []DiscoveredProfileCandidate
	files, _ := resolveSnapshotCheckFiles(target)
	for _, f := range files {
		c, cErr := discoverSingleFileCandidate(f)
		if cErr != nil {
			continue
		}

		candidates = append(candidates, c...)
	}

	return candidates
}

func runChromeProfileImportCheck(args []string) error {
	checkHelp("import-check", args)
	opts := parseChromeTransferOptions(args)
	target := resolveCheckTarget(args)

	candidates, err := DiscoverProfileCandidates(target)
	if err != nil || len(candidates) == 0 {
		candidates = collectSnapshotFileCandidates(target)
	}

	if len(candidates) == 0 && opts.Fnf {
		return fmt.Errorf("no profile snapshot files found to check in %q (--fnf asserted)", target)
	}

	if len(candidates) == 0 {
		fmt.Printf("No snapshot files found to check in %q\n", target)

		return nil
	}

	params := ChromePreviewOutputParams{
		Candidates: candidates,
		Target:     target,
		IsJSON:     opts.IsJSON || hasJSONFlag(args),
		FilePath:   opts.FilePath,
		TempFile:   opts.TempFile,
		Fnf:        opts.Fnf,
	}

	return dispatchPreviewOutput(params)
}

func resolveCheckTarget(args []string) string {
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" || a == "-j" || a == "--fnf" || isHelpFlag(a) || isPreflightInspectArg(a) {
			continue
		}

		if isValuedFlag(a) && !strings.Contains(a, "=") {
			i++
			continue
		}

		if isValuedFlag(a) {
			continue
		}

		if a == "*.*" || a == "*" {
			return "."
		}

		return a
	}

	return "."
}

func resolveSnapshotCheckFiles(target string) ([]string, error) {
	info, err := os.Stat(target)
	if err == nil && !info.IsDir() {
		return []string{target}, nil
	}

	files := scanSnapshotFiles(target)
	files = fallbackCheckFiles(files, target)

	return files, nil
}

func fallbackCheckFiles(files []string, target string) []string {
	if len(files) > 0 || target != "." {
		return files
	}

	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if !chromeProfilePathExists(gitmapChromeDir) {
		return files
	}

	return scanSnapshotFiles(gitmapChromeDir)
}

func listDiscoveredSnapshotsInDir(dir string) {
	files := scanSnapshotFiles(dir)
	files, dir = resolveFallbackSnapshots(files, dir)
	if len(files) == 0 {
		return
	}

	fmt.Printf("\n\033[1;96mDiscovered Backup / Snapshot Files in %s:\033[0m\n", dir)
	for i, f := range files {
		meta, err := readSnapshotMetadata(f)
		if err != nil {
			fmt.Printf("  [%d] %s\n", i+1, filepath.Base(f))
			continue
		}

		emailStr := meta.Email
		if emailStr == "" {
			emailStr = "(none)"
		}

		fmt.Printf("  [%d] \033[1;97m%-16s\033[0m (Display: %q, Email: %s, Bookmarks: %d, Ext: %d)\n",
			i+1, meta.FileName, meta.DisplayName, emailStr, meta.BookmarksCount, meta.ExtensionsCount)
	}

	fmt.Println()
}

func resolveFallbackSnapshots(files []string, dir string) ([]string, string) {
	if len(files) > 0 {
		return files, dir
	}

	gitmapChromeDir := filepath.Join(constants.GitMapDir, "chrome")
	if dir == gitmapChromeDir || !chromeProfilePathExists(gitmapChromeDir) {
		return files, dir
	}

	return scanSnapshotFiles(gitmapChromeDir), gitmapChromeDir
}
