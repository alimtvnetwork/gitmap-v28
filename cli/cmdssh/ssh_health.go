package cmdssh

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// SSHHealthOptions specifies parameters for machine health and reachability checks.
type SSHHealthOptions struct {
	Target  string
	Port    int
	Timeout time.Duration
}

// SSHHealthResult represents connectivity health metrics of a machine.
type SSHHealthResult struct {
	Status   string
	IsOnline bool
	Alias    string
	IP       string
	User     string
	Port     int
	Latency  time.Duration
	Details  string
}

func resolveHealthPort(port int) int {
	if port > 0 {
		return port
	}
	return 22
}

func resolveHealthTimeout(d time.Duration) time.Duration {
	if d > 0 {
		return d
	}
	return 1500 * time.Millisecond
}

func classifyProbeError(err error) string {
	if err == nil {
		return "reachable"
	}
	msg := strings.ToLower(err.Error())
	if strings.Contains(msg, "refused") {
		return "connection refused"
	}
	if strings.Contains(msg, "timeout") {
		return "connection timed out"
	}
	if strings.Contains(msg, "no route") {
		return "no route to host"
	}
	return err.Error()
}

func buildOnlineResult(host store.SSHHost, port int, latency time.Duration) SSHHealthResult {
	return SSHHealthResult{
		Status:   "ONLINE",
		IsOnline: true,
		Alias:    host.Alias,
		IP:       host.IP,
		User:     host.Username,
		Port:     port,
		Latency:  latency,
		Details:  "reachable",
	}
}

func buildOfflineResult(host store.SSHHost, port int, reason string) SSHHealthResult {
	return SSHHealthResult{
		Status:   "OFFLINE",
		IsOnline: false,
		Alias:    host.Alias,
		IP:       host.IP,
		User:     host.Username,
		Port:     port,
		Latency:  0,
		Details:  reason,
	}
}

func probeHostHealth(ctx context.Context, host store.SSHHost, port int, timeout time.Duration) SSHHealthResult {
	addr := net.JoinHostPort(host.IP, strconv.Itoa(port))
	dialer := net.Dialer{Timeout: timeout}
	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return buildOfflineResult(host, port, classifyProbeError(err))
	}
	_ = conn.Close()
	return buildOnlineResult(host, port, time.Since(start).Round(time.Millisecond))
}

func findHostByTarget(hosts []store.SSHHost, target string) (store.SSHHost, bool) {
	for _, h := range hosts {
		if h.Alias == target || h.IP == target {
			return h, true
		}
	}
	return store.SSHHost{}, false
}

func loadAllRegisteredHosts(ctx context.Context) ([]store.SSHHost, *apperror.AppError) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return nil, apperror.New("loadAllRegisteredHosts", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	hosts, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		return nil, apperror.WrapSimple(err, "loadAllRegisteredHosts")
	}
	return hosts, nil
}

func createAdhocHost(target string) store.SSHHost {
	return store.SSHHost{
		Alias:    "-",
		IP:       target,
		Username: "-",
	}
}

func matchOrAdhocHost(hosts []store.SSHHost, target string) []store.SSHHost {
	match, isFound := findHostByTarget(hosts, target)
	if isFound {
		return []store.SSHHost{match}
	}
	return []store.SSHHost{createAdhocHost(target)}
}

func resolveTargetHosts(ctx context.Context, target string) ([]store.SSHHost, *apperror.AppError) {
	hosts, err := loadAllRegisteredHosts(ctx)
	if target == "" {
		return hosts, err
	}
	if err != nil {
		return []store.SSHHost{createAdhocHost(target)}, nil
	}
	return matchOrAdhocHost(hosts, target), nil
}

func probeAllHostsConcurrent(ctx context.Context, hosts []store.SSHHost, port int, timeout time.Duration) []SSHHealthResult {
	results := make([]SSHHealthResult, len(hosts))
	var wg sync.WaitGroup
	for i, h := range hosts {
		wg.Add(1)
		go func(idx int, targetHost store.SSHHost) {
			defer wg.Done()
			results[idx] = probeHostHealth(ctx, targetHost, port, timeout)
		}(i, h)
	}
	wg.Wait()
	return results
}

func printHealthTableHeader(w io.Writer) {
	fmt.Fprintln(w, "STATUS\tALIAS\tIP\tUSER\tPORT\tLATENCY\tDETAILS")
}

func formatLatencyString(res SSHHealthResult) string {
	if res.IsOnline {
		return res.Latency.String()
	}
	return "-"
}

func printHealthTableRow(w io.Writer, res SSHHealthResult) {
	fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\t%s\t%s\n",
		res.Status, res.Alias, res.IP, res.User, res.Port,
		formatLatencyString(res), res.Details)
}

func printHealthTable(out io.Writer, results []SSHHealthResult) error {
	w := tabwriter.NewWriter(out, 0, 0, 2, ' ', 0)
	printHealthTableHeader(w)
	for _, res := range results {
		printHealthTableRow(w, res)
	}
	return w.Flush()
}

func handleEmptyHostsList(out io.Writer) ([]SSHHealthResult, *apperror.AppError) {
	fmt.Fprintln(out, "No registered SSH machines found.\nEnroll a machine: gitmap sj <ip> [alias]")
	return []SSHHealthResult{}, nil
}

// ExecuteHealthCheck evaluates network health and latency for target or all registered hosts.
func ExecuteHealthCheck(ctx context.Context, out io.Writer, opts SSHHealthOptions) ([]SSHHealthResult, *apperror.AppError) {
	port := resolveHealthPort(opts.Port)
	timeout := resolveHealthTimeout(opts.Timeout)
	hosts, err := resolveTargetHosts(ctx, opts.Target)
	if err != nil {
		return nil, err
	}
	if len(hosts) == 0 {
		return handleEmptyHostsList(out)
	}
	results := probeAllHostsConcurrent(ctx, hosts, port, timeout)
	if flushErr := printHealthTable(out, results); flushErr != nil {
		return nil, apperror.WrapSimple(flushErr, "ExecuteHealthCheck")
	}
	return results, nil
}
