package cmdssh

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	dbpkg "github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var (
	runSSHJoinFn = RunSSHJoinCLI
	spawnSSHFn   = SpawnSSHWithPassword
	openSSHDB    = openDB
)

// SSHLoginCmd represents the gitmap ssh login command.
var SSHLoginCmd = &cobra.Command{
	Use:   "login [target] [password]",
	Short: "Login via SSH to the specified target",
	Args:  cobra.RangeArgs(1, 2),
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

func isNodesSubcommand(cmd string) bool {
	return cmd == "nodes" || cmd == "node" || cmd == "ls"
}

func extractLoginPassword(args []string) string {
	hasPassword := len(args) >= 2
	if hasPassword {
		return args[1]
	}
	return ""
}

//nolint:revive
func runSSHLogin(cmd *cobra.Command, args []string, ctx context.Context) error {
	isEmptyArgs := len(args) < 1
	if isEmptyArgs {
		return apperror.NewSimple("runSSHLogin", "E_INTERNAL_ERROR")
	}

	if isJoinSubcommand(args[0]) {
		return runSSHJoinFn(args[1:])
	}

	if isNodesSubcommand(args[0]) {
		return RunSSHNodesCLI(ctx, args[1:])
	}

	pass := extractLoginPassword(args)
	return executeSSHLoginWithPassword(ctx, args[0], pass, false)
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
		return
	}
	hasHostIP := sshTarget != nil && sshTarget.IP != "" && net.ParseIP(sshTarget.IP) != nil
	if hasHostIP {
		resolveIPCredentials(ctx, sshTarget.IP, sshTarget)
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

func isDefaultOrEmptyUser(u string) bool {
	return u == "" || u == "root"
}

func assignHostToSSHTarget(host store.SSHHost, sshTarget *SSHTarget) {
	if isDefaultOrEmptyUser(sshTarget.Username) && host.Username != "" {
		sshTarget.Username = host.Username
	}
	sshTarget.Port = host.Port
	if sshTarget.EncryptedPassword == "" {
		sshTarget.EncryptedPassword = host.EncryptedPassword
	}
}

func assignConnToSSHTarget(conn *dbpkg.SSHConnection, sshTarget *SSHTarget) {
	if isDefaultOrEmptyUser(sshTarget.Username) && conn.Username != "" {
		sshTarget.Username = conn.Username
	}
	if sshTarget.EncryptedPassword == "" {
		sshTarget.EncryptedPassword = conn.EncryptedPassword
	}
}

func applyEnrolledIPHost(ctx context.Context, target string, db *store.DB, sshTarget *SSHTarget) {
	host, err := store.GetHostByIP(ctx, target, db.Conn())
	if err == nil {
		assignHostToSSHTarget(host, sshTarget)
		return
	}
	conn, connErr := dbpkg.GetSSHConnectionByIP(ctx, db.Conn(), target)
	if connErr == nil && conn != nil {
		assignConnToSSHTarget(conn, sshTarget)
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

func saveExplicitPassword(ctx context.Context, target string, sshTarget *SSHTarget, pass string) {
	encPass, err := EncryptSSHPassword(pass)
	if err != nil {
		return
	}
	db, err := openSSHDB()
	if err != nil {
		return
	}
	defer db.Close()
	updateAndInsertPassword(ctx, target, sshTarget, encPass, db)
}

func updateAndInsertPassword(ctx context.Context, target string, sshTarget *SSHTarget, encPass string, db *store.DB) {
	rows1, _ := store.UpdateHostPassword(ctx, target, encPass, db.Conn())
	_, _ = dbpkg.UpdateSSHConnectionPassword(ctx, db.Conn(), target, encPass)
	hasExtraIP := sshTarget != nil && sshTarget.IP != "" && sshTarget.IP != target
	if hasExtraIP {
		rows2, _ := store.UpdateHostPassword(ctx, sshTarget.IP, encPass, db.Conn())
		_, _ = dbpkg.UpdateSSHConnectionPassword(ctx, db.Conn(), sshTarget.IP, encPass)
		rows1 += rows2
	}
	isUnregistered := rows1 == 0 && sshTarget != nil
	if isUnregistered {
		insertSSHHostFallback(ctx, target, sshTarget, encPass, db)
	}
}

func insertSSHHostFallback(ctx context.Context, target string, sshTarget *SSHTarget, encPass string, db *store.DB) {
	host := store.SSHHost{
		ID:                fmt.Sprintf("host-%s", sshTarget.IP),
		Alias:             target,
		IP:                sshTarget.IP,
		Username:          sshTarget.Username,
		Port:              sshTarget.Port,
		EncryptedPassword: encPass,
		ClusterRole:       "worker",
		CreatedAt:         time.Now().UTC(),
	}
	_ = store.InsertSSHHost(ctx, host, db.Conn())
}

func resolvePassword(ctx context.Context, target string, sshTarget *SSHTarget, explicitPass string) string {
	stored := resolveTargetPassword(sshTarget)
	if stored != "" {
		return stored
	}
	if explicitPass != "" {
		saveExplicitPassword(ctx, target, sshTarget, explicitPass)
		return explicitPass
	}
	return ""
}

func executeSSHLogin(ctx context.Context, target string, force bool) error {
	return executeSSHLoginWithPassword(ctx, target, "", force)
}

func executeSSHLoginWithPassword(ctx context.Context, target string, explicitPass string, force bool) error {
	sshTarget, err := ParseSSHTarget(target, "root", 22)
	if err != nil {
		return err
	}

	checkAndResolveIP(ctx, target, sshTarget)
	if err := checkAndResolveAlias(ctx, target, sshTarget); err != nil {
		return err
	}

	autoTrustTargetHostFn(ctx, sshTarget)
	password := resolvePassword(ctx, target, sshTarget, explicitPass)
	probeAndEnsureNodeProfile(ctx, target, sshTarget, password)

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
	conn, connErr := dbpkg.GetSSHConnectionByAlias(ctx, db.Conn(), target)
	if connErr == nil && conn != nil {
		assignConnToSSHTarget(conn, sshTarget)
		sshTarget.IP = conn.IPAddress
		sshTarget.Port = 22
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
	_ = RenderSSHHostsTable(&sb, hosts)
	return sb.String()
}

func formatJoinExamples(target string) string {
	var sb strings.Builder
	sb.WriteString("To enroll this machine in your SSH registry:\n")
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join user@<ip> %s\n", target))
	sb.WriteString("  gitmap ssh-join user@<ip>\n")
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join <ip> %s\n", target))
	sb.WriteString(fmt.Sprintf("  gitmap ssh-join add user@<ip> %s\n", target))
	sb.WriteString(fmt.Sprintf("  gitmap ssh join user@<ip> %s\n", target))
	sb.WriteString("\nTo recall an existing registered host:\n")
	sb.WriteString("  gitmap ssh <alias>\n")
	sb.WriteString("  gitmap ssh nodes\n")
	sb.WriteString("  gitmap ssh ls\n")
	sb.WriteString("  gitmap ssh-join ls\n")
	return sb.String()
}

func suggestSSHSubcommand(target string) string {
	low := strings.ToLower(target)
	switch low {
	case "hosts", "host", "machines", "vms", "node":
		return "nodes"
	case "updat", "up", "upgrade":
		return "update"
	case "exe", "cmd", "run":
		return "exec"
	case "inst", "install", "setup":
		return "install-exec"
	case "joi", "enroll":
		return "join"
	case "err", "error", "logs":
		return "error-logs"
	case "stat", "check", "ping":
		return "health"
	}
	return ""
}

func formatAliasNotFoundMessage(target string, hosts []store.SSHHost) string {
	header := fmt.Sprintf("SSH host alias '%s' not found in registry.\n\n", target)
	if suggestion := suggestSSHSubcommand(target); suggestion != "" {
		header = fmt.Sprintf("SSH host alias '%s' not found in registry.\n  Did you mean: gitmap ssh %s?\n\n", target, suggestion)
	}
	table := formatRegisteredHostsTable(hosts)
	examples := formatJoinExamples(target)
	return header + table + "\n" + examples
}
