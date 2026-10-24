// Package cmdos — os_machine_alias.go provides cross-platform machine identity, OS hostname, and fleet alias management.
package cmdos

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// MachineIdentity represents the resolved identity of a local or remote SSH fleet machine.
type MachineIdentity struct {
	Sequence         int            `json:"sequence"`
	NodeID           string         `json:"nodeId"`
	IPAddress        string         `json:"ipAddress"`
	Alias            string         `json:"alias"`
	MachineName      string         `json:"machineName"`
	OSHostname       string         `json:"osHostname"`
	PreviousName     string         `json:"previousName,omitempty"`
	PreviousAlias    string         `json:"previousAlias,omitempty"`
	OSPlatform       string         `json:"osPlatform"`
	OSVersion        string         `json:"osVersion,omitempty"`
	Architecture     string         `json:"architecture,omitempty"`
	CPUCores         int            `json:"cpuCores,omitempty"`
	CurrentUser      string         `json:"currentUser,omitempty"`
	GitMapVersion    string         `json:"gitmapVersion,omitempty"`
	Scope            string         `json:"scope"`
	OSVersionDetails *OSVersionInfo `json:"osVersionDetails,omitempty"`
}

// RunMachineCLI executes the top-level or OS-nested machine command.
func RunMachineCLI(args []string) error {
	return dispatchMachineOrAlias("machine", args)
}

// RunAliasCLI executes the top-level or OS-nested alias command.
func RunAliasCLI(args []string) error {
	return dispatchMachineOrAlias("alias", args)
}

func dispatchMachineOrAlias(mode string, args []string) error {
	isSSH := hasFlag(args, "--ssh") || hasFlag(args, "-s")
	isJSON := hasFlag(args, "--json") || hasFlag(args, "-j")
	clean := filterPositionalMachineArgs(args)
	if len(clean) == 0 && isSSH {
		return renderFleetMachineAliasList(mode, isJSON)
	}
	if len(clean) == 0 {
		return renderSingleMachineIdentity(resolveLocalMachineIdentity(), isJSON)
	}
	if isOSHelpArg(clean[0]) {
		RenderMachineAliasHelp(mode)
		return executeMachineAliasList(mode, isSSH, isJSON)
	}
	return routeMachineAliasSubcmd(mode, strings.ToLower(clean[0]), clean[1:], isSSH, isJSON)
}

func routeMachineAliasSubcmd(mode, sub string, rest []string, isSSH, isJSON bool) error {
	switch sub {
	case "ls", "list", "status", "st", "show", "get":
		return executeMachineAliasList(mode, isSSH, isJSON)
	case "set", "change", "rename", "update":
		return executeMachineAliasSet(mode, rest, isSSH, isJSON)
	case "revert", "rollback", "undo":
		return executeMachineAliasRevert(mode, rest, isSSH, isJSON)
	default:
		return executeMachineAliasSet(mode, append([]string{sub}, rest...), isSSH, isJSON)
	}
}

func filterPositionalMachineArgs(args []string) []string {
	var out []string
	for _, a := range args {
		if !isSkippedMachineArg(a, len(out)) {
			out = append(out, a)
		}
	}
	return out
}

func isSkippedMachineArg(arg string, count int) bool {
	if isOSHelpArg(arg) {
		return false
	}
	if strings.HasPrefix(arg, "-") {
		return true
	}
	return count >= 2 && (strings.EqualFold(arg, "y") || strings.EqualFold(arg, "yes"))
}

func resolveLocalMachineIdentity() MachineIdentity {
	ip := resolveLocalIPv4()
	osHost, _ := os.Hostname()
	if len(osHost) == 0 {
		osHost = constants.Localhost
	}
	machName := readGlobalOrFallback("machine.name", osHost)
	alias := readGlobalOrFallback("machine.alias", ip)
	prevName, _ := config.GetVariable("global", "machine.previous_name")
	prevAlias, _ := config.GetVariable("global", "machine.previous_alias")
	osInfo := GetLocalOSInfo()
	user := os.Getenv("USERNAME")
	if user == "" {
		user = os.Getenv("USER")
	}
	detailedOS := probeDetailedOS()
	return MachineIdentity{
		Sequence: 1, NodeID: "local-01", IPAddress: ip, Alias: alias,
		MachineName: machName, OSHostname: osHost, PreviousName: prevName,
		PreviousAlias: prevAlias, OSPlatform: formatOSPlatform(runtime.GOOS),
		OSVersion: osInfo.OSVersion, Architecture: osInfo.Architecture,
		CPUCores: osInfo.NumCPU, CurrentUser: user, GitMapVersion: constants.Version,
		Scope: "local", OSVersionDetails: &detailedOS,
	}
}

