# JWT Standards (AI Execution Prompt)

> **/goal** Enforce strict JWT token lifecycle management, ephemeral access tokens, mandatory refresh token rotation, and secure browser cookie storage.
> **/learn** Internalize the risks of token replay attacks, ban `localStorage` for JWT tokens, enforce `HttpOnly` and `SameSite=Strict` flags, and prevent PII leakage in payload claims.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Configure short token expiration times (maximum 15 minutes) for all JWT access tokens.
- [ ] `/learn` Never store access or refresh tokens in `localStorage` or unencrypted client storage.
- [ ] `/goal` Mandate refresh token rotation and store refresh tokens exclusively in `HttpOnly`, `Secure`, `SameSite=Strict` cookies.
- [ ] `/learn` Restrict JWT payloads to non-sensitive claims (e.g., `UserId`, role enums) and strictly exclude PII.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

## 1. Token Lifecycles

- **Access Tokens:** Must have a short expiration time (e.g., 15 minutes).
- **Refresh Tokens:** Used to obtain new access tokens. Must be rotated upon use to detect and prevent replay attacks.

## 2. Storage

- **Access Tokens:** Should be stored in memory or in a short-lived closure. Never store access tokens in `localStorage`.
- **Refresh Tokens:** MUST be stored in `HttpOnly`, `Secure`, `SameSite=Strict` cookies. They should never be accessible via JavaScript.

```typescript
// ❌ FORBIDDEN: Storing JWT tokens in browser localStorage
localStorage.setItem("authToken", token);

// ✅ REQUIRED: In-memory store for access token, HttpOnly cookie for refresh token
authMemoryStore.setToken(token);
res.cookie("refreshToken", refreshToken, {
  httpOnly: true,
  secure: true,
  sameSite: "strict",
});
```

## 3. Payload Constraints

- Never include PII (Personally Identifiable Information) or sensitive secrets in the JWT payload. The payload is Base64 encoded and can be read by anyone. Include only necessary identifiers like `UserId` and role claims.

```typescript
// ❌ FORBIDDEN: Including PII in JWT payload
const badPayload = { userId: "u123", email: "user@example.com", ssn: "000-12-3456" };

// ✅ REQUIRED: Minimal non-sensitive claims only
const goodPayload = { sub: "u123", role: UserRole.Member };
```

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-SEC-001: JWT Token Lifecycle & HttpOnly Cookie Storage

**Given** Authentication services handling JWT generation and validation
**When** Evaluating token lifetimes, storage mechanisms, and payload claims
**Then** Access tokens expire within 15 minutes, refresh tokens are stored in HttpOnly cookies, and payloads contain no PII

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/11-security/02-jwt-standards.md --check-only
```
**Expected:** exit 0. Zero violations.
