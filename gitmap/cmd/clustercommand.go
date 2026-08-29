package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cluster"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runClusterCommand is the orchestrator for delegated cluster commands.
// It parses flags, resolves nodes, prints preflight, confirms,
// inserts a ClusterRun, and dispatches the sub-commands via the bounded worker pool.
func runClusterCommand(selector cluster.TargetSelectorType, args []string) error {
	flags, positional, err := ParseClusterFlags(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	subCmds, err := ParseSubCommands(positional)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing sub-commands: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	isEmptySubCmds := len(subCmds) == 0
	if isEmptySubCmds {
		fmt.Fprintln(os.Stderr, "No sub-commands provided.")
		cliexit.HandleError(nil, 1)
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
	if hasExceptClause {
		filter.Except = strings.Split(flags.ExceptClause, ",")
	}

	effective, err := cluster.ResolveTargetNodes(selector, filter, allNodes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error resolving target nodes: %v\n", err)
		cliexit.HandleError(nil, 1)
	}

	isEmptyEffectiveNodes := len(effective) == 0
	if isEmptyEffectiveNodes {
		fmt.Fprintln(os.Stderr, constants.ErrClusterNoNodes)
		cliexit.HandleError(nil, 1)
	}

	// Stub run ref generator if DB is nil
	runRef := generateRunRef(dbConn)

	cmdStr := strings.Join(positional, " ")

	if err := performPreflight(flags, selector, effective, cmdStr, runRef); err != nil {
		return err
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

	if hasExceptClause {
		run.ExceptClause = &flags.ExceptClause
	}

	runId = insertRun(ctx, dbConn, run)

	resultCh := make(chan db.ClusterExecResult, len(effective)*len(subCmds))

	if flags.DryRun {
		fmt.Fprintln(os.Stdout, "Dry run enabled. The following commands would be executed:")
		for _, sc := range subCmds {
			fmt.Printf("  %s\n", sc.Kind.String())
		}
		cliexit.HandleError(nil, 0)
	}

	verbose := flags.Verbose
	cluster.RunPool(ctx, effective, subCmds, dbConn, runId, 10, resultCh, verbose)

	succeeded := 0
	failed := 0
	skipped := 0
	for res := range resultCh {
		switch res.ResultStatus {
		case db.ResultStatusSucceeded:
			succeeded++
		case db.ResultStatusFailed, db.ResultStatusRequiresAuth:
			failed++
		case db.ResultStatusSkipped, db.ResultStatusPending, db.ResultStatusDeferred:
			skipped++
		}

		statusStr := res.ResultStatus.String()
		fmt.Printf("[%s] %s: %s\n", res.NodeId, res.SubCommand, statusStr)
	}

	updateRunCounts(ctx, dbConn, runId, totalNodes, succeeded, failed, skipped)

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
	return nil
}

func generateRunRef(dbConn *sql.DB) string {
	if dbConn == nil {
		return "RUN-YYYYMMDD-001"
	}
	runRef, err := cluster.RunRefGenerator(dbConn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating run ref: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	return runRef
}

func performPreflight(
	flags ClusterFlags,
	selector cluster.TargetSelectorType,
	effective []cluster.ClusterNode,
	cmdStr string,
	runRef string,
) *apperror.AppError {
	if flags.NoPreflight {
		return nil
	}
	confirmed, err := cluster.PrintPreflight(selector, effective, cmdStr, runRef, flags.AutoConfirm)
	if err != nil {
		return apperror.WrapSimple(err, "Preflight error")
	}

	isConfirmed := confirmed
	if !isConfirmed {
		return apperror.NewSimple("Operation aborted", "E9000")
	}
	return nil
}

func insertRun(ctx context.Context, dbConn *sql.DB, run db.ClusterRun) int64 {
	if dbConn == nil {
		return 0
	}
	runId, err := db.InsertClusterRun(ctx, dbConn, run)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inserting ClusterRun: %v\n", err)
		cliexit.HandleError(nil, 1)
	}
	return runId
}

func updateRunCounts(
	ctx context.Context,
	dbConn *sql.DB,
	runId int64,
	totalNodes,
	succeeded,
	failed,
	skipped int,
) {
	if dbConn == nil {
		return
	}
	now := time.Now()
	err := db.UpdateClusterRun(ctx, dbConn, runId, &now, &totalNodes, &succeeded, &failed, &skipped)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error updating ClusterRun counts: %v\n", err)
	}
}
