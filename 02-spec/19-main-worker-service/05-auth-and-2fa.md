# 05 — Authentication and 2FA

**Spec:** `19-main-worker-service`
**Version:** 2.1.0

> **v2.0.0 (Phase 3 — Users moved off Main).** Main is now a **credential-blind reverse proxy** for `/Auth/*` traffic. It owns the routing index `UserDirectory` (`03-main-db-schema.md` §2.4) and nothing else. Password hashes, TOTP secrets, backup codes, and role assignments live exclusively on the assigned Worker's split-DB App tier (`AppUser`, `AppUserRole`). The flows in §5–§7 below are rewritten accordingly. The previous v1.0 flow — where Main verified passwords locally — is **removed**.

Auth is a **first-class given** in BOTH Main and Worker tiers. This file defines the contract; implementer chooses Laravel Sanctum / Passport / custom JWT as long as the contract is honored.

---

## 1. Capability matrix (v2.0.0)

| Capability | Main | Worker |
|-----------|------|--------|
| Email + password sign-up (HTTP entry point) | ✅ proxy | ✅ authoritative |
| Email + password sign-in (HTTP entry point) | ✅ proxy | ✅ authoritative |
| Password hashing / verification | ❌ never | ✅ |
| TOTP enroll + verify | ❌ never | ✅ |
| TOTP backup-code storage | ❌ never | ✅ |
| `UserDirectory` routing index | ✅ | ❌ (mirrored read-only via bootstrap) |
| Session management (UI cookie) | ✅ | optional |
| Worker-JWT issuance | ❌ | ✅ (Worker mints; Main forwards) |
| Cookie-based session for React UI | ✅ | optional |
| Password reset (email link) | ✅ proxy | ✅ authoritative |
| Sign-out (single + all sessions) | ✅ | ✅ |

> **Why Main proxies instead of authenticates.** Per locked decision D5, Main MUST be a thin catalog. Storing password hashes on Main would re-introduce the cross-tenant blast radius the split-DB architecture exists to eliminate. The proxy pattern keeps Main credential-blind while preserving a single public entry URL.

---

## 2. Two Authentication Surfaces

### 2.1 User → Main → Worker (UI-facing, credential proxy)

- **Entry:** browser POSTs `/API/V1/Auth/SignIn` to Main with `{ Email, Password, TotpCode? }`.
- **Main step 1 — routing lookup:** Main reads `UserDirectory WHERE UserEmail = LOWER(:email)` to obtain `WorkerNodeId`. **Constant-time response on miss** (no email enumeration): if no row, Main forwards to a synthetic "null Worker" handler that returns the same generic 401 envelope after the same wall-clock budget as a real verification.
- **Main step 2 — forward:** Main proxies the original body verbatim over a mutual-TLS internal channel to `POST {WorkerEndpoint}/API/V1/Auth/InternalSignIn` with headers `X-Forwarded-For-User: <hash(email)>`, `X-Correlation-Id`, and an OAuth client-credentials Bearer (per §2.3). **Main does NOT log the password** (scrubbed per §3).
- **Worker step:** Worker reads its `AppUser` row, runs the password verifier (Argon2id / bcrypt per §3), checks TOTP if enrolled, then mints the Worker-JWT (per `12-jwt-delivery-contract.md`).
- **Main step 3 — session cookie:** On Worker 200, Main stamps a session cookie (HTTPOnly, Secure, SameSite=Lax) bound to the Worker-JWT and returns the JWT body to the browser. On Worker 401, Main returns 401 unchanged.
- **Main never sees the cleartext password after forwarding.** The proxied body buffer MUST be zeroed (`sodium_memzero` or equivalent) immediately after the forward call returns.

### 2.2 React → Worker (data-facing, after sign-in)

- **Mechanism:** short-lived JWT minted by **the Worker** (not Main, as in v1.0), accepted by that same Worker. JWT claims unchanged from v1.0 except `iss = WorkerEndpoint` (was Main URL).
- **JWT claims:**
  - `sub` = `AppUserId` (Worker-local)
  - `cmp` = `CompanyId`
  - `wnk` = `WorkerNodeId` (so Worker rejects misrouted tokens)
  - `iss` = Worker URL
  - `aud` = Worker URL (self-issued)
  - `exp` = issued + `MainWorker.Auth.WorkerJwtTtlSeconds` (default per `15-tunable-constants.md` §2.4 = 15 min)
  - `iat` = epoch
  - `roles` = array of `RoleCode` strings (Worker computed cascading union per Phase 5)
