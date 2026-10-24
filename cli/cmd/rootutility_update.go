package cmd

import (
	"context"
	"os"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstaller"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
	"github.com/alimtvnetwork/gitmap-v28/cli/heatmap"
)

func runUpdateHelp() error {
	cmdName := subcommandName()
	args := argsTail()
	if handled, err := checkEarlyUpdateRouting(cmdName, args); handled {
		return err
	}
	checkHelp("update", argsTail())
	if handled, err := checkVersionListRouting(args); handled {
		return err
	}
	return executeUpdateHelpFlow(args)
}

func checkEarlyUpdateRouting(cmdName string, args []string) (bool, error) {
	if cmdupdate.IsFleetUpdateCommand(cmdName, args) {
		return true, cmdupdate.RunFleetUpdateDispatch(cmdName, args)
	}
	if len(args) > 0 && (args[0] == "ssh" || args[0] == "remote") {
		return true, cmdssh.RunSSHUpdateCLI(args[1:])
	}
	if hasSSHArg(args) {
		return true, cmdssh.RunSSHUpdateCLI(stripSSHFlag(args))
	}
	return false, nil
}

func checkVersionListRouting(args []string) (bool, error) {
	if isVersionListRequest(args) && isAgmUpdateTarget(args) {
		return true, cmdagy.RunAGMVersionTagsLS()
	}
	if isVersionListRequest(args) {
		return true, RunGitMapVersionTagsLS()
	}
	return false, nil
}

func executeUpdateHelpFlow(args []string) error {
	remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		return cmdssh.RunSSHUpdateCLI([]string{resolveUpdatePackage(cleanArgs), remoteTarget})
	}
	targetVer := extractVersionFromArgs(args)
	applyPinIfRequested(args, targetVer)
	if isAgmUpdateTarget(args) {
		return runUpdateAgManagerTarget(args)
	}
	if isAntigravityIdeUpdateTarget(args) {
		return runUpdateAntigravityIdeTarget(args)
	}
	configureTargetVersion(targetVer)
	if len(args) > 0 && !isKnownUpdateTargetOrFlag(args[0]) {
		handleUnknownUpdateTarget(args[0])
		return nil
	}
	return cmdupdate.RunUpdate()
}

func configureTargetVersion(targetVer string) {
	eff := resolveEffectiveVersion("gitmap", targetVer)
	if eff != "" {
		cmdupdate.SetTargetVersion(eff)
	}
}

func extractRemoteUpdateTarget(args []string) (string, []string) {
	for i := 0; i < len(args); i++ {
		if target, clean, ok := matchRemoteArg(args, i); ok {
			return target, clean
		}
	}
	return "", args
}

func matchRemoteArg(args []string, i int) (string, []string, bool) {
	arg := args[i]
	if slices.Contains([]string{"--remote", "-r", "--node", "-n"}, arg) && i+1 < len(args) {
		return args[i+1], append(append([]string{}, args[:i]...), args[i+2:]...), true
	}
	for _, pfx := range []string{"--remote=", "--node="} {
		if strings.HasPrefix(arg, pfx) {
			return strings.TrimPrefix(arg, pfx), append(append([]string{}, args[:i]...), args[i+1:]...), true
		}
	}
	return "", nil, false
}

func applyPinIfRequested(args []string, targetVer string) {
	if !hasPinArg(args) || targetVer == "" {
		return
	}
	appName := "gitmap"
	if isAgmUpdateTarget(args) {
		appName = "agm"
	} else if isAntigravityIdeUpdateTarget(args) {
		appName = "antigravity"
	}
	_ = cmdinstaller.PinVersion(appName, targetVer)
}

func resolveEffectiveVersion(appName, requestedVer string) string {
	if requestedVer != "" {
		return requestedVer
	}
	pinned, err := cmdinstaller.GetPinnedVersion(appName)
	if err == nil && pinned != "" {
		return pinned
	}
	return ""
}

func runUpdateAgManagerTarget(args []string) error {
	remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		return cmdssh.RunSSHUpdateCLI([]string{"agm", remoteTarget})
	}
	if isVersionListRequest(cleanArgs) {
		return cmdagy.RunAGMVersionTagsLS()
	}
	ver := resolveEffectiveVersion("agm", extractVersionFromArgs(cleanArgs))
	opts := cmdinstall.InstallOptions{
		DryRun:  hasDryRunArg(cleanArgs),
		Version: ver,
	}
	return cmdinstall.RunUpdateAgManagerWithOpts(opts)
}

func runUpdateAntigravityIdeTarget(args []string) error {
	remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
	if remoteTarget != "" {
		return cmdssh.RunSSHUpdateCLI([]string{"antigravity", remoteTarget})
	}
	ver := resolveEffectiveVersion("antigravity", extractVersionFromArgs(cleanArgs))
	if ver == "" {
		ver = resolveEffectiveVersion("ide", extractVersionFromArgs(cleanArgs))
	}
	opts := cmdinstall.InstallOptions{
		DryRun:         hasDryRunArg(cleanArgs),
		Verbose:        hasVerboseArg(cleanArgs),
		Yes:            hasYesArg(cleanArgs),
		Version:        ver,
		Prefix:         extractPrefixFromArgs(cleanArgs),
		Force:          true,
		IsDownloadMust: hasDownloadMustArg(cleanArgs),
	}
	return cmdinstall.RunInstallAntigravityWithOpts(opts)
}

const (
	CmdHeatmap      = "heatmap"
	CmdHeatmapAlias = "hm"
)

func init() {
	hookHeatmapDispatch()
}

func hookHeatmapDispatch() {
	if len(os.Args) < 2 || !isHeatmapCommand(os.Args[1]) {
		return
	}

	if err := heatmap.DispatchHeatmapSubsystem(context.Background(), os.Args[2:], os.Stdout); err != nil {
		cliexit.HandleError(err, 1)
	}

	cliexit.Exit(0)
}

func isHeatmapCommand(cmd string) bool {
	clean := strings.ToLower(strings.TrimSpace(cmd))
	return clean == CmdHeatmap || clean == CmdHeatmapAlias
}

// RunHeatmapCommand executes the heatmap subsystem programmatically.
func RunHeatmapCommand(args []string) error {
	return heatmap.DispatchHeatmapSubsystem(context.Background(), args, os.Stdout)
}
