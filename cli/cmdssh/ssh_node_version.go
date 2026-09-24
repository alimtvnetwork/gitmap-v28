package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"golang.org/x/crypto/ssh"
)

func isNodeVersionFlag(arg string) bool {
	return arg == "-v" || arg == "--version" || arg == "-version" || arg == "version"
}

func isNodesVersionRequest(args []string) bool {
	for i, a := range args {
		if isNodeVersionFlag(a) {
			return true
		}
		if a == "gitmap" && i+1 < len(args) && isNodeVersionFlag(args[i+1]) {
			return true
		}
	}
	return false
}

func hasJSONOutputFlag(args []string) bool {
	for _, a := range args {
		if a == "--json" {
			return true
		}
	}
	return false
}

func resolveRemoteVersionCommand(osType string) (string, string) {
	if isWindowsOS(osType) {
		cmd := `where.exe gitmap >nul 2>nul && gitmap.exe --version 2>nul || powershell -NoProfile -Command "gitmap --version" 2>$null || if exist "%LOCALAPPDATA%\gitmap-cli\gitmap.exe" ("%LOCALAPPDATA%\gitmap-cli\gitmap.exe" --version) else if exist "%LOCALAPPDATA%\gitmap\bin\gitmap.exe" ("%LOCALAPPDATA%\gitmap\bin\gitmap.exe" --version) else (exit 1)`
		return cmd, "cmd"
	}

	cmd := wrapUnixPath(`which gitmap >/dev/null 2>&1 && gitmap --version 2>/dev/null || [ -x "$HOME/.local/bin/gitmap" ] && "$HOME/.local/bin/gitmap" --version 2>/dev/null || [ -x "/usr/local/bin/gitmap" ] && /usr/local/bin/gitmap --version 2>/dev/null || exit 1`)
	return cmd, "sh"
}

func parseRemoteVersionOutput(raw string, err error) (string, bool) {
	if err != nil {
		return "not installed", false
	}
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return "not installed", false
	}
	lines := strings.Split(trimmed, "\n")
	for _, l := range lines {
		clean := strings.TrimSpace(l)
		if clean == "" {
			continue
		}
		lower := strings.ToLower(clean)
		if strings.Contains(lower, "not recognized") || strings.Contains(lower, "not found") {
			return "not installed", false
		}
		if strings.HasPrefix(lower, "gitmap v") || strings.HasPrefix(lower, "gitmap version") || strings.HasPrefix(lower, "v") {
			return clean, true
		}
	}
	lastLine := strings.TrimSpace(lines[len(lines)-1])
	if strings.Contains(strings.ToLower(lastLine), "error") || strings.Contains(strings.ToLower(lastLine), "not recognized") {
		return "not installed", false
	}
	return lastLine, true
}

func queryNodeVersionViaSSH(client *ssh.Client, osType string) (string, bool) {
	cmd, shell := resolveRemoteVersionCommand(osType)
	out, err := crypto.RunCommand(client, cmd, shell)
	return parseRemoteVersionOutput(out, err)
}

func probeOnlineNodeVersion(c db.SSHConnection) NodeVersionInfo {
	info := NodeVersionInfo{
		Alias:    c.Alias,
		Host:     fmt.Sprintf("%s:22", c.IPAddress),
		Role:     "worker",
		Status:   "● ready",
		IsOnline: true,
	}

	header := fmt.Sprintf("[%s]", c.Alias)
	client, isConnected := connectSSHNode(c, header)
	if !isConnected {
		info.Version = "not installed"
		info.ErrorMsg = "connection failed"
		return info
	}
	defer client.Close()

	osType := c.OS
	if probed := probeRemoteOSType(client); probed != "" {
		osType = probed
	}

	ver, isInstalled := queryNodeVersionViaSSH(client, osType)
	info.Version = ver
	info.IsInstalled = isInstalled
	info.Comparison = resolveNodeVersionComparison(ver, isInstalled)
	return info
}

