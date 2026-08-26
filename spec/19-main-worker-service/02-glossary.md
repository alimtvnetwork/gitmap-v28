# 02 — Glossary

**Spec:** `19-main-worker-service`
**Version:** 1.1.0

Canonical terms. When this spec or any sibling spec uses these words, this is what they mean.

---

| Term | Definition |
|------|------------|
| **Main Server** | The coordinator deployment. Serves the React UI and edge REST endpoints. Runs no business logic. Owns the worker registry and tenant→worker mapping. Also called "Main", "Coordinator", "Master Node" (avoid the last; "Master" is reserved for human-facing language only). |
| **Worker Server / Worker Node** | An independent deployment that runs business logic and owns a split-DB. Has its own auth stack but no UI. Synonyms: "Worker", "Worker Node". |
| **Tenant Root** | The top-level multi-tenant entity. In this spec it is **`Company`**. |
| **Company** | The tenant root entity. Owns users. Mapped 1:N to a Worker. |
| **Power Admin** | Highest-privilege role. Held by application owners (Riseup Asia LLC). Can enable/disable endpoints and configure system-wide settings. |
| **Admin User** | Paying customer admin. Has an admin panel scoped to their `Company`. Cannot configure system-wide settings. |
| **EnumPage** | Compile-time enum identifying a logical "page" or capability gate (e.g., `EnumPage.AdminPage`, `EnumPage.PowerAdminPage`, `EnumPage.BillingPage`). Used in access checks. |
| **Access Check** | The authorization pattern `User has access to {EnumPage}`. Never `if user.role == X`. See `07-role-based-dashboards.md`. |
| **Split-DB** | The hierarchical SQLite layering defined in `spec/05-split-db-architecture/` (6 tiers total). Per FU-1 in [`11-split-db-tier-reconciliation.md`](./11-split-db-tier-reconciliation.md): **Worker uses 4 tiers** (Root / Settings / App / Session); **Main uses 3 tiers** (Root / Settings / Session, no App). Cache + Document tiers are unused in v1.0. Older "Root/App/Session" wording is stale. |
| **Seedable-Config** | The config-seeding mechanism in `spec/06-seedable-config-architecture/`. Replaces the legacy term `CW configuration`. |
| **Worker Registry** | Table on Main Server listing every known Worker (endpoint, identity, title, version). |
| **Worker Selection Strategy** | The algorithm Main uses to pick a Worker for a new tenant (round-robin, least-loaded, etc.). See `04-worker-routing.md`. |
| **Push Update** | Main-initiated update sent to one or all Workers. See `01-architecture.md` §3.3. |
| **Self-Update** | Worker- or Main-initiated update of its own deployment. POINTER ONLY in this spec — see `09-self-update-pointer.md`. |
| **Correlation ID** | Opaque per-request string (UUID v4) carried in `X-Correlation-Id` for cross-tier tracing. NOT a DB key. |
| **Idempotency Key** | Opaque client-supplied string in `X-Idempotency-Key` for mutating cross-tier calls. |
| **Edge Endpoint** | A REST endpoint hosted on Main that the React frontend hits first. |
| **Worker Endpoint** | A REST endpoint hosted on a Worker that React hits directly after Main resolves the worker. |
| **gitmap** | Replaces the term `git map`. The repo→snapshot system used elsewhere. |
| **Quarantined** | A `WorkerNodeStatus` value (per `03-` §2.2). Set by Main when a Worker reports `WorkerVersionMismatch`, fails repeated handshakes, or is administratively isolated. Quarantined Workers receive **no new routing** and **no heartbeat-driven recovery**; restart + clean Register call required to leave the state. (F-A-36) |
| **Draining** | A `WorkerNodeStatus` value (per `03-` §2.2). Set when a Worker is being gracefully removed (push-update, decommission). Existing tenants keep being served; new tenants are routed elsewhere. The Worker leaves the state by transitioning to `Offline` after its last in-flight request completes. (F-A-37) |
| **Seedable-Config** | The config-seeding mechanism defined in [`spec/06-seedable-config-architecture/`](../06-seedable-config-architecture/). SemVer-merged via GORM at boot; Categories block holds tunables, Tables block holds reference rows. See also `mem://architecture/seedable-configuration`. (F-A-38; superset of the row 8 lines above which is retained for backward search.) |
| **`apperror` package** | The error-wrapping convention defined in [`spec/03-error-manage/`](../03-error-manage/) and `mem://architecture/error-handling`. Every error crossing a layer boundary MUST be wrapped with explicit file path + operation name; raw `String(err)` is forbidden. (F-A-39) |
| **Power Admin** vs **`PowerAdmin`** | **`Power Admin`** (with space) is the human-facing label — use in UI strings, dashboard names, and prose. **`PowerAdmin`** (PascalCase, no space) is the code identifier — use in `RoleCode` rows, enum values, route middleware tokens, and any machine-parsed string. The two are NEVER interchangeable; linters MAY enforce the distinction. (F-A-40) |

---

## Reserved / Forbidden Terms

| Avoid | Use instead | Reason |
|-------|-------------|--------|
| `CW configuration` | `Seedable-Config` | Per verbatim §Important.3a |
| `git map` | `gitmap` | Per verbatim §Important.3b |
| `Master/Slave` | `Main/Worker` | Inclusive language, also matches verbatim |
| "CEO" for Md. Alim Ul Karim | **Chief Software Engineer** | Per memory `mem://project/author-attribution` |

---

*Glossary v1.1.0 — 2026-05-04 (added Quarantined, Draining, Seedable-Config link, apperror, Power Admin/PowerAdmin distinction; closes F-A-36/37/38/39/40)*
