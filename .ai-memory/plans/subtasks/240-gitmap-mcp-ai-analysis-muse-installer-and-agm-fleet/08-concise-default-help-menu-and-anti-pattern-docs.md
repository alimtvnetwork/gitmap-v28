# Subtask 08: Concise Default Help Menu & Anti-Pattern Governance

> **Subtask ID:** Subtask-08  
> **Parent Plan:** `.ai-memory/plans/240-gitmap-mcp-ai-analysis-muse-installer-and-agm-fleet.md`  
> **Target Subsystems:** `cli/cmd/`, `02-spec/02-coding-guidelines/`, `03-ai-scripts/`  
> **Owned Files:**  
> - `cli/cmd/rootusage_concise.go`  
> - `cli/cmd/rootusage.go`  
> - `cli/cmd/root.go`  
> - `cli/cmd/rootusage_test.go`  
> - `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md`  
> - `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`  
> - `03-ai-scripts/40-check-anti-patterns.py`  

---

## 1. Concrete Objectives

1. **Concise Default Help Layout (< 25 Lines):**
   - Refactor `Run()` in `cli/cmd/root.go` to split bare execution (`len(os.Args) < 2`) from explicit help requests (`gitmap help`, `gitmap -h`, `gitmap --help`).
   - Implement `printConciseUsage()` in `cli/cmd/rootusage_concise.go` rendering:
     - Header: GitMap version, OS/Arch, and active commit hash.
     - Workspace Context: Active repository directory, current Git branch, and status badge (`● Clean` / `▲ Dirty`).
     - Core 6 Commands:
       * `gitmap aum search <query>` — Sub-millisecond indexed codebase search
       * `gitmap find <pattern>` — Locate files across repository instantly
       * `gitmap cpc "<module> - msg"` — Atomic chore commit with safety gates
       * `gitmap log [-n 20]` — High-speed cached commit history (<5ms)
       * `gitmap diff-branch [a] [b]` — Compare branch divergence and file changes
       * `gitmap repo create <name>` — Smart repository creation with clone fallback
     - Footer: `"💡 Tip: Run 'gitmap help' for full 80+ command reference, or 'gitmap <command> --help'."`
   - Total printed line count must not exceed 25 lines.
2. **Preservation of Comprehensive Full Help:**
   - Retain full categorized 380-line menu (`printFullUsage()`) when invoked via `gitmap help`, `gitmap --help`, or `gitmap -h`.
3. **Anti-Pattern Documentation & Guidelines:**
   - Author `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` detailing the strict prohibition of:
     - `Get-ChildItem -Filter` / `dir -Filter` -> mandate `gitmap find` / `gitmap scan`
     - `git grep` / `rg` / `ripgrep` / `grep` -> mandate `gitmap aum search`
     - `Remove-Item -Force` / `rm -rf` -> mandate `gitmap rm --task <id>`
   - Update `02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md` adding Mistake #16 (`Get-ChildItem -Filter` Shell Bypasses) and Mistake #17 (`git grep` instead of `gitmap aum search`).
4. **Automated Anti-Pattern Linter Script:**
   - Implement `03-ai-scripts/40-check-anti-patterns.py` to scan repository scripts, prompts, and plans for prohibited commands.
5. **Unit Testing:**
   - Implement `cli/cmd/rootusage_test.go` verifying that bare invocation renders strictly <= 25 lines and contains required workspace context and pointers.

---

## 2. Core Domain Types & Structs

```go
package cmd

// ConciseWorkspaceContext captures active repository telemetry for help display.
type ConciseWorkspaceContext struct {
	RepoName      string
	CurrentBranch string
	IsWorkTree    bool
	IsClean       bool
}

// CoreCommandPointer describes a primary command highlighted in the concise menu.
type CoreCommandPointer struct {
	Syntax      string
	Description string
}
```

---

## 3. Implementation Checklist

- [ ] **1. Concise Usage Renderer (`cli/cmd/rootusage_concise.go`):**
  - Implement `resolveConciseWorkspaceContext() ConciseWorkspaceContext`.
  - Use `gitutil.IsInsideWorkTree()` and `gitutil.CurrentBranch()` to query repository status.
  - Implement `printConciseUsage()` with strict line budgeting:
    - 2 lines: Header & spacing.
    - 2 lines: Workspace context & branch badge.
    - 1 line: Empty separator.
    - 1 line: "Core AI & Developer Commands:" section title.
    - 6 lines: Core command entries formatted with ANSI color alignment.
    - 1 line: Empty separator.
    - 2 lines: Footer notes and help instruction.
  - Total line budget: 15–18 lines (guaranteed <= 25 lines).
- [ ] **2. Entry Point Separation (`cli/cmd/root.go`, `cli/cmd/rootusage.go`):**
  - In `cli/cmd/root.go`:
    - When `len(os.Args) < 2`, call `printConciseUsage()` instead of `printUsage()`.
    - When `os.Args[1]` is `"help"`, `"-h"`, or `"--help"`, call `printFullUsage()`.
  - Rename legacy `printUsage()` in `cli/cmd/rootusage.go` to `printFullUsage()`.
- [ ] **3. Anti-Pattern Specification Document (`02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md`):**
  - Document prohibited commands, reasons for ban, and mandatory GitMap alternatives.
  - Provide concrete polyglot examples (PowerShell, Bash, Python, Go).
  - Define CI/CD gate requirements and regex scanning rules.
- [ ] **4. Update Common AI Mistakes (`02-spec/02-coding-guidelines/06-ai-optimization/04-common-ai-mistakes.md`):**
  - Add Mistake #16: Running `Get-ChildItem -Filter` instead of `gitmap find`.
  - Add Mistake #17: Running `git grep` / `ripgrep` instead of `gitmap aum search`.
  - Add Mistake #18: Running `Remove-Item -Force` instead of `gitmap rm --task <id>`.
- [ ] **5. Anti-Pattern Static Linter (`03-ai-scripts/40-check-anti-patterns.py`):**
  - Implement cross-platform Python linter scanning repository for:
    - `Get-ChildItem.*-Filter`
    - `git grep\b`
    - `\brg\s+` / `\bripgrep\s+`
    - `Remove-Item.*-Force`
  - Output violation counts and remediation guidance.
- [ ] **6. Help Output Line Count Verification Test (`cli/cmd/rootusage_test.go`):**
  - Capture stdout of `printConciseUsage()`.
  - Count newlines and assert total lines <= 25.
  - Assert presence of `gitmap aum search`, `gitmap find`, and `'gitmap help'`.

---

## 4. Acceptance Criteria

- [x] Executing bare `gitmap` prints fewer than 25 lines of output to stdout.
- [x] Executing bare `gitmap` displays the active repository name, current branch, and clean/dirty badge if run inside a Git repository.
- [x] Executing `gitmap help` prints the full comprehensive categorized catalog.
- [x] `02-spec/02-coding-guidelines/06-ai-optimization/10-anti-pattern-replacements.md` is published and referenced in coding guidelines.
- [x] `03-ai-scripts/40-check-anti-patterns.py` runs cleanly and reports zero false positives on compliant files.
- [x] All references use strictly relative Git paths.

---

## 5. Verification Commands

```bash
# Run unit tests
go test -v ./cli/cmd/... -run TestConciseHelpOutput

# Verify compilation
go build -v ./cli/...

# Test bare help command output line count
gitmap | Measure-Object -Line

# Test explicit help output
gitmap help

# Run anti-pattern linter
python 03-ai-scripts/40-check-anti-patterns.py
```
