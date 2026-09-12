// Package cmd — clihelpers.go: shared command line and terminal helpers.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcg"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzip"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
	"github.com/spf13/cobra"
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

func runFixRepo(args []string) error {
	return cmdfixrepo.RunFixRepo(args)
}

type fixRepoBackupManifest = cmdfixrepo.FixRepoBackupManifest

type fixRepoIdentity struct {
	root    string
	host    string
	owner   string
	base    string
	current int
}

func resolveFixRepoIdentity() fixRepoIdentity {
	id := cmdfixrepo.ResolveFixRepoIdentity()
	return fixRepoIdentity{
		root:    id.Root,
		host:    id.Host,
		owner:   id.Owner,
		base:    id.Base,
		current: id.Current,
	}
}

func copyFileForBackup(src, dst string) error {
	return cmdfixrepo.CopyFileForBackup(src, dst)
}

const gofmtArgvOverhead = cmdfixrepo.GofmtArgvOverhead

func chunkPathsForGofmt(paths []string, maxCmdLen int) [][]string {
	return cmdfixrepo.ChunkPathsForGofmt(paths, maxCmdLen)
}

func batchCmdLen(batch []string) int {
	return cmdfixrepo.BatchCmdLen(batch)
}

var rewriteFixRepoFile = cmdfixrepo.RewriteFixRepoFile

var CountUnguardedTokenHits = cmdfixrepo.CountUnguardedTokenHits

var ScanUnguardedTokenHits = cmdfixrepo.ScanUnguardedTokenHits

func runDB(args []string) error {
	return cmddb.RunDB(args)
}

func runStartFresh(args []string) error {
	return cmddb.RunStartFresh(args)
}

var formatBytes = cmddb.FormatBytes
var confirmOrSkip = cmddb.ConfirmOrSkip
var isInteractiveStdin = cmddb.IsInteractiveStdin
var hasConfirmFlag = cmddb.HasConfirmFlag
var parseConfirmFlag = cmddb.ParseConfirmFlag
var truncateStr = cmddb.TruncateStr

func runPipeline(args []string) error {
	return cmdpipeline.RunPipeline(args)
}

func runPipelineAI(args []string) error {
	return cmdpipeline.RunPipelineAI(args)
}

var HandlePipelineLastFailedLogs = cmdpipeline.HandlePipelineLastFailedLogs
var handlePipelineDB = cmdpipeline.HandlePipelineDB
var IsNegativeIndexToken = cmdpipeline.IsNegativeIndexToken
var resolveTempDir = cmdpipeline.ResolveTempDir

type PipelineStatusPayload = cmdpipeline.PipelineStatusPayload
type PipelineErrorLogsPayload = cmdpipeline.PipelineErrorLogsPayload
type ErrorLogOutputParams = cmdpipeline.ErrorLogOutputParams
type SectionFailure = cmdpipeline.SectionFailure
type FailedRunItem = cmdpipeline.FailedRunItem
type FailedJobItem = cmdpipeline.FailedJobItem
type CICDCheckResult = cmdpipeline.CICDCheckResult
type ghRunItem = cmdpipeline.GhRunItem

var buildErrorLogsPayload = cmdpipeline.BuildErrorLogsPayload
var recordRunInSplitDb = cmdpipeline.RecordRunInSplitDb
var saveParsedFailedJobs = cmdpipeline.SaveParsedFailedJobs

func hasArgFlag(args []string, flagName string) bool {
	for _, a := range args {
		if a == flagName || strings.HasPrefix(a, flagName+"=") {
			return true
		}
	}

	return false
}

func extractFlagVal(args []string, flagName string) string {
	for i, arg := range args {
		if arg == flagName && i+1 < len(args) {
			return args[i+1]
		}

		if strings.HasPrefix(arg, flagName+"=") {
			return strings.TrimPrefix(arg, flagName+"=")
		}
	}

	return ""
}

