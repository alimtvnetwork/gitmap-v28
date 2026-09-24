package cmdssh

import (
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// FleetParallelOptions defines execution parameters for parallel fleet operations.
type FleetParallelOptions struct {
	Target   string
	Except   string
	TaskName string
	IsDryRun bool
}

// FleetNodeResult represents the execution outcome on a single SSH node.
type FleetNodeResult struct {
	Alias      string
	IP         string
	Success    bool
	Status     string
	Output     string
	Error      error
	DurationMs int64
}

// FleetWorkerFunc defines the function signature for executing work on an individual node.
type FleetWorkerFunc func(c db.SSHConnection) (string, error)

// RunParallelFleetExecution executes workerFunc concurrently across all matching SSH connections.
func RunParallelFleetExecution(conns []db.SSHConnection, opts FleetParallelOptions, worker FleetWorkerFunc) []FleetNodeResult {
	targets := prepareFleetTargets(conns, opts)
	if len(targets) == 0 {
		printFleetNoTargets(opts.TaskName)
		return nil
	}
	results := dispatchParallelWorkers(targets, opts, worker)
	PrintFleetSummary(opts.TaskName, results)
	return results
}

func prepareFleetTargets(conns []db.SSHConnection, opts FleetParallelOptions) []db.SSHConnection {
	filtered := filterConnectionsByTarget(conns, opts.Target)
	if opts.Except != "" {
		filtered = filterSSHConns(filtered, opts.Except)
	}
	return filtered
}

func dispatchParallelWorkers(targets []db.SSHConnection, opts FleetParallelOptions, worker FleetWorkerFunc) []FleetNodeResult {
	results := make([]FleetNodeResult, len(targets))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for idx, c := range targets {
		wg.Add(1)
		go executeSingleWorker(idx, c, opts, worker, &results, &mu, &wg)
	}
	wg.Wait()
	return results
}

func executeSingleWorker(idx int, c db.SSHConnection, opts FleetParallelOptions, worker FleetWorkerFunc, results *[]FleetNodeResult, mu *sync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()
	PrintFleetStart(c.Alias, c.IPAddress, opts.TaskName)
	start := time.Now()
	out, err := worker(c)
	dur := time.Since(start).Milliseconds()

	res := buildFleetNodeResult(c, out, err, dur)
	PrintFleetDone(res)

	mu.Lock()
	(*results)[idx] = res
	mu.Unlock()
}

func buildFleetNodeResult(c db.SSHConnection, out string, err error, dur int64) FleetNodeResult {
	status := "success"
	success := true
	if err != nil {
		status = "failed"
		success = false
	}
	return FleetNodeResult{
		Alias:      c.Alias,
		IP:         c.IPAddress,
		Success:    success,
		Status:     status,
		Output:     out,
		Error:      err,
		DurationMs: dur,
	}
}
