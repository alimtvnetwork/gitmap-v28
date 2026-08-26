# 21 — Backup Endpoints Contract

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative (spec-only, per plan §Mode)
**Resolves:** Phase 9 of the Backup-Nodes initiative — wire surface for everything sealed in Phase 8.
**Depends on:** [`05-auth-and-2fa.md`](./05-auth-and-2fa.md) §S2S, [`06-core-api-endpoints.md`](./06-core-api-endpoints.md) §1, [`08-error-contract.md`](./08-error-contract.md), [`13-error-codes.md`](./13-error-codes.md), [`15-tunable-constants.md`](./15-tunable-constants.md), [`18-backup-nodes.md`](./18-backup-nodes.md), [`19-incremental-backup-sync.md`](./19-incremental-backup-sync.md), [`20-backup-encryption-and-keys.md`](./20-backup-encryption-and-keys.md).

---

## Keywords

`backup-endpoints` · `s2s-oauth` · `incremental-diff` · `rotate-keys` · `restore-by-date` · `snapshots` · `health`

---

## 1. Purpose

Pin the five HTTP endpoints that connect Main, primary Workers, and Backup nodes together. Apply logic (decrypt → unzip → SQLite dispatch) is Phase 10; snapshot storage / retention is Phase 11. This file is **wire only** — request shape, response shape, auth row, error codes, and timeout cite.

---

## 2. Endpoint catalog (paste-ready merge into `06-core-api-endpoints.md` §2)

### 2.X Backup (Backup-node side, Main-triggered)

| # | Method | Path | Hosted on | Auth | Purpose |
|---|---|---|---|---|---|
| BE-1 | POST | `/API/V1/Backup/IncrementalDiff` | Backup node | S2S OAuth + `X-Primary-Worker-Identity` | Receive a sealed envelope (Phase 7+8 output). Returns ACK with `LastAcceptedSyncOpSeq`. |
| BE-2 | POST | `/API/V1/Backup/RotateKeys` | Backup node | S2S OAuth (Main only) | Step S2/S3/S6 of the rotation flow per `20-…` §7.2. |
| BE-3 | POST | `/API/V1/Backup/RestoreByDate` | Backup node | S2S OAuth (Main only) | Begin a date-anchored full restore — backup unpacks the stored snapshot and ships it back. |
| BE-4 | GET | `/API/V1/Backup/Snapshots` | Backup node | S2S OAuth (Main only) | List available snapshot dates (machine-readable inventory). |
| BE-5 | GET | `/API/V1/Backup/Health` | Backup node | S2S OAuth (Main only) | Liveness + lag + key-epoch state for dashboards. |

**Defence-in-depth:** every backup endpoint MUST also assert the request reached the node via the dedicated `s2s` reverse-proxy hostname per `18-backup-nodes.md` §1 D9. Inbound user-facing traffic returns `421 Misdirected Request` + `MAIN-800-04 TrafficOnBackupRejected` at the proxy layer **before** routing reaches these handlers.

---

## 3. Common request envelope

Per `06-core-api-endpoints.md` §1 — every request:

| Header | Required | Notes |
|---|---|---|
| `Authorization: Bearer <S2S-Token>` | yes | Per `05-auth-and-2fa.md` §S2S; client-credentials grant. |
| `X-Correlation-Id` | yes | ULID. Echo in response + every log line. |
| `X-Idempotency-Key` | yes (POST/PUT/PATCH) | TTL pinned by `MainWorker.IdempotencyKeyTtlSeconds` (§2.2 of `15-…`). |
| `X-Primary-Worker-Identity` | BE-1 only | Asserts which primary is shipping. Backup MUST verify it matches the registered pairing (`KnownBackupNode` per `18-…` §6). |

Bodies are JSON with PascalCase keys **except** BE-1 whose body is a `multipart/form-data` upload (the binary zip envelope + a small metadata JSON part). Envelope binary stays sealed per Phase 8.

---

## 4. BE-1 — `POST /API/V1/Backup/IncrementalDiff`

**Purpose:** Deliver one sealed envelope produced by Phase 7+8 to its target backup.

### 4.1 Request

`Content-Type: multipart/form-data; boundary=…`

| Part name | Content-Type | Body |
|---|---|---|
| `Metadata` | `application/json` | `{ "EnvelopeId": "...", "KeyEpoch": 7, "EnvelopeTimestampEpoch": 1746528000, "RowCount": 142, "PrimaryWorkerNodeId": 12, "BackupWorkerNodeId": 27 }` |
| `Envelope` | `application/zip` | The AES-256-zipped sealed envelope per `20-…` §4 step 9. |

