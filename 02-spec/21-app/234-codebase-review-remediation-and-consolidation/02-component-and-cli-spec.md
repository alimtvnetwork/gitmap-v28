# Component & CLI Specification: Codebase Review Remediation & Consolidation

> **Document ID:** `APP-SPEC-234-02`  
> **Plan Slug:** `234-codebase-review-remediation-and-consolidation`  
> **Specification Path:** `02-spec/21-app/234-codebase-review-remediation-and-consolidation/02-component-and-cli-spec.md`  
> **Master Ledger:** [00-master-audit-ledger.md](00-master-audit-ledger.md)  
> **Architecture Spec:** [01-architecture-spec.md](01-architecture-spec.md)  
> **Target Version:** `v6.497.0` (Preserved — zero unauthorized version churn)  
> **Status:** `APPROVED`  
> **Created Date:** `2026-10-06`  
> **Author:** MD ALIM UL KARIM / Multi-Agent Spec Engine  

---

## 1. Executive Scope & Component Boundaries

This specification defines the component and CLI-level architectural contracts for resolving the remaining actionable items identified in the codebase review (`docs/review/honest-feedback.md` and `docs/review/action-checklist.md`), specifically focusing on:

1. **Task-06: Clone Packages Mapping & Architectural Strategy:** Unify and map the five distinct clone packages (`cli/cloner`, `cli/clonefrom`, `cli/clonenow`, `cli/clonepick`, and `cli/clonenext`) under a cohesive architecture document (`docs/commands/cloning-architecture.md`) without collapsing them into a single monolithic file.
2. **Task-07: Error Exit Standardization:** Standardize CLI command exits across `cli/cmd/` by eliminating ad-hoc `cliexit.HandleError(nil, code)` invocations and transitioning to typed `*apperror.AppError` returns and explicit `cliexit` exit helpers (`HandleUsageError`, `HandleValidationError`, `HandleSuccess`).
3. **Task-08: Security Token Verification & `repo-secrets/` Non-Destructive Audit:** Catalog all 21 configuration and migration artifacts in `repo-secrets/`, establish the invariant of zero deletions without explicit user confirmation, and specify the automated security token verification protocol across state, logs, and caches.
4. **Task-10: Review Directory Removal & Objective AI Counter-Feedback:** Define the safe removal workflow for the temporary `docs/review/` directory and formulate an objective, grounded AI counter-review addressing News Spark's feedback.

---

## 2. Task-06: Clone Packages Mapping & Architectural Strategy

### 2.1 The Five Clone Engines & Modular Separation

The review suggested collapsing `cli/cloner`, `cli/clonefrom`, `cli/clonenow`, `cli/clonepick`, and `cli/clonenext` into a single engine. Per explicit architectural mandate, **monolithic single-file collapse is rejected**. Each package represents a specialized domain responsibility with distinct inputs, validation contracts, and runtime semantics.

```mermaid
flowchart TD
    subgraph GitMap CLI Clone Subsystem
        A[CLI Command Entrypoint] --> B{Clone Dispatcher}
        
        B -->|gitmap scan --clone| C[cli/cloner]
        B -->|gitmap clone-from file| D[cli/clonefrom]
        B -->|gitmap clone-now file| E[cli/clonenow]
        B -->|gitmap clone-pick url paths| F[cli/clonepick]
        B -->|gitmap cn / clone-next| G[cli/clonenext]
    end

    subgraph Data Models
        C --- M1[model.ScanRecord]
        D --- M2[clonefrom.Plan]
        E --- M3[clonenow.Plan]
        F --- M4[clonepick.Plan]
        G --- M5[clonenext.ParsedRepo]
    end

    subgraph Execution Characteristics
        C --> O1[Scan-Pipeline Cloner: In-memory DB-backed repo hierarchy]
        D --> O2[Plan-Driven Cloner: Arbitrary user-defined JSON/CSV plans]
        E --> O3[Round-Trip Cloner: GitMap scan output replay, preserving exact folder tree]
        F --> O4[Sparse-Checkout Cloner: Materialize selected paths & DB persistence]
        G --> O5[Version-Bumping Cloner: Increments vN, checks GitHub API, batch walks]
    end
```