- **Signing:** asymmetric (RS256). Worker holds private key; Main holds the Worker's public key (rotatable via Seedable-Config) so it can validate JWTs on session-cookie refresh.
- **Refresh:** React calls Worker `/API/V1/Auth/RefreshToken` directly (Main bypassed) when JWT is within `MainWorker.Auth.JwtRefreshLeadSeconds` of expiry (default per `15-tunable-constants.md` §2.4 = 60 s).

### 2.3 Main → Worker (orchestration: push-update, registry sync, credential proxy)

- **Mechanism:** OAuth 2.0 client-credentials grant OR pre-shared API key (configurable).
- **Default:** OAuth client-credentials per Worker, secrets stored via Seedable-Config (encrypted at rest).
- The credential-proxy channel (§2.1) reuses this same client-credential token; the Worker enforces that `/API/V1/Auth/InternalSignIn` is callable **only** from Main's IP allowlist + valid Bearer.
- Per-endpoint flexibility per verbatim §Main Server Concept 3c is **resolved** — see §8 below for the pinned `PATCH /API/V1/Settings/EndpointAuth` contract (OQ-1 RESOLVED 2026-05-04).

---

## 3. Password Storage (Worker tier only — v2.0.0)

> **Authority moved.** As of v2.0.0 (Phase 3), all password material lives in the Worker's `AppUser` table. Main MUST NOT persist `PasswordHash`, `PasswordSalt`, or `PasswordPepper`. Tunable keys keep the `MainWorker.Auth.*` namespace for compatibility but are now resolved on the Worker.

Per verbatim §Login 3:

- **Hash:** Argon2id (preferred) or bcrypt. Bcrypt cost is **environment-pinned** to remove ambiguity (resolves F-A-03): `MainWorker.Auth.BcryptCost` defaults to `12` when `Env=dev|test`, `14` when `Env=prod|staging`. Implementations MUST refuse cost < 12 and MUST NOT exceed 14 unless overridden by Power Admin via Seedable-Config (caps at 16).
- **Salt:** unique per user, stored alongside hash on the **Worker** as `AppUser.PasswordSalt`. The chosen hash function may also embed a salt; storing it explicitly keeps the contract stack-portable.
- **Pepper:** per-Worker pepper from Worker-local Seedable-Config secret `MainWorker.Auth.PasswordPepper`. **MUST be set when `Env=prod|staging`**. In `dev|test` MAY be empty; if empty, the Worker MUST log a one-shot WARN at startup so drift is visible. Each Worker MAY hold a distinct pepper — this is a feature, not a bug, because it limits cross-Worker hash portability.
- **No plaintext anywhere — including Main's proxy buffer.** Logs MUST scrub `password`, `confirmPassword`, `currentPassword` on **both** tiers. Main's proxy buffer MUST be zeroed after forward (§2.1).
- **No retrieval.** Reset is replace-only.
- **Breach check:** Implementations MUST verify against a HIBP-style API on sign-up and password change when `MainWorker.Auth.EnableBreachCheck=true` (default `true` in prod, `false` in dev). The breach check runs on the **Worker** (so the cleartext password never reaches an external service from Main).

---

## 4. Two-Factor Authentication (2FA) — Worker tier only (v2.0.0)

> **Authority moved.** TOTP enrollment, verification, and backup-code storage all live on the **Worker** as of v2.0.0 (Phase 3). Main proxies the QR-render request and the verify POST to the assigned Worker, but stores no TOTP material itself.

