package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/config"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/detector"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/mapper"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/scanner"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runScan handles the "scan" subcommand.
func runScan(args []string) {
	checkHelp("scan", args)
	dir, cfgPath, mode, output, outFile, outputPath, relativeRoot, defaultBranch, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, workers, maxDepth, probeOpts := parseScanFlags(args)
	cfg, err := config.LoadFromFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrConfigLoad, cfgPath, err)
		os.Exit(1)
	}
	cfg = config.MergeWithFlags(cfg, mode, output, outputPath)
	cache := model.ScanCache{
		Dir: dir, ConfigPath: cfgPath, Mode: mode, Output: output,
		OutFile: outFile, OutputPath: outputPath,
		GithubDesktop: ghDesktop, OpenFolder: openFolder, Quiet: quiet,
	}
	executeScan(dir, cfg, outFile, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors, workers, maxDepth, cache, probeOpts, relativeRoot, defaultBranch)
}

// executeScan performs the directory scan and outputs results.
//
// Each phase is wrapped in a benchmark.Phase call so that
// .gitmap/output/scan-benchmark.log captures wall-clock timings for every
// stage. This is the file users should attach when reporting "scan is
// slow" — it pinpoints which phase (walk, DB upsert, project detection,
// release import, desktop sync, …) actually consumed the time.
func executeScan(dir string, cfg model.Config, outFile string, ghDesktop, openFolder, quiet, noVSCodeSync, noAutoTags, reportErrors bool, workers, maxDepth int, cache model.ScanCache, probeOpts ScanProbeOptions, relativeRoot, defaultBranch string) {
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
	var repos []scanner.RepoInfo
	var err error
	bench.Phase("scan.walk", func() {
		repos, err = scanner.ScanDirWithOptions(absDir, scanner.ScanOptions{
			ExcludeDirs: cfg.ExcludeDirs,
			Workers:     workers,
			MaxDepth:    maxDepth,
			Progress:    progress.Callback(),
			OnDirError:  scanDirErrorCallback(errCollector),
		})
	})
	if err != nil {
		failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrScanFailed, absDir, err))
		fmt.Fprintf(os.Stderr, constants.ErrScanFailed, absDir, err)
		exitWith(1)
	}
	var records []model.ScanRecord
	relRootBase := resolveRelativeRoot(relativeRoot, absDir, quiet)
	bench.Phase("scan.buildRecords", func() {
		records = mapper.BuildRecordsWithOptions(repos, mapper.BuildOptions{
			Mode:          cfg.DefaultMode,
			DefaultNote:   cfg.Notes,
			RelRoot:       relRootBase,
			DefaultBranch: defaultBranch,
		})
	})
	outputDir := resolveOutputDir(cfg.OutputDir, absDir)
	fmt.Printf(constants.MsgSectionArtifacts, outputDir)
	bench.Phase("scan.writeOutputs", func() {
		writeAllOutputs(records, outputDir, outFile, quiet)
	})
	bench.Phase("scan.saveCache", func() {
		saveScanCache(outputDir, cache)
	})
	fmt.Print(constants.MsgSectionDatabase)
	bench.Phase("scan.dbUpsertRepos", func() {
		upsertToDB(records, outputDir)
	})
	bench.Phase("scan.tagScanFolder", func() {
		tagReposWithScanFolder(absDir, records, quiet)
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

	// Mark scan task as completed after all steps succeed.
	completePendingTask(taskDB, taskID)
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

	idsByPath := make(map[string]int64, len(repos))
	for _, repo := range repos {
		idsByPath[repo.AbsolutePath] = repo.ID
	}

	aligned := make([]model.ScanRecord, 0, len(records))
	for _, rec := range records {
		if id, ok := idsByPath[rec.AbsolutePath]; ok {
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
