# 26 — Trust Boundaries & One-Way Isolation

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-08
**Status:** Authoritative (Phase 14)
**Scope:** Defines the **non-negotiable, one-way trust gradients** between Main, Worker, and Backup tiers so that compromise of any lower-trust node CANNOT escalate into a higher-trust node.

> **CODE RED.** This chapter is normative. Any future endpoint, CDC payload, push-update, or operational tool that violates the matrix in §2 is itself the bug — the matrix wins. No exceptions, no "temporary" reversed channels, no "admin override" backdoors.

---

## 1. Why this chapter exists

Earlier chapters define the *happy-path* relationships:

- `04-worker-routing.md` — Main routes tenants to Workers.
- `18-backup-nodes.md` — Backups passively mirror their primary Worker's App tier.
- `19-incremental-backup-sync.md` — Primary Worker pushes CDC envelopes to Backup.
- `28-worker-push-instruction.md` (`spec/14-update`) — Main pushes update instructions to Workers.

What was **never explicitly written down** until now: **every one of those channels is one-way at the trust layer**, and the reverse direction MUST NOT exist. A worker cannot reach back into Main. A backup cannot reach back into its primary worker. A primary worker cannot reach into its own backup beyond the narrow CDC seal/push surface.

This chapter consolidates those invariants into one auditable table so a mediocre AI implementer cannot accidentally open a back-channel.

---

## 2. Trust matrix (authoritative)

Read as: *Row* may perform *Action* against *Column*.

| From ↓ / To → | **Main Server** | **Primary Worker** | **Backup Node** | **Git Backup Target** (ch. 27) |
|---|---|---|---|---|
| **Main Server** | self | ✅ Push update instructions (ch. 14/28); ✅ Mint Worker JWT (ch. 12); ✅ Issue routing decisions (ch. 4) | ✅ Issue failover-promotion instruction; ✅ Read heartbeat | ❌ Never (Main does not hold git creds) |
| **Primary Worker** | ❌ **NO write of any kind.** Only: heartbeat (status only), JWT introspection callback (read-only), update-ack POST. See §3.1. | self | ✅ Seal + push CDC envelopes (ch. 19/22); ❌ MUST NOT mutate Backup's Root/Settings/Session tiers | ✅ Push to its **own** repo namespace only (ch. 27); ❌ MUST NOT touch any other worker's repo |
| **Backup Node** | ❌ **NO contact whatsoever** (D9 invariant + §3.2) | ❌ **NO outbound to primary.** Backup cannot push, instruct, query, or invalidate the primary. See §3.3. | self (rebuilds local apply state only) | ❌ Never (only the primary holds git creds; backup cannot push to git) |
| **Git Backup Target** | ❌ Read-only artifact store; cannot initiate calls | ❌ Cannot call back into worker | ❌ Cannot call back into backup | self |

Legend: ✅ = explicitly allowed and specified elsewhere. ❌ = MUST be physically impossible (no credential, no route, no firewall hole, no listening endpoint).

---

## 3. The three isolation invariants

### 3.1 Worker → Main: read-only beacon, never a writer

**Invariant W→M.** A compromised Worker MUST NOT be able to mutate ANY Main-side state (no row insert/update/delete, no config change, no role grant, no other-tenant routing, no JWT minting, no update-channel pin, no other-worker disable).

Allowed Worker→Main calls (the **complete** list):

| Endpoint | Purpose | Why it cannot escalate |
|---|---|---|
| `POST /API/V1/Worker/Heartbeat` | Liveness + load metrics | Body fields are numeric/enum only; Main treats them as *hints*, never as commands. |
| `POST /API/V1/Worker/UpdateAck` | Acknowledge a push-update instruction | Acks reference an instruction Main itself issued; Worker cannot fabricate new instructions. |
| `GET  /API/V1/Worker/JwtKeySet` | Fetch current JWT signing pubkey | Read-only; rotation is Main-driven. |

Forbidden (compromise-attempt error codes, allocated for future detection middleware in `13-error-codes.md`):

- `MAIN-960-01 WorkerWriteAttempt` — Worker JWT used against any non-allow-listed endpoint.
- `MAIN-960-02 WorkerCrossTenantQuery` — Worker queries data for a tenant not assigned to it.
- `MAIN-960-03 WorkerImpersonatesMain` — Worker presents a token claiming role `Power` / `MainAdmin`.

**Network posture.** Workers connect *outbound* to Main over mTLS; Main does NOT keep an inbound socket open from Worker → Main beyond these three endpoints. Reverse SSH tunnels, WebSocket upgrades, and long-poll channels are explicitly forbidden.

### 3.2 Backup → Primary Worker: zero outbound

**Invariant B→W.** A compromised Backup Node MUST NOT be able to reach its primary Worker at all. Not to query, not to invalidate, not to "request resync", not to push a fake CDC reversal.

