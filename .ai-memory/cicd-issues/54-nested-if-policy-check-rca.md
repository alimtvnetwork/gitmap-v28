# CI/CD Issue 54: Nested If Policy Check Violations in AGY Runner & Help RCA

## 1. Symptom
In GitHub Actions CI run `#35310257604` on commit `3f459d7`, two CI jobs (`Boolean & Enum Linter` and `Nested If Linter`) failed at step `Run ./.github/actions/policy-check` with:

```text
❌ FAIL: Found 2 nested-if / anti-compression violation(s) across 2 file(s):
    cli/cmdagy/agy_help.go:16: Nested if statement found (depth 2 inside conditional block): if renderSubcommandHelp(cmd) {
    cli/cmdpipeline/pipeline_fix_agy_runner.go:16: Nested if statement found (depth 2 inside conditional block): if len(rest) == 0 {
Process completed with exit code 1.
```

---

## 2. Root Cause
1. **`pipeline_fix_agy_runner.go:16`**: Inside `isAgyCompound(first, rest)`, the check `if len(rest) == 0` was directly nested within `if first == "agy"`, violating the maximum nesting depth constraint of 1.
2. **`agy_help.go:16`**: Inside `renderAgyHelp(cmd, args)`, the check `if renderSubcommandHelp(cmd)` was nested within `if cmd != nil && cmd != AgyCmd`.

---

## 3. Resolution
1. **Flattened `isAgyCompound` in `cli/cmdpipeline/pipeline_fix_agy_runner.go`**:
   Used guard clause to return early when `first != "agy"`:
   ```go
   func isAgyCompound(first string, rest []string) bool {
       if first != "agy" {
           return false
       }
       if len(rest) == 0 {
           return true
       }

       return hasFixOrErrorsTarget(rest) || hasHelpToken(rest)
   }
   ```
2. **Flattened Subcommand Help in `cli/cmdagy/agy_help.go`**:
   Extracted helper `handleAgySubcommandHelp` and transitioned `agy_help.go` to the unified `cli/termhelp` framework:
   ```go
   func handleAgySubcommandHelp(cmd *cobra.Command) bool {
       if cmd == nil || cmd == AgyCmd {
           return false
       }
       if renderSubcommandHelp(cmd) {
           return true
       }
       _ = cmd.Usage()

       return true
   }
   ```
3. Verified zero violations with `check-nested-ifs.py`, `check-boolean-guidelines.py`, `check-error-management.py`, and `check-enum-guidelines.py`.

---

## 4. Prevention & Learnings
- **Guard-First Pattern**: Always convert multi-condition nested blocks into early guard returns (`if condition != expected { return ... }`).
- **Policy Enforcement**: Pre-commit policy linters prevent nesting anti-patterns from breaking remote CI gates.
