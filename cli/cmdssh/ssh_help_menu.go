package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderSSHHelp displays the styled two-column SSH help menu.
func RenderSSHHelp() {
	termhelp.RenderMenu(buildSSHHelpMenu())
}

func buildSSHHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "SSH & Multi-Node Cluster Fleet (gitmap ssh)",
		UsageLines: []string{
			"gitmap ssh [command] [args]",
			"gitmap sj [command] [args]",
		},
		Sections: []termhelp.HelpSection{
			buildSSHKeySection(),
			buildSSHNodeSection(),
			buildSSHClusterSection(),
			buildSSHSecuritySection(),
			buildSSHRecoverySection(),
			buildSSHRemoteToolsSection(),
		},
		FooterFlags: []termhelp.CommandEntry{
			{Command: "-h, --help", Description: "Show this SSH help menu"},
			{Command: "-y, --yes", Description: "Bypass confirmation prompts for rm/reset"},
		},
		Tips: []string{
			"Run 'gitmap ssh <command> --help' for details on any subcommand.",
			"Use 'gitmap ssh join user@ip alias' to enroll a remote machine.",
		},
	}
}

func buildSSHKeySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Key Management",
		Entries: []termhelp.CommandEntry{
			{Command: "create [name]", Description: "Generate or reuse an SSH key pair"},
			{Command: "list (ls)", Description: "List all managed SSH keys"},
			{Command: "copy (cp) [name]", Description: "Copy public key to clipboard"},
			{Command: "cat (view) [name]", Description: "Print public key to terminal"},
			{Command: "delete (rm) <name>", Description: "Remove a key from GitMap and disk"},
			{Command: "config", Description: "Rebuild ~/.ssh/config for all managed keys"},
		},
	}
}

func buildSSHNodeSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Node & Machine Management",
		Entries: []termhelp.CommandEntry{
			{Command: "nodes (ls)", Description: "List all registered SSH nodes and machine status"},
			{Command: "login <user@ip>", Description: "Connect & install environment on remote host"},
			{Command: "check [target]", Description: "Probe connectivity, port 22, and health"},
			{Command: "alias <sub>", Description: "Manage custom SSH host aliases", HasSubcommands: true},
			{Command: "scan", Description: "Probe reachability across registered SSH fleet"},
			{Command: "status (st)", Description: "Check ssh-agent and connection status"},
		},
	}
}
