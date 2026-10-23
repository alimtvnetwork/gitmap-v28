# Muse Memory IO and GitMap Integration — Execution Plan

- Request slug: `muse-memory-io-and-gitmap-integration`
- Status: `ACTIVE` | Phase: 2 | Wave: 0 / 2 | Step: ~25 / 120
- Machine ID: pending
- Branch: `main` | Tree at start: DIRTY (untouched files below — another in-flight run; DO NOT TOUCH)

## User Request (Verbatim, abridged)

Build "Muse Memory IO": (1) compile everything done with Muse in the last 7 days (2026-10-10→17: new conversations, implementations, new repos) into sequenced digest files; (2) store under a Muse/ folder in the repo cache with a root Muse/README.md plus numbered reader files that reference each other, so any future Muse can read them, rebuild chat context, and recreate chats from the last 5 days from seed files; (3) establish the "Muse memory IO / say READ" convention: in ANY chat the user says "Muse memory IO read" and Muse pulls repo-secrets, diffs what's new since last read, and gets up to speed (asking ONLY for the email password when go-email-reader setup is needed — never anything else secret); (4) update coding-guidelines-v24 Muse prompts (01-prompts/27-muse-prompts/) + skills so private items (which repos to read etc.) resolve from the repo cache, not hardcoded; (5) document go-email-reader's email info location, sending, and default emailer; (6) implement `gitmap llm train muse [--days N]` in gitmap-v28: auto-clone repo-secrets into the work-folder root if absent else pull, diff the last N days vs previous snapshot, print summary to terminal plus raw gitmap/AI traversal commands as a file-by-file checklist for LLM self-training; (7) register the new command in gitmap's new-command docs WITH changelog entries (changelog explicitly requested — allowed despite zero-release rule; NO version bumps); (8) commit+push every touched repo.

## Research Findings (synthesized)

- R-A (gitmap-v28): `train muse` is a positional sub-target of `llm train`, NOT a top-level command. Implementation: strip positional `muse` before `parseTrainFlags()` in `cli/cmd/llm/llm_train.go`; add `TrainTarget string` + `DaysN int` to `cli/cmd/llm/llm_types.go:TrainOptions`; new `cli/cmd/llm/llm_train_muse.go` (`runTrainMuse`) using store special-repo "rs" entry + os/exec git (circular-import trap: `cli/cmd/llm` cannot import `cli/cmd`; reimplement clone-or-pull inside package `llm`); usage lines in `printTrainUsage()` + `cli/cmdllm/llm_help_menu.go`; docs in `cli/helpdoc/llm.md` + `cli/helpdoc/examples_dev.go` (`reg["llm-train"]`); changelog under `## [Unreleased]` (no new version section, no bumps); next spec number 277.
- R-B (coding-guidelines-v24): `01-prompts/27-muse-prompts/` has 6 files; NONE reference repo-secrets/repo-cache or the "say READ" convention. Sync script `03-ai-scripts/50-sync-prompts-skills-to-targets.py` (prompts+skills only, NEVER spec; run after edits).
- R-B (go-email-reader): store = `$XDG_CONFIG_HOME/go-email-reader/root.db` (path constant `internal/store/store.go` + `internal/crypto/rsa.go:30` `ConfigDir()`); RSA keypair `<configDir>/id_rsa`(0600)/`id_rsa.pub`(0644), passwords as RSA-OAEP-SHA256 base64 ciphertext; default emailer = `cmdDefault` in `cmd_account.go:181` (`is_default` INTEGER); docs live in `README.md`; standing repo rule: save a dated note under `04-conversation/` BEFORE the task.
- R-B (repo-secrets): cloned to workspace root sibling dir (repo: `alimtvnetwork/repo-secrets`, private); 16 dirs + 8 files at top; NO `Muse/` yet. Digest sources confirmed: `~/memory/2026-10-10.md`(575) … `2026-10-17.md`(327), `~/MEMORY.md`(69).
- Digest secrets rule: ZERO secret values in digest files — pointers only.

## Task Table