func readGlobalOrFallback(key, fallback string) string {
	val, isFound := config.GetVariable("global", key)
	if isFound && len(strings.TrimSpace(val)) > 0 {
		return strings.TrimSpace(val)
	}
	return fallback
}

func resolveLocalIPv4() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, addr := range addrs {
		ip := extractNonLoopbackIPv4(addr)
		if len(ip) > 0 {
			return ip
		}
	}
	return "127.0.0.1"
}

func extractNonLoopbackIPv4(addr net.Addr) string {
	ipNet, isIPNet := addr.(*net.IPNet)
	if !isIPNet || ipNet.IP.IsLoopback() {
		return ""
	}
	v4 := ipNet.IP.To4()
	if v4 == nil {
		return ""
	}
	return v4.String()
}

var osPlatformLabels = map[string]string{
	"windows": "Windows (windows)",
	"darwin":  "macOS (darwin)",
	"macos":   "macOS (darwin)",
	"linux":   "Ubuntu / Linux (linux)",
	"ubuntu":  "Ubuntu / Linux (linux)",
}

func formatOSPlatform(goos string) string {
	if label, ok := osPlatformLabels[strings.ToLower(goos)]; ok {
		return label
	}

	return goos
}

func executeMachineAliasList(mode string, isSSH, isJSON bool) error {
	if isSSH {
		return renderFleetMachineAliasList(mode, isJSON)
	}
	return renderSingleMachineIdentity(resolveLocalMachineIdentity(), isJSON)
}

func executeMachineAliasSet(mode string, args []string, isSSH, isJSON bool) error {
	fmt.Printf("  %sSample format: dev-win-01, ubuntu-node-02, mac-studio-01%s\n", constants.ColorYellow, constants.ColorReset)
	targetArgs, err := resolveMachineSetArgs(mode, args, isJSON)
	if err != nil {
		return err
	}

	if isSSH {
		return executeSSHMachineAliasSet(mode, targetArgs, isJSON)
	}

	newVal := sanitizeMachineToken(targetArgs[0])

	return applyLocalMachineOrAliasSet(mode, newVal, isJSON)
}

func applyLocalMachineOrAliasSet(mode, newVal string, isJSON bool) error {
	cur := resolveLocalMachineIdentity()
	if mode == "alias" {
		_ = config.SetVariable("global", "machine.previous_alias", cur.Alias)
		_ = config.SetVariable("global", "machine.alias", newVal)
	} else {
		_ = config.SetVariable("global", "machine.previous_name", cur.MachineName)
		_ = config.SetVariable("global", "machine.name", newVal)
		_ = config.SetVariable("global", "machine.previous_alias", cur.Alias)
		_ = config.SetVariable("global", "machine.alias", newVal)
		applyOSHostnameChange(newVal)
	}
	updated := resolveLocalMachineIdentity()
	fmt.Printf("%s✔ Updated %s to %q (previous saved for 'gitmap %s revert')%s\n", constants.ColorGreen, mode, newVal, mode, constants.ColorReset)
	return renderIfJSON(updated, isJSON)
}

func executeMachineAliasRevert(mode string, args []string, isSSH, isJSON bool) error {
	if isSSH {
		return executeSSHMachineAliasRevert(mode, args, isJSON)
	}
	cur := resolveLocalMachineIdentity()
	return applyLocalRevert(mode, cur, isJSON)
}

func applyLocalRevert(mode string, cur MachineIdentity, isJSON bool) error {
	prev := cur.PreviousName
	if mode == "alias" || len(prev) == 0 {
		prev = cur.PreviousAlias
	}
	if len(prev) == 0 {
		prev = cur.OSHostname
	}
	_ = config.SetVariable("global", "machine.name", prev)
	_ = config.SetVariable("global", "machine.alias", prev)
	if mode == "machine" {
		applyOSHostnameChange(prev)
	}
	fmt.Printf("%s✔ Reverted %s back to %q%s\n", constants.ColorGreen, mode, prev, constants.ColorReset)
	return renderIfJSON(resolveLocalMachineIdentity(), isJSON)
}

func sanitizeMachineToken(raw string) string {
	return strings.ReplaceAll(strings.TrimSpace(raw), " ", "-")
}

func applyOSHostnameChange(newName string) {
	if os.Getenv("GITMAP_SKIP_OS_HOSTNAME_EXEC") == "1" {
		return
	}
	switch runtime.GOOS {
	case "darwin":
		_ = exec.Command("scutil", "--set", "HostName", newName).Run()
	case "linux":
		_ = exec.Command("hostnamectl", "set-hostname", newName).Run()
	default:
		_ = os.Setenv("GITMAP_MACHINE_NAME", newName)
	}
}
