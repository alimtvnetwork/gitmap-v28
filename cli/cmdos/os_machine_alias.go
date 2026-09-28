// Package cmdos — os_machine_alias.go provides cross-platform machine identity, OS hostname, and fleet alias management.
package cmdos

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// MachineIdentity represents the resolved identity of a local or remote SSH fleet machine.
type MachineIdentity struct {
	Sequence      int    `json:"sequence"`
	NodeID        string `json:"nodeId"`
	IPAddress     string `json:"ipAddress"`
	Alias         string `json:"alias"`
	MachineName   string `json:"machineName"`
	OSHostname    string `json:"osHostname"`
	PreviousName  string `json:"previousName,omitempty"`
	PreviousAlias string `json:"previousAlias,omitempty"`
	OSPlatform    string `json:"osPlatform"`
	Scope         string `json:"scope"`
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
	if len(clean) == 0 || isOSHelpArg(clean[0]) {
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
		if !strings.HasPrefix(a, "-") {
			out = append(out, a)
		}
	}
	return out
}

func resolveLocalMachineIdentity() MachineIdentity {
	ip := resolveLocalIPv4()
	osHost, _ := os.Hostname()
	if len(osHost) == 0 {
		osHost = "localhost"
	}
	machName := readGlobalOrFallback("machine.name", osHost)
	alias := readGlobalOrFallback("machine.alias", ip)
	prevName, _ := config.GetVariable("global", "machine.previous_name")
	prevAlias, _ := config.GetVariable("global", "machine.previous_alias")
	return MachineIdentity{
		Sequence: 1, NodeID: "local-01", IPAddress: ip, Alias: alias,
		MachineName: machName, OSHostname: osHost, PreviousName: prevName,
		PreviousAlias: prevAlias, OSPlatform: formatOSPlatform(runtime.GOOS), Scope: "local",
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

func formatOSPlatform(goos string) string {
	switch strings.ToLower(goos) {
	case "windows":
		return "Windows (windows)"
	case "darwin", "macos":
		return "macOS (darwin)"
	case "linux", "ubuntu":
		return "Ubuntu / Linux (linux)"
	default:
		return goos
	}
}

func executeMachineAliasList(mode string, isSSH, isJSON bool) error {
	if isSSH {
		return renderFleetMachineAliasList(mode, isJSON)
	}
	local := resolveLocalMachineIdentity()
	if isJSON {
		return printIdentityJSON(local)
	}
	renderSingleMachineIdentity(mode, local)
	return nil
}

func renderSingleMachineIdentity(mode string, id MachineIdentity) {
	title := "Machine Identity & OS Hostname"
	if mode == "alias" {
		title = "Machine Alias & Network Identifier"
	}
	fmt.Printf("%s● %s%s\n", constants.ColorCyan, title, constants.ColorReset)
	fmt.Printf("  Machine IP:      %s%s%s\n", constants.ColorGreen, id.IPAddress, constants.ColorReset)
	fmt.Printf("  Machine Alias:   %s%s%s (auto-defaults to IP if unset)\n", constants.ColorYellow, id.Alias, constants.ColorReset)
	fmt.Printf("  Machine Name:    %s%s%s (OS Hostname: %s)\n", constants.ColorCyan, id.MachineName, constants.ColorReset, id.OSHostname)
	fmt.Printf("  Previous Value:  name=%q | alias=%q\n", id.PreviousName, id.PreviousAlias)
	fmt.Printf("  OS Platform:     %s\n", id.OSPlatform)
}

func renderFleetMachineAliasList(mode string, isJSON bool) error {
	items := collectSSHFluidIdentities()
	if isJSON {
		return printIdentityJSON(items)
	}
	fmt.Printf("%s● SSH Fleet & Local Machine Identifiers (%s --ssh)%s\n", constants.ColorCyan, mode, constants.ColorReset)
	fmt.Printf("  %-4s %-12s %-16s %-20s %-20s %-22s\n", "SEQ", "NODE ID", "IP ADDRESS", "ALIAS", "MACHINE NAME", "OS PLATFORM")
	fmt.Printf("  %s\n", strings.Repeat("─", 98))
	for _, it := range items {
		fmt.Printf("  #%-3d %-12s %-16s %-20s %-20s %-22s\n",
			it.Sequence, it.NodeID, it.IPAddress, it.Alias, it.MachineName, it.OSPlatform)
	}
	return nil
}

func collectSSHFluidIdentities() []MachineIdentity {
	items := []MachineIdentity{resolveLocalMachineIdentity()}
	rootDB, err := openDB()
	if err != nil || rootDB == nil {
		return items
	}
	defer rootDB.Close()
	res := db.GetSSHConnections(context.Background(), rootDB.Conn())
	if res.IsFailure() {
		return items
	}
	for idx, conn := range res.Data {
		items = append(items, convertSSHConnToIdentity(idx+2, conn))
	}
	return items
}

func convertSSHConnToIdentity(seq int, conn db.SSHConnection) MachineIdentity {
	alias := strings.TrimSpace(conn.Alias)
	if len(alias) == 0 {
		alias = conn.IPAddress
	}
	nodeID := fmt.Sprintf("ssh-%02d", seq-1)
	prev, _ := config.GetVariable("ssh_prev", conn.IPAddress)
	return MachineIdentity{
		Sequence: seq, NodeID: nodeID, IPAddress: conn.IPAddress, Alias: alias,
		MachineName: alias, OSHostname: alias, PreviousAlias: prev,
		OSPlatform: formatOSPlatform(conn.OS), Scope: "ssh",
	}
}

func executeMachineAliasSet(mode string, args []string, isSSH, isJSON bool) error {
	fmt.Printf("  %sSample format: dev-win-01, ubuntu-node-02, mac-studio-01%s\n", constants.ColorYellow, constants.ColorReset)
	if len(args) == 0 {
		return apperror.NewSimple("missing target name/alias (e.g. gitmap "+mode+" set dev-win-01 -y)", "E_MISSING_MACHINE_NAME")
	}
	if isSSH {
		return executeSSHMachineAliasSet(mode, args, isJSON)
	}
	newVal := sanitizeMachineToken(args[0])
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
	fmt.Printf("%s✔ Updated %s to %q (previous saved for 'gitmap %s revert')%s\n",
		constants.ColorGreen, mode, newVal, mode, constants.ColorReset)
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

func executeSSHMachineAliasSet(mode string, args []string, isJSON bool) error {
	target, newVal := resolveSSHTargetAndValue(args)
	rootDB, err := openDB()
	if err != nil || rootDB == nil {
		return apperror.NewSimple("unable to open SSH database", "E_SSH_DB")
	}
	defer rootDB.Close()
	updateSSHConnRecord(rootDB, target, newVal)
	fmt.Printf("%s✔ Updated SSH fleet %s for %q -> %q%s\n", constants.ColorGreen, mode, target, newVal, constants.ColorReset)
	return renderFleetIfJSON(isJSON)
}

func executeSSHMachineAliasRevert(mode string, args []string, isJSON bool) error {
	target := "local-01"
	if len(args) > 0 {
		target = args[0]
	}
	prev, hasPrev := config.GetVariable("ssh_prev", target)
	if !hasPrev || len(prev) == 0 {
		prev = target
	}
	return executeSSHMachineAliasSet(mode, []string{target, prev}, isJSON)
}

func resolveSSHTargetAndValue(args []string) (string, string) {
	if len(args) >= 2 {
		return args[0], sanitizeMachineToken(args[1])
	}
	return "local-01", sanitizeMachineToken(args[0])
}

func updateSSHConnRecord(rootDB *store.DB, target, newVal string) {
	ctx := context.Background()
	_ = config.SetVariable("ssh_prev", target, target)
	_, _ = rootDB.Conn().ExecContext(ctx, `UPDATE SSHConnection SET Alias = ? WHERE Alias = ? OR IPAddress = ?`, newVal, target, target)
	_, _ = rootDB.Conn().ExecContext(ctx, `UPDATE ssh_hosts SET alias = ? WHERE alias = ? OR ip = ?`, newVal, target, target)
}

func sanitizeMachineToken(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.ReplaceAll(cleaned, " ", "-")
	return cleaned
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

func renderIfJSON(v interface{}, isJSON bool) error {
	if !isJSON {
		return nil
	}
	return printIdentityJSON(v)
}

func renderFleetIfJSON(isJSON bool) error {
	if !isJSON {
		return nil
	}
	return printIdentityJSON(collectSSHFluidIdentities())
}

func printIdentityJSON(v interface{}) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal machine identity json")
	}
	fmt.Println(string(data))
	return nil
}

// RenderMachineAliasHelp renders the boxed help menu for gitmap machine and gitmap alias.
func RenderMachineAliasHelp(mode string) {
	title := "Machine Identity & Hostname Management (gitmap machine)"
	if mode == "alias" {
		title = "Machine Network Alias & Fleet Identifier (gitmap alias)"
	}
	termhelp.RenderMenu(termhelp.HelpMenu{
		Title: title,
		UsageLines: []string{
			fmt.Sprintf("gitmap %s [ls|set|change|revert|help] [--ssh] [-y] [--json]", mode),
			fmt.Sprintf("gitmap os %s [ls|set|change|revert|help] [--ssh] [-y] [--json]", mode),
		},
		Sections:    buildMachineAliasHelpSections(mode),
		FooterFlags: buildMachineAliasFooterFlags(),
		Tips: []string{
			"If alias is not explicitly set, GitMap automatically defaults to the machine's Local IPv4.",
			fmt.Sprintf("Use 'gitmap %s revert' to restore the previous machine name or alias at any time.", mode),
		},
	})
}

func buildMachineAliasHelpSections(mode string) []termhelp.HelpSection {
	return []termhelp.HelpSection{
		{
			Title: "Subcommands",
			Entries: []termhelp.CommandEntry{
				{Command: "ls (list, status)", Description: "Display Machine IP, Alias, OS Hostname, Previous Value, and Platform"},
				{Command: "set <name> [-y]", Description: "Set machine name/alias (sample: dev-win-01, ubuntu-node-02)"},
				{Command: "change <name> [-y]", Description: "Change machine name/alias and record previous value for revert"},
				{Command: "revert (undo)", Description: "Restore previous machine name or alias from snapshot"},
				{Command: "help", Description: "Show boxed help menu and current machine identity summary"},
			},
		},
		{
			Title: "SSH Fleet Delegation (--ssh)",
			Entries: []termhelp.CommandEntry{
				{Command: fmt.Sprintf("gitmap %s ls --ssh", mode), Description: "Inspect IP, Alias, Hostname, and OS across all joined SSH fleet nodes"},
				{Command: fmt.Sprintf("gitmap %s set <node> <name> --ssh", mode), Description: "Update remote SSH machine alias/name in fleet registry"},
				{Command: fmt.Sprintf("gitmap %s revert [node] --ssh", mode), Description: "Revert remote SSH machine alias/name to previous value"},
			},
		},
	}
}

func buildMachineAliasFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "-s, --ssh", Description: "Operate across all joined SSH fleet machines and local node"},
		{Command: "-y, --yes", Description: "Apply name/alias changes automatically without confirmation"},
		{Command: "-j, --json", Description: "Output machine identity or SSH fleet table as structured JSON"},
		{Command: "-h, --help", Description: "Display help and current machine summary"},
	}
}
