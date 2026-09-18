package cmdssh

import (
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func resolveExecTargetAndArgs(conns []db.SSHConnection, opts seOptions) ([]db.SSHConnection, []string) {
	args := opts.Args
	target := opts.Target
	if target == "" && opts.IP != "" {
		target = opts.IP
	}
	if target == "" {
		target, args = extractTargetPrefix(conns, args)
	}
	args = resolveIPCommandArgs(args)
	if target != "" && target != "all" {
		conns = filterConnectionsByTarget(conns, target)
	}
	return conns, args
}

func extractTargetPrefix(conns []db.SSHConnection, args []string) (string, []string) {
	if len(args) <= 1 || !isTargetExpression(conns, args[0]) {
		return "", args
	}
	var targets []string
	idx := 0
	for idx < len(args)-1 && isTargetExpression(conns, args[idx]) {
		targets = append(targets, args[idx])
		idx++
	}
	return strings.Join(targets, ","), args[idx:]
}

func isTargetExpression(conns []db.SSHConnection, candidate string) bool {
	if isKnownTarget(conns, candidate) || net.ParseIP(candidate) != nil {
		return true
	}
	if strings.Contains(candidate, ",") {
		return hasAnyKnownTarget(conns, strings.Split(candidate, ","))
	}
	return false
}

func hasAnyKnownTarget(conns []db.SSHConnection, parts []string) bool {
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" && (isKnownTarget(conns, trimmed) || net.ParseIP(trimmed) != nil) {
			return true
		}
	}
	return false
}

func isKnownTarget(conns []db.SSHConnection, candidate string) bool {
	for _, c := range conns {
		if strings.EqualFold(c.Alias, candidate) || c.IPAddress == candidate {
			return true
		}
	}

	return false
}

func resolveIPCommandArgs(args []string) []string {
	if len(args) == 1 && strings.EqualFold(args[0], "ip") {
		return []string{"sh", "-c", "ip -br a 2>/dev/null || ip a 2>/dev/null || hostname -I 2>/dev/null || ifconfig"}
	}

	return args
}
