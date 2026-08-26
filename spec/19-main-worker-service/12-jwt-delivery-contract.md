# 12 — JWT Delivery Contract (XSS-Safe)

**Spec:** `19-main-worker-service`
**Version:** 1.2.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** audit findings F-A-12, F-D-04, F-B-05 (top-10 fix #3). Closes AC-4.
**Authority:** Canonical contract for how JWTs travel between Main, the React UI, and Workers. On any conflict with `05-auth-and-2fa.md`, diagram files, or other prose, **this file wins**.

---

## 1. The two tokens (do not confuse)

| Token | Carrier | Audience | Storage |
|---|---|---|---|
| **Main session cookie** | `HTTPOnly`, `Secure`, `SameSite=Lax` cookie set by Main | React UI ↔ Main | Browser cookie jar (no JS access) |
| **Worker JWT (RS256)** | JSON response body field `WorkerJwt` | React UI ↔ Worker | **In-memory only** (React state) |

The Main session cookie is **not a JWT**. The Worker JWT is **not delivered as a cookie**. This separation is intentional — see §3.

---

## 2. Authoritative delivery decision

### Worker JWT delivery: **JSON response body, in-memory storage**

**Decision:** The Worker JWT is returned in the **JSON body** of `GET /API/V1/Company/{CompanySlug}/Resolve` (per `06-core-api-endpoints.md` §3.2). React stores it **in memory only** (closure / `useRef` / Zustand non-persistent store). It is **never** written to `localStorage`, `sessionStorage`, IndexedDB, or any cookie.

**Decision is binding.** Implementers MUST NOT:

- ❌ Place the Worker JWT in a cookie (any flags).
- ❌ Persist it across page reloads.
- ❌ Pass it via URL query string.
- ❌ Store it in `localStorage`/`sessionStorage`.

After page reload, React calls `/API/V1/Company/{CompanySlug}/Resolve` again (cheap; uses Main session cookie) to get a fresh JWT.

---

## 3. Why not cookies for the Worker JWT?

| Concern | Cookie delivery | JSON-body + in-memory delivery (chosen) |
|---|---|---|
| XSS exfiltration | HTTPOnly blocks JS reads ✅ | JS can read it ⚠️ |
| CSRF | Auto-sent by browser ⚠️ | Never auto-sent ✅ |
| Cross-domain (React → Worker on different host) | Requires `SameSite=None; Secure` + CORS credentials ⚠️ | Just an `Authorization` header ✅ |
| Re-issue cost | Sticky until expiry ⚠️ | Trivial (refetch on resolve) ✅ |
| Survives page reload | Yes ⚠️ | No ✅ (forces re-resolve, re-checks tenant→worker mapping) |

**Trade-off accepted:** XSS risk is mitigated by (a) strict CSP per §5, (b) short TTL (15 min default per `15-tunable-constants.md` §2.4 `MainWorker.Auth.WorkerJwtTtlSeconds`), and (c) JWT being scoped to one `WorkerNodeId` + one company. CSRF is the larger systemic risk for cookie-borne API tokens, so we eliminate it.

---

## 4. End-to-end token flow

```
Browser (React)              Main                       Worker
     |                         |                          |
1. POST /Auth/SignIn --------->|                          |
2. <-- 200 + Set-Cookie:       |                          |
       MainSession=...; HTTPOnly; Secure; SameSite=Lax    |
3. GET /Company/{slug}/Resolve >|                          |
   (browser auto-attaches      |                          |
    MainSession cookie)        |                          |
4. <-- 200 JSON                |                          |
   { WorkerEndpoint, WorkerJwt, JwtExpiresAt, ... }       |
5. React stores WorkerJwt in   |                          |
   useRef / Zustand (memory)   |                          |
6. GET /Company/{slug} ----------------------->|          |
   Authorization: Bearer <WorkerJwt>           |          |
   X-Correlation-Id: ...                       |          |
                                               |--verify  |
                                               |  RS256   |
                                               |  + claim |
                                               |  wnk     |
   <-----------------200 JSON-------------------|          |
```

---

## 5. Mandatory Content Security Policy

To make the in-memory storage genuinely safer, both Main and Worker tiers MUST serve a CSP that blocks the common XSS exfiltration vectors:

```
Content-Security-Policy:
  default-src 'self';
  script-src 'self';
  connect-src 'self' https://*.example.com;
  frame-ancestors 'none';
  base-uri 'none';
  object-src 'none';
  form-action 'self';
```

`connect-src` MUST list the worker host pattern explicitly. No `unsafe-inline`, no `unsafe-eval`. Hashes / nonces required for any inline `<script>`.

This is non-negotiable — without CSP, in-memory storage gives no real XSS advantage over `localStorage`.

---

## 6. Refresh & rotation

| Event | React behavior |
|---|---|
| JWT within `MainWorker.Auth.JwtRefreshLeadSeconds` (default 60 s per `15-tunable-constants.md` §2.4) of `JwtExpiresAt` | Call `POST /API/V1/Auth/RefreshWorkerToken` (uses Main session cookie); replace in-memory token. |
| Page reload / new tab | In-memory token is gone; React calls `/Company/{slug}/Resolve` to get a fresh one. |
| Main session cookie expires | All Worker JWTs become unrefreshable; React redirects to `/sign-in`. |
| Worker `kid` mismatch (signing key rotated) | Worker returns `WORKER-100-02 KID_UNKNOWN`; React calls `/Resolve` to get a JWT signed by the new key. |
| User clicks "Sign out" | Call `/Auth/SignOut` (clears Main cookie); React drops in-memory JWT. |
| User clicks "Sign out everywhere" | Call `/Auth/SignOutAll`; same as above + invalidates all sessions server-side. |

---

## 7. Worker-side JWT verification (canonical)

Worker MUST verify, in order:

1. RS256 signature against `JwtPublicKeyPem` from bootstrap (`10-worker-bootstrap-protocol.md` §3.2).
2. `kid` matches stored `JwtSigningKeyId`. Mismatch → `WORKER-100-02 KID_UNKNOWN` (401).
3. `exp` is in the future, allowing `MainWorker.Auth.ClockSkewToleranceSeconds` (default 60 s per `15-tunable-constants.md` §2.4) of clock skew.
4. `wnk` claim equals this worker's `WorkerNodeId`. Mismatch → `WORKER-100-03 WRONG_WORKER` (403).
5. `iss` claim equals the configured `MainBaseUrl` host.
6. `aud` claim equals `worker`.

Any failure → registered error code (per `spec/03-error-manage/03-error-code-registry/`); never silently fall through.

---

## 8. JWT claim contract (canonical)

```jsonc
{
  "iss": "main.example.com",       // Main host
  "aud": "worker",                  // fixed string
  "sub": 42,                        // UserId
  "wnk": 3,                         // WorkerNodeId
  "cmp": "riseup-asia",             // CompanySlug
  "rls": ["PowerAdmin"],            // role codes
  "iat": 1714824000,                // RFC 7519 epoch
  "exp": 1714824900,                // iat + 900s default
  "jti": "01J...ULID",              // unique per token
  "kid": "main-jwt-2026-q2"         // matches JwtSigningKeyId
}
```

No PII beyond `UserId`. No password hashes. No 2FA secrets.

---

## 9. Test cases (mandatory for CI)

| # | Scenario | Expected |
|---|---|---|
| T-1 | Worker JWT placed in `localStorage` by React | Lint/test fail |
| T-2 | Worker JWT in cookie | Lint/test fail (no `Set-Cookie` from Main carries `WorkerJwt`) |
| T-3 | Page reload retains in-memory JWT | Fail (must be gone) |
| T-4 | CSP missing or contains `unsafe-inline` | Fail |
| T-5 | JWT TTL > `MainWorker.Auth.WorkerJwtTtlSeconds` (default 900 s per `15-tunable-constants.md` §2.4) | Fail |
| T-6 | Worker accepts JWT with wrong `wnk` | Fail (must return 403) |
| T-7 | Worker accepts JWT with expired `exp` | Fail (must return 401) |
| T-8 | URL contains `?token=` or `?jwt=` | Fail |
| T-9 | JWT logged in any access log | Fail (must be redacted) |

These tests live in `spec/19-main-worker-service/97-acceptance-criteria.md` AC-4.

---

## 10. Diagram corrections (apply during next pass)

| Diagram | Stale element | Replace with |
|---|---|---|
| `seq-login-routing.mmd` | "JWT in JSON body" with no XSS note | Add caption: "JWT delivered in JSON body; React stores in-memory only per `12-jwt-delivery-contract.md`. CSP per §5 mandatory." |
| `seq-company-creation.mmd` | implicit JWT delivery | Add same caption. |

Tracked as follow-up FU-7 / FU-8 (extends §8 list of `11-split-db-tier-reconciliation.md`).

---

## 11. Open Questions — formalised dispositions

This section formalises the dispositions of v1.0 open questions so dumb-AI implementers do not re-litigate them mid-implementation. Each disposition is **binding** for v1.0; the Future-Work block at the bottom names what a v2.0 spec MUST cover if a disposition is reopened.

### 11.1 OQ-12-1 — Separate `RefreshToken` cookie (rotation pattern)

**Question.** Should the Worker-JWT refresh path issue a dedicated `RefreshToken` cookie (rotated on every refresh) instead of relying on the long-lived Main session cookie?

**Disposition for v1.0:** ❌ **Deferred to v2.0. Implementers MUST NOT introduce a `RefreshToken` cookie in v1.0.**

**Rationale.**

| Concern | Main session cookie (chosen for v1.0) | `RefreshToken` cookie (deferred) |
|---|---|---|
| Refresh authority | Main session cookie already proves the user is signed in; refresh is a one-line lookup | Adds a parallel auth surface — cookie + DB row + rotation table |
| Theft window | Cookie carries the same risk profile already accepted in §3 | Adds a second long-lived secret with a different lifetime — twice the rotation surface to audit |
| Server state | Zero new tables | Requires `RefreshTokenRotation` table with `(TokenHash, IssuedAtEpoch, RotatedFromHash, RevokedAtEpoch)`, sweeper, and replay-detection family of error codes |
| Sign-out-everywhere | Already works via `/Auth/SignOutAll` invalidating the session row | Requires cascading revocation across rotation chain |
| CODE RED footprint | Adds 0 functions | Adds ≥4 functions (issue / rotate / revoke / detect-replay), each must obey the 8–15 line + 0-nesting + ≤2-operand caps |
| What it would actually buy | nothing in the current threat model — the Main session cookie's TTL is already the controlling factor | Defence against a stolen refresh token being usable for the full session lifetime; only matters once we ship long-lived "remember me" sessions |

**Trigger conditions that would force OQ-12-1 to be reopened in v2.0** (any one is sufficient):

1. The product introduces opt-in **"remember me"** sessions whose lifetime exceeds `MainWorker.Auth.MaxSessionLifetimeSeconds` (currently uncapped — see §11.3 future-work).
2. A formal threat model adopts **session-hijack containment** as a v2.0 requirement (currently it is not — XSS containment per §3 + CSP per §5 is the chosen mitigation).
3. Compliance regime (e.g. SOC2 control change) requires **per-refresh credential rotation evidence** in the audit log.

**Forbidden v1.0 implementations** (any of these is a CODE RED violation):

- ❌ Adding a `Set-Cookie: RefreshToken=…` header anywhere in the codebase.
- ❌ Persisting refresh tokens in any storage (DB row, KV, or file).
- ❌ Rotating the Main session cookie value on `/Auth/RefreshWorkerToken` (the cookie value MUST be stable for the session's lifetime; only the in-memory Worker JWT rotates).
- ❌ Inventing a `WORKER-100-04 REFRESH_REPLAY` style error code "for future use" — error codes are added when the feature ships, not before.

### 11.2 OQ-12-2 — Service-worker-assisted token storage (off-thread)

**Question.** Should the Worker JWT live in a Service Worker's memory (off-main-thread, inaccessible to page scripts) instead of in React state?

**Disposition for v1.0:** ❌ **Rejected. Out of scope for v1.0 and not on the v2.0 roadmap.**

**Rationale.**

| Concern | React state (chosen) | Service Worker storage (rejected) |
|---|---|---|
| XSS access path | `window.tokenRef.current` reachable from any same-origin script | `postMessage` to SW is reachable from any same-origin script — same blast radius |
| Reload survival | None (forces re-resolve — a feature, per §3) | Survives reload (re-introduces the persistence we explicitly chose against) |
| Browser support matrix | 100% | SW lifecycle and termination behavior varies across browsers; non-trivial debugging surface |
| Network interception | None | SW intercepts `fetch` — every API call now traverses an extra worker boundary |
| What it would actually buy | nothing in the current threat model — the SW is not a privilege boundary against same-origin XSS | A perceived "off-thread" benefit that the spec authors of `12` consider security theatre absent a cross-origin isolation story |

**Trigger conditions that would force OQ-12-2 to be reopened:**

1. Adoption of a formal **cross-origin-isolated** delivery model (`Cross-Origin-Opener-Policy: same-origin` + `Cross-Origin-Embedder-Policy: require-corp`) where the SW can be promoted to a privilege boundary.
2. Browser standards landing a **"trusted SW token store"** primitive that is genuinely unreachable from page script (does not currently exist).

Until **both** trigger conditions hold, this OQ stays rejected — re-opening it without them is a CODE RED violation of "no swallowed reasons" (the swallowed reason being "it sounds more secure").

### 11.3 Future-work catalogue (v2.0 prerequisites)

If/when the product reopens §11.1 or §11.2, the v2.0 spec MUST specify, in this order:

1. **`MainWorker.Auth.MaxSessionLifetimeSeconds`** tunable in `15-tunable-constants.md` §2.4 (currently uncapped — adding a cap is the precondition for "remember me").
2. **Threat-model document** at `spec/19-main-worker-service/24-threat-model.md` (does not yet exist — placeholder slot 24 is reserved by this clause).
3. Net-new error-code family in `13-error-codes.md` for the chosen rotation/revocation flow.
4. Acceptance-criterion row in `97-acceptance-criteria.md` proving the new path is actually exercised in CI.

**Until that v2.0 spec lands, the four "Forbidden v1.0 implementations" in §11.1 and the two "trigger conditions" in §11.2 are the binding contract.**

---

## 12. Cross-references

- `spec/19-main-worker-service/05-auth-and-2fa.md` — high-level auth contract.
- `spec/19-main-worker-service/06-core-api-endpoints.md` §3.2 — `/Resolve` response shape.
- `spec/19-main-worker-service/10-worker-bootstrap-protocol.md` §3.2 — `JwtPublicKeyPem` + `JwtSigningKeyId` source.
- `spec/03-error-manage/03-error-code-registry/` — MUST register `WORKER-100-02 KID_UNKNOWN`, `WORKER-100-03 WRONG_WORKER`.

---

*JWT delivery contract v1.2.0 — 2026-05-06 (Phase 12.4 — §11 OQ-12-1 / OQ-12-2 dispositions formalised: forbidden v1.0 patterns + v2.0 trigger conditions + future-work catalogue).*

---

## 13. Backup-tier S2S tokens (Phase 12 stub)

**Added:** Phase 12 (Backup-tier consolidation). **Authority for full contract:** `21-backup-endpoints.md` §3.

The Backup tier (BE-1..BE-6 endpoints) does **not** use the React-facing Worker JWT defined in §1–§8. It uses a separate S2S OAuth client-credentials token with its own audience and a mandatory `PairingId` claim.

### 13.1 Backup token shape (canonical)

```jsonc
{
  "iss": "main.example.com",          // Main host (unchanged)
  "aud": "Backup",                     // disjoint from "worker"
  "sub": "PairingId:01J...ULID",       // S2S principal — the pairing, not a user
  "PairingId": "01J...ULID",           // MUST equal sub's pairing; mandatory
  "scope": "Backup.Sync.Write",        // one of the four capabilities
  "iat": 1714824000,
  "exp": 1714824900,                   // same TTL as orchestration tokens
  "jti": "01J...ULID",
  "kid": "main-s2s-2026-q2"
}
```

### 13.2 Verification rules (in order)

1. RS256 signature against the orchestration signing key (reuses §7 verification path with `aud="Backup"` branch).
2. `aud` claim equals **exactly** `"Backup"`. UI tokens (`aud="worker"`) and orchestration tokens (`aud="main-orchestration"`) MUST be rejected here — no audience downgrade.
3. `PairingId` claim is present **and** matches a row in the receiving node's `BackupPairing` table (per `19-incremental-backup-sync.md` §5).
4. `scope` claim matches the endpoint's required scope (per `21-backup-endpoints.md` §3 table).
5. `exp` in future with `MainWorker.Auth.ClockSkewToleranceSeconds` tolerance.

Any failure → **HTTP 421 Misdirected Request** with error code `MAIN-800-04 BackupPairingMismatch` (per `13-error-codes.md`). CODE RED: no silent passthrough.

### 13.3 Why disjoint from the Worker JWT

| Concern | Worker JWT (`aud="worker"`) | Backup S2S (`aud="Backup"`) |
|---------|------------------------------|------------------------------|
| Carrier | JSON body, in-memory React | OAuth `Authorization: Bearer`, S2S only |
| Principal | End user (`sub=AppUserId`) | Pairing (`sub=PairingId:...`) |
| TTL | 15 min (UI refreshable) | 15 min (machine refreshable via client-creds) | <!-- TUNABLE-WAIVER: 15 min = MainWorker.Auth.WorkerJwtTtlSeconds (900s); table cell, key cited in body §2.4 -->
| XSS exposure | Designed-around per §3 | None — never reaches a browser |
| Cross-node misroute | Detected by `wnk` claim | Detected by `PairingId` claim → 421 |

### 13.4 Test cases (extend §9)

| # | Scenario | Expected |
|---|----------|----------|
| T-10 | UI Worker JWT (`aud="worker"`) sent to BE-1 | 421 `MAIN-800-04` |
| T-11 | Backup token sent to a non-paired Backup node | 421 `MAIN-800-04` |
| T-12 | Backup token missing `PairingId` claim | 421 `MAIN-800-04` |
| T-13 | Backup token with valid `PairingId` but wrong `scope` for the endpoint | 403 (per `21-backup-endpoints.md` §3) |

---

*JWT delivery contract v1.1.0 — 2026-05-06 (Phase 12: §13 Backup S2S tokens stub added)*
