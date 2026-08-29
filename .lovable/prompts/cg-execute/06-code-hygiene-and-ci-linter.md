# Instruction (must follow): Coding Guideline Execution — Code Hygiene, File Caps & Master CI Linters

Trigger Keywords & Aliases: `cg-hygiene`, `cg-execute hygiene`, `fix code hygiene`, `execute linter guidelines`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to scanning the codebase for file size violations, committed artifacts, un-mirrored guidelines, and missing CI linters, writing the master execution plan to `.lovable/plans/pending/XX-code-hygiene-audit.md`, and decomposing it into subtasks in `.lovable/plans/subtasks/XX-code-hygiene/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until all file sizes, repository hygiene rules, and automated linters are 100% compliant and connected to CI/CD.
- [ ] /goal **Linter Mandate**: You MUST ensure that every coding guideline category has an automated linter check integrated into the CI/CD runner (`.lovable/ai-fix-scripts/03-cicd-local-runner.py` and GitHub Actions). If any linter is missing, write it and wire it up.
- [ ] /learn Ingest `.lovable/coding-guidelines/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and `.lovable/memory/00-index.md` before touching any code.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. Code Hygiene & Linter Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across the entire codebase:

### A. File Size & Structure Caps
1. **File Size Cap**: Any source file must NOT exceed 300 lines maximum.
2. **Struct / Class Cap**: Any struct or class definition file must NOT exceed 120 lines maximum.
3. **Dedicated Definition Files**: Types, enums, constants, and interfaces must live in dedicated files (e.g. `types.go`, `constants.go`), never defined inline next to business logic.
4. **DRY Rule**: Duplicate logic across two or more sites MUST be extracted into a shared helper immediately.

### B. Repository Hygiene & Artifact Protection
1. **No Committed Generated Code or Binaries**: Never commit compiled binaries (`.exe`, `.dll`, `.so`), cache directories (`__pycache__`, `*.pyc`), test report artifacts, or generated code. Ensure `.gitignore` ignores all build output.
2. **Single Source of Truth Versioning**: Root `version.json` is the sole source of truth for version information. Never hardcode version strings across independent files.
3. **Guideline Synchronization**: Maintain byte-for-byte synchronization across canonical guideline files (`spec/17-consolidated-guidelines/31-compiled-simple-coding-guidelines.md`, `.lovable/coding-guidelines/coding-guidelines.md`, and `.cursorrules`).

### C. Master CI/CD Linter Integration
1. **Automated Linter Coverage**: The CI local runner and GitHub Actions pipeline must execute automated checks for:
   - Style and line-gap compliance.
   - Centralized error handling (zero bare `panic` or `os.Exit`).
   - Boolean prefix and semantic naming compliance.
   - Database schema naming and migration integrity.
   - File size caps (<= 300 lines per file).
2. **Deterministic Exit Codes**: All linters must fail with non-zero exit codes upon detecting violations to protect main branch integrity.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Hygiene Audit**: Identify files exceeding 300 lines, inline type definitions, untracked artifacts, and gaps in automated CI linters.
2. **Master Plan**: Write a detailed execution plan to `.lovable/plans/pending/XX-code-hygiene-audit.md`.
3. **Subtask Files**: Decompose into subtask files in `.lovable/plans/subtasks/XX-code-hygiene/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Construction**: Build or update missing linter scripts in `scripts/` or `.lovable/ai-fix-scripts/` and wire them into the local runner.

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Split Oversized Files**: Modularize files > 300 lines into cohesive, domain-specific packages.
2. **Sanitize Git Hygiene**: Proactively update `.gitignore` and remove any generated artifacts.
3. **Execute Full CI Suite**: Run the local CI runner (`python .lovable/ai-fix-scripts/03-cicd-local-runner.py`) and verify all quality gates exit 0.
4. **Update Status**: Mark completed tasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] All source files are <= 300 lines (classes/structs <= 120 lines).
- [ ] Definitions live in dedicated files.
- [ ] No generated binaries or caches are tracked in git.
- [ ] All guideline mirrors are synchronized byte-for-byte.
- [ ] All automated CI/CD linters pass with exit 0.
