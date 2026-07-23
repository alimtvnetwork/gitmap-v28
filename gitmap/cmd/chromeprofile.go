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
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runChromeProfileCopy implements `gitmap chrome-profile-copy`.
func runChromeProfileCopy(args []string) {
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
		return
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
	if !isRunning {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrChromeProfileChromeOpen, src, dst)
	os.Exit(constants.ExitChromeProfileCopyFailed)
}

// emitChromeSnapshots writes the JSON + CSV companions for a profile
// and prints both paths in a consistent Artifacts block. Used by cpc
// and cpe so the output is identical and copy-paste friendly.
func emitChromeSnapshots(srcPath, name string) chromeExportRecord {
	jsonPath := defaultChromeExportPath(name)
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
	fmt.Printf(constants.MsgChromeProfileArtifactRow, "json:", artifactValue(rec.JSONPath))
	fmt.Printf(constants.MsgChromeProfileArtifactRow, "csv:", artifactValue(rec.CSVPath))
}

func artifactValue(path string) string {
	if path == "" {
		return constants.MsgChromeProfileArtifactNA
	}
	return path
}

// runChromeProfileExport implements `gitmap chrome-profile-export`.
func runChromeProfileExport(args []string) {
	checkHelp(constants.CmdChromeProfileExport, args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageExport)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	name := args[0]
	outPath := defaultChromeExportPath(name)
	if len(args) >= 2 {
		outPath = args[1]
	}
	srcPath, ok := resolveChromeProfileDir(name)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, name, srcPath)
		printAvailableChromeProfilesWithDisplay()
		os.Exit(constants.ExitChromeProfileNotFound)
	}
	jsonBytes, err := writeChromeExport(srcPath, name, outPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	csvPath := outPath
	if ext := constants.ExtJSON; len(csvPath) > len(ext) && csvPath[len(csvPath)-len(ext):] == ext {
		csvPath = csvPath[:len(csvPath)-len(ext)] + constants.ExtCSV
	} else {
		csvPath += constants.ExtCSV
	}
	csvBytes, csvErr := writeChromeExportCSV(srcPath, name, csvPath)
	if csvErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, csvErr)
		csvPath = ""
	}
	rec := chromeExportRecord{
		JSONPath: outPath, JSONSize: jsonBytes,
		CSVPath: csvPath, CSVSize: csvBytes,
	}
	printChromeArtifacts(rec)
	persistChromeProfile(name, srcPath, rec)
}

// runChromeProfileImport implements `gitmap chrome-profile-import`.
// Accepts both .json (full snapshot) and .csv (lossy: extension IDs +
// known preferences only — bookmarks omitted).
func runChromeProfileImport(args []string) {
	checkHelp(constants.CmdChromeProfileImport, args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageImport)
		os.Exit(constants.ExitChromeProfileUsage)
	}
	srcFile := args[0]
	exp, err := loadChromeImport(srcFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	name := exp.Name
	if len(args) >= 2 {
		name = args[1]
	}
	dstPath := chromeProfilePath(name)
	if err := applyChromeExport(exp, dstPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		os.Exit(constants.ExitChromeProfileCopyFailed)
	}
	fmt.Printf(constants.MsgChromeProfileImportOk, srcFile, name)
}

// runChromeProfileList implements `gitmap chrome-profile-list`.
func runChromeProfileList(args []string) {
	checkHelp(constants.CmdChromeProfileList, args)
	root := chromeUserDataDir()
	entries := chromeProfileEntries()
	if len(entries) == 0 {
		fmt.Printf(constants.MsgChromeProfileListEmpty, root)
		listChromeProfilesFromDB()
		return
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
}

// defaultChromeExportPath builds the default JSON output location
// under .gitmap/chrome/<name>.json (cwd-relative).
func defaultChromeExportPath(name string) string {
	return filepath.Join(constants.GitMapDir, "chrome", name+constants.ExtJSON)
}

// readChromeExport loads a JSON export file from disk.
func readChromeExport(path string) (*chromeExport, error) {
	raw, err := os.ReadFile(path) //nolint:gosec // user-supplied path
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", path, err)
	}
	var exp chromeExport
	if err := json.Unmarshal(raw, &exp); err != nil {
		return nil, fmt.Errorf("parse %s: %w", path, err)
	}
	return &exp, nil
}
