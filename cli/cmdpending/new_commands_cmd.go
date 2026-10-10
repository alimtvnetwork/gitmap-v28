// Package cmdpending provides command implementations for pending commits and discovery.
package cmdpending

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RunNewCommands executes the new-commands discovery tool.
func RunNewCommands(args []string) error {
	opts, appErr := parseNewCommandsOptions(args)
	if appErr != nil {
		return appErr
	}
	if opts.IsHelp {
		printNewCommandsHelp()
		return nil
	}
	return executeNewCommands(opts)
}

func executeNewCommands(opts NewCommandsOptions) error {
	catalog := buildNewCommandsCatalog()
	payload := buildNewCommandsPayload(catalog, opts)
	if opts.IsJSON {
		return renderNewCommandsJSON(payload)
	}
	return renderNewCommandsTerminal(payload)
}

func parseNewCommandsOptions(args []string) (NewCommandsOptions, *appfault.AppError) {
	opts := DefaultNewCommandsOptions()
	for i := 0; i < len(args); i++ {
		nextIdx, appErr := parseNextCommandOption(args, i, &opts)
		if appErr != nil {
			return opts, appErr
		}
		i = nextIdx
	}
	clampNewCommandsLimit(&opts)
	return opts, nil
}

func parseNextCommandOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isHelpFlag(arg) {
		opts.IsHelp = true
		return idx, nil
	}
	if isJSONFlag(arg) {
		opts.IsJSON = true
		return idx, nil
	}
	return parseValuedOption(args, idx, opts)
}

func parseValuedOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	arg := args[idx]
	if isCategoryFlag(arg) {
		return parseCategoryOption(args, idx, opts)
	}
	if isFilterFlag(arg) {
		return parseFilterOption(args, idx, opts)
	}
	if isLimitFlag(arg) {
		return parseLimitOption(args, idx, opts)
	}
	return idx, nil
}

func parseCategoryOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --category")
	}
	opts.Category = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseFilterOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --filter")
	}
	opts.Filter = strings.TrimSpace(args[idx+1])
	return idx + 1, nil
}

func parseLimitOption(args []string, idx int, opts *NewCommandsOptions) (int, *appfault.AppError) {
	if idx+1 >= len(args) {
		return idx, appfault.NewValidation("new-commands", "E9003", "missing value for --limit")
	}
	lim, err := strconv.Atoi(args[idx+1])
	if err != nil {
		return idx, appfault.NewValidation("new-commands", "E9003", "invalid numeric value for --limit: "+args[idx+1])
	}
	opts.Limit = lim
	return idx + 1, nil
}

func clampNewCommandsLimit(opts *NewCommandsOptions) {
	if opts.Limit <= 0 {
		opts.Limit = 100
	}
}

func isHelpFlag(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func isJSONFlag(arg string) bool {
	return arg == "--json" || arg == "-j"
}

func isCategoryFlag(arg string) bool {
	return arg == "--category" || arg == "-c"
}

func isFilterFlag(arg string) bool {
	return arg == "--filter" || arg == "-f" || arg == "-q"
}

func isLimitFlag(arg string) bool {
	return arg == "--limit" || arg == "-n"
}

func buildNewCommandsPayload(catalog []NewCommandEntry, opts NewCommandsOptions) NewCommandsPayload {
	filtered := filterNewCommands(catalog, opts)
	return NewCommandsPayload{
		Timestamp:        time.Now().UTC(),
		TotalCommands:    len(catalog),
		FilteredCommands: len(filtered),
		Limit:            opts.Limit,
		Category:         opts.Category,
		Filter:           opts.Filter,
		Commands:         filtered,
	}
}

func filterNewCommands(catalog []NewCommandEntry, opts NewCommandsOptions) []NewCommandEntry {
	res := make([]NewCommandEntry, 0, len(catalog))
	for _, entry := range catalog {
		if isCommandMatched(entry, opts) {
			res = append(res, entry)
		}
		if len(res) >= opts.Limit {
			break
		}
	}
	return res
}

func isCommandMatched(entry NewCommandEntry, opts NewCommandsOptions) bool {
	if !isCategoryMatch(entry.Category, opts.Category) {
		return false
	}
	if !isTextFilterMatch(entry, opts.Filter) {
		return false
	}
	return true
}

func isCategoryMatch(entryCat, filterCat string) bool {
	trimmed := strings.TrimSpace(filterCat)
	if trimmed == "" {
		return true
	}
	return strings.EqualFold(entryCat, trimmed)
}

func isTextFilterMatch(entry NewCommandEntry, filter string) bool {
	query := strings.ToLower(strings.TrimSpace(filter))
	if query == "" {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Name), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Alias), query) {
		return true
	}
	return isEntryDetailMatch(entry, query)
}

