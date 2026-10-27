# 01 — Overview: Main / Worker Service Architecture

> **/goal** Master and enforce the architectural standards, specifications, and CI/CD validation rules for 19 Main Worker Service.
> **/learn** Read the sequentially ordered specification files in this directory, follow the actionable CI/CD checklist, and apply mandatory rules before generating code.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Read and understand all numbered specifications under `19-main-worker-service/`.
- [ ] `/learn` Adhere strictly to `.ai-memory/folder-structure.md` and `.ai-memory/strictly-avoid.md`.
- [ ] `/goal` Verify zero explicit `true` boolean evaluations and no mixed-polarity conditionals.
- [ ] `/learn` Run all local verification linters via `python 03-ai-scripts/06-cicd-local-runner.py`.

. **CRITICAL AI INSTRUCTION:** This `readme.md` file is the primary entry point for this directory. AI agents MUST read this file first before exploring other files in this folder.

**Spec:** `19-main-worker-service`
**Version:** 1.1.0
**Updated:** 2026-05-06
**Default stack:** Laravel (PHP) — stack-agnostic by design.

---

## 0. Mental Model — Kubernetes-style Control Plane

Treat this topology as a **lightweight Kubernetes analogy**:

| Kubernetes concept | This spec's analogue |
|---|---|
| Control plane (api-server + scheduler) | **Main Server** — catalog, routing, admission, push-update fan-out |
| Worker node (kubelet) | **Worker Server** — runs all business logic under its own split-DB |
| Pod scheduling | Tenant→Worker assignment via `05-worker-routing.md` |
| Node heartbeat / `kubelet → api-server` | Worker `POST /Heartbeat` to Main |
| `kubectl apply` / push deploy | Power-Admin push-update (`17-update-channels.md` §2) |
| Reconciliation loop / desired state | Worker pull-update poll loop (`17-update-channels.md` §3) |
| Image registry pull | Worker fetching release zip from a known URL (`17-update-channels.md` §4) |

Three update channels exist (full spec in `17-update-channels.md`):

1. **Pull from Main (Kubernetes-style reconciliation)** — Worker periodically polls `GET /API/V1/SelfUpdate/Desired` on Main; if the desired version differs from its running version, Worker self-updates from a Main-issued URL. Always available.
2. **Pull from a known release URL** — Worker reads its `UpdateSourceUrl` (Seedable-Config) and fetches the latest release JSON+zip directly. Used as Main-independent fallback.
3. **Push from Main (dev/debug only)** — Power Admin uploads a zip and Main fans it out to Workers. **Disabled in production** (`Env=Production` rejects the endpoint with HTTP 403). Allowed when `Env ∈ {Development, Debug, Staging}`.

The Main↔Worker relationship is **declarative**: Main stores desired state, Workers reconcile toward it. The push channel is an explicit override for dev/staging convenience.

---

## 0.1 Author Mindmaps (source of intent)

The architecture in this folder formalizes these author-drawn mindmaps.
The single best one-page summary is image 04. See [`images/readme.md`](./images/readme.md) for the full index.

![Top-level main / self-update split](./images/01-main-worker-topology.png)

![Full EndPointService overview](./images/04-endpoint-service-full-overview.png)

---

Define a two-tier server topology where a **Main Server** acts as a coordinator (Kubernetes-master analogy) and one or more **Worker Servers** hold all business logic. The Main serves the UI and the React frontend's edge endpoints; Workers do the heavy lifting under their own split-DB.

This spec is the contract any implementer (AI or human) follows to build the topology. Details live in numbered files `01-`…`09-`. Diagrams live in `diagrams/`.

---

## 2. Scope

| In scope | Out of scope |
|----------|--------------|
| Main↔Worker topology, routing, auth handshake | Split-DB internals (see `02-spec/05-split-db-architecture/`) |
| Tenant→Worker mapping in main DB | Seedable-Config internals (see `02-spec/06-seedable-config-architecture/`) |
| Push-update mechanism (main → worker) | Self-update mechanism (pointer only — see `10-self-update-pointer.md`) |
| Role model + `User has access to {EnumPage}` pattern | Generic error handling (see `02-spec/03-error-manage/`) |
| Core API endpoint surface | Per-endpoint business logic |