- **Standard:** TOTP (RFC 6238), 30s window <!-- TUNABLE-WAIVER: RFC 6238 mandates 30s; not a MainWorker tunable -->, 6 digits.
- **Enrollment:** the user's browser POSTs `/API/V1/Auth/2FA/Enroll` to Main; Main resolves `Email → Worker` via `UserDirectory` and forwards. The Worker generates the TOTP secret, writes `AppUser.TotpSecret` (encrypted at rest with key from Worker-local Seedable-Config), returns the `otpauth://` URI. Main proxies the URI back to the browser unmodified. The browser submits one TOTP to confirm via the same Main-proxied `/API/V1/Auth/2FA/Confirm` path.
- **Backup codes:** generate 10 single-use codes at enrollment (stored as bcrypt hashes in `AppUser.TotpBackupCodesHash` on the Worker). When the count of unused codes reaches **0**, the user MUST be forced to regenerate at next sign-in: the Worker returns `Error.SubCode = TotpBackupExhausted` with HTTP 403 and `X-Auth-Action: RegenerateBackupCodes` header (resolves F-A-05). Main forwards the response unchanged. Regeneration invalidates the prior batch. Power Admin override path: `POST /API/V1/Auth/2FA/RegenerateBackupCodes` (audit logged on the Worker, mirrored to Main `EndpointAuthAuditEvent`).
- **Verification points:** sign-in, password change, 2FA disable, role escalation.
- **Recovery:** Power Admin can reset 2FA for any user (audit logged on the Worker).

---

## 5. Sign-Up Flow (Main as proxy, Worker as authority)

1. Browser POSTs `/API/V1/Auth/SignUp` to **Main** with `{ Email, Password, CompanySlug? }`.
2. Main runs cheap guards locally (no credentials required): `IsEmailWellFormed`, `IsPasswordWellFormed` (length / charset only — strength check runs on Worker), `IsCompanySlugAvailable` (read of `Company.Slug`).
3. Main resolves the target Worker:
   - **Existing company:** `WorkerNodeId = Company.WorkerNodeId WHERE Slug = :companySlug`.
   - **New company:** Main runs worker selection (`04-worker-routing.md`), creates `Company`, creates `UserDirectory` row pointing at the chosen Worker.
4. Main forwards the body to `POST {WorkerEndpoint}/API/V1/Auth/InternalSignUp` over the credential-proxy channel (§2.3). The Worker:
   - Runs `IsPasswordStrongEnough` (Worker-side policy), HIBP breach check, hashes the password, creates `AppUser` row with `PasswordHash` / `PasswordSalt`.
   - Returns `{ AppUserId, NeedsEmailVerification: bool }`.
5. Main writes `UserDirectory.LastSeenAt = now()`, zeroes the proxy buffer, returns 201 to the browser plus a session cookie (or returns "verify email" status if email-confirm flag is on).

**Failure modes:**
- Worker rejects password (weak / breached) → Main returns the Worker's 400 envelope unchanged, **does not** persist a `UserDirectory` row.
- Worker unreachable → Main returns 503 `Error.SubCode = WorkerUnreachable` with retry-after; `UserDirectory` row is created only after the Worker ACK.

---

## 6. Sign-In Flow (Main as proxy, Worker as authority)

1. Browser POSTs `/API/V1/Auth/SignIn` to **Main** with `{ Email, Password, TotpCode? }`.
2. Main reads `UserDirectory WHERE UserEmail = LOWER(:email)`. **Constant-time miss handling** per §2.1.
3. Main forwards the body to `POST {WorkerEndpoint}/API/V1/Auth/InternalSignIn` over the credential-proxy channel.
4. **Worker** verifies the password hash, checks TOTP status:
   - If TOTP enrolled and no `TotpCode` provided: Worker returns 200 with `{ "Step": "AwaitTotp", "ChallengeId": "..." }`. Main forwards verbatim. The browser POSTs `/API/V1/Auth/Verify2FA` with `{ ChallengeId, TotpCode }` → Main proxies to Worker `/API/V1/Auth/InternalVerify2FA`.
   - If credentials valid (with TOTP if required): Worker mints the Worker-JWT (claims per §2.2), returns:
     ```json
     {
       "WorkerEndpoint": "https://w3.example.com",
       "WorkerJwt": "<RS256 token>",
       "JwtExpiresAt": 1746360900
     }
     ```
     (`JwtExpiresAt` is now epoch seconds per Rule 7.1 v2.)
5. Main updates `UserDirectory.LastSeenAt`, zeroes the proxy buffer, sets the browser session cookie bound to the JWT, returns the body to the browser.
6. React stores the JWT in memory (NOT localStorage), uses it for direct Worker calls (Main bypassed for data-tier traffic).

---

## 7. Worker JWT Validation (Worker side)

Every Worker request validates:
1. Signature against the Worker's own public key (Worker is the issuer in v2.0.0; Main holds the public key only for session-cookie refresh).
2. `exp` not passed.
3. `aud` matches this Worker's URL.
4. `wnk` claim matches this Worker's `WorkerNodeId`.
5. `cmp` claim matches the `CompanyId` resolved from the requested resource.

