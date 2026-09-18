# Rejog the Memory v1: Reliability & Failure-Chance Report

Generated: 2026-07-23 14:06 UTC
Scope: read-only synthesis of `.ai-memory/` and `02-spec/`. No code or spec content was modified.

## 1. Corpus inventory

### 1.1 `.ai-memory/` folders

| Folder | Notes |
|---|---|
| `.ai-memory/overview.md` | Project summary. **Stale version: says v3.1.0.** |
| `.ai-memory/strictly-avoid.md` | Hard prohibitions. Some phrased with negation that conflicts with the "positive logic" rule (rules about code, so acceptable). |
| `.ai-memory/prompt.md` | Index of onboarding prompts. |
| `.ai-memory/suggestions.md` | Legacy freeform log. References v3.12.1. |
| `.ai-memory/cicd-index.md` | Points to `cicd-issues/`. |
| `.ai-memory/issues/` | 3 files. Open: `01-ssh-repo-cloned-as-https.md`. |
| `.ai-memory/cicd-issues/` | 2 files. |
| `.ai-memory/solved-issues/` | 1 archive file (v2.74 to v2.76). |
| `.ai-memory/pending-issues/` | 1 file; flags missing unit tests for `task`, `env`, `install` since v2.49.0. |
| `.ai-memory/question-and-ambiguity/` | 11 decision logs from No-Questions Mode inferences. |
| `.ai-memory/audits/` | 2 pickers audits. |
| `.ai-memory/plans/pending/` | 3 pending plans (SSH-aware clone, reclone transport + VS Code open, CFR CG OS-aware guidelines). |
| `.ai-memory/plans/subtasks/` | Present, not surveyed in depth. |
| `.ai-memory/02-spec/commands/` | 6 command specs; possibly redundant with root `02-spec/01-app/`. |
| `01-prompts/` | 20 files. Files 03 to 20 all named `next-task.md` — high stale/duplicate risk. |
| `.ai-memory/memory/` | See 1.2. |

### 1.2 `.ai-memory/memory/` folders

| Folder | Files | Notes |
|---|---:|---|
| `index.md` | 101 lines | Core rules + memory index. **Stale current version: says v5.9.0.** |
| `constraints/` | 3 | CI/release untouchable, constants ownership, strictly prohibited. |
| `features/` | 60 | Ships-log per feature. Dense and current. |
| `plans/` | 9 | Multi-step plans, mostly complete. |
| `project/`, `tech/`, `style/`, `workflow/` | present | Not enumerated in depth. |
| `issues/`, `suggestions/` | present | `suggestions/01-suggestions.md` is the only entry. |
| Session logs | 3 | `01-replace-command-plan.md`, `02-v15-legacy-compat-audit.md`, `03-v3.12.1-session.md`. |

### 1.3 `02-spec/` root

| Folder | Purpose | Volume |
|---|---|---|
| `02-spec/01-app/` | gitmap CLI feature specs | 116+ files, heavy duplicate numbering |
| `02-spec/02-app-issues/` | Post-mortems | 34+ files, some duplicate prefixes |
| `02-spec/03-commit-in/` | Commit-in command | 5 files |
| `02-spec/03-general/` | Shared CLI/script patterns | 20 files |
| `02-spec/03-tasks/` | Ad-hoc task briefs | 2 files |
| `02-spec/04-generic-cli/` | Generic CLI blueprint | 37 files |
| `02-spec/05-coding-guidelines/` | Detailed coding rules | 31 files |
| `02-spec/06-design-system/` | Docs site visual tokens | ~5 files |
| `02-spec/07-generic-release/` | Release pipeline blueprint | ~10 files |
| `02-spec/08-generic-update/` | Self-update blueprint | ~10 files |
| `02-spec/08-json-schemas/` | JSON output schemas | 2 files |
| `02-spec/09-pipeline/` | CI pipeline blueprint | 11 files |
| `02-spec/12-consolidated-guidelines/` | Refined policy set | 19 files |
| `02-spec/15-research/` | Advanced git research | 1 file |

### 1.4 Duplicate numeric prefixes (02-spec/01-app + 02-app-issues)

