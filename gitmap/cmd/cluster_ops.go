package cmd

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func runClusterHistory(args []string) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()
	conn := storeDB.Conn()

	if len(args) == 0 {
		runs, err := db.ListClusterRuns(ctx, conn, 50)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to list runs: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("%-20s | %-15s | %-15s | %-5s | %-4s | %-4s | %s\n", "RunRef", "CommandKind", "TargetSelector", "Nodes", "OK", "FAIL", "StartedAt")
		for _, r := range runs {
			nodes := 0
			if r.TotalNodes != nil {
				nodes = *r.TotalNodes
			}
			ok := 0
			if r.SucceededNodes != nil {
				ok = *r.SucceededNodes
			}
			fail := 0
			if r.FailedNodes != nil {
				fail = *r.FailedNodes
			}
			fmt.Printf("%-20s | %-15s | %-15s | %-5d | %-4d | %-4d | %s\n", r.RunRef, r.CommandKind.String(), r.TargetSelector, nodes, ok, fail, r.StartedAt.Format(time.RFC3339))
		}
	} else {
		runRef := args[0]
		run, err := db.SelectClusterRun(ctx, conn, runRef)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get run %s: %v\n", runRef, err)
			os.Exit(1)
		}

		results, err := db.SelectClusterExecResultsByRunId(ctx, conn, run.ClusterRunId)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get results for run %s: %v\n", runRef, err)
			os.Exit(1)
		}

		fmt.Printf("RunRef: %s\nCommand: %s\n\n", run.RunRef, run.RawCommand)
		fmt.Printf("%-20s %-20s %-15s %-8s %s\n", "Node", "SubCommand", "Result", "ExitCode", "DurationMs")
		for _, res := range results {
			nodeId := res.NodeId
			node, err := db.GetClusterNode(ctx, conn, res.NodeId)
			displayStr := nodeId
			if err == nil {
				displayStr = fmt.Sprintf("[%d/%s]", node.DisplayId, node.Alias)
			}

			exitCode := "-"
			if res.ExitCode != nil {
				exitCode = fmt.Sprintf("%d", *res.ExitCode)
			}
			duration := "-"
			if res.DurationMs != nil {
				duration = fmt.Sprintf("%d", *res.DurationMs)
			}
			fmt.Printf("%-20s %-20s %-15s %-8s %s\n", displayStr, res.SubCommand, res.ResultStatus.String(), exitCode, duration)
		}
	}
}

func runClusterExport(args []string) {
	format := "json"
	output := ""

	for i := 0; i < len(args); i++ {
		if args[i] == "--format" && i+1 < len(args) {
			format = args[i+1]
			i++
		} else if args[i] == "--output" && i+1 < len(args) {
			output = args[i+1]
			i++
		}
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()
	conn := storeDB.Conn()

	nodes, err := db.ListClusterNodes(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list nodes: %v\n", err)
		os.Exit(1)
	}

	for i := range nodes {
		nodes[i].PasswordHash = nil
	}

	var data []byte
	if format == "csv" {
		var buf strings.Builder
		writer := csv.NewWriter(&buf)
		writer.Write([]string{"NodeId", "Alias", "DisplayId", "IPAddress", "NodeRole", "OS", "JoinedAt", "Status"})
		for _, n := range nodes {
			writer.Write([]string{n.NodeId, n.Alias, strconv.Itoa(n.DisplayId), n.IPAddress, n.NodeRole, n.OS, n.JoinedAt.Format(time.RFC3339), n.Status})
		}
		writer.Flush()
		data = []byte(buf.String())
	} else {
		data, err = json.MarshalIndent(nodes, "", "  ")
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to marshal nodes: %v\n", err)
			os.Exit(1)
		}
	}

	if output != "" {
		err = os.WriteFile(output, data, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to write output: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Println(string(data))
	}
}

func runClusterImport(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "missing input file")
		os.Exit(1)
	}
	file := args[0]
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read file: %v\n", err)
		os.Exit(1)
	}

	var nodes []db.ClusterNode
	err = json.Unmarshal(data, &nodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse json: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()
	conn := storeDB.Conn()

	inserted := 0
	updated := 0
	skipped := 0

	existing, _ := db.ListClusterNodes(ctx, conn)
	existingMap := make(map[string]bool)
	for _, e := range existing {
		existingMap[e.NodeId] = true
	}

	for _, n := range nodes {
		isUpdate := existingMap[n.NodeId]
		err := db.InsertOrUpdateClusterNode(ctx, conn, n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to import node %s: %v\n", n.NodeId, err)
			skipped++
		} else {
			if isUpdate {
				updated++
			} else {
				inserted++
			}
		}
	}

	fmt.Printf("Inserted: %d, Updated: %d, Skipped: %d\n", inserted, updated, skipped)
}

