package cmd

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
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
	runs, err := db.ListClusterRuns(ctx, conn, constants.ClusterDefaultHistoryLimit)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list runs: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	fmt.Printf("%-20s | %-15s | %-15s | %-5s | %-4s | %-4s | %s\n",
		constants.ClusterHeaderRunRef, constants.ClusterHeaderCommandKind,
		constants.ClusterHeaderTargetSelector, constants.ClusterHeaderNodes,
		constants.ClusterHeaderOK, constants.ClusterHeaderFAIL, constants.ClusterHeaderStartedAt)
	for _, r := range runs {
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
	results, err := db.SelectClusterExecResultsByRunId(ctx, conn, run.ClusterRunId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get results for run %s: %v\n", runRef, err)
		cliexit.HandleError(nil, 1)
	}
	fmt.Printf("RunRef: %s\nCommand: %s\n\n", run.RunRef, run.RawCommand)
	fmt.Printf("%-20s %-20s %-15s %-8s %s\n",
		constants.ClusterHeaderNode, constants.ClusterHeaderSubCommand,
		constants.ClusterHeaderResult, constants.ClusterHeaderExitCode, constants.ClusterHeaderDurationMs)
	for _, res := range results {
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
		if args[i] == constants.FlagClusterFormat && i+1 < len(args) {
			format = args[i+1]
			i++
		} else if args[i] == constants.FlagClusterOutput && i+1 < len(args) {
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
	nodes, err := db.ListClusterNodes(ctx, storeDB.Conn())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list nodes: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	data := formatClusterExportNodes(nodes, format)
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

func runClusterImport(args []string) error {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "missing input file")
		cliexit.HandleError(nil, 1)
	}
	nodes := loadClusterImportNodes(args[0])
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer storeDB.Close()
	importClusterNodes(ctx, storeDB.Conn(), nodes)
	return nil
}

func loadClusterImportNodes(file string) []db.ClusterNode {
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	var nodes []db.ClusterNode
	if err := json.Unmarshal(data, &nodes); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse json: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	return nodes
}

func importClusterNodes(ctx context.Context, conn *sql.DB, nodes []db.ClusterNode) {
	existingMap := getExistingClusterNodesMap(ctx, conn)
	inserted, updated, skipped := 0, 0, 0
	for _, n := range nodes {
		isUpdate := existingMap[n.NodeId]
		if err := db.InsertOrUpdateClusterNode(ctx, conn, n); err != nil {
			fmt.Fprintf(os.Stderr, "failed to import node %s: %v\n", n.NodeId, err)
			skipped++
		} else if isUpdate {
			updated++
		} else {
			inserted++
		}
	}
	fmt.Printf("Inserted: %d, Updated: %d, Skipped: %d\n", inserted, updated, skipped)
}

func getExistingClusterNodesMap(ctx context.Context, conn *sql.DB) map[string]bool {
	existing, _ := db.ListClusterNodes(ctx, conn)
	m := make(map[string]bool, len(existing))
	for _, e := range existing {
		m[e.NodeId] = true
	}
	return m
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
	confirm := false
	for i := 0; i < len(args); i++ {
		if args[i] == constants.FlagClusterID && i+1 < len(args) {
			id = args[i+1]
			i++
		} else if args[i] == constants.FlagClusterConfirm {
			confirm = true
		}
	}
	return id, confirm
}

func runClusterResetPassword(args []string) error {
	id, confirm := parseClusterConfirmArgs(args)
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		cliexit.HandleError(nil, 1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
	}
	updateClusterNodePasswordInDB(id, nil)
	fmt.Println("Password reset successfully.")
	return nil
}

func hasClusterJSONFlag(args []string) bool {
	for _, a := range args {
		if a == constants.FlagClusterJSON {
			return true
		}
	}
	return false
}

func runClusterNodes(args []string) error {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer storeDB.Close()
	nodes, err := db.ListClusterNodes(ctx, storeDB.Conn())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list nodes: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	displayClusterNodes(nodes, hasClusterJSONFlag(args))
	return nil
}

func displayClusterNodes(nodes []db.ClusterNode, asJson bool) {
	if asJson {
		data, _ := json.MarshalIndent(nodes, "", constants.JSONIndent)
		fmt.Println(string(data))
		return
	}
	printClusterNodesTable(nodes)
}

func printClusterNodesTable(nodes []db.ClusterNode) {
	fmt.Printf("%-10s | %-15s | %-15s | %-10s | %-10s | %-10s | %s\n", "DisplayId", "Alias", "IP", "OS", "Role", "Status", "LastHeartbeat")
	var unreachable []db.ClusterNode
	for _, n := range nodes {
		hb := "-"
		if n.LastHeartbeat != nil {
			hb = n.LastHeartbeat.Format(time.RFC3339)
		}
		if strings.EqualFold(n.Status, constants.ClusterStatusOffline) || strings.EqualFold(n.Status, constants.ClusterStatusUnreachable) {
			unreachable = append(unreachable, n)
		}
		fmt.Printf("%-10d | %-15s | %-15s | %-10s | %-10s | %-10s | %s\n", n.DisplayId, n.Alias, n.IPAddress, n.OS, n.NodeRole, n.Status, hb)
	}
	warnUnreachableNodes(unreachable)
}

func warnUnreachableNodes(unreachable []db.ClusterNode) {
	if len(unreachable) == 0 {
		return
	}
	fmt.Println()
	fmt.Printf("  ▲ WARNING: %d cluster machine(s) offline or unreachable:\n", len(unreachable))
	for _, u := range unreachable {
		fmt.Printf("     • Node %d [%s] (%s) - status: %s\n", u.DisplayId, u.Alias, u.IPAddress, u.Status)
	}
	fmt.Println("  These are the machines that cannot connect or find.")
	fmt.Println()
}

func runClusterRemove(args []string) error {
	id, confirm := parseClusterConfirmArgs(args)
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		cliexit.HandleError(nil, 1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
	}
	deleteClusterNodeInDB(id)
	fmt.Println("Node deleted successfully.")
	return nil
}

func deleteClusterNodeInDB(id string) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	defer storeDB.Close()
	if err := db.DeleteClusterNode(ctx, storeDB.Conn(), id); err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete node: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
}

func parseClusterAuditCleanArgs(args []string) (string, bool) {
	beforeStr := ""
	confirm := false
	for i := 0; i < len(args); i++ {
		if args[i] == constants.FlagClusterBefore && i+1 < len(args) {
			beforeStr = args[i+1]
			i++
		} else if args[i] == constants.FlagClusterConfirm {
			confirm = true
		}
	}
	return beforeStr, confirm
}

func runClusterAuditClean(args []string) error {
	beforeStr, confirm := parseClusterAuditCleanArgs(args)
	if beforeStr == "" {
		fmt.Fprintln(os.Stderr, "missing --before")
		cliexit.HandleError(nil, 1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		cliexit.HandleError(nil, 1)
	}
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
	count, err := db.DeleteClusterRunsBefore(ctx, storeDB.Conn(), before)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to clean audit records: %v\n", err)
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
	stats, err := db.GetClusterStats(ctx, storeDB.Conn())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get cluster stats: %v\n", err)
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