### 2.2 Data Model Comparison & Responsibilities

| Package | Canonical Command | Core Data Model | Input Source | Primary Responsibility & Distinction |
| :--- | :--- | :--- | :--- | :--- |
| `cli/cloner` | `gitmap scan --clone` | `model.ScanRecord` | In-memory scan results / SQLite DB | Executes parallel clones directly from discovery scans. Tightly coupled with the database and scanner hierarchy. |
| `cli/clonefrom` | `gitmap clone-from <file>` | `clonefrom.Plan`<br>`clonefrom.Row` | External JSON or CSV file | Plan-driven batch cloner for arbitrary external repo lists. Validates rows, computes target branches, handles depth, and renders previews. |
| `cli/clonenow` | `gitmap clone-now <file>` | `clonenow.Plan`<br>`clonenow.Row` | Scan output (`gitmap.json`, `gitmap.csv`) | Re-runs clones from GitMap's own scan output. Preserves exact relative directory structures (`RelativePath`) and toggles HTTPS/SSH transports. |
| `cli/clonepick` | `gitmap clone-pick <url> <paths>` | `clonepick.Plan`<br>`clonepick.Result` | CLI arguments or DB replay | Sparse-checkout cloner. Fetches only targeted paths/subdirectories, supports interactive picking, and persists selections for `--replay`. |
| `cli/clonenext` | `gitmap cn [target]`<br>`gitmap clone-next` | `clonenext.ParsedRepo`<br>`clonenext.LocalRepoState` | Local repo or batch CSV/directory | Next-version bumping cloner. Parses versioned repo slugs (`repo-v1` -> `repo-v2`), interacts with GitHub API, and handles directory-wide batch migrations. |

### 2.3 Detailed Data Model Definitions

#### 1. `model.ScanRecord` (`cli/model/record.go`)
```go
type ScanRecord struct {
    ID                  int64    `json:"id" csv:"id"`
    Slug                string   `json:"slug" csv:"slug"`
    RepoID              string   `json:"repoId" csv:"repoId"`
    RepoName            string   `json:"repoName" csv:"repoName"`
    HTTPSUrl            string   `json:"httpsUrl" csv:"httpsUrl"`
    SSHUrl              string   `json:"sshUrl" csv:"sshUrl"`
    DiscoveredURL       string   `json:"discoveredUrl" csv:"discoveredUrl"`
    Branch              string   `json:"branch" csv:"branch"`
    BranchSource        string   `json:"branchSource" csv:"branchSource"`
    RelativePath        string   `json:"relativePath" csv:"relativePath"`
    AbsolutePath        string   `json:"absolutePath" csv:"absolutePath"`
    CloneInstruction    string   `json:"cloneInstruction" csv:"cloneInstruction"`
    Notes               string   `json:"notes" csv:"notes"`
    Depth               int      `json:"depth" csv:"depth"`
    Transport           string   `json:"transport" csv:"transport"`
    IdentifiedTransport string   `json:"identifiedTransport" csv:"identifiedTransport"`
}
```

#### 2. `clonefrom.Plan` (`cli/clonefrom/clonefrom.go`)
```go
type Plan struct {
    Source string // Source file path (JSON/CSV)
    Format string // "json" | "csv"
    Rows   []Row  // Validated rows to clone
}

type Row struct {
    URL      string // Clone source URL
    Dest     string // Destination relative to cwd
    Branch   string // Optional branch name
    Depth    int    // Optional shallow clone depth
    Checkout string // Working tree checkout mode ("auto", "skip", "force")
}
```

