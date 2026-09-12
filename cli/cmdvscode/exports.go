package cmdvscode

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

// VSCodeCmd runs the main vscode command.
func VSCodeCmd(args []string) error {
	return runVSCode(args)
}

// RunVSCode runs the vscode CLI entrypoint.
func RunVSCode(args []string) error {
	return runVSCode(args)
}

// RunVSCodePMSync runs the vscode-pm-sync CLI entrypoint.
func RunVSCodePMSync(args []string) error {
	return runVSCodePMSync(args)
}

// RunVSCodePMPath runs the vscode-pm-path CLI entrypoint.
func RunVSCodePMPath(args []string) error {
	return runVSCodePMPath(args)
}

// RunVSCodeWorkspace runs the vscode-workspace CLI entrypoint.
func RunVSCodeWorkspace(args []string) error {
	return runVSCodeWorkspace(args)
}

// RunFindDuplicates discovers duplicate projects across VS Code configurations.
func RunFindDuplicates() error {
	return runFindDuplicatesVSCode()
}

// IsVSCodeSyncDisabled checks if automatic VS Code synchronization is disabled.
func IsVSCodeSyncDisabled() bool {
	return isVSCodeSyncDisabled()
}

// ReportVSCodePMSoftError logs soft errors during VS Code PM sync without failing CLI.
func ReportVSCodePMSoftError(err error) {
	reportVSCodePMSoftError(err)
}

// SyncClonedReposToVSCodePM syncs a slice of cloned repositories to VS Code PM.
func SyncClonedReposToVSCodePM(pairs []vscodepm.Pair, skip bool) {
	syncClonedReposToVSCodePM(pairs, skip)
}

// SyncSingleClonedRepoToVSCodePM syncs a single cloned repository to VS Code PM.
func SyncSingleClonedRepoToVSCodePM(absPath, repoName string, skip bool) {
	syncSingleClonedRepoToVSCodePM(absPath, repoName, skip)
}

// BuildClonePMPair constructs a vscodepm.Pair from repo attributes.
func BuildClonePMPair(absPath, repoName string) vscodepm.Pair {
	return buildClonePMPair(absPath, repoName)
}

// CanonicalizePMPath returns the canonical path for VS Code Project Manager.
func CanonicalizePMPath(absPath string) string {
	return canonicalizePMPath(absPath)
}

// SyncRecordsToVSCodePM syncs scan records to VS Code PM.
func SyncRecordsToVSCodePM(records []model.ScanRecord, noVSCodeSync, noAutoTags bool) {
	syncRecordsToVSCodePM(records, noVSCodeSync, noAutoTags)
}

// RenameVSCodePMByPath renames a project manager entry by its path.
func RenameVSCodePMByPath(absPath, newName string) {
	renameVSCodePMByPath(absPath, newName)
}

// PrintVSCodeOptimizeResult prints the result of VS Code optimization.
func PrintVSCodeOptimizeResult(s vscodepm.OptimizeSummary, isDryRun bool) {
	printVSCodeOptimizeResult(s, isDryRun)
}

// RunGitHubDesktopGroup runs GitHub Desktop grouping.
func RunGitHubDesktopGroup(args []string) error {
	return runGitHubDesktopGroup(args)
}

// StripVSCodeSyncDisabledFlag removes --no-vscode-sync from args.
func StripVSCodeSyncDisabledFlag(args []string) []string {
	return stripVSCodeSyncDisabledFlag(args)
}

// StripVSCodeTagFlags removes vscode tag flags from args.
func StripVSCodeTagFlags(args []string) []string {
	return stripVSCodeTagFlags(args)
}

// ApplyDebugPathsEnv sets or unsets GITMAP_DEBUG_PATHS.
func ApplyDebugPathsEnv(isOn bool) {
	applyDebugPathsEnv(isOn)
}
