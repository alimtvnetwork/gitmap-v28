# 27 — Git Backup Targets (GitHub / Git-SSH Repos)

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-08
**Status:** Authoritative (Phase 14)
**Scope:** Defines an **additional, optional backup destination class**: a Worker may push its sealed App-tier backup stream into one or more **Git repositories** (typically GitHub, but any git-over-SSH host) when an SSH deploy key with write permission has been provisioned for that worker.

> **Complement, not replacement.** Git backup targets are *additional* to the primary→backup-node CDC channel (ch. 18/19). They never replace it. Disaster recovery still flows through `BackupNode` first; git is a cold, off-site, history-preserving secondary.

---

## 1. Why git as a backup target

| Property | Backup Node (ch. 18) | Git Backup Target (this chapter) |
|---|---|---|
| Latency | Seconds (CDC) | Minutes–hours (batched commit) |
| Failover candidate | YES (Power-Admin promotion) | NO (cold archive only) |
| History | Last K envelopes + snapshots | Full commit history (effectively unbounded) |
| Off-site | Optional (operator choice) | YES (always external host) |
| Cost | Self-hosted infra | Free/cheap on GitHub free tier |
| Tamper-evidence | KeyEpoch HMAC | Git commit hash chain + signed tags |

Use cases:

- Off-site cold archive for compliance ("we keep last 7 years of nightly backups").
- Multi-region redundancy without provisioning a second `BackupNode`.
- Tamper-evident chain: each commit's SHA covers the entire prior history.
- Cheap "bring-your-own-storage" tier for small-tenant deployments.

---

## 2. Permission model (one-way, narrow)

A Git Backup Target is configured **per primary Worker** by a Power-Admin operator. The configuration consists of:

| Field | Type | Notes |
|---|---|---|
| `GitBackupTargetId` | INTEGER PK | `{TableName}Id` rule. |
| `WorkerNodeId` | INTEGER FK → `WorkerNode` | MUST be `IsBackup=0` (primary only). |
| `RepoSshUrl` | TEXT | e.g. `git@github.com:acme-backups/worker-eu-1.git` |
| `SshKeyRef` | TEXT | Reference to a secret in the Worker's local secret store; private key NEVER leaves the Worker host. |
| `RepoNamespacePrefix` | TEXT | e.g. `acme-backups/`. Worker MAY auto-create repos under this prefix only (§4). |
| `IsAutoCreateAllowed` | INTEGER (0/1) | If 1 and the SSH key has org-level repo-create scope, Worker may create repos under `RepoNamespacePrefix`. |
| `Description` | TEXT NULL | Operator notes (memory rule 11). |

**Trust gradient (anchored in ch. 26 §2):**

- The SSH private key lives **only on the primary Worker host**, in its local secret store. Main never holds it. Backup Nodes never hold it.
- The deploy key's scope MUST be limited to the namespace `RepoNamespacePrefix`. Org-wide admin keys are explicitly forbidden.
- The git remote authorisation MUST be **push-only**. Read access from third parties is allowed (it is just an encrypted blob store), but no other system may push.
- **A compromised Worker can corrupt only its own git namespace.** It cannot reach another worker's repos because each worker holds a different SSH key with a different namespace scope.

---

## 3. Payload format (encrypted, append-only)

Each commit on a git backup target represents one "git backup epoch" — typically one per snapshot interval (default: every 6 hours, tunable via `GitBackup.IntervalSeconds` in `15-tunable-constants.md` — to be added in the wire-up task, not this chapter).

Repo layout:

```
acme-backups/worker-eu-1/
├── readme.md                          # Auto-generated, see §3.1
├── manifest.json                      # Latest GitBackupManifest (see §3.2)
├── envelopes/
│   ├── 0000000001.env.bin             # Sealed BackupOutboxEnvelope (ch. 19)
│   ├── 0000000002.env.bin
│   └── ...
└── snapshots/
    ├── 2026-05-08T00-00-00Z.snap.bin  # Sealed BackupSnapshot (ch. 23)
    └── 2026-05-08T06-00-00Z.snap.bin
```

### 3.1 README contract

Auto-generated; MUST contain only:

- The string `# Worker Backup Archive — DO NOT EDIT` (literal H1).
- The `WorkerNodeId`, `RepoNamespacePrefix`, current Active `KeyEpochId`.
- A pointer to this spec chapter for decoding.

