# Pre-Flight Guard and Error Interceptor Specification

> **Task Slug:** `370-ai-execution-and-error-handling-instructions`  
> **Subsystem:** PowerShell Command Dispatcher (`cli/cmd/`)  
> **Date:** 31-Oct-26  
> **Status:** APPROVED  
> **AI Confidence:** Production-Ready  
> **Path Hygiene:** Strict Relative Paths Only  

---

## Table of Contents
1. [Executive Summary & Architectural Objective](#1-executive-summary--architectural-objective)
2. [Problem Statement & AI Anti-Pattern Analysis](#2-problem-statement--ai-anti-pattern-analysis)
3. [Two-Layer Defense Architecture](#3-two-layer-defense-architecture)
4. [Component Implementation Specifications](#4-component-implementation-specifications)
   - [4.1 `cli/cmd/cmd_powershell.go`](#41-clicmdcmd_powershellgo)
   - [4.2 `cli/cmd/powershell_advice.go`](#42-clicmdpowershell_advicego)
   - [4.3 `cli/cmd/powershell_interceptor.go`](#43-clicmdpowershell_interceptorgo)
   - [4.4 `cli/cmd/cmd_powershell_test.go`](#44-clicmdcmd_powershell_testgo)
5. [Path Detection Rules & Classification Engine](#5-path-detection-rules--classification-engine)
   - [5.1 Directory Target Detection](#51-directory-target-detection)
   - [5.2 Document File Detection (.md, .json, .txt, etc.)](#52-document-file-detection-md-json-txt-etc)
   - [5.3 Executable Script Allowlist](#53-executable-script-allowlist)
   - [5.4 Unknown Cmdlet & Token Syntax Classification](#54-unknown-cmdlet--token-syntax-classification)
6. [Formatting Contracts (Tokyo Night Diagnostic Card)](#6-formatting-contracts-tokyo-night-diagnostic-card)
   - [6.1 Tokyo Night Palette & Terminal Primitives](#61-tokyo-night-palette--terminal-primitives)
   - [6.2 Card Layout & ANSI Typography](#62-card-layout--ansi-typography)
   - [6.3 AI Diagnostic Inquiry Framework](#63-ai-diagnostic-inquiry-framework)
   - [6.4 Native Replacement Recommendation Matrix](#64-native-replacement-recommendation-matrix)
   - [6.5 3-Step Immediate Learning Action Guide](#65-3-step-immediate-learning-action-guide)
7. [Clean Error Bubbling & Stack Trace Suppression](#7-clean-error-bubbling--stack-trace-suppression)
   - [7.1 Leaf Error Sovereignty & AppError Architecture](#71-leaf-error-sovereignty--apperror-architecture)
   - [7.2 Abort Context & Suppression Mechanics (`reported: true`)](#72-abort-context--suppression-mechanics-reported-true)
   - [7.3 Integration with CLI Exit Handlers & DB Audit Logging](#73-integration-with-cli-exit-handlers--db-audit-logging)
   - [7.4 Coding Guideline Enforcements (Limits & Guard Clauses)](#74-coding-guideline-enforcements-limits--guard-clauses)
8. [Cross-Reference & Integration Contracts](#8-cross-reference--integration-contracts)
9. [Acceptance Criteria](#9-acceptance-criteria)

---

## 1. Executive Summary & Architectural Objective

Autonomous AI agents frequently execute CLI commands to inspect project status, explore directory hierarchies, and read subtask specifications. In environments where GitMap provides the `gitmap ps` command (a wrapper around PowerShell Core `pwsh` or Windows PowerShell `powershell.exe`), agents repeatedly confuse `gitmap ps` with directory listing (`ls`, `dir`), process inspection (`ps aux`), or task inspection (`gitmap agent subtask ls`).

When an agent mistakenly invokes:
```bash
gitmap ps .ai-memory/plans/subtasks/240-running-and-queued-prompt-detection-in-multi-instance
```
PowerShell attempts to interpret the directory path as an executable script or cmdlet name, terminating with an unhelpful error:
```text
.ai-memory/plans/...: The term '.ai-memory/plans/...' is not recognized as a name of a cmdlet...
gitmap: [E9000:EXECUTION] execute powershell: exit status 1
```
This failure often triggers a cascade of noisy Go runtime stack traces, unformatted terminal spam, and repetitive agent retries that consume context window and pollute audit logs.

This specification defines the architectural design, component contracts, path classification algorithms, Tokyo Night card formatting, and clean error bubbling protocols for:
1. `cli/cmd/cmd_powershell.go`: Core execution flow with integrated pre-flight guard and exit interceptor.
2. `cli/cmd/powershell_advice.go`: Path classification, AI self-reflection inquiries, and Tokyo Night diagnostic card rendering.
3. `cli/cmd/powershell_interceptor.go`: Post-execution exit code normalization, unrecognized cmdlet capture, and stack trace suppression.

---

## 2. Problem Statement & AI Anti-Pattern Analysis

### 2.1 The Semantic Confusion Matrix
AI agents exhibit three distinct patterns of confusion regarding `gitmap ps`:

| Anti-Pattern | Agent Assumption | Reality in GitMap | Resulting Error |
|---|---|---|---|
| **Directory Navigation** | Believes `gitmap ps <path>` lists directory contents like `ls` or `dir`. | `gitmap ps` executes arguments as a PowerShell expression (`-Command <arg>`). | Exit code 1: Term not recognized as cmdlet. |
| **Document Reading** | Believes `gitmap ps <file.md>` renders markdown or reads file contents like `cat`. | PowerShell attempts to invoke the document as an executable script. | Exit code 1: Term not recognized or file association error. |
| **Process Status Confusion** | Believes `gitmap ps` functions like UNIX `ps` (displaying running OS processes). | `gitmap ps` launches PowerShell process execution, not process listing. | Unexpected shell launch or syntax failure. |

### 2.2 Cascading Failure Symptoms
Without an intelligent guard, the command failure produces:
1. **Raw Go Runtime Stack Traces:** Standard error wrappers bubble up as generic execution faults, printing full stack traces (`goroutine 1 [running]...`) if unhandled.
2. **Context Budget Depletion:** The agent ingests 40-80 lines of Go error text and PowerShell internal parser errors.
3. **Repeated Retries:** Lacking structured guidance on the proper native GitMap command, the agent frequently retries the exact same command or variations with invalid flags.

---

## 3. Two-Layer Defense Architecture

To provide zero-regression execution while guaranteeing immediate AI learning, GitMap deploys a two-layer defense architecture:

```
User / AI Invocation: `gitmap ps <args>`
               │
               ▼
┌────────────────────────────────────────────────────────┐
│ Layer 1: Pre-Flight Guard (Synchronous In-Memory)     │
│ - Strip modifier flags (--force-all, -f)               │
│ - Inspect candidate argument (cleanArgs[0])           │
│ - Stat filesystem: IsDirectory? IsDocumentFile?       │
│                                                        │
│ [MATCH DETECTED]               [NO INVALID TARGET]     │
│        │                               │               │
│        ▼                               ▼               │
│ Render Tokyo Night Card       Invoke PowerShell Core   │
│ Bubble E1000 Abort Error      `pwsh -NoProfile -c ...` │
│ Exit Code 2 (Clean)                    │               │
└────────────────────────────────────────┼───────────────┘
                                         ▼
┌────────────────────────────────────────────────────────┐
│ Layer 2: Post-Execution Interceptor (Exit Guard)      │
│ - Catch `exec.ExitError` on non-zero exit code        │
│ - Detect "is not recognized as a name of a cmdlet"    │
│ - Detect "Cannot find path"                           │
│                                                        │
│ [CMDLET / PATH NOT FOUND]       [VALID SCRIPT ERROR]   │
│        │                               │               │
│        ▼                               ▼               │
│ Render Advice Fallback        Wrap Exit Code as Abort  │
│ Bubble E9001 Abort Error      Suppress Stack Trace     │
│ Exit Code 2 (Clean)           Audit DB Logged (Clean)  │
└────────────────────────────────────────────────────────┘
```

1. **Layer 1: Pre-Flight Guard (Pre-Invocation)**  
   Evaluates arguments *before* any OS process is spawned. If the argument is an existing filesystem directory or a non-executable document file (`.md`, `.json`, `.txt`, `.csv`, `.yml`, etc.), execution halts instantly. Zero process overhead, sub-millisecond response.

2. **Layer 2: Post-Execution Interceptor (Post-Invocation)**  
   If the command passes pre-flight (e.g., dynamic PowerShell expression, quoted script, or unknown cmdlet name) but exits with non-zero status containing PowerShell cmdlet resolution errors, the interceptor catches the `*exec.ExitError`, renders the educational guidance card, and bubbles a structured abort error with `reported: true`.

---

## 4. Component Implementation Specifications

### 4.1 `cli/cmd/cmd_powershell.go`

#### Architectural Responsibilities
- Primary orchestrator for the `ps` CLI command.
- Discovers the PowerShell binary (`pwsh` on PATH, fallback to `powershell` on Windows).
- Strips administrative flags (`--force-all`, `-f`).
- Enforces pre-flight inspection before process launch.
- Dispatches execution to `powershell_interceptor.go`.

#### Go Function Signatures & Contracts
```go
// RunPowerShell coordinates pre-flight checks, executable lookup, and process execution.
// Conforms strictly to leaf error sovereignty: returns error, never calls cliexit or os.Exit.
func RunPowerShell(args []string) error

// inspectPowerShellPreflight analyzes cleaned arguments and halts invalid directory/document targets.
func inspectPowerShellPreflight(args []string) error

// executePowerShellWithInterceptor executes the process and delegates error post-processing.
func executePowerShellWithInterceptor(psPath string, cleanArgs []string) error
```

#### Refactored Control Flow
```go
func RunPowerShell(args []string) error {
    psPath, hasPS := findPowerShellExecutable()
    if !hasPS {
        printPowerShellMissingAdvice()

        return apperror.NewSimple("powershell not found", "E_POWERSHELL_NOT_FOUND")
    }

    cleanArgs := stripForceFlag(args)
    if len(cleanArgs) == 0 {
        return executeInteractivePowerShell(psPath)
    }

    if err := inspectPowerShellPreflight(cleanArgs); err != nil {
        return err
    }

    return executePowerShellWithInterceptor(psPath, cleanArgs)
}
```

---

### 4.2 `cli/cmd/powershell_advice.go`

#### Architectural Responsibilities
- Encapsulates path classification logic and target categorization.
- Renders high-fidelity Tokyo Night ANSI cards via `cli/termout`.
- Houses the AI Self-Reflection Questions, execution explanations, and native replacement recommendations.

#### Target Kind Enumeration
```go
type PowerShellTargetKind string

const (
    TargetKindDirectory     PowerShellTargetKind = "DIRECTORY"
    TargetKindDocument      PowerShellTargetKind = "DOCUMENT"
    TargetKindUnknownCmdlet PowerShellTargetKind = "UNKNOWN_CMDLET"
    TargetKindValidCommand  PowerShellTargetKind = "VALID_COMMAND"
)
```

#### Diagnostic Configuration Struct
```go
type PowerShellAdviceOptions struct {
    RawTarget    string
    TargetKind   PowerShellTargetKind
    ResolvedPath string
    ContextNote  string
}
```

#### Key Functions
```go
// detectPowerShellTarget evaluates a raw CLI token and determines its TargetKind.
func detectPowerShellTarget(token string) (PowerShellTargetKind, string)

// isDocumentFileExtension checks if an extension matches non-executable document formats.
func isDocumentFileExtension(ext string) bool

// isExecutableScriptExtension checks if an extension matches executable scripts (.ps1, .bat, .cmd, .exe).
func isExecutableScriptExtension(ext string) bool

// renderPowerShellDiagnosticCard builds and prints the Tokyo Night card to os.Stderr.
func renderPowerShellDiagnosticCard(opts PowerShellAdviceOptions)

// buildAdviceQuestions formats the AI reflection inquiries.
func buildAdviceQuestions() string

// buildNativeReplacements formats the command replacement matrix based on TargetKind.
func buildNativeReplacements(opts PowerShellAdviceOptions) string
```

---

### 4.3 `cli/cmd/powershell_interceptor.go`

#### Architectural Responsibilities
- Executes the configured `*exec.Cmd` while buffering or intercepting standard streams.
- Intercepts process exit codes and error output.
- Matches PowerShell cmdlet failure signatures.
- Constructs standardized `apperror.AppError` instances with `apperror.ErrorTypeAbort` and `map[string]any{"reported": true}` to ensure clean exit code propagation without Go runtime stack traces.

#### Key Functions
```go
// executePowerShellProcess runs the command and passes errors through normalizePowerShellExitError.
func executePowerShellProcess(psPath string, cleanArgs []string) error

// normalizePowerShellExitError inspects run errors and returns clean abort AppErrors.
func normalizePowerShellExitError(runErr error, cleanArgs []string) error

// hasPowerShellCmdletNotFoundError inspects standard error text for missing cmdlet patterns.
func hasPowerShellCmdletNotFoundError(errOutput string) bool

// buildPowerShellAbortError creates a structured abort AppError with reported: true.
func buildPowerShellAbortError(code int, msg string) error
```

---

### 4.4 `cli/cmd/cmd_powershell_test.go`

#### Test Coverage Requirements
1. `TestPowerShell_PreflightDirectory`: Verifies that passing an existing directory (e.g. `.ai-memory/plans` or temporary test directory) halts immediately without calling `exec.Command` and returns an abort error with `reported: true`.
2. `TestPowerShell_PreflightDocumentFile`: Verifies that passing a `.md`, `.json`, or `.txt` file halts with `TargetKindDocument` diagnostic output.
3. `TestPowerShell_AllowValidPs1Script`: Verifies that passing a `.ps1` script file bypasses the pre-flight guard and proceeds to execution.
4. `TestPowerShell_AllowPowerShellExpression`: Verifies that commands like `gitmap ps "Get-ChildItem -Path '.'"` pass pre-flight inspection cleanly.
5. `TestPowerShell_CleanErrorBubbling`: Verifies that intercepted errors produce `apperror.ErrorTypeAbort`, prevent stack trace generation in `cliexit.WriteAppErrorReport`, and carry exit code `2`.

---

## 5. Path Detection Rules & Classification Engine

### 5.1 Directory Target Detection

An argument candidate is classified as `TargetKindDirectory` if:
1. `candidate` is non-empty and does not begin with a command flag (`-` or `--`).
2. `cleanCandidate := strings.Trim(candidate, "\"'")` is evaluated.
3. `os.Stat(cleanCandidate)` succeeds without error AND `fi.IsDir()` is `true`.
4. **Normalized Fallback:** If `cleanCandidate` contains path separators (`/` or `\`), or starts with `.`, `./`, `../`, or drive letters (`C:`), and exists on the filesystem as a directory.

```go
func isDirectoryTarget(token string) bool {
    clean := strings.Trim(token, "\"'")
    if clean == "" {
        return false
    }

    fi, err := os.Stat(clean)
    if err != nil {
        return false
    }

    return fi.IsDir()
}
```

### 5.2 Document File Detection (.md, .json, .txt, etc.)

An argument candidate is classified as `TargetKindDocument` if:
1. `os.Stat(cleanCandidate)` succeeds without error AND `fi.IsDir()` is `false`.
2. The file extension (converted to lowercase) matches the document format list:
   - Documentation & Text: `.md`, `.markdown`, `.txt`, `.rtf`, `.doc`, `.docx`
   - Structured Data: `.json`, `.jsonl`, `.yaml`, `.yml`, `.toml`, `.xml`, `.csv`, `.tsv`
   - Source Specifications: `.sql`, `.graphql`, `.html`, `.css`, `.log`
3. The extension is NOT an executable script or binary format.

```go
var documentExtensions = map[string]bool{
    ".md": true, ".markdown": true, ".txt": true,
    ".json": true, ".jsonl": true, ".csv": true, ".tsv": true,
    ".yml": true, ".yaml": true, ".toml": true, ".xml": true,
    ".sql": true, ".html": true, ".css": true, ".log": true,
}

func isDocumentFileTarget(token string) bool {
    clean := strings.Trim(token, "\"'")
    fi, err := os.Stat(clean)
    if err != nil || fi.IsDir() {
        return false
    }

    ext := strings.ToLower(filepath.Ext(clean))

    return documentExtensions[ext]
}
```

### 5.3 Executable Script Allowlist

The following script extensions are recognized as valid executable targets and MUST NOT be blocked by pre-flight:
- `.ps1`: PowerShell script
- `.psm1`: PowerShell module
- `.bat`: Windows Batch script
- `.cmd`: Windows Command script
- `.exe`: Windows executable binary
- `.sh`: Shell script (when executing in bash-compatible environments)

```go
func isExecutableScript(token string) bool {
    ext := strings.ToLower(filepath.Ext(strings.Trim(token, "\"'")))

    return ext == ".ps1" || ext == ".psm1" || ext == ".bat" || ext == ".cmd" || ext == ".exe"
}
```

### 5.4 Unknown Cmdlet & Token Syntax Classification

If an argument fails pre-flight (e.g., target does not exist as a local directory or document), but post-execution produces standard error output containing any of the following patterns:
1. `The term '...' is not recognized as a name of a cmdlet`
2. `is not recognized as the name of a cmdlet`
3. `The term '...' is not recognized as a function, script file, or operable program`
4. `Cannot find path '...' because it does not exist`

The post-execution interceptor classifies the error as `TargetKindUnknownCmdlet` and triggers the fallback educational card.

---

## 6. Formatting Contracts (Tokyo Night Diagnostic Card)

### 6.1 Tokyo Night Palette & Terminal Primitives

All visual cards adhere strictly to GitMap's Tokyo Night design system (`cli/termout`):

| Palette Token | Hex Code | Terminal Usage in Card |
|---|---|---|
| `TokyoBase` | `#1A1B26` | Card background fill (where supported) |
| `TokyoSurface`| `#24283B` | Card inner pill background |
| `TokyoPrimary`| `#7AA2F7` | Card title, borders, section headers |
| `TokyoAccent` | `#BB9AF7` | Command examples, highlighted code tokens |
| `TokyoOK`     | `#9ECE6A` | Recommended replacement commands |
| `TokyoWarn`   | `#E0AF68` | AI reflection questions, caution indicators |
| `TokyoError`  | `#F7768E` | Invalid target indicator, error status |
| `TokyoMuted`  | `#565F89` | Dividing lines, explanatory metadata notes |
| `TokyoFG`     | `#C0CAF5` | Standard body prose and descriptions |

### 6.2 Card Layout & ANSI Typography

The card is enclosed within `termout.TokyoBoxStyle` with clear, scannable visual hierarchy:

```text
╭─────────────────────────────────────────────────────────────────────────────╮
│  ⚡ GITMAP AI GUARD: Invalid PowerShell Target Detected                     │
│                                                                             │
│  Target:  .ai-memory/plans/subtasks/240-running-and-queued-prompt...        │
│  Class:   DIRECTORY (Filesystem folder cannot be executed as script)       │
├─────────────────────────────────────────────────────────────────────────────┤
│  🤔 AI Reflection Questions:                                                │
│  1. Why did you invoke 'gitmap ps <path>' on this target?                   │
│  2. Did you mistake 'gitmap ps' for directory listing or task inspection?   │
│  3. What was your intended outcome (list files, read doc, run script)?      │
├─────────────────────────────────────────────────────────────────────────────┤
│  ℹ️ What 'gitmap ps' Actually Does:                                         │
│  • Executes commands via PowerShell Core ('pwsh -NoProfile -Command ...').  │
│  • Expects valid PowerShell expressions, cmdlets, or executable scripts.    │
│  • Does NOT navigate folders, display directory contents, or view files.    │
├─────────────────────────────────────────────────────────────────────────────┤
│  💡 Recommended Native Replacements:                                        │
│  • List Directory:     gitmap dir <path>   OR   gitmap ls <path>           │
│  • Inspect Subtasks:   gitmap agent subtask ls                              │
│  • Read File / Spec:   gitmap cat <filepath>   OR   gitmap heatmap <file>   │
│  • PowerShell Query:   gitmap ps "Get-ChildItem -Path '<path>'"             │
├─────────────────────────────────────────────────────────────────────────────┤
│  🚀 3-Step Immediate Learning Action:                                       │
│  Step 1: Identify your real intent (list directory vs execute script).     │
│  Step 2: Use the specialized native GitMap command shown above.             │
│  Step 3: If running PowerShell, enclose full expression in quotes ("...").  │
╰─────────────────────────────────────────────────────────────────────────────╯
```

### 6.3 AI Diagnostic Inquiry Framework

Every intercepted error MUST prominently ask the AI agent to self-reflect before making its next tool call:
1. **Intention Check:** `"Why did you run this? What was your expected outcome?"`
2. **Mental Model Correction:** Explaining why a folder or markdown document cannot be executed by an operating system shell.
3. **Hypothesis Verification:** Asking the agent whether it intended to list files, read content, or query tasks.

### 6.4 Native Replacement Recommendation Matrix

Based on `TargetKind`, the card presents tailored replacements:

```go
func getReplacementsForTarget(kind PowerShellTargetKind, target string) []ReplacementOption {
    switch kind {
    case TargetKindDirectory:
        return []ReplacementOption{
            {Command: "gitmap dir " + target, Description: "Fast directory tree listing"},
            {Command: "gitmap ls " + target, Description: "Standard directory listing"},
            {Command: "gitmap agent subtask ls", Description: "Inspect pending and active subtasks"},
            {Command: fmt.Sprintf("gitmap ps \"Get-ChildItem -Path '%s'\"", target), Description: "PowerShell directory enumeration"},
        }
    case TargetKindDocument:
        return []ReplacementOption{
            {Command: "gitmap heatmap " + target, Description: "Read file with cognitive Heatmap touch logging"},
            {Command: "gitmap cat " + target, Description: "Fast file inspection with Tokyo Night rendering"},
            {Command: "gitmap cat " + target + " -L 1:50", Description: "Read bounded line range"},
        }
    default:
        return []ReplacementOption{
            {Command: "gitmap help", Description: "View all native GitMap commands"},
            {Command: fmt.Sprintf("gitmap ps \"%s\"", target), Description: "Enclose PowerShell command in quotes"},
        }
    }
}
```

### 6.5 3-Step Immediate Learning Action Guide

The bottom of the card provides a 3-step action guide:
- **Step 1:** Stop and parse the AI Reflection Questions.
- **Step 2:** Select the exact replacement command from the list above.
- **Step 3:** Avoid raw unquoted commands when passing expressions to `gitmap ps`.

---

## 7. Clean Error Bubbling & Stack Trace Suppression

### 7.1 Leaf Error Sovereignty & AppError Architecture

Under GitMap coding guidelines (`02-spec/02-coding-guidelines/` and `.agents/rules/gitmap_guidelines.md`):
- Leaf functions (`RunPowerShell`, `inspectPowerShellPreflight`, `executePowerShellWithInterceptor`) MUST return errors directly.
- Leaf functions MUST NEVER call `cliexit.HandleError` or `os.Exit`.
- Only root command dispatchers in `cli/cmd/root.go` handle errors and call exit handlers.

### 7.2 Abort Context & Suppression Mechanics (`reported: true`)

To completely prevent Go runtime stack traces and redundant error logs, errors generated by pre-flight or post-execution interception MUST be constructed as:

```go
func buildPreflightAbortError(target string, kind PowerShellTargetKind) error {
    msg := fmt.Sprintf("invalid powershell target %s: %s cannot be executed directly", target, strings.ToLower(string(kind)))

    return apperror.NewWithDetails(
        "powershell",
        "E1000",
        msg,
        "cmdpowershell",
        apperror.ErrorTypeAbort,
        apperror.SeverityWarn,
        map[string]any{
            "reported": true,
            "target":   target,
            "kind":     string(kind),
            "exitCode": int(constants.ExitUsageError),
        },
    )
}
```

#### How Stack Trace Suppression Operates in `cliexit/handle.go`:
```go
// From cli/cliexit/handle.go
func isStackTraceEnabled(e *apperror.AppError) bool {
    // If GITMAP_DEBUG is explicitly enabled, developer receives stack trace.
    if os.Getenv("GITMAP_DEBUG") == "1" || os.Getenv("DEBUG") == "1" {
        return true
    }

    // Suppress stack trace when Type is ErrorTypeAbort or reported context is true.
    if e.Type == apperror.ErrorTypeAbort || isReportedContext(e.Ctx) {
        return false
    }

    return e.Severity == apperror.SeverityFatal
}
```
When `reported: true` is set, `cliexit.HandleError`:
1. Executes output pipe flushers (`runFlushers()`).
2. Directly exits with `exitFunc(code)` without writing stack traces or internal diagnostic dumps.

### 7.3 Integration with CLI Exit Handlers & DB Audit Logging

Errors bubble up to `root.go` where `finishCommandAudit` logs command execution in SQLite (`repodb/pipeline.db` and audit tables):
- Command: `ps`
- Arguments: `[<target>]`
- Status: `"failed"`
- Exit Code: `2` (`constants.ExitUsageError`)
- Error Code: `"E1000"` (Validation / Usage Error)
- Audit Record: Preserved with full telemetry and zero data corruption.

### 7.4 Coding Guideline Enforcements (Limits & Guard Clauses)

All implemented Go code adheres to:
1. **15-Line Function Limit:** Every function must be <= 15 lines (8-12 lines preferred). Larger workflows are decomposed into helper functions.
2. **Zero Nested Ifs:** All conditionals use early return guard clauses.
3. **Boolean Naming:** Booleans use `is` or `has` prefixes only (`hasPS`, `isDir`, `isDocument`, `hasUnrecognizedCmdlet`, `isWin`). Negative booleans are strictly forbidden.
4. **Spacing Rules:** Exactly one blank line after `}` and one blank line before `return`.
5. **No Magic Integers:** Always use `constants.ExitUsageError` (2) or strongly-typed enums.

---

## 8. Cross-Reference & Integration Contracts

| Component / File | Relationship | Role in Architecture |
|---|---|---|
| `cli/cmd/cmd_powershell.go` | Core CLI File | Entry point for `gitmap ps`; calls pre-flight and interceptor. |
| `cli/cmd/powershell_advice.go` | New Helper File | Path detection rules, reflection inquiries, and Tokyo Night card rendering. |
| `cli/cmd/powershell_interceptor.go` | New Helper File | Process runner, exit code capture, and `reported: true` error generation. |
| `cli/cmd/cmd_powershell_test.go` | Test Suite | Unit tests verifying pre-flight, card output, and clean error bubbling. |
| `cli/termout/theme_lipgloss.go` | UI Subsystem | Provides `TokyoBoxStyle`, `TokyoHeader`, `TokyoPill`, and Tokyo Night palette. |
| `cli/cliexit/handle.go` | Exit Orchestrator | Suppresses stack traces for `ErrorTypeAbort` with `reported: true`. |
| `.agents/skills/gitmap/skill.md` | Skill System | Documents native `ps` rules and anti-patterns for AI agents. |
| `02-spec/21-app/370-ai-execution-and-error-handling-instructions/readme.md` | Module Master | Index for 370 specification bundle. |

---

## 9. Acceptance Criteria

- [ ] **AC-1 (Pre-Flight Directory Interception):** Executing `gitmap ps <existing_dir>` intercepts before process launch, renders the Tokyo Night card to `os.Stderr`, and returns exit code 2 without spawning `pwsh`.
- [ ] **AC-2 (Pre-Flight Document Interception):** Executing `gitmap ps <file.md>` or `gitmap ps <file.json>` intercepts before process launch and recommends `gitmap heatmap <file>` or `gitmap cat <file>`.
- [ ] **AC-3 (Valid PowerShell Execution):** Executing `gitmap ps "Get-Date"` or `gitmap ps script.ps1` executes normally through PowerShell without interruption.
- [ ] **AC-4 (AI Reflection Prompts):** The rendered card explicitly contains the 3 AI reflection questions asking why the command was run and what was expected.
- [ ] **AC-5 (Clean Error Bubbling):** No Go runtime stack traces (`goroutine 1...`, `runtime/debug.Stack`) appear on stderr. Exit code is strictly `2`.
- [ ] **AC-6 (Audit Log Preservation):** The failure is captured cleanly in SQLite command audit logs with exit code 2.
- [ ] **AC-7 (Coding Guidelines Compliance):** All Go functions in `cmd_powershell.go`, `powershell_advice.go`, and `powershell_interceptor.go` are <= 15 lines, contain zero nested if statements, and use `is*`/`has*` boolean naming.
- [ ] **AC-8 (Strict Relative Paths):** All paths in code, tests, and specifications use relative paths only.
- [ ] **AC-9 (Unit Test Verification):** `go test -v ./cli/cmd -run TestPowerShell` passes with 100% green status.
