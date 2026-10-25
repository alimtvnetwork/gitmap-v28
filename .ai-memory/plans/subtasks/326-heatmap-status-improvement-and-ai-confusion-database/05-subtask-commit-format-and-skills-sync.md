# Subtask Plan 05: Commit Format Remediation and Skills Synchronization

- **Subtask Slug:** `05-subtask-commit-format-and-skills-sync`
- **Parent Plan:** `326-heatmap-status-improvement-and-ai-confusion-database`
- **Specification Reference:** `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/03-migration-and-verification-spec.md`
- **Target Repository:** `gitmap-v28`
- **Status:** `COMPLETED`
- **Created At:** 2026-10-25

---

## 1. Objectives & Scope

This subtask remediates ambiguous and malformed commit instructions across GitMap skills and prompts, and fortifies the commit formatting engine in `cli/cmdcommit/commit_formatter.go`:

1. **Remediate Prompt & Skill Files:**
   - Update `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` to eliminate malformed bracketed and type-module examples (`cpf "(MODULE): ..."`, `cpr "RELEASE - ..."`).
   - Update `.agents/skills/gitmap/skill.md` to reinforce the strict 5-verb taxonomy (`cpr`, `cpc`, `cbf`, `cfr`, `cpf`) and explicitly ban pseudo-module stuttering.
   - Update `.agents/skills/muse-master-prompt/skill.md` to clarify the canonical hyphen standard (`gitmap <verb> "MODULE - <summary>"`) with true architectural subsystems.
2. **Define Subsystem Module Standard:**
   - Explicitly define `MODULE` as an architectural subsystem (`CLI`, `HEATMAP`, `STORE`, `FIX`, `RC`, `PULL`, `SCANNER`, `NODES`, `AUM`, `SYNC`, `DOCS`).
   - Strictly prohibit using commit verbs/types as the module name (`FEATURE`, `RELEASE`, `BUG`, `FIX`, `CHORE`, `SPEC`).
3. **Defensive Formatting in `cli/cmdcommit/commit_formatter.go`:**
   - Implement redundant prefix stripping (`Feature: `, `FEAT: `, `Release: `, `Bug: `, `Fix: `, `Refactor: `, `Chore: `).
   - Implement pseudo-module collision detection: when an input uses a banned commit type as the module (e.g. `gitmap cpf "FEATURE - add flag"`), strip the redundant type token to prevent stuttered outputs like `Feature: FEATURE - add flag`.
   - Guarantee Title Case verb prefixing (`Feature: `, `Bug: `, `Refactor: `, `Chore: `, `Release: `).
4. **Unit Test Updates (`cli/cmdcommit/commit_formatter_test.go`):**
   - Author comprehensive test cases covering prefix deduplication, pseudo-module stripping, and title-case enforcement.
5. **Coding Guideline Conformance:**
   - Maintain function lengths $\le 15$ lines.
   - Zero nested `if` statements.
   - Positive boolean naming exclusively (`hasModule`, `hasHyphenSeparator`, `hasDuplicatePrefix`, `isTypeModule`, `isConventional`).

---

## 2. Target Files & Disjoint Ownership

| Target File Path | Purpose | Modification Category |
| :--- | :--- | :--- |
| `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md` | V6 execution prompt & core agent rules | Skill / Documentation |
| `.agents/skills/gitmap/skill.md` | Master GitMap CLI reference & verb rules | Skill / Documentation |
| `.agents/skills/muse-master-prompt/skill.md` | Muse master onboarding and execution prompt | Skill / Documentation |
| `cli/cmdcommit/commit_formatter.go` | Commit message parser, sanitizer, and formatter | Core Go Engine |
| `cli/cmdcommit/commit_formatter_test.go` | Formatter unit test assertions | Go Test Suite |

---

## 3. Atomic Implementation Checklist

### Phase 1: Skill Files Remediation

- [x] **1.1**: Update `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`:
  - Locate Section 8 (R8/R9 Atomic Commit & Push via GitMap).
  - Replace ambiguous parentheses examples e.g. `gitmap cpf "(MODULE): <feature summary>"` with canonical hyphen format:
    ```bash
    gitmap cpf "MODULE - <feature summary>"
    gitmap cbf "MODULE - <fix summary>"
    gitmap cfr "MODULE - <refactor summary>"
    gitmap cpr "MODULE - <release summary>"
    gitmap cpc "MODULE - <chore summary>"
    ```
  - Add explicit note defining `MODULE` as the codebase subsystem (`CLI`, `HEATMAP`, `STORE`, `FIX`, `RC`, `PULL`, `CLONE`, `SCANNER`, `NODES`, `AUM`, `SYNC`, `DOCS`).
  - Add explicit negative guardrail: **TOTAL BAN on commit types as MODULE** (`FEATURE`, `RELEASE`, `BUG`, `FIX`, `CHORE`, `SPEC`).
  - Reinforce the 5-verb rule: never use `cpf` for chores or releases.

- [x] **1.2**: Update `.agents/skills/gitmap/skill.md`:
  - Locate Section 7 (GitMap Commit & Release Commands).
  - Clarify the prefix generation table: GitMap automatically prepends Title Case prefixes (`Feature: `, `Bug: `, `Refactor: `, `Chore: `, `Release: `).
  - Add approved subsystem module taxonomy table.
  - Update usage examples to show proper subsystem tags (`CLI`, `HEATMAP`, `STORE`, `FIX`, `RC`).

