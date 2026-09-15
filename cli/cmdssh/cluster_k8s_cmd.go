package cmdssh

import (
	"context"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ClusterK8sCmd represents the gitmap cluster k8s command.
var ClusterK8sCmd = &cobra.Command{
	Use:     "k8s <subcommand> [args...]",
	Aliases: []string{"kube", "kubernetes"},
	Short:   "Manage Kubernetes cluster lifecycle, runtime, and components",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunClusterK8sCLI(args)
	},
}

func printClusterK8sSubcommands() {
	fmt.Println("  prereq <target>                                              Configure overlay, br_netfilter and sysctl")
	fmt.Println("  install <target> [--version <ver>] [--hostname <name>]       Install CRI-O, kubelet, kubeadm, kubectl")
	fmt.Println("  init <target> [--pod-cidr <cidr>]                            Initialize control plane with kubeadm")
	fmt.Println("  cni <target> [--plugin weave|calico]                         Deploy pod network add-on")
	fmt.Println("  join-command <target>                                        Print join command from control plane")
	fmt.Println("  join <target> [--command \"<cmd>\"]                            Join worker node(s) to cluster")
}

func printClusterK8sAddonSubcommands() {
	fmt.Println("  nfs <target> [--export-dir <dir>]                            Install NFS kernel server and export dir")
	fmt.Println("  helm-install <target> [--version <ver>]                      Install Helm binary on target")
	fmt.Println("  helm-nfs <target> [--nfs-server <ip>] [--export-dir <dir>]   Deploy NFS subdir provisioner via Helm")
	fmt.Println("  status <target>                                              Show cluster and service status")
	fmt.Println("  reset <target>                                               Tear down node with kubeadm reset --force")
}

func showClusterK8sHelp() error {
	fmt.Println("Usage: gitmap cluster k8s <subcommand> [args...]\n\nSubcommands:")
	printClusterK8sSubcommands()
	printClusterK8sAddonSubcommands()
	fmt.Println("\nFlags:\n  -h, --help    Show this help message")
	return nil
}

func isClusterK8sHelp(args []string) bool {
	if len(args) == 0 {
		return true
	}
	for _, arg := range args {
		if isClusterHelpFlag(arg) {
			return true
		}
	}
	return false
}

func getNextK8sArg(args []string, idx int) string {
	if idx+1 < len(args) {
		return args[idx+1]
	}
	return ""
}

func parseStringFlag(arg, nextArg, longFlag, shortFlag string) (string, int) {
	prefix := "--" + longFlag + "="
	if strings.HasPrefix(arg, prefix) {
		return strings.TrimPrefix(arg, prefix), 1
	}
	isMatch := arg == "--"+longFlag || (shortFlag != "" && arg == "-"+shortFlag)
	if isMatch && nextArg != "" {
		return nextArg, 2
	}
	return "", 0
}

func parseK8sFlagStep(args []string, idx int, longFlag, shortFlag string) (string, int) {
	arg, next := args[idx], getNextK8sArg(args, idx)
	return parseStringFlag(arg, next, longFlag, shortFlag)
}

func extractTargetAndFlag(args []string, longFlag, shortFlag string) ([]string, string) {
	var pos []string
	var val string
	for idx := 0; idx < len(args); idx++ {
		fVal, n := parseK8sFlagStep(args, idx, longFlag, shortFlag)
		if n > 0 {
			val, idx = fVal, idx+n-1
			continue
		}
		pos = append(pos, args[idx])
	}
	return pos, val
}

func parseTargetAndFlag(args []string, longFlag, shortFlag, usage string) (string, string, error) {
	pos, val := extractTargetAndFlag(args, longFlag, shortFlag)
	if len(pos) < 1 {
		return "", "", apperror.NewValidationError(usage)
	}
	return pos[0], val, nil
}

func extractTwoFlagsStep(args []string, idx int, f1, s1, f2, s2 string) (string, string, int) {
	if v, n := parseK8sFlagStep(args, idx, f1, s1); n > 0 {
		return v, "", n
	}
	if v, n := parseK8sFlagStep(args, idx, f2, s2); n > 0 {
		return "", v, n
	}
	return "", "", 0
}

func pickFlagVal(current, newVal string) string {
	if newVal != "" {
		return newVal
	}
	return current
}

func extractTargetAndTwoFlags(args []string, f1, s1, f2, s2 string) ([]string, string, string) {
	var pos []string
	var v1, v2 string
	for idx := 0; idx < len(args); idx++ {
		r1, r2, n := extractTwoFlagsStep(args, idx, f1, s1, f2, s2)
		if n > 0 {
			v1, v2, idx = pickFlagVal(v1, r1), pickFlagVal(v2, r2), idx+n-1
			continue
		}
		pos = append(pos, args[idx])
	}
	return pos, v1, v2
}