Failure → 401 with `08-error-contract.md` envelope. NEVER 500 on auth failure.

---

## 8. Endpoint Authentication Defaults

| Endpoint pattern | Default auth | Configurable? |
|------------------|--------------|----------------|
| `/API/V1/Auth/*` | None (sign-up/in) or session (sign-out) | No |
| `/API/V1/Status` | None | Yes (admin can require auth) |
| `/API/V1/Version` | None | Yes |
| `/API/V1/Company/*` | Session (Main) or Worker JWT (Worker) | Per-endpoint via Seedable-Config |
| `/API/V1/Workers/*` | Session + `User has access to EnumPage.PowerAdminPage` | No (always protected) |
| `/API/V1/SelfUpdate` | OAuth client-credentials | No (always protected) |

> ✅ **Open Question OQ-1 — RESOLVED 2026-05-04 (task #39).**
> Per-endpoint auth-mechanism overrides ARE supported. Contract: single-row whole-replace `PATCH /API/V1/Settings/EndpointAuth` keyed by `EndpointPathPattern`, with `AcceptedMechanisms[]` allow-list, 7 validation rules, lock-list for `Workers/*` + `SelfUpdate`, and idempotent re-application. Full schema + semantics in `06-core-api-endpoints.md` §5. Every successful PATCH writes one `EndpointAuthAuditEvent` row (`03-main-db-schema.md` §2.6.4) inside the same transaction (FU-17 — RESOLVED 2026-05-05). Closes audit finding F-M-10.

---

## 9. Anti-patterns (CODE RED)

- ❌ `if user.role === 'admin'` — use `User has access to EnumPage.X` (`07-role-based-dashboards.md`).
- ❌ Storing JWTs in `localStorage` (XSS exposure).
- ❌ Long-lived worker JWTs (> 1 hour) <!-- TUNABLE-WAIVER: anti-pattern threshold, not a tunable; canonical TTL is MainWorker.Auth.WorkerJwtTtlSeconds in 15-tunable-constants.md §2.4 -->.
- ❌ Symmetric JWT signing across tiers (key sharing risk).
- ❌ Returning 500 on bad credentials (info leak; use 401 + neutral message).
- ❌ `if (!isAuthenticated)` — invert to `if (isAuthenticated)` and use early-return guards.

---

*Auth and 2FA v2.0.0 — 2026-05-06 (Phase 3: Main is credential-blind proxy; auth authority moved to Worker)*

---

## §11 — Backup S2S Audience (Phase 12 stub)

**Added:** Phase 12 (Backup-tier consolidation). **Authority for full contract:** `21-backup-endpoints.md` §3 + `12-jwt-delivery-contract.md` §13.

In addition to the OAuth client-credentials surface in §2.3 (Main → Worker orchestration), Workers acting as **Primary** or **Backup** nodes MUST accept a separate S2S audience for Backup-tier traffic (BE-1..BE-6 per `21-backup-endpoints.md`):

| Token field | Value | Purpose |
|-------------|-------|---------|
| `aud` | `Backup` | Disjoint from `worker` (UI JWT) and `main-orchestration` (push-update). |
| `scope` | One of: `Backup.Sync.Write`, `Backup.Ack.Read`, `Backup.Restore.Write`, `Backup.Restore.Apply` | Per-endpoint capability. |
| `PairingId` | Mandatory string claim | MUST match a row in `BackupPairing` on the receiving node, or request is rejected with **HTTP 421 Misdirected Request** (error code `MAIN-800-04`). |

**Rules:**
1. The `Backup` audience MUST NOT be granted to UI clients, React, or end-user sessions — S2S only.
2. A token missing `PairingId`, or carrying a `PairingId` not registered on the receiving node, MUST be rejected at the proxy layer **before** application logic runs (CODE RED: no silent fallback to other audiences).
3. Token signing key, TTL, and rotation cadence reuse the existing OAuth client-credentials infrastructure from §2.3; only the `aud` and `scope`/`PairingId` claims are new.

*See `21-backup-endpoints.md` §3 for the full enforcement table and `97-acceptance-criteria.md` AC "Backup-tier acceptance" for tests.*

---

*Auth and 2FA v2.1.0 — 2026-05-06 (Phase 12: §11 Backup S2S audience stub added)*
