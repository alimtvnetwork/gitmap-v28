// Package cmd — clihelpers.go: shared command line and terminal helpers.
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcg"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdchrome"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddb"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmddoctor"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixgit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixrepo"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdmacro"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpipeline"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdschedule"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsetup"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvhost"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdvscode"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdworkdir"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdzip"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/spf13/cobra"
)

func isTerminalInput() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}

	return (stat.Mode() & os.ModeCharDevice) != 0
}

func isExistingFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func isFileFlagWithArg(arg string) bool {
	return arg == "--file" || arg == "--out" || arg == "-o" || arg == "-f"
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

func syncRecordsToVSCodePM(records []model.ScanRecord, isSkipVSCodeSync, isSkipAutoTags bool) {
	cmdvscode.SyncRecordsToVSCodePM(records, isSkipVSCodeSync, isSkipAutoTags)
}

func reportVSCodePMSoftError(err error) {
	cmdvscode.ReportVSCodePMSoftError(err)
}

func renameVSCodePMByPath(absPath, newName string) {
	cmdvscode.RenameVSCodePMByPath(absPath, newName)
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

func runUpdate() error {
	return cmdupdate.RunUpdate()
}

func runUpdateCleanup() error {
	return cmdupdate.RunUpdateCleanup()
}

func runUpdateRunner() error {
	return cmdupdate.RunUpdateRunner()
}

func scheduleDeployedCleanupHandoff() {
	cmdupdate.ScheduleDeployedCleanupHandoff()
}

func initRunnerVerbose() {
	cmdupdate.InitRunnerVerbose()
}

func expandTilde(path string) string {
	return cmdupdate.ExpandTilde(path)
}

func createHandoffCopy(selfPath string) string {
	return cmdupdate.CreateHandoffCopy(selfPath)
}

func hasFlag(flagName string) bool {
	return cmdupdate.HasFlag(flagName)
}

func handleHandoffError(err error) {
	cmdupdate.HandleHandoffError(err)
}

func writeScriptToTemp(script string) (string, error) {
	return cmdupdate.WriteScriptToTemp(script)
}

func normalizeRepoPath(path string) string {
	return cmdupdate.NormalizeRepoPath(path)
}

func saveRepoPathToDB(path string) {
	cmdupdate.SaveRepoPathToDB(path)
}

func runPull(args []string) error {
	return cmdpull.RunPull(args)
}

func runPullAll(args []string) error {
	return cmdpull.RunPullAll(args)
}

func runPullReleaseCD(args []string) error {
	return cmdpull.RunPullReleaseCD(args)
}

func runPush(args []string) error {
	return cmdpull.RunPush(args)
}

func findBySlug(records []model.ScanRecord, slug string) []model.ScanRecord {
	return cmdpull.FindBySlug(records, slug)
}

func isGitRepoCWD() bool {
	return cmdpull.IsGitRepoCWD()
}

func pullOneRepo(rec model.ScanRecord) {
	cmdpull.PullOneRepo(rec)
}

func ResolvePullDirectoryTargets(dirPath string) []model.ScanRecord {
	return cmdpull.ResolvePullDirectoryTargets(dirPath)
}

func ExtractTransportFlags(args []string) (bool, bool, []string) {
	return cmdpull.ExtractTransportFlags(args)
}

func repoNameFromURL(rawURL string) string {
	return cmdclone.RepoNameFromURL(rawURL)
}

func extractRepoName(rawURL string) string {
	return cmdclone.ExtractRepoName(rawURL)
}

func resolveCloneNextFolder(token string) (string, error) {
	return cmdclone.ResolveCloneNextFolder(token)
}

func openInVSCode(absPath string) {
	cmdclone.OpenInVSCode(absPath)
}

func registerSingleDesktop(name, absPath string) {
	cmdclone.RegisterSingleDesktop(name, absPath)
}

func runClone(args []string) error {
	return cmdclone.RunClone(args)
}

func runCloneFixRepo(args []string) error {
	return cmdclone.RunCloneFixRepo(args)
}

func runCloneFixRepoPub(args []string) error {
	return cmdclone.RunCloneFixRepoPub(args)
}

func runCloneFrom(args []string) error {
	return cmdclone.RunCloneFrom(args)
}

func runCloneNext(args []string) error {
	return cmdclone.RunCloneNext(args)
}

func runCloneNow(args []string) error {
	return cmdclone.RunCloneNow(args)
}

func runClonePick(args []string) error {
	return cmdclone.RunClonePick(args)
}

func runCloneSync() error {
	return cmdclone.RunCloneSync()
}

func ConvertURLToSSH(rawURL string) (string, bool) {
	return cmdclone.ConvertURLToSSH(rawURL)
}

func ConvertURLToHTTPS(rawURL string) (string, bool) {
	return cmdclone.ConvertURLToHTTPS(rawURL)
}

func ResolveCloneFixRepoName(absPath string) string {
	return cmdclone.ResolveCloneFixRepoName(absPath)
}

func RunRepoReclone(target string, yes bool) error {
	return cmdclone.RunRepoReclone(target, yes)
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

	cmdupdate.RunPostUpdateMigrateFn = runPostUpdateMigrate
	cmdupdate.RequireOnlineFn = requireOnline
	cmdupdate.ResolveDeployedAndConfigPathsFn = resolveDeployedAndConfigPaths
	cmdupdate.PrintGitmapIdentityBlockLongFn = printGitmapIdentityBlockLong

	cmdpull.LoadAllRecordsDBFn = loadAllRecordsDB
	cmdpull.LoadRecordsByGroupFn = loadRecordsByGroup
	cmdpull.CreatePendingTaskFn = createPendingTask
	cmdpull.CompletePendingTaskFn = completePendingTask
	cmdpull.FailPendingTaskFn = failPendingTask
	cmdpull.RequireOnlineFn = requireOnline
	cmdpull.CheckHelpFn = checkHelp
	cmdpull.ApplyTransportFlagFn = ApplyTransportFlag
	cmdpull.HasAliasFn = HasAlias
	cmdpull.GetAliasSlugFn = GetAliasSlug
	cmdpull.GetAliasPathFn = GetAliasPath
	cmdpull.RunStatusFn = runStatus
	cmdpull.PrintRemediationSummaryNoPromptFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummaryNoPrompt(cmdItems)
	}
	cmdpull.PrintRemediationSummaryAutoFixFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummaryAutoFix(cmdItems)
	}
	cmdpull.PrintRemediationSummaryFn = func(items []cmdpull.RemediationItem) {
		cmdItems := make([]RemediationItem, len(items))
		for i, it := range items {
			cmdItems[i] = RemediationItem{
				RepoPath:      it.RepoPath,
				RepoName:      it.RepoName,
				SummaryReason: it.SummaryReason,
				Recipes:       it.Recipes,
				Files:         it.Files,
			}
		}
		PrintRemediationSummary(cmdItems)
	}

	cmdclone.CreatePendingTaskFn = createPendingTask
	cmdclone.CompletePendingTaskFn = completePendingTask
	cmdclone.FailPendingTaskFn = failPendingTask
	cmdclone.RequireOnlineFn = requireOnline
	cmdclone.CheckHelpFn = checkHelp
	cmdclone.WriteShellHandoffFn = WriteShellHandoff
	cmdclone.EscapeCwdIfInsideFn = escapeCwdIfInside
	cmdclone.FinalizeErrorReportFn = cmdscan.FinalizeErrorReport
	cmdclone.RunCodingGuidelinesInstallFn = func(dir string) error {
		return RunCodingGuidelinesInstall(CodingGuidelinesOpts{WorkingDir: dir})
	}
	cmdclone.CommitCodingGuidelinesFn = func(dir string, isSkipCommit, isSkipPush bool) error {
		return CommitCodingGuidelines(CGCommitOpts{WorkingDir: dir, IsSkipCommit: isSkipCommit, IsSkipPush: isSkipPush})
	}
	cmdclone.RunCFRPPriorVersionPrivatizeFn = runCFRPPriorVersionPrivatize
	cmdclone.RunGitHubDesktopOptimizeFn = runGitHubDesktopOptimize
	cmdclone.ResolveEndpointStringFn = resolveEndpointString
	cmdclone.ResolveReleaseAliasPathFn = resolveReleaseAliasPath
	cmdclone.RunStatusFn = runStatus

	cmdscan.CreatePendingTaskFn = createPendingTask
	cmdscan.CompletePendingTaskFn = completePendingTask
	cmdscan.FailPendingTaskFn = failPendingTask
	cmdscan.SyncRecordsToVSCodePMFn = syncRecordsToVSCodePM
	cmdscan.RunPruneStaleDBFn = runPruneStaleDB
	cmdscan.CheckHelpFn = checkHelp

	cmdos.RunPowerNeverSleepFn = runPowerNeverSleep
	cmdos.RunPowerSetFn = runPowerSet
	cmdos.RunPowerResetFn = runPowerReset
	cmdos.CheckHelpFn = checkHelp

	cmdschedule.CheckHelpFn = checkHelp

	cmdworkdir.CheckHelpFn = checkHelp
}

