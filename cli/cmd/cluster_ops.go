package cmd

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runClusterHistory(args []string) error {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	defer storeDB.Close()
	conn := storeDB.Conn()

	if len(args) == 0 {
		printClusterHistoryList(ctx, conn)

		return nil
	}

	printClusterRunDetails(ctx, conn, args[0])

	return nil
}

func printClusterHistoryList(ctx context.Context, conn *sql.DB) {
	runsRes := db.ListClusterRuns(ctx, conn, constants.ClusterDefaultHistoryLimit)
	if runsRes.IsFailure() {
		fmt.Fprintf(os.Stderr, "failed to list runs: %v\n", runsRes.AppError())
		cliexit.HandleError(nil, 1)
	}

	fmt.Printf("%-20s | %-15s | %-15s | %-5s | %-4s | %-4s | %s\n",
		constants.ClusterHeaderRunRef, constants.ClusterHeaderCommandKind,
		constants.ClusterHeaderTargetSelector, constants.ClusterHeaderNodes,
		constants.ClusterHeaderOK, constants.ClusterHeaderFAIL, constants.ClusterHeaderStartedAt)
	for _, r := range runsRes.Data {
		printClusterRunRow(r)
	}
}

func printClusterRunRow(r db.ClusterRun) {
	nodes, ok, fail := countClusterRunNodes(r)
	fmt.Printf("%-20s | %-15s | %-15s | %-5d | %-4d | %-4d | %s\n",
		r.RunRef, r.CommandKind.String(), r.TargetSelector, nodes, ok, fail, r.StartedAt.Format(time.RFC3339))
}

func countClusterRunNodes(r db.ClusterRun) (int, int, int) {
	nodes, ok, fail := 0, 0, 0
	if r.TotalNodes != nil {
		nodes = *r.TotalNodes
	}

	if r.SucceededNodes != nil {
		ok = *r.SucceededNodes
	}

	if r.FailedNodes != nil {
		fail = *r.FailedNodes
	}

	return nodes, ok, fail
}

func printClusterRunDetails(ctx context.Context, conn *sql.DB, runRef string) {
	run, err := db.SelectClusterRun(ctx, conn, runRef)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get run %s: %v\n", runRef, err)
		cliexit.HandleError(nil, 1)
	}

	resultsRes := db.SelectClusterExecResultsByRunId(ctx, conn, run.ClusterRunId)
	if resultsRes.IsFailure() {
		fmt.Fprintf(os.Stderr, "failed to get results for run %s: %v\n", runRef, resultsRes.AppError())
		cliexit.HandleError(nil, 1)
	}

	fmt.Printf("RunRef: %s\nCommand: %s\n\n", run.RunRef, run.RawCommand)
	fmt.Printf("%-20s %-20s %-15s %-8s %s\n",
		constants.ClusterHeaderNode, constants.ClusterHeaderSubCommand,
		constants.ClusterHeaderResult, constants.ClusterHeaderExitCode, constants.ClusterHeaderDurationMs)
	for _, res := range resultsRes.Data {
		printClusterExecResultRow(ctx, conn, res)
	}
}

func printClusterExecResultRow(ctx context.Context, conn *sql.DB, res db.ClusterExecResult) {
	displayStr := formatClusterNodeDisplay(ctx, conn, res.NodeId)
	exitCode, duration := formatClusterExecMetrics(res)
	fmt.Printf("%-20s %-20s %-15s %-8s %s\n", displayStr, res.SubCommand, res.ResultStatus.String(), exitCode, duration)
}

func formatClusterNodeDisplay(ctx context.Context, conn *sql.DB, nodeId string) string {
	node, err := db.GetClusterNode(ctx, conn, nodeId)
	if err == nil {
		return fmt.Sprintf("[%d/%s]", node.DisplayId, node.Alias)
	}

	return nodeId
}

func formatClusterExecMetrics(res db.ClusterExecResult) (string, string) {
	exitCode := "-"
	if res.ExitCode != nil {
		exitCode = fmt.Sprintf("%d", *res.ExitCode)
	}

	duration := "-"
	if res.DurationMs != nil {
		duration = fmt.Sprintf("%d", *res.DurationMs)
	}

	return exitCode, duration
}

func parseClusterExportArgs(args []string) (string, string) {
	format := constants.FormatJSON
	output := ""
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == constants.FlagClusterFormat && i+1 < len(args):
			format = args[i+1]
			i++
		case args[i] == constants.FlagClusterOutput && i+1 < len(args):
			output = args[i+1]
			i++
		}
	}

	return format, output
}

func runClusterExport(args []string) error {
	format, output := parseClusterExportArgs(args)
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	defer storeDB.Close()
	nodesRes := db.ListClusterNodes(ctx, storeDB.Conn())
	if nodesRes.IsFailure() {
		fmt.Fprintf(os.Stderr, "failed to list nodes: %v\n", nodesRes.AppError())
		cliexit.HandleError(nil, 1)
	}

	data := formatClusterExportNodes(nodesRes.Data, format)
	writeClusterExportData(data, output)

	return nil
}