func isEntryDetailMatch(entry NewCommandEntry, query string) bool {
	if strings.Contains(strings.ToLower(entry.Description), query) {
		return true
	}
	if strings.Contains(strings.ToLower(entry.Example), query) {
		return true
	}
	return false
}

func renderNewCommandsJSON(payload NewCommandsPayload) error {
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return appfault.WrapExecution(err, "failed to serialize new commands JSON")
	}
	fmt.Println(string(data))
	return nil
}

func renderNewCommandsTerminal(payload NewCommandsPayload) error {
	renderTerminalHeader(payload)
	grouped := groupCommandsByCategory(payload.Commands)
	renderGroupedCategories(grouped)
	renderTerminalFooter(payload)
	return nil
}

func renderTerminalHeader(payload NewCommandsPayload) {
	fmt.Println("╔══════════════════════════════════════════════════════════════════════════════╗")
	fmt.Println("║                 GITMAP RECENT COMMANDS DISCOVERY ENGINE                      ║")
	fmt.Println("╠══════════════════════════════════════════════════════════════════════════════╣")
	line := formatHeaderSummaryLine(payload)
	fmt.Println(line)
	fmt.Println("╚══════════════════════════════════════════════════════════════════════════════╝")
	fmt.Println()
}

func formatHeaderSummaryLine(payload NewCommandsPayload) string {
	info := fmt.Sprintf("║ Total: %d | Returned: %d | Limit: %d", payload.TotalCommands, payload.FilteredCommands, payload.Limit)
	if payload.Category != "" {
		info += " | Cat: " + payload.Category
	}
	if payload.Filter != "" {
		info += " | Query: " + payload.Filter
	}
	for len(info) < 79 {
		info += " "
	}
	return info + "║"
}

type categoryGroup struct {
	category string
	commands []NewCommandEntry
}

func groupCommandsByCategory(commands []NewCommandEntry) []categoryGroup {
	var groups []categoryGroup
	groupMap := make(map[string]int)
	for _, cmd := range commands {
		idx, hasCat := groupMap[cmd.Category]
		if hasCat {
			groups[idx].commands = append(groups[idx].commands, cmd)
			continue
		}
		groupMap[cmd.Category] = len(groups)
		groups = append(groups, categoryGroup{category: cmd.Category, commands: []NewCommandEntry{cmd}})
	}
	return groups
}

func renderGroupedCategories(groups []categoryGroup) {
	for _, grp := range groups {
		renderCategorySection(grp)
	}
}

func renderCategorySection(grp categoryGroup) {
	fmt.Printf("▶ %s (%d commands)\n", strings.ToUpper(grp.category), len(grp.commands))
	fmt.Println(strings.Repeat("─", 80))
	for _, cmd := range grp.commands {
		renderCommandEntry(cmd)
	}
	fmt.Println()
}

func renderCommandEntry(cmd NewCommandEntry) {
	title := formatCommandTitle(cmd)
	fmt.Println(title)
	fmt.Printf("    %s\n", cmd.Description)
	fmt.Printf("    Example: %s\n", cmd.Example)
	fmt.Println()
}

