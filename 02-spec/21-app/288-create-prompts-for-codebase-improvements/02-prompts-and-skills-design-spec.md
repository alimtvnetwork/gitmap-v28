# Specification 288: Prompts & Skills Detailed Design Specification

- **Sequence:** 288
- **Slug:** `288-create-prompts-for-codebase-improvements`
- **File:** `02-spec/21-app/288-create-prompts-for-codebase-improvements/02-prompts-and-skills-design-spec.md`
- **Status:** APPROVED & SPECIFIED

---

## 1. Prompt 1 Design: Find 100 Improvements in the UI Level

### 1.1 Core Objective
Instruct an autonomous AI agent to conduct a systematic, grounded audit of the terminal TUI, Unicode tables, progress indicators, status formatters, and web visualizers, producing a structured inventory of 100 UI improvements.

### 1.2 Mandatory Item Schema
Each item in the 100-item output must adhere to the following schema:
```markdown
### Item [NNN]: [Descriptive Improvement Title]
- **Target Component / File:** [Relative path, e.g. `cli/termout/themes.go` or `cli/tui/browser.go`]
- **Visual Defect / UX Rough Edge:** [What is currently broken, low contrast, misaligned, or poorly formatted]
- **What It Fixes:** [Specific issue addressed, e.g. eliminating ANSI escape line wrapping, fixing colorblind contrast]
- **How It Improves the UI:** [Concrete visual enhancement, ergonomic benefit, and hierarchy improvement]
- **Color Tokens & Swatches:**
  - Role: `RolePrimary` / `RoleAccent` / `RoleWarn` / `RoleError` / `RoleOK` / `RoleMuted`
  - Hex: `#XXXXXX`
  - ANSI SGR Code: `\033[38;5;NNNm` or `\033[38;2;R;G;Bm`
  - Lipgloss Token: `lipgloss.AdaptiveColor{Light: "...", Dark: "..."}`
- **Concrete Layout / TUI Implementation Pattern:** [Code snippet showing Lipgloss style or Bubble Tea layout]
```

### 1.3 10 Discovery Dimensions (10 items each = 100 total):
1. Terminal ANSI & Lipgloss Color Palette Authority
2. Progress Bars, Spinners & Multi-Threaded TUI Gauges
3. Interactive Fuzzy Selectors, Multi-Select Menus & Keybindings
4. Unicode Tables, Dynamic Column Auto-Padded Folding & Sorting
5. Hierarchical Tree Visualizers (Branching, Merge Diffs & Heatmaps)
6. Modern Database Inspection Drawer & Data Badges (REST / SQLite Visualizer)
7. Error Telemetry, Formatted Stack Traces & 4-Part RCA Diagnostics HUD
8. Multi-Node SSH Fleet Dashboard & Node Ping Connectivity HUD
9. Real-Time Cognitive Heatmap Timeline & AI Agent Activity Stream
10. Desktop Notifications, Audio Chimes & Terminal Sound HUD

---

## 2. Prompt 2 Design: Find 100 Improvements in the Codebase (Coding Guidelines & Error Management)

### 2.1 Core Objective
Instruct an autonomous AI agent to audit the codebase for 100 surgical code improvements strictly enforcing grounded coding guidelines and error management rules.

### 2.2 Mandatory Item Schema
```markdown
### Item [NNN]: [Refactoring Target Title]
- **File Path & Lines:** [Relative path and line range, e.g. `cli/cmdpull/pull_remediation_ssh.go:105-113`]
- **Function / Symbol Name:** [Exact function signature, e.g. `func updateRepoTransportInDB(repoPath, sshURL string)`]
- **Existing Guideline Violation:** [Exact citation from `02-spec/02-coding-guidelines/` or `02-spec/03-error-manage/`]
  - Error Swallowing (`_ = ...`, silent void return)
  - Function Length Cap (> 15 lines)
  - Nested Conditionals (nested `if` instead of guard clauses)
  - Negative or Non-Conforming Boolean (`can*`, `should*`, `isNot*`)
  - Direct SQL in Business Logic (violating Split-DB repository encapsulation)
  - Raw Magic String / Integer instead of typed `*Type` enum
- **Context & Architectural Rationale:** [Why this violation causes defects, unhandled states, or audit gaps]
- **Proposed Refactoring (Before vs. After):**
  - Before: [Original offending snippet]
  - After: [Clean, compliant replacement snippet]
```

### 2.3 10 Guideline Focus Categories (10 items each = 100 total):
1. Error Return Sovereignty & AppError Enveloping
2. Function Length Reduction (<= 8 lines preferred, <= 15 lines max)
3. Guard Clauses & Inverted Condition Flattening (Zero Nested Ifs)
4. Boolean Logic Modernization (`is*` and `has*` only, zero negative booleans)
5. Split-DB Repository Pattern Encapsulation (Zero Raw SQL in Business Logic)
6. Strongly-Typed Enums ending in `Type` & Constants Centralization
7. Git Command Centralization (`gitutil.RunGitInDir` DRY execution)
8. Struct Value Semantics & Pointer Elimination
9. Memory-Safe Subprocess Spawning & Environment Sanitization
10. Database Connection Pool & WAL Mode Governance

---

## 3. Prompt 3 Design: Find 100 Improvements in the Whole System (UX & Safety)

### 3.1 Core Objective
Instruct an autonomous AI agent to discover 100 whole-system improvements that enhance user experience, developer velocity, and operational safety without introducing regressions or breaking backward compatibility.

### 3.2 Mandatory Item Schema
```markdown
### Item [NNN]: [Systemic Improvement Title]
- **Subsystem & Touchpoint:** [e.g. `cli/cmdpull/` batch pull or `cli/cmdssh/` remote execution]
- **User Experience Friction Point / Symptom:** [Exact pain point encountered by the developer]
- **Root Cause Analysis (Why it happens):** [Underlying technical or design cause]
- **Systemic Effect & Impact on User Workflow:** [Consequences on productivity, latency, or clarity]
- **Non-Destructive Implementation Blueprint:**
  - Safe Default Strategy (opt-in flags, auto-stash protections, non-destructive fallbacks)
  - Backward Compatibility Guarantee (how existing scripts/workflows continue working)
  - Proposed Technical Enhancement (exact mechanism, daemon, or command addition)
- **Safety Verification & Rollback Protocol:** [How to verify safety and roll back if issues arise]
```

### 3.3 10 Systemic Categories (10 items each = 100 total):
1. High-Speed Repository Discovery & Filesystem Walker
2. Parallel Batch Cloner, Round-Trip Restoration & Clonenow Engine
3. Multi-Node SSH Fleet Orchestration & Host Management
4. Cluster Task Delegation, Cross-Node Remote Execution & RPC
5. Three-Tier Split-DB Architecture & Schema Migrations
6. Typed SQL Engine, WAL Mode Concurrency & Query Optimization
7. CI/CD Pipeline Telemetry, Triage & Self-Healing RCA Engine
8. Git Credential Management, DPAPI Resilience & Headless SSH
9. Autonomous AI Agent Task Engine & Heatmap Telemetry
10. Interactive Macro Engine, Terminal UI Styling & Accidental Data Loss Prevention (ADLP)
