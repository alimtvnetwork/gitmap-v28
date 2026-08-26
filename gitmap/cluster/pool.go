package cluster

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"

	"github.com/pterm/pterm"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
)

// RunPool executes sub-commands across multiple nodes using a bounded worker pool.
// Each worker processes one node's sub-command chain sequentially.
// It also tracks progress using a parallel-spinner UI and handles Ctrl+C.
func RunPool(ctx context.Context, nodes []ClusterNode, subCmds []ClusterSubCommand, dbConn *sql.DB, runId int64, maxWorkers int, resultCh chan<- db.ClusterExecResult, verbose bool) {
	isInvalidWorkers := maxWorkers <= 0
	if isInvalidWorkers == true {
		maxWorkers = 10
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	multi := pterm.DefaultMultiPrinter
	isMultiActive := pterm.Output
	if isMultiActive {
		multi.Start()
	}

	var mu sync.Mutex
	spinners := make(map[string]*pterm.SpinnerPrinter)

	workCh := make(chan ClusterNode, len(nodes))
	for _, node := range nodes {
		workCh <- node
	}
	close(workCh)

	var wg sync.WaitGroup
	totalNodes := len(nodes)
	var succeeded, failed, skipped int

	updateCounts := func(isCancelled bool) {
		mu.Lock()
		defer mu.Unlock()
		s := succeeded
		f := failed
		sk := skipped
		if isCancelled {
			sk = totalNodes - (s + f)
		}
		if dbConn != nil {
			now := time.Now()
			_ = db.UpdateClusterRun(ctx, dbConn, runId, &now, &totalNodes, &s, &f, &sk)
		}
	}

	go func() {
		select {
		case <-sigCh:
			cancel()
			updateCounts(true)
			if isMultiActive {
				multi.Stop()
			}

			fmt.Println("\n┌ Cluster Run Cancelled ─────────────────┐")
			mu.Lock()
			fmt.Printf("│ Nodes: %d  OK: %d  Failed: %d  Skipped: %d │\n", totalNodes, succeeded, failed, totalNodes-(succeeded+failed))
			mu.Unlock()
			fmt.Println("└────────────────────────────────────────┘")
		case <-ctx.Done():
		}
	}()

	worker := func() {
		defer wg.Done()
		for node := range workCh {
			if ctx.Err() != nil {
				mu.Lock()
				skipped++
				mu.Unlock()
				continue
			}

			nodeLabel := fmt.Sprintf("[%d/%s]", node.DisplayId, node.ID)

			cmdParts := []string{}
			for _, sc := range subCmds {
				cmdParts = append(cmdParts, sc.Kind.String())
			}
			displayCmd := strings.Join(cmdParts, " ")
			if len(subCmds) == 1 {
				displayCmd = subCmds[0].RawArg
			}
			if len(subCmds) == 1 && displayCmd == "" {
				displayCmd = subCmds[0].Kind.String()
			}

			var spinner *pterm.SpinnerPrinter
			if isMultiActive {
				mu.Lock()
				spinner, _ = pterm.DefaultSpinner.WithWriter(multi.NewWriter()).Start(fmt.Sprintf("%s Running %s\u2026", nodeLabel, displayCmd))
				spinners[node.ID] = spinner
				mu.Unlock()
			} else {
				pterm.Info.Printf("%s Running %s\u2026\n", nodeLabel, displayCmd)
			}

			var lastRes db.ClusterExecResult
			allOk := true

			for _, subCmd := range subCmds {
				if ctx.Err() != nil {
					allOk = false
					break
				}

				res := Dispatch(ctx, node, subCmd)
				res.ClusterRunId = runId
				lastRes = res

				id, dbErr := int64(0), error(nil)
				if dbConn != nil {
					id, dbErr = db.InsertClusterExecResult(ctx, dbConn, res)
				}
				if dbConn != nil && dbErr == nil {
					res.ClusterExecResultId = id
				}

				lines := []string{}
				if verbose && res.Stdout != nil && *res.Stdout != "" {
					lines = strings.Split(*res.Stdout, "\n")
				}
				for _, line := range lines {
					if line != "" {
						pterm.Info.Printf("%s %s\n", nodeLabel, line)
					}
				}

				resultCh <- res

				if res.ResultStatus != db.ResultStatusSucceeded {
					allOk = false
					break
				}
			}

			mu.Lock()
			durMs := 0
			if allOk && lastRes.DurationMs != nil {
				durMs = *lastRes.DurationMs
			}
			exitCode := 1
			if !allOk && lastRes.ExitCode != nil {
				exitCode = *lastRes.ExitCode
			}

			if ctx.Err() != nil {
				skipped++
				if spinner != nil {
					spinner.Warning(fmt.Sprintf("%s Skipped", nodeLabel))
				} else {
					pterm.Warning.Printf("%s Skipped\n", nodeLabel)
				}
			} else if allOk {
				succeeded++
				if spinner != nil {
					spinner.Success(fmt.Sprintf("%s %s (%dms, exit 0)", nodeLabel, displayCmd, durMs))
				} else {
					pterm.Success.Printf("%s %s (%dms, exit 0)\n", nodeLabel, displayCmd, durMs)
				}
			} else {
				failed++
				if spinner != nil {
					spinner.Fail(fmt.Sprintf("%s %s (exit %d)", nodeLabel, displayCmd, exitCode))
				} else {
					pterm.Error.Printf("%s %s (exit %d)\n", nodeLabel, displayCmd, exitCode)
				}
			}
			mu.Unlock()
		}
	}

	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go worker()
	}

	wg.Wait()

	if ctx.Err() == nil {
		updateCounts(false)
		if isMultiActive {
			multi.Stop()
		}
		// If zero padding is required for RUN-NNN, we format it. But runId might be just the DB ID.
		fmt.Printf("\n┌ Cluster Run RUN-%d ─────────────────┐\n", runId)
		fmt.Printf("│ Nodes: %d  OK: %d  Failed: %d  Skipped: %d │\n", totalNodes, succeeded, failed, skipped)
		fmt.Println("└────────────────────────────────────────┘")
	}
	close(resultCh)
}
