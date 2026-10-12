# Subtask 04 — acceptance-and-release (task update-antigravity-ide, multi-repo)

## Target Repositories
All repository paths are relative to workspace parent (`git-work`):
- `gitmap` (`git-work/gitmap`)
- `scripts-fixer` (`git-work/scripts-fixer`)
- `antigravity-manager` (`git-work/antigravity-manager`)

## File Box
All file paths relative to each respective repository root:
- `gitmap`:
  - `version.json`
  - `package.json`
  - `changelog.md`
  - `02-spec/21-app/update-antigravity-ide/readme.md`
  - `02-spec/21-app/update-antigravity-ide/01-architecture-and-gitmap-command.md`
  - `02-spec/21-app/update-antigravity-ide/02-scripts-fixer-and-live-system-update.md`
  - `.ai-memory/temp-agents/01-update-antigravity-ide/agent-task.db`
  - `.ai-memory/temp-agents/01-update-antigravity-ide/ledger.md`
- `scripts-fixer`:
  - `version.json`
  - `package.json`
  - `changelog.md`
  - `scripts-linux/69-install-antigravity/manifest.json`
  - `scripts/69-install-antigravity/manifest.json`
  - `scripts-linux/69-install-antigravity/run.sh`
  - `scripts/os/ubuntu/install-antigravity.sh`
- `antigravity-manager`:
  - `version.json`
  - `package.json`
  - `changelog.md`
  - `src-tauri/tauri.conf.json`

## Objective & Scope
Subtask 04 orchestrates the comprehensive end-to-end verification, live Ubuntu host update acceptance, multi-repository version bumps, changelog documentation, and atomic release push across all participating repositories using GitMap acceleration tooling.

## Specification References
- `02-spec/21-app/update-antigravity-ide/01-architecture-and-gitmap-command.md`
- `02-spec/21-app/update-antigravity-ide/02-scripts-fixer-and-live-system-update.md`

---

## Detailed Execution Phases

### Phase 1: Live Antigravity IDE Host Verification
1. **Host Environment Sanity Inspection**:
   - Verify that `~/.local/share/antigravity-ide/` contains the updated binary:
     ```bash
     test -x "$HOME/.local/share/antigravity-ide/antigravity"
     file "$HOME/.local/share/antigravity-ide/antigravity" | grep -q "ELF"
     ```
   - Verify that safety backup exists:
     ```bash
     ls -d "$HOME/.local/share/antigravity-ide.bak-"*
     ```
2. **Runner Wrapper & Environment Flags Verification**:
   - Inspect `~/.local/share/antigravity-ide/antigravity-ide.run` and `antigravity-runner.sh`:
     - Must include `export ELECTRON_OZONE_PLATFORM_HINT="auto"`
     - Must execute with `--no-sandbox` to prevent AppArmor unconfined namespace crashes on modern Linux.
     - Must have executable permissions (`chmod +x`).
3. **Symlink Chain Verification**:
   - Validate launcher symlinks in `~/.local/bin/`:
     ```bash
     test -L "$HOME/.local/bin/antigravity"
     test -L "$HOME/.local/bin/agy"
     test -L "$HOME/.local/bin/antigravity-ide"
     ```
   - Validate directory alias symlink:
     ```bash
     test -L "$HOME/.local/share/antigravity"
     ```
4. **Desktop Launcher & User State Verification**:
   - Verify `~/.local/share/applications/antigravity.desktop` exists and references the canonical runner.
   - Confirm user settings directories (`~/.config/Antigravity/` and `~/.antigravity/`) were not overwritten or modified during the upgrade.

---

### Phase 2: Multi-Repository Version Bumps

Follow the single-source-of-truth rule where `version.json` governs all versioning.

1. **`gitmap` Repository**:
   - Inspect current version in `version.json` (`6.547.0`).
   - Increment patch version: `"Version": "6.548.0"`.
   - Update `"releaseDate": "2026-10-12"`.
   - Synchronize dependent manifests (`package.json`, `cli/constants/constants.go`) via repo sync script or targeted update.
2. **`scripts-fixer` Repository**:
   - Inspect current version in `version.json` (`1.57.2`).
   - Increment patch version: `"version": "1.57.3"`.
   - Update `"releaseDate": "2026-10-12"`.
   - Update `package.json` version to match (`1.57.3`).