| Task-ID | Subtask | Owner | Owned files | Status | Evidence |
|---|---|---|---|---|---|
| Task-01 | W-1: gitmap-v28 spec files | Worker 01 | `02-spec/21-app/277-muse-memory-io-and-gitmap-integration/…`; `.ai-memory/plans/subtasks/muse-memory-io-and-gitmap-integration/NN-*.md` | DONE | spec 277 3 files + 2 subtask files on disk |
| Task-02 | W-1: `llm train muse` implementation | Worker 01 | `cli/cmd/llm/llm_train_muse.go` (new), `llm_train.go`, `llm_types.go`, `llm_train_test.go`, `cli/cmdllm/llm_help_menu.go` | DONE | autofixer exit 0 (10 files), gofmt clean |
| Task-03 | W-1: new-command docs + changelog | Worker 01 | `cli/helpdoc/llm.md`, `cli/helpdoc/examples_dev.go`, `changelog.md`, `cli/cmdllm/llmdocsgroups.go` | DONE | [Unreleased] entries, no version bump |
| Task-04 | W-2: repo-secrets Muse/ digest | Worker 02 | `Muse/README.md`, `Muse/0N-*.md`, `Muse/chat-seeds/*.md`, `Muse/last-read.md` | DONE | 16 files on disk, pointers-only |
| Task-05 | W-2: coding-guidelines-v24 prompt + skill updates | Worker 02 | `01-prompts/27-muse-prompts/*`, `.agents/skills/muse-master-prompt/skill.md`, `.cursor/skills/muse-master-prompt/skill.md` | DONE | sync script exit 0 (7 synced, 613 skipped) |
| Task-06 | W-2: sync script run | Worker 02 | (targets under `~/workspace/custom-skills`) | DONE | exit 0, fail-safe gate proven |
| Task-07 | W-2: go-email-reader email docs + conversation note | Worker 02 | `04-conversation/2026-10-17-1045-muse-memory-io.md`, `README.md` | DONE | note written first, README extended |
| Task-08 | Lead: commits + pushes per repo | Lead | (gitmap cpf/cpb + git ls-remote verify) | IN PROGRESS | - |

## Conflicts

- None. Parent directive folded the spec step into Phase-2 worker boxes; the lead's unified plan (this file) is written by the lead per the parent's explicit order.

## Assumptions

- `--days N` default 7 per the user brief. Digest window 2026-10-10→17; chat seeds 2026-10-12→17 (last 5 days).
- `gitmap llm train muse` treats a missing local repo-secrets as clone (gh, private repo, fail-fast on missing auth, never hang on credential prompts); present → pull; missing `Muse/` dir → clear message, exit 0.
- Changelog entries go under `## [Unreleased]`; no version bumps, no `version.json` edits.

## Race note (2026-10-17)

Program 275's lint commit `942a674a6` ("revert rogue worker files") deleted in-flight files of this program mid-run (plan file, index rows, parts of worker edits). Workers re-applied and re-verified; files re-created and committed here. If reverts recur, the two programs' file boxes must be serialized by the parent.

## Stage list

gitmap-v28: `.ai-memory/plans/muse-memory-io-and-gitmap-integration.md` (lead), `02-spec/21-app/277-…`, `.ai-memory/plans/subtasks/muse-memory-io-and-gitmap-integration/…`, `cli/cmd/llm/llm_train_muse.go`, `cli/cmd/llm/llm_train.go`, `cli/cmd/llm/llm_types.go`, `cli/cmd/llm/llm_train_test.go`, `cli/cmdllm/llm_help_menu.go`, `cli/helpdoc/llm.md`, `cli/helpdoc/examples_dev.go`, `changelog.md`, indexes (lead).
repo-secrets: `Muse/README.md`, `Muse/0N-*.md`, `Muse/chat-seeds/*.md`, `Muse/last-read.md`.
coding-guidelines-v24: `01-prompts/27-muse-prompts/*`, skills, sync targets.
go-email-reader: `04-conversation/2026-10-17-1045-muse-memory-io.md`, `README.md`.

## Follow-ups (2026-10-17, same day — all DONE, committed+pushed, ls-remote-verified)

- F1: numbered DECISION TRACE in terminal output (each step carries its why) + why-annotations on checklist steps; bugfix: existing repo-secrets dir at work-folder root reused without `rs` record.
- F2: last-read watermark moved OUT of MD into gitmap itself (vault HEAD SHA + UTC timestamp JSON in `store.BinaryDataDir()`, 0600); SINCE LAST READ section; `Muse/last-read.md` deleted.
- F3: digest folder renamed `Muse/` → `12-muse-ai/` (repo NN-sequence convention; user couldn't find it).
- F4: email config documented from `repo-secrets/11-go-email-reader/` (password-only rule); new-muse onboarding instruction written (now 9-step report-back checklist in `repo-cache/12-muse-ai/00-new-muse-onboarding.md`).
- F5: digest MOVED to `alimtvnetwork/repo-cache` ("GitMap distributed shared cache and AI memory vault", created 2026-10-17); deleted from repo-secrets (secrets stay). `llm train muse` reconciled with another worker's landed dual-vault refactor: kept their architecture, flipped priority to repo-cache-first (special-repo key `rc`), repo-secrets as fallback. Rebase conflict resolved without discarding landed work.