The Worker MUST refuse to push if the README has been hand-edited (SHA mismatch on the prior commit's README blob). This is a tamper tripwire.

### 3.2 `manifest.json`

```json
{
  "WorkerNodeId": 12,
  "GitBackupTargetId": 4,
  "LatestEnvelopeSeq": 9421,
  "LatestSnapshotUtc": "2026-05-08T06:00:00Z",
  "ActiveKeyEpochId": 17,
  "PreviousKeyEpochId": 16,
  "GitBackupSchemaVersion": "1.0.0"
}
```

PascalCase per the project naming rule (memory: Naming).

### 3.3 Encryption

- Envelopes and snapshots are written **already-sealed** under the Active KeyEpoch (ch. 20). The git host sees ciphertext only.
- The KeyEpoch material is **never** committed to git. Recovery requires both the git repo *and* the operator-held KeyEpoch keystore (separation of duties).
- A leaked git repo on its own MUST NOT enable plaintext recovery.

---

## 4. Auto-creation flow (when `IsAutoCreateAllowed = 1`)

If the SSH key has GitHub org-level `repo:create` scope under `RepoNamespacePrefix`, the Worker MAY create repos on demand. The flow:

1. Worker computes target repo name: `{RepoNamespacePrefix}{WorkerNodeId}-{ShardKey}` (deterministic; idempotent).
2. Worker calls GitHub API `POST /orgs/{org}/repos` (or generic git-host equivalent) with the SSH-deploy-key's bearer token. **Private** visibility is mandatory; public creation MUST fail closed.
3. On success, Worker inserts a `GitBackupRepo` row (schema in §5) and proceeds to first push.
4. On 4xx/5xx, Worker emits `WORKER-870-01 GitBackupRepoCreateFailed` (new code, to be allocated in `13-error-codes.md` during the wire-up task).

**Auto-creation is disabled by default.** Operators MUST opt-in per target.

---

## 5. New main-DB tables (sketch — final binding lives in `03-main-db-schema.md` next pass)

| Table | Purpose | Rule-12 column |
|---|---|---|
| `GitBackupTarget` | Per-worker git-push configuration (§2). | `Description TEXT NULL` |
| `GitBackupRepo` | Tracks each auto-created or operator-registered repo under a target. | `Description TEXT NULL` |
| `GitBackupPushLog` | Append-only push history (transactional). | `Notes TEXT NULL`, `Comments TEXT NULL` |

All PKs are `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT` (memory: DB Schema). No UUIDs.

> **SPEC-ONLY NOTE.** This chapter is binding spec text only. The schema rows above MUST be added to `03-main-db-schema.md` in a follow-up task before any implementation begins, and a new error code block MUST be allocated in `13-error-codes.md`. Per the project memory's Spec/19 rule, no implementation code is to be written.

---

## 6. New endpoints (sketch — final binding lives in `06-core-api-endpoints.md` next pass)

| Method + Path | Caller | Purpose |
|---|---|---|
| `POST /API/V1/Admin/GitBackupTarget` | Power-Admin | Register a new target for a worker. |
| `GET  /API/V1/Admin/GitBackupTarget/{Id}` | Power-Admin | Inspect target + last push status. |
| `POST /API/V1/Admin/GitBackupTarget/{Id}/Test` | Power-Admin | Dry-run SSH handshake + (optional) repo-create probe. |
| `POST /API/V1/Worker/GitBackup/Trigger` | Power-Admin → routed to Worker | Force an immediate push outside the schedule. |

No backup-node-side endpoints. Backups do NOT push to git (ch. 26 §2 matrix).

---

## 7. Failure semantics

| Condition | Code | Behavior |
|---|---|---|
| SSH key rejected | `WORKER-870-02 GitBackupAuthFailed` | Retry with exponential backoff up to `GitBackup.MaxRetries`; then alert. |
| Repo missing and auto-create disabled | `WORKER-870-03 GitBackupRepoMissing` | Halt pushes for this target; alert operator. |
| Push rejected (non-fast-forward) | `WORKER-870-04 GitBackupPushRejected` | Indicates external tampering — DO NOT force-push. Alert + halt. |
| README tamper detected (§3.1) | `WORKER-870-05 GitBackupReadmeTampered` | Halt + alert + audit event. |
| Network unreachable | `WORKER-870-06 GitBackupNetworkUnreachable` | Backoff + retry. |

All codes will be formally allocated in `13-error-codes.md` during the wire-up task; this chapter reserves them.

---

## 8. Trust-boundary cross-check (ch. 26)

| Question | Answer |
|---|---|
| Can a compromised Worker corrupt other workers' git backups? | **NO** — each worker holds a different SSH key scoped to its own namespace prefix. |
| Can a compromised Backup Node push to git? | **NO** — backups never hold the SSH key (ch. 26 §3.2). |
| Can a compromised Git host force-push back into the Worker? | **NO** — git is a target, not a caller. Worker only initiates outbound git-push. |
| Can a leaked git repo enable plaintext recovery? | **NO** — KeyEpoch material is held separately by the operator (§3.3). |

These four MUST hold for any future change to this chapter.

---

## 9. Open questions deferred to v2.0

| OQ | Question | Why deferred |
|---|---|---|
| **OQ-27-1** | Should snapshot commits be **GPG-signed** with a Power-Admin key for stronger tamper evidence? | Requires KMS chapter (not yet started). |
| **OQ-27-2** | Should we support `git-lfs` for snapshots > 100 MiB to stay under GitHub's per-blob limit? | Depends on snapshot pyramid decision (OQ-23-1). |
| **OQ-27-3** | Should the Worker rotate its SSH key on every `KeyEpoch` rotation? | Requires SSH-key-rotation tooling not yet specified. |

---

## 10. Cross-references

- `18-backup-nodes.md` — primary backup destination class.
- `19-incremental-backup-sync.md` — envelope source format (reused for git-pushed envelopes).
- `20-backup-encryption-and-keys.md` — KeyEpoch material that seals git-bound payloads.
- `23-snapshot-storage-and-restore.md` — snapshot format (reused for git-pushed snapshots).
- `26-trust-boundaries-and-isolation.md` — authoritative trust matrix that this chapter conforms to.

---

*Phase 14 (2026-05-08) — adds Git/SSH repos as an optional cold backup destination, fully governed by the ch. 26 trust matrix.*