func parseTargetAndTwoFlags(args []string, f1, s1, f2, s2, usage string) (string, string, string, error) {
	pos, v1, v2 := extractTargetAndTwoFlags(args, f1, s1, f2, s2)
	if len(pos) < 1 {
		return "", "", "", apperror.NewValidationError(usage)
	}
	return pos[0], v1, v2, nil
}

func parseSingleK8sTarget(subcmd string, args []string) (string, error) {
	if len(args) < 1 {
		return "", apperror.NewValidationError(fmt.Sprintf("usage: gitmap cluster k8s %s <target>", subcmd))
	}
	return args[0], nil
}

func cleanJoinOutput(output string) string {
	res := strings.ReplaceAll(output, "\\\r\n", " ")
	return strings.ReplaceAll(res, "\\\n", " ")
}

func findKubeadmJoinLine(output string) string {
	cleaned := cleanJoinOutput(output)
	for _, line := range strings.Split(cleaned, "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "kubeadm join") {
			return trimmed
		}
	}
	return ""
}

func extractKubeadmJoinCommand(output string) (string, error) {
	cmd := findKubeadmJoinLine(output)
	if cmd != "" {
		return cmd, nil
	}
	idx := strings.Index(output, "kubeadm join")
	if idx != -1 {
		return strings.TrimSpace(cleanJoinOutput(output[idx:])), nil
	}
	return "", apperror.NewExecutionError("failed to extract kubeadm join command from output")
}

func resolveControlPlaneHosts(ctx context.Context) ([]store.SSHHost, error) {
	dbConn, err := openClusterDBFunc()
	if err != nil {
		return nil, apperror.WrapSimple(err, "openClusterDB")
	}
	defer dbConn.Close()
	hosts, err := store.ListHostsByRole(ctx, "control", dbConn.Conn())
	if err != nil {
		return nil, apperror.Wrap(err, "ListHostsByRole", map[string]any{"role": "control"})
	}
	if len(hosts) == 0 {
		return nil, apperror.NewNotFoundError("no control plane hosts found")
	}
	return hosts, nil
}

func runK8sScriptOnHosts(ctx context.Context, hosts []store.SSHHost, script string) []ClusterRunResult {
	results := make([]ClusterRunResult, 0, len(hosts))
	for _, host := range hosts {
		res := executeNodeCmdFn(ctx, host, script, true)
		results = append(results, res)
	}
	return results
}

func executeClusterK8sScript(ctx context.Context, target string, script string) error {
	hosts, err := resolveClusterHosts(ctx, target)
	if err != nil {
		return err
	}
	results := runK8sScriptOnHosts(ctx, hosts, script)
	PrintClusterSummaryTable(results)
	if hasExecutionFailures(results) {
		return apperror.NewExecutionError("one or more cluster nodes failed k8s execution")
	}
	return nil
}

func fetchKubeadmJoinCmd(ctx context.Context) (string, error) {
	ctrlHosts, err := resolveControlPlaneHosts(ctx)
	if err != nil {
		return "", err
	}
	res := executeNodeCmdFn(ctx, ctrlHosts[0], "kubeadm token create --print-join-command", true)
	if res.ExitCode != 0 || res.Err != nil {
		return "", apperror.NewExecutionError(fmt.Sprintf("failed to get join command from control plane: %s", res.Stderr))
	}
	return extractKubeadmJoinCommand(res.Stdout)
}

func runK8sAutoJoin(ctx context.Context, target string) error {
	joinCmd, err := fetchKubeadmJoinCmd(ctx)
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateK8sJoinScript(joinCmd))
}

func generateK8sStatusScript() string {
	return `#!/bin/bash
if [ -f /etc/kubernetes/admin.conf ]; then
  export KUBECONFIG=/etc/kubernetes/admin.conf
  kubectl cluster-info
  echo ""
  kubectl get nodes -o wide
else
  systemctl status kubelet --no-pager || true
fi
`
}

func runK8sPrereq(ctx context.Context, args []string) error {
	target, err := parseSingleK8sTarget("prereq", args)
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateK8sPrereqScript())
}

