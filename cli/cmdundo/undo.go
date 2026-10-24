package cmdundo

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdundo/enums"
)

// undoOptions captures parsed CLI flags.
type undoOptions struct {
	isList   bool
	isDryRun bool
	snapshot string
}

// RunUndo is the CLI dispatcher for restoring backup snapshots.
func RunUndo(args []string) error {
	checkHelp(constants.CmdUndo, args)

	opts, err := parseUndoArgs(args)
	if err != nil {
		return err
	}

	return executeUndo(opts)
}

func executeUndo(opts undoOptions) error {
	identity := cmdfixrepo.ResolveFixRepoIdentity()
	baseDir := buildUndoBaseDir(identity)
	snapshots := listUndoSnapshots(baseDir)
	if opts.isList {
		printUndoSnapshotList(baseDir, snapshots)

		return nil
	}

	return executeUndoRestore(snapshots, opts, baseDir, identity.Root)
}

func buildUndoBaseDir(identity cmdfixrepo.FixRepoIdentity) string {
	versionTag := "v" + strconv.Itoa(identity.Current)

	return filepath.Join(identity.Root, constants.GitMapDir,
		constants.SubdirBackup,
		identity.Base+"-"+versionTag,
		versionTag,
		constants.CmdFixRepo)
}

func executeUndoRestore(snapshots []string, opts undoOptions, baseDir, repoRoot string) error {
	chosen := pickUndoSnapshot(snapshots, opts.snapshot)
	if chosen == "" {
		fmt.Fprintf(os.Stderr, constants.ErrNoSnapshotFmt, baseDir)

		return apperror.NewValidationError(fmt.Sprintf(constants.ErrNoSnapshotFmt, baseDir))
	}

	snapDir := filepath.Join(baseDir, chosen)

	return restoreUndoSnapshot(snapDir, repoRoot, opts.isDryRun)
}

func parseUndoArgs(args []string) (undoOptions, error) {
	var opts undoOptions
	for i := 0; i < len(args); i++ {
		nextI, err := applyUndoArg(args, i, &opts)
		if err != nil {
			return undoOptions{}, err
		}

		i = nextI
	}

	return opts, nil
}

func applyUndoArg(args []string, i int, opts *undoOptions) (int, error) {
	arg := args[i]
	switch arg {
	case constants.FlagListLong, constants.FlagListShort1, constants.FlagListShort2:
		opts.isList = true
	case constants.FlagDryRunLong, constants.FlagDryRunShort1, constants.FlagDryRunShort2:
		opts.isDryRun = true
	case constants.FlagSnapshotLong, constants.FlagSnapshotShort1, constants.FlagSnapshotShort2:
		return parseSnapshotFlag(args, i, opts)
	default:
		return handleUnknownArg(arg, i)
	}

	return i, nil
}

func handleUnknownArg(arg string, i int) (int, error) {
	fmt.Fprintf(os.Stderr, constants.ErrBadFlagFmt, arg)

	return i, apperror.NewValidationError(fmt.Sprintf(constants.ErrBadFlagFmt, arg))
}

func parseSnapshotFlag(args []string, i int, opts *undoOptions) (int, error) {
	if i+1 < len(args) {
		opts.snapshot = args[i+1]

		return i + 1, nil
	}

	return i, nil
}

func listUndoSnapshots(baseDir string) []string {
	entries, err := os.ReadDir(baseDir)
	if err != nil {
		return nil
	}

	out := filterDirNames(entries)
	sort.Sort(sort.Reverse(sort.StringSlice(out)))

	return out
}

func filterDirNames(entries []os.DirEntry) []string {
	out := make([]string, 0, len(entries))
	for _, e := range entries {
		if e.IsDir() {
			out = append(out, e.Name())
		}
	}

	return out
}

func printUndoSnapshotList(baseDir string, snapshots []string) {
	if len(snapshots) == 0 {
		fmt.Printf(constants.MsgNoSnapshotsFmt, baseDir)

		return
	}

	fmt.Printf(constants.MsgListHeaderFmt, baseDir, len(snapshots))
	renderSnapshotRows(baseDir, snapshots)
}

func renderSnapshotRows(baseDir string, snapshots []string) {
	for i, ts := range snapshots {
		count := countUndoFiles(filepath.Join(baseDir, ts))
		marker := snapshotMarker(i)
		fmt.Printf(constants.MsgListRowFmt, marker, ts, count)
	}
}

