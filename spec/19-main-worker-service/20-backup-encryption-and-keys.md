# 20 — Backup Encryption and Key Rotation

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative (spec-only, per plan §Mode)
**Resolves:** Locked decision **D13** (RSA pair shared between Worker and its Backups; Main issues rotation; zip password follows known pattern). Closes open question **OQ-A3** (zip password derivation).
**Depends on:** [`18-backup-nodes.md`](./18-backup-nodes.md), [`19-incremental-backup-sync.md`](./19-incremental-backup-sync.md), [`05-auth-and-2fa.md`](./05-auth-and-2fa.md), [`13-error-codes.md`](./13-error-codes.md), [`15-tunable-constants.md`](./15-tunable-constants.md).

---

## Keywords

`backup-encryption` · `rsa-key-pair` · `key-rotation` · `zip-password` · `hkdf` · `key-state-machine`

---

## 1. Purpose

Define how a Worker (the **primary**) encrypts each incremental envelope (Phase 7) before shipping it to one of its attached **backup** nodes, how that envelope is sealed inside a password-protected zip, and how the Main node orchestrates **rotation** of the underlying key material without dropping in-flight traffic.

This file is encryption only. Wire endpoints land in Phase 9 (`21-backup-endpoints.md`). Apply logic lands in Phase 10. Snapshot/restore lands in Phase 11.

---

## 2. Threat model (in-scope / out-of-scope)

In-scope:

- Confidentiality of envelope contents in transit and at rest on the backup node.
- Authenticity — a backup MUST reject any envelope not signed by the paired primary.
- Forward secrecy across rotation — a leaked retired key MUST NOT decrypt envelopes minted under the new active key.
- Defence against operator-side mistakes (wrong primary uploading to wrong backup).

Out-of-scope (explicitly):

- Network-level transport security — assumed handled by mTLS at the reverse-proxy layer (`05-auth-and-2fa.md` §S2S).
- Hardware tamper resistance on the backup node — operational concern.
- Long-term archival cryptography (post-quantum) — review scheduled v3.0.

---

## 3. Key material inventory

Three distinct artefacts. Do **not** conflate.

| # | Artefact | Lifetime | Storage location | Holder |
|---|---|---|---|---|
| K1 | **Pair-RSA** — RSA-4096 key pair, one pair per `(PrimaryWorkerNodeId, BackupWorkerNodeId)` couple. | One **KeyEpoch**. Rotated on Main instruction. | Private half on primary + backup; public half mirrored to Main. | Primary, Backup, Main (public only). |
| K2 | **Envelope-AES** — AES-256-GCM session key, freshly generated per envelope. | One envelope. Discarded after sealing. | In-memory only. | Primary (write), Backup (read). |
| K3 | **Zip-Password** — derived ASCII password used by the outer ZIP container's AES-256 mode. | Per envelope. Re-derived on demand by the holder of the shared secret. | Never persisted. | Primary (write), Backup (read). |

**Invariant E1.** A backup node MUST hold exactly one **Active** Pair-RSA private key per primary at any moment (plus optionally one **Pending** and one **Retired** during rotation; see §6).

**Invariant E2.** Envelope-AES keys (K2) are never reused. A primary that crashes mid-envelope MUST mint a fresh K2 on retry, even when re-sending the same `EnvelopeId`.

**Invariant E3.** Zip-Password (K3) is a deterministic function of the Pair-RSA epoch identifier and the envelope timestamp — it MUST be reproducible by the receiver from envelope metadata alone.

---

## 4. Envelope sealing pipeline (primary side)

Per envelope produced by Phase 7 (`19-incremental-backup-sync.md` §6):

```
1. PHASE7_OUTPUT  : SQLite file `envelope-<EnvelopeId>.sqlite`
2. K2 = RandomBytes(32)                              // AES-256 key
3. IV = RandomBytes(12)                              // GCM nonce
4. CIPHERTEXT, GcmTag = AES-256-GCM(K2, IV, PHASE7_OUTPUT, AAD)
   AAD = concat(EnvelopeId || KeyEpoch || PrimaryWorkerNodeId || BackupWorkerNodeId)
5. WRAPPED_K2 = RSA-OAEP-SHA256(BackupPublicKey_K1, K2)
6. SIGNATURE  = RSA-PSS-SHA256 (PrimaryPrivateKey_K1, SHA256(CIPHERTEXT || AAD))
7. INNER_BLOB = JSON {
        "EnvelopeId":            …,
        "KeyEpoch":              …,
        "PrimaryWorkerNodeId":   …,
        "BackupWorkerNodeId":    …,
        "Iv":                    Base64(IV),
        "GcmTag":                Base64(GcmTag),
        "WrappedKey":            Base64(WRAPPED_K2),
        "Signature":             Base64(SIGNATURE),
        "Ciphertext":            Base64(CIPHERTEXT)
   }
8. ZIP_PASSWORD = DerivePassword(KeyEpoch, EnvelopeTimestampEpoch)   // §5
9. OUTPUT       = AES-256-ZIP(INNER_BLOB, ZIP_PASSWORD)              // standard zipcrypto NOT permitted
```

