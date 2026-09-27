package cmdagy

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var (
	agyLogJSON    bool
	agyLogSSH     bool
	agyLogLimit   int
	agyLogProject string
	agyLogCommand string
)

// AgyLogCmd provides inspection of Antigravity decision audit logs.
var AgyLogCmd = &cobra.Command{
	Use:     "log",
	Aliases: []string{"logs", "decision-log", "decision-logs"},
	Short:   "View Antigravity decision audit logs and operational traces",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunAgyLogCLI(args)
	},
}

func init() {
	AgyLogCmd.Flags().BoolVar(&agyLogJSON, "json", false, "Output logs in JSON format")
	AgyLogCmd.Flags().BoolVar(&agyLogSSH, "ssh", false, "Query decision logs across remote SSH cluster nodes")
	AgyLogCmd.Flags().IntVarP(&agyLogLimit, "limit", "n", 20, "Maximum number of log entries to display")
	AgyLogCmd.Flags().StringVarP(&agyLogProject, "project", "P", "", "Filter decision logs by target project")
	AgyLogCmd.Flags().StringVarP(&agyLogCommand, "command", "c", "", "Filter decision logs by command name")
	AgyCmd.AddCommand(AgyLogCmd)
}

// RunAgyLogCLI executes the decision log query and rendering.
func RunAgyLogCLI(args []string) error {
	if checkAgyLogHelp(args) {
		printAgyLogHelp()
		return nil
	}
	opts := store.AgyLogQueryOptions{
		Project: agyLogProject,
		Command: agyLogCommand,
		Limit:   agyLogLimit,
	}
	logs, err := queryLocalAgyLogs(opts)
	if err != nil {
		return err
	}
	if agyLogSSH {
		remoteLogs := fetchClusterSSHAgyLogs(opts)
		logs = mergeAgyLogs(logs, remoteLogs, opts.Limit)
	}
	if agyLogJSON {
		return outputAgyLogsJSON(logs)
	}
	RenderAgyLogsTable(logs)
	return nil
}

func checkAgyLogHelp(args []string) bool {
	if len(args) == 0 {
		return false
	}
	sub := strings.ToLower(args[0])
	return sub == "help" || sub == "--help" || sub == "-h"
}

func queryLocalAgyLogs(opts store.AgyLogQueryOptions) ([]store.AgyDecisionLogRecord, error) {
	logDB, err := store.OpenAgyLogSplitDB("")
	if err != nil {
		return nil, err
	}
	defer logDB.Close()
	return logDB.QueryAgyDecisionLogs(opts)
}

func mergeAgyLogs(local, remote []store.AgyDecisionLogRecord, limit int) []store.AgyDecisionLogRecord {
	seen := make(map[string]bool)
	var merged []store.AgyDecisionLogRecord
	for _, l := range local {
		if !seen[l.LogId] {
			seen[l.LogId] = true
			merged = append(merged, l)
		}
	}
	for _, r := range remote {
		if !seen[r.LogId] {
			seen[r.LogId] = true
			merged = append(merged, r)
		}
	}
	sort.SliceStable(merged, func(i, j int) bool {
		return merged[i].CreatedAt > merged[j].CreatedAt
	})
	if limit > 0 && len(merged) > limit {
		return merged[:limit]
	}
	return merged
}

func fetchClusterSSHAgyLogs(opts store.AgyLogQueryOptions) []store.AgyDecisionLogRecord {
	if SSHConnectionsFetcher == nil {
		return nil
	}
	conns, err := SSHConnectionsFetcher()
	if err != nil || len(conns) == 0 {
		return nil
	}
	var aggregated []store.AgyDecisionLogRecord
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			records := querySingleNodeAgyLogs(conn, opts)
			mu.Lock()
			aggregated = append(aggregated, records...)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return aggregated
}

func querySingleNodeAgyLogs(c db.SSHConnection, opts store.AgyLogQueryOptions) []store.AgyDecisionLogRecord {
	if !probeNodeOnline(c.IPAddress, 400*time.Millisecond) {
		return nil
	}
	client, err := dialSSHNodeClient(c)
	if err != nil {
		return nil
	}
	defer client.Close()
	cmdStr := fmt.Sprintf("gitmap agy log --json -n %d", opts.Limit)
	if opts.Project != "" {
		cmdStr += fmt.Sprintf(" -P %s", opts.Project)
	}
	if opts.Command != "" {
		cmdStr += fmt.Sprintf(" -c %s", opts.Command)
	}
	out, runErr := crypto.RunCommand(client, cmdStr, resolveNodeShell(c.OS))
	if runErr != nil {
		return nil
	}
	var records []store.AgyDecisionLogRecord
	if jsonErr := json.Unmarshal([]byte(strings.TrimSpace(out)), &records); jsonErr != nil {
		return nil
	}
	assignRemoteNodeAliases(records, c.Alias)
	return records
}

func assignRemoteNodeAliases(records []store.AgyDecisionLogRecord, alias string) {
	for i := range records {
		if records[i].Node == "" || records[i].Node == "local" {
			records[i].Node = alias
		}
	}
}

func outputAgyLogsJSON(logs []store.AgyDecisionLogRecord) error {
	data, err := json.MarshalIndent(logs, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal agy logs JSON")
	}
	fmt.Println(string(data))
	return nil
}

func printAgyLogHelp() {
	fmt.Println()
	fmt.Println("  Antigravity Decision Audit Logs")
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap agy log [flags]")
	fmt.Println("    gitmap agy logs [flags]")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    --json          Output in structured JSON format")
	fmt.Println("    --ssh           Query and merge decision logs from cluster SSH nodes")
	fmt.Println("    -n, --limit     Maximum number of entries to display (default: 20)")
	fmt.Println("    -P, --project   Filter by target project name or path")
	fmt.Println("    -c, --command   Filter by AGY command (rerun, ls, inject, etc.)")
	fmt.Println()
}
