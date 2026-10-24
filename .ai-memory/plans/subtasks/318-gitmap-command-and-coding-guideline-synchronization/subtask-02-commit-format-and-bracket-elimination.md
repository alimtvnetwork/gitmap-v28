# Subtask Plan 02: Commit Format Standardization and Third Bracket Elimination

- **Subtask Slug:** `subtask-02-commit-format-and-bracket-elimination`
- **Parent Plan:** `318-gitmap-command-and-coding-guideline-synchronization`
- **Specification Reference:** `02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/02-commit-format-and-bracket-elimination-spec.md`
- **Status:** Completed
- **Target Repository:** `gitmap-v28`

---

## 1. Objectives & Scope

This subtask executes the technical refactoring and documentation updates required to eliminate square brackets (`[` `]`) from all commit messages, formatters, tests, help documentation, and AI coding guidelines across GitMap:

1. **Refactor Commit Formatter Engine (`cli/cmdcommit/commit_formatter.go`)**:
   - Eliminate `bracketRegex` and `tryParseBracketSyntax`.
   - Eliminate `buildBracketRecord` and forced bracket output generation.
   - Implement uppercase module normalization and hyphen separation (`Feature: MODULE - summary`).
   - Cleanse user input containing colons or lingering brackets into canonical hyphen format.
   - Enforce strictly positive booleans (`hasModule`, `hasHyphenSeparator`, `isConventional`).
   - Maintain strict 15-line function limits and zero nested conditionals.
2. **Update Formatter Unit Tests & Terminal Help UI**:
   - Update `cli/cmdcommit/commit_formatter_test.go` to test the hyphen standard and assert bracket eradication.
   - Update `cli/cmdcommit/commit_help_menu.go` to display `"MODULE - description"` in usage examples and tips.
3. **Revert Bracketed Rules in Spec 317**:
   - Update `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md` to remove bracket syntax recommendations.
4. **Synchronize AI Coding Guideline 12**:
   - Update `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md` Rule 4 to mandate the uppercase hyphen standard and document the ban on third brackets.
5. **Verify Zero Bracket Invariant & Path Hygiene**:
   - Verify no bracket generation logic remains in `cli/cmdcommit/`.
   - Verify relative path hygiene with `python linter-scripts/check-relative-paths.py`.

---

## 2. Target Files & Component Map

| Target File Path | Purpose | Change Category |
| :--- | :--- | :--- |
| `cli/cmdcommit/commit_formatter.go` | Core commit message parser and formatter | Refactoring / Code |
| `cli/cmdcommit/commit_formatter_test.go` | Unit tests for parser and formatter | Refactoring / Test |
| `cli/cmdcommit/commit_help_menu.go` | Terminal UI help cards and tips | Refactoring / UI |
| `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md` | Prior workgroup spec containing legacy bracket rules | Spec Update |
| `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md` | Master AI commit provenance guideline | Guideline Update |

---

## 3. Atomic Execution Checklist

### Phase 1: Commit Formatter Refactoring (`cli/cmdcommit/commit_formatter.go`)
- [x] **1.1**: Update `ParsedCommitMessage` struct:
  - Replace `HasBracketSyntax bool` with `HasHyphenSeparator bool`.
  - Ensure all booleans use `is` or `has` prefixes exclusively (`hasModule`, `hasHyphenSeparator`, `isConventional`).
- [x] **1.2**: Remove `bracketRegex`.
- [x] **1.3**: Update `conventionalRegex` to standard conventional commit format without brackets:
  ```go
  conventionalRegex = regexp.MustCompile(`^([a-z]+)\(([A-Za-z0-9_\-\./]+)\):\s*(.+)$`)
  ```
- [x] **1.4**: Define `hyphenRegex` to match both space-hyphen and colon variations for robust normalization:
  ```go
  hyphenRegex = regexp.MustCompile(`^([A-Za-z0-9_\-\./]+)\s*[:-]\s*(.+)$`)
  ```