### 4.2 Response — `200 OK`

```jsonc
{
  "EnvelopeId":              "01J...ULID",
  "Accepted":                true,
  "LastAcceptedSyncOpSeq":   918374,
  "ReceivedAtEpoch":         1746528002,
  "BackupWorkerNodeId":      27
}
```

`LastAcceptedSyncOpSeq` is the watermark the primary will store in `BackupSyncWatermark.LastAckedSyncOpSeq` (per `19-…` §3.1) and use for the next pass.

### 4.3 Errors

| Trigger | Code |
|---|---|
| `X-Primary-Worker-Identity` does not match `KnownBackupNode` row | `MAIN-800-04 TrafficOnBackupRejected` (proxy) → fallback `WORKER-300-04 BackupNotRoutable` |
| Envelope sealed under unknown / discarded / pending epoch | `WORKER-920-03 / 02 / 01` per `20-…` §9 |
| Outer zip not AES-256 | `WORKER-920-04 ZipCipherTooWeak` |
| RSA-PSS signature verify failure | `WORKER-920-05 EnvelopeSignatureInvalid` |
| Watermark regression (`SyncOpSeq` ≤ `LastAcceptedSyncOpSeq`) | `WORKER-910-01 BackupSyncWatermarkInconsistent` |
| Envelope build read-side failure on receipt persistence | `WORKER-910-02 BackupEnvelopeBuildFailed` |

### 4.4 Idempotency

Replays of the same `X-Idempotency-Key` MUST return the original `200 OK` body verbatim from the worker idempotency store (per `06-core-api-endpoints.md` §1). A replay with a **different** envelope binary returns `MAIN-300-04 IdempotencyBodyMismatch` (409).

---

## 5. BE-2 — `POST /API/V1/Backup/RotateKeys`

**Purpose:** Drive one step of the Pair-RSA rotation flow per `20-…` §7.2. Main calls this on the **backup** for steps S3 (mint pair) and S6 (activate). The matching primary endpoint is hosted on the primary Worker via the existing worker push channel — not duplicated here.

### 5.1 Request

```jsonc
{
  "RotationId":              "01J...ULID",
  "Step":                    "Generate" | "Activate",
  "KeyEpoch":                8,
  "PrimaryWorkerNodeId":     12,
  "BackupWorkerNodeId":      27,
  "PeerPublicKeyPem":        "-----BEGIN PUBLIC KEY-----\n...",   // present on Step=Generate
  "Reason":                  "Scheduled" | "OperatorForced" | "Compromise"
}
```

### 5.2 Response — `200 OK`

```jsonc
{
  "RotationId":              "01J...ULID",
  "Step":                    "Generate",
  "KeyEpoch":                8,
  "OwnPublicKeyPem":         "-----BEGIN PUBLIC KEY-----\n...",   // present on Step=Generate
  "AttestationSignature":    "Base64(...)",
  "BackupKeyEpochState":     "Pending" | "Active",
  "AcknowledgedAtEpoch":     1746528100
}
```

### 5.3 Errors

| Trigger | Code |
|---|---|
| Step ACK budget exceeded (server-side guard) | `MAIN-820-01 RotationStepTimeout` |
| `Pending` row already exists for this pairing on `Step=Generate` | `MAIN-820-02 RotationAlreadyInProgress` |
| `Activate` arrives but counterpart never reached `Pending` | `MAIN-820-03 RotationActivationSplitBrain` |

### 5.4 Timeouts

| Tunable | Default | Purpose |
|---|---:|---|
| `MainWorker.Backup.Endpoint.RotateKeysTimeoutSeconds` | **30** | Per-call HTTP timeout on the Main side, well below `RotationAckTimeoutSeconds=120`. |

---

## 6. BE-3 — `POST /API/V1/Backup/RestoreByDate`

**Purpose:** Main initiates a date-anchored full restore. Backup decompresses the named snapshot and POSTs it to the primary's restore-inbox endpoint (defined in Phase 11). This call returns immediately with a `RestoreJobId`; final status is polled or pushed via standard worker job channels.

### 6.1 Request

