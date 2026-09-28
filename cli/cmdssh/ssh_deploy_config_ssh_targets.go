// Package cmdssh — ssh_deploy_config_ssh_targets.go parses flags and filters targets for deploy config ssh.
package cmdssh

import (
	"net"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func resolveDeployConfigTargets(conns []db.SSHConnection, target, except string) []db.SSHConnection {
	if target != "" && strings.ToLower(target) != "all" {
		var matched []db.SSHConnection
		for _, c := range conns {
			if strings.EqualFold(c.Alias, target) || strings.EqualFold(c.IPAddress, target) {
				matched = append(matched, c)
			}
		}
		if len(matched) > 0 {
			return FilterSSHConnectionsByExcept(matched, except)
		}
	}
	filtered := FilterSSHConnectionsByExcept(conns, except)
	return filterOutLocalFleetIPs(filtered)
}

func filterOutLocalFleetIPs(conns []db.SSHConnection) []db.SSHConnection {
	localIPs := getLocalHostIPSet()
	var out []db.SSHConnection
	for _, c := range conns {
		if c.IPAddress == "127.0.0.1" || strings.EqualFold(c.IPAddress, "localhost") || localIPs[c.IPAddress] {
			continue
		}
		out = append(out, c)
	}
	return out
}

func getLocalHostIPSet() map[string]bool {
	ips := make(map[string]bool)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ips[ipnet.IP.String()] = true
			}
		}
	}
	return ips
}

func parseDeployConfigSSHFlags(args []string) (string, string, string, bool, bool, bool) {
	var target, filePath string
	var exceptParts []string
	isDryRun, isJSON, isHelp := false, false, false

	for i := 0; i < len(args); i++ {
		a := args[i]
		low := strings.ToLower(a)
		if low == "config" || low == "ssh" || low == "node-config" || low == "nodeconfig" || low == "nc" || low == "deploy" {
			continue
		}
		if low == "--help" || low == "-h" || low == "help" {
			isHelp = true
			continue
		}
		if low == "--dry-run" || low == "-n" {
			isDryRun = true
			continue
		}
		if low == "--json" || low == "-j" {
			isJSON = true
			continue
		}
		if (low == "--file" || low == "-f") && i+1 < len(args) {
			filePath = args[i+1]
			i++
			continue
		}
		if strings.HasPrefix(low, "--file=") {
			filePath = a[strings.IndexByte(a, '=')+1:]
			continue
		}
		if strings.HasPrefix(low, "-f=") {
			filePath = a[3:]
			continue
		}
		if isExceptOrExcepFlag(low) {
			for i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				exceptParts = append(exceptParts, args[i+1])
				i++
			}
			continue
		}
		if strings.HasPrefix(low, "--except=") || strings.HasPrefix(low, "--excep=") || strings.HasPrefix(low, "--accept=") || strings.HasPrefix(low, "--exclude=") {
			idx := strings.IndexByte(a, '=')
			exceptParts = append(exceptParts, a[idx+1:])
			continue
		}
		if !strings.HasPrefix(a, "-") && target == "" {
			target = a
		}
	}
	if target == "" {
		target = "all"
	}
	return target, filePath, strings.Join(exceptParts, ","), isDryRun, isJSON, isHelp
}
