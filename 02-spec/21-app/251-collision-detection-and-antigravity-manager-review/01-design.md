# Spec 251 — Collision benchmark, agent ls/show, Antigravity review

## Part A — gitmap collision detection & agent observability

### A1. SQLite-as-cache confirmation (must document)
- The collision cache IS SQLite: Tier-1 central DB at
  `<repo>/.ai-memory/temp-agents/ai_agents.db` (resolved via
  `store.ResolveMasterAgentDbPath`), tables `FileClaim` (write claims, one row
  per file) and `CollisionEvent` (history). Written at `claim-files` time;
  released on subtask complete/fail.
- `agent collisions` reads ONLY the central DB (one open + one query). Measured:
  detection core **1.7ms** (instrumented). No Tier-2 scans.
- Honest note: end-to-end command is ~250ms due to ~200ms pre-existing `agent`
  subtree dispatch overhead (affects all agent subcommands, predates 250).
  The dispatch cost is out of scope for this task; the detection itself is ms.

### A2. Benchmark-first-line (new requirement)
- `agent collisions`, `agent stats`, `agent heatmap`, `agent ps` MUST print the
  measured elapsed time as the FIRST line, e.g. `[7ms] agent collisions`.
- Timing covers the command's own work (open DB → render), measured with
  `time.Now()`/`time.Since`, printed before any other output (including tables).
- `--json` mode: include `"elapsed_ms"` as the first key.

### A3. `gitmap agent ls` (new)
- Lists the last 10–15 parent tasks (default 12, `--limit N`), high-level one
  line per task: `ID | SLUG | STATUS | AGENTS | SUBTASKS done/total | UPDATED`.
- `--tree` flag (default ON for ls?): subtasks shown as an indented subtree
  under each parent with their IDs and statuses.
- Every row shows IDs (parent task ID + subtask IDs) for traceability.
- Status values surfaced: ACTIVE, COMPLETED, plus subtask PENDING/IN_PROGRESS/
  DONE/FAILED; stale/broken DBs reported as `BROKEN` (not crashed).
- All from SQLite (Tier-1 + Tier-2), no network. Must stay fast (<300ms for
  15 tasks on typical hardware).

### A4. `gitmap agent show <id-or-slug>` (new)
- Full detail for one parent task: ID, slug, name, status, run directory,
  budget, timestamps, then the subtask list as a subtree (ID, code, title,
  agent role, status, file-box summary, evidence excerpt).
- Accepts parent task ID or slug (canonical resolution, like other commands).

## Part B — Antigravity Manager investigation

### B1. openai.rs attribution (answer in spec)
- Upstream: `lbjlaq/Antigravity-Manager`. Fork: `alimtvnetwork/Antigravity-Manager`.
- `openai.rs` (7,583 lines) authorship: 22/28 commits by "Jeik"
  (jeikliu@outlook.com). Determine: is Jeik the upstream maintainer (lbjlaq
  org) or a fork contributor? Check upstream's openai.rs for the same shape.
- "7.5K lines means what": the OpenAI handler has become an image-processing
  monolith (382 image-related matches, 2 pub fns) — protocol adaptation logic
  that belongs in shared pipeline stages lives in the adapter, violating the
  repo's own "Pipeline First" rule.

### B2. Rescue plan for the leakage (design)
- Extract image normalization/validation (`validate_input_image_limits`,
  `parse_image_data_url*`, `build_image_contents`, `NormalizedInputImage`)
  into `proxy/pipeline/` (new `pipeline/image.rs`) or `proxy/common/`,
  protocol-agnostic.
- Keep in the adapter ONLY: OpenAI-specific parameter mapping.
- No behavior change: extract + rewire, verify with existing proxy tests.

### B3. Open PR #6 "Sync v5" review
- PR: `lbjlaq:main` → `alimtvnetwork:main`, 140 files, +31,206/−9,802.
- Review for: pipeline improvements worth cherry-picking, adapter fixes,
  breaking changes to our fork's divergences. NEVER merge blindly; extract
  ideas, verify logic, adapt.
- Record every considered upstream commit in memory (like the 250 as-built notes).

### B4. Scope spiral fix (permission granted)
- Inventory the 54 modules; classify: core (protocol gateway) vs adjacent
  (email, telegram, supabase, ssh, cloudflared).
- Fix = explicit module tiers + dependency guardrails (adjacent modules may
  not be imported by core proxy code), NOT deletion (hard rule: never delete
  unless asked). Document the tier map.
- Codebase improvements: replace `unwrap()`/`expect()` in proxy hot paths
  with proper error returns (start with handlers/).

## Part C — Release
- gitmap 6.519.0 via `37-bump-version.py`, tag, push. `gitmap update` cannot
  download in this sandbox (blocked) — rebuild binary from source instead and
  report honestly.
