package cmdscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/formatter"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// writeAllOutputs writes terminal, CSV, JSON, text, folder structure, and clone fspath.
func writeAllOutputs(records []model.ScanRecord, outputDir, outFile string, quiet, compact bool) {
	records = injectCompanionRecords(records)
	writeTerminalOutput(records, outputDir, quiet)
	writeCSVOutput(records, outputDir, outFile)
	writeJSONOutput(records, outputDir)
	writeCWDJSONOutput(records, compact)
	writeTextOutput(records, outputDir)
	writeFolderStructure(records, outputDir)
	writeCloneScript(records, outputDir)
	writeDirectCloneScript(records, outputDir)
	writeDirectCloneSSHScript(records, outputDir)
	writeDesktopScript(records, outputDir)
}

// writeCWDJSONOutput writes a portable JSON file to CWD.
func writeCWDJSONOutput(records []model.ScanRecord, compact bool) {
	records = injectCompanionRecords(records)
	cwd, err := os.Getwd()
	if err != nil {
		return
	}

	path := filepath.Join(cwd, constants.DefaultJSONFile)
	writeFormattedFile(path, func(f *os.File) error {
		if compact {
			return formatter.WriteJSONCompact(f, records)
		}
		return formatter.WriteJSON(f, records)
	}, "CWD JSON")
}

// writeTerminalOutput renders records to stdout.
func writeTerminalOutput(records []model.ScanRecord, outputDir string, quiet bool) {
	err := formatter.Terminal(os.Stdout, records, outputDir, quiet)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrOutputFailed, err)
	}
}

func writeFormattedFile(path string, writeFn func(*os.File) error, desc string) {
	file, err := createOutputFile(path)
	if err != nil {
		return
	}
	defer file.Close()
	if err := writeFn(file); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not write %s to %s: %v\n", desc, path, err)
	}
}

// writeCSVOutput writes records to a CSV file.
func writeCSVOutput(records []model.ScanRecord, outputDir, outFile string) {
	path := resolveOutFile(outFile, outputDir, constants.DefaultCSVFile)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteCSV(f, records) }, "CSV")
}

// writeJSONOutput writes records to a JSON file.
func writeJSONOutput(records []model.ScanRecord, outputDir string) {
	records = injectCompanionRecords(records)
	path := filepath.Join(outputDir, constants.DefaultJSONFile)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteJSON(f, records) }, "JSON")
}

// writeTextOutput writes records as plain text clone commands.
func writeTextOutput(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultTextFile)
	writeFormattedFile(path, func(f *os.File) error {
		formatter.WriteText(f, records)
		return nil
	}, "text")
}

// writeFolderStructure writes a Markdown file showing the repo tree.
func writeFolderStructure(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultStructureFile)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteStructure(f, records) }, "structure")
}

// writeCloneScript writes a PowerShell clone script.
func writeCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultCloneScript)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteCloneScript(f, records) }, "clone script")
}

// writeDirectCloneScript writes a plain PS1 with one git clone per line.
func writeDirectCloneScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneScript)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteDirectCloneScript(f, records) }, "direct clone script")
}

// writeDirectCloneSSHScript writes a plain SSH PS1 with one git clone per line.
func writeDirectCloneSSHScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDirectCloneSSHScript)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteDirectCloneSSHScript(f, records) }, "SSH clone script")
}

// writeDesktopScript writes a PowerShell script to register repos with GitHub Desktop.
func writeDesktopScript(records []model.ScanRecord, outputDir string) {
	path := filepath.Join(outputDir, constants.DefaultDesktopScript)
	writeFormattedFile(path, func(f *os.File) error { return formatter.WriteDesktopScript(f, records) }, "desktop script")
}

// resolveOutFile determines the output file path.
func resolveOutFile(outFile, outputDir, defaultName string) string {
	if len(outFile) > 0 {
		return outFile
	}

	return filepath.Join(outputDir, defaultName)
}

// createOutputFile ensures the directory exists and creates the file.
func createOutputFile(path string) (*os.File, error) {
	err := os.MkdirAll(filepath.Dir(path), constants.DirPermission)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateDir, filepath.Dir(path), err)

		return nil, err
	}

	file, err := os.Create(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCreateFile, path, err)

		return nil, err
	}

	return file, nil
}

