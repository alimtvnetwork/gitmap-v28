package cmdpending

func catalogEntries6() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("os dock", "dock", "v6.517.0", "os", "Configure taskbar and dock alignment across Windows and Linux", `gitmap os dock bottom`),
		newCmd("os panel", "panel", "v6.518.0", "os", "Configure system display panels and multi-monitor geometry", `gitmap os panel --primary=1`),
		newCmd("apps", "app", "v6.516.0", "os", "Installed application auditor and uninstaller framework", `gitmap apps list`),
		newCmd("fix-link", "fl", "v6.512.0", "os", "Inspect and repair broken symlinks and VMware shared mounts", `gitmap fix-link --repair`),
		newCmd("vpm", "vmware-power", "v6.514.0", "os", "VMware guest shared folder mounting and open-vm-tools manager", `gitmap vpm mount`),
		newCmd("which-format", "wf", "v6.519.0", "os", "Inspect file line ending (CRLF/LF) and encoding formats", `gitmap which-format file.go`),
		newCmd("os-dns", "dns", "v6.513.0", "os", "Inspect, benchmark, and switch system DNS servers", `gitmap os-dns switch 1.1.1.1`),
		newCmd("os-theme", "theme", "v6.511.0", "os", "Switch desktop appearance between Dark Mode and Light Mode", `gitmap os-theme dark`),
		newCmd("power", "pwr", "v6.515.0", "os", "Screen timeout and power plan management with state tracking", `gitmap power sleep 30`),
		newCmd("autologin", "os-autologin", "v6.510.0", "os", "Configure OS auto-login credentials for Windows and Ubuntu", `gitmap autologin --user=admin`),
		newCmd("cd", "cd", "v6.527.0", "os", "Autonomous default workdir resolution and persistence", `gitmap cd`),
	}
}

func catalogEntries7() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("spec issue", "si", "v6.522.0", "spec", "Audit specification gap issues and generate remediation steps", `gitmap spec issue --spec=02-spec/21-app`),
		newCmd("spec next", "spec-next", "v6.520.0", "spec", "Issue next free spec number atomically via DB-backed counter", `gitmap spec next`),
		newCmd("spec-author", "sa", "v6.520.0", "spec", "Scaffold structured architecture and component specifications", `gitmap spec-author --slug=my-feature`),
		newCmd("spec-audit", "spa", "v6.519.0", "spec", "Conduct blind-AI readiness audits on specification documents", `gitmap spec-audit --target=02-spec`),
		newCmd("user", "usr", "v6.521.0", "spec", "Audit and configure git user name and email per repository", `gitmap user set --name="Dev" --email="dev@co.com"`),
		newCmd("changelog", "clog", "v6.508.0", "spec", "Generate SemVer changelog entries from recent atomic commits", `gitmap changelog generate`),
		newCmd("release-notes", "rn", "v6.509.0", "spec", "Extract markdown release notes for tagged versions", `gitmap release-notes v6.523.0`),
		newCmd("seowrite", "seo", "v6.510.0", "spec", "Format repository descriptions and keywords for SEO", `gitmap seowrite --keyword=git`),
		newCmd("help-json", "hj", "v6.515.0", "spec", "Emit machine-readable JSON schema for CLI help topics", `gitmap help --json --filter=task`),
		newCmd("prompt-template", "pt", "v6.512.0", "spec", "Manage reusable AI prompt prefix and verification templates", `gitmap prompt-template list`),
		newCmd("prompt list", "plst", "v6.518.0", "spec", "List registered canonical AI prompts in 01-prompts directory", `gitmap prompt list`),
		newCmd("summary", "summary", "v6.525.0", "spec", "Summary engine with full treeview and nodes delegation", `gitmap summary`),
	}
}