#### 3. `clonenow.Plan` (`cli/clonenow/clonenow.go`)
```go
type Plan struct {
    Source     string
    Format     string // "json" | "csv" | "text"
    Mode       string // "https" | "ssh"
    OnExists   string // "skip" | "update" | "force"
    Rows       []Row
    CoerceURL  func(string) string
    PersistURL func(string)
}

type Row struct {
    RepoName     string
    HTTPSUrl     string
    SSHUrl       string
    Branch       string
    RelativePath string
}
```

#### 4. `clonepick.Plan` (`cli/clonepick/clonepick.go`)
```go
type Plan struct {
    Name            string
    RepoCanonicalId string
    RepoUrl         string
    Mode            string // "https" | "ssh"
    Branch          string
    Depth           int
    Cone            bool // Sparse-checkout cone mode
    KeepGit         bool // Preserve .git directory
    DestDir         string
    Paths           []string // Repo-relative paths to fetch
    UsedAsk         bool
    DryRun          bool
    Quiet           bool
    Force           bool
    PreClonedSrc    string
}
```

#### 5. `clonenext.ParsedRepo` (`cli/clonenext/version.go`)
```go
type ParsedRepo struct {
    BaseName       string // e.g., "gitmap"
    CurrentVersion int    // e.g., 28
    HasVersion     bool   // true if name had -vN suffix
}

type LocalRepoState struct {
    OriginURL string // Extracted from .git/config
    HeadSHA   string // Resolved HEAD commit hash
}
```

### 2.4 Unification via Architecture Specification

Rather than collapsing the source code into a bloated monolithic file, the five engines are unified through a comprehensive developer-facing document: `docs/commands/cloning-architecture.md`.

This document specifies:
- When to use each command.
- How CLI flags map across the engines.
- The shared subprocess execution engine (`os/exec` wrappers, progress streaming, error handling).
- The distinction between workspace-preserving clones (`clonenow`, `cloner`), plan-driven clones (`clonefrom`), partial checkouts (`clonepick`), and version evolutions (`clonenext`).

---

## 3. Task-07: Error Exit Standardization

### 3.1 Problem Analysis: The Ad-Hoc Exit Anti-Pattern

An audit of the `cli/` codebase identified approximately 80 occurrences of:
```go
cliexit.HandleError(nil, code)
```

Passing `nil` to `HandleError` is an ad-hoc pattern where developers used `HandleError` purely as a process termination trigger. This introduces three major issues:
1. **Masked Diagnostics:** Passing `nil` bypasses the structured error envelope, leaving stderr empty or misleading without operational context.
2. **Conflated Intent:** `HandleError` was designed for error propagation, not clean process exits or usage warnings.
3. **Violated Guidelines:** Coding guidelines require all errors to flow through typed `*apperror.AppError` wrappers with explicit `Code`, `Op`, `Type`, and `Severity`.

### 3.2 Target Architecture & Typed Exit Helpers

The error exit path is standardized into three distinct tiers:

```mermaid
flowchart TD
    A[CLI Command Exit Site] --> B{Exit Reason}
    
    B -->|Clean / Intentional Exit| C[cliexit.Exit code or cliexit.HandleSuccess]
    B -->|CLI Flag / Input Validation Failure| D[cliexit.HandleUsageError / HandleValidationError]
    B -->|Domain / Runtime Execution Error| E[cliexit.HandleAppError with *apperror.AppError]
    
    C --> F[Flush Output Buffers & Terminate Process]
    D --> G[Format E1000 Diagnostic Envelope & Terminate]
    E --> H[Format Structured [Code:Type] Diagnostic & Terminate]
```

### 3.3 Helper Signatures in `cli/cliexit/handle.go`

To eliminate `HandleError(nil, code)`, `cli/cliexit/handle.go` is augmented with dedicated, typed exit functions:

