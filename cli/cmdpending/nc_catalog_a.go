package cmdpending
func catalogEntries1() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("cpar", "commit-push-all-repos", "v6.520.0", "commits", "Commit and push all workspace repositories", `gitmap cpar "wip: save changes"`),
		newCmd("cpf", "commit-push-feature", "v6.500.0", "commits", "Atomic commit and push feature branch", `gitmap cpf "auth - add jwt validation"`),
		newCmd("cpb", "commit-push-bug", "v6.500.0", "commits", "Atomic commit and push bugfix branch", `gitmap cpb "cache - fix ttl expiration"`),
		newCmd("cpr", "commit-push-release", "v6.500.0", "commits", "Commit, push, and trigger release chore", `gitmap cpr "v6.524.0"`),
		newCmd("cin", "commit-in", "v6.510.0", "commits", "Commit in replay with JSON author rotation and AST heuristics", `gitmap cin --source repo-a --dest repo-b`),
		newCmd("pending-commits", "pc", "v6.523.0", "commits", "Discover uncommitted changes and unpushed commits with tree remediation", `gitmap pc`),
		newCmd("backup-branch", "bb", "v6.518.0", "commits", "Create automated backup branch for workspace repositories", `gitmap backup-branch "feat-login"`),
		newCmd("fix all", "fix-all", "v6.522.0", "commits", "Batch stash, commit, or discard pending changes across fleet", `gitmap fix all --action=wip`),
		newCmd("pcp", "pull-commit-push", "v6.505.0", "commits", "Pull rebase, stage all, commit with message, and push", `gitmap pcp "sync upstream"`),
		newCmd("commons", "co", "v6.515.0", "commits", "Apply standardized repo configurations and git hooks", `gitmap commons`),
	}
}

func catalogEntries2() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("pe", "pipeline-errors", "v6.510.0", "diagnostics", "Pipeline error analyzer with dynamic runner extraction", `gitmap pe -t`),
		newCmd("pe all", "pipeline-errors-all", "v6.515.0", "diagnostics", "Extract failing step logs across all recent workflow runs", `gitmap pe all`),
		newCmd("te all", "pipeline-errors-all", "v6.522.0", "diagnostics", "Shortcut for pipeline errors all: aggregate pipeline errors across all repos", `gitmap te all --file=errors.md`),
		newCmd("sug", "shutdown-until-green", "v6.512.0", "diagnostics", "Pause execution until CI/CD quality gates are green", `gitmap sug`),
		newCmd("rerun", "rr", "v6.511.0", "diagnostics", "Rerun failed pipeline jobs or Antigravity prompts", `gitmap rerun --step=lint`),
		newCmd("regoldens", "rg-gold", "v6.517.0", "diagnostics", "Regenerate golden test fixture outputs across test suites", `gitmap regoldens`),
		newCmd("doctor", "doc", "v6.508.0", "diagnostics", "Comprehensive workstation and toolchain diagnostic health check", `gitmap doctor`),
		newCmd("audit-legacy", "al", "v6.514.0", "diagnostics", "Scan codebase for deprecated APIs and legacy function signatures", `gitmap audit-legacy`),
		newCmd("fastgate", "fg", "v6.519.0", "diagnostics", "Run fast preflight lint and formatting verification gates", `gitmap fastgate`),
		newCmd("pipeline-ai", "pl-ai", "v6.513.0", "diagnostics", "Check CI/CD workflow state and wait dynamically on ETA", `gitmap pipeline-ai status -t 120`),
	}
}

func catalogEntries3() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("scan export", "s-exp", "v6.521.0", "scanner", "Export scanned repository metadata to JSON or CSV manifest", `gitmap scan export --format=json`),
		newCmd("scan merge", "s-mrg", "v6.521.0", "scanner", "Merge external scan databases into central repository index", `gitmap scan merge --src=other.db`),
		newCmd("rescan", "rsc", "v6.502.0", "scanner", "Perform delta rescan of workspace filesystem changes", `gitmap rescan`),
		newCmd("rescan-subtree", "rss", "v6.504.0", "scanner", "Narrowly re-run scan against an at-cap repository subtree", `gitmap rescan-subtree --max-depth=5`),
		newCmd("clone-only-missing", "com", "v6.506.0", "scanner", "Clone only un-cloned repositories from workspace manifest", `gitmap com --manifest=repos.json`),
		newCmd("clone-sync", "cs", "v6.507.0", "scanner", "Synchronize and clone all remote fleet repositories", `gitmap clone-sync`),
		newCmd("repo-create", "repoc", "v6.503.0", "scanner", "Create a new local repository with standard folder scaffold", `gitmap repo-create my-service`),
		newCmd("dedupe", "dd", "v6.509.0", "scanner", "Detect duplicate repository clones and directories across disks", `gitmap dedupe`),
		newCmd("size", "sz", "v6.501.0", "scanner", "Inspect repository disk size and identify large bloated objects", `gitmap size --top=10`),
		newCmd("orphans", "orph", "v6.510.0", "scanner", "Discover orphaned git directories not indexed in split-db", `gitmap orphans`),
	}
}

