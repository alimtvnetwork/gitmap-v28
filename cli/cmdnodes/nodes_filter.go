// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"context"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// NodeFilterOptions specifies criteria for selecting and filtering fleet nodes.
type NodeFilterOptions struct {
	Target      string   `json:"target,omitempty"`
	Except      []string `json:"except,omitempty"`
	Include     []string `json:"include,omitempty"`
	IncludeMain bool     `json:"include_main,omitempty"`
	OpenOnly    bool     `json:"open_only,omitempty"`
}

// ParseNodeFilterOptions extracts node filtering options from command-line arguments.
func ParseNodeFilterOptions(args []string) NodeFilterOptions {
	var opts NodeFilterOptions

	for i := 0; i < len(args); i++ {
		step := parseSingleNodeFilterFlag(&opts, args, i)
		if step > 0 {
			i += (step - 1)
			continue
		}

		arg := args[i]
		if !strings.HasPrefix(arg, "-") && opts.Target == "" {
			opts.Target = strings.TrimSpace(arg)
		}
	}

	return opts
}

func parseSingleNodeFilterFlag(opts *NodeFilterOptions, args []string, i int) int {
	arg := args[i]
	if arg == "--include-main" {
		opts.IncludeMain = true
		return 1
	}

	if arg == "--open-only" {
		opts.OpenOnly = true
		return 1
	}

	if step := parseTargetParamFlag(opts, args, i); step > 0 {
		return step
	}

	if step := parseExceptParamFlag(opts, args, i); step > 0 {
		return step
	}

	return parseIncludeParamFlag(opts, args, i)
}

func parseTargetParamFlag(opts *NodeFilterOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--target=") {
		opts.Target = strings.TrimSpace(strings.TrimPrefix(arg, "--target="))
		return 1
	}

	if strings.HasPrefix(arg, "-t=") {
		opts.Target = strings.TrimSpace(strings.TrimPrefix(arg, "-t="))
		return 1
	}

	if (arg == "--target" || arg == "-t") && i+1 < len(args) {
		opts.Target = strings.TrimSpace(args[i+1])
		return 2
	}

	return 0
}

func parseExceptParamFlag(opts *NodeFilterOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--except=") {
		appendSplitTokens(&opts.Except, strings.TrimPrefix(arg, "--except="))
		return 1
	}

	if strings.HasPrefix(arg, "-e=") {
		appendSplitTokens(&opts.Except, strings.TrimPrefix(arg, "-e="))
		return 1
	}

	if strings.HasPrefix(arg, "--exclude=") {
		appendSplitTokens(&opts.Except, strings.TrimPrefix(arg, "--exclude="))
		return 1
	}

	if (arg == "--except" || arg == "-e" || arg == "--exclude") && i+1 < len(args) {
		appendSplitTokens(&opts.Except, args[i+1])
		return 2
	}

	return 0
}

func parseIncludeParamFlag(opts *NodeFilterOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--include=") {
		appendSplitTokens(&opts.Include, strings.TrimPrefix(arg, "--include="))
		return 1
	}

	if strings.HasPrefix(arg, "--accept=") {
		appendSplitTokens(&opts.Include, strings.TrimPrefix(arg, "--accept="))
		return 1
	}

	if (arg == "--include" || arg == "--accept") && i+1 < len(args) {
		appendSplitTokens(&opts.Include, args[i+1])
		return 2
	}

	return 0
}

func appendSplitTokens(target *[]string, raw string) {
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			*target = append(*target, trimmed)
		}
	}
}

// FilterFleetNodes filters fleet connections according to target, except, include, and local exclusions.
func FilterFleetNodes(conns []db.SSHConnection, opts NodeFilterOptions) []db.SSHConnection {
	var candidates []db.SSHConnection

	for _, c := range conns {
		if !isConnectionCandidate(c, opts) {
			continue
		}

		candidates = append(candidates, c)
	}

	if !opts.OpenOnly || len(candidates) == 0 {
		return candidates
	}

	return filterOnlineNodesConcurrently(candidates, 1500*time.Millisecond)
}

func isConnectionCandidate(c db.SSHConnection, opts NodeFilterOptions) bool {
	if isLocalMachineConnection(c) {
		return false
	}

	if !opts.IncludeMain && isMainNode(c) {
		return false
	}

	if opts.Target != "" && !matchesTargetNode(c, opts.Target) {
		return false
	}

	if isNodeExcluded(c, opts.Except) {
		return false
	}

	if len(opts.Include) > 0 && !isNodeIncluded(c, opts.Include) {
		return false
	}

	return true
}

func isMainNode(c db.SSHConnection) bool {
	return strings.EqualFold(c.Alias, "main")
}

func isNodeExcluded(c db.SSHConnection, exceptList []string) bool {
	for _, exc := range exceptList {
		clean := strings.TrimSpace(exc)
		if clean == "" {
			continue
		}

		if strings.EqualFold(c.Alias, clean) || strings.EqualFold(c.IPAddress, clean) {
			return true
		}
	}

	return false
}

func isNodeIncluded(c db.SSHConnection, includeList []string) bool {
	for _, inc := range includeList {
		clean := strings.TrimSpace(inc)
		if clean == "" {
			continue
		}

		if strings.EqualFold(c.Alias, clean) || strings.EqualFold(c.IPAddress, clean) {
			return true
		}
	}

	return false
}

type nodeProbeResult struct {
	conn     db.SSHConnection
	isOnline bool
}

func filterOnlineNodesConcurrently(conns []db.SSHConnection, timeout time.Duration) []db.SSHConnection {
	results := make([]nodeProbeResult, len(conns))
	var wg sync.WaitGroup
	ctx, cancel := context.WithTimeout(context.Background(), timeout+500*time.Millisecond)
	defer cancel()

	for i, c := range conns {
		wg.Add(1)
		go func(idx int, conn db.SSHConnection) {
			defer wg.Done()
			online := probeNodeTCPConnection(ctx, conn, timeout)
			results[idx] = nodeProbeResult{conn: conn, isOnline: online}
		}(i, c)
	}

	wg.Wait()

	var onlineNodes []db.SSHConnection
	for _, r := range results {
		if r.isOnline {
			onlineNodes = append(onlineNodes, r.conn)
		}
	}

	return onlineNodes
}

func probeNodeTCPConnection(ctx context.Context, conn db.SSHConnection, timeout time.Duration) bool {
	host, port := resolveNodeHostAndPort(conn.IPAddress)
	addr := net.JoinHostPort(host, port)
	d := net.Dialer{Timeout: timeout}

	rawConn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return false
	}

	_ = rawConn.Close()
	return true
}

func resolveNodeHostAndPort(ipAddr string) (string, string) {
	trimmed := strings.TrimSpace(ipAddr)
	if host, port, err := net.SplitHostPort(trimmed); err == nil && host != "" && port != "" {
		return host, port
	}

	return trimmed, "22"
}