```go
// HandleSuccess flushes stdout/stderr buffers and terminates with exit code 0.
func HandleSuccess() {
    Exit(0)
}

// HandleUsageError creates a structured validation error and exits with usage error code.
func HandleUsageError(command, message string, defaultCode ...int) {
    appErr := apperror.NewValidationError(message).
        WithContext("command", command).
        WithContext("category", "usage")
    code := resolveExitCode(defaultCode...)
    if len(defaultCode) == 0 {
        code = int(ExitCodeUsageError)
    }
    dispatchError(appErr, code)
}

// HandleValidationError creates an AppError with op, subject, and reason, then exits.
func HandleValidationError(command, op, subject, message string, defaultCode ...int) {
    appErr := apperror.NewValidation(op, "E1000", message).
        WithContext("command", command).
        WithContext("subject", subject)
    code := resolveExitCode(defaultCode...)
    if len(defaultCode) == 0 {
        code = int(ExitCodeValidationFailed)
    }
    dispatchError(appErr, code)
}

// HandleAppError explicitly handles an AppError instance with full diagnostics.
func HandleAppError(appErr *apperror.AppError, defaultCode ...int) {
    HandleError(appErr, defaultCode...)
}
```

### 3.4 Migration Rules for Existing Call Sites

| Current Pattern | Replacement Pattern | Semantics |
| :--- | :--- | :--- |
| `cliexit.HandleError(nil, 0)` | `cliexit.HandleSuccess()` or `cliexit.Exit(0)` | Clean successful termination with drained output pipes. |
| `cliexit.HandleError(nil, 1)` (on flag error) | `cliexit.HandleUsageError("scan", "missing required --root flag")` | Structured validation error (`E1000:validation`) with usage code. |
| `cliexit.HandleError(nil, 2)` (validation failure) | `cliexit.HandleValidationError("clone-from", "parse", path, "invalid row")` | Structured validation error with operation and subject context. |
| `cliexit.HandleError(err, code)` (where `err` is raw) | `cliexit.HandleError(apperror.WrapSimple(err, op), code)` | Automatic wrapping in typed `*apperror.AppError` before dispatch. |

---

## 4. Task-08: Security Token Verification & `repo-secrets/` Non-Destructive Audit

### 4.1 Invariant: Non-Destructive Audit & Zero Deletions

The codebase review questioned the presence of `repo-secrets/` at the repository root. Per explicit user mandate:
> *"Remove the committed artifact from the repo secret. I don't know what kind of artifacts have been committed. You can ask me so that we can confirm it. Do not remove anything unless we confirm from the repo secrets. Audit the repo secrets if there is anything. Repo secrets will contain secrets. No worries on this. Verify security token part."*

**Non-Negotiable Invariant:**
No file within `repo-secrets/` shall be deleted, truncated, or removed during autonomous execution. Any modification or deletion requires explicit user confirmation.

### 4.2 Complete Catalog of 21 `repo-secrets/` Artifacts

The `repo-secrets/` directory serves as an operational companion vault for multi-workstation migration, fleet synchronization, and Cursor/Antigravity profile backups. Below is the full catalog of the 21 tracked files:

