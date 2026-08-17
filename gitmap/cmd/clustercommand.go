package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

// runClusterCommand is the orchestrator for delegated cluster commands.
// It parses flags, resolves nodes, prints preflight, confirms,
// inserts a ClusterRun, and dispatches the sub-commands via the bounded worker pool.
func runClusterCommand(selector cluster.TargetSelectorType, args []string) {
	flags, positional, err := ParseClusterFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	subCmds, err := ParseSubCommands(positional)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing sub-commands: %v\n", err)
		os.Exit(1)
	}

	hasNoSubCmds := len(subCmds) == 0
	if hasNoSubCmds == true {
		fmt.Fprintln(os.Stderr, "No sub-commands provided.")
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var dbConn *sql.DB // Stub DB for now

	// Stub node retrieval
	allNodes := []cluster.ClusterNode{
		{ID: "node-1", DisplayId: 1, IP: "192.168.1.10", IsServer: true},
		{ID: "node-2", DisplayId: 2, IP: "192.168.1.11", IsServer: false},
	}

	filter := cluster.NodeFilter{
		Except: []string{},
		IPs:    flags.OnlyIPs,
		IDs:    flags.OnlyIDs,
	}

	hasExceptClause := flags.ExceptClause != ""
	if hasExceptClause == true {
		filter.Except = strings.Split(flags.ExceptClause, ",")
	}

	effective, err := cluster.ResolveTargetNodes(selector, filter, allNodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving target nodes: %v\n", err)
		os.Exit(1)
	}

	hasNoEffectiveNodes := len(effective) == 0
	if hasNoEffectiveNodes == true {
		fmt.Fprintln(os.Stderr, constants.ErrClusterNoNodes)
		os.Exit(1)
	}

	// Stub run ref generator if DB is nil
	runRef := "RUN-YYYYMMDD-001"
	if dbConn != nil {
		runRef, err = cluster.RunRefGenerator(dbConn)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error generating run ref: %v\n", err)
			os.Exit(1)
		}
	}

	cmdStr := strings.Join(positional, " ")

	if flags.NoPreflight == false {
		confirmed, err := cluster.PrintPreflight(selector, effective, cmdStr, runRef, flags.AutoConfirm)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Preflight error: %v\n", err)
			os.Exit(1)
		}

		isConfirmed := confirmed == true
		if isConfirmed == false {
			fmt.Fprintln(os.Stderr, "Operation aborted.")
			os.Exit(1)
		}
	}

	var runId int64
	totalNodes := len(effective)

	selectorStr := "ServersClients"
	if selector == cluster.ClientsOnly {
		selectorStr = "ClientsOnly"
	}

	run := db.ClusterRun{
		RunRef:         runRef,
		CommandKind:    subCmds[0].Kind,
		RawCommand:     cmdStr,
		TargetSelector: selectorStr,
		StartedAt:      time.Now(),
		TotalNodes:     &totalNodes,
	}

	if hasExceptClause == true {
		run.ExceptClause = &flags.ExceptClause
	}

	if dbConn != nil {
		runId, err = db.InsertClusterRun(ctx, dbConn, run)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting ClusterRun: %v\n", err)
			os.Exit(1)
		}
	}

	resultCh := make(chan db.ClusterExecResult, len(effective)*len(subCmds))

	verbose := false // or maybe parse it from flags later if needed
	cluster.RunPool(ctx, effective, subCmds, dbConn, runId, 10, resultCh, verbose)

	succeeded := 0
	failed := 0
	skipped := 0
	for res := range resultCh {
		switch res.ResultStatus {
		case db.ResultStatusSucceeded:
			succeeded++
		case db.ResultStatusFailed:
			failed++
		case db.ResultStatusSkipped:
			skipped++
		}

		statusStr := res.ResultStatus.String()
		fmt.Printf("[%s] %s: %s\n", res.NodeId, res.SubCommand, statusStr)
	}

	if dbConn != nil {
		now := time.Now()
		err = db.UpdateClusterRun(ctx, dbConn, runId, &now, &totalNodes, &succeeded, &failed, &skipped)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error updating ClusterRun counts: %v\n", err)
		}
	}

	summaryTitle := fmt.Sprintf(" Cluster Run %s ", runRef)
	summaryData := fmt.Sprintf(" Nodes: %d  OK: %d  Failed: %d  Skipped: %d ", totalNodes, succeeded, failed, skipped)

	boxWidth := len(summaryTitle)
	if len(summaryData) > boxWidth {
		boxWidth = len(summaryData)
	}

	titlePad := boxWidth - len(summaryTitle)
	dataPad := boxWidth - len(summaryData)

	fmt.Printf("┌%s%s┐\n", summaryTitle, strings.Repeat("─", titlePad))
	fmt.Printf("│%s%s│\n", summaryData, strings.Repeat(" ", dataPad))
	fmt.Printf("└%s┘\n", strings.Repeat("─", boxWidth))
}
