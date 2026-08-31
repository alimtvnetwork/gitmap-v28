package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const PublicLlmSpecURL = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md"

const llmMarkdownSpec = `# Gitmap LLM Specification & AI Agent Guidelines

> Public Instruction URL: ` + PublicLlmSpecURL + `

Gitmap is a powerful CLI designed for autonomous agents, LLMs, and developers to navigate, analyze, and modify codebases efficiently.

## Core Capabilities

1. **Repository Discovery & Scanning**:
   - ` + "`gitmap scan`" + `: Automatically locate and index Git repositories.
   - ` + "`gitmap rescan`" + `: Refresh indexed local repositories.

2. **Repository Cloning & Synchronization**:
   - ` + "`gitmap clone <source>`" + `, ` + "`gitmap clone-next`" + `: Manage bulk cloning and automated directory migration.

3. **Intelligent File Search & Discovery**:
   - ` + "`gitmap find-files \"<name>\" -ext \"<exts>\"`" + `: Find files matching exact filename.
   - ` + "`gitmap find-files-any \"<substring>\" -ext \"<exts>\"`" + `: Find files containing substring with extension filters.
   - ` + "`gitmap find-files-startswith \"<prefix>\" -ext \"<exts>\"`" + `: Find files by name prefix.
   - ` + "`gitmap find-files-endswith \"<suffix>\" -ext \"<exts>\"`" + `: Find files by name suffix.
   - ` + "`gitmap find \"<wildcard*>\" -ext \"<exts>\"`" + `: Universal glob search.
   - ` + "`gitmap list-files [pattern]`" + `: List indexed files.
   - ` + "`gitmap search <query>`" + `: Execute SplitDB indexed searches.
   - ` + "`gitmap repo-search-json <query>`" + ` (alias ` + "`rsj`" + `): Retrieve search results as structured JSON.

4. **CI/CD Pipeline Telemetry & Error Recovery (For Autonomous Self-Healing)**:
   - ` + "`gitmap pipeline status --json`" + `: Live pipeline execution state, active workflow, ETA seconds, and pending PRs.
   - ` + "`gitmap eta`" + ` (alias ` + "`gitmap pipeline eta`" + `, ` + "`gitmap pipeline waittime`" + `): Remaining wait time in seconds as an integer for timer scheduling.
   - ` + "`gitmap pipeline error-logs --json --tempfile \"ci-failure.json\"`" + ` (alias ` + "`gitmap error-logs`" + `): Directly extract failing step logs to temp file for 4-part RCA and automatic error recovery.
   - ` + "`gitmap logs`" + ` (alias ` + "`gitmap pipeline logs`" + `): Display full workflow run logs.

5. **Regex & Replacement**:
   - ` + "`gitmap replace <query> <replacement>`" + `: Find and replace exact phrases.
   - ` + "`gitmap replace-regex <regex> <replacement>`" + `: Regex replacement with automatic history tracking.
   - ` + "`gitmap replace history`" + `: View replacement operations.

6. **Releases & Commits**:
   - ` + "`gitmap release`" + `: Automate semantic version bumping and changelog generation.
   - ` + "`gitmap commit-in`" + `: Perform robust, orchestrator-driven batch commits.

## Instructions for LLMs & AI Coding Assistants

- Always prefer JSON output commands (` + "`rsj`" + `, ` + "`gitmap pipeline status --json`" + `, ` + "`gitmap pipeline error-logs --json`" + `) when parsing data programmatically.
- After pushing commits or releasing tags, use ` + "`gitmap eta`" + ` or ` + "`gitmap pipeline status --json`" + ` to monitor CI/CD progress. Never poll in a busy loop; use the returned ` + "`etaSeconds`" + ` to schedule non-blocking sleep/timer intervals.
- When CI/CD fails, use ` + "`gitmap pipeline error-logs --json --tempfile \"ci-failure.json\"`" + ` to dump the failed step's stdout/stderr directly into a workspace file for instant root cause analysis (RCA).
- Avoid modifying the root SQLite DB manually; rely on the CLI commands.
- Before running heavy regex operations across a large repo, use ` + "`gitmap search`" + ` or ` + "`gitmap find-files-any`" + ` to verify matches.
- All file references in plans, specs, and subtasks MUST be strictly relative to the repository root.

## Alternative Commands for AI (Instead of Raw Git)

When checking repository status or logs, DO NOT use raw 'git status && git log -1'. Instead, use:
- ` + "`gitmap status`" + `: Checks dirty/clean status, ahead/behind, and stash for all tracked repos.
- ` + "`gitmap history`" + `: View a rich audit log of command executions.
- ` + "`gitmap changelog`" + `: See concise release notes and recent commits.

When checking CI/CD and deployment health, DO NOT guess or loop blind. Use:
- ` + "`gitmap pipeline status`" + `: Human-readable live workflow status and ETA.
- ` + "`gitmap eta`" + `: Bare integer ETA in seconds for timer scheduling.
- ` + "`gitmap error-logs`" + `: Diagnostic failure logs.
- ` + "`gitmap logs`" + `: Full workflow step trace.

For committing and pushing, NEVER use raw 'git commit' and 'git push'. Use the semantic commit-push suite:
- ` + "`gitmap commit-push-feature \"<message>\"`" + ` (alias ` + "`gitmap cpf`" + `): Commit and push a feature.
- ` + "`gitmap commit-push-bug \"<message>\"`" + ` (alias ` + "`gitmap cpb`" + `): Commit and push a bug fix.
- ` + "`gitmap commit-push-release \"<message>\"`" + ` (alias ` + "`gitmap cpr`" + `): Commit and push a release chore.
- ` + "`gitmap pull-commit-push \"<message>\"`" + ` (alias ` + "`gitmap pcp`" + `): Pull latest, then commit and push.
`

// Run executes the llm command.
func Run(args []string) *apperror.AppError {
	fs := flag.NewFlagSet("llm", flag.ExitOnError)
	isUrl := fs.Bool("url", false, "Output the URL to the LLM spec")
	isInstruction := fs.Bool("instruction", false, "Output the full markdown instructions")

	if err := fs.Parse(args); err != nil {
		return apperror.WrapSimple(err, "parse flags")
	}

	if *isUrl {
		fmt.Println(PublicLlmSpecURL)
		return nil
	}

	if *isInstruction {
		fmt.Print(llmMarkdownSpec)
		return nil
	}

	fmt.Print(llmMarkdownSpec)
	return nil
}
