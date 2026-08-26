# 28 — Worker Push-Update Instruction Document

**Spec:** `14-update`
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** spec/19 audit findings F-X-14, F-X-15, F-X-17 (cross-spec); top-10 fix #5
**Authority:** This file IS the canonical schema for the JSON instruction document referenced by `spec/19-main-worker-service/09-self-update-pointer.md`. Any conflicting prose in spec/19 defers to this file.

---

## 1. Purpose

When the Power Admin triggers `POST /API/V1/Workers/{WorkerNodeId}/Update` (single-target) or `POST /API/V1/Workers/All/Update` (fan-out), Main does **not** stream the binary itself. Instead, Main writes a **JSON Instruction Document** (JID) and POSTs it to each worker's `POST /API/V1/SelfUpdate` endpoint. The worker uses the JID to fetch, verify, stage, and swap its own binary using the standard self-update workflow defined in `spec/14-update/22-update-command-workflow.md` and `spec/14-update/03-rename-first-deploy.md`.

The JID is the **contract** between Main (issuer) and Worker (executor). It MUST be:

- Idempotent — same `InstructionId` re-delivered = no-op.
- Self-describing — worker needs no out-of-band knowledge.
- Verifiable — checksum + signature mandatory.
- Auditable — signed, timestamped, traceable to a single Power Admin user.

---

## 2. Transport

| Aspect | Value |
|---|---|
| HTTP Method | `POST` |
| Path | `/API/V1/SelfUpdate` (on the target worker) |
| Auth | OAuth client-credentials (Main → Worker) per `spec/19/05-auth-and-2fa.md` |
| Content-Type | `application/json; charset=utf-8` |
| Mandatory Headers | `X-Correlation-Id`, `X-Idempotency-Key` (= `InstructionId`), `X-Auth-Action: SelfUpdate` |
| Body | A single JID object (see §3) |
| Max body size | 16 KB |
| Timeout | 30 s handshake (Main → Worker); pinned in `spec/19-main-worker-service/15-tunable-constants.md` §2.5 (`MainWorker.Routing.HttpHandshakeTimeoutSeconds`). |

---

## 3. JID Schema (PascalCase, JSON)

```jsonc
{
  "InstructionId": "01J...ULID",                  // ULID, primary idempotency key
  "InstructionVersion": "1.0.0",                  // SemVer of THIS schema
  "IssuedAtUtc": "2026-05-04T12:00:00Z",          // RFC3339 UTC
  "IssuedByUserId": 7,                            // Power Admin user id
  "IssuedByUserDisplayName": "alim",              // for audit logs
  "TargetWorkerNodeId": 3,                        // 0 = "any" (fan-out broadcast)
  "TargetCurrentVersion": "1.4.2",                // worker MUST match or refuse
  "TargetMinimumVersion": "1.4.0",                // refuse if worker < this
  "PayloadKind": "ZipBundle",                     // see §4
  "PayloadUrl": "https://main.example.com/Releases/worker/1.5.0/worker-1.5.0.zip",
  "PayloadSizeBytes": 18234567,
  "PayloadSha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "PayloadSignatureAlg": "RS256",
  "PayloadSignatureBase64": "MEUCIQDx...==",      // signature over PayloadSha256
  "SigningKeyId": "main-release-2026-q2",         // matches a key in worker's trust store
  "TargetVersion": "1.5.0",                       // version the worker WILL run after swap
  "DeployStrategy": "RenameFirst",                // see §5
  "ExecutionWindow": {
    "EarliestStartUtc": "2026-05-04T16:00:00Z",
    "LatestStartUtc":   "2026-05-04T20:00:00Z",
    "MaxRunDurationSeconds": 600
  },
  "OnFailure": {
    "RollbackPolicy": "AutoRevertOnHandoffFailure", // see §6
    "MaxRetries": 3,
    "RetryBackoffSeconds": [30, 120, 300]
  },
  "Notes": "Quarterly security patch.",
  "Description": null
}
```

### 3.1 Field-by-field