- [x] **1.3**: Update `.agents/skills/muse-master-prompt/skill.md`:
  - Locate Section 2, Step 8 (Verify, Pre-Commit Secrets Gate & Atomic Commit).
  - Clarify commit command syntax: `gitmap cpf "MODULE - <summary>"`.
  - Ensure all 5 verbs are referenced (`cpf`, `cbf`, `cfr`, `cpc`, `cpr`).

---

### Phase 2: Formatter Defensive Updates (`cli/cmdcommit/commit_formatter.go`)

- [x] **2.1**: Audit and update `resolveVerbPrefix`:
  - Ensure title-case prefixes are returned:
    - `cpf`, `feature`, `commit-push-feature` $\rightarrow$ `"Feature: "`
    - `cbf`, `cpb`, `bug`, `fix`, `commit-push-bug` $\rightarrow$ `"Bug: "` (or `"Fix: "` if matching fix alias)
    - `cfr`, `refactor`, `commit-push-refactor` $\rightarrow$ `"Refactor: "`
    - `cpc`, `chore`, `commit-push-chore` $\rightarrow$ `"Chore: "`
    - `cpr`, `release`, `commit-push-release` $\rightarrow$ `"Release: "`

- [x] **2.2**: Implement `stripRedundantCommitPrefixes(s string) string`:
  - Detect and remove leading redundant prefixes like `(?i)^(feature|feat|release|bug|fix|refactor|chore):\s*`.
  - Handle case-insensitive prefix variations without regex overhead where string slicing suffices.

- [x] **2.3**: Implement `isCommitTypeToken(mod string) bool`:
  - Return `true` if `mod` (uppercased and trimmed) matches any reserved commit type token:
    - `"FEATURE"`, `"FEAT"`, `"RELEASE"`, `"REL"`, `"BUG"`, `"FIX"`, `"CHORE"`, `"REFACTOR"`, `"REFACT"`, `"SPEC"`.
  - Use positive boolean naming (`isCommitTypeToken`).

- [x] **2.4**: Update `buildHyphenRecord(raw, prefix, mod, item string) ParsedCommitMessage`:
  - Check if `isCommitTypeToken(mod)` is `true`.
  - If `true`, inspect `item`:
    - If `item` itself contains a nested hyphen with a subsystem (e.g. `CLI - add flag`), extract `CLI` as the true `mod` and `add flag` as `item`.
    - Otherwise, strip the redundant `mod` and construct a clean output without stuttering (`prefix + cleanItem`).
  - If `false`, format standard canonical output: `prefix + upperMod + " - " + cleanItem`.

- [x] **2.5**: Coding guideline compliance:
  - Verify every function in `commit_formatter.go` is $\le 15$ lines.
  - Verify zero nested conditionals.
  - Verify proper blank line spacing around control flow and returns.

---

### Phase 3: Formatter Unit Test Expansion (`cli/cmdcommit/commit_formatter_test.go`)

- [x] **3.1**: Add test cases for redundant prefix elimination:
  - `cpf "Feature: CLI - add flag"` $\rightarrow$ `Feature: CLI - add flag`
  - `cpf "FEAT: CLI - add flag"` $\rightarrow$ `Feature: CLI - add flag`
  - `cpr "Release: V6.605.0 - release notes"` $\rightarrow$ `Release: V6.605.0 - release notes`
  - `cbf "Fix: STORE - handle nil db"` $\rightarrow$ `Bug: STORE - handle nil db`

- [x] **3.2**: Add test cases for pseudo-module stripping:
  - `cpf "FEATURE - add flag"` $\rightarrow$ `Feature: add flag`
  - `cpf "FEATURE - CLI - add flag"` $\rightarrow$ `Feature: CLI - add flag`
  - `cpr "RELEASE - v6.605.0"` $\rightarrow$ `Release: v6.605.0`
  - `cpf "SPEC - author specs"` $\rightarrow$ `Feature: author specs`

- [x] **3.3**: Add test cases verifying standard canonical behavior:
  - `cpf "HEATMAP - visual overhaul"` $\rightarrow$ `Feature: HEATMAP - visual overhaul`
  - `cbf "FIX - prevent non-git commit"` $\rightarrow$ `Bug: FIX - prevent non-git commit`
  - `cfr "STORE - optimize query"` $\rightarrow$ `Refactor: STORE - optimize query`
  - `cpc "DOCS - update readme"` $\rightarrow$ `Chore: DOCS - update readme`

---

### Phase 4: Quality Verification & Boundary Defense

- [x] **4.1**: Relative Path Verification:
  - Run relative path checker to ensure 0 absolute paths exist in modified files:
    ```powershell
    python linter-scripts/check-relative-paths.py
    ```
- [x] **4.2**: Boolean Hygiene Verification:
  - Verify all added booleans use `is` or `has` prefixes with positive polarity.
- [x] **4.3**: Clean File Boundaries:
  - Ensure modifications are strictly confined to the 5 designated target files.
