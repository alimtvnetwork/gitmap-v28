# Subtask 06: Upstream Coding Guidelines, V6 Prompt & GitMap Skills Sync

- **Subtask ID:** `subtask-06`
- **Slug:** `coding-guidelines-update-and-sync`
- **Parent Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../pending/294-code-refactoring-and-coding-guideline-update.md)
- **Component Spec:** [02-component-spec.md](../../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md)
- **Migration & Verification Spec:** [03-migration-and-verification-spec.md](../../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/03-migration-and-verification-spec.md)
- **Status:** Completed
- **Wave:** Wave 3
- **Assigned Worker:** Autonomous Subagent (Worker 06)

---

## 1. Context & Objectives

During past refactoring cycles, four critical coding standard regressions occurred repeatedly across repositories:
1. Ad-hoc date-time formatting calls (`time.Now().UTC().Format(time.RFC3339)`) scattered throughout business logic instead of calling centralized utilities (`timeutil.NowUTCISO()`).
2. Function calls embedded inline within conditional `if` statements (`if err := fn(a, b); err != nil`), making guard clauses crowded and violating statement simplicity.
3. Multi-argument function declarations and invocations with $>2$ arguments crammed into single unwieldy lines instead of vertical multiline formatting.
4. Direct raw SQL execution in UI renderers and leaf command files, coupled with swallowed errors (`_ = db.Close()`) and missing `*apperror.AppError` return sovereignty.

Furthermore, machine-to-machine AI agent communication was hampered by plain-text table rendering for subtasks rather than structured, collapsed JSON payloads.

This subtask updates the canonical `coding-guidelines` repository and the GitMap skills documentation to enshrine these rules as permanent, non-negotiable standards.

---

## 2. Disjoint File Boundary

This subtask exclusively owns and modifies the following file set:
- `coding-guidelines/01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`
- `coding-guidelines/02-spec/02-coding-guidelines/01-cross-language/37-datetime-and-timeutil-standards.md`
- `coding-guidelines/02-spec/02-coding-guidelines/01-cross-language/38-no-inline-function-calls-in-if.md`
- `coding-guidelines/02-spec/02-coding-guidelines/01-cross-language/39-multiline-parameter-and-argument-formatting.md`
- `coding-guidelines/02-spec/02-coding-guidelines/01-cross-language/40-database-repository-pattern-and-error-sovereignty.md`
- `coding-guidelines/02-spec/02-coding-guidelines/readme.md`
- `coding-guidelines/.agents/skills/gitmap/skill.md`
- `.agents/skills/gitmap/skill.md` (local mirror in `gitmap`)

Zero overlap with Subtasks 01 through 05 or Subtask 07.

---

## 3. Step-by-Step Implementation Instructions

### Step 3.1: Update `02-execute-parent-task-with-n-steps-v6.md`

In `coding-guidelines/01-prompts/14-execute/02-execute-parent-task-with-n-steps-v6.md`, add the following core discipline sections under both the prompt invariants and the code quality checklist:

1. **Centralized `timeutil` Mandate:**
   ```markdown
   - **Centralized Date-Time Utility (`timeutil`):**
     - NEVER construct ad-hoc timestamp strings via `time.Now().UTC().Format(...)` or string concatenations.
     - ALWAYS call centralized time providers: `timeutil.NowUTCISO()` in Go, `DateUtil::nowUtcIso()` in PHP, `now_utc_iso()` in Python, `nowUtcIso()` in TypeScript.
     - Single point of change: all timestamp formats must update everywhere through utility helpers.
   ```

2. **Ban on Inline Function Invocations Inside `if` Statements:**
   ```markdown
   - **Zero Inline Function Calls in `if` Statements:**
     - NEVER call functions inside conditional `if` expressions:
       ```go
       // ❌ FORBIDDEN
       if err := setupSSHDirAndKey(keyPath, email, assumeYes, force); err != nil {
           return err
       }
       ```
     - ALWAYS invoke the function on a separate line above the `if`, store the result in a variable, and evaluate the simple variable in the guard clause:
       ```go
       // ✅ MANDATORY
       err := setupSSHDirAndKey(
           keyPath,
           email,
           assumeYes,
           force,
       )
       if err != nil {
           return err
       }
       ```
   ```

