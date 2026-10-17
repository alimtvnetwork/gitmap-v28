package cmdfleet

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os/exec"
	"sort"
	"strconv"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

const nodesVersionHelp = `nodes version — version skew table across fleet nodes.

Description:
  Queries ` + "`gitmap --version`" + ` on every fleet node over ssh and
  renders a skew table (node, version, status). Nodes that do not answer
  are reported as unreachable instead of failing the command.

Usage:
  gitmap nodes version [--timeout <seconds>]
  gitmap nv [--timeout <seconds>]

Flags:
  --timeout <n>   ssh connect timeout in seconds (default: 6)
  --help, -h      show this help

Example:
  gitmap nodes version
`

// nodeVersion is one probed node's result.
type nodeVersion struct {
	name    string
	version string
	status  string
}

// RunNodesVersion implements `gitmap nodes version` (alias: nv).
func RunNodesVersion(args []string) *appfault.AppError {
	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Print(nodesVersionHelp)
			return nil
		}
	}
	fs := flag.NewFlagSet("nodes version", flag.ContinueOnError)
	timeout := fs.Int("timeout", 6, "ssh connect timeout in seconds")
	if parseErr := fs.Parse(args); parseErr != nil {
		fmt.Print(nodesVersionHelp)
		return apperror.WrapSimple(parseErr, "nodes.version.flags")
	}

	nodes, note, loadErr := loadFleetNodes()
	if loadErr != nil {
		return loadErr
	}
	if note != "" {
		fmt.Printf("nodes version: %s\n", note)
		return nil
	}
	if len(nodes) == 0 {
		fmt.Println("nodes version: fleet registry is empty; no nodes to probe.")
		return nil
	}

	results := make([]nodeVersion, 0, len(nodes))
	for _, n := range nodes {
		ver, isReachable := probeNodeVersion(n, *timeout)
		status := "unreachable"
		if isReachable {
			status = "ok"
		}
		results = append(results, nodeVersion{name: n.Name, version: ver, status: status})
	}

	latest := latestVersion(results)
	for i, r := range results {
		if r.status == "ok" && r.version != latest && latest != "" {
			results[i].status = "behind"
		}
		if r.status == "ok" && r.version == latest {
			results[i].status = "latest"
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].name < results[j].name })
	fmt.Printf("  %-28s %-16s %s\n", "NODE", "VERSION", "STATUS")
	fmt.Printf("  %-28s %-16s %s\n", "----", "-------", "------")
	for _, r := range results {
		ver := r.version
		if ver == "" {
			ver = "-"
		}
		fmt.Printf("  %-28s %-16s %s\n", r.name, ver, r.status)
	}
	if latest != "" {
		fmt.Printf("fleet latest: %s\n", latest)
	}

	return nil
}

// probeNodeVersion runs `gitmap --version` on the node over ssh.
func probeNodeVersion(n FleetNode, timeout int) (string, bool) {
	sshArgs := []string{
		"-o", "BatchMode=yes",
		"-o", "ConnectTimeout=" + strconv.Itoa(timeout),
		"-o", "StrictHostKeyChecking=accept-new",
		sshTarget(n), "gitmap --version",
	}
	cmd := exec.Command("ssh", sshArgs...)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	if runErr := cmd.Run(); runErr != nil {
		return "", false
	}
	version := parseVersionToken(stdout.String())
	if version == "" {
		return "", false
	}

	return version, true
}

// parseVersionToken extracts the first vX.Y.Z-like token from output.
func parseVersionToken(output string) string {
	for _, field := range strings.Fields(output) {
		trimmed := strings.Trim(field, ",;()")
		if strings.HasPrefix(trimmed, "v") && strings.Contains(trimmed, ".") {
			return trimmed
		}
	}

	return ""
}

// latestVersion returns the highest semantic-ish version seen.
func latestVersion(results []nodeVersion) string {
	best := ""
	for _, r := range results {
		if r.status != "ok" || r.version == "" {
			continue
		}
		if best == "" || compareVersions(r.version, best) > 0 {
			best = r.version
		}
	}

	return best
}

// compareVersions compares v-prefixed dotted versions numerically.
func compareVersions(a, b string) int {
	pa := strings.Split(strings.TrimPrefix(a, "v"), ".")
	pb := strings.Split(strings.TrimPrefix(b, "v"), ".")
	for i := 0; i < len(pa) && i < len(pb); i++ {
		na, _ := strconv.Atoi(pa[i])
		nb, _ := strconv.Atoi(pb[i])
		if na != nb {
			if na > nb {
				return 1
			}
			return -1
		}
	}
	if len(pa) != len(pb) {
		if len(pa) > len(pb) {
			return 1
		}
		return -1
	}

	return 0
}
