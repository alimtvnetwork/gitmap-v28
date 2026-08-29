# Instruction (must follow): Coding Guideline Execution — Style, Formatting & Line-Gaps

Trigger Keywords & Aliases: `cg-style`, `cg-execute style`, `fix style guidelines`, `execute style guidelines`

```text
N = 100
```

N = total self-loop steps budget. The user may override this number when triggering the prompt.

- [ ] /goal First `N/2` steps (Phase 1) are dedicated to reading the codebase, identifying all style, formatting, function length, and line-gap violations, writing the master execution plan to `.lovable/plans/pending/XX-style-guidelines-audit.md`, and decomposing it into microscopic subtasks in `.lovable/plans/subtasks/XX-style-guidelines/`.
- [ ] /goal Second `N/2` steps (Phase 2) are dedicated to executing those subtasks in an autonomous self-loop until every style violation is resolved, unit tests pass, and linters exit 0.
- [ ] /goal **Linter Mandate**: If an automated linter for style and line-gap rules does not exist in the codebase, you MUST create an advanced linter script (e.g. in `scripts/` or `linters/`) and connect it directly to the CI/CD local runner and workflows.
- [ ] /learn Ingest `.lovable/coding-guidelines/coding-guidelines.md`, `.lovable/strictly-avoid.md`, and `.lovable/memory/00-index.md` before touching any code.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after the user sets them. Never change them mid-execution.

---

## 1. Style & Formatting Non-Negotiable Checklist

You MUST audit and strictly enforce every rule below across the entire codebase:

### A. Function Length & Complexity
1. **Function Length**: 8 lines preferred, 15 lines hard cap (excluding blank lines and comments). Waiver allowed only via inline comment `// lint-allow: function-length reason="..." max=N`.
2. **Nested `if` Flattening**: No nested `if` statements. Flatten with early returns and guard clauses.
3. **Simple Positive Conditions**: `if` conditions must be positive and simple. No `!`, no double negatives. Extract positively named booleans if negation is needed.
4. **No Mixed Polarity in Joins**: Never mix positive and negative conditions in a single conditional join (`if (a && !b)` is forbidden; use all positive named variables). Max one logical join (two operands) per `if`.

### B. Line-Gap and Whitespace Rules
1. **Single Blank Line Before `return` / `throw`**: Exactly one blank line before every `return` or `throw`, unless it is the only statement in the block.
2. **Single Blank Line After `}`**: Exactly one blank line after a closing `}`, unless the next line is another `}`, `else`, `case`, or `catch`.
3. **Never Consecutive Blank Lines**: Never allow two or more blank lines in a row anywhere.
4. **No Blank Lines at Block Boundaries**: No blank lines immediately after `{` or immediately before `}`.
5. **Top-Level Separation**: Exactly one blank line between top-level declarations (functions, structs, classes, exported constants).
6. **Import Grouping**: Group imports with exactly one blank line between groups: standard library, third-party, first-party absolute, first-party relative.
7. **Trailing Whitespace & EOF**: Zero trailing whitespace on any line. Exactly one trailing newline at the end of the file.

### C. Markdown Document Style
1. **Header Spacing (MD022)**: Exactly one blank line before and after every Markdown heading (`#`, `##`, `###`).
2. **List Spacing (MD032)**: Lists must be surrounded by blank lines.

---

## 2. Phase 1: Planning, Audit & Subtask Decomposition (Steps 1 .. N/2)

1. **Automated Audit**: Scan the codebase using scripts or ripgrep for functions exceeding 15 lines, nested `if` statements, and line-gap violations.
2. **Master Plan**: Write a comprehensive execution plan to `.lovable/plans/pending/XX-style-guidelines-audit.md` detailing every file needing refactoring.
3. **Subtask Files**: Decompose the plan into microscopic subtask files in `.lovable/plans/subtasks/XX-style-guidelines/` (e.g. `01-task.md`, `02-task.md`, ...).
4. **Linter Connection**: Check if a style/formatting linter is registered in the CI/CD local runner (`.lovable/ai-fix-scripts/03-cicd-local-runner.py` or `.github/workflows/`). If missing or incomplete, write a dedicated linter script and wire it into the pipeline.

---

## 3. Phase 2: Autonomous Execution Loop (Steps N/2+1 .. N)

1. **Iterate Through Subtasks**: Process each subtask sequentially, refactoring functions to conform to the 15-line limit and line-gap style rules.
2. **Surgical Refactoring**: Extract helper functions with clear domain names. Do not use generic names (`temp`, `data`, `obj`, `helper1`).
3. **Verification**: Run `go vet ./...`, `go test ./...`, markdown linters, and the CI local runner.
4. **Status Tracking**: Mark completed subtasks as `DONE`, move completed plans to `.lovable/plans/completed/`, and update `.lovable/plans/index.md`.

---

## 4. Pre-Commit Verification Checklist

- [ ] All functions are <= 15 lines (or contain valid `lint-allow` comments).
- [ ] No nested `if` statements remain.
- [ ] Line-gap rules (single blank lines before returns, after closing braces) are 100% compliant.
- [ ] Markdown files comply with MD022 and MD032.
- [ ] Style linter is integrated into CI/CD and exits 0.
- [ ] Working tree is clean and builds pass without errors.
