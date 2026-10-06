# Subtask 02: GitMap Installer and WorkDir Dollar Variables Navigation

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/01-architecture-spec.md)  
> **Status:** `QUEUED`  
> **Target Files:**  
> - `cli/cmd/cd_workdir_resolver.go`  
> - `cli/completion/cdfunction.go`  
> - `cli/constants/constants_cd.go`  
> - `cli/cmdworkdir/workdir_cmd.go`  
> - `cli/cmdworkdir/workdir_test.go`  

---

## 1. Technical Objective

Upgrade the `gitmap cd` destination resolver and shell completion installer to recognize `$work`, `$def`, and related workspace tokens regardless of whether they are expanded by the parent shell or passed as literal strings. Verify multi-workdir registration, listing, and switching across the SQLite Split-DB layer.

---

## 2. Implementation Specifications

1. **Keyword Resolution Enhancement (`cli/cmd/cd_workdir_resolver.go`):**
   - Update `isWorkDirKeyword(name string) bool` to sanitize input:
     - Trim surrounding whitespace.
     - Strip optional leading dollar sign (`$`).
     - Lowercase token comparison.
   - Support keywords:
     - `work`, `$work`
     - `def`, `$def`
     - `default`, `$default`
     - `workdir`, `$workdir`
     - `wd`, `$wd`
   - In `resolveCDWorkDirPath(name string)`:
     - Strip leading `$` before passing to `findWorkDirByNameOrLabel` if no default keyword matches.

2. **Shell Wrapper & Expansion Safety (`cli/completion/cdfunction.go` & `constants_cd.go`):**
   - When a user types unquoted `gitmap cd $work` in a shell where `$work` is unset, the shell evaluates it to empty string (`""`).
   - `handleBareCD()` in `cli/cmd/cd.go` already falls back to `resolveDefaultWorkDirPath()`.
   - When passed quoted (`gitmap cd "$work"` or `gitmap cd '$def'`), the Go CLI receives the literal string with the dollar sign and resolves it via `isWorkDirKeyword()`.
   - Ensure shell wrapper functions (`gcd` and `gitmap`) correctly preserve and forward quoted parameters.

3. **Multi-Workdir Command Verification (`cli/cmdworkdir/`):**
   - Validate CLI commands for managing multiple developer directories:
     - `gitmap workdir add <path> [--label <label>]`: Registers new path.
     - `gitmap workdir ls`: Displays all registered paths and marks active default.
     - `gitmap workdir set <id|path|label>`: Changes active default directory.
     - `gitmap cd <label>`: Navigates directly to any registered work directory by label.

---

## 3. Verification Protocol

```bash
# 1. Test literal dollar resolution
gitmap cd '$work'
gitmap cd '$def'

# 2. Test standard keyword resolution
gitmap cd work
gitmap cd def

# 3. Test multi-workdir registration and listing
gitmap workdir ls
gitmap workdir add /tmp/test-workspace --label test
gitmap cd test
gitmap workdir rm test
```

---

## 4. Acceptance Criteria

- [ ] `gitmap cd $work` and `gitmap cd '$work'` resolve to the active default workspace.
- [ ] `gitmap cd $def` and `gitmap cd '$def'` resolve to the active default workspace.
- [ ] Registered work directory labels resolve directly via `gitmap cd <label>`.
- [ ] Multi-workdir CLI commands (`add`, `ls`, `set`, `rm`) function properly without database locks.
- [ ] Coding guidelines enforced: zero nesting, guard inversion, positive boolean naming (`hasWorkDir`, `isDefault`).
