# CI Issue 79: Pipeline Misspell and Negative Boolean Guidelines RCA

## 1. Reproduction & Symptoms
- **CI Run ID:** [36108775417](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36108775417)
- **Commit:** `5f4dd0b5`
- **Failing Jobs:**
  1. `Spell Check (misspell, US locale)`: Step 'Run misspell on changed files' failed on `cancelled`.
  2. `Boolean Guidelines Linter`: Step 'Run ./.github/actions/policy-check' failed on `check-boolean-guidelines.py`.

## 2. Root Cause Analysis
1. **Misspell Linter Violation:**
   - In `cli/cmdpipeline/pipeline_stages.go` and `pipeline_stages_render.go`, the literal `"cancelled"` was used for conclusion string matching. The `misspell` auditor enforces US-English spelling (`canceled`).
2. **Boolean Guidelines Violation:**
   - `check-boolean-guidelines.py` bans negative boolean variable prefixes (`hasNo*`, `isNot*` per regex `\b(isNot[A-Z]\w*|hasNo[A-Z]\w*)\b`).
   - The variables `hasNoConvs`, `hasNoQueues`, `hasNoJobs`, and `hasNoRun` violated this rule.

## 3. Code Fix
1. Refactored `pipeline_stages.go` and `pipeline_stages_render.go` to use `strings.HasPrefix(conclusion, "cancel")` and format badges with `"canceled"`.
2. Renamed negative booleans to positive framing (`isEmpty`, `isMissing`).
3. Ran `python linter-scripts/check-boolean-guidelines.py` and `python 03-ai-scripts/27-misspell-auditor.py` to confirm zero violations.

## 4. Prevention
- Always run `python linter-scripts/check-boolean-guidelines.py` locally on all newly created and modified files before committing.
- Avoid British-spelled literals in conclusion badges and comparisons.