Rationale: backups are write-receivers. They have no legitimate reason to call upstream. If the primary needs to recover from a backup, it is a **Power-Admin failover-promotion event** (§3.5 of ch. 18), executed via Main — never as a backup-initiated call.

Allowed Backup outbound calls:

| Target | Endpoint | Purpose |
|---|---|---|
| Main | `POST /API/V1/Worker/Heartbeat` (with `IsBackup=1`) | Liveness only |
| Main | `POST /API/V1/Backup/Apply/Ack` | ACK envelope apply (per `21-backup-endpoints.md`) |

Everything else — including DNS resolution of the primary worker's FQDN — SHOULD be firewalled at the host level. See `24-threat-model.md` (when materialised) for the firewall recipe.

### 3.3 Primary Worker → Backup: narrow CDC seal channel only

**Invariant W→B.** A compromised primary Worker MUST NOT be able to corrupt, hijack, or wipe its Backup Node beyond the CDC envelope channel (which is itself sealed + idempotent + apply-ordered).

What this means in practice:

- The Worker may call **only** `POST /API/V1/Backup/Envelope` (BE-1) and `GET /API/V1/Backup/Watermark` (BE-2). No "DELETE all rows", no "DROP table", no "set IsBackup=0".
- The Backup node MUST reject any envelope that fails the seal check (HMAC over current Active KeyEpoch, ch. 20). A compromised worker that does not hold the current key cannot forge envelopes at all.
- Backup's Root/Settings/Session tiers are **backup-local** (per ch. 18 §1) and there is no endpoint to write them remotely. A compromised primary cannot poison the backup's local config.
- Snapshot/restore endpoints (BE-3..BE-6) are **operator-only**, gated by Power-Admin JWT, and never reachable from a Worker JWT.

If a worker compromise *is* later discovered, the Backup remains a clean recovery source as long as the **previous KeyEpoch** has not yet been retired — which is exactly the rolling-window guarantee in `20-backup-encryption-and-keys.md`.

### 3.4 Symmetric statement

These three invariants together enforce the user's stated rule, verbatim:

> *"If a worker gets infected, it cannot hijack the worker backups. The same way, if a backup gets interrupted or hacked, it cannot interrupt the real worker."*

Compromise propagation is **strictly downhill** in the trust gradient: Main → Worker → Backup → Git target. Never the reverse.

---

## 4. Enforcement layers

| Layer | Mechanism | Failure mode |
|---|---|---|
| **Network** | mTLS + per-role allow-listed endpoints; backup hosts firewall outbound except Main. | Connection refused at TCP layer. |
| **AuthN** | Distinct cert SANs for Main / Worker-primary / Worker-backup; JWTs signed by Main only. | `WORKER-401-01 InvalidPeerCert`. |
| **AuthZ** | Endpoint allow-list keyed on `(role, IsBackup)` matrix in §2; deny-by-default. | `MAIN-960-0x WorkerWriteAttempt` family (§3.1). |
| **Crypto** | CDC envelope HMAC under Active KeyEpoch (ch. 20). | Backup rejects envelope; alerts ApplyDeadLetter. |
| **Audit** | Every `MAIN-960-*` raises an `AuditEvent` row tagged `IsolationViolation`. | Operator-visible in dashboard (ch. 7). |
| **Lint** | `linter-scripts/check-trust-matrix.py` (future, post-v1.0) parses every endpoint chapter and asserts the matrix in §2 holds. | CI fail. |

---

## 5. Open questions deferred to v2.0

| OQ | Question | Why deferred |
|---|---|---|
| **OQ-26-1** | Should we add an *automated* tripwire that quarantines a Worker when `MAIN-960-*` exceeds threshold N within window T? | Requires threat-model chapter (`24-threat-model.md`) to define N, T, and quarantine semantics. |
| **OQ-26-2** | Should the Backup hold a *signed transcript* of the last K applied envelopes so a post-incident audit can prove non-divergence? | Pyramid-snapshot decision (OQ-23-1) must land first. |

Until those land, the matrix in §2 + the three invariants in §3 are the authoritative posture.

---

## 6. Cross-references

- `01-architecture.md` — high-level tier diagram (to be annotated in next pass).
- `04-worker-routing.md` §1.4 — `IsPrimary` guard (already excludes backups from routing).
- `12-jwt-delivery-contract.md` — JWT minting is Main-only.
- `18-backup-nodes.md` §3 — failover-promotion is Power-Admin-driven, not backup-initiated.
- `19-incremental-backup-sync.md` — defines the BE-1/BE-2 channel referenced in §3.3.
- `20-backup-encryption-and-keys.md` — KeyEpoch rotation that bounds primary-worker compromise blast radius.
- `27-git-backup-targets.md` — companion chapter for the *outbound-only* git-push lane.

---

*Phase 14 (2026-05-08) — locks the worker→main, backup→worker, worker→backup non-modification rules.*
