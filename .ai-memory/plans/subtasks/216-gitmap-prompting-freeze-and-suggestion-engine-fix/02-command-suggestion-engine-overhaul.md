# Subtask 216.02: Command Suggestion Engine Overhaul

## 1. Context & Objective

GitMap provides an intelligent typo suggestion engine for unrecognized commands invoked at the root CLI level (e.g. `gitmap <mistyped-command>`). However, users reported that typing slight variations or misspellings of critical commands like `install`, `uninstall`, `apps`, `commit`, `search`, and `find` yielded no suggestions or unrelated keyword fallbacks.

The root cause is that `collectTopCommandCandidates()` in `cli/cmd/rootsuggest_calc.go` only aggregated commands from `primaryTopCommands` (a manually maintained 101-command slice in `cli/cmd/rootsuggest.go`), `"run"`, `"run-until"`, and user macros. It completely omitted `completion.AllCommands()`, which contains the complete inventory of 512 commands and aliases.

The objective of this subtask is to overhaul `collectTopCommandCandidates()` by merging `completion.AllCommands()`, `primaryTopCommands`, `"run"`, `"run-until"`, and `macro.ListMacros()`, deduplicating candidates with a `seen` map, and verifying Levenshtein matching across all 512 commands and aliases.

---

## 2. Target Files

- `cli/cmd/rootsuggest_calc.go`
- `cli/cmd/rootsuggest.go`
- `cli/cmd/rootsuggest_test.go`

---

## 3. Technical Investigation & Scope

1. **Candidate Pool Gap in `cli/cmd/rootsuggest_calc.go`:**
   - Currently, `collectTopCommandCandidates()` executes:
     ```go
     func collectTopCommandCandidates() []string {
         candidates := append([]string{}, primaryTopCommands...)
         candidates = append(candidates, "run", "run-until")
         macroList := macro.ListMacros()
         ...
     ```
   - It ignores `completion.AllCommands()`, leaving hundreds of commands and aliases out of the candidate pool.
2. **Missing Core Commands in `primaryTopCommands`:**
   - `primaryTopCommands` in `cli/cmd/rootsuggest.go` omitted `install`, `uninstall`, `apps`, `commit`, `search`, `find`, `login`, `setup`, `cpf`, `cpb`, `cpr`.
3. **Levenshtein Distance Evaluation:**
   - `rankCandidateCommands()` calculates distance and similarity. With an incomplete pool, commands like `instlal` (distance 1 from `install`) could not be matched because `install` was not in the candidate slice.

---

## 4. Implementation Steps

### Step 1: Update `collectTopCommandCandidates()` in `cli/cmd/rootsuggest_calc.go`
- Import `github.com/alimtvnetwork/gitmap-v28/cli/completion`.
- Create a `seen := make(map[string]bool)` map for deduplication.
- Append commands from `completion.AllCommands()`.
- Append commands from `primaryTopCommands`.
- Append `"run"` and `"run-until"`.
- If `macroList := macro.ListMacros()` succeeds, append macro names.
- Ensure all items added are trimmed and non-empty.

### Step 2: Expand `primaryTopCommands` in `cli/cmd/rootsuggest.go`
- Add missing core commands to `primaryTopCommands`:
  - `install`, `uninstall`, `apps`, `commit`, `search`, `find`, `login`, `setup`, `cpf`, `cpb`, `cpr`.
- This ensures high priority for core commands in direct ranking evaluations.

### Step 3: Implement Comprehensive Test Coverage in `cli/cmd/rootsuggest_test.go`
- Test that `collectTopCommandCandidates()` returns $\ge 500$ unique items with 0 duplicates.
- Test Levenshtein distance suggestions:
  - `gitmap instlal` $\to$ suggestions include `install`
  - `gitmap seach` $\to$ suggestions include `search`
  - `gitmap fnd` $\to$ suggestions include `find`
  - `gitmap commti` $\to$ suggestions include `commit`
  - `gitmap ap` $\to$ suggestions include `apps`
- Verify that `suggestTopLevelCommands()` correctly handles exact matches and edge cases.

---

## 5. Acceptance Criteria

- [ ] `collectTopCommandCandidates()` unions `completion.AllCommands()`, `primaryTopCommands`, `"run"`, `"run-until"`, and `macro.ListMacros()`.
- [ ] Candidate pool is deduplicated via `seen` map and contains no empty or whitespace strings.
- [ ] All 512 commands and aliases from `completion.AllCommands()` are valid suggestion targets.
- [ ] `gitmap instlal` suggests `install`.
- [ ] `gitmap seach` suggests `search`.
- [ ] `gitmap fnd` suggests `find`.
- [ ] `gitmap commti` suggests `commit`.
- [ ] `gitmap ap` suggests `apps`.
- [ ] All unit tests in `cli/cmd/rootsuggest_test.go` pass with zero regressions.
