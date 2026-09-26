// Package commitin — help.go prints detailed multi-section usage documentation.
package commitin

import (
	"strings"
)

// PrintCommitInHelp returns detailed documentation and JSON examples for commit-in.
func PrintCommitInHelp() string {
	var b strings.Builder
	writeHelpHeader(&b)
	writeHelpOverview(&b)
	writeHelpExclusions(&b)
	writeHelpPRWorkflow(&b)
	writeHelpFinalSnapshot(&b)
	writeHelpUsageExamples(&b)
	writeHelpJSONExamples(&b)

	return b.String()
}

func writeHelpHeader(b *strings.Builder) {
	b.WriteString("gitmap commit-in (cin) / pr in: Chronological Commit & PR Replay Engine\n")
	b.WriteString("========================================================================\n\n")
}

func writeHelpOverview(b *strings.Builder) {
	b.WriteString("OVERVIEW:\n")
	b.WriteString("  Walks one or more SOURCE repositories chronologically and APPENDS each\n")
	b.WriteString("  commit into a TARGET repository, preserving both AuthorDate and CommitterDate.\n")
	b.WriteString("  Idempotent via ShaMap in SQLite (.gitmap/data/commit-in/<slug>/sql.db).\n\n")
}

func writeHelpExclusions(b *strings.Builder) {
	b.WriteString("EXCLUSION & FILTERING RULES:\n")
	b.WriteString("  --exclude <glob>          Skip files matching glob patterns (e.g. '*.log', 'vendor/**')\n")
	b.WriteString("  --message-exclude <regex> Skip commits whose subject matches regex (e.g. '^WIP', '^Merge')\n")
	b.WriteString("  --function-intel on|off   FuncIntel: AST parsing of Go/TS/Python to filter out comment churn\n")
	b.WriteString("  SEO Commit Scheduling:    Generates formatted commit messages with custom templates\n")
	b.WriteString("  Default exclusions ignore .git/, node_modules/, and .gitmap/ internal state.\n\n")
}

func writeHelpPRWorkflow(b *strings.Builder) {
	b.WriteString("PR & MERGE WORKFLOW (--pr):\n")
	b.WriteString("  When --pr merges is enabled, every merge commit detected in the source history\n")
	b.WriteString("  is simulated through a local feature branch (feature/<slug> or pr/<id>), generates\n")
	b.WriteString("  a comprehensive Markdown PR description, merges into target mainline (--no-ff),\n")
	b.WriteString("  and registers records in .gitmap/data/pr/<slug>/sql.db.\n\n")
}

func writeHelpFinalSnapshot(b *strings.Builder) {
	b.WriteString("DETERMINISTIC FINAL SNAPSHOT SYNCHRONIZATION:\n")
	b.WriteString("  At the end of execution, the engine runs a final mirror snapshot synchronization\n")
	b.WriteString("  between source and target. Any stray target-only files are removed and all source\n")
	b.WriteString("  files are synchronized, guaranteeing an exact byte-for-byte repository match.\n\n")
}

func writeHelpUsageExamples(b *strings.Builder) {
	b.WriteString("COMMAND USAGE:\n")
	b.WriteString("  gitmap commit-in <target> <input-1> <input-2> ... [flags]\n")
	b.WriteString("  gitmap commit-in \"New Target Repo\" <input-1> ... (auto-provisions target)\n")
	b.WriteString("  gitmap commit-in <target> \"https://github.com/alimtvnetwork/gitmap-v{2..28}\" --pr merges --sponsor\n")
	b.WriteString("  gitmap commit-pull <target> gitmap-v2..v28 --tree --sponsor\n")
	b.WriteString("  gitmap cpull <target> <inputs...> --pr merges\n")
	b.WriteString("  gitmap pr in <target> <input-1> <input-2> ... [flags]\n")
	b.WriteString("  gitmap cin <target> all --since 2026-01-01 --pr merges\n")
	b.WriteString("  gitmap pr-clean <target> --yes\n\n")
	b.WriteString("FLAGS:\n")
	b.WriteString("  --sponsor                 Inject Rise Up Asia LLC sponsor templates into commit & PR bodies\n")
	b.WriteString("  --seo-template <name>     Use specific SEO template category (e.g. 'riseup')\n")
	b.WriteString("  --tree                    Render preflight PR/branch tree before execution\n")
	b.WriteString("  --pr merges               Simulate feature branches and PR merges\n\n")
}

func writeHelpJSONExamples(b *strings.Builder) {
	b.WriteString("JSON CONFIGURATION / REPLAY EXAMPLES:\n")
	b.WriteString("  {\n")
	b.WriteString("    \"command\": \"commit-pull\",\n")
	b.WriteString("    \"target\": \"./unified-gitmap\",\n")
	b.WriteString("    \"range\": \"v2..v28\",\n")
	b.WriteString("    \"inputs\": [\"https://github.com/alimtvnetwork/git-repo-navigator\", \"https://github.com/alimtvnetwork/gitmap-v{2..28}\"],\n")
	b.WriteString("    \"prMode\": \"merges\",\n")
	b.WriteString("    \"sponsor\": true,\n")
	b.WriteString("    \"seoTemplate\": \"riseup\",\n")
	b.WriteString("    \"tree\": true,\n")
	b.WriteString("    \"excludePatterns\": [\"*.log\", \"tmp/**\", \"node_modules/**\"],\n")
	b.WriteString("    \"messageExclusions\": [\"^WIP:\", \"^chore\\\\(sync\\\\):\"],\n")
	b.WriteString("    \"finalSnapshotSync\": true,\n")
	b.WriteString("    \"dryRun\": false\n")
	b.WriteString("  }\n\n")
	b.WriteString("JSON PR EVENT SCHEMA:\n")
	b.WriteString("  {\n")
	b.WriteString("    \"prNumber\": 101,\n")
	b.WriteString("    \"title\": \"feat(auth): token refresh rotation\",\n")
	b.WriteString("    \"sourceBranch\": \"feature/auth-refresh\",\n")
	b.WriteString("    \"targetBranch\": \"main\",\n")
	b.WriteString("    \"status\": \"merged\",\n")
	b.WriteString("    \"mergeCommitSha\": \"a7b8c9d0e1f2\",\n")
	b.WriteString("    \"commitCount\": 3,\n")
	b.WriteString("    \"filesChangedCount\": 8\n")
	b.WriteString("  }\n")
}
