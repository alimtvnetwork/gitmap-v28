package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"sync/atomic"
	"text/tabwriter"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var openSSHDBFunc = openDB

const msgMissingJoinTarget = `missing target host

Usage:
  gitmap ssh-join <user@ip|ip> [alias] [flags]
  gitmap ssh-join add <user@ip|ip> [alias] [flags]
  gitmap ssh-join add-with-pass <user@ip|ip> [password] [alias] [flags]
  gitmap ssh join <user@ip|ip> [alias] [flags]
  gitmap sj <user@ip|ip> [alias] [flags]

Examples:
  gitmap ssh-join user@192.168.1.14
  gitmap ssh-join root@192.168.1.14 prod-server
  gitmap ssh-join add dev@192.168.1.50 devbox
  gitmap ssh-join add-with-pass alim@192.168.1.14 secret123 devbox
  gitmap ssh-join 192.168.1.14
  gitmap ssh join alim@192.168.1.14 devbox
  gitmap ssh join ubuntu@192.168.1.14:2222 prod --auth`

func executeSSHJoin(ctx context.Context, target string, history store.SSHHistory) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to open db", "err": err.Error()})
	}

	defer dbConn.Close()

	return runJoinTransaction(ctx, dbConn.SQL(), target, history)
}

func executeJoinInTx(ctx context.Context, wrapper *dbengine.DbWrapper, target string, history store.SSHHistory) error {
	appErr := wrapper.WithTransaction(ctx, func(tx *dbengine.TxWrapper) *apperror.AppError {
		return insertJoinRecords(ctx, tx, target, history)
	})
	if appErr != nil {
		return appErr
	}

	return nil
}

func runJoinTransaction(ctx context.Context, db *sql.DB, target string, history store.SSHHistory) error {
	wrapper, appErr := dbengine.WrapDb(db, dbengine.DbSQLite)
	if appErr != nil {
		return appErr
	}

	if txErr := executeJoinInTx(ctx, wrapper, target, history); txErr != nil {
		return txErr
	}

	fmt.Println("Joined successfully")
	return nil
}

func insertJoinRecords(ctx context.Context, tx *dbengine.TxWrapper, target string, history store.SSHHistory) *apperror.AppError {
	host := store.SSHHost{
		ID:       history.ID,
		IP:       history.HostIP,
		Username: history.User,
		Alias:    target,
	}

	if err := store.InsertSSHHost(ctx, host, tx.Tx()); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "insert host error", "err": err.Error()})
	}

	return logSSHJoinInTx(ctx, tx, history)
}

func logSSHJoinInTx(ctx context.Context, tx *dbengine.TxWrapper, history store.SSHHistory) *apperror.AppError {
	query := `INSERT INTO ssh_history (id, host_ip, joined_at, user) VALUES (?, ?, ?, ?)`
	if _, err := tx.Exec(ctx, query, history.ID, history.HostIP, history.JoinedAt, history.User); err != nil {
		return apperror.New("executeSSHJoin", "E_INTERNAL_ERROR", map[string]any{"msg": "log join error", "err": err.Error()})
	}

	return nil
}

func isSJClusterImportSubcommand(sub string) bool {
	return sub == "import-cluster" || sub == "import" || sub == "import-config"
}

func isSJAddWithPassSubcommand(sub string) bool {
	return sub == "add-with-pass" || sub == "add-pass" || sub == "add-password"
}

func isSJAddSubcommand(sub string) bool {
	return sub == "add" || sub == "join" || sub == "new" || sub == "enroll"
}

func isSJBootstrapSubcommand(sub string) bool {
	return sub == "bootstrap" || sub == "bs"
}

func isSJInstallSubcommand(sub string) bool {
	return sub == "install"
}

func routeSJClusterSpecialSubcommand(args []string) (bool, error) {
	if isSJInstallSubcommand(args[0]) {
		return true, RunClusterInstallCLI(args[1:])
	}
	if isSJBootstrapSubcommand(args[0]) {
		return true, RunClusterBootstrapCLI(args[1:])
	}
	if isSJClusterImportSubcommand(args[0]) {
		return true, RunClusterImportCLI(args[1:])
	}
	return false, nil
}

func routeSSHJoinSpecialSubcommand(cmd *cobra.Command, args []string) (bool, error) {
	if isHandled, err := routeSJClusterSpecialSubcommand(args); isHandled {
		return true, err
	}
	if isSJAddWithPassSubcommand(args[0]) {
		return true, executeEnrollWithPassCLI(resolveContext(cmd.Context()), args[1:])
	}
	if isSJAddSubcommand(args[0]) {
		return true, executeEnrollCLI(resolveContext(cmd.Context()), args[1:])
	}
	return false, nil
}

func routeSSHJoinCmd(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return RunSSHJoinCLI(args)
	}
	if isHandled, err := routeSSHJoinSpecialSubcommand(cmd, args); isHandled {
		return err
	}
	if !isSJSubcommand(args[0]) {
		return executeEnrollCLI(resolveContext(cmd.Context()), args)
	}
	return RunSSHJoinCLI(args)
}

var SSHJoinCmd = &cobra.Command{
	Use:     "ssh-join [user@ip|ip] [alias]",
	Aliases: []string{"sj", "ssh-joined", "ssh-joiner"},
	Short:   "Join an SSH machine by user@ip or IP address",
	Args:    cobra.ArbitraryArgs,
	RunE:    routeSSHJoinCmd,
}

//nolint:revive
func runSSHJoin(cmd *cobra.Command, args []string, ctx context.Context) error {
	return RunSSHJoinCLI(args)
}

var SJRmCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove machine-alias or ip",
	RunE: func(cmd *cobra.Command, args []string) error {
		return runSJRm(cmd, args, cmd.Context())
	},
}