| Field | Type | Nullable | Notes |
|---|---|---|---|
| `InstructionId` | TEXT (ULID 26 chars) | No | Idempotency key. Worker MUST persist + dedupe for ≥ 14 days. |
| `InstructionVersion` | TEXT (SemVer) | No | Worker rejects `Major` mismatch. |
| `IssuedAtUtc` | TEXT (RFC3339 UTC) | No | Worker rejects skew > `WorkerPushUpdate.IssuedSkewSeconds` (default 5 min per `spec/19-main-worker-service/15-tunable-constants.md` §2.7). |
| `IssuedByUserId` | INTEGER | No | FK in Main; opaque to Worker (audit only). |
| `IssuedByUserDisplayName` | TEXT | No | Display only. |
| `TargetWorkerNodeId` | INTEGER | No | `0` = fan-out (worker still verifies own id matches OR is `0`). |
| `TargetCurrentVersion` | TEXT (SemVer) | No | Worker refuses if its current version ≠ this. Prevents stale instructions. |
| `TargetMinimumVersion` | TEXT (SemVer) | No | Lower bound — worker must be ≥ this OR refuse with `WORKER_VERSION_TOO_OLD`. |
| `PayloadKind` | ENUM | No | `ZipBundle` \| `SingleBinary` \| `Patch`. v1.0 only `ZipBundle` is implemented; others reserved. |
| `PayloadUrl` | TEXT (HTTPS URL) | No | MUST be HTTPS. Worker fetches with OAuth bearer (same creds as inbound). |
| `PayloadSizeBytes` | INTEGER | No | Worker aborts if downloaded bytes ≠ this. |
| `PayloadSha256` | TEXT (64 hex) | No | Worker computes locally, MUST match. |
| `PayloadSignatureAlg` | ENUM | No | `RS256` only in v1.0. |
| `PayloadSignatureBase64` | TEXT | No | Signature over `PayloadSha256` bytes (NOT the URL, NOT the JID). |
| `SigningKeyId` | TEXT | No | Worker looks up in local trust store (seeded via Seedable-Config). |
| `TargetVersion` | TEXT (SemVer) | No | The version that will be installed. Worker writes to `latest.json` after handoff. |
| `DeployStrategy` | ENUM | No | `RenameFirst` only in v1.0; future: `BlueGreen`, `InPlace`. |
| `ExecutionWindow.EarliestStartUtc` | TEXT (RFC3339) | No | Worker queues until this time. |
| `ExecutionWindow.LatestStartUtc` | TEXT (RFC3339) | No | Worker drops the JID and reports `INSTRUCTION_EXPIRED` if started after. |
| `ExecutionWindow.MaxRunDurationSeconds` | INTEGER | No | Worker self-aborts if exceeded; triggers rollback. |
| `OnFailure.RollbackPolicy` | ENUM | No | `AutoRevertOnHandoffFailure` \| `ManualOnly` \| `RetryOnly`. |
| `OnFailure.MaxRetries` | INTEGER | No | Single canonical value, default **3**. Pinned in `spec/19-main-worker-service/15-tunable-constants.md` §2.1 (`WorkerPushUpdate.MaxRetries`). Resolves audit F-A-15. |
| `OnFailure.RetryBackoffSeconds` | INTEGER[] | No | Length MUST equal `MaxRetries`. Defaults `[30, 120, 300]` per `spec/19-main-worker-service/15-tunable-constants.md` §2.1. |
| `Notes` | TEXT | No | Free text shown in worker logs. (Per Code Red Rule 11 — transactional table.) |
| `Description` | TEXT | Yes | Optional per Code Red Rule 10. |

---

## 4. PayloadKind catalog

| Kind | v1.0 Implemented? | Description |
|---|---|---|
| `ZipBundle` | ✅ Yes | Standard release zip per `spec/14-update/13-release-assets.md`. Worker unzips, runs `RenameFirst` deploy. |
| `SingleBinary` | ❌ Reserved | Single executable; no zip. Future. |
| `Patch` | ❌ Reserved | Binary diff (bsdiff). Future. |

A worker receiving an unimplemented kind MUST refuse with `INSTRUCTION_KIND_UNSUPPORTED` and NOT delete the JID (so a future worker version can replay it).

---

## 5. DeployStrategy = `RenameFirst` (v1.0 default)

The worker MUST execute the standard rename-first flow defined in `spec/14-update/03-rename-first-deploy.md`:

1. Download `PayloadUrl` to `staging/<TargetVersion>.zip`.
2. Verify size + SHA256 + RS256 signature. Abort on any mismatch → `PAYLOAD_VERIFICATION_FAIL`.
3. Unzip into `staging/<TargetVersion>/`.
4. Rename current binary directory: `current/` → `previous/` (atomic on same filesystem).
5. Rename `staging/<TargetVersion>/` → `current/`.
6. Spawn `current/<binary> --post-update-handoff` (per `spec/14-update/05-handoff-mechanism.md`).
7. On successful handoff: write `current/latest.json` with `TargetVersion`; report success to Main via `POST /API/V1/Workers/{Id}/Heartbeat` with `LastUpdateInstructionId = InstructionId`.
8. On handoff failure within `WorkerPushUpdate.HandoffTimeoutSeconds` (default 60 s per `spec/19-main-worker-service/15-tunable-constants.md` §2.7): revert (rename `current/` → `failed-<TargetVersion>/`, `previous/` → `current/`), report `HANDOFF_FAILED`.

---

## 6. Failure semantics

| Worker outcome | Reported error code (registered per audit F-X-08) | Idempotent retry? |
|---|---|---|
| Signature invalid | `WORKER-600-01 PAYLOAD_VERIFICATION_FAIL` | No — Main must reissue with new JID. |
| SHA256 mismatch | `WORKER-600-02 PAYLOAD_CHECKSUM_FAIL` | No. |
| Size mismatch | `WORKER-600-03 PAYLOAD_SIZE_FAIL` | No. |
| Worker version stale | `WORKER-400-01 WORKER_VERSION_TOO_OLD` | No. |
| Worker version mismatch | `WORKER-400-02 VERSION_MISMATCH` | No. |
| Window expired | `WORKER-400-03 INSTRUCTION_EXPIRED` | No. |
| Disk full / IO | `WORKER-700-01 DEPLOY_IO_FAIL` | Yes (until `MaxRetries`). |
| Handoff timeout | `WORKER-500-01 HANDOFF_FAILED` | Yes (rollback first). |
| Duplicate `InstructionId` | `WORKER-200-01 INSTRUCTION_ALREADY_APPLIED` (200 OK) | N/A — successful no-op. |

