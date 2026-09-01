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
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
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
		cliexit.HandleError(nil, constants.ExitChromeProfileUsage)
	}
	srcProfile, ok := resolveChromeProfile(pos[0])
	dstProfile := chromeProfileDestination(pos[1])
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, pos[0], srcProfile.Path)
		printAvailableChromeProfilesWithDisplay()
		cliexit.HandleError(nil, constants.ExitChromeProfileNotFound)
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
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
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
	isNonRunning := !isRunning
	if isNonRunning {
		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrChromeProfileChromeOpen, src, dst)
	cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
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

func parseChromeExportArgs(args []string) (string, []string) {
	format := ""
	var positional []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "--format=") {
			format = strings.TrimPrefix(arg, "--format=")
		} else if !strings.HasPrefix(arg, "-") {
			positional = append(positional, arg)
		}
	}
	return format, positional
}

// runChromeProfileExport implements `gitmap chrome-profile-export` / `gitmap cpe`.
func runChromeProfileExport(args []string) error {
	checkHelp(constants.CmdChromeProfileExport, args)
	format, positional := parseChromeExportArgs(args)
	if len(positional) == 0 {
		return runExportAllProfilesToPath("", format)
	}
	target := positional[0]
	if isAllProfilesTarget(target, positional) {
		outPath := resolveAllProfilesOutPath(positional)
		return runExportAllProfilesToPath(outPath, format)
	}
	return runExportSingleProfileNamed(target, positional, format)
}

func isAllProfilesTarget(target string, positional []string) bool {
	if target == "all" || target == "--all" {
		return true
	}
	if _, isProfile := resolveChromeProfileDir(target); isProfile {
		return false
	}
	return true
}

func resolveAllProfilesOutPath(positional []string) string {
	if len(positional) >= 2 && (positional[0] == "all" || positional[0] == "--all") {
		return positional[1]
	}
	if len(positional) >= 1 && positional[0] != "all" && positional[0] != "--all" {
		return positional[0]
	}
	return filepath.Join(constants.GitMapDir, "chrome")
}

func runExportSingleProfileNamed(name string, positional []string, format string) error {
	srcPath, isProfile := resolveChromeProfileDir(name)
	if !isProfile {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileSrcMissing, name, srcPath)
		printAvailableChromeProfilesWithDisplay()
		cliexit.HandleError(nil, constants.ExitChromeProfileNotFound)
		return nil
	}
	outPath := defaultChromeExportPath(name, format)
	if len(positional) >= 2 {
		outPath = positional[1]
	}
	format = inferExportFormatFromPath(outPath, format)
	return executeSingleProfileExport(format, srcPath, name, outPath)
}

func executeSingleProfileExport(format, srcPath, name, outPath string) error {
	rec, exportErr := exportChromeFormat(format, srcPath, name, outPath)
	if exportErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, exportErr)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
		return exportErr
	}
	printChromeArtifacts(rec)
	persistChromeProfile(name, srcPath, rec)
	return nil
}

func exportChromeFormat(format, srcPath, name, outPath string) (chromeExportRecord, error) {
	if format == constants.OutputZIP {
		bytes, err := writeChromeExportZIP(srcPath, name, outPath)
		return chromeExportRecord{ZIPPath: outPath, ZIPSize: bytes}, err
	}
	if format == constants.OutputSQLite {
		bytes, err := writeAllChromeProfilesSQLite([]string{name}, outPath)
		return chromeExportRecord{SQLitePath: outPath, SQLiteSize: bytes}, err
	}
	if format == constants.OutputYAML {
		bytes, err := writeAllChromeProfilesYAML([]string{name}, outPath)
		return chromeExportRecord{JSONPath: outPath, JSONSize: bytes}, err
	}
	return exportChromeJSONAndCSV(srcPath, name, outPath)
}

func exportChromeJSONAndCSV(srcPath, name, outPath string) (chromeExportRecord, error) {
	jsonBytes, err := writeChromeExport(srcPath, name, outPath)
	if err != nil {
		return chromeExportRecord{}, err
	}

	csvPath := resolveChromeCSVPath(outPath)
	csvBytes, csvErr := writeChromeExportCSV(srcPath, name, csvPath)
	if csvErr != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileExportFail, csvErr)
		csvPath = ""
	}

	return chromeExportRecord{
		JSONPath: outPath, JSONSize: jsonBytes,
		CSVPath: csvPath, CSVSize: csvBytes,
	}, nil
}

func resolveChromeCSVPath(outPath string) string {
	ext := constants.ExtJSON
	if len(outPath) > len(ext) && outPath[len(outPath)-len(ext):] == ext {
		return outPath[:len(outPath)-len(ext)] + constants.ExtCSV
	}
	return outPath + constants.ExtCSV
}

// runChromeProfileImport implements `gitmap chrome-profile-import` / `gitmap cpi`.
func runChromeProfileImport(args []string) error {
	checkHelp(constants.CmdChromeProfileImport, args)
	if len(args) < 1 {
		fmt.Fprint(os.Stderr, constants.ErrChromeProfileUsageImport)
		cliexit.HandleError(nil, constants.ExitChromeProfileUsage)
	}
	srcFile := args[0]
	name := ""
	if len(args) >= 2 {
		name = args[1]
	}
	if err := importChromeSnapshot(srcFile, name); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
		return err
	}
	return nil
}

func importChromeArchive(srcFile, name string) error {
	if name == "" {
		base := filepath.Base(srcFile)
		name = strings.TrimSuffix(base, filepath.Ext(base))
	}
	dstPath := chromeProfilePath(name)
	if err := applyChromeExportZIP(srcFile, dstPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
	}
	fmt.Printf(constants.MsgChromeProfileImportOk, srcFile, name)
	return nil
}

func applyImportedProfile(exp *chromeExport, srcFile, name string) error {
	if name == "" {
		name = exp.Name
	}
	dstPath := chromeProfilePath(name)
	if err := applyChromeExport(exp, dstPath); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrChromeProfileImportFail, err)
		cliexit.HandleError(nil, constants.ExitChromeProfileCopyFailed)
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
		ext = constants.ExtZIP
	} else if format == constants.OutputSQLite {
		ext = constants.ExtDB
	} else if format == constants.OutputYAML {
		ext = constants.ExtYAML
	}
	return filepath.Join(constants.GitMapDir, "chrome", name+ext)
}
