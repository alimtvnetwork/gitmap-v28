# Subtask Plan 05: End-to-End Verification, Dual-Routing Parity & CI Quality Gates

> **Task Reference:** `210-gitmap-push-fix-command-and-auth-recovery`  
> **Subtask ID:** `05-e2e-verification-and-ci-parity`  
> **Status:** Pending Execution  
> **Assigned Agent Role:** Quality Assurance & CI/CD Verification Engineer  

---

## 1. Overview & Objectives

Execute rigorous end-to-end verification, command uniqueness validation, AST parity checks, and CI quality gates across the `gitmap push-fix` implementation. Ensure that all dual-routing variations execute seamlessly, flag parsing handles edge cases, `--dry-run` guarantees zero mutations, and all repository coding standards are strictly verified.

### Key Deliverables:
1. Verify CI constants uniqueness and AST parity guards in `cli/constants/`.
2. Author end-to-end unit tests verifying dual-routing dispatch parity across all 6 entry points (`push-fix`, `pf`, `pushfix`, `fix-push`, `push fix`, `ph fix`).
3. Author verification tests for `--dry-run` ensuring no local or remote git mutations take place.
4. Verify non-interactive execution behavior when `--yes` is supplied or CI environment variables are active.
5. Author automated tests for JSON help schema output (`gitmap help --json`).
6. Execute CI linter scripts and ensure 100% adherence to repository coding guidelines.

---

## 2. Step-by-Step Implementation Steps

### Step 5.1: Verify AST Parity & Constants Uniqueness Guards
1. **Execute Parity Tests**:
   - Run `go test -v ./cli/constants -run TestTopLevelCmdRegistryMatchesAST`.
   - Run `go test -v ./cli/constants -run TestTopLevelCmdConstantsAreUnique`.
   - Run `go test -v ./cli/constants -run TestTopLevelCmdAliasesAreUnique`.
2. **Failure Triage & Resolution**:
   - If `TestTopLevelCmdRegistryMatchesAST` fails, ensure every constant in `constants_cli.go` under `// gitmap:cmd top-level` has a corresponding entry in `topLevelCmds()` in `cmd_constants_test.go`.
   - If `TestTopLevelCmdAliasesAreUnique` reports duplicate `"pf"`, verify that `CmdProfileAlias` in `constants_profile.go` was updated to `"prf"`.

### Step 5.2: Implement Dual-Routing Dispatch Parity Tests
1. **Create `cli/cmdpull/push_fix_e2e_test.go`**:
   - Test command line invocation variations:
     ```go
     package cmdpull

     import (
         "testing"
     )

     func TestDualRoutingEquivalence(t *testing.T) {
         testCases := []struct {
             name string
             args []string
             wantSubcommandStripped bool
         }{
             {"TopLevelLong", []string{"push-fix", "--dry-run"}, false},
             {"TopLevelShort", []string{"pf", "-n"}, false},
             {"TopLevelSolid", []string{"pushfix", "-n"}, false},
             {"TopLevelInverted", []string{"fix-push", "-n"}, false},
             {"SubcommandStandard", []string{"push", "fix", "-n"}, true},
             {"SubcommandAlias", []string{"ph", "fix", "-n"}, true},
         }

         for _, tc := range testCases {
             t.Run(tc.name, func(t *testing.T) {
                 // Verify that all variations resolve to the push-fix handler
                 // with identical effective flags (DryRun: true).
                 var effectiveArgs []string
                 if tc.wantSubcommandStripped {
                     effectiveArgs = tc.args[2:]
                 } else {
                     effectiveArgs = tc.args[1:]
                 }
                 flags, err := ParsePushFixFlags(effectiveArgs)
                 if err != nil {
                     t.Fatalf("unexpected error parsing flags: %v", err)
                 }
                 if !flags.DryRun {
                     t.Errorf("expected DryRun=true for %v, got false", tc.args)
                 }
             })
         }
     }
     ```

### Step 5.3: Verify `--dry-run` Zero-Mutation Invariant
1. **In `cli/cmdpull/push_fix_e2e_test.go`**, add `TestDryRunGuaranteesNoMutation`:
   - Initialize a temporary Git repository fixture using `t.TempDir()`.
   - Setup remote origin pointing to a local bare repository fixture.
   - Run `RunPushFix([]string{"--dry-run"})`.
   - Assert:
     1. Remote URL remains identical before and after.
     2. Remote HEAD commit ref is unchanged.
     3. Local working tree status and reflog are unchanged.
     4. Exit code is `0`.

### Step 5.4: Verify Non-Interactive CI Automation & Flags
1. **In `cli/cmdpull/push_fix_e2e_test.go`**, add `TestNonInteractiveExecution`:
   - Set environment variable `CI=true`.
   - Run push-fix flags parser without `-y`.
   - Assert that non-interactive execution is recognized.
   - Test `-f` / `--force` properly translates to `--force-with-lease` rather than raw `--force`.
   - Test `--ssh` and `--https` precedence when both flags are supplied (SSH wins with stderr notice).

### Step 5.5: Helptext & JSON Schema Verification
1. **Create `cli/cmdpull/push_fix_schema_test.go`**:
   - Verify `helptext.ReadRaw("push_fix")` returns the complete Markdown file.
   - Verify that topic aliases `"pf"`, `"pushfix"`, and `"fix-push"` in `helptext.Print` successfully load `push_fix.md`.
   - Validate help text against expected sections: `## Usage`, `## Flags`, `## Examples`.

### Step 5.6: Run CI Linters & Coding Guidelines Quality Gates
1. **Relative Path Hygiene**:
   - Run `python linter-scripts/check-relative-paths.py` to ensure no absolute paths exist in code or comments.
2. **Coding Guidelines Checker**:
   - Run `python 03-ai-scripts/05-guideline-autofixer.py cli/constants cli/cmd cli/cmdpull cli/helptext --check-only` to ensure compliance with boolean naming, variable naming, and function length.

---

## 3. Verification & Quality Commands Matrix

| Check / Gate | Target Scope | Command | Success Criteria |
|---|---|---|---|
| **AST Parity** | `cli/constants` | `go test -v ./cli/constants -run TestTopLevelCmd` | All 3 tests PASS without duplicate or missing constants |
| **Flag Unit Tests** | `cli/cmdpull` | `go test -v ./cli/cmdpull -run TestPushFixFlags` | 100% flag edge cases pass |
| **Dual-Routing Tests** | `cli/cmdpull` | `go test -v ./cli/cmdpull -run TestDualRouting` | All 6 invocation paths resolve identically |
| **Dry-Run Invariant** | `cli/cmdpull` | `go test -v ./cli/cmdpull -run TestDryRun` | Confirms zero mutations on disk/remotes |
| **Helptext Embedding** | `cli/helptext` | `go test -v ./cli/helptext -run TestCoverage` | Markdown loads cleanly for all aliases |
| **Relative Path Guard** | Entire Repo | `python linter-scripts/check-relative-paths.py` | 0 absolute path violations |
| **Guideline Autofixer** | Touched Packages | `python 03-ai-scripts/05-guideline-autofixer.py cli/cmdpull --check-only` | Clean report (0 violations) |
