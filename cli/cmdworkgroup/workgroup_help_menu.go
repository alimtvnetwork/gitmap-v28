package cmdworkgroup

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderWorkgroupHelp displays the styled two-column Workgroup help menu.
func RenderWorkgroupHelp() {
	termout.RenderMenu(buildWorkgroupHelpMenu())
}

func buildWorkgroupHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title:       "GitMap Workgroup & Folder SSH Automation (gitmap workgroup)",
		UsageLines:  buildWorkgroupUsageLines(),
		Sections:    buildWorkgroupHelpSections(),
		FooterFlags: buildWorkgroupFooterFlags(),
		Tips:        buildWorkgroupTips(),
	}
}

func buildWorkgroupUsageLines() []string {
	return []string{
		"gitmap workgroup [command] [arguments...]",
		"gitmap wg <command>          # Short alias",
		"gitmap group <command>       # Grouping alias",
	}
}

func buildWorkgroupHelpSections() []termout.HelpSection {
	return []termout.HelpSection{
		buildWorkgroupCommandsSection(),
		buildCommitFormatterSection(),
	}
}

func buildWorkgroupCommandsSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "WORKGROUP COMMANDS",
		Entries: []termout.CommandEntry{
			{Command: "add <name> <dir>", Description: "Register named workgroup for directory root"},
			{Command: "bind <dir> <key>", Description: "Bind folder to specific SSH private key"},
			{Command: "unbind <dir>", Description: "Remove folder-specific SSH key binding"},
			{Command: "ls, list", Description: "Display all registered workgroups & bindings"},
			{Command: "show [name|dir]", Description: "Inspect detailed configuration & key fingerprint"},
			{Command: "test-ssh [name|dir]", Description: "Test remote SSH handshake against GitHub/GitLab"},
			{Command: "rm <name|id>", Description: "Safely remove workgroup registration"},
		},
	}
}

func buildCommitFormatterSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "COMMIT FORMATTER ENHANCEMENTS (cpf / cbf)",
		Entries: []termout.CommandEntry{
			{Command: "cpf \"[auth]: add jwt\"", Description: "Commit with Feature: [auth]: prefix"},
			{Command: "cbf \"[db]: fix lock\"", Description: "Commit with Bug: [db]: prefix (alias: cpb)"},
			{Command: "cpc \"[ci]: update go\"", Description: "Commit with Chore: [ci]: prefix"},
			{Command: "cpr \"[v1.0]: release\"", Description: "Commit with Release: [v1.0]: prefix"},
		},
	}
}

func buildWorkgroupFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--ssh-key <path>", Description: "Path to private SSH key for directory"},
		{Command: "--desc <text>", Description: "Human-readable description of workgroup"},
		{Command: "--user-name <name>", Description: "Custom git commit author user.name"},
		{Command: "--user-email <email>", Description: "Custom git commit author user.email"},
		{Command: "--default", Description: "Set workgroup as default active workspace"},
		{Command: "--host <hostname>", Description: "Remote SSH host for handshake test (default: github.com)"},
		{Command: "-j, --json", Description: "Output results in JSON format"},
		{Command: "-h, --help", Description: "Show this workgroup help menu"},
	}
}

func buildWorkgroupTips() []string {
	return []string{
		"Use 'gitmap wg add client-a D:\\work\\client-a --ssh-key ~/.ssh/id_client' for 1-step setup.",
		"Repositories inside child folders automatically inherit parent folder SSH keys.",
		"Use 'gitmap cbf \"[module]: fix\"' as a fast alias for bug fix commit and push.",
		"Use 'gitmap wg test-ssh' to verify SSH authentication before pushing.",
	}
}

func buildWorkgroupRow(r WorkgroupRecord) termout.Row {
	isDefStr := ""
	if r.IsDefault {
		isDefStr = "yes"
	}
	author := r.GitUserName
	if r.GitUserEmail != "" {
		author += " <" + r.GitUserEmail + ">"
	}

	return termout.Row{
		Cells: []string{
			fmt.Sprintf("%d", r.ID),
			r.Name,
			r.RootDirectory,
			r.SSHKeyPath,
			author,
			isDefStr,
		},
	}
}

func buildWorkgroupTableConfig(records []WorkgroupRecord) termout.TableConfig {
	rows := make([]termout.Row, 0, len(records))
	for _, rec := range records {
		rows = append(rows, buildWorkgroupRow(rec))
	}

	return termout.TableConfig{
		Columns: []termout.ColumnConfig{
			{Title: "ID", MinWidth: 4},
			{Title: "NAME", MinWidth: 14},
			{Title: "ROOT DIRECTORY", MinWidth: 28},
			{Title: "SSH KEY", MinWidth: 20},
			{Title: "GIT AUTHOR", MinWidth: 20},
			{Title: "DEFAULT", MinWidth: 8},
		},
		Rows:       rows,
		HasBorders: true,
	}
}

// RenderWorkgroupTable prints a formatted table of all workgroups.
func RenderWorkgroupTable(records []WorkgroupRecord) {
	if len(records) == 0 {
		fmt.Println("No workgroups registered. Run 'gitmap wg add <name> <dir>' to create one.")

		return
	}
	cfg := buildWorkgroupTableConfig(records)
	termout.PrintTable(cfg)
}

// RenderWorkgroupDetailCard displays a detailed summary card for a workgroup.
func RenderWorkgroupDetailCard(rec WorkgroupRecord, eff EffectiveWorkgroupCredentials) {
	fmt.Printf("\n── Workgroup Details ──────────────────────────────────────\n")
	fmt.Printf("  ID             : %d\n", rec.ID)
	fmt.Printf("  Name           : %s\n", rec.Name)
	fmt.Printf("  Root Directory : %s\n", rec.RootDirectory)
	fmt.Printf("  SSH Key Path   : %s\n", rec.SSHKeyPath)
	fmt.Printf("  Description    : %s\n", rec.Description)
	fmt.Printf("  Git Author     : %s <%s>\n", rec.GitUserName, rec.GitUserEmail)
	fmt.Printf("  Is Default     : %t\n", rec.IsDefault)
	fmt.Printf("  Inherited Path : %t\n", eff.IsInherited)
	fmt.Printf("───────────────────────────────────────────────────────────\n\n")
}