---

## 3. Stack Flexibility (READ THIS)

The default reference implementation is **Laravel (PHP)**. Every rule in this spec is written stack-agnostically.

> **Future stacks explicitly supported:** .NET, Go (Golang), Python, Node.js, additional PHP frameworks (e.g., Symfony, raw PHP), WordPress as a host.

Implementer obligations regardless of stack:

1. Implement the same REST API surface (`07-core-api-endpoints.md`).
2. Honor the same auth contract (`06-auth-and-2fa.md`).
3. Use the same main-DB schema (`04-main-db-schema.md`) — column names PascalCase, PKs `{TableName}Id INTEGER AUTOINCREMENT`, no UUIDs.
4. Use the same error contract (`09-error-contract.md`) for main↔worker calls.

If a stack-specific deviation is unavoidable, document it in that stack's `99-consistency-report.md` and link back here.

---

## 4. Tenant Root Model

**Decision:** Company-as-root.

- Top-level entity: `Company`.
- A `Company` has many `User`s.
- Worker assignment is per-`Company`. All users of a company route to the same worker.
- Single-user products are modeled as `1 Company : 1 User` (degenerate case). No schema change.

Rationale: matches the verbatim's worked example (`POST /API/V1/Company`, "company-to-worker mapping"). See `29-plan.md` §Decisions.

---

## 5. Two-Tier Topology (one-paragraph version)

**Main Server** serves the React UI, holds a thin SQLite catalog (workers, tenant→worker map, versions), and routes business requests. It runs **no** business logic. **Worker Servers** are independent deployments of the same backend stack; each owns a split-DB (Root / App / Session per `02-spec/05-split-db-architecture/`) and runs all business logic. Both tiers ship with auth, 2FA, session, sign-up, sign-in, and JWT/cookie support out of the box; only the Main serves UI.

Full diagrams in `diagrams/` and details in `02-architecture.md`.

---

## 6. Document Map

A literal AI MUST read every numbered file 00→24 before scoping work. The map below is exhaustive — no file in this folder is optional reading.

| File | Purpose |
|------|---------|
| `29-plan.md` | Phased task list, decisions, open questions |
| `readme.md` | This file |
| `100-improvement-possibilities.md` | Comprehensive catalog of 100 architectural improvement possibilities across 10 functional clusters |
| `02-architecture.md` | Topology, request flow, boundaries |
| `03-glossary.md` | Canonical terms |
| `04-main-db-schema.md` | Main-server SQLite schema (single source of truth for `WorkerNode`, `Company`, `UserDirectory`, audit tables) |
| `05-worker-routing.md` | Selection strategies, caching, failover (no endpoint catalogue — see `06-`) |
| `06-auth-and-2fa.md` | Auth flows, 2FA, JWT/cookie, main↔worker handshake |
| `07-core-api-endpoints.md` | **Single source of truth** for the REST surface (paths, payloads, auth) |
| `08-role-based-dashboards.md` | Roles + `User has access to {EnumPage}` pattern |
| `09-error-contract.md` | **Single source of truth** for the JSON error envelope shape |
| `10-self-update-pointer.md` | Pointer-only doc (no implementation) |
| `11-worker-bootstrap-protocol.md` | Worker first-boot registration + heartbeat (defers WorkerNode schema to `03-`) |
| `12-split-db-tier-reconciliation.md` | Authoritative Main↔spec/05 tier mapping (applied stub) |
| `13-jwt-delivery-contract.md` | JWT issuance, claims, verification, rotation |
| `14-error-codes.md` | Error code catalogue (codes only — envelope owned by `08-`) |
| `15-rbac-and-status-seed.md` | RBAC + status enum seed values |
| `16-tunable-constants.md` | **Single source of truth** for every numeric tunable + cache TTL |
| `17-update-channels.md` | Three update channels: pull-from-main, pull-from-url, dev-only push |
| `18-cascading-roles-and-cache-bin.md` | Role cache bin + cascading roles union |
| `19-backup-nodes.md` | Backup-node registration, propagation, lifecycle |
| `20-incremental-backup-sync.md` | Incremental CDC backup sync (Shape A inline column) |
| `21-backup-encryption-and-keys.md` | Pair-RSA encryption + key rotation |
| `22-backup-endpoints.md` | Backup-tier S2S endpoint catalogue (BE-1..BE-6) |
| `23-backup-apply-logic.md` | Backup apply pipeline (stages, transactions, DLQ) |
| `24-snapshot-storage-and-restore.md` | Snapshot build + restore flow |
| `25-threat-model.md` | Threat-model pointer (no implementation) |
| `26-inherited-rules.md` | Inline blind-AI reference for rules inherited from `02-spec/03-`, `04-`, `05-`, `06-` (Phase 13.3) |
| `27-trust-boundaries-and-isolation.md` | Non-negotiable one-way trust gradients and reverse channel denial matrix |
| `28-git-backup-targets.md` | Isolated Git repositories as out-of-band backup targets |
| `96-linter-audit.md` | ⚠ Historical linter snapshot — see banner at top |
| `97-acceptance-criteria.md` | Verbatim acceptance criteria mapping |
| `98-changelog.md` | Spec version history |
| `99-consistency-report.md` | ⚠ Historical consistency snapshot — see banner at top |
| `audit/` | Independent blind-AI readiness audits (latest = `07-…-v2.md`); earlier files marked STALE |
| `diagrams/` | Mermaid ERDs + sequence diagrams (non-authoritative projections) |

