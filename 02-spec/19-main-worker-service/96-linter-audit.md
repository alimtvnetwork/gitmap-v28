# 96 — Linter Audit

> **⚠ STALE — historical snapshot.** This file's "all green" claims pre-date Phase 13.2 hardening. The current authoritative blind-AI score lives in `audit/07-blind-ai-readiness-2026-05-06-v2.md` (and any successor). Do **not** treat this file as a green light to skip verification — re-run the linters before relying on it.

**Spec:** `19-main-worker-service`
**Audited:** 2026-05-04
**Audit Version:** 1.0.0

---

## 1. Scope

Manual + automated compliance pass against memory-defined SQL rules and the
two repo-resident spec linters.

| Rule / Linter | Source | Status |
|---|---|---|
| MISSING-DESC-001 | mem://sessions/2026-04-sql-linter-rules | ✅ PASS (manual) |
| DB-FREETEXT-001 | mem://sessions/2026-04-sql-linter-rules | ✅ PASS (manual) |
| Rule 10 — nullable, no DEFAULT | mem://index Core | ✅ PASS |
| Rule 11 — entity/ref `Description` | mem://index Core | ✅ PASS |
| Rule 12 — transactional `Notes` + `Comments` | mem://index Core | ✅ PASS |
| `check-spec-cross-links.py` | linter-scripts/ | ✅ PASS |
| `check-spec-folder-refs.py` | linter-scripts/ | ✅ PASS |

---

## 2. Per-Table Audit (`03-main-db-schema.md`)

| Table | Class | Required Cols | Present | Verdict |
|---|---|---|---|---|
| `WorkerNode` | entity | Description | ✅ | OK |
| `WorkerNodeStatus` | ref | Description | ✅ | OK |
| `WorkerNodeKind` | ref | Description | ✅ | OK |
| `Company` | entity | Description | ✅ | OK |
| `User` | entity | Description | ✅ | OK |
| `UserRole` | join | (exempt) | n/a | OK |
| `Role` | ref | Description | ✅ | OK |
| `WorkerVersion` | transactional | Notes + Comments | ✅ | OK |
| `WorkerSelectionEvent` | transactional | Notes + Comments | ✅ | OK |
| `WorkerSelectionStrategy` | ref | Description | ✅ | OK |

All `Description` / `Notes` / `Comments` columns are `TEXT NULL` with no
DEFAULT clause, satisfying Rule 10.

---

## 3. Free-Text Inline Strings (DB-FREETEXT-001)

No inline `Type` / `Status` / `Category` / `Kind` text columns. All such
discriminators are FK references to ref tables (`WorkerNodeStatus`,
`WorkerNodeKind`, `WorkerSelectionStrategy`, `Role`).

---

## 4. Naming Convention Audit

- ✅ All tables PascalCase.
- ✅ All columns PascalCase.
- ✅ PKs named `{TableName}Id`, all `INTEGER PRIMARY KEY AUTOINCREMENT`.
- ✅ No UUID columns.
- ✅ FKs named `{ReferencedTable}Id`.
- ✅ Indexes prefixed `IX_`.

---

## 5. Cross-Link Audit

`python3 linter-scripts/check-spec-cross-links.py` →
`OK All internal spec cross-references resolve.`

`python3 linter-scripts/check-spec-folder-refs.py` →
`✅ All spec/NN-name references resolve or are allowlisted.`

---

## 6. Outcome

**No findings.** Spec `19-main-worker-service` is compliant with all active
SQL and cross-link linter rules as of audit date.

---

*Linter audit v1.0.0 — 2026-05-04*
