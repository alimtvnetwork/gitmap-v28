# 16 — Update Channels (Pull, Push, Reconcile)

**Spec:** `19-main-worker-service`
**Version:** 1.0.0
**Created:** 2026-05-06
**Status:** Authoritative
**Authority:** This file is the canonical contract for **how** a Worker receives a new version. It complements (does NOT replace) `09-self-update-pointer.md` (intent/pause note) and `spec/14-update/28-worker-push-instruction.md` (JSON Instruction Document schema).

---

## 1. Mental Model — Kubernetes Reconciliation

This spec deliberately mirrors the Kubernetes control-plane / kubelet relationship:

- **Main = control plane.** Holds *desired state* per Worker (which version each Worker should run).
- **Worker = kubelet.** Reports *actual state* via heartbeat; reconciles toward desired state on a loop.
- **Channels are additive, not exclusive.** A correctly built Worker MUST support all three; operators pick which to enable per environment via Seedable-Config.

| Channel | Initiator | Direction | Allowed in `Env=Production`? | Default cadence |
|---|---|---|---|---|
| **§2 — Push (dev-only)** | Power Admin (Main) | Main → Worker | ❌ NO (HTTP 403) | On-demand |
| **§3 — Pull from Main (reconciliation loop)** | Worker | Worker → Main | ✅ YES | Every `WorkerUpdate.PullFromMain.PollIntervalSeconds` (default 300 s) |
| **§4 — Pull from known release URL** | Worker | Worker → external HTTPS | ✅ YES | Every `WorkerUpdate.PullFromUrl.PollIntervalSeconds` (default 3600 s) |

All three channels converge on the same execution path: a **JSON Instruction Document** per `spec/14-update/28-worker-push-instruction.md` is constructed (or fetched), validated, and applied via `RenameFirst` deploy.

---

## 2. Channel A — Push from Main (DEV/STAGING ONLY)

### 2.1 Purpose

Allow a Power Admin to upload a zip and force every Worker (or one Worker) to run it immediately. Useful for in-house testing, hotfix rehearsal, and reproducing customer environments.

### 2.2 Hard environment gate

```text
IF Settings.Env IN {Production, Prod}
    REFUSE all push endpoints with HTTP 403 + ErrorCode WORKER-403-01 PUSH_DISABLED_IN_PRODUCTION.
ELSE IF Settings.Env IN {Development, Debug, Staging, QA}
    ALLOW.
ELSE
    REFUSE (default deny).
```

Implementer obligations:
1. The check MUST live in middleware, not the handler.
2. The check MUST read `Settings.Env` from Seedable-Config (key `MainWorker.Environment`); never from a request-scoped value.
3. The check MUST be covered by a regression test that asserts 403 when `Env=Production`.
4. Audit row written to `AccessDenialEvent` (per `08-error-contract.md`) on every refusal.

### 2.3 Endpoints (already declared in `06-core-api-endpoints.md` §2.5)

| Endpoint | Purpose |
|---|---|
| `POST /API/V1/Workers/PublishZip` | Power Admin uploads zip → Main stores under `/Releases/worker/<TargetVersion>/`. |
| `POST /API/V1/Workers/{WorkerNodeId}/Update` | Push to one Worker. |
| `POST /API/V1/Workers/All/Update` | Fan-out to every `Active` Worker. |

### 2.4 Wire flow

1. Power Admin `POST /Workers/PublishZip` (multipart, `application/zip`).
2. Main verifies `Env ≠ Production`, stores zip, computes SHA256, signs with `MainReleaseSigningKey`.
3. Main constructs a JID per `spec/14-update/28-worker-push-instruction.md` §3 with:
   - `PayloadUrl` = the stored zip URL (HTTPS, OAuth-protected).
   - `IssuedByUserId` = the Power Admin.
   - `ExecutionWindow.EarliestStartUtc` = `Now()` (immediate).
4. Main `POST /API/V1/SelfUpdate` to each target Worker.
5. Workers respond `202 Accepted`, then execute per the JID contract.
6. Main returns `207 Multi-Status` aggregating per-worker outcomes.

### 2.5 Why this is dev-only

- A push channel bypasses the Worker's pull-loop safety net (windowed retries, channel selection).
- In production, all updates MUST flow through Channel B or C so they are auditable, schedulable, and rollbackable via the same path operators rehearse daily.

---

## 3. Channel B — Pull from Main (Reconciliation Loop, ALWAYS ON)

### 3.1 Purpose

