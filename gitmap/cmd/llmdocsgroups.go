// Package cmd — llmdocsgroups.go defines command groups for LLM doc generation.
package cmd

// buildScanningGroup returns the scanning commands group.
func buildScanningGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Scanning & Discovery",
		commands: []llmCmdEntry{
			{"scan", "s", "Recursively scan a directory tree for Git repos", "gitmap scan ~/projects"},
			{"rescan", "rsc", "Re-scan previously scanned directories using cached config", "gitmap rescan"},
			{"rescan-subtree", "rss", "Deep-rescan a single subtree (e.g. an at-cap row's absolutePath) in one step", "gitmap rescan-subtree /home/me/work/monorepo"},
			{"list", "ls", "Show all tracked repos (filterable by type)", "gitmap ls go"},
			{"go-repos", "gr", "List Go projects detected by go.mod", "gitmap go-repos --json"},
			{"node-repos", "nr", "List Node.js projects detected by package.json", "gitmap node-repos"},
			{"react-repos", "rr", "List React projects", "gitmap react-repos"},
			{"cpp-repos", "cr", "List C++ projects", "gitmap cpp-repos"},
			{"csharp-repos", "csr", "List C# projects", "gitmap csharp-repos"},
		},
	}
}

// buildCloningGroup returns the cloning commands group.
func buildCloningGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Cloning",
		commands: []llmCmdEntry{
			{"clone", "c", "Clone repos from a scan output file (JSON/CSV/text)", "gitmap clone json --target-dir ./restored"},
			{"clone-next", "cn", "Clone the next versioned iteration of the current repo", "gitmap cn v++ --delete"},
			{"desktop-sync", "ds", "Register tracked repos with GitHub Desktop", "gitmap desktop-sync"},
		},
	}
}

// buildGitOpsGroup returns the git operations commands group.
func buildGitOpsGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Git Operations",
		commands: []llmCmdEntry{
			{"pull", "p", "Pull a specific repo by name (or all in group)", "gitmap pull my-api --group work"},
			{"exec", "x", "Run any git command across all tracked repos", "gitmap exec fetch --prune"},
			{"status", "st", "Show dirty/clean status for all repos", "gitmap status --group work"},
			{"watch", "w", "Live-refresh status dashboard", "gitmap watch --interval 10 --group work"},
			{"has-any-updates", "hau", "Check if remote has commits you haven't pulled", "gitmap hau"},
			{"latest-branch", "lb", "Find most recently updated remote branch", "gitmap lb --top 5 --json"},
		},
	}
}

// buildNavigationGroup returns the navigation commands group.
func buildNavigationGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Navigation & Groups",
		commands: []llmCmdEntry{
			{"cd", "go", "Navigate to a repo by slug name", "gitmap cd my-api"},
			{"group", "g", "Create/manage repo groups and activate for batch ops", "gitmap g work"},
			{"multi-group", "mg", "Select multiple groups for batch operations", "gitmap mg backend,frontend"},
			{"alias", "a", "Assign short names to repos (used with -A flag)", "gitmap alias set api my-api-gateway"},
			{"diff-profiles", "dp", "Compare repos across two scan profiles", "gitmap diff-profiles dev prod"},
		},
	}
}

// buildReleaseGroup returns the release workflow commands group.
func buildReleaseGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Release & Versioning",
		commands: []llmCmdEntry{
			{"release", "r", "Create release: branch, tag, push, cross-compile binaries", "gitmap release --bump patch --bin --compress --checksums"},
			{"release-self", "rs", "Release gitmap itself from any directory", "gitmap release-self --bump minor"},
			{"release-branch", "rb", "Create a release branch without tagging", "gitmap release-branch v2.50.0"},
			{"temp-release", "tr", "Create lightweight temp release branches", "gitmap tr 10 v1.$$ -s 5"},
			{"pull-release", "pr", "git pull (ff/rebase/merge), then release current repo", "gitmap pull-release --rebase v1.4.0"},
		},
	}
}

// buildReleaseInfoGroup returns the release info commands group.
func buildReleaseInfoGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Release Info",
		commands: []llmCmdEntry{
			{"changelog", "cl", "View changelog entries", "gitmap changelog --latest"},
			{"changelog-generate", "cg", "Auto-generate changelog from commit messages", "gitmap changelog-generate --write"},
			{"list-versions", "lv", "List all Git release tags", "gitmap list-versions --limit 10 --json"},
			{"list-releases", "lr", "List release metadata from database", "gitmap list-releases --json"},
			{"release-pending", "rp", "Show unreleased commits since last tag", "gitmap release-pending"},
			{"revert", "—", "Revert to a specific release version", "gitmap revert v2.48.0"},
			{"clear-release-json", "crj", "Remove orphaned release metadata files", "gitmap clear-release-json"},
			{"prune", "prn", "Delete stale release branches that have been tagged", "gitmap prune"},
		},
	}
}

