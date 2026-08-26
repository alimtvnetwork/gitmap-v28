# Fixture Encoding Conventions

These rules apply to **every** JSON fixture in this folder. A blind AI MUST honor them verbatim — they resolve the ambiguities flagged in audit 08 §2.2.

---

## 1. JSON

- **Encoding:** UTF-8, no BOM, LF line endings.
- **Indentation:** 2 spaces. Trailing newline at EOF. No trailing commas.
- **Key casing:** `PascalCase` for every field name (matches `mem://style/naming-conventions`). No `snake_case`, no `camelCase`. Acronyms are fully uppercased (`UserId`, `JwtExpiresAt`, `WorkerJwt`, `CompanyId`).
- **Key ordering:** Insertion order matters for golden-file diffs. Order shown in each fixture is canonical: identifying fields first (`*Id`, `CorrelationId`), then payload, then metadata (`OccurredAt`, `EnvelopeVersion`).
- **Numbers:** Integers for IDs and epoch seconds. No quoted numerics.

## 2. Timestamps

| Context | Format | Example |
| --- | --- | --- |
| **Wire (JSON body)** | RFC 3339 / ISO 8601 with `Z` suffix, second precision | `"2026-05-06T10:22:14Z"` |
| **DB columns** | Epoch seconds, UTC, INTEGER | `1762423334` |
| **Logs** | RFC 3339 with millisecond precision | `2026-05-06T10:22:14.812Z` |

JSON body NEVER carries epoch seconds. DB rows NEVER carry ISO strings. This split is intentional and binding.

## 3. Null vs omitted

- All envelope extension fields (`OperationId`, `SubCode`, `FieldErrors`) are **always present** with explicit `null` when not applicable. Never omit.
- All optional request fields default to `null` over omission when the server tolerates absence. Fixtures show `null` explicitly so a literal AI does not invent omission.
- Arrays default to `[]`, never `null`, when "no items" is meaningful.

## 4. Headers (shown alongside JSON in `request.*.json`)

Each request fixture's top-level `_headers` block is **non-wire metadata** (the underscore prefix marks it as such) showing the HTTP headers a literal client must send. The body proper sits under `_body`. Splitting body and headers in one file keeps the fixture self-contained without needing a curl script.

Mandatory headers on every Main↔Worker call (per `06-` §1):
- `X-Correlation-Id` — ULID (Crockford base32, 26 chars). Server generates if missing on inbound.
- `X-Idempotency-Key` — UUID v4, mandatory on POST/PUT/PATCH. Optional on GET/DELETE.
- `Content-Type: application/json; charset=utf-8` on bodies.
- `Authorization: Bearer <jwt>` when the endpoint row in `06-` §2 lists `JWT (on Worker)` as the auth method.

## 5. Error envelope

Every error response uses the shape from `08-` §2 verbatim. Examples in `errors/` are minimal-but-complete — no field omitted, no field invented.

## 6. JWT fixtures

`jwt/worker-jwt.rs256.json` and `jwt/partial-auth-jwt.json` show the **decoded** header and payload as separate JSON objects. The companion `worker-jwt.example.txt` shows the three-part dot-joined string with **fake** signatures (real RS256 signatures depend on key material outside the spec).

Claim ordering inside the payload follows `12-` §6: `iss`, `sub`, `aud`, `exp`, `iat`, `jti`, `cmp`, `wnk`, `kid`, `roles`. A literal verifier MUST NOT assume claim order, but golden files preserve it for stable diffs.