| # | Relative Path | Type | Size | Operational Function & Classification |
| :---: | :--- | :--- | :--- | :--- |
| 1 | `repo-secrets/readme.md` | Doc | 8.7 KB | Master vault documentation, folder naming rules (`xx-*`), and `.ai-memory/` separation protocol. |
| 2 | `repo-secrets/04-ubuntu-migration/cursor-profile-transfer-notes.md` | Doc | 2.1 KB | Architectural notes for migrating Cursor profiles across Ubuntu fleet nodes. |
| 3 | `repo-secrets/04-ubuntu-migration/step-by-step-transition-guide.md` | Doc | 3.4 KB | Step-by-step operational transition runbook for fleet workstations. |
| 4 | `repo-secrets/04-ubuntu-migration/cursor-fleet-status.json` | Manifest | 1.8 KB | Machine state manifest tracking Cursor profile sync status across cluster nodes. |
| 5 | `repo-secrets/05-scripts/sync-cursor-profile.sh` | Shell | 2.8 KB | Bash script to sync Cursor settings and extensions to Ubuntu workstations. |
| 6 | `repo-secrets/05-scripts/setup-cursor-ubuntu.sh` | Shell | 3.1 KB | Ubuntu workstation setup automation for developer environment provisioning. |
| 7 | `repo-secrets/05-scripts/sync-cursor-profile.py` | Python | 4.2 KB | Python cross-platform runner for Cursor configuration synchronization. |
| 8 | `repo-secrets/05-scripts/heal-u1-pull-errors.py` | Python | 2.6 KB | Automated recovery script for resolving git pull divergence on Ubuntu Node 1. |
| 9 | `repo-secrets/05-scripts/migrate-cursor-memories-conversations.sh` | Shell | 1.9 KB | Shell migration utility for transferring Cursor chat memories and workspace SQLite DBs. |
| 10 | `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py` | Python | 3.8 KB | Python migration engine for Cursor conversation SQLite database schemas. |
| 11 | `repo-secrets/05-scripts/setup-cursor-ubuntu.py` | Python | 3.5 KB | Python installer for configuring extensions and dock settings on Ubuntu. |
| 12 | `repo-secrets/05-scripts/heal-u1-pull-errors.sh` | Shell | 1.4 KB | Shell wrapper for the Ubuntu Node 1 pull divergence healer. |
| 13 | `repo-secrets/09-antigravity-backup/readme.md` | Doc | 1.5 KB | Antigravity IDE backup architecture specification and vault restore instructions. |
| 14 | `repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.sh` | Shell | 1.7 KB | Verification script validating integrity of Antigravity backup archives. |
| 15 | `repo-secrets/09-antigravity-backup/scripts/verify-antigravity-backup.ps1` | PS1 | 2.2 KB | PowerShell verification script for Windows Antigravity backup validation. |
| 16 | `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.sh` | Shell | 2.9 KB | Automated shell script to restore Antigravity settings, skills, and plugins. |
| 17 | `repo-secrets/09-antigravity-backup/scripts/restore-antigravity-all.ps1` | PS1 | 3.4 KB | PowerShell restore script for Antigravity environment restoration. |
| 18 | `repo-secrets/09-antigravity-backup/vault/settings-manifest.json` | Manifest | 1.2 KB | Structured configuration manifest of Antigravity global settings. |
| 19 | `repo-secrets/09-antigravity-backup/vault/plugins-and-skills.json` | Manifest | 4.6 KB | Inventory of active Antigravity skills, tools, and extensions. |
| 20 | `repo-secrets/09-antigravity-backup/vault/projects-manifest.json` | Manifest | 3.1 KB | Workspace project catalog linking URIs to CorpusNames across nodes. |
| 21 | `repo-secrets/09-antigravity-backup/vault/pinned-projects.json` | Manifest | 1.5 KB | Pinned repository references for rapid IDE context switching. |

*(Note: `repo-secrets/09-antigravity-backup/temp/.gitkeep` preserves empty temporary directory structure).*

### 4.3 Security Token Purge Verification Protocol

To verify that the repository is completely clean of leaked authentication tokens, an automated audit is executed across three operational scopes:

```mermaid
flowchart TD
    A[Start Security Token Audit] --> B[Scope 1: Tracked Codebase & Documentation]
    A --> C[Scope 2: Runtime State & Local Caches]
    A --> D[Scope 3: Environment & Configuration Files]

    B --> E{Check Token Patterns}
    C --> E
    D --> E

    E -->|Regex Match Detected| F[FLAG AS CRITICAL: Immediate Quarantine & Rotation]
    E -->|Zero Matches Found| G[Token Audit Passed: Repository Clean]
```