```jsonc
{
  "SnapshotDate":            "2026-05-01",            // YYYY-MM-DD (UTC)
  "TargetPrimaryWorkerNodeId": 12,
  "Reason":                  "OperatorRecovery" | "DriftRepair" | "Compliance",
  "ForceOverwrite":          false                    // Phase 11 will pin the safety semantics
}
```

### 6.2 Response — `202 Accepted`

```jsonc
{
  "RestoreJobId":            "01J...ULID",
  "SnapshotDate":            "2026-05-01",
  "EstimatedSizeBytes":      4194304,
  "AcceptedAtEpoch":         1746528200
}
```

### 6.3 Errors

| Trigger | Code |
|---|---|
| Snapshot for date not present | `MAIN-830-01 SnapshotNotFound` (Phase 11 will register the catalogue entry) |
| `TargetPrimaryWorkerNodeId` does not match `BackupOf` pairing | `MAIN-800-04 TrafficOnBackupRejected` (proxy) |
| Concurrent restore already in progress for this pairing | `MAIN-830-02 RestoreAlreadyInProgress` (reserved Phase 11) |

### 6.4 Note on `MAIN-830-*` placeholders

Codes `MAIN-830-01` / `MAIN-830-02` are **reserved here, defined in Phase 11**. They are pre-allocated in §11 below to lock the wire contract; their final messages and HTTP statuses are owned by `22-snapshot-storage-and-restore.md`.

---

## 7. BE-4 — `GET /API/V1/Backup/Snapshots`

**Purpose:** Machine-readable inventory of available snapshots for the dashboard.

### 7.1 Query string

| Name | Required | Default | Notes |
|---|---|---:|---|
| `From` | no | _(min)_ | Inclusive ISO date `YYYY-MM-DD`. |
| `To` | no | _(today)_ | Inclusive. |
| `Limit` | no | **90** | Hard ceiling regardless of date range. |

### 7.2 Response — `200 OK`

```jsonc
{
  "BackupWorkerNodeId":      27,
  "PrimaryWorkerNodeId":     12,
  "RetentionDays":           30,
  "Snapshots": [
    {
      "SnapshotDate":        "2026-05-01",
      "SizeBytes":           4194304,
      "KeyEpoch":            8,
      "BuiltAtEpoch":        1746460800,
      "Sha256":              "..."
    }
  ]
}
```

### 7.3 Errors

| Trigger | Code |
|---|---|
| `From` > `To` | `WORKER-300-03 RequestBodyInvalid` |
| `Limit` exceeds hard ceiling | `WORKER-300-03 RequestBodyInvalid` |

---

## 8. BE-5 — `GET /API/V1/Backup/Health`

**Purpose:** Single dashboard call — never throws on degradation; surfaces it in the body.

### 8.1 Response — `200 OK`

```jsonc
{
  "BackupWorkerNodeId":      27,
  "PrimaryWorkerNodeId":     12,
  "Status":                  "BackupAttached" | "BackupLagging" | "Provisioning",
  "LastAcceptedSyncOpSeq":   918374,
  "LastAcceptedAtEpoch":     1746528002,
  "LagSeconds":              42,
  "ActiveKeyEpoch":          8,
  "PendingKeyEpoch":         null,
  "RetiredKeyEpoch":         7,
  "DiskFreeBytes":           107374182400,
  "RetentionDays":           30
}
```

### 8.2 Errors

| Trigger | Code |
|---|---|
| Endpoint cannot read its own `BackupKeyEpoch` table | `WORKER-500-03 SplitDbWriteFail` |

---

## 9. Auth — S2S OAuth recap

Per `05-auth-and-2fa.md` §S2S:

- Client credentials grant (`grant_type=client_credentials`).
- Audience MUST be `Backup` (new value to add in Phase 11 alongside `Main`, `Worker`).
- Scope MUST contain at least one of `Backup.Diff.Write`, `Backup.Rotate.Write`, `Backup.Restore.Write`, `Backup.Read`.
- Tokens expire per `MainWorker.Auth.WorkerJwtTtlSeconds` (default 900 s).
- The reverse proxy in front of every backup node strips any `Cookie` header before handing off — backups never read session cookies.

Endpoint ↔ scope matrix:

| Endpoint | Required scope |
|---|---|
| BE-1 IncrementalDiff | `Backup.Diff.Write` |
| BE-2 RotateKeys | `Backup.Rotate.Write` |
| BE-3 RestoreByDate | `Backup.Restore.Write` |
| BE-4 Snapshots | `Backup.Read` |
| BE-5 Health | `Backup.Read` |