func runClusterSetPassword(args []string) {
	id := ""
	for i := 0; i < len(args); i++ {
		if args[i] == "--id" && i+1 < len(args) {
			id = args[i+1]
			i++
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		os.Exit(1)
	}

	var pass1, pass2 string
	fmt.Print("Enter password: ")
	fmt.Scanln(&pass1)
	fmt.Print("Confirm password: ")
	fmt.Scanln(&pass2)

	if pass1 != pass2 {
		fmt.Fprintln(os.Stderr, "passwords do not match")
		os.Exit(1)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pass1), 12)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to hash password: %v\n", err)
		os.Exit(1)
	}
	hashStr := string(hash)

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()

	err = db.UpdateClusterNodePassword(ctx, storeDB.Conn(), id, &hashStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to update password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Password updated successfully.")
}

func runClusterResetPassword(args []string) {
	id := ""
	confirm := false
	for i := 0; i < len(args); i++ {
		if args[i] == "--id" && i+1 < len(args) {
			id = args[i+1]
			i++
		} else if args[i] == "--confirm" {
			confirm = true
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		os.Exit(1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		os.Exit(1)
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()

	err = db.UpdateClusterNodePassword(ctx, storeDB.Conn(), id, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to reset password: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Password reset successfully.")
}

func runClusterNodes(args []string) {
	asJson := false
	for _, a := range args {
		if a == "--json" {
			asJson = true
		}
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()
	conn := storeDB.Conn()

	nodes, err := db.ListClusterNodes(ctx, conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to list nodes: %v\n", err)
		os.Exit(1)
	}

	if asJson {
		data, _ := json.MarshalIndent(nodes, "", "  ")
		fmt.Println(string(data))
		return
	}

	fmt.Printf("%-10s | %-15s | %-15s | %-10s | %-10s | %-10s | %s\n", "DisplayId", "Alias", "IP", "OS", "Role", "Status", "LastHeartbeat")
	for _, n := range nodes {
		hb := "-"
		if n.LastHeartbeat != nil {
			hb = n.LastHeartbeat.Format(time.RFC3339)
		}
		fmt.Printf("%-10d | %-15s | %-15s | %-10s | %-10s | %-10s | %s\n", n.DisplayId, n.Alias, n.IPAddress, n.OS, n.NodeRole, n.Status, hb)
	}
}

func runClusterRemove(args []string) {
	id := ""
	confirm := false
	for i := 0; i < len(args); i++ {
		if args[i] == "--id" && i+1 < len(args) {
			id = args[i+1]
			i++
		} else if args[i] == "--confirm" {
			confirm = true
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "missing --id")
		os.Exit(1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		os.Exit(1)
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()

	err = db.DeleteClusterNode(ctx, storeDB.Conn(), id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to delete node: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Node deleted successfully.")
}

func runClusterAuditClean(args []string) {
	beforeStr := ""
	confirm := false
	for i := 0; i < len(args); i++ {
		if args[i] == "--before" && i+1 < len(args) {
			beforeStr = args[i+1]
			i++
		} else if args[i] == "--confirm" {
			confirm = true
		}
	}
	if beforeStr == "" {
		fmt.Fprintln(os.Stderr, "missing --before")
		os.Exit(1)
	}
	if !confirm {
		fmt.Fprintln(os.Stderr, "missing --confirm")
		os.Exit(1)
	}

	before, err := time.Parse(time.RFC3339, beforeStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid date format, use RFC3339: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()

	count, err := db.DeleteClusterRunsBefore(ctx, storeDB.Conn(), before)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to clean audit records: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Cleaned %d cluster run records older than %s.\n", count, beforeStr)
}

func runClusterStats(args []string) {
	ctx := context.Background()
	storeDB, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to open db: %v\n", err)
		os.Exit(1)
	}
	defer storeDB.Close()

	stats, err := db.GetClusterStats(ctx, storeDB.Conn())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to get cluster stats: %v\n", err)
		os.Exit(1)
	}

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
