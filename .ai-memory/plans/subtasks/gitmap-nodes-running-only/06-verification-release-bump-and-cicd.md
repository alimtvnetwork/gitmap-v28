# Subtask 06: Verification, Release Bump (`v6.555.0`) & CI/CD Telemetry

- **Subtask ID:** `SUBTASK-06`
- **Parent Plan:** `.ai-memory/plans/gitmap-nodes-running-only.md`
- **Component Spec:** `02-spec/21-app/gitmap-nodes-running-only/02-component-spec.md`
- **Status:** `PENDING`
- **Priority:** High
- **Dependencies:** Subtasks 01 through 05

---

## 1. Objective & Context

This subtask defines the comprehensive quality assurance, static verification, version increment, Git release ceremony, and CI/CD telemetry monitoring for the `gitmap-nodes-running-only` feature package.

The completion of this subtask establishes production readiness by:
1. Executing automated linter scripts across all modified and newly authored Go, TypeScript, and Markdown files.
2. Verifying unit tests for the 10-second liveness cache and open-only filtering logic.
3. Performing a minor version bump from `v6.554.0` to `v6.555.0` synchronized across `version.json`, `package.json`, `cli/constants/constants.go`, and `changelog.md`.
4. Executing an atomic Git commit and push adhering to project standards (`gitmap cpf`).
5. Monitoring remote CI/CD pipeline execution using `gitmap pe -t` telemetry until a green passing state is confirmed.

---

## 2. Owned Files & Target Modifications

| File Path | Description of Changes |
|---|---|
| `version.json` | Increment root `"Version"` field from `"6.554.0"` to `"6.555.0"`. |
| `package.json` | Synchronize `"version"` field to `"6.555.0"`. |
| `cli/constants/constants.go` | Update `AppVersion` constant to `"6.555.0"`. |
| `changelog.md` | Prepend release section for `v6.555.0` documenting `gitmap nodes-running-only` (`nro`), liveness caching, and running-only delegation commands. |

---

## 3. Verification Protocol & Quality Gates

### Phase 1: Static Linter Gates

Execute the canonical Python linter suite to verify code health, nesting limits, boolean conventions, and path hygiene:

```bash
# 1. Check nested-if nesting depth (strict maximum depth <= 1)
python 03-ai-scripts/18-check-nested-ifs.py

# 2. Check positive boolean conventions (is*, has*, can*, should*; no negative naming; no == true)
python 03-ai-scripts/09-check-boolean-guidelines.py

# 3. Check relative path hygiene (zero absolute paths, forward slashes only)
python 03-ai-scripts/14-check-relative-paths.py
```

All three linters must exit with code 0. Any violation must be remediated immediately prior to proceeding.

---

### Phase 2: Unit Testing & Verification Gates

Once code authoring across Subtasks 01–05 is complete, execute targeted unit tests to validate:
- **Cache TTL Expiration**: Verify entries expire exactly after 10 seconds.
- **Cache Persistence**: Ensure JSON serialization and deserialization in OS temp directory succeeds without data corruption.
- **Open-Only Filtering**: Verify `FilterOnlineNodes` excludes offline records and preserves only `is_online: true`.
- **JSON Envelope Parity**: Ensure `gitmap nro --json` output matches the schema of `gitmap nodes --json`.

```bash
# Run targeted unit tests for nodes package
go test -v ./cli/cmdnodes/... -run "TestLivenessCache|TestOpenOnlyFilter"
```

---

### Phase 3: Minor Version Bump Ceremony (`v6.555.0`)

Synchronize the minor version bump across canonical files:

1. **`version.json`**:
   ```json
   {
     "Version": "6.555.0",
     ...
   }
   ```
2. **`package.json`**:
   ```json
   {
     "version": "6.555.0",
     ...
   }
   ```
3. **`cli/constants/constants.go`**:
   ```go
   const AppVersion = "6.555.0"
   ```
4. **`changelog.md`**:
   Prepend release entry:
   ```markdown
   ## [6.555.0] - 2026-10-16

   ### Added
   - `gitmap nodes-running-only` (`gitmap nro`) command returning strictly active, online fleet machines.
   - Cross-process 10-second node liveness caching engine in `cli/cmdnodes/nodes_liveness_cache.go` to eliminate repeated TCP probe timeout penalties.
   - Running-only delegation subcommands (`nro exec`, `nro tasks`, `nro resources`, `nro ping`, `nro pull`, `nro cp all`, `nro scan`).
   - `--running-only` and `-ro` flag support across `gitmap nodes tasks` and `gitmap nodes resources`.
   - Complete Markdown documentation in `cli/helpdoc/nodes-running-only.md` and Web UI catalog integration in `src/data/commands.ts`.
   ```

Alternatively, invoke the version sync automation script:
```bash
python 03-ai-scripts/37-bump-version.py --minor
```

---

### Phase 4: Atomic Git Commit & Push Ceremony

Adhere to the mandatory "No Push = Not Done" completion invariant:

```bash
# Stage and commit feature changes with standard gitmap prefix
gitmap cpf "feat(nodes): implement nodes-running-only nro command with 10s liveness cache and delegation parity"
```

Verify that the working tree is completely clean:
```bash
git status --porcelain
```

---

### Phase 5: CI/CD Telemetry Monitoring via `gitmap pe -t`

Monitor remote GitHub Actions pipeline progress with dynamic waiting:

```bash
# Query CI/CD pipeline telemetry and failure analysis
gitmap pe -t
```

#### Monitoring Protocol
1. Interrogate `gitmap pe -t` to track active run status (`in_progress`, `queued`, `completed`).
2. If the pipeline is running, sleep dynamically according to estimated completion time (typically 30–60s) before re-checking.
3. If any step fails:
   - Extract failure logs using `gitmap pe -t`.
   - Perform 4-part Root Cause Analysis (RCA): Issue Summary, Root Cause, Action Plan, Prevention.
   - Apply surgical fixes and re-commit.
4. Verify all workflow jobs conclude with a green `success` state.

---

## 4. Invariants & Acceptance Gates

- [ ] All three Python static linters exit cleanly with zero violations.
- [ ] Working tree is clean after release commit.
- [ ] Version `6.555.0` is consistently reflected in `version.json`, `package.json`, `constants.go`, and `changelog.md`.
- [ ] CI/CD pipeline monitored to full green completion via `gitmap pe -t`.
