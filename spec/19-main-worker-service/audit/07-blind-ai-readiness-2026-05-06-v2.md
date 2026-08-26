# 07 — Blind-AI Implementation Readiness Audit (Independent Re-Run)

> **Spec under audit:** `spec/19-main-worker-service/` (29 numbered files, ~370 KB / ~7,000 lines, excluding `98-changelog.md`).
> **Audit date:** 2026-05-06
> **Auditor model:** `google/gemini-2.5-pro` via Lovable AI Gateway, fed the **full unabridged corpus** (no per-file truncation).
> **Persona simulated:** the dumbest plausible AI coder — never asks, picks the first matching rule, treats MUST as hard / SHOULD as ignored, cannot reconcile contradictions, cannot infer.
> **Why a re-run:** The prior audit (`06-blind-ai-readiness-2026-05-06.md`) scored **92/100**. That run scored the spec as the *author* of the changelog believes it to be. This run scores it as a hostile literal reader actually sees it.

---

## TL;DR

| Dimension                  | Score (this run) | Prior audit (06-…) |
| -------------------------- | ---------------: | -----------------: |
| Completeness               |               30 |                 94 |
| Determinism                |               10 |                 91 |
| Consistency                |                5 |                 96 |
| Testability                |               45 |                  — |
| Blind-buildability         |               20 |                 88 |
| **Overall blind-AI score** |       **22 / 100 (F)** |       92 / 100 (A−) |
| Estimated dumb-AI implements correctly | **~20 %** | ~92 % |
| Estimated dumb-AI fails or builds wrong thing | **~80 %** | ~8 % |

**Headline:** the spec is *conceptually* solid but operationally hostile to a literal AI. The previous audit measured intent; this one measures the artifact. A blind AI handed only this folder will build a working backup subsystem and a broken core.

---

## 1. Why the score dropped 70 points

The prior audit:

- accepted "follow-up task" markers as if they were already applied;
- treated reconciliation prose in `11-split-db-tier-reconciliation.md` as if it had rewritten the files it points at;
- counted the `99-consistency-report.md` and `96-linter-audit.md` "all green" claims as evidence of consistency;
- did not stress-test the same artifact (e.g. the `WorkerNode` table) appearing in three files.

This audit ran the same files through a literal reader. Every contradiction the prior audit catalogued as "to be fixed later" became a hard implementer fork.

---

## 2. What the spec actually gets right (the 20 % that ships)

1. **Pointer-only files** (`09-self-update-pointer.md`, `24-threat-model.md`) are unambiguous negative constraints — a blind AI will correctly *not* implement them.
2. **Centralised registries exist** for error codes (`13-`) and tunables (`15-`). The *concept* is buildable; the execution has bugs (see §3.5).
3. **Explicit forbidden-implementation lists** in `12-jwt-delivery-contract.md` and `23-snapshot-storage-and-restore.md` are best-in-class for blind AIs.
4. **Backup / snapshot / restore subsystem** (files 18–23) is the most determinism-friendly part of the spec; a literal AI will get most of it right because the contracts are stated once and not contradicted elsewhere.
5. **Top-level Main/Worker mental model** (`00-`, `01-`) is conceptually clear at the K8s analogy level.

---

## 3. The 80 % that fails — itemised

### 3.1 Critical contradictions (≈65 % of all failures)

| # | File(s) | Failure mode for a dumb AI | One-line fix |
|---|---------|----------------------------|--------------|
| C-1 | `10-§8` vs `03-§2.1` — `WorkerNode` table | Two schemas with different column names (`WorkerNodeDisplayName` vs `WorkerNodeTitle`) and different normalisation (inline vs separate `WorkerVersion` table). AI builds the first one it reads. | Delete the `WorkerNode` definition from `10-` and reference `03-` as SoT. |
| C-2 | `04-§7.2` vs `06-` vs `10-` vs `18-` — `POST /API/V1/Workers/Register` | Four endpoint catalogues, each with a different request body. AI implements whichever file it parses first. | Make `06-core-api-endpoints.md` the only endpoint catalogue; delete the others. |
| C-3 | `06-§5.1` `EndpointAuthSetting.UpdatedByUserId REFERENCES User(UserId)` | The `User` table was deleted in `03-`; the FK is unresolvable and SQL DDL fails at create-time. | Repoint the FK to `UserDirectory(UserDirectoryId)`. |
| C-4 | `15-§2` (prose keys) vs `15-§4` (JSON seed keys) | Prose says `MainWorker.RateLimit.AuthEndpointsPerMinutePerIp`; seed JSON says `RateAuthPerMinutePerIp`. Lookup at runtime returns nothing → AI silently uses hard-coded defaults. | Make the two key sets byte-identical. |
| C-5 | `13-§6` vs `08-` — JSON error envelope | `13-` adds `ErrorCodeFlat` and `ErrorName`; `08-` does not. AI emits whichever envelope it saw first; clients break either way. | Define the envelope once in `08-` and have `13-` link to it. |
| C-6 | `11-§3` & `§8` — "follow-up tasks" still listed | The reconciliation file *describes* what should change in other files instead of *applying* the change. Literal AI reads the unchanged target files and implements the stale rule. | Apply the follow-up edits inline; replace `11-` with a "done" stub. |

### 3.2 Information not present in the corpus (≈20 %)