func runK8sInstall(ctx context.Context, args []string) error {
	target, ver, host, err := parseTargetAndTwoFlags(args, "version", "v", "hostname", "",
		"usage: gitmap cluster k8s install <target> [--version <ver>] [--hostname <name>]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateK8sInstallScript(ver, host))
}

func runK8sInit(ctx context.Context, args []string) error {
	target, cidr, err := parseTargetAndFlag(args, "pod-cidr", "",
		"usage: gitmap cluster k8s init <target> [--pod-cidr <cidr>]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateK8sInitScript(cidr))
}

func runK8sCNI(ctx context.Context, args []string) error {
	target, plugin, err := parseTargetAndFlag(args, "plugin", "p",
		"usage: gitmap cluster k8s cni <target> [--plugin weave|calico]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateCNIScript(plugin))
}

func runK8sJoinCommand(ctx context.Context, args []string) error {
	target, err := parseSingleK8sTarget("join-command", args)
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, "kubeadm token create --print-join-command")
}

func runK8sJoin(ctx context.Context, args []string) error {
	target, cmd, err := parseTargetAndFlag(args, "command", "c",
		"usage: gitmap cluster k8s join <target> [--command \"<cmd>\"]")
	if err != nil {
		return err
	}
	if cmd != "" {
		return executeClusterK8sScript(ctx, target, GenerateK8sJoinScript(cmd))
	}
	return runK8sAutoJoin(ctx, target)
}

func runK8sNFS(ctx context.Context, args []string) error {
	target, dir, err := parseTargetAndFlag(args, "export-dir", "d",
		"usage: gitmap cluster k8s nfs <target> [--export-dir <dir>]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateNFSServerScript(dir))
}

func runK8sHelmInstall(ctx context.Context, args []string) error {
	target, ver, err := parseTargetAndFlag(args, "version", "v",
		"usage: gitmap cluster k8s helm-install <target> [--version <ver>]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateHelmInstallScript(ver))
}

func runK8sHelmNFS(ctx context.Context, args []string) error {
	target, server, dir, err := parseTargetAndTwoFlags(args, "nfs-server", "s", "export-dir", "d",
		"usage: gitmap cluster k8s helm-nfs <target> [--nfs-server <ip>] [--export-dir <dir>]")
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateHelmNFSScript(server, dir))
}

func runK8sStatus(ctx context.Context, args []string) error {
	target, err := parseSingleK8sTarget("status", args)
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, generateK8sStatusScript())
}

func runK8sReset(ctx context.Context, args []string) error {
	target, err := parseSingleK8sTarget("reset", args)
	if err != nil {
		return err
	}
	return executeClusterK8sScript(ctx, target, GenerateK8sResetScript())
}

func routeClusterK8sBootstrap(ctx context.Context, sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "prereq":
		return result.MatchWrapper(runK8sPrereq(ctx, rest))
	case "install":
		return result.MatchWrapper(runK8sInstall(ctx, rest))
	case "init":
		return result.MatchWrapper(runK8sInit(ctx, rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterK8sNetwork(ctx context.Context, sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "cni":
		return result.MatchWrapper(runK8sCNI(ctx, rest))
	case "join-command":
		return result.MatchWrapper(runK8sJoinCommand(ctx, rest))
	case "join":
		return result.MatchWrapper(runK8sJoin(ctx, rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterK8sAddons(ctx context.Context, sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "nfs":
		return result.MatchWrapper(runK8sNFS(ctx, rest))
	case "helm-install":
		return result.MatchWrapper(runK8sHelmInstall(ctx, rest))
	case "helm-nfs":
		return result.MatchWrapper(runK8sHelmNFS(ctx, rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterK8sAdmin(ctx context.Context, sub string, rest []string) result.ErrorWrapper {
	switch sub {
	case "status":
		return result.MatchWrapper(runK8sStatus(ctx, rest))
	case "reset":
		return result.MatchWrapper(runK8sReset(ctx, rest))
	default:
		return result.UnmatchedWrapper()
	}
}

func routeClusterK8sCommand(ctx context.Context, sub string, rest []string) result.ErrorWrapper {
	resBoot := routeClusterK8sBootstrap(ctx, sub, rest)
	if resBoot.IsMatched() {
		return resBoot
	}

	resNet := routeClusterK8sNetwork(ctx, sub, rest)
	if resNet.IsMatched() {
		return resNet
	}

	resAddon := routeClusterK8sAddons(ctx, sub, rest)
	if resAddon.IsMatched() {
		return resAddon
	}

	return routeClusterK8sAdmin(ctx, sub, rest)
}

// RunClusterK8sCLI dispatches Kubernetes cluster lifecycle recipes across target hosts.
func RunClusterK8sCLI(args []string) error {
	if isClusterK8sHelp(args) {
		return showClusterK8sHelp()
	}

	ctx := context.Background()
	res := routeClusterK8sCommand(ctx, args[0], args[1:])
	if res.IsMatched() {
		return res.AppError()
	}

	return apperror.NewValidationError(fmt.Sprintf("unknown cluster k8s subcommand: %s", args[0]))
}
