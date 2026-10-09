package cmduser

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/termout"
)

// RenderUserHelp displays the styled two-column User help menu.
func RenderUserHelp() {
	termout.RenderMenu(buildUserHelpMenu())
}

func buildUserHelpMenu() termout.HelpMenu {
	return termout.HelpMenu{
		Title: "Git Identity, Multi-User Profiles & System Accounts (gitmap user)",
		UsageLines: []string{
			"gitmap user <subcommand> [flags]",
			"gitmap user info [--json]",
			"gitmap user switch <alias> [--global|--project]",
			"gitmap user list",
			"gitmap user project bind <alias>",
		},
		Sections: []termout.HelpSection{
			buildGitUserSection(),
			buildUserAccountSection(),
			buildUserSSHAuthSection(),
		},
		FooterFlags: buildUserFooterFlags(),
		Tips: []string{
			"Run 'gitmap user info' to inspect both GitHub CLI auth state and Git author configs.",
			"Switch profiles effortlessly per project or machine using 'switch' and 'project bind'.",
			"Provision system users and deploy SSH authorization keys securely.",
		},
	}
}

func buildGitUserSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Git Profile & Context Actions",
		Entries: []termout.CommandEntry{
			{Command: "info / status", Description: "Display consolidated Git and GitHub CLI user identity card"},
			{Command: "list / ls", Description: "List configured Git profiles with active and bound indicators"},
			{Command: "switch / use <alias>", Description: "Switch active Git profile (--global or --project)"},
			{Command: "project [bind|unbind]", Description: "Manage per-repository profile bindings"},
			{Command: "config [global|local]", Description: "Inspect or configure Git user.name and user.email"},
			{Command: "sync", Description: "Apply bound project user profile to current repository"},
		},
	}
}

func buildUserFooterFlags() []termout.CommandEntry {
	return []termout.CommandEntry{
		{Command: "--global", Description: "Apply profile switch or configuration globally"},
		{Command: "--project", Description: "Bind profile and apply configuration to current repository"},
		{Command: "--name <string>", Description: "Specify author name for Git profile or configuration"},
		{Command: "--email <string>", Description: "Specify author email for Git profile or configuration"},
		{Command: "-j, --json", Description: "Output results in structured JSON format"},
		{Command: "-h, --help", Description: "Show this user help menu"},
	}
}

func buildUserAccountSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "Operating System User Accounts",
		Entries: []termout.CommandEntry{
			{Command: "add <username>", Description: "Create new operating system user with home directory"},
			{Command: "rm / delete <user>", Description: "Remove operating system user account and home directory"},
			{Command: "create-root <user>", Description: "Create administrative user with passwordless sudo access"},
			{Command: "kill <username>", Description: "Terminate all running processes owned by user"},
		},
	}
}

func buildUserSSHAuthSection() termout.HelpSection {
	return termout.HelpSection{
		Title: "SSH Credential Deployment",
		Entries: []termout.CommandEntry{
			{Command: "add-ssh-key <u...> <key>", Description: "Install public SSH key into target user's authorized_keys"},
		},
	}
}
