# Subtask 02: Case Collision Detector, Normalization Recipe & Fix-All Integration

- **Subtask ID**: `02-subtask-storytelling-readme-fix`
- **Parent Task**: `gitmap-enhancement-and-issue-resolution`
- **Component Spec**: `02-spec/21-app/gitmap-enhancement-and-issue-resolution/02-component-spec.md`
- **Owned Files**:
  - `cli/gitutil/case_collision_detector.go`
  - `cli/gitutil/case_collision_recipe.go`
  - `cli/gitutil/case_collision_test.go`
  - `cli/cmdremediation/remediation_box.go`
- **Status**: COMPLETED

---

## 1. Objectives & Context

On Windows NTFS and macOS APFS, case-insensitive file lookups fold `story/README.md` and `story/readme.md` to the same physical inode on disk. When Git tracks both files in its binary index, restoring or discarding changes triggers a perpetual dirty loop where one file's checkout marks the other as modified or deleted.

This subtask implements:
1. `DetectCaseCollisions`: Parse `git ls-files -s` to discover duplicate index entries differing only in case.
2. `GenerateCaseNormalizationRecipe`: Purge the uppercase duplicate from the Git index (`git rm --cached -f <uppercase>`) while strictly preserving disk content, stage the canonical lowercase entry, and create a normalization commit.
3. Integration with `gitmap fix-all` (`fa`) in `cli/cmdremediation/remediation_box.go`.
4. Persistence of remediation telemetry to `~/.gitmap/output/last_remediation.json` and `.gitmap/repo-cache/remediation_history.json`.

---

## 2. Step-by-Step Implementation Plan

### Step 2.1: Data Models & Index Parsing (`cli/gitutil/case_collision_detector.go`)
- Define `IndexEntry` struct (`Mode`, `BlobSHA`, `Stage`, `Path`, `LowerPath`).
- Define `CaseCollision` struct (`CanonicalPath`, `CollidingPaths`, `ActiveEntries`, `HasDivergentSHA`).
- Implement `parseGitLsFilesOutput(output string) []IndexEntry` (<= 15 lines):
  - Split lines, parse octal mode, 40-char SHA, stage, and tab-separated path.
- Implement `DetectCaseCollisions(repoPath string) ([]CaseCollision, error)` (<= 15 lines):
  - Run `git ls-files -s` in `repoPath`.
  - Group entries by lowercased path in a hash map.
  - Extract buckets with `len(bucket) > 1` and differing path casing.
  - Identify canonical path according to lowercase readme conventions.

### Step 2.2: Recipe Generation & Index Purge (`cli/gitutil/case_collision_recipe.go`)
- Implement `GenerateCaseNormalizationRecipe(repoPath string, c CaseCollision) ([]string, error)` (<= 15 lines):
  - Construct shell commands:
    1. `git rm --cached -f "<uppercase_path>"`
    2. `git add -- "<canonical_path>"`
    3. `git commit -m "fix(case): normalize <uppercase> to canonical <lowercase>"`
- Implement `ExecuteCaseNormalization(repoPath string, c CaseCollision) error` (<= 15 lines):
  - Backup working tree file to memory buffer.
  - Execute `git rm --cached -f` for each colliding uppercase entry.
  - Verify disk file exists and write back buffer if needed.
  - Execute `git add --` and commit change cleanly.
  - Return direct error or `nil`.

### Step 2.3: Fix-All & Remediation UI Integration (`cli/cmdremediation/remediation_box.go`)
- Update `RemediationItem` to capture case collision recipes:
  - Add `CaseCollisions []gitutil.CaseCollision` to `RemediationItem`.
- Update `PrintRemediationBox` and `SaveRemediationState`:
  - When case collisions are detected in a repo, surface option `[4] Normalize Case Collisions` or automatically resolve during `fix-all`.
  - Write detailed recipe actions into `last_remediation.json`.
  - Append audit record to `.gitmap/repo-cache/remediation_history.json`.

### Step 2.4: Unit & Integration Tests (`cli/gitutil/case_collision_test.go`)
- Test parsing raw `git ls-files -s` output with synthetic dual entries (`story/README.md` and `story/readme.md`).
- Test detection returns exact canonical path (`story/readme.md`) and uppercase colliding path (`story/README.md`).
- Test recipe generation returns exact `git rm --cached -f` command strings.
- Test that single-entry paths return no collision.

---

## 3. Coding Guidelines Invariants Checklist

- [ ] Error Sovereignty: All functions returning `error` return `err` or `apperror.Wrap(...)` directly; no exit handlers called.
- [ ] Boolean Naming: All booleans prefixed with `is*` or `has*` only (`hasCollision`, `isCanonicalPath`, `hasDivergentSHA`).
- [ ] Function Size: All functions <= 15 lines (<= 8 lines preferred).
- [ ] Conditional Flattening: Zero nested `if` statements; guard clauses with early returns used throughout.
- [ ] Relative Path Hygiene: Zero absolute paths; strictly relative paths across all code and tests.

---

## 4. Verification & Validation Commands

1. **Unit Tests Execution**:
   ```bash
   go test -v ./cli/gitutil/ -run TestDetectCaseCollisions
   go test -v ./cli/gitutil/ -run TestCaseNormalizationRecipe
   ```
2. **Storytelling Repository Simulation**:
   ```bash
   gitmap fa --repo ../story-telling --dry-run
   cat ~/.gitmap/output/last_remediation.json
   ```
3. **Repository Status Check**:
   ```bash
   cd ../story-telling && git status
   ```
   Verify working tree is clean and `story/readme.md` is the sole tracked entry in `git ls-files -s`.
