# Subtask 01: GitMap CLI Antigravity IDE Update Command Implementation & Tests

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/update-antigravity-ide.md`  
> **Target Path:** `.ai-memory/plans/subtasks/update-antigravity-ide/01-gitmap-ide-update.md`  
> **Status:** READY  
> **Estimated Effort:** 3 Engineering Units  
> **Dependencies:** `02-spec/21-app/update-antigravity-ide/01-architecture-and-gitmap-command.md`

---

## 1. Objectives & Scope

1. Implement first-class CLI routing for updating Google Antigravity Desktop IDE across all supported command aliases:
   - `gitmap update antigravity`
   - `gitmap update ide`
   - `gitmap update antigravity-ide`
   - `gitmap update agy-ide`
   - `gitmap update antigravity-desktop`
2. Implement shortcut interception for `gitmap agy update ide` to route directly to IDE application updating rather than failing on project ID lookups.
3. Verify behavioral consistency with `gitmap in ide -f` / `gitmap install antigravity --force`.
4. Update argument parsing, suggestions, and help text in `cli/cmd/rootutility_update_args.go`.
5. Connect dispatch logic in `cli/cmd/rootutility_update.go` to `cmdinstall.RunInstallAntigravityWithOpts(cmdinstall.InstallOptions{Force: true})`.
6. Author comprehensive unit test suite in `cli/cmd/rootutility_update_test.go` or `cli/cmd/cmd_update_ide_test.go`.
7. Enforce zero absolute paths and adhere strictly to repository coding guidelines.

---

## 2. Detailed Technical Plan

### 2.1 Argument Parsing & Token Recognition (`cli/cmd/rootutility_update_args.go`)

- **Define Token Predicates**:
  - Implement `isAntigravityIdeToken(arg string) bool`:
    - Matches case-insensitively: `"antigravity"`, `"ide"`, `"antigravity-ide"`, `"agy-ide"`, `"antigravity-desktop"`.
  - Implement `isAntigravityIdeUpdateTarget(args []string) bool`:
    - Iterates over positional non-flag arguments and detects presence of any valid IDE token.
- **Update Target Whitelist**:
  - In `isKnownUpdateTargetOrFlag(token string) bool`:
    - Add `"antigravity"`, `"ide"`, `"antigravity-ide"`, `"agy-ide"`, `"antigravity-desktop"` to the switch cases.
- **Enhance Suggestion Engine**:
  - In `suggestUpdateTarget(token string) string`:
    - Map common prefixes and abbreviations (`"anti"`, `"antigrav"`, `"id"`, `"idea"`, `"desktop"`) to `"ide"`.
- **Package Resolution for Fleet / SSH**:
  - In `resolveUpdatePackage(args []string) string`:
    - Return `"antigravity"` when `isAntigravityIdeUpdateTarget(args)` evaluates to true.
- **Update Help Text**:
  - In `updateUsageHelpText`:
    - Add `gitmap update ide - Update Google Antigravity Desktop IDE (alias: antigravity)`.

### 2.2 Routing & Dispatch Logic (`cli/cmd/rootutility_update.go`)

- **Wire Dispatch in `executeUpdateHelpFlow(args []string)`**:
  - Place `isAntigravityIdeUpdateTarget(args)` check immediately before `isAgmUpdateTarget(args)`.
  - Ensure remote target flags (`--remote`, `--node`, `-r`, `-n`) are extracted cleanly prior to target evaluation.
- **Implement Target Runner `runUpdateAntigravityIdeTarget(args []string) error`**:
  ```go
  func runUpdateAntigravityIdeTarget(args []string) error {
      remoteTarget, cleanArgs := extractRemoteUpdateTarget(args)
      if remoteTarget != "" {
          return cmdssh.RunSSHUpdateCLI([]string{"ide", remoteTarget})
      }
      if isVersionListRequest(cleanArgs) {
          return cmdinstall.RunListAntigravityVersions()
      }
      ver := resolveEffectiveVersion("antigravity", extractVersionFromArgs(cleanArgs))
      opts := cmdinstall.InstallOptions{
          Tool:    "antigravity",
          Force:   true,
          DryRun:  hasDryRunArg(cleanArgs),
          Version: ver,
          Verbose: hasVerboseArg(cleanArgs),
      }
      return cmdinstall.RunInstallAntigravityWithOpts(opts)
  }
  ```
- **Version Pinning Support**:
  - In `applyPinIfRequested(args, targetVer)`:
    - Add check for `isAntigravityIdeUpdateTarget(args)` and record pin for app name `"antigravity"`.

### 2.3 Antigravity Shortcut Interception (`cli/cmdagy/agy_cmd.go`)

- In `tryDispatchAgyShortcut(args []string) (error, bool)`:
  - Add inspection check for compound command `agy update ide` or `agy update antigravity`:
  ```go
  if len(args) >= 2 && strings.EqualFold(args[0], "update") && isAntigravityIdeToken(args[1]) {
      opts := cmdinstall.InstallOptions{
          Tool:    "antigravity",
          Force:   true,
          DryRun:  hasDryRunArg(args[2:]),
          Version: extractVersionFromArgs(args[2:]),
          Verbose: hasVerboseArg(args[2:]),
      }
      return cmdinstall.RunInstallAntigravityWithOpts(opts), true
  }
  ```

### 2.4 Unit Test Suite (`cli/cmd/cmd_update_ide_test.go`)

Author dedicated test cases covering:
1. `TestIsAntigravityIdeToken`: Validate positive and negative token matches.
2. `TestIsAntigravityIdeUpdateTarget`: Validate flag permutations, order of arguments, and clean extraction.
3. `TestSuggestUpdateTarget_Ide`: Validate suggestion responses for typos like `anti` and `idea`.
4. `TestIsKnownUpdateTarget_Ide`: Ensure `isKnownUpdateTargetOrFlag` accepts all IDE tokens without error.
5. `TestUpdateIdeDryRunDispatch`: Execute `runUpdateAntigravityIdeTarget([]string{"ide", "--dry-run"})` and verify no filesystem side effects.
6. `TestUpdateIdeVersionOverride`: Verify version string is parsed and forwarded to `InstallOptions`.

---

## 3. Verification Criteria & Acceptance Checklist

- [ ] All IDE update aliases (`antigravity`, `ide`, `antigravity-ide`, `agy-ide`) recognized by `isAntigravityIdeUpdateTarget`.
- [ ] `gitmap update ide --dry-run` executes without error and displays planned download/destination.
- [ ] `gitmap update unknown-token` displays suggestion `Did you mean: gitmap update ide?` when given `anti`.
- [ ] `gitmap agy update ide` successfully intercepted and directed to `RunInstallAntigravityWithOpts`.
- [ ] `cli/cmd/rootutility_update.go` and `cli/cmd/rootutility_update_args.go` pass all unit tests.
- [ ] Zero absolute paths: Verified via `python3 linter-scripts/check-relative-paths.py`.

---

## 4. File Bounding Box & Impact Summary

| File Path (Strictly Relative) | Action | Scope / Change Description |
| :--- | :--- | :--- |
| `cli/cmd/rootutility_update.go` | Modify | Add `runUpdateAntigravityIdeTarget`, update `executeUpdateHelpFlow`, pin handling |
| `cli/cmd/rootutility_update_args.go` | Modify | Add `isAntigravityIdeToken`, `isAntigravityIdeUpdateTarget`, update whitelist/suggestions/help |
| `cli/cmdagy/agy_cmd.go` | Modify | Intercept `update ide` shortcut in `tryDispatchAgyShortcut` |
| `cli/cmd/cmd_update_ide_test.go` | Create | Unit test coverage for update ide token parsing, suggestions, and dry-run |