Duplicates observed at prefixes: `26-`, `27-` (three files), `89-`, `90-` (three files), `95-`, `96-`, `100-`, `108-`, `109-`, `110-` (three files), `111-` (three files). Full list in the sub-agent report (see corpus survey). Consequence: unstable ordering, ambiguous cross-references, and grep-first workflows break.

### 1.5 Version drift

| Source | Version |
|---|---|
| `cli/constants/constants.go` (source of truth) | v6.79.0 |
| `.ai-memory/memory/index.md` | v5.9.0 |
| `.ai-memory/overview.md` | v3.1.0 |
| `.ai-memory/suggestions.md` | v3.12.1 |

Every AI onboarding prompt reads `overview.md` and `index.md`, so this drift ships bad context on turn one.

## 2. Success probability by tier

Assumptions common to all tiers:
- The AI honors "No-Questions Mode" and logs ambiguities to `.ai-memory/question-and-ambiguity/`.
- The AI enforces the Core rules from `mem://index.md` (constants centralization, positive logic, file/function length, zero-swallow errors, DB rules, SQLite `SetMaxOpenConns(1)`).
- The release pipeline and `.gitmap/release/` remain untouched.

| Tier | Example specs | Estimated first-pass success | Rationale |
|---|---|---:|---|
| Simple: single-file spec, isolated CLI flag | `110-clone-ssh-flag`, `19-list-versions`, `27-bookmarks` | 80 to 90% | Small surface, mature patterns, strong constants discipline, plenty of prior art in `features/`. |
| Medium: multi-file feature | `104-clone-multi`, `111-cn-folder-arg`, `114-committransfer-idempotence` | 60 to 75% | Cross-package coordination + tests. Duplicate spec numbering and outdated `index.md` version raise the risk of misapplied precedents. |
| Complex agentic: self-update, cross-platform install, history rewrite | `108-cross-platform-install-update`, `110-update-remote-installer`, `115-v6-migration`, `02-spec/15-research/*` | 35 to 55% | Windows locking, PATH management across three OSes, GitHub Actions coupling. Rules exist in `02-spec/08-generic-update/` but external state (registry, wrapper shims, handoff sentinel) is invisible to a fresh AI. |
| End-to-end: v6 migration, docs-site parity, bulk visibility, commit-in | `115-v6-migration`, `35-docs-site`, `116-bulk-visibility`, `02-spec/03-commit-in/` | 25 to 45% | Multi-repo effects, DB migrations, non-idempotent operations. `commit-in` is spec-only and gated on user typing `next`; easy to violate. |

## 3. Failure map

| Area | Likely failure | Why | Symptom |
|---|---|---|---|
| Onboarding | AI trusts stale version in `index.md`/`overview.md` | Drift documented in 1.5 | Wrong feature gates, "already shipped" features re-planned. |
| Spec lookup | Numeric-prefix collisions in `02-spec/01-app/` | See 1.4 | AI cites `108-*` for install and lands on cross-platform doc when it wanted install-quick-auto-source, or vice versa. |
| Constants | Local constants block added despite "all CLI IDs in `constants_cli.go`" | Rule split across `spec/05` and `spec/12` (relaxed to domain-local) | Duplicate command IDs, silent alias shadows. |
| Function length | `02-spec/05/01` says 25, `02-spec/05/02` and `02-spec/12/02` say 15 | Contradiction | Reviewer rejects PR / auto-lint fails. |
| Release pipeline | AI edits `.github/workflows/release.yml` when artifacts missing | Root cause is upstream build failure, not workflow | Broken releases; explicit prohibition in `mem://constraints/ci-release-pipeline-untouchable`. |
| Self-update on Windows | AI removes handoff copy shim | Documented in `02-spec/08-generic-update/05-handoff-mechanism.md` | Update loop / locked binary. |
| Windows install | PATH snippet marker mistaken for wrapper marker | Documented in `features/command-wrapper-marker-separation.md` | Users think install worked; wrapper not loaded. |
| Hosted docs fallback | Someone re-adds `os.Exit(1)` in `hd` when docs missing | Contract pinned by tests, but easy to break in refactor | `gitmap hd` crashes offline. |
| JSON stability | Field reordering after Go upgrade | `stablejson` package required, per index Core | Snapshot diffs across environments. |
| Suggestions workflow | New suggestions dumped into legacy `01-suggestions.md` | No template or naming rule enforced | Suggestions vanish or get overwritten. |
| Prompt hygiene | New AI reads all 20 identically named `next-task.md` files | Duplicate naming | Wrong or superseded task picked. |

