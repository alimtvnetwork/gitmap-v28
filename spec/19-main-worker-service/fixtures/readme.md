# spec/19 — Golden-File Fixtures

**Purpose:** Wire-format ground truth for blind-AI implementers. Every fixture in this folder is **authoritative for shape, field names, ordering, and encoding** of the example it depicts.

**Authority precedence:** When prose in `spec/19-…/*.md` and a fixture here disagree on a wire-format detail (field name, casing, ordering, null-vs-omitted, timestamp encoding), the **fixture wins**. When they disagree on *behavior* (algorithm, validation rule, retry policy), the **prose wins**.

**Created:** 2026-05-06 (Phase 13.3 — closes audit 08 §2.2 fixture gap).

---

## Layout

```
fixtures/
├── readme.md                       (this file)
├── conventions.md                  (encoding rules — read first)
├── endpoints/                      (one folder per endpoint family in 06-)
│   ├── auth-signin/
│   │   ├── request.success.json
│   │   ├── response.success.json
│   │   └── response.2fa-required.json
│   ├── company-resolve/
│   │   ├── request.success.json
│   │   └── response.success.json
│   ├── company-create/
│   │   ├── request.success.json
│   │   └── response.success.json
│   ├── workers-register/
│   │   ├── request.success.json
│   │   └── response.success.json
│   ├── workers-heartbeat/
│   │   ├── request.success.json
│   │   └── response.success.json
│   ├── settings-endpointauth-patch/
│   │   ├── request.success.json
│   │   └── response.success.json
│   └── version/
│       └── response.success.json
├── errors/                         (one example per error family in 13-)
│   ├── worker-unreachable.json
│   ├── worker-version-mismatch.json
│   ├── splitdb-write-fail.json
│   ├── validation-fail.json
│   ├── idempotency-conflict.json
│   ├── totp-backup-exhausted.json
│   └── endpoint-auth-locked.json
└── jwt/
    ├── worker-jwt.example.txt      (3-part HS-style for illustration)
    ├── worker-jwt.rs256.json       (decoded header + payload, RS256)
    └── partial-auth-jwt.json       (decoded header + payload, 2FA gate)
```

## Coverage matrix (audit 08 §2.2)

| Source spec | Fixture surface | Status |
| --- | --- | ---: |
| `06-` §2.1 Auth | `endpoints/auth-signin/` | ✅ |
| `06-` §2.2 Company | `endpoints/company-resolve/`, `endpoints/company-create/` | ✅ |
| `06-` §2.4 Status/Version | `endpoints/version/` | ✅ |
| `06-` §2.5 Workers | `endpoints/workers-register/`, `endpoints/workers-heartbeat/` | ✅ |
| `06-` §2.7 Settings | `endpoints/settings-endpointauth-patch/` | ✅ |
| `08-` §3 + `13-` Worker tier | `errors/worker-unreachable.json`, `errors/worker-version-mismatch.json`, `errors/splitdb-write-fail.json` | ✅ |
| `08-` §3 + `13-` Main tier | `errors/validation-fail.json`, `errors/idempotency-conflict.json`, `errors/totp-backup-exhausted.json`, `errors/endpoint-auth-locked.json` | ✅ |
| `12-` JWT | `jwt/worker-jwt.rs256.json`, `jwt/partial-auth-jwt.json`, `jwt/worker-jwt.example.txt` | ✅ |

Endpoints not covered (intentional — trivial wrappers over the patterns above): `SignOut`, `SignOutAll`, `RefreshWorkerToken`, `Status`, `SelfUpdate`, `Workers/{id}/Update`, `Workers/All/Update`, `Workers/PublishZip`, all CRUD on `/Users`, `/Company` `PATCH`/`DELETE`. These reuse the envelope + auth headers documented in `conventions.md`.