All codes MUST be registered in `spec/03-error-manage/03-error-code-registry/` under prefix `WORKER` before this document ships to implementation.

---

## 7. Worker-side persistence

> **Tier correction (FU-5, applied 2026-05-04):** `WorkerUpdateInstruction` is **worker-wide**, not company-scoped, so it lives in the **Settings tier** — NOT the App tier. Authority: [`spec/19-main-worker-service/11-split-db-tier-reconciliation.md`](../19-main-worker-service/11-split-db-tier-reconciliation.md) §5. Earlier drafts of this file placed it in the App tier; that placement is now retracted.

Worker stores each JID in the **Settings tier** DB (per `spec/05-split-db-architecture/`; concrete tier mapping in `spec/19/11-split-db-tier-reconciliation.md` §5):

```sql
CREATE TABLE WorkerUpdateInstruction (
    WorkerUpdateInstructionId    INTEGER PRIMARY KEY AUTOINCREMENT,
    InstructionUlid              TEXT NOT NULL UNIQUE,
    InstructionVersion           TEXT NOT NULL,
    IssuedAtUtc                  TEXT NOT NULL,
    IssuedByUserId               INTEGER NOT NULL,
    TargetWorkerNodeId           INTEGER NOT NULL,
    TargetCurrentVersion         TEXT NOT NULL,
    TargetVersion                TEXT NOT NULL,
    PayloadUrl                   TEXT NOT NULL,
    PayloadSha256                TEXT NOT NULL,
    SigningKeyId                 TEXT NOT NULL,
    DeployStrategyCode           TEXT NOT NULL,
    OutcomeCode                  TEXT NULL,
    OutcomeReportedAtUtc         TEXT NULL,
    Notes                        TEXT NULL,
    Comments                     TEXT NULL,
    Description                  TEXT NULL
);
CREATE INDEX IxWorkerUpdateInstructionUlid
    ON WorkerUpdateInstruction (InstructionUlid);
```

Compliant with Code Red Schema Rules 10/11/12: `Notes`, `Comments`, `Description` all NULL-able with no DEFAULT; no UUIDs; PK is `{TableName}Id`.

Retention: minimum **14 days** (`WorkerPushUpdate.InstructionRetentionDays`); pinned in `spec/19-main-worker-service/15-tunable-constants.md` §2.7 and overridable via Seedable-Config key `WorkerUpdateInstructionRetentionDays`.

---

## 8. End-to-end sequence (informative)

```
Power Admin           Main Server                   Worker
     |                     |                          |
     |--POST Workers/N/Update------->|                |
     |                     |--build JID-------------->|
     |                     |--POST /SelfUpdate ------>|
     |                     |       (JID + headers)    |
     |                     |<--202 Accepted-----------|
     |                     |                          |--persist JID
     |                     |                          |--wait for window
     |                     |                          |--download PayloadUrl
     |                     |                          |--verify size/sha/sig
     |                     |                          |--RenameFirst deploy
     |                     |                          |--handoff
     |                     |<--POST /Heartbeat--------|
     |                     |   {LastUpdateInstructionId, OutcomeCode}
     |<--200 OK------------|                          |
```

**Authority:** prose in this file wins over the diagram on conflict.

---

## 9. Open Questions (logged, not blocking)

- **OQ-28-1** Should `0` (fan-out) be a magic value or a separate `IsBroadcast: true` field? Inferred: keep `0` for wire compactness; revisit if we add multi-tenant routing.
- **OQ-28-2** Should `PayloadSignatureBase64` cover the entire JID instead of just `PayloadSha256`? Inferred: SHA256-only for v1.0 (simpler key rotation); upgrade in v2.0 by bumping `InstructionVersion` major.

---

## 10. Cross-references

- `spec/19-main-worker-service/09-self-update-pointer.md` — calls into this contract.
- `spec/14-update/03-rename-first-deploy.md` — deploy mechanics.
- `spec/14-update/05-handoff-mechanism.md` — post-swap handoff.
- `spec/14-update/22-update-command-workflow.md` — analogous CLI flow (sibling, not authoritative for workers).
- `spec/03-error-manage/03-error-code-registry/` — MUST register `WORKER-*` codes from §6.
- `spec/05-split-db-architecture/` — tier-semantic source of truth.
- `spec/19-main-worker-service/11-split-db-tier-reconciliation.md` — `WorkerUpdateInstruction` placement: **Settings tier** (FU-5).
- `spec/06-seedable-config-architecture/` — `WorkerUpdateInstructionRetentionDays` and `SigningKeyId` trust store.

---

*Worker push-update instruction document v1.0.0 — 2026-05-04*
