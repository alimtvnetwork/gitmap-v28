package cmd

// `gitmap undo` — restore the latest `gitmap fix-repo` snapshot for
// the current repo + current version (v5.40.0+).
//
// Usage:
//
//   gitmap undo                  # restore latest snapshot in CWD repo
//   gitmap undo --list           # list snapshots, newest first
//   gitmap undo --snapshot <ts>  # restore a specific UTC timestamp dir
//   gitmap undo --dry-run        # show what would be restored
//
// Backups live at `<repoRoot>/.gitmap/backup/<repo>/v<N>/fix-repo/<ts>/`
// and are produced automatically by every `gitmap fix-repo` write.

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// undoOptions captures parsed CLI flags.
type undoOptions struct {
	isList   bool
	isDryRun bool
	snapshot string
}

// runUndo is the CLI dispatcher.
func runUndo(args []string) {
	checkHelp(constants.CmdUndo, args)
	opts := parseUndoArgs(args)
	identity := resolveFixRepoIdentity()
	baseDir := filepath.Join(identity.root, constants.GitMapDir,
		constants.FixRepoBackupSubdir,
		identity.base+"-v"+strconv.Itoa(identity.current),
		"v"+strconv.Itoa(identity.current),
		constants.CmdFixRepo)
	snapshots := listUndoSnapshots(baseDir)
	if opts.isList {
		printUndoSnapshotList(baseDir, snapshots)
		os.Exit(constants.FixRepoExitOk)
	}
	chosen := pickUndoSnapshot(snapshots, opts.snapshot)
	if chosen == "" {
		fmt.Fprintf(os.Stderr, constants.UndoErrNoSnapshotFmt, baseDir)
		os.Exit(constants.FixRepoExitBadFlag)
	}
	restoreUndoSnapshot(filepath.Join(baseDir, chosen), identity.root, opts.isDryRun)
}

// parseUndoArgs is a tiny flag walker (no os.Args dep). Unknown flags
// abort with E_BAD_FLAG to keep parity with fix-repo's strict parser.
func parseUndoArgs(args []string) undoOptions {
	var opts undoOptions
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--list", "-list", "-l":
			opts.isList = true
		case "--dry-run", "-dry-run", "-DryRun":
			opts.isDryRun = true
		case "--snapshot", "-snapshot", "-s":
			if i+1 < len(args) {
				opts.snapshot = args[i+1]
				i++
			}
		default:
			fmt.Fprintf(os.Stderr, constants.UndoErrBadFlagFmt, args[i])
			os.Exit(constants.FixRepoExitBadFlag)
		}
	}

	return opts
}

// listUndoSnapshots returns timestamp dir names sorted DESC (newest
// first). Missing baseDir → empty slice, never an error.
func listUndoSnapshots(baseDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil
	}
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}
	sort.Sort(sort.Reverse(sort.StringSlice(out)))

	return out
}

// printUndoSnapshotList emits the table for `gitmap undo --list`.
func printUndoSnapshotList(baseDir string, snapshots []string) {
	if len(snapshots) == 0 {
		fmt.Printf(constants.UndoMsgNoSnapshotsFmt, baseDir)

		return
	}
	fmt.Printf(constants.UndoMsgListHeaderFmt, baseDir, len(snapshots))
	for i, ts := range snapshots {
		count := countUndoFiles(filepath.Join(baseDir, ts))
		marker := " "
		if i == 0 {
			marker = "*" // latest
		}
		fmt.Printf(constants.UndoMsgListRowFmt, marker, ts, count)
	}
}

// pickUndoSnapshot returns the explicit snapshot when set, else the
// newest. Empty string means nothing usable.
func pickUndoSnapshot(snapshots []string, explicit string) string {
	if explicit != "" {
		for _, s := range snapshots {
			if s == explicit {
				return s
			}
		}

		return ""
	}
	if len(snapshots) == 0 {
		return ""
	}

	return snapshots[0]
}

// restoreUndoSnapshot reads the manifest and copies each file back.
// Dry-run reports without writing.
func restoreUndoSnapshot(snapDir, repoRoot string, isDryRun bool) {
	manifest, ok := readUndoManifest(snapDir)
	if !ok {
		os.Exit(constants.FixRepoExitBadConfig)
	}
	mode := constants.FixRepoModeWrite
	if isDryRun {
		mode = constants.FixRepoModeDryRun
	}
	fmt.Printf(constants.UndoMsgRestoreHeaderFmt, snapDir, len(manifest.Files), mode)
	restored, failed := walkUndoRestore(snapDir, repoRoot, manifest.Files, isDryRun)
	fmt.Printf(constants.UndoMsgRestoreSummaryFmt, restored, failed)
	if failed > 0 {
		os.Exit(constants.FixRepoExitWriteFailed)
	}
}

// walkUndoRestore copies each manifest entry back. Returns counters.
func walkUndoRestore(snapDir, repoRoot string, files []string, isDryRun bool) (int, int) {
	restored, failed := 0, 0
	for _, rel := range files {
		src := filepath.Join(snapDir, constants.FixRepoBackupFilesSubdir, rel)
		dst := filepath.Join(repoRoot, rel)
		if isDryRun {
			fmt.Printf(constants.UndoMsgRestoreRowFmt, "[dry-run]", rel)
			restored++

			continue
		}
		if err := copyFileForBackup(src, dst); err != nil {
			fmt.Fprintf(os.Stderr, constants.UndoMsgRestoreErrFmt, rel, err)
			failed++

			continue
		}
		fmt.Printf(constants.UndoMsgRestoreRowFmt, "restored", rel)
		restored++
	}

	return restored, failed
}

// readUndoManifest decodes manifest.json. Missing/invalid → false.
func readUndoManifest(snapDir string) (fixRepoBackupManifest, bool) {
	path := filepath.Join(snapDir, constants.FixRepoBackupManifestName)
	body, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.UndoErrManifestMissingFmt, path, err)

		return fixRepoBackupManifest{}, false
	}
	var m fixRepoBackupManifest
	if err := json.Unmarshal(body, &m); err != nil {
		fmt.Fprintf(os.Stderr, constants.UndoErrManifestBadFmt, path, err)

		return fixRepoBackupManifest{}, false
	}

	return m, true
}

// countUndoFiles loads the manifest only to count entries; on any
// failure returns 0 so the listing still renders.
func countUndoFiles(snapDir string) int {
	m, ok := readUndoManifestQuiet(snapDir)
	if !ok {
		return 0
	}

	return len(m.Files)
}

// readUndoManifestQuiet is the noisy-error-free variant used by --list.
func readUndoManifestQuiet(snapDir string) (fixRepoBackupManifest, bool) {
	path := filepath.Join(snapDir, constants.FixRepoBackupManifestName)
	f, err := os.Open(path)
	if err != nil {
		return fixRepoBackupManifest{}, false
	}
	defer f.Close()
	body, err := io.ReadAll(f)
	if err != nil {
		return fixRepoBackupManifest{}, false
	}
	var m fixRepoBackupManifest
	if err := json.Unmarshal(body, &m); err != nil {
		return fixRepoBackupManifest{}, false
	}

	return m, true
}