3. **`antigravity-manager` Repository**:
   - Inspect current version in `version.json` (`4.190.0`).
   - Increment patch version: `"Version": "4.190.1"`, `"version": "4.190.1"`.
   - Update `"releaseDate": "2026-10-12"`.
   - Synchronize `package.json` and `src-tauri/tauri.conf.json` (`4.190.1`).

---

### Phase 3: Changelog Documentation Across Repositories

Add structured changelog entries under date `2026-10-12` in each repository's `changelog.md`:

1. **`gitmap/changelog.md`**:
   ```markdown
   ## [6.548.0] - 2026-10-12

   ### Added
   - `gitmap update antigravity` and `gitmap update ide` commands for automated IDE upgrades.
   - Antigravity tools updater integration and fleet delegation support.
   - Specification and execution plans in `02-spec/21-app/update-antigravity-ide/` and `.ai-memory/plans/subtasks/update-antigravity-ide/`.
   ```

2. **`scripts-fixer/changelog.md`**:
   ```markdown
   ## [1.57.3] - 2026-10-12

   ### Added
   - `update` verb in `scripts-linux/69-install-antigravity/run.sh` and manifest registration.
   - In-place `--update` upgrade workflow in `scripts/os/ubuntu/install-antigravity.sh` with automated backup, atomic unpack, and symlink preservation.
   ```

3. **`antigravity-manager/changelog.md`**:
   ```markdown
   ## [4.190.1] - 2026-10-12

   ### Changed
   - Synchronized Antigravity IDE updater action delegation with GitMap and scripts-fixer toolchain.
   ```

---

### Phase 4: Quality Gates & Compilation Verification

1. **`gitmap` Build & Linters**:
   - Run Go build check:
     ```bash
     go build ./...
     ```
   - Verify exit code `0`. Do not run unneeded test suites (`go test` prohibited in fast gate).
2. **`scripts-fixer` Shell Syntax Checks**:
   - Run syntax checks:
     ```bash
     bash -n scripts-linux/69-install-antigravity/run.sh
     bash -n scripts/os/ubuntu/install-antigravity.sh
     ```
   - Verify exit code `0`.
3. **JSON Manifests Linting**:
   - Verify all modified JSON files parse cleanly without syntax errors.

---

### Phase 5: GitMap Accelerated Atomic Release (`gitmap cpf`)

Commit and push each repository sequentially using the accelerated `gitmap cpf` command. Never invoke raw `git add` or `git commit`.

1. **`scripts-fixer` Commit & Push**:
   ```bash
   gitmap cpf "scripts-fixer - add update verb to 69-install-antigravity and safe ubuntu in-place upgrade"
   ```
2. **`antigravity-manager` Commit & Push**:
   ```bash
   gitmap cpf "antigravity-manager - align antigravity ide updater delegation and release"
   ```
3. **`gitmap` Commit & Push**:
   ```bash
   gitmap cpf "gitmap - implement antigravity ide update command, docs, and subtask plans"
   ```

---

### Phase 6: Task State & Ledger Synchronization

1. **SQLite Database Update**:
   - Update `.ai-memory/temp-agents/01-update-antigravity-ide/agent-task.db`:
     - Mark subtasks as `DONE`.
     - Record final status in `ParentTask`.
2. **Ledger Update**:
   - Update `.ai-memory/temp-agents/01-update-antigravity-ide/ledger.md`:
     - Set `Status: COMPLETED`.
     - Record `Last completed action: Subtask 04 Acceptance and Release Completed`.
     - Record `Workers in flight: none`.

---

## Non-Negotiable Rules

- **Strictly Relative Paths**: All path references must remain relative to repository roots or use portable home paths (`~/.local/...`). Zero absolute paths.
- **Zero Raw Git Commits**: Releases must be executed exclusively via `gitmap cpf`.
- **Atomic Verification**: No version is bumped or pushed without successful build and validation gates.

## Deliverables
- Fully validated host Antigravity IDE deployment.
- Synchronized version manifests and changelogs across all 3 codebases.
- Atomic commit hashes and clean `gitmap cpf` execution confirmation.