// CompareSemverStrings compares two version strings (e.g. "v6.326.0" vs "v6.324.0").
// Returns 1 if v1 > v2, -1 if v1 < v2, and 0 if equal.
func CompareSemverStrings(v1, v2 string) int {
	clean1 := strings.TrimPrefix(strings.TrimSpace(v1), "v")
	clean2 := strings.TrimPrefix(strings.TrimSpace(v2), "v")
	if clean1 == clean2 {
		return 0
	}
	p1 := strings.Split(clean1, ".")
	p2 := strings.Split(clean2, ".")
	maxLen := len(p1)
	if len(p2) > maxLen {
		maxLen = len(p2)
	}
	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(p1) {
			_, _ = fmt.Sscanf(p1[i], "%d", &n1)
		}
		if i < len(p2) {
			_, _ = fmt.Sscanf(p2[i], "%d", &n2)
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}

func resolveNodeVersionComparison(ver string, isInstalled bool) string {
	if !isInstalled {
		return "○ NOT INSTALLED"
	}
	cmp := CompareSemverStrings(ver, constants.Version)
	if cmp == 0 {
		return "● IDENTICAL (SAME)"
	}
	if cmp > 0 {
		return "▲ ABOVE (NEWER)"
	}
	return "▼ BELOW (OLDER)"
}

func probeNodeVersion(c db.SSHConnection) NodeVersionInfo {
	isOnline, _ := CheckConnLiveness(context.Background(), c.IPAddress, 22, 0)
	if !isOnline {
		return NodeVersionInfo{
			Alias:       c.Alias,
			Host:        fmt.Sprintf("%s:22", c.IPAddress),
			Role:        "worker",
			Status:      "○ offline",
			Version:     "(machine is off)",
			Comparison:  "○ OFFLINE",
			IsOnline:    false,
			IsInstalled: false,
		}
	}

	return probeOnlineNodeVersion(c)
}

func collectFleetNodeVersions(conns []db.SSHConnection) []NodeVersionInfo {
	results := make([]NodeVersionInfo, len(conns))
	var wg sync.WaitGroup

	for i, c := range conns {
		wg.Add(1)
		go func(idx int, conn db.SSHConnection) {
			defer wg.Done()
			results[idx] = probeNodeVersion(conn)
		}(i, c)
	}

	wg.Wait()
	return results
}

func renderNodeVersionTable(infos []NodeVersionInfo) {
	fmt.Println()
	fmt.Printf("  %-16s %-20s %-10s %-12s %-18s %s\n",
		"ALIAS", "HOST (IP:PORT)", "ROLE", "STATUS", "GITMAP VERSION", "COMPARISON (LOCAL: v"+constants.Version+")")
	fmt.Println("  " + strings.Repeat("-", 100))

	for _, info := range infos {
		statusColor := resolveNodeStatusColor(info.IsOnline)
		versionColor := resolveNodeVersionColor(info.IsInstalled, info.IsOnline)
		cmpColor := resolveComparisonColor(info.Comparison)
		fmt.Printf("  %-16s %-20s %-10s %s%-12s%s %s%-18s%s %s%s%s\n",
			info.Alias,
			info.Host,
			info.Role,
			statusColor, info.Status, constants.ColorReset,
			versionColor, info.Version, constants.ColorReset,
			cmpColor, info.Comparison, constants.ColorReset,
		)
	}
	fmt.Println()
}

func resolveComparisonColor(comparison string) string {
	if strings.Contains(comparison, "IDENTICAL") {
		return constants.ColorGreen
	}
	if strings.Contains(comparison, "ABOVE") {
		return constants.ColorCyan
	}
	if strings.Contains(comparison, "BELOW") {
		return constants.ColorYellow
	}
	return constants.ColorDim
}

func resolveNodeStatusColor(isOnline bool) string {
	if isOnline {
		return constants.ColorGreen
	}
	return constants.ColorDim
}

func resolveNodeVersionColor(isInstalled, isOnline bool) string {
	if !isOnline {
		return constants.ColorDim
	}
	if isInstalled {
		return constants.ColorCyan
	}
	return constants.ColorYellow
}

func renderNodeVersionJSON(infos []NodeVersionInfo) error {
	data, err := json.MarshalIndent(infos, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal node version json")
	}
	fmt.Println(string(data))
	return nil
}

func runSSHNodesVersion(args []string) error {
	conns, err := fetchAllSSHConnections()
	if err != nil {
		return apperror.WrapSimple(err, "fetchAllSSHConnections")
	}

	target := extractTargetFilter(args)
	filtered := filterConnectionsByTarget(conns, target)
	infos := collectFleetNodeVersions(filtered)

	if hasJSONOutputFlag(args) {
		return renderNodeVersionJSON(infos)
	}

	renderNodeVersionTable(infos)
	return nil
}

func extractTargetFilter(args []string) string {
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if (arg == "-t" || arg == "--target") && i+1 < len(args) {
			return args[i+1]
		}
		if strings.HasPrefix(arg, "--target=") {
			return strings.TrimPrefix(arg, "--target=")
		}
	}
	return "all"
}

// ResolveNodeVersionComparisonForTest exports resolveNodeVersionComparison for testing.
func ResolveNodeVersionComparisonForTest(ver string, isInstalled bool) string {
	return resolveNodeVersionComparison(ver, isInstalled)
}
