package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const PublicLlmSpecURL = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md"

const llmMarkdownSpec = `
# Gitmap LLM Specification & AI Agent Guidelines

> Public Instruction URL: ` + PublicLlmSpecURL + `

Gitmap is a high-performance CLI designed for autonomous AI agents, LLMs, and developers to explore, search, refactor, and self-heal codebases with maximum speed and reliability.


---

## 1. AI Agent Standard Operating Procedure (Order of Commands)

When an autonomous AI Agent is assigned a coding, refactoring, or debugging task, it MUST follow this structured 5-phase lifecycle in exact sequential order:

` + "```" + `
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 1: Discovery & Context Gathering                     │
  │  ➔ gitmap find-files, find-files-any, search, list-files   │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 2: Code Modification & Refactoring                   │
  │  ➔ gitmap replace, replace-regex, targeted code edits       │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 3: Local Verification & Linting                      │
  │  ➔ python .github/scripts/go-format-check.py, go test       │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 4: Semantic Commit & Push                            │
  │  ➔ gitmap commit-push-feature (cpf), commit-push-bug (cpb)  │
  └──────────────────────────────┬──────────────────────────────┘
                                 ▼
  ┌─────────────────────────────────────────────────────────────┐
  │  Phase 5: Non-Blocking CI Telemetry & Self-Healing Loop     │
  │  ➔ gitmap pipeline status --json, gitmap eta, error-logs   │
  └─────────────────────────────────────────────────────────────┘
` + "```" + `


---

## 2. Step-by-Step Workflows with Real-World Examples


### Workflow A: Autonomous CI/CD Self-Healing Telemetry Loop

When a remote pipeline is triggered or encounters an error, the AI Agent must follow this non-blocking telemetry loop:

1. **Check Live Pipeline Status & ETA**:
   ` + "```bash" + `
   gitmap pipeline status --json
   ` + "```" + `
   *Output:*
   ` + "```json" + `
   {
     "status": "RUNNING",
     "workflow": "Release",
     "etaSeconds": 75,
     "pendingPipelines": 1,
     "runUrl": "https://github.com/alimtvnetwork/gitmap-v28/actions/runs/12345"
   }
   ` + "```" + `

2. **Schedule Non-Blocking Wait Timer**:
   ` + "```bash" + `
   gitmap eta
   ` + "```" + `
   *Output:*
   ` + "```" + `
   75
   ` + "```" + `
   *(Use the returned integer to schedule a non-blocking timer in your harness; NEVER busy-poll).*

3. **Extract Failing CI Diagnostic Error Logs to File**:
   ` + "```bash" + `
   gitmap pipeline error-logs --json --tempfile "ci-failure.json"
   ` + "```" + `
   *Behavior:* Dumps the exact failing step, command line, exit code, and stderr into ` + "`ci-failure.json`" + ` for instant 4-part Root Cause Analysis (RCA).

4. **Verify Locally & Push Semantic Fix**:
   ` + "```bash" + `
   python .github/scripts/go-format-check.py
   gitmap commit-push-bug "fix(ci): fix gofmt formatting in cmd/root.go"
   ` + "```" + `


---

### Workflow B: Fast Code Search & File Discovery

Instead of scanning entire directory trees or reading large files into context:

- **Find Exact Filename**:
  ` + "```bash" + `
  gitmap find-files "constants.go" -ext "go"
  gitmap ff "package.json"
  ` + "```" + `

- **Find Files Containing Substring (with Extension Filter)**:
  ` + "```bash" + `
  gitmap find-files-any "runner" -ext "py,go"
  gitmap ffa "pipeline" -ext "go,json"
  ` + "```" + `

- **Find Files by Prefix or Suffix**:
  ` + "```bash" + `
  gitmap find-files-startswith "01-" -ext "md"
  gitmap find-files-endswith "_test.go" -ext "go"
  gitmap ffe ".spec.ts"
  ` + "```" + `

- **Universal Wildcard / Glob Search**:
  ` + "```bash" + `
  gitmap find "*config*.json"
  gitmap list-files "*pipeline*"
  gitmap f "*docker*" --limit 5
  ` + "```" + `

- **Fast Indexed Global Keyword Search**:
  ` + "```bash" + `
  gitmap search "Resolve-Version"
  gitmap repo-search-json "error-logs"   # alias: rsj (returns structured JSON matches)
  ` + "```" + `


---

### Workflow C: Batch Refactoring with History & Rollback

- **Exact Literal String Replacement**:
  ` + "```bash" + `
  gitmap replace "oldEndpoint" "newEndpoint"
  ` + "```" + `

- **Regex Replacement with Audit Trail**:
  ` + "```bash" + `
  gitmap replace-regex "v6\.\d+\.\d+" "v6.155.7"
  ` + "```" + `

- **View Replacement History**:
  ` + "```bash" + `
  gitmap replace history
  ` + "```" + `


---

### Workflow D: Workspace & Background Automation

- **Antigravity Workspaces**:
  ` + "```bash" + `
  gitmap agy add /path/to/repo       # Register repository in Antigravity workspaces
  gitmap agy sync                     # Synchronize agent configs across projects
  gitmap agy stats                    # View token and session statistics
  ` + "```" + `

- **VS Code Project Manager**:
  ` + "```bash" + `
  gitmap vsc add /path/to/repo        # Register workspace in VS Code Project Manager
  gitmap vsc ls                       # List tracked workspaces
  ` + "```" + `

- **Background Scheduler & Daemons**:
  ` + "```bash" + `
  gitmap schedule add "nightly-sync" --interval "24h"
  gitmap schedule status
  ` + "```" + `


---

## 3. Alternative Commands & Aliases Cheat Sheet (Instead of Raw Git)

| Action / Goal | Standard CLI Command | Fast AI Shortcut / Alias | Notes |
| :--- | :--- | :--- | :--- |
| **Pipeline Status** | ` + "`gitmap pipeline status --json`" + ` | ` + "`gitmap pipeline status`" + ` | Live execution state, ETA seconds, run URL |
| **Get ETA Seconds** | ` + "`gitmap pipeline waittime`" + ` | ` + "`gitmap eta`" + ` | Bare integer for non-blocking timer |
| **Extract Error Logs**| ` + "`gitmap pipeline error-logs`" + ` | ` + "`gitmap error-logs`" + ` | Dumps failing step logs to workspace file |
| **Full CI Run Logs**| ` + "`gitmap pipeline logs`" + ` | ` + "`gitmap logs`" + ` | Complete workflow execution trace |
| **Find File (Exact)** | ` + "`gitmap find-files <name>`" + ` | ` + "`gitmap ff <name>`" + ` | Match exact filename with optional -ext |
| **Find File (Any)** | ` + "`gitmap find-files-any <str>`" + ` | ` + "`gitmap ffa <str>`" + ` | Match substring in filename |
| **Find File (Prefix)**| ` + "`gitmap find-files-startswith`" + ` | ` + "`gitmap ffs <prefix>`" + ` | Match prefix in filename |
| **Find File (Suffix)**| ` + "`gitmap find-files-endswith`" + ` | ` + "`gitmap ffe <suffix>`" + ` | Match suffix in filename (e.g. _test.go) |
| **Commit Feature** | ` + "`gitmap commit-push-feature`" + ` | ` + "`gitmap cpf \"<msg>\"`" + ` | Stage, commit, and push feature branch |
| **Commit Bug Fix** | ` + "`gitmap commit-push-bug`" + ` | ` + "`gitmap cpb \"<msg>\"`" + ` | Stage, commit, and push bugfix branch |
| **Commit Release** | ` + "`gitmap commit-push-release`" + ` | ` + "`gitmap cpr \"<msg>\"`" + ` | Stage, commit, and push release chore |
| **Pull & Push** | ` + "`gitmap pull-commit-push`" + ` | ` + "`gitmap pcp \"<msg>\"`" + ` | Pull latest, commit, and push |
| **Antigravity** | ` + "`gitmap antigravity`" + ` | ` + "`gitmap agy`, `gitmap ag`" + ` | AI agent workspaces and config sync |
| **VS Code PM** | ` + "`gitmap vscode`" + ` | ` + "`gitmap vsc`" + ` | VS Code Project Manager integrations |
| **Scheduler** | ` + "`gitmap schedule`" + ` | ` + "`gitmap sc`" + ` | Background cron and interval scheduler |
| **Repo Status** | ` + "`gitmap status`" + ` | ` + "`gitmap st`" + ` | Dirty, ahead, behind across all repos |
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