Each Worker independently asks Main "what version should I be running?" and self-heals if behind. Mirrors `kubelet` polling `kube-apiserver` for pod spec.

### 3.2 Endpoint

`GET /API/V1/SelfUpdate/Desired?WorkerNodeId={id}&CurrentVersion={semver}`

- **Auth:** OAuth client-credentials (Worker → Main) per `05-auth-and-2fa.md` §2.3.
- **Headers:** `X-Correlation-Id` (mandatory), `X-Worker-Identity` (mirror of `WorkerNodeIdentity`).
- **Response 200 (no update needed):**
  ```json
  { "DesiredVersion": "1.5.0", "CurrentVersion": "1.5.0", "InstructionRequired": false }
  ```
- **Response 200 (update needed):** Full JID per `spec/14-update/28-worker-push-instruction.md` §3, with `PayloadUrl` pointing at Main's release store.
- **Response 304 Not Modified:** when `If-None-Match: <ETag>` matches; saves bandwidth.

### 3.3 Worker poll loop

```text
LOOP every WorkerUpdate.PullFromMain.PollIntervalSeconds (default 300):
    response = GET /API/V1/SelfUpdate/Desired
    IF response.InstructionRequired == true:
        validate JID
        execute RenameFirst deploy per spec/14/28 §5
        report outcome via POST /API/V1/Workers/{id}/Heartbeat
    ELSE:
        no-op
END LOOP
```

### 3.4 Desired-state source on Main

Main computes `DesiredVersion` per Worker from this priority list:

1. Explicit pin in `WorkerNode.PinnedVersion` column (set by Power Admin).
2. Channel-level default in `Settings.WorkerUpdate.DefaultDesiredVersion` (Seedable-Config).
3. Latest release tagged `Stable` in Main's release index.

The first non-null value wins. This is the **declarative spec** for that Worker — same idea as a Kubernetes `Deployment.spec.template.spec.containers[].image`.

### 3.5 Tunable constants (added to `15-tunable-constants.md` §2.10)

| Key | Default | Min | Max |
|---|---|---|---|
| `WorkerUpdate.PullFromMain.Enabled` | `true` | — | — |
| `WorkerUpdate.PullFromMain.PollIntervalSeconds` | 300 | 60 | 86400 |
| `WorkerUpdate.PullFromMain.JitterSeconds` | 30 | 0 | 300 |
| `WorkerUpdate.PullFromMain.FailureBackoffSeconds` | `[60, 300, 900]` | — | — |

Jitter prevents thundering-herd when many Workers boot together.

---

## 4. Channel C — Pull from Known Release URL (Main-Independent)

### 4.1 Purpose

Allow a Worker to keep updating itself even when Main is unreachable. Source of truth is a single HTTPS URL pointing to a release manifest (analogous to a Kubernetes container image registry).

### 4.2 Configuration

| Seedable-Config key | Type | Default | Notes |
|---|---|---|---|
| `WorkerUpdate.PullFromUrl.Enabled` | bool | `false` | Disabled by default; opt-in per deployment. |
| `WorkerUpdate.PullFromUrl.ManifestUrl` | string | null | HTTPS only. e.g. `https://releases.example.com/worker/latest.json`. |
| `WorkerUpdate.PullFromUrl.SigningKeyId` | string | null | Must resolve to a public key in worker's trust store. |
| `WorkerUpdate.PullFromUrl.PollIntervalSeconds` | int | 3600 | Min 300, max 86400. |
| `WorkerUpdate.PullFromUrl.AllowedHostsAllowlist` | string[] | `[]` | Worker refuses any `PayloadUrl` host not on this list. |

### 4.3 Manifest format (`latest.json`)

```jsonc
{
  "ManifestVersion": "1.0.0",
  "LatestVersion": "1.5.0",
  "MinimumWorkerVersion": "1.4.0",
  "PayloadUrl": "https://releases.example.com/worker/1.5.0/worker-1.5.0.zip",
  "PayloadSha256": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855",
  "PayloadSizeBytes": 18234567,
  "PayloadSignatureBase64": "MEUCIQDx...==",
  "SigningKeyId": "main-release-2026-q2",
  "PublishedAtUtc": "2026-05-04T12:00:00Z"
}
```

### 4.4 Worker behavior

1. Fetch `ManifestUrl` (HTTPS, optional bearer per Seedable-Config).
2. If `LatestVersion > CurrentVersion` AND `CurrentVersion ≥ MinimumWorkerVersion`:
   - Synthesize a JID locally (Worker is both issuer and executor for this channel).
   - Set `IssuedByUserId = 0` (system), `IssuedByUserDisplayName = "PullFromUrl"`.
   - Validate signature against `SigningKeyId` from Worker's trust store.
   - Execute the same `RenameFirst` deploy path as Channels A and B.
