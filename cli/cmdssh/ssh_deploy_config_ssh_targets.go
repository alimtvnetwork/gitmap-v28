// Package cmdssh — ssh_deploy_config_ssh_targets.go parses flags and filters targets for deploy config ssh.
package cmdssh

import (
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func parseDeployConfigSSHFlags(args []string) DeployConfigSSHOptions {
	opts := DeployConfigSSHOptions{Target: "all"}
	for i := 0; i < len(args); i++ {
		i = consumeDeployConfigFlag(args, i, &opts)
	}
	return opts
}

func consumeDeployConfigFlag(args []string, i int, opts *DeployConfigSSHOptions) int {
	a := args[i]
	low := strings.ToLower(a)
	if isConfigSkipToken(low) {
		return i
	}
	if low == "--dry-run" || low == "-n" {
		opts.IsDryRun = true
		return i
	}
	if low == "--json" || low == "-j" {
		opts.IsJSON = true
		return i
	}
	if low == "--force" {
		opts.IsForce = true
		return i
	}
	return consumeValueOrTargetDeployConfig(args, i, opts)
}

func consumeValueOrTargetDeployConfig(args []string, i int, opts *DeployConfigSSHOptions) int {
	a := args[i]
	low := strings.ToLower(a)
	if isFileFlag(low) && i+1 < len(args) {
		opts.FilePath = args[i+1]
		return i + 1
	}
	if strings.HasPrefix(low, "--file=") || strings.HasPrefix(low, "-f=") {
		opts.FilePath = a[strings.IndexByte(a, '=')+1:]
		return i
	}
	if isExceptOrExcepFlag(low) && i+1 < len(args) {
		opts.Except = appendToken(opts.Except, args[i+1])
		return i + 1
	}
	if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--exclude=") {
		opts.Except = appendToken(opts.Except, a[strings.IndexByte(a, '=')+1:])
		return i
	}
	if !strings.HasPrefix(a, "-") && (opts.Target == "all" || opts.Target == "") {
		opts.Target = a
	}
	return i
}

func isConfigSkipToken(s string) bool {
	return s == "config" || s == "ssh" || s == "deploy" || s == "node-config" || s == "nc" || s == "all"
}

func isFileFlag(s string) bool {
	return s == "--file" || s == "-f" || s == "--from-file"
}

func appendToken(existing, addition string) string {
	if existing == "" {
		return addition
	}
	return existing + "," + addition
}

func filterConfigDeployTargets(targets []db.SSHConnection, opts DeployConfigSSHOptions) []db.SSHConnection {
	filtered := FilterSSHConnectionsByExcept(targets, opts.Except)
	if !isDeployAllTarget(opts.Target) {
		return filtered
	}
	var remoteOnly []db.SSHConnection
	for _, c := range filtered {
		if isLocalMachineIP(c.IPAddress) {
			continue
		}
		remoteOnly = append(remoteOnly, c)
	}
	return remoteOnly
}

func isLocalMachineIP(ip string) bool {
	if ip == "" || ip == "127.0.0.1" || strings.EqualFold(ip, "localhost") {
		return true
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ipnet, ok := a.(*net.IPNet)
		if ok && !ipnet.IP.IsLoopback() && ipnet.IP.String() == ip {
			return true
		}
	}
	return false
}
