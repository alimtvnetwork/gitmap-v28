package cmdssh

import (
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func resolveExecTargetAndArgs(conns []db.SSHConnection, opts seOptions) ([]db.SSHConnection, []string) {
	args := opts.Args
	target := opts.Target
	if target == "" && opts.IP != "" {
		target = opts.IP
	}

	if target == "" && len(args) > 1 && isKnownTarget(conns, args[0]) {
		target = args[0]
		args = args[1:]
	}

	args = resolveIPCommandArgs(args)
	if target != "" && target != "all" {
		conns = filterConnectionsByTarget(conns, target)
	}

	return conns, args
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
