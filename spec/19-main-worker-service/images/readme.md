# Images — Main↔Worker Service

Author's source mindmaps that seeded this spec. Kept verbatim as
authoritative intent; the spec text is the formal interpretation.

| # | File | Subject | Referenced In |
|---|------|---------|---------------|
| 01 | [`01-main-worker-topology.png`](./01-main-worker-topology.png) | Top-level `n-worker` split: `main` (recalltime.com, `/api/v1/.../company` CRUD, admin/users, poweradmin, push-updates → workers) and `self-update` (endpoint → JSON instructions → 10-part file). | `00-overview.md`, `01-architecture.md`, `09-self-update-pointer.md` |
| 02 | [`02-endpoint-service-worker-pattern.png`](./02-endpoint-service-worker-pattern.png) | EndPointService → WorkerPattern decomposition: "Give me new" (which is more free / endpoint of communication) and "Knowledge" (total workers, machine/websites, who is free). Underpins routing strategy decisions. | `04-worker-routing.md` |
| 03 | [`03-worker-subdomain-routing.png`](./03-worker-subdomain-routing.png) | Per-tenant subdomain layout: `alimkarim.com` → `w1..w10.alimkarim.com`, each exposing `api/v1/company` (POST/PUT/DELETE/GET) plus the per-worker `self-update` channel. | `03-main-db-schema.md`, `06-core-api-endpoints.md` |
| 04 | [`04-endpoint-service-full-overview.png`](./04-endpoint-service-full-overview.png) | Full combined EndPointService mindmap (images 02 + 03 in one canvas). The single best one-page summary. | `00-overview.md`, `01-architecture.md` |

## Conventions

- File names are `NN-kebab-case.png` where `NN` is sequential.
- Originals are PNG mindmap exports. Do not edit in place — add a new
  numbered file if a revision is needed.
- All Mermaid recreations live in [`../diagrams/`](../diagrams/). When
  the diagrams diverge from the mindmaps, the mindmaps are the source
  of intent and the diagrams must be updated.

---

*Images index v1.0.0 — 2026-05-04*