Steps 4–6 use widely-available primitives in every default-stack target (`openssl_*` in PHP, `crypto.subtle` in TS, `crypto/rsa` in Go, `RSACryptoServiceProvider` in C#, `ring` / `RustCrypto` in Rust). The outer zip is **only** there to give operators a familiar artefact they can move around with normal tooling — its password is not a security boundary on its own; the inner GCM + RSA sign/wrap layer is the real crypto.

---

## 5. Zip-Password derivation (resolves OQ-A3)

```
SharedSecret    = HKDF-SHA256(
                      ikm     = PrimaryPrivateKey_K1.Modulus  XOR  BackupPublicKey_K1.Modulus,
                      salt    = "BackupZip/v1",
                      info    = concat(PrimaryWorkerNodeId || BackupWorkerNodeId || KeyEpoch)
                  ).Take(32 bytes)

DerivePassword(KeyEpoch, EnvelopeTimestampEpoch) =
    Hex( HMAC-SHA256( SharedSecret, BigEndian64(EnvelopeTimestampEpoch) ) ).Substring(0, 32)
```

Why this shape:

- **Deterministic** — receiver reproduces password from `KeyEpoch` + `EnvelopeTimestampEpoch` (both in clear envelope metadata).
- **Per-envelope** — different timestamp ⇒ different password ⇒ leak of one zip doesn't open siblings.
- **No external KDF service** — every default-stack ships HKDF + HMAC-SHA256.
- **Fixed length 32 hex chars** — copy-pasteable in `unzip -P …` for ops, no shell-escaping landmines.
- **Rotation-bound** — `KeyEpoch` change invalidates every previously derived password automatically.

**Forbidden alternatives:**

- Plain `SharedSecret` as password (no per-envelope variance).
- ZipCrypto (the legacy weak cipher) — receivers MUST refuse it (`WORKER-920-04`).
- Passwords < 32 chars or non-hex (linter `BACKUP-ZIP-001`).

---

## 6. Pair-RSA key state machine

Per `(PrimaryWorkerNodeId, BackupWorkerNodeId)` couple, exactly one row per state in the new `BackupKeyEpoch` table on **both** primary and backup (mirrored).

```
                 RotationInstruction
   Active  ──────────────────────────►  Pending(new) + Retired(old)
                                                    │
                                AckBoth + GraceWindow elapses
                                                    │
                                                    ▼
                                              Active(new) + Retired(old)
                                                    │
                                          KeepRetiredForLagWindow
                                                    │
                                                    ▼
                                                 Discarded
```

States:

| State | Meaning | Allowed envelope direction |
|---|---|---|
| `Pending` | New pair generated; both sides hold private halves; not yet sealing envelopes with it. | Decrypt only (in case of clock skew on Main). |
| `Active` | The single epoch new envelopes are sealed under. | Seal + decrypt. |
| `Retired` | Previous epoch; receiver still accepts envelopes (in-flight when rotation fired). | Decrypt only. |
| `Discarded` | Past `MainWorker.Backup.RetiredKeyGraceSeconds`. Private material wiped. | Reject. |

Concurrency rule: **at most one Pending + at most one Active + at most one Retired** at any time. A second rotation attempted while `Pending` already exists returns `MAIN-820-02 RotationAlreadyInProgress`.

---

## 7. Rotation orchestration (Main → Worker → Backup)

CODE RED budget for the orchestrator: ≤ 15 lines per step. State transitions are persisted before any side-effect.

### 7.1 Trigger sources

1. **Scheduled** — Main cron evaluates `MainWorker.Backup.MaxKeyAgeSeconds` per couple. Default 90 days.
2. **Operator-forced** — admin endpoint (Phase 9 will expose `POST /API/V1/Backup/RotateKeys`).
3. **Compromise-forced** — explicit `Reason="Compromise"` skips the grace window (sets `RetiredKeyGraceSeconds` override = 0 for that rotation only).

### 7.2 Step-by-step

```
S1  Main allocates RotationId (uuid) + new KeyEpoch (monotonic int).
S2  Main commands the primary: "GenerateNewPairFor(BackupId, KeyEpoch)".
    Primary mints pair, stores private half in `BackupKeyEpoch (state=Pending)`,
    returns public half + attestation signature.
S3  Main forwards the same KeyEpoch + primary-attested public half to the backup
    via a separate Main→Backup call.
    Backup mints its own pair, stores private half in `BackupKeyEpoch (state=Pending)`,
    returns its public half + attestation.
S4  Main cross-publishes the two public halves so each side knows the other's K1 public.
S5  When BOTH sides ACK Pending, Main flips local row to "ActivationArmed".
S6  After `MainWorker.Backup.RotationActivationDelaySeconds` (default 60s)
    Main issues "Activate(KeyEpoch)" to BOTH sides in parallel.
    On receipt:
        old Active → Retired
        Pending    → Active
S7  Both sides ACK Activate. Main marks RotationId=Completed.
S8  Background sweeper on each side discards Retired epochs older than
    `MainWorker.Backup.RetiredKeyGraceSeconds` (default 86400s / 24h).
```

Failure handling:

- Any step times out → `MAIN-820-01 RotationStepTimeout` + automatic rollback (delete Pending, leave Active untouched).
- ACK mismatch (primary acked Pending, backup did not within `RotationAckTimeoutSeconds`) → rotation aborts.
- Activation step is the only step that is **not** rolled back automatically — once one side has activated, Main MUST drive the other side to activation manually (operator alert `MAIN-820-03 RotationActivationSplitBrain`).

---

## 8. New table: `BackupKeyEpoch` (App tier, both primary and backup)

```
CREATE TABLE BackupKeyEpoch (
    BackupKeyEpochId          INTEGER PRIMARY KEY AUTOINCREMENT,
    PrimaryWorkerNodeId       INTEGER NOT NULL,
    BackupWorkerNodeId        INTEGER NOT NULL,
    KeyEpoch                  INTEGER NOT NULL,            -- monotonic, allocated by Main
    State                     TEXT    NOT NULL,            -- Pending|Active|Retired|Discarded
    PrivateKeyPem             TEXT    NULL,                -- NULL once Discarded
    PeerPublicKeyPem          TEXT    NULL,                -- counterpart's public half
    CreatedAt                 INTEGER NOT NULL,            -- epoch seconds (D2)
    ActivatedAt               INTEGER NULL,
    RetiredAt                 INTEGER NULL,
    DiscardedAt               INTEGER NULL,
    Description               TEXT    NULL,                -- Rule 12: nullable, no DEFAULT
    UNIQUE (PrimaryWorkerNodeId, BackupWorkerNodeId, KeyEpoch)
);
```

Notes:

- PascalCase + `{TableName}Id` PK (Memory: Naming, DB Schema).
- `Description` nullable (Memory: DB Schema Rule 10).
- DateTime columns are INTEGER epoch seconds (D2).
- `PrivateKeyPem` MUST be wiped (`UPDATE … SET PrivateKeyPem=NULL`) when `State` flips to `Discarded`. Linter `BACKUP-KEY-001` to enforce.

The matching `BackupKeyEpoch` row on the **Main** side (mirror) does NOT carry `PrivateKeyPem` — Main holds public halves only.

---

## 9. Verification on the backup (decrypt path)

```
1. Receive zip + envelope metadata (timestamp, KeyEpoch).
2. ZIP_PASSWORD = DerivePassword(KeyEpoch, EnvelopeTimestamp).
3. UNZIP — refuse if cipher != AES-256 (WORKER-920-04 ZipCipherTooWeak).
4. Look up BackupKeyEpoch row by KeyEpoch.
   - State=Active or Retired → continue.
   - State=Pending          → reject WORKER-920-01 KeyEpochNotYetActive.
   - State=Discarded        → reject WORKER-920-02 KeyEpochDiscarded.
   - Not found              → reject WORKER-920-03 UnknownKeyEpoch.
5. RSA-PSS verify Signature against (Ciphertext || AAD) using PeerPublicKeyPem.
6. RSA-OAEP unwrap WrappedKey to recover K2.
7. AES-256-GCM decrypt → SQLite envelope file.
8. Hand to Phase-10 apply pipeline (out of scope here).
```

Failure on any step short-circuits with the named error code; partial state MUST NOT be persisted.

---

## 10. Tunables introduced (mirrored verbatim into `15-tunable-constants.md` §2.12)

| Key | Default | Unit | Notes |
|---|---:|---|---|
| `MainWorker.Backup.MaxKeyAgeSeconds` | **7776000** | seconds (90d) | Triggers scheduled rotation. |
| `MainWorker.Backup.RotationAckTimeoutSeconds` | **120** | seconds | S2/S3/S6 ACK budget. |
| `MainWorker.Backup.RotationActivationDelaySeconds` | **60** | seconds | S5→S6 settle window. |
| `MainWorker.Backup.RetiredKeyGraceSeconds` | **86400** | seconds (24h) | How long Retired private material is kept for in-flight envelopes. |
| `MainWorker.Backup.RsaKeySizeBits` | **4096** | bits | RSA-OAEP / RSA-PSS modulus. |

Override via Seedable-Config; same linter rules (T1/T2/T3) apply.

---

## 11. Errors introduced (mirrored verbatim into `13-error-codes.md`)

Worker tier (decrypt failures, range `WORKER-920-*` / 21095-21099):

| Prefixed | Flat | HTTP | Meaning |
|---|---:|---:|---|
| `WORKER-920-01` | 21095 | 409 | `KeyEpochNotYetActive` — envelope received under a Pending epoch. |
| `WORKER-920-02` | 21096 | 410 | `KeyEpochDiscarded` — Retired grace window expired. |
| `WORKER-920-03` | 21097 | 404 | `UnknownKeyEpoch` — `BackupKeyEpoch` row missing. |
| `WORKER-920-04` | 21098 | 415 | `ZipCipherTooWeak` — outer zip used ZipCrypto / non-AES. |
| `WORKER-920-05` | 21099 | 422 | `EnvelopeSignatureInvalid` — RSA-PSS verify failed. |

Main tier (rotation orchestration, range `MAIN-820-*` / 21186-21188):

| Prefixed | Flat | HTTP | Meaning |
|---|---:|---:|---|
| `MAIN-820-01` | 21186 | 504 | `RotationStepTimeout` — S2/S3/S6 ACK missed. |
| `MAIN-820-02` | 21187 | 409 | `RotationAlreadyInProgress` — Pending epoch already exists. |
| `MAIN-820-03` | 21188 | 500 | `RotationActivationSplitBrain` — one side activated, other did not. |

§4 (range table) of `13-error-codes.md` updates: future-expansion `WORKER-21095-21099` collapses to **fully allocated**; `MAIN-21186-21199` narrows to `MAIN-21189-21199`.

---

## 12. Linter hooks queued for Phase 12

| ID | Rule |
|---|---|
| `BACKUP-KEY-001` | Every `BackupKeyEpoch` row with `State='Discarded'` MUST have `PrivateKeyPem IS NULL` and `DiscardedAt IS NOT NULL`. |
| `BACKUP-KEY-002` | Per `(PrimaryWorkerNodeId, BackupWorkerNodeId)`: at most one row each in `Pending` / `Active`. |
| `BACKUP-ZIP-001` | Any zip artefact under `var/backup-outbox/` MUST be AES-256 (`unzip -lv` parse). |

---

## 13. Cross-references

- Decision register: `.lovable/plan.md` D13.
- `18-backup-nodes.md` §3 — primary/backup pairing.
- `19-incremental-backup-sync.md` §6 — envelope SQLite shape (input to step 1 of §4).
- `05-auth-and-2fa.md` §S2S — transport mTLS context.
- `13-error-codes.md` §4 — reserved-range table (updated).
- `15-tunable-constants.md` §2.12 — mirrored tunables.

---

## 14. Open Questions (logged, non-blocking)

- **OQ-20-1** Should `MAIN-820-03 RotationActivationSplitBrain` page on-call automatically? Inferred: yes, severity P1, but routing is an ops concern out of spec.
- **OQ-20-2** Is RSA-4096 a defensible default vs. moving to Ed25519 (signature) + X25519 (key wrap)? Deferred to v2.0 design review — interop with PHP 8.1 default OpenSSL build is the deciding factor.

---

*Backup encryption + key rotation v1.0.0 — 2026-05-06 (Phase 8). Resolves D13 + OQ-A3.*