func printJSON(v any) error {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}

	fmt.Println(string(b))

	return nil
}

// RunFixGit delegates to cmdfixgit.RunFixGit.
func runFixGit(args []string) error {
	return cmdfixgit.RunFixGit(args)
}

// CGOptions is an alias to cmdcg.CGOptions.
type CGOptions = cmdcg.CGOptions

// runCG delegates to cmdcg.RunCG.
func runCG(args []string) error {
	return cmdcg.RunCG(args)
}

// parseCGFlags delegates to cmdcg.ParseCGFlags.
func parseCGFlags(args []string) CGOptions {
	return cmdcg.ParseCGFlags(args)
}

// SSHJoinCmd is re-exported from cmdssh.
var SSHJoinCmd = cmdssh.SSHJoinCmd

// SJAddAuthCmd is re-exported from cmdssh.
var SJAddAuthCmd = cmdssh.SJAddAuthCmd

// SJHistCmd is re-exported from cmdssh.
var SJHistCmd = cmdssh.SJHistCmd

// SJLsCmd is re-exported from cmdssh.
var SJLsCmd = cmdssh.SJLsCmd

// runSSH delegates to cmdssh.RunSSH.
func runSSH(args []string) error {
	return cmdssh.RunSSH(args)
}

// runSSHExec delegates to cmdssh.RunSSHExec.
func runSSHExec(args []string) error {
	return cmdssh.RunSSHExec(args)
}

// runSSHBind delegates to cmdssh.RunSSHBind.
func runSSHBind(args []string) error {
	return cmdssh.RunSSHBind(args)
}

// runSJAddAuth delegates to cmdssh.RunSJAddAuth.
func runSJAddAuth(cmd *cobra.Command, args []string, ctx context.Context) error {
	return cmdssh.RunSJAddAuth(cmd, args, ctx)
}

// runSJHistory delegates to cmdssh.RunSJHistory.
func runSJHistory(cmd *cobra.Command, args []string, ctx context.Context) error {
	return cmdssh.RunSJHistory(cmd, args, ctx)
}

// runSJLs delegates to cmdssh.RunSJLs.
func runSJLs(cmd *cobra.Command, args []string, ctx context.Context) error {
	return cmdssh.RunSJLs(cmd, args, ctx)
}

// runSJRm delegates to cmdssh.RunSJRm.
func runSJRm(cmd *cobra.Command, args []string, ctx context.Context) error {
	return cmdssh.RunSJRm(cmd, args, ctx)
}

// runSSHJoin delegates to cmdssh.RunSSHJoin.
func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	return cmdssh.RunSSHJoin(cmd, args, ctx)
}

// parseSEFlags delegates to cmdssh.ParseSEFlags.
type SEOptions = cmdssh.SEOptions

func parseSEFlags(args []string) SEOptions {
	return cmdssh.ParseSEFlags(args)
}

// parseSJFlags delegates to cmdssh.ParseSJFlags.
type SJOptions = cmdssh.SJOptions

func parseSJFlags(args []string) SJOptions {
	return cmdssh.ParseSJFlags(args)
}

// ensureSSHDir delegates to cmdssh.EnsureSSHDir.
func ensureSSHDir(dir string) error {
	return cmdssh.EnsureSSHDir(dir)
}

// copyPubKeyAndAnnounce delegates to cmdssh.CopyPubKeyAndAnnounce.
func copyPubKeyAndAnnounce(pub string) {
	cmdssh.CopyPubKeyAndAnnounce(pub)
}

// resolveGitEmail delegates to cmdssh.ResolveGitEmail.
func resolveGitEmail() string {
	return cmdssh.ResolveGitEmail()
}

// validateSSHKeygen delegates to cmdssh.ValidateSSHKeygen.
func validateSSHKeygen() error {
	return cmdssh.ValidateSSHKeygen()
}

