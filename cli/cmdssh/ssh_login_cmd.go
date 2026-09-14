package cmdssh

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	runSSHJoinFn = RunSSHJoinCLI
	spawnSSHFn   = SpawnSSHWithPassword
	openSSHDB    = openDB
)

// SSHLoginCmd represents the gitmap ssh login command.
var SSHLoginCmd = &cobra.Command{
	Use:   "login [target]",
	Short: "Login via SSH to the specified target",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSSHLogin(cmd, args, cmd.Context())
	},
}

// RunSSHLogin executes the ssh login subcommand.
//
//nolint:revive
func RunSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error {
	return runSSHLogin(cmd, args, ctx)
}

func isJoinSubcommand(cmd string) bool {
	return cmd == "join" || cmd == "sj"
}

//nolint:revive
func runSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error {
	if len(args) < 1 {
		return apperror.NewSimple("runSSHLogin", "E_INTERNAL_ERROR")
	}

	if isJoinSubcommand(args[0]) {
		return runSSHJoinFn(args[1:])
	}

	return executeSSHLogin(ctx, args[0], false)
}

func isAliasTarget(target string) bool {
	hasAtSign := strings.Contains(target, "@")
	if hasAtSign {
		return false
	}
	isIP := net.ParseIP(target) != nil
	if isIP {
		return false
	}
	return true
}

func isPlainIPTarget(target string) bool {
	hasAtSign := strings.Contains(target, "@")
	if hasAtSign {
		return false
	}
	return net.ParseIP(target) != nil
}

func checkAndResolveIP(ctx context.Context, target string, sshTarget *SSHTarget) {
	isIP := isPlainIPTarget(target)
	if isIP {
		resolveIPCredentials(ctx, target, sshTarget)
	}
}

func resolveIPCredentials(ctx context.Context, target string, sshTarget *SSHTarget) {
	db, err := openSSHDB()
	if err != nil {
		return
	}
	defer db.Close()
	applyEnrolledIPHost(ctx, target, db, sshTarget)
}

func applyEnrolledIPHost(ctx context.Context, target string, db *store.DB, sshTarget *SSHTarget) {
	host, err := store.GetHostByIP(ctx, target, db.Conn())
	if err == nil {
		sshTarget.Username = host.Username
		sshTarget.Port = host.Port
		sshTarget.EncryptedPassword = host.EncryptedPassword
	}
}

func resolveTargetPassword(sshTarget *SSHTarget) string {
	if sshTarget == nil || sshTarget.EncryptedPassword == "" {
		return ""
	}
	plain, err := DecryptSSHPassword(sshTarget.EncryptedPassword)
	if err == nil {
		return plain
	}
	return ""
}

func executeSSHLogin(ctx context.Context, target string, force bool) error {
	sshTarget, err := ParseSSHTarget(target, "root", 22)
	if err != nil {
		return err
	}
	checkAndResolveIP(ctx, target, sshTarget)
	if err := checkAndResolveAlias(ctx, target, sshTarget); err != nil {
		return err
	}
	password := resolveTargetPassword(sshTarget)
	return spawnSSHFn(ctx, *sshTarget, nil, password)
}

func checkAndResolveAlias(ctx context.Context, target string, sshTarget *SSHTarget) error {
	isAlias := isAliasTarget(target)
	if isAlias {
		return resolveAliasOrReport(ctx, target, sshTarget)
	}
	return nil
}

func resolveAliasOrReport(ctx context.Context, target string, sshTarget *SSHTarget) error {
	db, err := openSSHDB()
	if err != nil {
		return apperror.Wrap(err, "resolveAliasOrReport", map[string]any{"target": target})
	}
	defer db.Close()
	return lookupSSHHostOrReport(ctx, target, db, sshTarget)
}

func lookupSSHHostOrReport(ctx context.Context, target string, db *store.DB, sshTarget *SSHTarget) error {
	host, err := store.GetHostByAlias(ctx, target, db.Conn())
	if err == nil {
		sshTarget.Username = host.Username
		sshTarget.IP = host.IP
		sshTarget.Port = host.Port
		sshTarget.EncryptedPassword = host.EncryptedPassword
		return nil
	}
	return reportAliasNotFound(ctx, target, db)
}

func reportAliasNotFound(ctx context.Context, target string, db *store.DB) error {
	hosts, _ := store.ListHosts(ctx, db.Conn())
	msg := formatAliasNotFoundMessage(target, hosts)
	return apperror.NewNotFoundError(msg)
}

func formatRegisteredHostsTable(hosts []store.SSHHost) string {
	if len(hosts) == 0 {
		return "Currently registered hosts:\n  (none)\n"
	}
	var sb strings.Builder
	sb.WriteString("Currently registered hosts:\n")
	sb.WriteString(fmt.Sprintf("  %-16s %-16s %-12s\n", "ALIAS", "IP", "USER"))
	for _, h := range hosts {
		sb.WriteString(fmt.Sprintf("  %-16s %-16s %-12s\n", h.Alias, h.IP, h.Username))
	}
	return sb.String()
}

func formatJoinExamples(target string) string {
	var sb strings.Builder
	sb.WriteString("To enroll this machine in your SSH registry:\n")
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join user@<ip> %s\n", target))
	sb.WriteString("  gitmap ssh-join user@<ip>\n")
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join <ip> %s\n", target))
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join add user@<ip> %s\n", target))
	sb.WriteString("\nTo recall an existing registered host:\n")
	sb.WriteString("  gitmap ssh <alias>\n")
	sb.WriteString("  gitmap ssh-join ls\n")
	return sb.String()
}

func formatAliasNotFoundMessage(target string, hosts []store.SSHHost) string {
	header := fmt.Sprintf("SSH host alias '%s' not found in registry.\n\n", target)
	table := formatRegisteredHostsTable(hosts)
	examples := formatJoinExamples(target)
	return header + table + "\n" + examples
}

func resolveSSHHostTarget(ctx context.Context, target string, sshTarget *SSHTarget) {
	_ = resolveAliasOrReport(ctx, target, sshTarget)
}
