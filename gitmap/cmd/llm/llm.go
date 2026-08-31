package llm

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

const PublicLlmSpecURL = "https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/llm.md"

const llmMarkdownSpec = `# Gitmap LLM Specification & AI Agent Guidelines

> Public Instruction URL: ` + PublicLlmSpecURL + `

Gitmap is a high-performance CLI designed for autonomous agents, LLMs, and developers to navigate, analyze, modify, and self-heal codebases efficiently.

---

## 1. AI Agent Standard Operating Procedure (Execution Order)

When an AI Coding Agent operates in a repository, it MUST follow this structured 5-phase lifecycle:

` + "```" + `
  ┌─────────────────────────┐
  │ 1. Discovery & Context  │ ➔ find-files, find-files-any, search, list
  └────────────┬────────────┘
               ▼
  ┌─────────────────────────┐
  │ 2. Edit & Refactor      │ ➔ replace, replace-regex, code edits
  └────────────┬────────────┘
               ▼
  ┌─────────────────────────┐
  │ 3. Local Verification   │ ➔ go-format-check.py, go test, linting
  └────────────┬────────────┘
               ▼
  ┌─────────────────────────┐
  │ 4. Semantic Commit/Push │ ➔ commit-push-feature (cpf), commit-push-bug (cpb)
  └────────────┬────────────┘
               ▼
  ┌─────────────────────────┐
  │ 5. CI Telemetry & Loop  │ ➔ pipeline status --json, eta (sleep), error-logs
  └─────────────────────────┘
` + "```" + `

---

## 2. Core Workflows with Concrete Examples

### Workflow A: Autonomous CI/CD Self-Healing Loop
When a remote pipeline fails or is running:
1. **Check Live Status & ETA**:
   ` + "```bash" + `
   gitmap pipeline status --json
   # Returns: {"status":"RUNNING","workflow":"Release","etaSeconds":75,"pendingPipelines":1}
   ` + "```" + `
2. **Schedule Non-Blocking Wait**:
   ` + "```bash" + `
   gitmap eta
   # Returns bare integer (e.g. 75). Use this value to schedule timers without busy polling.
   ` + "```" + `
3. **Extract Diagnostic Failure Logs**:
   ` + "```bash" + `
   gitmap pipeline error-logs --json --tempfile "ci-failure.json"
   # Extracts the exact failing step, stderr, and exit code directly into ci-failure.json.
   ` + "```" + `
4. **Fix Code & Push Semantic Fix**:
   ` + "```bash" + `
   python .github/scripts/go-format-check.py
   gitmap commit-push-bug "fix(ci): fix gofmt formatting in cmd/root.go"
   ` + "```" + `

---

### Workflow B: Fast Code Search & File Discovery
Instead of scanning entire directories or reading unnecessary files:
- **Exact File Match**:
  ` + "```bash" + `
  gitmap find-files "constants.go" -ext "go"
  gitmap ff "package.json"
  ` + "```" + `
- **Substring Match with Multi-Extension Filter**:
  ` + "```bash" + `
  gitmap find-files-any "runner" -ext "py, go"
  gitmap ffa "pipeline" -ext "go, json"
  ` + "```" + `
- **Prefix and Suffix Match**:
  ` + "```bash" + `
  gitmap find-files-startswith "01-" -ext "md"
  gitmap find-files-endswith "_test.go" -ext "go"
  ` + "```" + `
- **Glob / Wildcard Pattern Search**:
  ` + "```bash" + `
  gitmap find "*config*.json"
  gitmap list-files "*pipeline*"
  ` + "```" + `
- **Fast Indexed Global Keyword Search**:
  ` + "```bash" + `
  gitmap search "Resolve-Version"
  gitmap repo-search-json "error-logs"   # alias: rsj (returns structured JSON matches)
  ` + "```" + `

---

### Workflow C: Batch Refactoring with History & Rollback
- **Exact String Replacement**:
  ` + "```bash" + `
  gitmap replace "oldEndpoint" "newEndpoint"
  ` + "```" + `
- **Regex Replacement with Audit Trail**:
  ` + "```bash" + `
  gitmap replace-regex "v6\.\d+\.\d+" "v6.155.5"
  ` + "```" + `
- **View Replacement History**:
  ` + "```bash" + `
  gitmap replace history
  ` + "```" + `

---

### Workflow D: Workspace & Background Automation
- **Antigravity Workspaces**:
  ` + "```bash" + `
  gitmap agy add /path/to/repo       # Add project to Antigravity workspaces
  gitmap agy sync                     # Sync all agent configs
  gitmap agy stats                    # View token and run stats
  ` + "```" + `
- **VS Code Project Manager**:
  ` + "```bash" + `
  gitmap vsc add /path/to/repo        # Register repository in VS Code Project Manager
  gitmap vsc ls                       # List tracked workspaces
  ` + "```" + `
- **Background Scheduler & Daemons**:
  ` + "```bash" + `
  gitmap schedule add "nightly-sync" --interval "24h"
  gitmap schedule status
  ` + "```" + `

---

## 3. Alternative Commands & Aliases Cheat Sheet (Instead of Raw Git)

| Task | Standard Command | Fast Shortcut / Alias | Notes |
| :--- | :--- | :--- | :--- |
| **Pipeline Status** | ` + "`gitmap pipeline status --json`" + ` | ` + "`gitmap pipeline status`" + ` | Live execution state, ETA, run URL |
| **Get Wait Time** | ` + "`gitmap pipeline waittime`" + ` | ` + "`gitmap eta`" + ` | Returns integer seconds for timer |
| **Extract Error Logs**| ` + "`gitmap pipeline error-logs`" + ` | ` + "`gitmap error-logs`" + ` | Dumps failing CI step to file |
| **Full CI Logs** | ` + "`gitmap pipeline logs`" + ` | ` + "`gitmap logs`" + ` | Displays complete CI run stdout/stderr |
| **Find File (Exact)** | ` + "`gitmap find-files <name>`" + ` | ` + "`gitmap ff <name>`" + ` | Filters by name and -ext |
| **Find File (Contains)**| ` + "`gitmap find-files-any <str>`" + ` | ` + "`gitmap ffa <str>`" + ` | Matches substring in filename |
| **Find File (Prefix)**| ` + "`gitmap find-files-startswith`" + ` | ` + "`gitmap ffs <prefix>`" + ` | Matches prefix in filename |
| **Find File (Suffix)**| ` + "`gitmap find-files-endswith`" + ` | ` + "`gitmap ffe <suffix>`" + ` | Matches suffix (e.g. _test.go) |
| **Commit Feature** | ` + "`gitmap commit-push-feature`" + ` | ` + "`gitmap cpf \"<msg>\"`" + ` | Auto commits and pushes feature |
| **Commit Bugfix** | ` + "`gitmap commit-push-bug`" + ` | ` + "`gitmap cpb \"<msg>\"`" + ` | Auto commits and pushes bug fix |
| **Commit Release** | ` + "`gitmap commit-push-release`" + ` | ` + "`gitmap cpr \"<msg>\"`" + ` | Auto commits and pushes release chore |
| **Pull & Push** | ` + "`gitmap pull-commit-push`" + ` | ` + "`gitmap pcp \"<msg>\"`" + ` | Pulls latest, commits, and pushes |
| **Antigravity** | ` + "`gitmap antigravity`" + ` | ` + "`gitmap agy`, `gitmap ag`" + ` | Workspace & AI agent manager |
| **VS Code PM** | ` + "`gitmap vscode`" + ` | ` + "`gitmap vsc`" + ` | VS Code Project Manager sync |
| **Scheduler** | ` + "`gitmap schedule`" + ` | ` + "`gitmap sc`" + ` | Background task scheduler |
| **Repo Status** | ` + "`gitmap status`" + ` | ` + "`gitmap st`" + ` | Shows dirty/ahead/behind across repos |
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
