// Package cmd — clihelpers.go: shared command line and terminal helpers.
package cmd

import (
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzip"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func isFileFlagWithArg(arg string) bool {
	return arg == "--file" || arg == "--out" || arg == "-o" || arg == "-f"
}

func matchFlagWithVal(arg string, names ...string) bool {
	for _, n := range names {
		if arg == n || strings.HasPrefix(arg, n+"=") {
			return true
		}
	}

	return false
}

func extractFlagValue(idx *int, args []string) string {
	a := args[*idx]
	if strings.Contains(a, "=") {
		parts := strings.SplitN(a, "=", 2)

		return parts[1]
	}

	if *idx+1 < len(args) {
		*idx++

		return args[*idx]
	}

	return ""
}

func runCatCmd(args []string) error {
	return cmdmacro.RunCatCmd(args)
}

func runTouchCmd(args []string) error {
	return cmdmacro.RunTouchCmd(args)
}

func runMkfileCmd(args []string) error {
	return cmdmacro.RunMkfileCmd(args)
}

func parseExecOptions(args []string) macro.ExecOptions {
	return cmdmacro.ParseExecOptions(args)
}

func outputStructuredData(data interface{}, opts macro.ExecOptions) error {
	return cmdmacro.OutputStructuredData(data, opts)
}

func extractMacroNameAndFlags(args []string) (string, []string) {
	return cmdmacro.ExtractMacroNameAndFlags(args)
}

func parseDurationArg(val string, fallback time.Duration) time.Duration {
	return cmdmacro.ParseDurationArg(val, fallback)
}

func runVSCode(args []string) error {
	return cmdvscode.RunVSCode(args)
}

func runVSCodePMSync(args []string) error {
	return cmdvscode.RunVSCodePMSync(args)
}

func runVSCodePMPath(args []string) error {
	return cmdvscode.RunVSCodePMPath(args)
}

func runVSCodeWorkspace(args []string) error {
	return cmdvscode.RunVSCodeWorkspace(args)
}

func runFindDuplicatesVSCode() error {
	return cmdvscode.RunFindDuplicates()
}

func syncRecordsToVSCodePM(records []model.ScanRecord, noVSCodeSync, noAutoTags bool) {
	cmdvscode.SyncRecordsToVSCodePM(records, noVSCodeSync, noAutoTags)
}

func syncClonedReposToVSCodePM(pairs []vscodepm.Pair, skip bool) {
	cmdvscode.SyncClonedReposToVSCodePM(pairs, skip)
}

func syncSingleClonedRepoToVSCodePM(absPath, repoName string, skip bool) {
	cmdvscode.SyncSingleClonedRepoToVSCodePM(absPath, repoName, skip)
}

func buildClonePMPair(absPath, repoName string) vscodepm.Pair {
	return cmdvscode.BuildClonePMPair(absPath, repoName)
}

func canonicalizePMPath(absPath string) string {
	return cmdvscode.CanonicalizePMPath(absPath)
}

func isVSCodeSyncDisabled() bool {
	return cmdvscode.IsVSCodeSyncDisabled()
}

func reportVSCodePMSoftError(err error) {
	cmdvscode.ReportVSCodePMSoftError(err)
}

func renameVSCodePMByPath(absPath, newName string) {
	cmdvscode.RenameVSCodePMByPath(absPath, newName)
}

func printVSCodeOptimizeResult(s vscodepm.OptimizeSummary, isDryRun bool) {
	cmdvscode.PrintVSCodeOptimizeResult(s, isDryRun)
}

func runGitHubDesktopGroup(args []string) error {
	return cmdvscode.RunGitHubDesktopGroup(args)
}

func stripVSCodeSyncDisabledFlag(args []string) []string {
	return cmdvscode.StripVSCodeSyncDisabledFlag(args)
}

func stripVSCodeTagFlags(args []string) []string {
	return cmdvscode.StripVSCodeTagFlags(args)
}

func applyDebugPathsEnv(isOn bool) {
	cmdvscode.ApplyDebugPathsEnv(isOn)
}

func runVHost(args []string) error {
	return cmdvhost.RunVHost(args)
}

func runNginxRm(args []string) error {
	return cmdvhost.RunVHostRm(args)
}

func runVHostEnable(args []string) error {
	return cmdvhost.RunVHostEnable(args)
}

func runVHostDisable(args []string) error {
	return cmdvhost.RunVHostDisable(args)
}

func runVHostCreate(args []string) error {
	return cmdvhost.RunVHostCreate(args)
}

func runVHostTest(args []string) error {
	return cmdvhost.RunVHostTest(args)
}

func runVHostReload(args []string) error {
	return cmdvhost.RunVHostReload(args)
}

func runVHostList(args []string) error {
	return cmdvhost.RunVHostList(args)
}

type VHostConfig = cmdvhost.VHostConfig

const (
	VHostSiteTypeLaravel   = cmdvhost.VHostSiteTypeLaravel
	VHostSiteTypeWordpress = cmdvhost.VHostSiteTypeWordpress
)

func RenderVHostConfig(cfg VHostConfig) (string, *apperror.AppError) {
	return cmdvhost.RenderVHostConfig(cfg)
}

func runZip(args []string) error {
	return cmdzip.RunZip(args)
}

func runUnzipCompact(args []string) error {
	return cmdzip.RunUnzipCompact(args)
}

func runZipGroup(args []string) error {
	return cmdzip.RunZipGroup(args)
}