// encodeSSHListJSON delegates to cmdssh.EncodeSSHListJSON.
func encodeSSHListJSON(w io.Writer, keys []model.SSHKey) error {
	return cmdssh.EncodeSSHListJSON(w, keys)
}

// FixGitOptions aliases cmdfixgit.FixGitOptions.
type FixGitOptions = cmdfixgit.FixGitOptions

// FixGitIssue aliases cmdfixgit.FixGitIssue.
type FixGitIssue = cmdfixgit.FixGitIssue

// RemediateGitIndex delegates to cmdfixgit.RemediateGitIndex.
func RemediateGitIndex(repoRoot, gitDir string, opts FixGitOptions) ([]FixGitIssue, error) {
	return cmdfixgit.RemediateGitIndex(repoRoot, gitDir, opts)
}

// CGMetadata aliases cmdcg.CGMetadata.
type CGMetadata = cmdcg.CGMetadata

// WriteCGMetadata delegates to cmdcg.WriteCGMetadata.
func WriteCGMetadata(repoPath string, meta CGMetadata) error {
	return cmdcg.WriteCGMetadata(repoPath, meta)
}

// ReadCGMetadata delegates to cmdcg.ReadCGMetadata.
func ReadCGMetadata(repoPath string) (*CGMetadata, error) {
	return cmdcg.ReadCGMetadata(repoPath)
}

// ResolveCGTarget delegates to cmdcg.ResolveCGTarget.
func ResolveCGTarget(target string) (string, bool) {
	return cmdcg.ResolveCGTarget(target)
}

// SSHExecutor is forwarded from cmdssh.
var SSHExecutor = cmdssh.SSHExecutor

// RunSSHLogin delegates to cmdssh.RunSSHLogin.
func RunSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error {
	cmdssh.SSHExecutor = SSHExecutor
	return cmdssh.RunSSHLogin(cmd, args, ctx)
}

// ParseMultiIPList delegates to cmdssh.ParseMultiIPList.
func ParseMultiIPList(raw string) []string {
	return cmdssh.ParseMultiIPList(raw)
}

// VersionInstallConfig aliases cmdcg.VersionInstallConfig.
type VersionInstallConfig = cmdcg.VersionInstallConfig

// DefaultVersionInstallConfig delegates to cmdcg.DefaultVersionInstallConfig.
func DefaultVersionInstallConfig(initialVersion string) VersionInstallConfig {
	return cmdcg.DefaultVersionInstallConfig(initialVersion)
}

// InstallVersionJSON delegates to cmdcg.InstallVersionJSON.
func InstallVersionJSON(repoPath string, cfg VersionInstallConfig, isDryRun bool) error {
	return cmdcg.InstallVersionJSON(repoPath, cfg, isDryRun)
}

// ExportInstallerFlags aliases cmdinstaller.ExportInstallerFlags.
type ExportInstallerFlags = cmdinstaller.ExportInstallerFlags

// InstallerTreeNode aliases cmdinstaller.InstallerTreeNode.
type InstallerTreeNode = cmdinstaller.InstallerTreeNode

// RunInstallerCLI delegates to cmdinstaller.RunInstallerCLI.
func RunInstallerCLI(args []string) error {
	return cmdinstaller.RunInstallerCLI(args)
}

// runChrome delegates to cmdchrome.RunChrome.
func runChrome(args []string) error {
	return cmdchrome.RunChrome(args)
}

// findChromeBinaryPath delegates to cmdchrome.FindChromeBinaryPath.
func findChromeBinaryPath() (string, error) {
	return cmdchrome.FindChromeBinaryPath()
}

// isHelpFlag returns true if the argument matches standard help flags.
func isHelpFlag(arg string) bool {
	return arg == "help" || arg == "--help" || arg == "-h"
}

// readChromeBackup delegates to cmdchrome.ReadChromeBackup.
func readChromeBackup(src, dstRoot string) (int, error) {
	return cmdchrome.ReadChromeBackup(src, dstRoot)
}

