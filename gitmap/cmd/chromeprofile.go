// Package cmd — chromeprofile.go: entry points for the Chrome profile
// copy/export/import/list pipeline.
//
//	cpc : copy a profile dir (offline, no sign-in tokens)
//	cpe : export profile to a JSON snapshot
//	cpi : import a JSON snapshot back into a profile dir
//	cpl : list profiles discovered under Chrome User Data
//
// Full spec: spec/04-generic-cli/40-chrome-profile-copy.md.
package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runChromeProfileCopy implements `gitmap chrome-profile-copy`.
func runChromeProfileCopy(args []string) error {
	checkHelp(constants.CmdChromeProfileCopy, args)
	fs := flag.NewFlagSet(constants.CmdChromeProfileCopy, flag.ExitOnError)
	registerOnly := fs.Bool("register-only", false, "skip copy; only (re)register the destination in Chrome's Local State")
	fs.BoolVar(registerOnly, "r", false, "alias for --register-only")
	_ = fs.Parse(args)
	pos := fs.Args()
	if len(pos) < 2 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageCopy)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	srcProfile, ok := resolveChromeProfile(pos[0])
	dstProfile := chromeProfileDestination(pos[1])
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, pos[0], srcProfile.Path)
		printAvailableChromeProfilesWithDisplay()
		os.Exit(constants.ExitChromeProfileNotFound)
	}
	if *registerOnly {
		fmt.Printf(constants.MsgChromeProfileRegOnly, pos[1])
		registerCopiedChromeProfile(srcProfile.Dir, dstProfile.Dir, pos[1])
		rec := emitChromeSnapshots(dstProfile.Path, pos[1])
		persistChromeProfile(pos[1], dstProfile.Path, rec)
		fmt.Printf(constants.MsgChromeProfileNextSteps, pos[1], pos[0], pos[1])
		return nil
	}
	guardChromeClosedOrExit(pos[0], pos[1])
	fmt.Printf(constants.MsgChromeProfileCopyStart, chromeProfileSummary(srcProfile), chromeProfileSummary(dstProfile), srcProfile.Path, dstProfile.Path)
	start := time.Now()
	chromeProfileLockSkipCount = 0
	files, err := copyChromeProfile(srcProfile.Path, dstProfile.Path)
	if err != nil {
		printChromeProfileCopyError(srcProfile, dstProfile, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	if chromeProfileLockSkipCount > 0 {
		fmt.Fprintf(os.Stderr, constants.MsgChromeProfileLockSummary, chromeProfileLockSkipCount)
	}
	fmt.Printf(constants.MsgChromeProfileCopyDone, files, time.Since(start).Round(time.Millisecond))
	registerCopiedChromeProfile(srcProfile.Dir, dstProfile.Dir, pos[1])
	rec := emitChromeSnapshots(dstProfile.Path, pos[1])
	persistChromeProfile(pos[1], dstProfile.Path, rec)
	fmt.Printf(constants.MsgChromeProfileNextSteps, pos[1], pos[0], pos[1])
	return nil
}

// registerCopiedChromeProfile makes the destination directory visible
// in Chrome's profile picker by (1) scrubbing the copied Preferences
// of source signed-in identity + stamping the picker name, and (2)
// adding the dir to Local State `profile.info_cache` + `profiles_order`.
// Step (1) is required: Chrome ignores Local State entries whose
// Preferences still carry the source GAIA fields and silently merges
// the tile back into the source identity on next launch.
func registerCopiedChromeProfile(srcDir, dstDir, displayName string) {
	dstPath := filepath.Join(chromeUserDataDir(), dstDir)
	if err := patchCopiedChromeProfilePreferences(dstPath, displayName); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnChromeProfileRegister, displayName, err)
	}
	if err := registerChromeProfileInLocalState(srcDir, dstDir, displayName); err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnChromeProfileRegister, displayName, err)
		return
	}
	fmt.Printf(constants.MsgChromeProfileRegistered, displayName)
}

