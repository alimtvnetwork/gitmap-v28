# 10 — Worker Bootstrap Protocol

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** audit findings F-B-01, F-B-02, F-B-03, F-X-08 (top-10 fix #1). Unblocks AC-1, AC-3, AC-4.
**Authority:** This file is the canonical contract for first-boot worker registration. Any diagram or prose elsewhere defers to this file on conflict.

---

## 1. Purpose

Defines exactly what a freshly-installed Worker process does between `process start` and `ready to serve traffic`. No prior knowledge of Main is assumed beyond a single configuration file. After bootstrap, the Worker can:

- Mint and verify Main-issued JWTs.
- Receive routed tenant requests.
- Accept push-update instructions (per `spec/14-update/28-worker-push-instruction.md`).
- Heartbeat to Main on schedule.

---

## 2. Required pre-boot configuration

Provisioned by the installer (see `spec/14-update/18-install-scripts.md`) into the Worker's **Settings tier** DB:

| Key | Type | Nullable | Source | Notes |
|---|---|---|---|---|
| `MainBaseUrl` | TEXT | No | Installer arg `--main-url` | e.g. `https://main.example.com`. HTTPS only. |
| `WorkerNodeDisplayName` | TEXT | No | Installer arg `--node-name` | Human label, e.g. `worker-kl-01`. |
| `WorkerOAuthClientId` | TEXT | No | Issued by Main Power Admin out-of-band | OAuth client-credentials id. |
| `WorkerOAuthClientSecret` | TEXT | No | Issued by Main, stored encrypted at rest | Rotated per `spec/05-split-db-architecture/`. |
| `WorkerEndpointPublic` | TEXT | No | Installer arg `--public-url` | URL Main will route tenant traffic to. |
| `WorkerVersionPin` | TEXT (SemVer) | No | Compiled in via `-ldflags` | The version Main MUST agree to talk to. |
| `Description` | TEXT | Yes | — | Per Code Red Rule 10. |

If any non-nullable key is missing, Worker MUST exit with `WORKER-000-01 BOOTSTRAP_CONFIG_MISSING` and NOT open any port.

---

## 3. Bootstrap sequence (8 steps, deterministic)

```
Worker boot                              Main
   |                                      |
1. Load Settings DB --------------->      |
2. Self-test split-DB tiers (R/S/A/Sn)    |
3. POST /API/V1/Workers/Register -------->|  (OAuth client_credentials)
   {WorkerNodeDisplayName,                |
    WorkerEndpointPublic,                 |
    WorkerVersionPin,                     |
    BootInstanceUlid}                     |
                                          |--validate version range
                                          |--upsert WorkerNode row
                                          |--mint WorkerNodeId
   <----------------------------------200 |  {WorkerNodeId,
                                          |   JwtPublicKeyPem,
                                          |   JwtSigningKeyId,
                                          |   AcceptedVersionRange,
                                          |   HeartbeatIntervalSeconds,
                                          |   ServerTimeUtc}
4. Verify clock skew ≤ 60 s              |
5. Persist WorkerNodeId + JwtPublicKey to Settings DB
6. Run schema-migrate on App + Session tiers
7. Open public listener on WorkerEndpointPublic
8. Start heartbeat loop (interval from §3 response)
```

### 3.1 `POST /API/V1/Workers/Register` request

```jsonc
{
  "WorkerNodeDisplayName": "worker-kl-01",
  "WorkerEndpointPublic":  "https://w1.example.com",
  "WorkerVersionPin":      "1.5.0",
  "BootInstanceUlid":      "01J...ULID",   // new every process boot; idempotency
  "Description":           null
}
```

Headers: `Authorization: Bearer <oauth-token>`, `X-Correlation-Id`, `X-Idempotency-Key: <BootInstanceUlid>`.

### 3.2 `200 OK` response

```jsonc
{
  "WorkerNodeId":              3,
  "JwtPublicKeyPem":           "-----BEGIN PUBLIC KEY-----\nMIIB...\n-----END PUBLIC KEY-----",
  "JwtSigningKeyId":           "main-jwt-2026-q2",
  "AcceptedVersionRange":      ">=1.4.0 <2.0.0",
  "HeartbeatIntervalSeconds":  30,                                              // default per `15-tunable-constants.md` §2.3 (`MainWorker.Heartbeat.IntervalSeconds`)
  "ServerTimeUtc":             "2026-05-04T12:00:00Z",
  "Description":               null
}
```

Worker MUST refuse to proceed if its `WorkerVersionPin` ∉ `AcceptedVersionRange` → exit `WORKER-400-02 VERSION_MISMATCH`.

---

## 4. JWT public-key fetch (canonical)

Worker fetches the JWT signing public key **only** through the registration response (§3.2 `JwtPublicKeyPem`). There is **no separate `/jwks` endpoint** in v1.0 — this prevents the dumb-AI implementer from inventing one.

Re-fetch triggers (Worker MUST re-register, NOT call any other endpoint):

| Trigger | Action |
|---|---|
| JWT verification fails with `kid` mismatch | Re-register; cache new `JwtPublicKeyPem`. |
| `JwtSigningKeyId` rotated by Main (signaled via heartbeat response) | Re-register on next heartbeat tick. |
| Worker restart | Re-register (idempotent via `BootInstanceUlid`). |

Trust store: only ONE active key at a time in v1.0. (OQ-10-1 reserved for multi-key rollover.)

---

## 5. Version pinning rules

| Compile-time `WorkerVersionPin` | Main `AcceptedVersionRange` | Outcome |
|---|---|---|
| `1.5.0` | `>=1.4.0 <2.0.0` | ✅ Register |
| `1.3.9` | `>=1.4.0 <2.0.0` | ❌ `WORKER_VERSION_TOO_OLD` |
| `2.0.0` | `>=1.4.0 <2.0.0` | ❌ `VERSION_MISMATCH` (major bump requires Main upgrade first) |

`AcceptedVersionRange` MUST be a single SemVer range string (npm-style). No comma-separated lists.

---

## 6. Failure modes

All codes registered in `spec/03-error-manage/03-error-code-registry/` under prefix `WORKER`:

| Code | Name | Phase | Worker action |
|---|---|---|---|
| `WORKER-000-01` | `BOOTSTRAP_CONFIG_MISSING` | Step 1 | Exit, do NOT listen. |
| `WORKER-000-02` | `SPLIT_DB_TIER_MISSING` | Step 2 | Exit. |
| `WORKER-100-01` | `OAUTH_HANDSHAKE_FAIL` | Step 3 | Exit, retry with `MainWorker.Bootstrap.RetryBackoffSeconds` (default `[10, 30, 90, 300]` per `15-tunable-constants.md` §2.9); after `MainWorker.Bootstrap.RetryMaxAttempts` exhaustion, exit. |
| `WORKER-400-01` | `WORKER_VERSION_TOO_OLD` | Step 3 | Exit, log loudly. |
| `WORKER-400-02` | `VERSION_MISMATCH` | Step 3 | Exit, log loudly. |
| `WORKER-500-01` | `CLOCK_SKEW_TOO_LARGE` | Step 4 | Exit. |
| `WORKER-700-01` | `SETTINGS_PERSIST_FAIL` | Step 5 | Exit. |
| `WORKER-700-02` | `SCHEMA_MIGRATE_FAIL` | Step 6 | Exit. |
| `WORKER-800-01` | `LISTENER_BIND_FAIL` | Step 7 | Exit. |

Worker MUST log the failing step number explicitly (per Code Red error-handling rule).

---

## 7. Heartbeat (post-bootstrap)

Once §3 step 8 starts, Worker pings Main every `HeartbeatIntervalSeconds`:

```
POST /API/V1/Workers/{WorkerNodeId}/Heartbeat
  Authorization: Bearer <oauth-token>
  X-Correlation-Id: ...
Body:
{
  "WorkerNodeId":               3,
  "BootInstanceUlid":           "01J...",
  "CurrentVersion":             "1.5.0",
  "ActiveTenantCount":          12,
  "LastUpdateInstructionId":    null,
  "LastUpdateOutcomeCode":      null,
  "Notes":                      null,
  "Comments":                   null
}
```

Response MAY signal:

```jsonc
{
  "JwtSigningKeyIdCurrent":  "main-jwt-2026-q3",   // if changed, worker re-registers
  "RequestReregister":       false,
  "ServerTimeUtc":           "2026-05-04T12:00:30Z"
}
```

Missed-heartbeat policy is owned by Main (per `04-worker-routing.md`). The threshold (`MainWorker.Heartbeat.MissedThreshold`, default **3**), grace window, and quarantine cooldown are pinned in `15-tunable-constants.md` §2.3.

---

## 8. WorkerNode table (Main-side) — see canonical schema

The canonical `WorkerNode` and `WorkerNodeStatus` table definitions live in **`03-main-db-schema.md` §2.1 and §2.2** and are the single source of truth. This file MUST NOT redefine them.

The bootstrap fields named in §3.1 (`WorkerNodeDisplayName`, `WorkerEndpointPublic`, `WorkerVersionPin`, `BootInstanceUlid`) are the **request-body field names** the Worker sends. They map onto the canonical schema columns as follows:

| Bootstrap request field | Canonical column in `03-` §2.1 |
|---|---|
| `WorkerNodeDisplayName` | `WorkerNodeTitle` |
| `WorkerEndpointPublic` | `WorkerNodeEndpoint` |
| `WorkerVersionPin` | `WorkerVersion.WorkerVersionSemver` (latest row, see `03-` §2.7) |
| `BootInstanceUlid` | not persisted on `WorkerNode`; used only for `X-Idempotency-Key` (see §3.1) |

The `WorkerNodeIdentity` column (`03-` §2.1) is server-derived from the OAuth client identity at registration time. Status enum seeded via `spec/06-seedable-config-architecture/`.

---

## 9. Worker-side Settings DB persistence

```sql
CREATE TABLE WorkerBootstrapState (
    WorkerBootstrapStateId  INTEGER PRIMARY KEY AUTOINCREMENT,
    WorkerNodeId            INTEGER NOT NULL,
    JwtSigningKeyId         TEXT NOT NULL,
    JwtPublicKeyPem         TEXT NOT NULL,
    AcceptedVersionRange    TEXT NOT NULL,
    HeartbeatIntervalSec    INTEGER NOT NULL,
    LastRegisteredAtUtc     TEXT NOT NULL,
    Notes                   TEXT NULL,
    Comments                TEXT NULL,
    Description             TEXT NULL
);
```

Only ONE row at a time. Re-registration UPDATEs in place.

---

## 10. Open Questions (logged, non-blocking)

- **OQ-10-1** Multi-active-key JWT rollover (currently single-key). Inferred: defer to v2.0; rotation = brief re-register window.
- **OQ-10-2** Should `BootInstanceUlid` collisions across workers be globally unique or per-worker? Inferred: ULIDs are globally unique by construction; treat globally.

---

## 11. Cross-references

- `spec/19-main-worker-service/04-worker-routing.md` — uses `WorkerNode` rows seeded here.
- `spec/19-main-worker-service/05-auth-and-2fa.md` — JWT verification uses `JwtPublicKeyPem` from §3.2.
- `spec/19-main-worker-service/06-core-api-endpoints.md` — §2.5 `/Workers/Register` and `/Workers/{Id}/Heartbeat`.
- `spec/14-update/28-worker-push-instruction.md` — push-update relies on `WorkerNodeId` minted here.
- `spec/03-error-manage/03-error-code-registry/` — MUST register all `WORKER-*` codes from §6.
- `spec/05-split-db-architecture/` — Settings tier hosts `WorkerBootstrapState`.
- `spec/06-seedable-config-architecture/` — seeds `WorkerNodeStatus` enum + bootstrap config keys.

---

*Worker bootstrap protocol v1.0.0 — 2026-05-04*