#### Token Signatures Audited:
1. **GitHub Personal Access Tokens:** `ghp_[0-9a-zA-Z]{36}`, `github_pat_[0-9a-zA-Z_]{82}`
2. **GitHub OAuth / App Tokens:** `gho_[0-9a-zA-Z]{36}`, `ghu_[0-9a-zA-Z]{36}`, `ghs_[0-9a-zA-Z]{36}`
3. **GitLab Personal Tokens:** `glpat-[0-9a-zA-Z\-]{20}`
4. **Slack Tokens:** `xox[baprs]-[0-9a-zA-Z]{10,48}`
5. **Private Key Headers:** `-----BEGIN OPENSSH PRIVATE KEY-----`, `-----BEGIN RSA PRIVATE KEY-----`
6. **Generic Bearer Tokens:** `Bearer\s+[A-Za-z0-9\-\._~\+\/]+=*`

Any token found in git-tracked files, `.ai-memory/`, or `cli/` triggers an immediate security alert.

---

## 5. Task-10: Review Directory Removal & Objective AI Counter-Feedback

### 5.1 Safe Deletion of `docs/review/`

The temporary review directory `docs/review/` was created by News Spark to house the external evaluation:
- `docs/review/honest-feedback.md`
- `docs/review/action-checklist.md`
- `docs/review/readme.md`

#### Deletion Pre-Conditions:
1. All actionable items (P0-01 to P0-06, P1-01, P1-02, P1-06) have been codified into approved specifications (`01-architecture-spec.md` and `02-component-and-cli-spec.md`).
2. Execution subtasks have been enqueued in `.ai-memory/plans/subtasks/234-codebase-review-remediation-and-consolidation/`.
3. The root `readme.md` consolidation, binary untracking, clutter deletion, and `version.json` corrections are fully planned or executed.
4. Objective AI counter-feedback has been permanently documented in the master ledger.

Once pre-conditions are verified, `docs/review/` is deleted entirely from the repository tree.

### 5.2 Objective AI Counter-Review (Addressing News Spark)

The external review by News Spark contained a mixture of high-value structural observations and flawed architectural recommendations. Below is the objective engineering counter-review:

```mermaid
flowchart TD
    subgraph News Spark Feedback Evaluation
        A[News Spark Review] --> B[Accurate Critiques Accepted]
        A --> C[Flawed Recommendations Refuted]
    end

    subgraph Accepted & Remediated
        B --> B1[Root Clutter: 22 scratch scripts & CSVs purged]
        B --> B2[Documentation Duplication: 2,522-line readme pair consolidated]
        B --> B3[Identity Drift: version.json updated to GitMap]
        B --> B4[Binary Tracking: cli/*.syso untracked, .gitignore updated]
        B --> B5[Ad-Hoc Exits: Standardized on AppError and typed cliexit helpers]
    end

    subgraph Refuted & Defended
        C --> C1[Single-File Clone Collapse: Refuted. Separate domains require modular packages]
        C --> C2[Function Sizing Relaxation: Refuted. 8-15 lines prevents AI context bloat]
        C --> C3[Positive Boolean Ban: Refuted. Positive logic prevents inversion bugs]
        C --> C4[Version Churn Accusation: Refuted. Preserving v6.497.0, micro-bumps gated by CI]
    end
```

#### Grounded Rebuttal Points:

1. **Defense of Modular Clone Packages (`cloner`, `clonefrom`, `clonenow`, `clonepick`, `clonenext`):**
   * *Critique:* News Spark viewed five packages as "unnecessary duplication" and recommended merging them into a single engine.
   * *Rebuttal:* This reveals an incomplete understanding of their operational contracts. `cloner` is deeply integrated with the discovery scanner and SQLite database; `clonefrom` handles arbitrary user-authored CSV/JSON plans; `clonenow` is a round-trip cloner preserving the exact hierarchy of `gitmap scan` output; `clonepick` is a partial sparse-checkout engine; and `clonenext` handles semantic version bumping and remote repo creation. Collapsing them into one file would create a 3,000-line monolithic mess with conflicting flag semantics and bloated imports. Unification belongs at the documentation and CLI UX layer (`docs/commands/cloning-architecture.md`), not code collapse.