---

## 6.1 Architectural Evolution: 100 Improvement Possibilities

Full detailed specification and remediation catalog: **[`100-improvement-possibilities.md`](./100-improvement-possibilities.md)**.

The catalog documents 100 concrete, grounded improvement possibilities across 10 functional clusters (10 items each) with exact deficiency analysis, concrete engineering remediation, and quantified architectural benefits:

1. **Cluster 01: Core Worker Architecture & Topology Hardening (Items 01–10):** Multi-region active-active worker federation, dynamic ingress admission rate limiting, mTLS inter-tier wire encryption, graceful worker drain protocol, cgroup resource isolation, dual-phase push staging and rollback, ephemeral secret zeroization, distributed static asset offload, cross-tier capability negotiation, and multi-process advisory lock pooling.
2. **Cluster 02: Split-DB Schema, Multi-Tenant Isolation & Migration Integrity (Items 11–20):** Schema version attestation matrix, zero-downtime WAL checkpointing, tenant schema mutation sandboxing, automated orphaned directory pruning, cross-tier migration rollback verification gates, parameterized query bind hardening, tier boundary violation detectors, monotonic rowid optimizations, storage quota admission controls, and pre-migration page integrity verification.
3. **Cluster 03: Authentication, Session Tokens & 2FA/TOTP Resilience (Items 21–30):** Constant-time synthetic authentication, single-use TOTP verification tickets, JWT `kid` grace-period keyring rotation, distributed token revocation Bloom filters, independent Argon2id salted backup codes, client transport fingerprint binding, refresh token family rotation, fine-grained M2M OAuth scopes, WebAuthn/FIDO2 hardware keys, and sliding-window anti-brute-force delays.
4. **Cluster 04: Structured Error Handling, Leaf Error Sovereignty & AppError Integration (Items 31–40):** Monadic Result wire envelopes, leaf error sovereignty preservation, structured field-level validation schemas, compile-time error code registry parity, automatic caller context trace injection, exponential backoff with decorrelated jitter, deep PII error redaction, non-breaking minor envelope versioning, distributed correlation ID context binding, and stack-agnostic AppFault wrappers.
5. **Cluster 05: Backup Replication, Key Rotation & Disaster Recovery (Items 41–50):** CDC sequence gap detection and backfill, hybrid RSA envelope key rotation, atomic transactional CDC application, poison-pill CDC isolation via dead-letter queues (DLQ), scoped append-only Git backup targets, differential snapshots with Merkle tree verification, zero-data-loss failover promotion state machine, at-rest HMAC storage envelope verification, tiered cold-storage compaction, and continuous synthetic disaster recovery drill harness.
6. **Cluster 06: Worker Registration, Heartbeat Leases & Health Checking (Items 51–60):** Cryptographic hardware machine attestation, adaptive sliding-window heartbeat leases, deep synthetic split-DB health probing, flapping node quarantine state machine, graceful worker de-registration drain, multi-dimensional capacity scoring, compressed heartbeat metrics batching, probabilistic early cache refresh (XFetch), split-brain safe mode fallback, and security taint tracking on worker retirement.
7. **Cluster 07: Subdomain Routing, Edge Proxying & Reverse Tunneling (Items 61–70):** Dynamic wildcard DNS-01 TLS provisioning, reverse tunneling gateway for private VPC workers, edge DNS SRV/Anycast direct routing, strict SNI proxy header sanitization, per-tenant dedicated egress IP pools, WebSocket and SSE bidirectional streaming tunnels, geo-distributed latency routing, edge circuit breakers for degraded workers, HTTP/3 QUIC listener upgrades, and vanity subdomain anti-squatting locks.
8. **Cluster 08: Observability, Cognitive Heatmap Telemetry & Structured Logging (Items 71–80):** Native cognitive heatmap telemetry lifecycle hooks, W3C Trace Context and OpenTelemetry span propagation, high-frequency SQLite lock contention counters, dynamic seedable-config log level reload, immutable JSON audit trails for privilege escalation, autonomous CDC latency anomaly detection, per-tenant bandwidth accounting, privacy-preserving PII log filters, W3C Server-Timing latency headers, and distributed cooperative cancellation early warning dumps.
9. **Cluster 09: Cross-Language Coding Guidelines, Positive Booleans & Monadic Results (Items 81–90):** AST linter ban on inverted booleans and explicit true, cross-language monadic Result and AppError envelopes, parameter struct refactoring for >2-3 arguments, mandatory vertical line gaps and clean flow, strict 8-15 line function and 150-300 line file sizing limits, universal PascalCase database and wire naming, compile-time enum exhaustiveness checks, elimination of bare void returns, value semantics over pointer indirection, and immutable request context pattern.
10. **Cluster 10: Chaos Engineering, Test Isolation & E2E Validation Boundaries (Items 91–100):** Synthetic Toxiproxy network partition injection, split-DB disk corruption and lock contention harness, clock skew and NTP drift simulation, malformed and out-of-order CDC fuzzing, injectable OS command mock executors, RAM-backed shared-cache SQLite test rig, blind-AI automated spec conformance suite, mass concurrent registration race stress harness, backup promotion split-brain chaos verification, and rolling zero-downtime canary deployment suites.