func catalogEntries4() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("ssh-bind", "sb", "v6.518.0", "fleet", "Bind local port forwarding tunnels over SSH to fleet nodes", `gitmap ssh-bind --remote=node-01 --port=8080`),
		newCmd("ssh-deploy", "sd", "v6.516.0", "fleet", "Deploy compiled binaries and configuration to remote SSH host", `gitmap ssh-deploy --node=prod-1`),
		newCmd("nodes", "cluster-nodes", "v6.512.0", "fleet", "List registered fleet nodes with reachability and latency", `gitmap nodes`),
		newCmd("cluster", "cl", "v6.511.0", "fleet", "Manage multi-node cluster topology and remote execution", `gitmap cluster status`),
		newCmd("sc", "servers-clients", "v6.514.0", "fleet", "Distributed fan-out execution across server-client topologies", `gitmap sc exec "uptime"`),
		newCmd("deploy-keys", "dk", "v6.515.0", "fleet", "Deploy authorized SSH public keys across all fleet machines", `gitmap deploy-keys --all`),
		newCmd("ssh-join", "sj", "v6.510.0", "fleet", "Enroll target machine into SSH host registry with key auth", `gitmap ssh-join user@192.168.1.50`),
		newCmd("cluster-run-script", "crs", "v6.513.0", "fleet", "Deploy and execute local script remotely across cluster nodes", `gitmap cluster-run-script deploy.sh`),
		newCmd("cluster-bootstrap", "cb", "v6.515.0", "fleet", "Bootstrap node with RSA keys, passwordless sudo, and tools", `gitmap cluster-bootstrap 192.168.1.51`),
		newCmd("ssh-scan", "ss", "v6.512.0", "fleet", "Scan local subnet for machines with open SSH port 22", `gitmap ssh-scan 192.168.1.0/24`),
	}
}

func catalogEntries5() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("agent task", "at", "v6.522.0", "ai", "Manage autonomous agent task lifecycle and action logs in SQLite", `gitmap task claim --db task.db --agent "Worker 01"`),
		newCmd("agy", "antigravity", "v6.510.0", "ai", "Google Antigravity workspace manager and prompt orchestrator", `gitmap agy deploy`),
		newCmd("agm", "antigravity-mgr", "v6.512.0", "ai", "Manage Antigravity Manager GUI application and fleet tools", `gitmap agm status`),
		newCmd("aum", "auto", "v6.511.0", "ai", "High-performance native Go automation and multi-core search", `gitmap aum search "pattern" cli`),
		newCmd("ai-analysis", "aa", "v6.514.0", "ai", "Perform automated codebase architecture and guideline analysis", `gitmap ai-analysis --rules=02-spec`),
		newCmd("macro", "mc", "v6.509.0", "ai", "Record, inspect, and replay interactive shell terminal macros", `gitmap macro replay build-all`),
		newCmd("llm-docs", "ld", "v6.515.0", "ai", "Generate consolidated markdown command matrix for LLMs", `gitmap llm-docs`),
		newCmd("llm train", "llm-train", "v6.516.0", "ai", "Run 4-stage LLM curriculum and generate developer skill", `gitmap llm train`),
		newCmd("prompt show --copy", "psc", "v6.520.0", "ai", "Render canonical prompt and copy directly to system clipboard", `gitmap prompt show --copy 01-prompts/v6.md`),
		newCmd("ai-fix", "af", "v6.513.0", "ai", "Run standardized repository autofix targets for coding guidelines", `gitmap ai-fix --target=naming`),
		newCmd("merge-ai", "merge-ai", "v6.525.0", "ai", "AI-assisted merge with conflict resolution and summary", `gitmap merge-ai`),
	}
}
