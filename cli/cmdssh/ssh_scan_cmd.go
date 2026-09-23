package cmdssh

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func runSSHScanCLI(args []string) error {
	ctx := context.Background()
	InvalidateLivenessCache()

	isRegisteredOnly, opts := parseScanArgs(args)
	if isRegisteredOnly {
		return runRegisteredHostsProbe(ctx)
	}

	_, appErr := ExecuteSubnetScan(ctx, os.Stdout, opts)
	if appErr != nil {
		return appErr
	}
	return nil
}

func parseScanPort(val string, fallback int) int {
	if p, err := strconv.Atoi(val); err == nil && p > 0 {
		return p
	}
	return fallback
}

func parseScanWorkers(val string, fallback int) int {
	if w, err := strconv.Atoi(val); err == nil && w > 0 {
		return w
	}
	return fallback
}

func parseScanTimeout(val string, fallback time.Duration) time.Duration {
	if d, err := time.ParseDuration(val); err == nil && d > 0 {
		return d
	}
	return fallback
}

func parseScanArgs(args []string) (bool, SSHScanOptions) {
	opts := SSHScanOptions{
		Port:    22,
		Timeout: 600 * time.Millisecond,
		Workers: 50,
	}
	isRegisteredOnly := false

	idx := 0
	for idx < len(args) {
		a := args[idx]
		if a == "--registered" || a == "-r" {
			isRegisteredOnly = true
			idx++
			continue
		}
		if (a == "--port" || a == "-p") && idx+1 < len(args) {
			opts.Port = parseScanPort(args[idx+1], opts.Port)
			idx += 2
			continue
		}
		if (a == "--workers" || a == "-w") && idx+1 < len(args) {
			opts.Workers = parseScanWorkers(args[idx+1], opts.Workers)
			idx += 2
			continue
		}
		if (a == "--timeout" || a == "-t") && idx+1 < len(args) {
			opts.Timeout = parseScanTimeout(args[idx+1], opts.Timeout)
			idx += 2
			continue
		}
		if !strings.HasPrefix(a, "-") && opts.Subnet == "" {
			opts.Subnet = a
			idx++
			continue
		}
		idx++
	}

	return isRegisteredOnly, opts
}

func runRegisteredHostsProbe(ctx context.Context) error {
	hosts, appErr := loadAllRegisteredHosts(ctx)
	if appErr != nil {
		return appErr
	}

	if len(hosts) == 0 {
		fmt.Printf("\n  %sNo SSH machines registered. Run 'gitmap ssh join <user@ip> [alias]' to add one.%s\n\n",
			constants.ColorYellow, constants.ColorReset)
		return nil
	}

	results := probeAllHosts(ctx, hosts)
	renderScanTable(results)
	return nil
}

func probeAllHosts(ctx context.Context, hosts []store.SSHHost) []SSHHealthResult {
	var results []SSHHealthResult
	timeout := 1200 * time.Millisecond

	for _, h := range hosts {
		res := probeHostHealth(ctx, h, resolveHealthPort(h.Port), timeout)
		setCachedLiveness(h.IP+":22", res.IsOnline, res.Details)
		results = append(results, res)
	}

	return results
}