// runSetup delegates to cmdsetup.RunSetup.
func runSetup(args []string) error {
	return cmdsetup.RunSetup(args)
}

// runSetupPerms delegates to cmdsetup.RunSetupPerms.
func runSetupPerms(args []string) error {
	return cmdsetup.RunSetupPerms(args)
}

// warnIfNoWrapper delegates to cmdsetup.WarnIfNoWrapper.
func warnIfNoWrapper() {
	cmdsetup.WarnIfNoWrapper()
}

// resolveSetupConfigPath delegates to cmdsetup.ResolveSetupConfigPath.
//
//nolint:unused
func resolveSetupConfigPath(configPath string, hasConfig bool) string {
	return cmdsetup.ResolveSetupConfigPath(configPath, hasConfig)
}

// isWrapperActive delegates to cmdsetup.IsWrapperActive.
//
//nolint:unused
func isWrapperActive() bool {
	return cmdsetup.IsWrapperActive()
}

type installOptions = cmdinstall.InstallOptions
type ProfileComposition = cmdinstall.ProfileComposition

func runInstall(args []string) error {
	return cmdinstall.RunInstall(args)
}

func runInstalledDir() error {
	return cmdinstall.RunInstalledDir()
}

func installTool(opts installOptions) {
	cmdinstall.InstallTool(opts)
}

func runInstallAdd(args []string) error {
	return cmdinstall.RunInstallAdd(args)
}

func runInstallChromeLinux(opts installOptions) error {
	return cmdinstall.RunInstallChromeLinux(opts)
}

func resolveProfileTree(slug string) (ProfileComposition, bool) {
	return cmdinstall.ResolveProfileTree(slug)
}

func printProfileTree(prof ProfileComposition) {
	cmdinstall.PrintProfileTree(prof)
}

func printProfileInstallSummary(slug string) {
	cmdinstall.PrintProfileInstallSummary(slug)
}

func resolvePowerShellBinaryWithLookPath(lookPath func(string) (string, error)) string {
	return cmdinstall.ResolvePowerShellBinaryWithLookPath(lookPath)
}

func getInstalledVersion(binary string) string {
	return cmdinstall.GetInstalledVersion(binary)
}

func validateToolName(tool string) {
	cmdinstall.ValidateToolName(tool)
}

func runUninstallCtx() {
	cmdinstall.RunUninstallCtx()
}

func resolvePackageManager(override, tool string) string {
	return cmdinstall.ResolvePackageManager(override, tool)
}

func resolvePackageName(tool, manager string) string {
	return cmdinstall.ResolvePackageName(tool, manager)
}

func runInstallCommand(args []string, opts installOptions) error {
	return cmdinstall.RunInstallCommand(args, opts)
}

func init() {
	cmdssh.JoinRunner = runJoin
	cmdssh.ProfileRunner = runProfile

	cmdinstaller.ResolveProfileTreeFn = func(s string) (any, bool) {
		p, ok := resolveProfileTree(s)
		return p, ok
	}
	cmdinstaller.PrintProfileTreeFn = func(p any) {
		if prof, ok := p.(ProfileComposition); ok {
			printProfileTree(prof)
		}
	}
	cmdinstaller.PrintProfileInstallSummaryFn = printProfileInstallSummary
	cmdinstaller.RunInstallAddFn = runInstallAdd

	cmdchrome.InstallChromeLinuxFn = func(isDryRun bool) error {
		return runInstallChromeLinux(installOptions{Tool: constants.ToolChrome, DryRun: isDryRun})
	}
	cmdchrome.InstallToolFn = func(tool string, isDryRun bool) {
		installTool(installOptions{Tool: tool, DryRun: isDryRun})
	}
	cmdchrome.RunFindDuplicatesFn = func(category string, args []string) error {
		return runFindDuplicates(category, args)
	}
	cmdchrome.CheckHelpFn = checkHelp
	cmdinstall.CheckHelpFn = checkHelp
}