func guardChromeClosedOrExit(src, dst string) {
	isRunning, err := isChromeRunning(runtime.GOOS)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.WarnChromeProfileCheckOpen, err)
		fmt.Fprint(os.Stderr, constants.MsgChromeProfileSkipChrome)
		return
	}
	isNonRunning := isRunning == false
	if isNonRunning {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrChromeProfileChromeOpen, src, dst)
	os.Exit(constants.ExitChromeProfileCopyFailed)
}

// emitChromeSnapshots writes the JSON + CSV companions for a profile
// and prints both paths in a consistent Artifacts block. Used by cpc
// and cpe so the output is identical and copy-paste friendly.
func emitChromeSnapshots(srcPath, name string) chromeExportRecord {
	jsonPath := defaultChromeExportPath(name, constants.OutputJSON)
	jsonBytes, err := writeChromeExport(srcPath, name, jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, err)
		return chromeExportRecord{}
	}
	csvPath := jsonPath[:len(jsonPath)-len(constants.ExtJSON)] + constants.ExtCSV
	csvBytes, err := writeChromeExportCSV(srcPath, name, csvPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, err)
		csvPath = ""
	}
	rec := chromeExportRecord{JSONPath: jsonPath, JSONSize: jsonBytes, CSVPath: csvPath, CSVSize: csvBytes}
	printChromeArtifacts(rec)
	return rec
}

// printChromeArtifacts prints the canonical Artifacts: block. Always
// emits both rows so callers can grep `json:`/`csv:` deterministically.
func printChromeArtifacts(rec chromeExportRecord) {
	fmt.Print(constants.MsgChromeProfileArtifactsHd)
	if rec.JSONPath != "" || rec.CSVPath != "" {
		fmt.Printf(constants.MsgChromeProfileArtifactRow, "json:", artifactValue(rec.JSONPath))
		fmt.Printf(constants.MsgChromeProfileArtifactRow, "csv:", artifactValue(rec.CSVPath))
	}
	if rec.ZIPPath != "" {
		fmt.Printf(constants.MsgChromeProfileArtifactRow, "zip:", artifactValue(rec.ZIPPath))
	}
	if rec.SQLitePath != "" {
		fmt.Printf(constants.MsgChromeProfileArtifactRow, "sqlite:", artifactValue(rec.SQLitePath))
	}
}

func artifactValue(path string) string {
	if path == "" {
		return constants.MsgChromeProfileArtifactNA
	}
	return path
}

// runChromeProfileExport implements `gitmap chrome-profile-export`.
func runChromeProfileExport(args []string) error {
	checkHelp(constants.CmdChromeProfileExport, args)

	format := constants.OutputJSON
	var positionalArgs []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--format=") {
			format = strings.TrimPrefix(arg, "--format=")
		} else {
			positionalArgs = append(positionalArgs, arg)
		}
	}

	if len(positionalArgs) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageExport)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	name := positionalArgs[0]
	outPath := defaultChromeExportPath(name, format)
	if len(positionalArgs) >= 2 {
		outPath = positionalArgs[1]
	}
	srcPath, ok := resolveChromeProfileDir(name)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, name, srcPath)
		printAvailableChromeProfilesWithDisplay()
		os.Exit(constants.ExitChromeProfileNotFound)
	}

	// We will hand off to the actual export logic based on format.
	// For now, it delegates to the existing JSON/CSV exporter.
	// Later steps will inject the ZIP and SQLite code.
	var jsonBytes, csvBytes int
	var err, csvErr error
	var csvPath string

	if format == constants.OutputJSON || format == constants.OutputCSV {
		jsonBytes, err = writeChromeExport(srcPath, name, outPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, err)
			os.Exit(constants.ExitChromeProfileCopyFailed)
		}
		csvPath = outPath
		if ext := constants.ExtJSON; len(csvPath) > len(ext) && csvPath[len(csvPath)-len(ext):] == ext {
			csvPath = csvPath[:len(csvPath)-len(ext)] + constants.ExtCSV
		} else {
			csvPath += constants.ExtCSV
		}
		csvBytes, csvErr = writeChromeExportCSV(srcPath, name, csvPath)
		if csvErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, csvErr)
			csvPath = ""
		}
	} else if format == constants.OutputZIP {
		var zipErr error
		jsonBytes, zipErr = writeChromeExportZIP(srcPath, name, outPath)
		if zipErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, zipErr)
			os.Exit(constants.ExitChromeProfileCopyFailed)
		}
	} else if format == constants.OutputSQLite {
		var sqErr error
		jsonBytes, sqErr = writeChromeExportSQLite(srcPath, name, outPath)
		if sqErr != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, sqErr)
			os.Exit(constants.ExitChromeProfileCopyFailed)
		}
	}

	rec := chromeExportRecord{
		JSONPath: outPath, JSONSize: jsonBytes,
		CSVPath: csvPath, CSVSize: csvBytes,
		ZIPPath: outPath, ZIPSize: jsonBytes,
		SQLitePath: outPath, SQLiteSize: jsonBytes,
	}
	// For ZIP/SQLite formats, outPath represents the ZIP/SQLite file, and we use JSONSize as a generic ByteSize field for now.
	if format == constants.OutputJSON || format == constants.OutputCSV {
		rec.ZIPPath = ""
		rec.ZIPSize = 0
		rec.SQLitePath = ""
		rec.SQLiteSize = 0
	} else if format == constants.OutputZIP {
		rec.JSONPath = ""
		rec.CSVPath = ""
		rec.SQLitePath = ""
	} else if format == constants.OutputSQLite {
		rec.JSONPath = ""
		rec.CSVPath = ""
		rec.ZIPPath = ""
	}

	printChromeArtifacts(rec)
	persistChromeProfile(name, srcPath, rec)
	return nil
}

