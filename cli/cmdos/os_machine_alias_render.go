package cmdos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

func renderSingleMachineIdentity(id MachineIdentity, isJSON bool) error {
	if isJSON {
		return printIdentityJSON([]MachineIdentity{id})
	}
	fmt.Printf("%s● Machine Identity:%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • %-17s%s\n", "Alias:", id.Alias)
	fmt.Printf("  • %-17s%s\n", "Hostname:", id.MachineName)
	fmt.Printf("  • %-17s%s\n", "IP Address:", id.IPAddress)
	fmt.Printf("  • %-17s%s\n", "OS Platform:", id.OSPlatform)
	renderDetailedOSVersionTree(id.OSVersionDetails, id.OSVersion)
	archDesc := formatArchWithCores(id.Architecture, id.CPUCores)
	fmt.Printf("  • %-17s%s\n", "Architecture:", archDesc)
	fmt.Printf("  • %-17s%s\n", "User:", id.CurrentUser)
	fmt.Printf("  • %-17s%s\n", "GitMap Version:", id.GitMapVersion)
	return nil
}

func renderDetailedOSVersionTree(det *OSVersionInfo, fallback string) {
	if det == nil {
		fmt.Printf("  • %-17s%s\n", "OS Version:", fallback)
		return
	}
	fmt.Printf("  • %-17s\n", "OS Version:")
	fmt.Printf("    ├─ %-16s%s\n", "Type:", formatOSType(det.Type))
	fmt.Printf("    ├─ %-16s%s\n", "Edition:", det.Edition)
	fmt.Printf("    ├─ %-16s%s\n", "Release:", det.Release)
	buildInfo := fmt.Sprintf("%d", det.Build)
	if det.KernelVersion != "" {
		buildInfo = fmt.Sprintf("%d (Kernel: %s)", det.Build, det.KernelVersion)
	}
	fmt.Printf("    ├─ %-16s%s\n", "Build:", buildInfo)
	fmt.Printf("    ├─ %-16s%s\n", "Windows Server:", formatYesNo(det.IsWindowsServer))
	fmt.Printf("    ├─ %-16s%s\n", "Windows Pro:", formatYesNo(det.IsWindowsPro))
	fmt.Printf("    └─ %-16s%s\n", "Windows Home:", formatYesNo(det.IsWindowsHome))
}

func formatOSType(t string) string {
	if strings.EqualFold(t, "windows") {
		return "Windows"
	}
	if strings.EqualFold(t, "darwin") || strings.EqualFold(t, "macos") {
		return "macOS"
	}
	if strings.EqualFold(t, "linux") {
		return "Linux"
	}
	return t
}

func formatYesNo(isEnabled bool) string {
	if isEnabled {
		return "Yes"
	}
	return "No"
}

func renderFleetMachineAliasList(mode string, isJSON bool) error {
	items := collectSSHFluidIdentities()
	if isJSON {
		return printIdentityJSON(items)
	}
	fmt.Printf("%s● SSH Fleet & Local Machine Identifiers (%s --ssh)%s\n", constants.ColorCyan, mode, constants.ColorReset)
	if hasFleetGitMapVersion(items) {
		return renderFleetVersionTable(items)
	}
	return renderFleetStandardTable(items)
}

func hasFleetGitMapVersion(items []MachineIdentity) bool {
	for _, it := range items {
		if it.GitMapVersion != "" {
			return true
		}
	}
	return false
}

func renderFleetVersionTable(items []MachineIdentity) error {
	fmt.Printf("  %-4s %-12s %-16s %-18s %-18s %-20s %-14s %s\n", "SEQ", "NODE ID", "IP ADDRESS", "ALIAS", "USER", "OS PLATFORM", "OS VERSION", "GITMAP VERSION")
	fmt.Printf("  %s\n", strings.Repeat("─", 125))
	for _, it := range items {
		fmt.Printf("  #%-3d %-12s %-16s %-18s %-18s %-20s %-14s %s\n",
			it.Sequence, it.NodeID, it.IPAddress, it.Alias, fallbackDash(it.CurrentUser), it.OSPlatform, fallbackDash(it.OSVersion), fallbackDash(it.GitMapVersion))
	}
	return nil
}

func renderFleetStandardTable(items []MachineIdentity) error {
	fmt.Printf("  %-4s %-12s %-16s %-18s %-18s %-20s %s\n", "SEQ", "NODE ID", "IP ADDRESS", "ALIAS", "USER", "OS PLATFORM", "OS VERSION")
	fmt.Printf("  %s\n", strings.Repeat("─", 110))
	for _, it := range items {
		fmt.Printf("  #%-3d %-12s %-16s %-18s %-18s %-20s %s\n",
			it.Sequence, it.NodeID, it.IPAddress, it.Alias, fallbackDash(it.CurrentUser), it.OSPlatform, fallbackDash(it.OSVersion))
	}
	return nil
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
	termout.RenderMenu(termout.HelpMenu{
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

func buildMachineAliasHelpSections(mode string) []termout.HelpSection {
	return []termout.HelpSection{
		{
			Title: "Subcommands",
			Entries: []termout.CommandEntry{
				{Command: "ls (list, status)", Description: "Display Machine IP, Alias, OS Hostname, Previous Value, and Platform"},
				{Command: "set <name> [-y]", Description: "Set machine name/alias (sample: dev-win-01, ubuntu-node-02)"},
				{Command: "change <name> [-y]", Description: "Change machine name/alias and record previous value for revert"},
				{Command: "revert (undo)", Description: "Restore previous machine name or alias from snapshot"},
				{Command: "help", Description: "Show boxed help menu and current machine identity summary"},
			},
		},
		{
			Title: "SSH Fleet Delegation (--ssh)",
			Entries: []termout.CommandEntry{
				{Command: fmt.Sprintf("gitmap %s ls --ssh", mode), Description: "Inspect IP, Alias, Hostname, and OS across all joined SSH fleet nodes"},
				{Command: fmt.Sprintf("gitmap %s set <node> <name> --ssh", mode), Description: "Update remote SSH machine alias/name in fleet registry"},
				{Command: fmt.Sprintf("gitmap %s revert [node] --ssh", mode), Description: "Revert remote SSH machine alias/name to previous value"},
			},
		},
	}
}

func buildMachineAliasFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "-s, --ssh", Description: "Operate across all joined SSH fleet machines and local node"},
		{Command: "-y, --yes", Description: "Apply name/alias changes automatically without confirmation"},
		{Command: "-j, --json", Description: "Output machine identity or SSH fleet table as structured JSON"},
		{Command: "-h, --help", Description: "Display help and current machine summary"},
	}
}

func formatArchWithCores(arch string, cores int) string {
	if cores <= 0 {
		return arch
	}
	if arch != "" {
		return fmt.Sprintf("%s (%d CPUs)", arch, cores)
	}
	return fmt.Sprintf("%d CPUs", cores)
}

func fallbackDash(val string) string {
	if val == "" {
		return "-"
	}
	return val
}