---

## 10. Endpoint timeouts (mirrored verbatim into `15-tunable-constants.md` §2.13)

| Key | Default | Unit | Used by | Notes |
|---|---:|---|---|---|
| `MainWorker.Backup.Endpoint.IncrementalDiffTimeoutSeconds` | **120** | seconds | BE-1 | Primary-side HTTP timeout. Larger than `MainWorker.Routing.HttpTimeoutSeconds=15` because envelope upload may be MB-scale. |
| `MainWorker.Backup.Endpoint.RotateKeysTimeoutSeconds` | **30** | seconds | BE-2 | Per §5.4. |
| `MainWorker.Backup.Endpoint.RestoreByDateTimeoutSeconds` | **60** | seconds | BE-3 | Accepts the job; the actual restore is a long-running worker job. |
| `MainWorker.Backup.Endpoint.SnapshotsTimeoutSeconds` | **15** | seconds | BE-4 | Catalogue read. |
| `MainWorker.Backup.Endpoint.HealthTimeoutSeconds` | **5** | seconds | BE-5 | Dashboard polling — keep tight to avoid masking real outages. |

---

## 11. Errors introduced by this file (mirrored into `13-error-codes.md`)

Reserved sub-range for Phase 9 wire-only failures: **`MAIN-830-*` / 21189-21190** (definitive entries land in this file; rows 21191-21199 reserved for Phase 11 snapshot/restore).

| Prefixed | Flat | HTTP | Name | Meaning |
|---|---:|---:|---|---|
| `MAIN-830-01` | 21189 | 404 | `SnapshotNotFound` | BE-3: requested `SnapshotDate` not present on backup filesystem (Phase 11 owns the storage contract; this row is the wire response). |
| `MAIN-830-02` | 21190 | 409 | `RestoreAlreadyInProgress` | BE-3: a `RestoreJobId` for the same `(BackupWorkerNodeId, TargetPrimaryWorkerNodeId)` is still running. |

§4 (range table) of `13-error-codes.md` updates: previously-pending `MAIN-21189-21199` narrows to `MAIN-21191-21199`.

No new `WORKER-*` codes — every decrypt / watermark / CDC failure on BE-1 is already catalogued in Phases 5–8.

---

## 12. CODE RED budget for handlers

Each handler implementation (out of scope here, future Phase 13+):

| Handler | Max LOC | Nesting | Operands |
|---|---:|---|---|
| BE-1 receive | 15 | zero | ≤2 per guard |
| BE-2 rotate step | 12 | zero | ≤2 |
| BE-3 enqueue restore | 10 | zero | ≤2 |
| BE-4 list | 8 | zero | ≤2 |
| BE-5 health | 10 | zero | ≤2 |

Linter `BACKUP-EP-001` (queued for Phase 12) asserts handler size + `_lib`-style guard extraction.

---

## 13. Cross-references

- `06-core-api-endpoints.md` §1 — header conventions (consumed verbatim).
- `06-core-api-endpoints.md` §2 — receives §2.X table merge (this file is the source of truth until merged).
- `08-error-contract.md` — JSON envelope shape (consumed).
- `13-error-codes.md` §3.10 — new `MAIN-830-*` rows.
- `15-tunable-constants.md` §2.13 — mirrored endpoint timeouts.
- `18-backup-nodes.md` §6 `KnownBackupNode` — pairing assertion source for BE-1.
- `19-incremental-backup-sync.md` §3.1 `BackupSyncWatermark` — destination for BE-1's ACK value.
- `20-backup-encryption-and-keys.md` §4 / §7.2 / §9 — sealed envelope + rotation flow consumed by BE-1 / BE-2.

---

## 14. Open Questions (logged, non-blocking)

- **OQ-21-1** Should BE-1 also accept a streaming `application/octet-stream` body instead of multipart, to avoid double-buffering on multi-GB envelopes? Inferred: defer until envelope sizes consistently exceed 100 MB; multipart is the dumb-AI-friendly default.
- **OQ-21-2** Should BE-5 require any auth scope at all, or be readable by an unauthenticated `/Health/Liveness` proxy probe? Inferred: keep `Backup.Read` scope; proxies use a separate `/Health/Liveness` (out of scope here, owned by `18-…`).

---

*Backup endpoints v1.0.0 — 2026-05-06 (Phase 9). Wire-only — apply logic Phase 10, snapshot/restore Phase 11.*