func sanitizeClusterNodes(nodes []db.ClusterNode) {
	for i := range nodes {
		nodes[i].PasswordHash = nil
	}
}

func formatClusterExportNodes(nodes []db.ClusterNode, format string) []byte {
	sanitizeClusterNodes(nodes)
	if format == constants.FormatCSV {
		return exportClusterNodesCSV(nodes)
	}

	data, err := json.MarshalIndent(nodes, "", constants.JSONIndent)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal nodes: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	return data
}

func exportClusterNodesCSV(nodes []db.ClusterNode) []byte {
	var buf strings.Builder
	writer := csv.NewWriter(&buf)
	writer.Write([]string{"NodeId", "Alias", "DisplayId", "IPAddress", "NodeRole", "OS", "JoinedAt", "Status"})
	for _, n := range nodes {
		writer.Write([]string{n.NodeId, n.Alias, strconv.Itoa(n.DisplayId), n.IPAddress, n.NodeRole, n.OS, n.JoinedAt.Format(time.RFC3339), n.Status})
	}

	writer.Flush()

	return []byte(buf.String())
}

func writeClusterExportData(data []byte, output string) {
	if output == "" {
		fmt.Println(string(data))

		return
	}

	if err := os.WriteFile(output, data, constants.FilePermission); err != nil {
		fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
}

func parseClusterNodeID(args []string) string {
	for i := 0; i < len(args); i++ {
		if args[i] == constants.FlagClusterID && i+1 < len(args) {
			return args[i+1]
		}
	}

	return ""
}

func promptClusterPassword() string {
	var pass1, pass2 string
	fmt.Print("Enter password: ")
	fmt.Scanln(&pass1)
	fmt.Print("Confirm password: ")
	fmt.Scanln(&pass2)
	if pass1 != pass2 {
		fmt.Fprintln(os.Stderr, "passwords do not match")
		cliexit.HandleError(nil, 1)
	}

	return pass1
}

func runClusterSetPassword(args []string) error {
	id := parseClusterNodeID(args)
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		cliexit.HandleError(nil, 1)
	}

	pass := promptClusterPassword()
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), constants.ClusterBcryptCost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	hashStr := string(hash)
	updateClusterNodePasswordInDB(id, &hashStr)
	fmt.Println("Password updated successfully.")

	return nil
}

func updateClusterNodePasswordInDB(id string, hash *string) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	defer storeDB.Close()
	if err := db.UpdateClusterNodePassword(ctx, storeDB.Conn(), id, hash); err != nil {
		fmt.Fprintf(os.Stderr, "failed to update password: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
}

func parseClusterConfirmArgs(args []string) (string, bool) {
	id := ""
	hasConfirm := false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == constants.FlagClusterID && i+1 < len(args):
			id = args[i+1]
			i++
		case args[i] == constants.FlagClusterConfirm:
			hasConfirm = true
		}
	}

	return id, hasConfirm
}

func runClusterResetPassword(args []string) error {
	id, hasConfirm := parseClusterConfirmArgs(args)
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		cliexit.HandleError(nil, 1)
	}

	if !hasConfirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
	}

	updateClusterNodePasswordInDB(id, nil)
	fmt.Println("Password reset successfully.")

	return nil
}

var openClusterStore = store.OpenDefault

func hasClusterJSONFlag(args []string) bool {
	for _, a := range args {
		if a == constants.FlagClusterJSON {
			return true
		}
	}

	return false
}

func runClusterNodes(args []string) error {
	if hasHelpFlag(args) {
		helptext.Print("cluster-nodes")
		return nil
	}
	ctx := context.Background()
	storeDB, err := openClusterStore()
	if err != nil {
		return apperror.WrapSimple(err, "runClusterNodes_OpenDB")
	}
	defer storeDB.Close()
	return executeClusterNodesList(ctx, storeDB.Conn(), hasClusterJSONFlag(args))
}

func executeClusterNodesList(ctx context.Context, conn *sql.DB, isJSON bool) error {
	hosts, err := store.ListHosts(ctx, conn)
	if err != nil {
		return apperror.WrapSimple(err, "executeClusterNodesList")
	}
	if isJSON {
		return printClusterHostsJSON(hosts)
	}
	return printClusterHostsTable(hosts)
}

func redactClusterHosts(hosts []store.SSHHost) []store.SSHHost {
	redacted := make([]store.SSHHost, len(hosts))
	copy(redacted, hosts)
	for i := range redacted {
		redacted[i].EncryptedPassword = ""
	}
	return redacted
}

func printClusterHostsJSON(hosts []store.SSHHost) error {
	redacted := redactClusterHosts(hosts)
	data, err := json.MarshalIndent(redacted, "", constants.JSONIndent)
	if err != nil {
		return apperror.WrapSimple(err, "printClusterHostsJSON")
	}
	fmt.Println(string(data))
	return nil
}