func snapshotMarker(index int) string {
	if index == 0 {
		return constants.MarkerLatest
	}

	return constants.MarkerOlder
}

func pickUndoSnapshot(snapshots []string, explicit string) string {
	if explicit != "" {
		return findExplicitSnapshot(snapshots, explicit)
	}

	if len(snapshots) == 0 {
		return ""
	}

	return snapshots[0]
}

func findExplicitSnapshot(snapshots []string, explicit string) string {
	for _, s := range snapshots {
		if s == explicit {
			return s
		}
	}

	return ""
}

func restoreUndoSnapshot(snapDir, repoRoot string, isDryRun bool) error {
	manifest, hasManifest := readUndoManifest(snapDir)
	if !hasManifest {
		return apperror.NewValidationError(fmt.Sprintf(constants.ErrManifestMissingFmt, snapDir, "missing or invalid manifest"))
	}

	mode := resolveRestoreMode(isDryRun)
	fmt.Printf(constants.MsgRestoreHeaderFmt, snapDir, len(manifest.Files), mode)
	restored, failed := walkUndoRestore(snapDir, repoRoot, manifest.Files, isDryRun)
	fmt.Printf(constants.MsgRestoreSummaryFmt, restored, failed)
	if failed > 0 {
		return apperror.NewExecutionError(fmt.Sprintf("undo: restore failed with %d error(s)", failed))
	}

	return nil
}

func resolveRestoreMode(isDryRun bool) string {
	if isDryRun {
		return constants.ModeDryRun
	}

	return constants.ModeWrite
}

func walkUndoRestore(snapDir, repoRoot string, files []string, isDryRun bool) (int, int) {
	restored, failed := 0, 0
	for _, rel := range files {
		isOk := restoreSingleEntry(snapDir, repoRoot, rel, isDryRun)
		if isOk {
			restored++
		} else {
			failed++
		}
	}

	return restored, failed
}

func restoreSingleEntry(snapDir, repoRoot, rel string, isDryRun bool) bool {
	if isDryRun {
		fmt.Printf(constants.MsgRestoreRowFmt, enums.UndoActionDryRun.String(), rel)

		return true
	}

	return copyAndLogEntry(snapDir, repoRoot, rel)
}

func copyAndLogEntry(snapDir, repoRoot, rel string) bool {
	src := filepath.Join(snapDir, constants.SubdirFiles, rel)
	dst := filepath.Join(repoRoot, rel)
	if err := cmdfixrepo.CopyFileForBackup(src, dst); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgRestoreErrFmt, rel, err)

		return false
	}

	fmt.Printf(constants.MsgRestoreRowFmt, enums.UndoActionRestored.String(), rel)

	return true
}

func readUndoManifest(snapDir string) (cmdfixrepo.FixRepoBackupManifest, bool) {
	path := filepath.Join(snapDir, constants.ManifestName)
	body, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrManifestMissingFmt, path, err)

		return cmdfixrepo.FixRepoBackupManifest{}, false
	}

	return unmarshalUndoManifest(body, path)
}

func unmarshalUndoManifest(body []byte, path string) (cmdfixrepo.FixRepoBackupManifest, bool) {
	var m cmdfixrepo.FixRepoBackupManifest
	if err := json.Unmarshal(body, &m); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrManifestBadFmt, path, err)

		return cmdfixrepo.FixRepoBackupManifest{}, false
	}

	return m, true
}

func countUndoFiles(snapDir string) int {
	m, hasManifest := readUndoManifestQuiet(snapDir)
	if !hasManifest {
		return 0
	}

	return len(m.Files)
}

func readUndoManifestQuiet(snapDir string) (cmdfixrepo.FixRepoBackupManifest, bool) {
	path := filepath.Join(snapDir, constants.ManifestName)
	f, err := os.Open(path)
	if err != nil {
		return cmdfixrepo.FixRepoBackupManifest{}, false
	}

	defer f.Close()

	return readManifestFromReader(f)
}

func readManifestFromReader(r io.Reader) (cmdfixrepo.FixRepoBackupManifest, bool) {
	body, err := io.ReadAll(r)
	if err != nil {
		return cmdfixrepo.FixRepoBackupManifest{}, false
	}

	var m cmdfixrepo.FixRepoBackupManifest
	if err := json.Unmarshal(body, &m); err != nil {
		return cmdfixrepo.FixRepoBackupManifest{}, false
	}

	return m, true
}
