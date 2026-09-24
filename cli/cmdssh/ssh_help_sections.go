package cmdssh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

func buildSSHClusterSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Cluster & Remote Operations",
		Entries: []termhelp.CommandEntry{
			{Command: "join <sub>", Description: "Join, broadcast, or enroll cluster nodes", HasSubcommands: true},
			{Command: "nodes export-json [file]", Description: "Export SSH nodes to JSON (default: gitmap-ssh-nodes.json)"},
			{Command: "nodes import-json [file]", Description: "Import SSH nodes from JSON or --base64 payload"},
			{Command: "export-oneliner", Description: "Generate portable 1-line command to import all SSH nodes"},
			{Command: "deploy node-config (nc)", Description: "Deploy local SSH node config across fleet (--except id,ip,alias)"},
			{Command: "exec <cmd>", Description: "Execute command across remote SSH nodes"},
			{Command: "install [target]", Description: "Install or update GitMap on remote machine(s)"},
			{Command: "update [agm|target]", Description: "Update GitMap or AGM across remote fleet (--except id,ip,alias)"},
			{Command: "bootstrap [target]", Description: "Setup passwordless sudo and deploy keys"},
			{Command: "compare (matrix)", Description: "Display architecture matrix: SSH vs Cluster vs SC"},
		},
	}
}

func buildSSHSecuritySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Known Hosts & Security",
		Entries: []termhelp.CommandEntry{
			{Command: "known-hosts (kh)", Description: "Manage, list, trust, and sync known_hosts", HasSubcommands: true},
			{Command: "trust <target>", Description: "Auto-scan and trust host key in known_hosts"},
			{Command: "untrust <target>", Description: "Remove machine from known_hosts and database"},
			{Command: "fix-auth <target>", Description: "Deploy public key to remote authorized_keys"},
		},
	}
}

func buildSSHRecoverySection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "State & Recovery",
		Entries: []termhelp.CommandEntry{
			{Command: "rm <nodes|keys>", Description: "Remove node(s) or key(s) with confirmation"},
			{Command: "reset [-y]", Description: "Wipe all registered nodes and reset SSH state"},
			{Command: "history [n]", Description: "Show SSH task and join history (default: 100)"},
			{Command: "undo", Description: "Undo last destructive node deletion or reset"},
			{Command: "redo", Description: "Redo last undone SSH operation"},
			{Command: "restore <id>", Description: "Restore nodes from a specific snapshot ID"},
			{Command: "error-logs [flags]", Description: "Display diagnostic logs and trace from last failed SSH action"},
		},
	}
}

func buildSSHRemoteToolsSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "Remote Tools & AGY Fleet",
		Entries: []termhelp.CommandEntry{
			{Command: "agy <subcmd> [--except]", Description: "Run any AGY command across remote SSH fleet (alias: gitmap agy ssh)"},
			{Command: "code <args>", Description: "Open remote folder in VS Code via SSH Remote"},
		},
	}
}
