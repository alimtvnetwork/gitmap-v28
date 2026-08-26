# 24 — Threat Model (Reserved Stub)

**Spec:** `19-main-worker-service`
**Version:** 0.1.0
**Updated:** 2026-05-06
**Status:** 🟡 **RESERVED — placeholder only.** Do **not** implement against this file.

---

## 1. Purpose of this stub

This file **reserves spec slot 24** within `19-main-worker-service/` so a future
spec author cannot accidentally claim slot 24 for an unrelated topic.

Slot 24 is reserved for a forthcoming **Threat Model** document that will
underpin the v2.0 reopen of two formally-deferred open questions:

| Origin OQ | Source disposition | What the threat model must cover |
|---|---|---|
| **OQ-12-1** — Separate `RefreshToken` cookie + rotation | [`12-jwt-delivery-contract.md` §11.3](./12-jwt-delivery-contract.md) — *Future-work catalogue* | Token-theft window, sign-out-everywhere, refresh-replay attack tree, cross-origin isolation assumptions. |
| **OQ-23-1** / **OQ-23-2** — Snapshot dedup pyramid + partial-tenant restore | [`23-snapshot-storage-and-restore.md` §14.4](./23-snapshot-storage-and-restore.md) — *Future-work catalogue* | Snapshot-chain integrity, per-tenant blast radius, re-seal/key-rotation interaction, audit-bypass risk of partial restore. |

Both deferred OQs declared `24-threat-model.md` as a **shared prerequisite**
for any v2.0 reopen. Until this file becomes a real spec, neither OQ may be
reopened.

---

## 2. What is **not** in this file (CODE RED)

Per the formalised dispositions in §11.3 / §14.4, this stub MUST NOT contain
any of the following — adding them here is itself a CODE RED violation:

- ❌ Any STRIDE / attack-tree content (premature without the v2.0 trigger).
- ❌ Any new error codes (`WORKER-100-04+`, `WORKER-940-05+`, `MAIN-830-04+`).
- ❌ Any schema sketches (no `RefreshToken` cookie shape, no
  `BasedOnSnapshotCatalogId` FK sketch, no `RestoreScope` column).
- ❌ Any acceptance criteria — `97-acceptance-criteria.md` must remain free
  of rows that depend on this file.
- ❌ Any "for future use" allocations of any kind.
- ❌ Any v1.0 implementation guidance — slot 24 is **v2.0-only**.

Forbidden patterns above mirror:
- [`12-jwt-delivery-contract.md` §11.1 / §11.3](./12-jwt-delivery-contract.md)
- [`23-snapshot-storage-and-restore.md` §14.1 / §14.2 / §14.4](./23-snapshot-storage-and-restore.md)

---

## 3. Promotion criteria (when this stub becomes a real spec)

This file may be promoted from **stub** (v0.x) to **draft** (v1.0+) only when
**all** of the following are true:

1. At least one v2.0 reopen trigger from §11.3 OR §14.4 has fired (logged with
   evidence in `98-changelog.md`).
2. A `MainWorker.Auth.MaxSessionLifetimeSeconds` cap has been added to
   `15-tunable-constants.md` (prerequisite for OQ-12-1).
3. A dedicated error-code family (`WORKER-100-04+` for auth or
   `WORKER-940-05+` / `MAIN-830-04+` for backup) has been formally allocated in
   `13-error-codes.md` — **not** before.
4. A corresponding row exists in `97-acceptance-criteria.md` covering the new
   threat-model-derived behaviour.

Until **all four** are satisfied, this file stays at v0.1.0 and contains
only this reservation notice.

---

## 4. Cross-references

| Reference | Location |
|---|---|
| OQ-12-1 / OQ-12-2 dispositions | [`12-jwt-delivery-contract.md` §11](./12-jwt-delivery-contract.md) |
| OQ-23-1 / OQ-23-2 dispositions | [`23-snapshot-storage-and-restore.md` §14](./23-snapshot-storage-and-restore.md) |
| Tunable constants (cap prerequisite) | [`15-tunable-constants.md`](./15-tunable-constants.md) |
| Error codes (allocation prerequisite) | [`13-error-codes.md`](./13-error-codes.md) |
| Acceptance criteria | [`97-acceptance-criteria.md`](./97-acceptance-criteria.md) |
| Changelog | [`98-changelog.md`](./98-changelog.md) |

---

*Stub v0.1.0 — 2026-05-06 — slot reservation only, no normative content.*
