# 125 - Unused injectColumnFor, Gencommands Top-Level Marker & Helpdoc Lemon RCA

## 1. What Failed
- **Commit:** `b37ca939004706545b7215fe2ca42eb19f77965e`
- **Workflows:**
  - `CI` (Run `#39150858184`):
    - `Generate Drift Check`: Step 5 `Run go generate` failed: `gencommands: found 1 Cmd* constant(s) in const block(s) lacking the '// gitmap:cmd top-level' marker: constants_commands_lemon.go:3 CmdLemon, CmdLemonAlias`.
    - `Lint Baseline Guard (unused)`: Step 7 `Unused diff` failed: `[unused] func 'injectColumnFor' is unused (cli/store/inject_idempotency.go:63)`.
    - `Full Suite Guard`:
      - `cli/constants/TestTopLevelCmdRegistryMatchesAST`: Failed because `CmdLemon` and `CmdLemonAlias` were not listed in `topLevelCmds()`.
      - `cli/helpdoc/TestEveryCmdIDHasHelpFile`: Failed because `cli/helpdoc/lemon.md` was missing.

---

## 2. Root Cause Analysis
1. **Unused Function Dead Code:**
   When `cli/store/inject_idempotency.go` was migrated from dynamic `fmt.Sprintf` SQL formatting to static typed query constants (`SQLUpdateInjectDesktopTimestamp` / `SQLUpdateInjectVSCodeTimestamp`), the helper `injectColumnFor` was no longer referenced by any caller, triggering the `unused` linter diff gate.
2. **Missing `gencommands` Opt-in Marker:**
   `cli/constants/constants_commands_lemon.go` declared a `const (...)` block containing `Cmd*` identifiers without the required `// gitmap:cmd top-level` directive comment, causing `gencommands` during `go generate` to exit with code 1.
3. **AST Parity and Helpdoc Coverage Contract:**
   The repository mandates that every top-level CLI command constant defined in `cli/constants/` must:
   - Be registered in `topLevelCmds()` in `cmd_constants_test.go` (or marked with `// gitmap:cmd skip`).
   - Have a corresponding markdown documentation file in `cli/helpdoc/<id>.md` containing an `## Examples` section with fenced code blocks.

---

## 3. Action Taken
1. **Removed Dead Code:**
   Removed `injectColumnFor` from `cli/store/inject_idempotency.go`. Verified with `golangci-lint run --no-config --disable-all --enable=unused ./...` (0 findings).
2. **Opt-in Marker for Code Generator:**
   Added `// gitmap:cmd top-level` comment to `cli/constants/constants_commands_lemon.go`.
3. **Test Registry Parity:**
   Registered `CmdLemon` and `CmdLemonAlias` in `topLevelCmdsPart2()` inside `cli/constants/cmd_constants_registry2_test.go`. Verified `TestTopLevelCmdRegistryMatchesAST` passes.
4. **Created Helpdoc Documentation:**
   Authored `cli/helpdoc/lemon.md` with usage syntax, flag documentation table, and fenced code block examples. Verified `TestEveryCmdIDHasHelpFile` and `TestEveryHelpFileHasExamples` pass.
5. **Pre-Commit Verification:**
   - `go generate ./...` in `cli/` runs cleanly without error.
   - `go test -v ./constants/... ./helpdoc/...` all pass.
   - `python 03-ai-scripts/50-fastgate.py` passed all 4 gates in 0.011s.

---

## 4. Invariant Preserved
- Zero unused dead code.
- 100% AST parity between `Cmd*` constants, code generation metadata, and helpdoc coverage.
- Strictly relative paths across all files.