2. **Defense of 8–15 Line Function Sizing and Modularity:**
   * *Critique:* News Spark claimed the 8–15 line limit "fragments simple logic and fights Go idioms."
   * *Rebuttal:* In an AI-native engineering environment, short, single-purpose functions are essential. They allow AI models to generate, refactor, and reason about code with zero context window overflow, prevent multi-branch cyclomatic complexity, and make unit testing trivial. Fragmented logic is prevented by clean caller composition.

3. **Defense of Positive Booleans and Ban on Double Negatives:**
   * *Critique:* News Spark criticized "banning negation" as artificial.
   * *Rebuttal:* Negative boolean naming (`isDisabled`, `isNotValid`, `unregistered`) combined with logical operators (`!isDisabled`) is one of the most frequent sources of logic errors in polyglot codebases. Enforcing positive naming (`IsEnabled`, `IsValid`, `IsRegistered`) ensures readability and machine checkability.

4. **Defense of Controlled Versioning:**
   * *Critique:* News Spark argued that version numbers reached `v6.496.0` through "indiscriminate bumps."
   * *Rebuttal:* While past CI/CD workflows used micro-bumps for atomic delivery tracking, the current governance policy explicitly mandates that version increments occur only for significant feature milestones. The target version `v6.497.0` is strictly preserved throughout this remediation without unnecessary churn.

---

## 6. Quality Gates & Verification Traceability

| Item ID | Target Deliverable | Verification Command / Check | Gate Criteria |
| :--- | :--- | :--- | :--- |
| **GATE-01** | `docs/commands/cloning-architecture.md` | `test -f docs/commands/cloning-architecture.md` | Architecture doc maps all 5 packages, models, and entrypoints. |
| **GATE-02** | `cli/cliexit/handle.go` Helpers | Unit test `handle_test.go` | `HandleSuccess`, `HandleUsageError`, `HandleValidationError` compiled and verified. |
| **GATE-03** | Ad-Hoc Exit Elimination | Automated code scan | Zero occurrences of `cliexit.HandleError(nil, ...)` in `cli/cmd/`. |
| **GATE-04** | `repo-secrets/` Integrity | `git status repo-secrets/` | Exactly 21 artifacts preserved intact; zero deletions. |
| **GATE-05** | Security Token Verification | Automated regex scan script | Zero token or secret patterns found across repository. |
| **GATE-06** | `docs/review/` Deletion | `test ! -d docs/review/` | Directory deleted and absent from git tracking. |
| **GATE-07** | Relative Git Paths | Automated path scanner | Zero absolute drive letters or absolute file paths across markdown docs. |

---

## 7. Subtask Index & Execution Allocation

This specification feeds directly into three downstream execution subtasks:

1. **Subtask 03:** [03-fix-version-json-identity-and-clone-architecture-spec.md](../../../.ai-memory/plans/subtasks/234-codebase-review-remediation-and-consolidation/03-fix-version-json-identity-and-clone-architecture-spec.md)  
   *Correct `version.json` identity to GitMap and author `docs/commands/cloning-architecture.md`.*
2. **Subtask 04:** [04-error-exit-audit-and-security-token-verification.md](../../../.ai-memory/plans/subtasks/234-codebase-review-remediation-and-consolidation/04-error-exit-audit-and-security-token-verification.md)  
   *Standardize error exit paths in `cli/cmd/` and execute automated security token check.*
3. **Subtask 05:** [05-repo-secrets-audit-and-honest-counter-review.md](../../../.ai-memory/plans/subtasks/234-codebase-review-remediation-and-consolidation/05-repo-secrets-audit-and-honest-counter-review.md)  
   *Perform non-destructive audit of `repo-secrets/`, clean removal of `docs/review/`, and document honest counter-review summary.*