- [x] **1.5**: Implement `cleanseBracketArtifacts(s string) string` helper to strip any accidental leading/trailing square brackets from module tokens.
- [x] **1.6**: Implement `buildHyphenRecord(raw, prefix, mod, item string) ParsedCommitMessage`:
  - Normalize `mod` to uppercase (`strings.ToUpper(strings.Trim(mod, "[] ")))`.
  - Format output as `prefix + upperMod + " - " + item`.
  - Set `HasModule: true`, `HasHyphenSeparator: true`.
- [x] **1.7**: Remove `buildBracketRecord` and `tryParseBracketSyntax`.
- [x] **1.8**: Implement `tryParseHyphenSyntax(verb, s string) (ParsedCommitMessage, bool)` using `buildHyphenRecord`.
- [x] **1.9**: Ensure all functions in `commit_formatter.go` are $\le 15$ lines max with guard clauses, early returns, and proper vertical spacing.

### Phase 2: Formatter Unit Tests & Help Menu Updates
- [x] **2.1**: Refactor `cli/cmdcommit/commit_formatter_test.go`:
  - Remove all test cases asserting bracket output (e.g. `Feature: [auth]: ...`).
  - Add test cases for uppercase hyphen output across semantic verbs:
    - `cpf "auth - add jwt"` $\rightarrow$ `Feature: AUTH - add jwt`
    - `cbf "db - fix deadlock"` $\rightarrow$ `Bug: DB - fix deadlock`
    - `cpc "skills - sync skill"` $\rightarrow$ `Chore: SKILLS - sync skill`
    - `cpr "v6.605.0 - release notes"` $\rightarrow$ `Release: V6.605.0 - release notes`
  - Add test cases for bracket stripping and sanitization:
    - `cpf "[auth]: add jwt"` $\rightarrow$ `Feature: AUTH - add jwt`
    - `cpf "auth: add jwt"` $\rightarrow$ `Feature: AUTH - add jwt`
  - Add test cases verifying duplicate prefix removal without bracket insertion.
- [x] **2.2**: Update `cli/cmdcommit/commit_help_menu.go`:
  - Update `buildCommitTips()`: change `"gitmap cpf \"[module]: description\""` to `"gitmap cpf \"MODULE - description\""`.
  - Update `buildCommitAIEntries()`: change descriptions from `commit [module]: feat` to `commit MODULE - feat`.

### Phase 3: Specification & Guideline Synchronization
- [x] **3.1**: Update `02-spec/21-app/317-system-for-folder-specific-work-directories-and-ssh-key-automation/03-cli-commands-and-commit-format-spec.md`:
  - In Section 1.2, remove bracketed module syntax descriptions.
  - In Section 4.2, update EBNF grammar to the uppercase hyphen format.
  - In Section 4.3 & 4.4, replace bracketed examples with approved hyphen format.
- [x] **3.2**: Update `02-spec/02-coding-guidelines/06-ai-optimization/12-commit-provenance-and-tracking.md`:
  - Rewrite Rule 4: replace `"Hyphen-Separated & Bracketed Module Commit Message Standards"` with `"Hyphen-Separated Uppercase Module Standard (Zero Third Brackets)"`.
  - Add explicit rationale explaining why square brackets are banned.
  - Insert the comprehensive Banned vs Approved examples table.

### Phase 4: Verification Contract & Quality Gates
- [x] **4.1**: Verify that zero functions in `cli/cmdcommit/` produce square brackets in commit headers.
- [x] **4.2**: Verify relative path hygiene:
  ```powershell
  python linter-scripts/check-relative-paths.py
  ```
- [x] **4.3**: Run FastGate linter:
  ```powershell
  python 03-ai-scripts/50-fastgate.py
  ```

---

## 4. Verification Contract & Acceptance Criteria

| Check | Target / Metric | Acceptance Threshold | Positive Boolean |
| :--- | :--- | :--- | :--- |
| **Bracket Elimination** | `cli/cmdcommit/commit_formatter.go` output | 0 square brackets in formatted message | `isBracketFree = true` |
| **Uppercase Normalization** | Module token in output | All alphabetical chars uppercase | `hasUppercaseModule = true` |
| **Hyphen Formatting** | Message separator | Space-padded hyphen (`" - "`) | `hasHyphenSeparator = true` |
| **Function Length** | Lines of code per function | $\le 15$ lines | `isWithinLineLimit = true` |
| **Relative Path Hygiene** | Repository-wide path checker | 0 absolute filesystem paths | `isPathHygieneClean = true` |
| **FastGate Linter** | `python 03-ai-scripts/50-fastgate.py` | 100% PASS (< 1.5s) | `isFastGatePassed = true` |
