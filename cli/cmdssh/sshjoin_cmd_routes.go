package cmdssh

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync/atomic"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/helpdoc"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

func routeSJSubcommands(ctx context.Context, args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	switch strings.ToLower(args[0]) {
	case "bootstrap", "bs":
		return true, RunClusterBootstrapCLI(args[1:])
	case "install":
		return true, RunClusterInstallCLI(args[1:])
	case "import-json", "importjson", "import":
		return true, RunSSHNodesImportJSON(args[1:])
	case "export-json", "exportjson", "export":
		return true, RunSSHNodesExportJSON(args[1:])
	case "export-oneliner", "oneliner", "eo":
		return true, RunSSHExportOnelinerCLI(args[1:])
	}
	if isSJNodesSubcommand(args) {
		return true, dispatchNodesArgs(ctx, args[1:])
	}
	if !isSJSubcommand(args[0]) {
		return false, nil
	}
	return true, dispatchSJSubcommand(ctx, args[0], args[1:])
}

func isSJNodesSubcommand(args []string) bool {
	if len(args) <= 1 {
		return false
	}
	sub := strings.ToLower(args[0])
	return sub == "nodes" || sub == "node"
}

func buildHostRecord(opts *SSHJoinOptions, now time.Time) store.SSHHost {
	return store.SSHHost{
		ID:        fmt.Sprintf("host-%s", opts.Target.IP),
		Alias:     opts.Alias,
		IP:        opts.Target.IP,
		Username:  opts.Target.Username,
		CreatedAt: now,
	}
}

var sshHistSeq uint64

func buildHistRecord(opts *SSHJoinOptions, now time.Time) store.SSHHistory {
	seq := atomic.AddUint64(&sshHistSeq, 1)
	return store.SSHHistory{
		ID:       fmt.Sprintf("hist-%s-%d-%d", opts.Target.IP, now.UnixNano(), seq),
		HostIP:   opts.Target.IP,
		JoinedAt: now,
		User:     opts.Target.Username,
	}
}

func buildHostAndHistory(opts *SSHJoinOptions) (store.SSHHost, store.SSHHistory) {
	now := time.Now().UTC()
	return buildHostRecord(opts, now), buildHistRecord(opts, now)
}

func pushAuthIfRequested(ctx context.Context, opts *SSHJoinOptions) error {
	if !opts.IsPushAuth {
		return nil
	}

	pubKey, err := getLocalPublicKey(ctx, "", false)
	if err != nil {
		return err
	}

	return appendKeyRemote(ctx, pubKey, *opts.Target)
}

func printEnrollSuccess(alias, target string) {
	fmt.Printf("✓ Machine '%s' (%s) joined successfully.\n", alias, target)
	fmt.Printf("  Recall anytime: gitmap ssh %s\n", alias)
	fmt.Printf("  Or connect directly: gitmap ssh %s\n", target)
	fmt.Println("  Undo anytime: gitmap ssh undo")
}

func resolveContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func enrollParsedTarget(ctx context.Context, opts *SSHJoinOptions) error {
	ctx = resolveContext(ctx)
	if opts.Target == nil {
		return apperror.NewValidationError(msgMissingJoinTarget)
	}

	return ExecuteSSHJoinEnrollment(ctx, opts)
}

func showJoinHelpAndExit() error {
	helpdoc.Print("ssh-join")
	return nil
}

func parseJoinArgs(args []string) (*SSHJoinOptions, error) {
	if len(args) == 0 {
		return nil, apperror.NewValidationError(msgMissingJoinTarget)
	}

	opts, err := parseSSHJoinOptions(args)
	if err != nil {
		return nil, apperror.NewValidationError(msgMissingJoinTarget)
	}

	return opts, nil
}

func findMultiTargetIndex(args []string) int {
	for i, a := range args {
		if !strings.HasPrefix(a, "-") && strings.Contains(a, ",") {
			return i
		}
	}
	return -1
}

func buildSingleTargetArgs(args []string, idx int, target string) []string {
	res := make([]string, len(args))
	copy(res, args)
	res[idx] = target
	return res
}

func executeMultiEnroll(ctx context.Context, args []string, targetIdx int) error {
	targets := ParseMultiIPList(args[targetIdx])
	for _, target := range targets {
		singleArgs := buildSingleTargetArgs(args, targetIdx, target)
		opts, err := parseJoinArgs(singleArgs)
		if err != nil {
			return err
		}
		if err := enrollParsedTarget(ctx, opts); err != nil {
			return err
		}
	}
	return nil
}

func isUserAtHostAddress(s string) bool {
	parts := strings.Split(s, "@")
	if len(parts) != 2 || len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	cleanHost := stripPortOrBrackets(parts[1])
	return net.ParseIP(cleanHost) != nil || strings.Contains(cleanHost, ".") || cleanHost == constants.Localhost
}

func isTargetAddress(s string) bool {
	if strings.Contains(s, "@") {
		return isUserAtHostAddress(s)
	}
	clean := stripPortOrBrackets(s)
	return net.ParseIP(clean) != nil
}

func extractTargetPositions(args []string) []string {
	var targets []string
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && isTargetAddress(a) {
			targets = append(targets, a)
		}
	}
	return targets
}

func extractFlagArgs(args []string) []string {
	var flags []string
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			flags = append(flags, a)
		}
	}
	return flags
}

func executeSpaceMultiEnroll(ctx context.Context, args []string) error {
	targets := extractTargetPositions(args)
	flags := extractFlagArgs(args)
	for _, target := range targets {
		singleArgs := append([]string{target}, flags...)
		opts, err := parseJoinArgs(singleArgs)
		if err != nil {
			return err
		}
		if err := enrollParsedTarget(ctx, opts); err != nil {
			return err
		}
	}
	return nil
}

func executeEnrollCLI(ctx context.Context, args []string) error {
	ctx = resolveContext(ctx)
	if hasHelpFlag(args) {
		return showJoinHelpAndExit()
	}

	targetIdx := findMultiTargetIndex(args)
	if targetIdx != -1 {
		return executeMultiEnroll(ctx, args, targetIdx)
	}

	targets := extractTargetPositions(args)
	if len(targets) > 1 {
		return executeSpaceMultiEnroll(ctx, args)
	}

	opts, err := parseJoinArgs(args)
	if err != nil {
		return err
	}

	if opts.IsShowHelp {
		return showJoinHelpAndExit()
	}

	return enrollParsedTarget(ctx, opts)
}

func runSSHJoinCLI(args []string) error {
	if hasHelpFlag(args) {
		return showJoinHelpAndExit()
	}

	ctx := context.Background()
	isHandled, err := routeSJSubcommands(ctx, args)
	if isHandled {
		return err
	}

	return executeEnrollCLI(ctx, args)
}

var SJAuthKeyCmd = &cobra.Command{
	Use:   "auth-key [deploy] [target] [flags]",
	Short: "Deploy SSH public key to target machines",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHAuthKeyDeployCLI(args)
	},
}
