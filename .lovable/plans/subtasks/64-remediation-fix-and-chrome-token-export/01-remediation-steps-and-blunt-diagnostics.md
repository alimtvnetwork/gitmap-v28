# Subtask 01 — Native Remediation Steps & Blunt Error Diagnostics

**Parent Plan:** [.lovable/plans/pending/64-remediation-fix-and-chrome-token-export.md](../../pending/64-remediation-fix-and-chrome-token-export.md)  
**Status:** COMPLETED  
**Files:**
- `gitmap/gitutil/remediation_generator.go`
- `gitmap/gitutil/remediation_commit.go`
- `gitmap/gitutil/remediation_stash.go`
- `gitmap/gitutil/remediation_discard.go`
- `gitmap/cmd/fix_cmd.go`

---

## Objective
Replace brittle shell chaining (`cmd /c "git -C ... && ..."`) with native structured `RemediationStep` slices, eliminating quote stripping on Windows (`pathspec 'local' did not match any file(s)`). Capture stdout/stderr, step exit status, blunt root-cause analysis, and known solution suggestions for failed git commands.

## Requirements
1. Add `RemediationStep` struct (`Name string`, `Args []string`) and field `Steps []RemediationStep` to `RemediationRecipe`.
2. Update `GenerateCommitRecipe`, `GenerateStashRecipe`, and `GenerateDiscardRecipe` to populate `Steps`.
3. In `gitmap/cmd/fix_cmd.go`, execute `recipe.Steps` sequentially without invoking `cmd /c`.
4. Capture stdout and stderr per step. If a step fails, print:
   - Full command string and arguments
   - Exit code
   - Raw output
   - Blunt root-cause diagnosis
   - Actionable known solutions
5. If `commit -m` outputs `nothing to commit`, treat as benign and proceed to `pull --rebase`.
