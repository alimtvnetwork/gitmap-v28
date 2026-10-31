# Subtask 08: Release Ceremony, Binary Swap & Fleet Sync Verification

- **Parent Plan:** [363-gitmap-synchronization-and-version-control](../../pending/363-gitmap-synchronization-and-version-control.md)
- **Subtask ID:** `TASK-08`
- **Slug:** `gitmap-synchronization-and-version-control - Task 08`
- **Target Package:** `version.json`, `pkg/constants/`, `changelog.md`, `release.md`
- **Status:** `pending`
- **Assigned Worker:** Release Orchestrator Subagent

---

## 1. Goal

Conduct the targeted verification pass, perform zero-downtime binary swap via `gitmap swap-binary`, execute the formal minor version bump release ceremony (`v6.636.0`), and verify cross-fleet synchronization across connected workspace nodes.

---

## 2. In Scope

1. **Targeted Verification Pass:**
   - Execute targeted unit verification covering:
     - Manifest generator and deterministic SHA-256 serialization.
     - 3-Way hash comparison and human conflict protection logic.
     - CLI dispatch routing for `gitmap cg sync`, `gitmap prompt sync`, `gitmap skill sync`.
     - On-the-fly execution prompt templates (`cg execute v6 / v7`).
     - Muse onboarding and help card renderers.
   - Strictly avoid running full repository builds or broad unrequested E2E test suites.

2. **Zero-Downtime Binary Swap:**
   - Compile new binary with single process optimization (`-p 1`).
   - Run `gitmap swap-binary` to safely replace the active binary in `AppData/Local` or `/usr/local/bin` without Windows `ERROR_SHARING_VIOLATION` file locking contention.
   - Verify active binary version via `gitmap --version`.

3. **Release Ceremony (`v6.636.0`):**
   - Update canonical version in `version.json`:
     - `"Version": "6.636.0"`
     - `"releaseDate": "2026-10-31"`
   - Propagate version string to `pkg/constants/version.go` and `cli/constants/constants.go`.
   - Update `changelog.md` and `release.md` documenting:
     - Manifest generation and deterministic SHA-256 hash trees.
     - Non-destructive 3-way hash verification with human conflict protection.
     - New CLI commands: `gitmap cg sync`, `gitmap prompt sync`, `gitmap skill sync`.
     - On-the-fly terminal prompts: `gitmap cg execute v6` and `v7`.
     - Muse onboarding and rich terminal help engine.
     - Tokyo Night themed Lipgloss card rendering.

4. **Semantic Release Commit & Fleet Sync:**
   - Commit changes using mandatory release verb:
     ```bash
     gitmap cpr "v6.636.0 gitmap synchronization and version control"
     ```
   - TOTAL BAN on using `cpf`, `cpc`, or raw `git commit` for release commits.
   - Push commit upstream to origin main.
   - Execute initial fleet sync broadcast using `gitmap cg sync` across connected workspaces.

---

## 3. Out of Scope

- Unrequested regression E2E tests.
- Modifying unrelated database schemas or Git remote SSH keys.

---

## 4. Architectural Invariants

- **Mandatory Semantic Commit Verb:** `gitmap cpr` strictly enforced.
- **Single Source of Truth:** `version.json` is the sole authoritative version definition.
- **Relative Git Paths:** Zero absolute filesystem paths in changelog or release documentation.
- **Positive Booleans:** `isCleanSwap`, `isVerified`, `hasBumpedVersion`.

---

## 5. Verification Checklist

- [ ] Targeted unit tests for sync engine and dispatchers pass with 100% clean status.
- [ ] `gitmap swap-binary` completes atomic replacement and `gitmap --version` reports `6.636.0`.
- [ ] `version.json`, `pkg/constants/`, and `changelog.md` are synchronized.
- [ ] Release committed using `gitmap cpr "v6.636.0 gitmap synchronization and version control"`.
- [ ] Changes pushed to remote tracking branch.
- [ ] Fleet sync verified across connected repositories.