func formatCommandTitle(cmd NewCommandEntry) string {
	res := "  ● gitmap " + cmd.Name
	if cmd.Alias != "" {
		res += " (" + cmd.Alias + ")"
	}
	if cmd.Version != "" {
		res += " [" + cmd.Version + "]"
	}
	return res
}

func renderTerminalFooter(payload NewCommandsPayload) {
	fmt.Println("────────────────────────────────────────────────────────────────────────────────")
	fmt.Printf("Tip: Use 'gitmap nc --category <cat>' or '--filter <query>' to narrow results.\n")
	fmt.Printf("     Use 'gitmap nc --json' for machine-readable JSON output.\n")
}

func printNewCommandsHelp() {
	fmt.Println("GitMap New Commands Discovery Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap new-commands [flags]")
	fmt.Println("  gitmap nc [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  --limit, -n <N>          Limit number of commands returned (default: 100)")
	fmt.Println("  --category, -c <cat>     Filter by functional category")
	fmt.Println("  --filter, -f, -q <query> Search commands by name, alias, description, or example")
	fmt.Println("  --json, -j               Emit output as formatted JSON")
	fmt.Println("  -h, --help               Show this help text")
}

func newCmd(name, alias, ver, cat, desc, ex string) NewCommandEntry {
	return NewCommandEntry{
		Name:        name,
		Alias:       alias,
		Version:     ver,
		Category:    cat,
		Description: desc,
		Example:     ex,
	}
}

func buildNewCommandsCatalog() []NewCommandEntry {
	res := make([]NewCommandEntry, 0, 100)
	res = append(res, catalogBatch1to5()...)
	return append(res, catalogBatch6to10()...)
}

func catalogBatch1to5() []NewCommandEntry {
	items := append(catalogEntries1(), catalogEntries2()...)
	items = append(items, catalogEntries3()...)
	items = append(items, catalogEntries4()...)
	return append(items, catalogEntries5()...)
}

func catalogBatch6to10() []NewCommandEntry {
	items := append(catalogEntries6(), catalogEntries7()...)
	items = append(items, catalogEntries8()...)
	items = append(items, catalogEntries9()...)
	return append(items, catalogEntries10()...)
}

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
		newCmd("te all", "test-errors-all", "v6.516.0", "diagnostics", "Aggregate failing unit test output across all packages", `gitmap te all`),
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
	}
}

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
	}
}

func catalogEntries7() []NewCommandEntry {
	return []NewCommandEntry{
		newCmd("spec issue", "si", "v6.522.0", "spec", "Audit specification gap issues and generate remediation steps", `gitmap spec issue --spec=02-spec/21-app`),
		newCmd("spec-author", "sa", "v6.520.0", "spec", "Scaffold structured architecture and component specifications", `gitmap spec-author --slug=my-feature`),
		newCmd("spec-audit", "spa", "v6.519.0", "spec", "Conduct blind-AI readiness audits on specification documents", `gitmap spec-audit --target=02-spec`),
		newCmd("user", "usr", "v6.521.0", "spec", "Audit and configure git user name and email per repository", `gitmap user set --name="Dev" --email="dev@co.com"`),
		newCmd("changelog", "clog", "v6.508.0", "spec", "Generate SemVer changelog entries from recent atomic commits", `gitmap changelog generate`),
		newCmd("release-notes", "rn", "v6.509.0", "spec", "Extract markdown release notes for tagged versions", `gitmap release-notes v6.523.0`),
		newCmd("seowrite", "seo", "v6.510.0", "spec", "Format repository descriptions and keywords for SEO", `gitmap seowrite --keyword=git`),
		newCmd("help-json", "hj", "v6.515.0", "spec", "Emit machine-readable JSON schema for CLI help topics", `gitmap help --json --filter=task`),
		newCmd("prompt-template", "pt", "v6.512.0", "spec", "Manage reusable AI prompt prefix and verification templates", `gitmap prompt-template list`),
		newCmd("prompt list", "plst", "v6.518.0", "spec", "List registered canonical AI prompts in 01-prompts directory", `gitmap prompt list`),
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
	}
}