## 4. Corrective actions (prioritized)

Priority = P0 (blocker), P1 (should fix before large work), P2 (hygiene).

| # | Priority | Action | Where | Reliability gain |
|---|---|---|---|---|
| 1 | P0 | Sync current version in `.ai-memory/overview.md` and `.ai-memory/memory/index.md` to the value in `cli/constants/constants.go` and document the source of truth. | 2 files | Removes onboarding lie; unblocks tier 3 and 4 correctly. |
| 2 | P0 | Add a `02-spec/01-app/README.md` mapping numeric prefixes to canonical files, and renumber or namespace duplicates (26, 27, 89, 90, 95, 96, 100, 108, 109, 110, 111). | `02-spec/01-app/`, `02-spec/02-app-issues/` | Removes ambiguity in every future spec reference. |
| 3 | P0 | Resolve function-length contradiction between `02-spec/05-coding-guidelines/01-*` (25) and `02-spec/05/02-go-code-style.md` + `02-spec/12/02-go-code-style.md` (15). Publish one number and delete the other. | 2 files | Stops PR/lint churn. |
| 4 | P1 | Split Core rules in `mem://index.md` into "Permanent" and "Session flags" so `NO-QUESTIONS MODE` and its 40-task budget cannot leak into future sessions. | `index.md` | Prevents future AI sessions from silently inheriting a mode. |
| 5 | P1 | Reconcile `02-spec/05-coding-guidelines/` constants rule ("all strings in one package") with `02-spec/12-consolidated-guidelines/` domain-local ownership. Pick one and cross-link. | 2 folders | Consistent constants placement. |
| 6 | P1 | Deduplicate `.ai-memory/02-spec/commands/` vs `02-spec/01-app/`. Keep root `02-spec/` as the single source; leave a stub redirect in `.ai-memory/spec/`. | 2 folders | Fewer collided sources for the AI. |
| 7 | P1 | Compress `01-prompts/03-*` through `20-next-task.md` to one canonical prompt + an archive folder; the user preference already forbids per-invocation mirrors. | `01-prompts/` | Aligns with user memory rule; reduces noise on read. |
| 8 | P1 | Formalize suggestions workflow (see Deliverable 2) so no new work uses the legacy log. | `.ai-memory/memory/suggestions/` | Deterministic backlog capture. |
| 9 | P2 | Document contract tests that guard high-risk invariants (`TestHostedDocsFallbackContract`, wrapper marker separation, stablejson encoding) in one page so refactors see them. | new file under `02-spec/01-app/` | Guards against regressions in complex tier. |
| 10 | P2 | Note open bug `.ai-memory/issues/01-ssh-repo-cloned-as-https.md` explicitly in the root plan so it is not lost. | `plan.md` | Prevents shipped SSH bug from being forgotten. |
| 11 | P2 | Add missing unit tests for `task`, `env`, `install` (flagged in `pending-issues/`). | `cli/cmd/*_test.go` | Restores CI baseline. |
| 12 | P2 | Delete `02-spec/03-general/10-strictly-prohibited.md` duplication with `.ai-memory/strictly-avoid.md` OR make one a link to the other. | 2 files | Single source for prohibitions. |

## 5. Readiness decision

**Not ready to hand off to another AI for large or complex work.** Small isolated tickets (tier 1) are safe. Before starting anything in tier 3 or 4, fix items 1, 2, 3 (all P0). Items 4 to 7 should follow immediately after.

Minimum unblock set:
- Item 1: version sync in `overview.md` and `index.md`.
- Item 2: prefix map + rename duplicates in `02-spec/01-app/` and `02-spec/02-app-issues/`.
- Item 3: single function-length number across `02-spec/05/` and `02-spec/12/`.

Once those three land, the corpus is ready for medium tier work; complex and end-to-end tier work still needs items 4 through 9.