// buildDataGroup returns the data/profile/bookmark commands group.
func buildDataGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Data, Profiles & Bookmarks",
		commands: []llmCmdEntry{
			{"export", "ex", "Export database to file", "gitmap export"},
			{"import", "im", "Import repos from file", "gitmap import gitmap-export.json"},
			{"profile", "pf", "Manage database profiles (create, switch, list)", "gitmap profile create work"},
			{"bookmark", "bk", "Save commands as named bookmarks and replay them", "gitmap bookmark save daily scan ~/projects"},
			{"rm", "remove, del", "Untrack repos from the database without deleting files", "gitmap rm my-repo"},
			{"db-reset", "—", "Reset the SQLite database (requires --confirm)", "gitmap db-reset --confirm"},
		},
	}
}

// buildHistoryGroup returns the history and stats commands group.
func buildHistoryGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "History & Analytics",
		commands: []llmCmdEntry{
			{"history", "hi", "Show CLI command execution history", "gitmap history --limit 10"},
			{"history-reset", "hr", "Clear command history", "gitmap history-reset"},
			{"stats", "ss", "Show aggregated usage and performance metrics", "gitmap stats --json"},
			{"dashboard", "db", "Generate interactive HTML dashboard for a repo", "gitmap dashboard"},
		},
	}
}

// buildAmendGroup returns the amend commands group.
func buildAmendGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Author Amendment",
		commands: []llmCmdEntry{
			{"amend", "am", "Rewrite commit author name/email across commits", "gitmap amend --name \"John\" --email \"john@co.com\" --force-push"},
			{"amend-list", "al", "List previous amendments", "gitmap amend-list --json"},
		},
	}
}

// buildProjectGroup returns the project detection commands group.
func buildProjectGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Project Detection",
		commands: []llmCmdEntry{
			{"go-repos", "gr", "List Go projects with metadata", "gitmap go-repos --json"},
			{"node-repos", "nr", "List Node.js projects", "gitmap node-repos"},
			{"react-repos", "rr", "List React projects", "gitmap react-repos"},
			{"cpp-repos", "cr", "List C++ projects", "gitmap cpp-repos"},
			{"csharp-repos", "csr", "List C# projects", "gitmap csharp-repos"},
		},
	}
}

// buildSSHGroup returns the SSH key management commands group.
func buildSSHGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "SSH Key Management",
		commands: []llmCmdEntry{
			{"ssh", "—", "Generate, list, display, and delete SSH keys for Git auth", "gitmap ssh --name work --path ~/.ssh/id_rsa_work"},
		},
	}
}

// buildZipGroup returns the zip group commands group.
func buildZipGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Zip Groups (Release Archives)",
		commands: []llmCmdEntry{
			{"zip-group", "z", "Manage named collections of files bundled into ZIP during releases", "gitmap z create docs-bundle ./README.md ./docs/"},
		},
	}
}

// buildEnvToolsGroup returns the env and install commands group.
func buildEnvToolsGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Environment & Tools",
		commands: []llmCmdEntry{
			{"env", "ev", "Manage persistent environment variables and PATH entries", "gitmap env set GOPATH /home/user/go"},
			{"install", "in", "Install developer tools via platform package manager", "gitmap install node"},
			{"uninstall", "un", "Uninstall a developer tool", "gitmap uninstall node"},
		},
	}
}

// buildTaskGroup returns the task commands group.
func buildTaskGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "File-Sync Tasks",
		commands: []llmCmdEntry{
			{"task", "tk", "Create and run one-way folder synchronization tasks", "gitmap task create my-sync --src ./src --dest ./backup"},
			{"pending", "—", "Show pending tasks awaiting execution", "gitmap pending"},
			{"do-pending", "dp", "Execute pending tasks", "gitmap do-pending"},
		},
	}
}

// buildUtilityGroup returns the utility commands group.
func buildUtilityGroup() llmCmdGroup {
	return llmCmdGroup{
		title: "Utilities",
		commands: []llmCmdEntry{
			{"setup", "—", "Interactive first-time configuration wizard", "gitmap setup"},
			{"doctor", "—", "Diagnose PATH, deploy, and version issues", "gitmap doctor --fix-path"},
			{"update", "—", "Self-update from source repo or via gitmap-updater", "gitmap update"},
			{"version", "v", "Show version number", "gitmap version"},
			{"completion", "cmp", "Generate shell tab-completion (PowerShell, Bash, Zsh)", "gitmap completion"},
			{"interactive", "i", "Full-screen interactive TUI", "gitmap interactive"},
			{"docs", "d", "Open documentation website in browser", "gitmap docs"},
			{"help-dashboard", "hd", "Serve the docs site locally in your browser", "gitmap hd"},
			{"gomod", "gm", "Rename Go module path across repo", "gitmap gomod old/path new/path"},
			{"seo-write", "sw", "Auto-commit SEO messages from CSV", "gitmap seo-write --csv data.csv"},
			{"llm-docs", "ld", "Generate this LLM.md reference file", "gitmap llm-docs"},
		},
	}
}