---

## 7. Compliance References

> **Blind-AI inline reference (Phase 13.3):** the operative subset of every inherited rule below — enough to implement without leaving this folder — is inlined in **[`26-inherited-rules.md`](./26-inherited-rules.md)**. Read that file first; the source-spec links remain authoritative on conflict.

This spec inherits and does not redefine:

- `.ai-memory/coding-guidelines.md` — function length, zero-nesting, positive booleans, PascalCase, enum-for-Type/Status/Category/Kind
- `02-spec/03-error-manage/` — generic error rules (catch → log → rethrow; `apperror` wrapping; explicit file/op context)
- `02-spec/04-database-conventions/` — DB schema + REST header rules (PascalCase, `{TableName}Id`, no UUIDs as PKs, `X-Correlation-Id` / `X-Idempotency-Key` / `X-Auth-Action`)
- `02-spec/05-split-db-architecture/` — 6-tier split-DB model (Worker uses 4, Main uses 3 per FU-1)
- `02-spec/06-seedable-config-architecture/` — SemVer-merged config seeding for both tiers

---

*Overview v1.0.0 — 2026-05-04*

---

## Verification

_Auto-generated section — see `02-spec/19-main-worker-service/97-acceptance-criteria.md` for the full criteria index._

### AC-GEN-001: General conformance check: Index

**Given** Run the spec health-check against this folder.
**When** Run the verification command shown below.
**Then** Cross-references resolve and the folder contains the required `readme.md`, `97-acceptance-criteria.md`, and `99-consistency-report.md`.

**Verification command:**

```bash
python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-08-30_
