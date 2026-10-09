package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderSSHHelp displays the styled two-column SSH help menu.
func RenderSSHHelp() {
	termout.RenderMenu(buildSSHHelpMenu())
}

func buildSSHHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "SSH & Multi-Node Cluster Fleet (gitmap ssh)",
		UsageLines: []string{
			"gitmap ssh [command] [args]",
			"gitmap sj [command] [args]",
		},
		Sections: []termout.HelpSection{
			buildSSHKeySection(),
			buildSSHDaemonSection(),
			buildSSHNodeSection(),
			buildSSHClusterSection(),
			buildSSHSecuritySection(),
			buildSSHRecoverySection(),
			buildSSHRemoteToolsSection(),
		},
		FooterFlags: []termout.CommandEntry{
			{Command: "-h, --help", Description: "Show this SSH help menu"},
			{Command: "-y, --yes", Description: "Bypass confirmation prompts for rm/reset/create"},
		},
		Tips: []string{
			"Run 'gitmap ssh <command> --help' for details on any subcommand.",
			"Use 'gitmap ssh view' to inspect public key and auto-copy to clipboard.",
			"Use 'gitmap ssh join user@ip alias' to enroll a remote machine.",
		},
	}
}

func buildSSHKeySection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Key Management",
		Entries: []termout.CommandEntry{
			{Command: "view (v, show) [name]", Description: "Display SSH public key card, copy to clipboard & hint regeneration"},
			{Command: "create [name] [-y]", Description: "Generate SSH key pair with overwrite guard, backup & undo/redo"},
			{Command: "copy (cp) [name]", Description: "Copy public key to clipboard"},
			{Command: "list (ls)", Description: "List all managed SSH keys"},
			{Command: "delete (rm) <name>", Description: "Remove a key from GitMap and disk"},
			{Command: "config", Description: "Rebuild ~/.ssh/config for all managed keys"},
		},
	}
}

func buildSSHDaemonSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Server Daemon, Firewall & Web UI",
		Entries: []termout.CommandEntry{
			{Command: "enable [--port <p>]", Description: "Enable OpenSSH Server daemon and start service"},
			{Command: "disable", Description: "Stop OpenSSH Server daemon and close firewall ports"},
			{Command: "port <ls|add|rm|set>", Description: "Manage SSH ports in sshd_config with cross-OS firewall sync"},
			{Command: "enable-public [port]", Description: "Open host firewall and bind 0.0.0.0 for public internet access"},
			{Command: "ui (web, dashboard)", Description: "Open interactive fleet management web dashboard in browser"},
			{Command: "troubleshoot (doctor)", Description: "Diagnose connection failures with 7-step guided remediation"},
		},
	}
}

func buildSSHNodeSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Node & Machine Management",
		Entries: []termout.CommandEntry{
			{Command: "nodes (ls)", Description: "List all registered SSH nodes and machine status"},
			{Command: "login <user@ip>", Description: "Connect & install environment on remote host"},
			{Command: "check [target]", Description: "Probe connectivity, port 22, and health"},
			{Command: "alias <sub>", Description: "Manage custom SSH host aliases", HasSubcommands: true},
			{Command: "scan", Description: "Probe reachability across registered SSH fleet"},
			{Command: "status (st)", Description: "Check ssh-agent and connection status"},
		},
	}
}