const msgNoClusterNodes = "No nodes currently registered in cluster. Enroll with: gitmap cluster add <user@ip|ip> [alias]"

func printClusterHostsTable(hosts []store.SSHHost) error {
	if len(hosts) == 0 {
		fmt.Println(msgNoClusterNodes)
		return nil
	}
	return renderClusterHostsASCII(os.Stdout, hosts)
}

func renderClusterHostsASCII(out io.Writer, hosts []store.SSHHost) error {
	return cmdssh.RenderSSHHostsTable(out, hosts)
}

func hasPositionalNodeTarget(args []string) bool {
	return len(args) > 0 && !strings.HasPrefix(args[0], "-")
}

func removeClusterHostByTarget(target string) error {
	ctx := context.Background()
	storeDB, err := openClusterStore()
	if err != nil {
		return apperror.WrapSimple(err, "removeClusterHostByTarget_OpenDB")
	}
	defer storeDB.Close()
	_, err = store.DeleteHostByAliasOrIP(ctx, target, storeDB.Conn())
	if err != nil {
		return apperror.WrapSimple(err, "removeClusterHostByTarget_Delete")
	}
	fmt.Printf("✓ Node '%s' removed from cluster registry.\n", target)
	return nil
}

func validateClusterRemoveLegacyArgs(args []string) (string, error) {
	id, hasConfirm := parseClusterConfirmArgs(args)
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		cliexit.HandleError(nil, 1)
		return "", apperror.NewSimple("missing --id", "E1001")
	}
	if !hasConfirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
		return "", apperror.NewSimple("missing --confirm", "E1001")
	}
	return id, nil
}

func runClusterRemoveLegacy(args []string) error {
	id, err := validateClusterRemoveLegacyArgs(args)
	if err != nil {
		return err
	}
	deleteClusterNodeInDB(id)
	fmt.Println("Node deleted successfully.")
	return nil
}

func runClusterRemove(args []string) error {
	if hasHelpFlag(args) {
		helptext.Print("cluster-remove")
		return nil
	}
	if hasPositionalNodeTarget(args) {
		return removeClusterHostByTarget(args[0])
	}
	return runClusterRemoveLegacy(args)
}

func deleteClusterNodeInDB(id string) {
	ctx := context.Background()
	storeDB, err := openClusterStore()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
		return
	}

	defer storeDB.Close()
	if err := db.DeleteClusterNode(ctx, storeDB.Conn(), id); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete node: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
}

func parseClusterAuditCleanArgs(args []string) (string, bool) {
	beforeStr := ""
	hasConfirm := false
	for i := 0; i < len(args); i++ {
		switch {
		case args[i] == constants.FlagClusterBefore && i+1 < len(args):
			beforeStr = args[i+1]
			i++
		case args[i] == constants.FlagClusterConfirm:
			hasConfirm = true
		}
	}

	return beforeStr, hasConfirm
}

func validateClusterAuditCleanArgs(beforeStr string, hasConfirm bool) {
	if beforeStr == "" {
		fmt.Fprintln(os.Stderr, "missing --before")
		cliexit.HandleError(nil, 1)
	}

	if !hasConfirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
	}
}

func runClusterAuditClean(args []string) error {
	beforeStr, hasConfirm := parseClusterAuditCleanArgs(args)
	validateClusterAuditCleanArgs(beforeStr, hasConfirm)
	cleanClusterAuditRecords(beforeStr)

	return nil
}

func cleanClusterAuditRecords(beforeStr string) {
	before, err := time.Parse(time.RFC3339, beforeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date format, use RFC3339: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	defer storeDB.Close()
	count, delErr := db.DeleteClusterRunsBefore(ctx, storeDB.Conn(), before)
	if delErr != nil {
		fmt.Fprintf(os.Stderr, "failed to clean audit records: %v\n", delErr)
		cliexit.HandleError(nil, 1)
	}

	fmt.Printf("Cleaned %d cluster run records older than %s.\n", count, beforeStr)
}

func runClusterStats(args []string) error {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	defer storeDB.Close()
	stats, statsErr := db.GetClusterStats(ctx, storeDB.Conn())
	if statsErr != nil {
		fmt.Fprintf(os.Stderr, "failed to get cluster stats: %v\n", statsErr)
		cliexit.HandleError(nil, 1)
	}

	printClusterStatsReport(stats)

	return nil
}

func printClusterStatsReport(stats db.ClusterStats) {
	rate := 0.0
	if stats.TotalCommands > 0 {
		rate = float64(stats.SuccessCommands) / float64(stats.TotalCommands) * 100
	}

	fmt.Println("Cluster Statistics:")
	fmt.Printf("Total Runs: %d\n", stats.TotalRuns)
	fmt.Printf("Total Commands Dispatched: %d\n", stats.TotalCommands)
	fmt.Printf("Success Rate: %.2f%%\n", rate)
	fmt.Printf("Most Targeted Node: %s\n", stats.MostTargetedNode)
	fmt.Printf("Most Used Sub-Command: %s\n", stats.MostUsedSubCmd)
}
