package cmdscan

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/detector"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitignoreagm"
	"github.com/alimtvnetwork/gitmap-v28/cli/scanpipe"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CheckSpecialReposOnScanFn is an optional callback to verify and initialize special repos (rs/rc) on scan.
var CheckSpecialReposOnScanFn func(workBaseDir string, isQuiet bool)

// runScan handles the "scan" subcommand.
func runScan(args []string) error {
	checkHelp("scan", args)
	dir, cfgPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch, forceInclude, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, compact, fix, workers, maxDepth, probeOpts := parseScanFlags(args)
	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrConfigLoad)
	}

	cfg = config.MergeWithFlags(cfg, mode, output, outputPath)
	cache := model.ScanCache{
		Dir: dir, ConfigPath: cfgPath, Mode: mode, Output: output,
		OutFile: outFile, OutputPath: outputPath,
		IsGithubDesktop: ghDesktop, IsOpenFolder: openFolder, IsQuiet: quiet,
	}

	forceIncludeDirs := parseForceIncludeDirs(forceInclude)

	return executeScan(dir, cfg, outFile, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, compact, fix, workers, maxDepth, cache, probeOpts, relativeRoot, defaultBranch, forceIncludeDirs)
}

// executeScan performs the directory scan and outputs results.
//
// Each phase is wrapped in a benchmark.Phase call so that
// .gitmap/output/scan-benchmark.log captures wall-clock timings for every
// stage. This is the file users should attach when reporting "scan is
// slow" — it pinpoints which phase (walk, DB upsert, project detection,
// release import, desktop sync, …) actually consumed the time.
func executeScan(
	dir string,
	cfg model.Config,
	outFile string,
	ghDesktop,
	openFolder,
	quiet,
	noVSCodeSync,
	noAutoTags,
	reportErrors,
	compact,
	fix bool,
	workers,
	maxDepth int,
	cache model.ScanCache,
	probeOpts ScanProbeOptions,
	relativeRoot,
	defaultBranch string,
	forceIncludeDirs []string,
) error {
	absDir := resolveScanTarget(dir)

	bench := newScanBenchmark(absDir)
	if !quiet {
		fmt.Printf("  ▶ gitmap scan v%s — %s\n", constants.Version, absDir)
	}

	// Enqueue scan as a pending task before execution.
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}

	cmdArgs := buildCommandArgs(append([]string{"scan"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeScan, absDir, workDir, "scan", cmdArgs)
	if taskDB != nil {
		defer taskDB.Close()
	}

	progress := newScanProgressRenderer(quiet)
	errCollector := newScanCollector(reportErrors)
	var repos []scanpipe.RepoInfo
	var err error
	bench.Phase("scan.walk", func() {
		repos, err = scanpipe.ScanDirWithOptions(absDir, scanpipe.ScanOptions{
			ExcludeDirs:      cfg.ExcludeDirs,
			ForceIncludeDirs: forceIncludeDirs,
			Workers:          workers,
			MaxDepth:         maxDepth,
			Progress:         progress.Callback(),
			OnDirError:       scanDirErrorCallback(errCollector),
		})
	})
	if err != nil {
		failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrScanFailed, absDir, err))
		fmt.Fprintf(os.Stderr, constants.ErrScanFailed, absDir, err)

		return apperror.WrapSimple(err, fmt.Sprintf(constants.ErrScanFailed, absDir, err))
	}

	var records []model.ScanRecord
	relRootBase := resolveRelativeRoot(relativeRoot, absDir, quiet)
	bench.Phase("scan.buildRecords", func() {
		records = scanpipe.BuildRecordsWithOptions(repos, scanpipe.BuildOptions{
			Mode:          cfg.DefaultMode,
			DefaultNote:   cfg.Notes,
			RelRoot:       relRootBase,
			DefaultBranch: defaultBranch,
		})
		records = enrichRecordsWithAliases(records)
	})
	outputDir := resolveOutputDir(cfg.OutputDir, absDir)
	bench.Phase("scan.writeOutputs", func() {
		writeAllOutputs(records, outputDir, outFile, quiet, compact)
	})
	bench.Phase("scan.saveCache", func() {
		saveScanCache(outputDir, cache)
	})
	fmt.Print(constants.MsgSectionDatabase)
	if fix {
		bench.Phase("scan.pruneStaleDB", func() {
			runPruneStaleDB(absDir, records)
		})
	}

	bench.Phase("scan.dbUpsertRepos", func() {
		upsertToDB(records, outputDir)
	})
	bench.Phase("scan.tagScanFolder", func() {
		tagReposWithScanFolder(absDir, records, quiet)
	})
	var wasFirstWorkDirRegistered bool
	bench.Phase("scan.autoRegisterWorkDir", func() {
		wasFirstWorkDirRegistered = autoRegisterFirstWorkDir(absDir, quiet)
	})
	bench.Phase("scan.autoAliases", func() {
		autoPopulateScanAliases(quiet)
	})
	bench.Phase("scan.specialReposCheck", func() {
		if CheckSpecialReposOnScanFn != nil {
			CheckSpecialReposOnScanFn(absDir, quiet)
		}
	})
	bench.Phase("scan.gitignoreAgmCheck", func() {
		checkAgmResumeTaskOnScan(records, quiet, fix)
	})
	bench.Phase("scan.alignDBIDs", func() {
		records = alignRecordsWithDB(records, outputDir)
	})
	// Background probe: kicked off here so it runs concurrently with
	// the project-detection / desktop-sync phases below. Drained
	// before the "Done" banner unless --no-probe-wait was passed.
	probeRunner := startBackgroundProbe(records, probeOpts, quiet, errCollector)
	fmt.Print(constants.MsgSectionProjects)
	var detected []detector.DetectionResult
	bench.Phase("scan.detectProjects", func() {
		detected = detectAllProjects(records)
	})
	bench.Phase("scan.writeProjectJSON", func() {
		writeProjectJSONFiles(detected, outputDir)
	})
	bench.Phase("scan.dbUpsertProjects", func() {
		upsertProjectsToDB(detected, records, outputDir)
	})
	bench.Phase("scan.importReleases", func() {
		importReleases(absDir, outputDir)
	})
	bench.Phase("scan.addToDesktop", func() {
		addToDesktop(records, ghDesktop)
	})
	bench.Phase("scan.vscodePMSync", func() {
		syncRecordsToVSCodePM(records, noVSCodeSync, noAutoTags)
	})
	bench.Phase("scan.ideSync", func() {
		syncRecordsToIDEs(records, os.Args, noVSCodeSync, quiet)
	})
	openOutputFolder(outputDir, openFolder)
	bench.WriteLog(outputDir)
	if !quiet {
		fmt.Printf("  📊 Benchmark log: %s\n", filepath.Join(outputDir, scanBenchmarkFile))
	}
	bench.Phase("scan.backgroundProbeWait", func() {
		drainBackgroundProbe(probeRunner, probeOpts, quiet)
	})
	finalizeErrorReport(errCollector, quiet)
	fmt.Print(constants.MsgSectionDone)
	if wasFirstWorkDirRegistered && !quiet {
		fmt.Printf("  ✓ First work directory registered and marked as default: %s\n", absDir)
	}
	// Mark scan task as completed after all steps succeed.
	completePendingTask(taskDB, taskID)
	return nil
}
// autoRegisterFirstWorkDir checks if any work directories are registered.
// If none exist, it registers absDir as the first work directory and marks it as default.
func autoRegisterFirstWorkDir(absDir string, quiet bool) bool {
	db, err := store.OpenDefault()
	if err != nil {
		return false
	}
	defer db.Close()
	dirs, errList := db.ListWorkDirs()
	if errList != nil || len(dirs) > 0 {
		return false
	}
	label := filepath.Base(absDir)
	if _, errEnsure := db.EnsureWorkDir(absDir, label, true); errEnsure != nil {
		return false
	}
	_ = db.SetDefaultWorkDir(absDir)
	return true
}
// tagReposWithScanFolder registers absDir as a ScanFolder and tags every
// just-scanned repo with the resulting ScanFolderId. Failures are reported
// to stderr but do NOT fail the scan — the underlying Repo rows still exist.
func tagReposWithScanFolder(absDir string, records []model.ScanRecord, quiet bool) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProbeOpenDB, err)
		return
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	folder, err := db.EnsureScanFolder(absDir, "", "")
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	paths := make([]string, 0, len(records))
	for _, r := range records {
		paths = append(paths, r.AbsolutePath)
	}
	if err := db.TagReposByScanFolder(folder.ID, paths); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
	if !quiet {
		fmt.Printf(constants.MsgScanFolderTagged, len(paths), folder.ID)
	}
}
// upsertToDB persists scan results into the SQLite database.
func upsertToDB(records []model.ScanRecord, outputDir string) {
	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		return
	}
	defer db.Close()
	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		return
	}
	if err := db.UpsertRepos(records); err != nil {
		fmt.Fprintf(os.Stderr, constants.MsgDBUpsertFailed, err)
		return
	}
	fmt.Printf(constants.MsgDBUpsertDone, len(records))
}
// alignRecordsWithDB rewrites record IDs to match persisted repo IDs by path.
func alignRecordsWithDB(records []model.ScanRecord, outputDir string) []model.ScanRecord {
	db, err := store.OpenDefault()
	if err != nil {
		return records
	}
	defer db.Close()
	repos, err := db.ListRepos()
	if err != nil {
		return records
	}
	idsByPath := make(map[string]int64, len(repos)*2)
	for _, repo := range repos {
		idsByPath[repo.AbsolutePath] = repo.ID
		norm := store.NormalizeStoragePath(repo.AbsolutePath)
		idsByPath[norm] = repo.ID
		idsByPath[strings.ToLower(norm)] = repo.ID
	}
	aligned := make([]model.ScanRecord, 0, len(records))
	for _, rec := range records {
		norm := store.NormalizeStoragePath(rec.AbsolutePath)
		lowerNorm := strings.ToLower(norm)
		if id, ok := idsByPath[lowerNorm]; ok {
			rec.ID = id
		} else if id, ok := idsByPath[norm]; ok {
			rec.ID = id
		} else if id, ok := idsByPath[rec.AbsolutePath]; ok {
			rec.ID = id
		}
		aligned = append(aligned, rec)
	}
	return aligned
}
// addToDesktop registers repos with GitHub Desktop if requested.
func addToDesktop(records []model.ScanRecord, enabled bool) {
	if enabled {
		summary := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, summary.Added, summary.Failed)
	}
}
// openOutputFolder opens the output directory in the OS file explorer.
func openOutputFolder(outputDir string, enabled bool) {
	if enabled {
		cmd := resolveOpenCommand(outputDir)
		_ = cmd.Start()
		fmt.Printf(constants.MsgOpenedFolder, outputDir)
	}
}
// resolveOpenCommand returns the OS-specific command to open a folder.
func resolveOpenCommand(dir string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return exec.Command(constants.CmdExplorer, dir)
	}
	if runtime.GOOS == constants.OSDarwin {
		return exec.Command(constants.CmdOpen, dir)
	}
	return exec.Command(constants.CmdXdgOpen, dir)
}
// resolveOutputDir determines the output directory relative to scan root.
func resolveOutputDir(cfgDir, scanDir string) string {
	if filepath.IsAbs(cfgDir) {
		return cfgDir
	}
	return filepath.Join(scanDir, constants.GitMapDir, constants.OutputDirName)
}
// autoPopulateScanAliases generates aliases for unaliased repositories found during scan.
func autoPopulateScanAliases(quiet bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return
	}
	defer db.Close()
	count, _ := cmdinstall.EnsureTrackedRepoAliases(db)
	if count > 0 && !quiet {
		fmt.Printf("  ✓ Auto-generated %d repository alias(es)\n", count)
	}
}
func checkAgmResumeTaskOnScan(records []model.ScanRecord, quiet, fix bool) {
	paths := make([]string, 0, len(records))
	for _, r := range records {
		paths = append(paths, r.AbsolutePath)
	}
	isAuto := fix || gitignoreagm.IsAutoRemediateScanEnabled()
	_ = gitignoreagm.CheckAndPromptRepos(paths, quiet, isAuto)
}
func parseForceIncludeDirs(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	parts := strings.Split(trimmed, ",")
	var result []string
	for _, p := range parts {
		item := strings.TrimSpace(p)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}