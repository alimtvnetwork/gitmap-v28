package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termhelp"
)

// RenderUserHelp displays the styled two-column User help menu.
func RenderUserHelp() {
	termhelp.RenderMenu(buildUserHelpMenu())
}

func buildUserHelpMenu() termhelp.HelpMenu {
	return termhelp.HelpMenu{
		Title: "Host OS User & SSH Credential Management (gitmap user)",
		UsageLines: []string{
			"gitmap user <subcommand> [flags]",
			"gitmap user add <username> [--password <pwd>]",
			"gitmap user add-ssh-key <username> <pubkey-path>",
		},
		Sections: []termhelp.HelpSection{
			buildUserAccountSection(),
			buildUserSSHAuthSection(),
		},
		FooterFlags: buildUserFooterFlags(),
		Tips: []string{
			"Provision system users and deploy SSH authorization keys securely.",
			"Use 'create-root' to set up administrative automation users on remote Linux nodes.",
		},
	}
}

func buildUserFooterFlags() []termhelp.CommandEntry {
	return []termhelp.CommandEntry{
		{Command: "--password <pwd>", Description: "Specify account password for newly created user"},
		{Command: "-k, --key <path>", Description: "Path to public key file for SSH authorization"},
		{Command: "-f, --force", Description: "Force user modification or process termination"},
		{Command: "-j, --json", Description: "Output results in structured JSON format"},
		{Command: "-h, --help", Description: "Show this user help menu"},
	}
}

func buildUserAccountSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "User Account Actions",
		Entries: []termhelp.CommandEntry{
			{Command: "add <username>", Description: "Create new operating system user with home directory"},
			{Command: "rm / delete <user>", Description: "Remove operating system user account and home directory"},
			{Command: "create-root <user>", Description: "Create administrative user with passwordless sudo access"},
			{Command: "kill <username>", Description: "Terminate all running processes owned by user"},
		},
	}
}

func buildUserSSHAuthSection() termhelp.HelpSection {
	return termhelp.HelpSection{
		Title: "SSH Credential Deployment",
		Entries: []termhelp.CommandEntry{
			{Command: "add-ssh-key <u...> <key>", Description: "Install public SSH key into target user's authorized_keys"},
		},
	}
}