func injectCompanionRecords(records []model.ScanRecord) []model.ScanRecord {
	probe, hasRepoCache := ProbeWorkspaceRepoCacheRecord()
	if !hasRepoCache || probe == nil {
		return records
	}
	for i := range records {
		if strings.EqualFold(records[i].RepoName, "repo-cache") {
			mergeRepoCacheRecord(&records[i], probe)
			return records
		}
	}
	return append(records, *probe)
}

func mergeRepoCacheRecord(existing, probe *model.ScanRecord) {
	if len(existing.DiscoveredURL) == 0 {
		existing.DiscoveredURL = probe.DiscoveredURL
		existing.HTTPSUrl = probe.HTTPSUrl
	}
	if len(existing.Notes) == 0 {
		existing.Notes = probe.Notes
	}
	if len(existing.PrimaryAlias) == 0 {
		existing.PrimaryAlias = probe.PrimaryAlias
	}
}

// ProbeWorkspaceRepoCacheRecord detects whether repo-cache exists and returns a synthesized ScanRecord.
func ProbeWorkspaceRepoCacheRecord() (*model.ScanRecord, bool) {
	root, remoteURL, hasRepoCache := probeRepoCacheDir()
	if !hasRepoCache || !isDir(filepath.Join(root, ".git")) {
		return nil, false
	}
	rec := synthesizeRepoCacheRecord(root, remoteURL)
	return &rec, true
}

func probeRepoCacheDir() (string, string, bool) {
	if ResolveRepoCacheRootFn != nil {
		hookRoot := ResolveRepoCacheRootFn()
		if isDir(hookRoot) {
			return hookRoot, "", true
		}
	}
	if root, remote, hasSpecial := probeSpecialReposDB(); hasSpecial {
		return root, remote, true
	}
	if root, hasDefault := probeDefaultWorkDir(); hasDefault {
		return root, "", true
	}
	if isDir(`D:\work\repo-cache`) {
		return `D:\work\repo-cache`, "", true
	}
	return "", "", false
}

func probeSpecialReposDB() (string, string, bool) {
	db, err := store.OpenSpecialReposSplitDB()
	if err != nil {
		return "", "", false
	}
	defer db.Close()
	rec, err := db.GetSpecialRepo("repo-cache")
	hasRec := err == nil && rec != nil && len(rec.LocalPath) > 0 && isDir(rec.LocalPath)
	if hasRec {
		return rec.LocalPath, rec.RemoteURL, true
	}
	return "", "", false
}

func probeDefaultWorkDir() (string, bool) {
	db, err := store.OpenDefault()
	if err != nil {
		return "", false
	}
	defer db.Close()
	wd, err := db.GetDefaultWorkDir()
	hasDefault := err == nil && wd != nil && len(wd.AbsolutePath) > 0
	if hasDefault {
		candidate := filepath.Join(wd.AbsolutePath, "repo-cache")
		return candidate, isDir(candidate)
	}
	return "", false
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func synthesizeRepoCacheRecord(root, remoteURL string) model.ScanRecord {
	if len(remoteURL) == 0 {
		remoteURL = readOriginRemoteURL(root)
	}
	relPath := resolveRepoCacheRelativePath(root)
	cloneCmd := "gitmap clone rc"
	if len(remoteURL) > 0 {
		cloneCmd = fmt.Sprintf("git clone %s %s", remoteURL, relPath)
	}
	return model.ScanRecord{
		RepoName: "repo-cache", Slug: "repo-cache", RepoID: "repo-cache",
		DiscoveredURL: remoteURL, HTTPSUrl: remoteURL, RelativePath: relPath,
		AbsolutePath: root, Branch: "main", CloneInstruction: cloneCmd,
		Notes: "companion:repo-cache", PrimaryAlias: "rc",
		Aliases: []string{"companion:repo-cache", "special:rc"},
	}
}

func resolveRepoCacheRelativePath(root string) string {
	cwd, err := os.Getwd()
	if err != nil {
		return "repo-cache"
	}
	rel, err := filepath.Rel(cwd, root)
	if err != nil || len(rel) == 0 {
		return "repo-cache"
	}
	return filepath.ToSlash(rel)
}

func readOriginRemoteURL(repoDir string) string {
	content, err := os.ReadFile(filepath.Join(repoDir, ".git", "config"))
	if err == nil {
		for _, line := range strings.Split(string(content), "\n") {
			trimmed := strings.TrimSpace(line)
			if strings.HasPrefix(trimmed, "url =") {
				return strings.TrimSpace(strings.TrimPrefix(trimmed, "url ="))
			}
		}
	}
	return ""
}