func isSJScanSubcommand(sub string) bool {
	return sub == "scan" || sub == "find" || sub == "discover" || sub == "probe"
}

func isSJStatusSubcommand(sub string) bool {
	return sub == "status" || sub == "ping" || sub == "health" || sub == "check"
}

func isSJSubcommand(sub string) bool {
	if isSJAddSubcommand(sub) || isSJAddWithPassSubcommand(sub) || isSJScanSubcommand(sub) || isSJStatusSubcommand(sub) || isSJClusterImportSubcommand(sub) || isSJBootstrapSubcommand(sub) || isSJInstallSubcommand(sub) {
		return true
	}

	return sub == "ls" || sub == "list" || sub == "rm" || sub == "remove" || sub == "delete" ||
		sub == "add-auth" || sub == "auth" || sub == "history" || sub == "hist"
}

func renderSJListTable(out io.Writer, hosts []store.SSHHost) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tALIAS\tIP\tUSERNAME\tCREATED_AT")

	for _, host := range hosts {
		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			host.ID, host.Alias, host.IP, host.Username, host.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	return w.Flush()
}

func executeSJList(ctx context.Context) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return printSJList(ctx, os.Stdout, 0)
	}

	defer dbConn.Close()

	hosts, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		return apperror.New("executeSJList", "E_INTERNAL_ERROR", map[string]any{"msg": "failed to list hosts", "err": err.Error()})
	}

	return renderSJListTable(os.Stdout, hosts)
}

func dispatchSJBasicSubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if isSJInstallSubcommand(sub) {
		return true, RunClusterInstallCLI(args)
	}
	if isSJClusterImportSubcommand(sub) {
		return true, RunClusterImportCLI(args)
	}
	if isSJAddWithPassSubcommand(sub) {
		return true, executeEnrollWithPassCLI(ctx, args)
	}
	if isSJAddSubcommand(sub) {
		return true, executeEnrollCLI(ctx, args)
	}
	return false, nil
}

func dispatchSJQuerySubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if isSJScanSubcommand(sub) {
		return true, RunSJScan(SJScanCmd, args, ctx)
	}
	if isSJStatusSubcommand(sub) {
		return true, RunSJStatus(SJStatusCmd, args, ctx)
	}
	if sub == "ls" || sub == "list" {
		return true, executeSJList(ctx)
	}
	return false, nil
}

func dispatchSJActionSubcommand(ctx context.Context, sub string, args []string) (bool, error) {
	if sub == "rm" || sub == "remove" || sub == "delete" {
		return true, runSJRm(nil, args, ctx)
	}
	if sub == "add-auth" || sub == "auth" {
		return true, runSJAddAuth(nil, args, ctx)
	}
	return false, nil
}

func dispatchSJSubcommand(ctx context.Context, sub string, args []string) error {
	if isHandled, err := dispatchSJBasicSubcommand(ctx, sub, args); isHandled {
		return err
	}
	if isHandled, err := dispatchSJQuerySubcommand(ctx, sub, args); isHandled {
		return err
	}
	if isHandled, err := dispatchSJActionSubcommand(ctx, sub, args); isHandled {
		return err
	}
	return runSJHistory(nil, args, ctx)
}

func routeSJSubcommands(ctx context.Context, args []string) (bool, error) {
	if len(args) == 0 {
		return false, nil
	}
	switch args[0] {
	case "bootstrap", "bs":
		return true, RunClusterBootstrapCLI(args[1:])
	case "install":
		return true, RunClusterInstallCLI(args[1:])
	}
	if !isSJSubcommand(args[0]) {
		return false, nil
	}
	return true, dispatchSJSubcommand(ctx, args[0], args[1:])
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

func persistEnrollment(ctx context.Context, db *sql.DB, opts *SSHJoinOptions) error {
	host, hist := buildHostAndHistory(opts)
	return store.EnrollSSHHost(ctx, host, hist, db)
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
}

func completeEnrollment(ctx context.Context, opts *SSHJoinOptions) error {
	if err := pushAuthIfRequested(ctx, opts); err != nil {
		return err
	}

	printEnrollSuccess(opts.Alias, opts.Target.String())
	return nil
}

func resolveContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}

	return ctx
}

func openAndPersist(ctx context.Context, opts *SSHJoinOptions) error {
	ctx = resolveContext(ctx)
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("enrollParsedTarget", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}

	defer dbConn.Close()

	return persistEnrollment(ctx, dbConn.SQL(), opts)
}

func enrollParsedTarget(ctx context.Context, opts *SSHJoinOptions) error {
	ctx = resolveContext(ctx)
	if opts.Target == nil {
		return apperror.NewValidationError(msgMissingJoinTarget)
	}

	if err := openAndPersist(ctx, opts); err != nil {
		return err
	}

	return completeEnrollment(ctx, opts)
}

func showJoinHelpAndExit() error {
	helptext.Print("ssh-join")
	cliexit.Exit(0)
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

func executeEnrollCLI(ctx context.Context, args []string) error {
	ctx = resolveContext(ctx)
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

func init() {
	SSHJoinCmd.AddCommand(SJAddCmd)
	SSHJoinCmd.AddCommand(SJAddWithPassCmd)
	SSHJoinCmd.AddCommand(SJRmCmd)
	SSHJoinCmd.AddCommand(SJAddAuthCmd)
	SSHJoinCmd.AddCommand(SJLsCmd)
	SSHJoinCmd.AddCommand(SJHistCmd)
	SSHJoinCmd.AddCommand(SJClusterImportCmd)
	SSHJoinCmd.AddCommand(ClusterBootstrapCmd)
	SSHJoinCmd.AddCommand(ClusterInstallCmd)
}
