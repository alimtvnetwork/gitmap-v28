# Subtask 07: CG Wrapped Error Plan Enqueueing, FastGate Verification & Atomic Commit

- **Subtask ID:** `subtask-07`
- **Slug:** `cg-wrapped-error-plan-enqueue-and-verification`
- **Parent Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../pending/294-code-refactoring-and-coding-guideline-update.md)
- **Component Spec:** [02-component-spec.md](../../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md)
- **Migration & Verification Spec:** [03-migration-and-verification-spec.md](../../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/03-migration-and-verification-spec.md)
- **Status:** Completed
- **Wave:** Wave 3
- **Assigned Worker:** Autonomous Subagent / Lead Orchestrator (Worker 07)

---

## 1. Context & Objectives

The user prompt mandates two essential concluding actions following the completion of the refactoring and guidelines update:
1. **Enqueue CG Execute Wrapped Error Plan:**
   - User directive: *"After running and executing this, you should also enqueue the CG execute wrapped error prompt so that it executes afterwards. There are several places that it is not respected."*
   - A dedicated pending master plan `.ai-memory/plans/pending/295-cg-execute-wrapped-error-remediation.md` must be registered in `.ai-memory/plans/readme.md`, establishing an actionable multi-batch backlog to eliminate unwrapped errors across all remaining Go packages.
2. **Quality Gate Verification & Atomic Commit:**
   - Run `python 03-ai-scripts/50-fastgate.py` to confirm 100% compliance with repository coding guidelines, relative path invariants, function line limits ($\le 15$ lines), positive booleans, no inline `if` calls, and multiline parameter formatting.
   - Execute the final hyphen-separated atomic commit via `gitmap cpf`.

---

## 2. Disjoint File Boundary

This subtask exclusively owns and modifies the following file set:
- `.ai-memory/plans/pending/295-cg-execute-wrapped-error-remediation.md`
- `.ai-memory/plans/readme.md` (Pending section update)
- Execution of verification runner: `03-ai-scripts/50-fastgate.py`

Zero overlap with Subtasks 01 through 06.

---

## 3. Step-by-Step Implementation Instructions

### Step 3.1: Author Plan 295 (`295-cg-execute-wrapped-error-remediation.md`)

Create `.ai-memory/plans/pending/295-cg-execute-wrapped-error-remediation.md` with the following structure:
- **Title:** `# Master Plan 295: Repository-Wide Wrapped Error & AppError Remediation (CG Execute)`
- **Status:** `PENDING`
- **Objective:** Systematically scan and remediate all remaining unwrapped Go errors, raw string `fmt.Errorf` calls without apperror wrapping, and functions returning generic `error` instead of strongly typed `*apperror.AppError`.
- **Target Subsystems:**
  - `cli/cmd/` (CLI root dispatchers and tooling commands).
  - `cli/store/` (Storage engine and query filters).
  - `cli/cmdpe/` (Pipeline execution and diagnostics).
  - `cli/cmdpipeline/` (Pipeline database operations and history rendering).
  - `cli/cmdpull/` (Pull coordination and target resolvers).
  - `cli/cmdexclude/` (Exclusion rules and scan filters).
- **Execution Workflow:** Follow `01-prompts/12-coding-guidelines/` micro-batching protocol (5–8 files per wave) with strict leaf error sovereignty and database audit logging.

### Step 3.2: Register Plan 295 in `.ai-memory/plans/readme.md`

Under Section 3 (`Pending Plans`) of `.ai-memory/plans/readme.md`, append the entry:
```markdown
- [295-cg-execute-wrapped-error-remediation.md](pending/295-cg-execute-wrapped-error-remediation.md) — Repository-Wide Wrapped Error & AppError Remediation (CG Execute)
```

---

### Step 3.3: Execute FastGate Static AST Verification

Invoke FastGate to validate all staged changes across the working tree:
```bash
python 03-ai-scripts/50-fastgate.py
```

**Verification Gate Criteria:**
- **Zero Absolute Paths:** Confirms no file in the diff contains drive letters (`C:\`, `D:\`) or URI schemes (`file:///`).
- **Function Line Limits:** Confirms every modified function has $\le 15$ lines of executable code.
- **File Sizing:** Confirms all files are $\le 300$ lines.
- **Positive Booleans:** Confirms booleans only use `is*` or `has*` prefixes, with zero negative booleans.
- **Zero Inline Calls in `if`:** Confirms no `if x := fn(); x != nil` patterns exist.
- **Multiline Arguments:** Confirms functions with $>2$ arguments use multiline parameter formatting.
- **Zero Builds & Tests:** Verifies that no `go build` or `go test` invocations took place.

---

### Step 3.4: Lead Orchestrator Atomic Commit

Upon 100% clean FastGate confirmation, the lead orchestrator executes an atomic commit adhering to the V6 hyphen-separated standard:
```bash
gitmap cpf "feat(refactor) - code refactoring and coding guideline update for task 294"
```

---

## 4. Acceptance Criteria & Final Sign-Off Checklist

- [x] `.ai-memory/plans/pending/295-cg-execute-wrapped-error-remediation.md` exists and details repository-wide wrapped error remediation.
- [x] `.ai-memory/plans/readme.md` is updated with Plan 295 in the pending registry.
- [x] All file references across all subtasks and specs are strictly relative Git paths.
- [x] All new files and directories are strictly lowercase.
- [ ] Working tree is clean and committed via `gitmap cpf` (to be performed by Lead Orchestrator).