| # | File | Missing reference | Failure mode |
|---|------|-------------------|--------------|
| M-1 | `00-§7` | `spec/03-error-manage/`, other sibling specs | Inherited rules can't be read → AI ignores them. |
| M-2 | `01-§5` | `mem://architecture/caching-policy` | URI unresolvable from inside the spec corpus → cache TTLs undefined. |
| M-3 | `04-§2` | TTLs `15 min` / `60 s` hard-coded *and* missing from `15-` | Violates the "single source of truth" rule the spec itself sets. |
| M-4 | `00-§6` Document Map | Lists files 00–09 only; omits 10–24 | A literal AI scoping work from the map will never read 60 % of the spec. |
| M-5 | `03-§2.1` | "Validated by trigger" without trigger SQL | Constraint can't be created. |
| M-6 | `05-§2.1` | "Synthetic null-Worker handler matching wall-clock budget" with no algorithm | AI cannot invent timing-attack mitigation. |

### 3.3 Misleading meta-documents (≈10 %)

- `96-linter-audit.md` and `99-consistency-report.md` claim the spec is clean. A trusting AI uses them as a green light and skips verification. A skeptical AI loses trust in the entire corpus.
- `06-blind-ai-readiness-2026-05-06.md` (the prior audit) scored 92/100 against the same artifact this run scores 22. Treat that file as historical, not current.

### 3.4 Ambiguity / required invention (≈5 %)

- `14-§3.1` mandates `@Role.PowerAdmin` syntax in seed JSON but the "seeder feature" that resolves it is undefined → AI inserts the literal string `@Role.PowerAdmin`.
- `07-§3` `AccessItem` catalog is given as PHP `enum` syntax in an otherwise stack-agnostic document.
- `13-§1` slot-overflow rule depends on undocumented sub-ranges (e.g. `2115x = DB`).

### 3.5 Determinism forks the prior audit caught — still open

- `19-§1` permits two `SyncOp` shapes (inline column / sidecar). Prose now says Shape A is the v1.0 default (per Phase 13.1) but the alternative is not deleted from §1, so a literal AI still sees a fork.
- `22-§Stage 4` WAL pragma decision: the new §6.4 forbids per-envelope WAL but Stage 4 prose still reads as a choice.

---

## 4. Failure mix (where the 80 % goes)

```
65 %  Inconsistent definitions across files
20 %  Missing / external references
10 %  Misleading meta-documents (stale audits, linter all-green claims)
 5 %  Ambiguity / required invention
```

---

## 5. Why a literal AI cannot just "pick one"

The spec repeatedly says: *"if A and B disagree, the later/more-specific file wins."* That rule is itself unenforceable for a dumb AI because:

1. The files are read in unspecified order (filesystem listing is not a contract).
2. There is no machine-readable manifest declaring which file is authoritative for which artifact.
3. Several "more-specific" files (`10-`, `13-`, `18-`) repeat shapes the "general" files (`03-`, `08-`) own — the dumb AI cannot tell which side is the override.

---

## 6. The minimum edit set to reach ≥ 90 / 100

If the goal is to lift the dumb-AI score above 90 without rewriting the spec, the *narrow* fix list is:

1. **Delete-and-reference**: remove duplicate `WorkerNode` schema from `10-`; remove duplicate endpoint catalogues from `04-`, `10-`, `18-`; remove duplicate error envelope from `13-`. Each deleted block becomes a one-line `See: <file>#<section>` link.
2. **Apply, don't list**: execute the §3 / §8 follow-ups inside `11-split-db-tier-reconciliation.md` against their target files, then replace `11-` with a 5-line "applied on YYYY-MM-DD" stub.
3. **Reconcile keys**: make `15-§2` prose keys byte-identical to `15-§4` JSON seed keys. Add the `04-` cache TTLs to `15-`.
4. **Fix the FK**: change `EndpointAuthSetting.UpdatedByUserId` to reference `UserDirectory`.
5. **Inline the externals**: copy the inherited rules from `spec/03-error-manage/` into `08-`/`13-`; define the caching policy locally in `15-` (drop the `mem://` URI from `01-`).
6. **Update `00-§6`** Document Map to list every file 00–24.
7. **Mark stale meta-docs**: prepend a `STALE — see audit/07-…` banner to `06-…audit.md`, `96-`, and `99-`.

These seven edits remove the entire "Inconsistent Definitions" (65 %) and "Misleading Meta-Documents" (10 %) categories, lifting the practical pass rate from ~20 % to ~95 %.

---

## 7. Verdict

**Current state, dumb-AI hand-off:** **22 / 100 — F. ~80 % of the spec will be implemented incorrectly or not at all.**

**Current state, competent-AI hand-off** (which is what `06-…audit.md` measured): around **88–92 %** is buildable because a competent AI silently picks the right file in each duplicate pair and applies the §11 follow-ups by hand.

**Recommendation:** treat the seven edits in §6 as the v5.22.0 hardening pass. They are mechanical, ≤ 15 line-touches each, and close the entire critical category. After they land, re-run `bun run spec:audit` and confirm the score crosses 90 for `spec/19-main-worker-service`.

---

## 8. Reproduction

```bash

# Re-run this audit:

node /tmp/audit19b.mjs   # or rebuild from scripts/spec-audit/run-audit.mjs

# Inputs: every *.md under spec/19-main-worker-service/ except 98-changelog.md

# Model:  google/gemini-2.5-pro via Lovable AI Gateway

# Mode:   full corpus (no per-file truncation), strict JSON response_format

```

---

*Blind-AI Audit v2 — independent re-run — 2026-05-06.*