func catalogEntries8() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("rs file", "rs-f", "v6.515.0", "storage", "Copy secret file into encrypted repo-secrets and auto-push", `gitmap rs file .env --repo=backend`),
		newCmd("rs folder", "rs-d", "v6.515.0", "storage", "Copy secret folder into repo-secrets and auto-push", `gitmap rs folder certs/`),
		newCmd("rs text", "rs-t", "v6.515.0", "storage", "Store sensitive string directly into repo-secrets store", `gitmap rs text "key_123" --slug=api-key`),
		newCmd("rc file", "rc-f", "v6.516.0", "storage", "Offload reusable test script into repo-cache and auto-push", `gitmap rc file test_harness.py`),
		newCmd("rc folder", "rc-d", "v6.516.0", "storage", "Store shared fixture directory into repo-cache", `gitmap rc folder fixtures/`),
		newCmd("rc text", "rc-t", "v6.516.0", "storage", "Write reusable automation test script into repo-cache", `gitmap rc text "script" --slug=test --ext=.ps1`),
		newCmd("storage", "stg", "v6.511.0", "storage", "Inspect disk volumes, split-DB sqlite sizes, and snapshots", `gitmap storage`),
		newCmd("clean-dev", "cld", "v6.517.0", "storage", "Purge IDE temporary files, dev caches, and lock buffers", `gitmap clean-dev`),
		newCmd("clear-terminal", "cls-term", "v6.518.0", "storage", "Clear stuck terminal processes, handles, and buffer pipes", `gitmap clear-terminal`),
		newCmd("purge-cache", "pgc", "v6.520.0", "storage", "Invalidate and purge stale SQLite status cache entries", `gitmap purge-cache --ttl=90s`),
	}
}

func catalogEntries9() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("sync", "sy", "v6.522.0", "sync", "Synchronize prompts, skills, and specs across 43 repositories", `gitmap sync --workers=8`),
		newCmd("sync --repo", "sy-r", "v6.522.0", "sync", "Synchronize canonical assets to single target repository", `gitmap sync --repo=movie-cli-v8`),
		newCmd("sync --dry-run", "sy-d", "v6.522.0", "sync", "Preview synchronization file diffs without disk mutations", `gitmap sync --dry-run`),
		newCmd("sync --list", "sy-l", "v6.522.0", "sync", "List registered fleet repositories connected to sync engine", `gitmap sync --list`),
		newCmd("pull-all-efficient", "pae", "v6.508.0", "sync", "Pull all repositories concurrently with JSON telemetry", `gitmap pae --json`),
		newCmd("reconcile", "recon", "v6.505.0", "sync", "Reconcile divergent branch heads and resolve conflicts", `gitmap reconcile`),
		newCmd("has-any-updates", "hau", "v6.509.0", "sync", "Check remote tracking branches for incoming commits", `gitmap hau`),
		newCmd("latest-branch", "lb", "v6.510.0", "sync", "Discover the most recently updated remote branch across fleet", `gitmap lb`),
		newCmd("desktop-sync", "ds", "v6.503.0", "sync", "Synchronize local repositories with GitHub Desktop state", `gitmap desktop-sync`),
		newCmd("watch", "w", "v6.504.0", "sync", "Live-refresh terminal dashboard monitoring repo changes", `gitmap watch`),
	}
}

func catalogEntries10() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("new-commands", "nc", "v6.524.0", "tooling", "Inspect recent commands added across releases with examples", `gitmap nc --limit=20`),
		newCmd("which-format", "whichfmt", "v6.519.0", "tooling", "Detect line endings, encoding BOM, and file permissions", `gitmap which-format ./cli`),
		newCmd("find-files", "ff", "v6.512.0", "tooling", "Find exact filename across workspace within milliseconds", `gitmap ff "types.go"`),
		newCmd("find-files-any", "ffa", "v6.512.0", "tooling", "Find files matching substring across 10,000+ files", `gitmap ffa "cache"`),
		newCmd("find-files-startswith", "ffs", "v6.512.0", "tooling", "Find files matching filename prefix", `gitmap ffs "test_"`),
		newCmd("find-files-endswith", "ffe", "v6.512.0", "tooling", "Find files matching filename extension or suffix", `gitmap ffe "_test.go"`),
		newCmd("list-files", "lf", "v6.511.0", "tooling", "Stream relative file paths matching pattern or folder", `gitmap lf cli/cmdpending`),
		newCmd("replace", "rep", "v6.513.0", "tooling", "Literal multi-file string replacement with audit trail", `gitmap replace "oldString" "newString"`),
		newCmd("replace-regex", "repr", "v6.513.0", "tooling", "Regex pattern replacement across repository code", `gitmap replace-regex "v[0-9]+" "v6.524.0"`),
		newCmd("cat", "ct", "v6.514.0", "tooling", "Zero-disk stdout stream of file content in CLI sessions", `gitmap cat cli/constants/constants_cli.go`),
		newCmd("tree", "tree", "v6.528.0", "tooling", "Tree export with search engine family and relative paths", `gitmap tree --export`),
	}
}