3. **Multiline Formatting for $>2$ Arguments:**
   ```markdown
   - **Multiline Parameter & Argument Formatting ($>2$ arguments):**
     - Whenever a function declaration or invocation contains more than two arguments, format each argument on its own dedicated line with trailing commas or closing parenthesis alignment:
       ```go
       // ✅ MANDATORY
       err := executePipeline(
           context.Background(),
           repoPath,
           branchName,
           true,
       )
       ```
   ```

4. **Database Repository Pattern & `*apperror.AppError` Return Sovereignty:**
   ```markdown
   - **Database Repository Architecture & Error Sovereignty:**
     - Raw SQL strings (`SELECT ... FROM ...`) must NEVER reside in UI rendering files, Cobra command files, or high-level coordinators.
     - Encapsulate all database operations within dedicated repository types (`RepoRepository`, `HeatmapRepository`) with named SQL constants.
     - Database close operations must NEVER swallow errors: `if err := db.Close(); err != nil { log.Warn(err) }` (eliminate `_ = db.Close()`).
     - Leaf repository methods must return `*apperror.AppError` directly, bubbling up to the top-level command dispatcher for database audit capture.
   ```

---

### Step 3.2: Author New Guideline Specifications in `02-spec/02-coding-guidelines/01-cross-language/`

Create four new canonical specification files:

1. **`37-datetime-and-timeutil-standards.md`:**
   - Define canonical time helper interfaces across Go, TypeScript, Python, PHP.
   - Mandate ISO8601 UTC layout (`YYYY-MM-DDTHH:MM:SSZ`).
   - Document mockable time interfaces for deterministic unit testing.

2. **`38-no-inline-function-calls-in-if.md`:**
   - Document the cognitive load and debuggability benefits of separating execution from condition checking.
   - Contrast `if err := ...; err != nil` vs. `err := ...\nif err != nil`.
   - Provide concrete examples in Go, Rust, and Python.

3. **`39-multiline-parameter-and-argument-formatting.md`:**
   - Define readability thresholds: $\le 2$ arguments may remain inline; $\ge 3$ arguments MUST be multiline.
   - Enforce trailing comma and paren alignment rules.

4. **`40-database-repository-pattern-and-error-sovereignty.md`:**
   - Detail the Split-DB repository abstraction pattern.
   - Enforce database metadata enums holding file paths and table names.
   - Ban on `_ = db.Close()` and raw `os.Exit` in database layers.

5. Update `coding-guidelines/02-spec/02-coding-guidelines/readme.md` with links to the 4 new specifications.

---

### Step 3.3: Update GitMap Skills Documentation

In both `coding-guidelines/.agents/skills/gitmap/skill.md` and `.agents/skills/gitmap/skill.md`:

1. **Add `gitmap agent subtask ls --json` Documentation:**
   - Document the `--json` / `-j` flag.
   - Provide the exact JSON output schema:
     ```json
     [
       {
         "id": "subtask-01",
         "name": "IDE Tokens & Workdir Constants",
         "status": "pending",
         "assignedTo": "Worker-01",
         "database": ".ai-memory/temp-agents/agent-task.db"
       }
     ]
     ```

2. **Add `gitmap agent subtask add` Grouped JSON Documentation:**
   - Document the `--json` output grouping the shared database path once:
     ```json
     {
       "database": ".ai-memory/temp-agents/agent-task.db",
       "addedCount": 2,
       "subtasks": [
         { "id": "task-01", "name": "Refactor tokens" },
         { "id": "task-02", "name": "Extract repository" }
       ]
     }
     ```

3. **Add AI Agent Communication Guidance:**
   - Document that AI agents communicating with GitMap or running subagent swarms must parse JSON outputs rather than ANSI tables to maximize parsing precision and eliminate regex scraping bugs.

---

## 4. Verification & Quality Gate Standards

1. Run `python 03-ai-scripts/50-fastgate.py` to confirm all newly authored and edited files:
   - Use strictly relative Git paths (zero absolute paths or URI schemes).
   - Use strictly lowercase filenames.
   - Meet markdown formatting and link integrity standards.
2. Verify that `gitmap` and `coding-guidelines` skills files remain identical across mirrors.
3. Zero compiler builds and zero unit tests executed.