// runChromeProfileImport implements `gitmap chrome-profile-import`.
func runChromeProfileImport(args []string) error {
	checkHelp(constants.CmdChromeProfileImport, args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageImport)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	srcFile := args[0]
	name := ""
	if len(args) >= 2 {
		name = args[1]
	}

	if strings.HasSuffix(srcFile, ".zip") || strings.HasSuffix(srcFile, ".sqlite") {
		if name == "" {
			// guess name from zip filename without extension
			base := filepath.Base(srcFile)
			name = strings.TrimSuffix(base, filepath.Ext(base))
		}
		dstPath := chromeProfilePath(name)
		if err := applyChromeExportZIP(srcFile, dstPath); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
			os.Exit(constants.ExitChromeProfileCopyFailed)
		}
		fmt.Printf(constants.MsgChromeProfileImportOk, srcFile, name)
		return nil
	}

	exp, err := loadChromeImport(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	if name == "" {
		name = exp.Name
	}
	dstPath := chromeProfilePath(name)
	if err := applyChromeExport(exp, dstPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	fmt.Printf(constants.MsgChromeProfileImportOk, srcFile, name)
	return nil
}

// runChromeProfileList implements `gitmap chrome-profile-list`.
func runChromeProfileList(args []string) error {
	checkHelp(constants.CmdChromeProfileList, args)
	root := chromeUserDataDir()
	entries := chromeProfileEntries()
	if len(entries) == 0 {
		fmt.Printf(constants.MsgChromeProfileListEmpty, root)
		listChromeProfilesFromDB()
		return nil
	}
	fmt.Printf(constants.MsgChromeProfileListHdr, root)
	for _, e := range entries {
		if e.DisplayName != "" {
			fmt.Printf("  - %s  (display: %q)\n", e.Dir, e.DisplayName)
			continue
		}
		fmt.Printf("  - %s\n", e.Dir)
	}
	listChromeProfilesFromDB()
	return nil
}

// defaultChromeExportPath builds the default output location
// under .gitmap/chrome/<name>.<format> (cwd-relative).
func defaultChromeExportPath(name, format string) string {
	ext := constants.ExtJSON
	if format == constants.OutputZIP {
		ext = ".zip"
	} else if format == constants.OutputSQLite {
		ext = ".sqlite"
	}
	return filepath.Join(constants.GitMapDir, "chrome", name+ext)
}

// readChromeExport loads a JSON export file from disk.
func readChromeExport(path string) (*chromeExport, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &exp, nil
}