func runScan(args []string) error {
	return cmdscan.RunScan(args)
}

func runRescan() error {
	return cmdscan.RunRescan()
}

func runRescanSubtree(args []string) error {
	return cmdscan.RunRescanSubtree(args)
}

func expandHome(p string) string {
	return cmdscan.ExpandHome(p)
}

func resolveOutFile(outFile, outputDir, defaultName string) string {
	return cmdscan.ResolveOutFile(outFile, outputDir, defaultName)
}

type CleanOptions = cmddoctor.CleanOptions
type CleanResult = cmddoctor.CleanResult

func runDoctor(args []string) error {
	return cmddoctor.RunDoctorCmd(args)
}

func runCleanCorrupted(args []string) error {
	return cmddoctor.RunCleanCorrupted(args)
}

func CleanCorruptedDirs(opts cmddoctor.CleanOptions) (cmddoctor.CleanResult, error) {
	return cmddoctor.CleanCorruptedDirs(opts)
}

func runOS(args []string) error {
	return cmdos.RunOS(args)
}

func runOSFixLink(args []string) error {
	return cmdos.RunOSFixLink(args)
}

func RunOSCLI(args []string) error {
	return cmdos.RunOSCLI(args)
}

func runSchedule(args []string) error {
	return cmdschedule.RunSchedule(args)
}

func runExportConfig(args []string) error {
	return cmdconfig.RunExportConfig(args)
}

func runImportConfig(args []string) error {
	return cmdconfig.RunImportConfig(args)
}

func runWorkDir(args []string) error {
	return cmdworkdir.RunWorkDir(args)
}