3. Report outcome to Main via `Heartbeat` (best-effort; failure to reach Main MUST NOT abort the local update).

### 4.5 Why allowlist is mandatory

`PayloadUrl` arriving from a remote manifest is an **attacker-controllable string** if the manifest host is compromised. The allowlist is the second line of defense (signature verification is the first). Worker MUST refuse `PayloadUrl` whose host is not on `AllowedHostsAllowlist` with `WORKER-403-02 PAYLOAD_HOST_NOT_ALLOWED`.

---

## 5. Channel Interaction Rules

### 5.1 Precedence

A Worker MAY have multiple channels enabled. When more than one fires in the same window, precedence is:

1. **Channel A (Push)** — wins instantly (operator intent overrides reconciliation).
2. **Channel B (Pull from Main)** — wins over Channel C if both want different versions.
3. **Channel C (Pull from URL)** — used only when (B) returned no instruction OR when `WorkerUpdate.PullFromMain.Enabled=false`.

### 5.2 Idempotency across channels

All three channels produce a JID with a unique `InstructionId` (ULID). Worker dedupes on `InstructionId` per `spec/14-update/28-worker-push-instruction.md` §3.1; the same logical update arriving via two channels is applied **once**.

### 5.3 No double-apply guard

Worker MUST persist a `LastAppliedInstructionId` in the Settings tier. On startup, Worker:
1. Reads `LastAppliedInstructionId`.
2. If a channel offers an InstructionId equal to this, skip with `INSTRUCTION_ALREADY_APPLIED` (200 OK).

---

## 6. Failure Semantics

| Channel | Failure mode | Code | Worker action |
|---|---|---|---|
| A — Push | `Env=Production` at Main | `WORKER-403-01 PUSH_DISABLED_IN_PRODUCTION` | n/a (caller sees 403) |
| A — Push | Worker rejects JID | per `spec/14/28` §6 | reply with code, no retry by Main beyond `MaxRetries` |
| B — Pull from Main | Main unreachable | `WORKER-503-01 MAIN_UNREACHABLE` | exponential backoff, fall through to Channel C if enabled |
| B — Pull from Main | Stale cache (ETag mismatch) | n/a (just refetch) | refetch with no `If-None-Match` |
| C — Pull from URL | Manifest unreachable | `WORKER-503-02 MANIFEST_UNREACHABLE` | retry per `PollIntervalSeconds` schedule |
| C — Pull from URL | `PayloadUrl` host not allowlisted | `WORKER-403-02 PAYLOAD_HOST_NOT_ALLOWED` | refuse, alert via Heartbeat |
| C — Pull from URL | Signature invalid | `WORKER-600-01 PAYLOAD_VERIFICATION_FAIL` | refuse, alert via Heartbeat |

All `WORKER-*` codes registered in `spec/03-error-manage/03-error-code-registry/` per `13-error-codes.md`.

---

## 7. Implementer Checklist

```
[ ] Channel A endpoints reject with 403 when Env=Production (test exists).
[ ] Channel B poll loop runs at PollIntervalSeconds ± JitterSeconds.
[ ] Channel B respects HTTP 304 Not Modified.
[ ] Channel C disabled by default; opt-in via Seedable-Config.
[ ] Channel C enforces AllowedHostsAllowlist before download.
[ ] All channels share the same RenameFirst execution path (no duplicate deploy code).
[ ] LastAppliedInstructionId persisted in Settings tier.
[ ] Heartbeat reports CurrentVersion and LastAppliedInstructionId on every beat.
[ ] Push (A) writes AccessDenialEvent on every Production refusal.
```

---

## 8. Cross-References

- `00-overview.md` §0 — Kubernetes mental model.
- `06-core-api-endpoints.md` §2.5 — push endpoints (Channel A wire surface).
- `09-self-update-pointer.md` — high-level intent + pause note (this file is the technical realization).
- `15-tunable-constants.md` §2.10 — defaults for all channel knobs.
- `spec/14-update/28-worker-push-instruction.md` — JID schema (shared by all three channels).
- `spec/03-error-manage/03-error-code-registry/` — `WORKER-*` codes.
- `spec/06-seedable-config-architecture/` — config keys home.
- `13-error-codes.md` — worker error code prefixes.

---

*Update channels v1.0.0 — 2026-05-06*
