# Subtask Plan 02: Implement `gitmap llm train muse [--days N]` + docs + changelog

- **Version:** 1.0.0
- **Updated:** 2026-10-17
- **Owner:** Worker 01
- **Status:** DONE
- **Subtask Slug:** `02-implementation`
- **Parent Task:** `muse-memory-io-and-gitmap-integration` (Spec 277)

---

## Goal

Implement `gitmap llm train muse [--days N]` as a positional sub-target of
`llm train` (NOT a top-level command), wire the `train muse` help surfaces,
and ship the new-command docs + changelog entries — all inside the Owned
Files list, with zero builds/tests run locally (R1; CI verifies).

## Inputs

- Spec 277 (`02-spec/21-app/277-muse-memory-io-and-gitmap-integration/`).
- `cli/cmd/llm/llm_train.go`, `cli/cmd/llm/llm_types.go`,
  `cli/cmd/llm/llm_train_test.go` (existing train code + test style).
- `cli/cmd/llm/llm_heal.go` imports (safe sibling packages:
  `cli/apperror`, `cli/cmdautofix`, `cli/cmdfix`, `cli/cmdremediation`).
- `cli/cmd/special_repos_ops.go:335` `runGitInSpecialRepo` os/exec pattern.
- Store: `OpenSpecialReposSplitDB().GetSpecialRepo("rs")` — short key `rs`
  verified in `cli/store/special_repos_split_db.go`.

## Steps

1. `cli/cmd/llm/llm_types.go`: add `TrainTarget string` and `DaysN int`
   fields to `TrainOptions`; add `DefaultTrainMuseDays = 7` const.
2. `cli/cmd/llm/llm_train.go`:
   - Add `stripTrainTarget(args)` — strips the FIRST bare `muse` token that is
     not a flag and not a flag value; returns `(target, rest)`.
   - `RunTrain` calls `stripTrainTarget` BEFORE `parseTrainFlags`, sets
     `opts.TrainTarget`, routes `TrainTarget == "muse"` to `runTrainMuse`;
     `muse --help` prints muse-specific usage via new `printTrainMuseUsage`.
   - `parseTrainFlags`: add `--days` int flag (default 7); clamp `<= 0` to 7.
   - `printTrainUsage`: add `train muse [--days N]` lines + `--days` entry.
3. NEW `cli/cmd/llm/llm_train_muse.go` (`runTrainMuse` + helpers):
   - Resolve via store `GetSpecialRepo("rs")`; absent → `gh repo clone
     https://github.com/alimtvnetwork/repo-secrets.git` with non-interactive
     stdin + `gh auth` pre-check; present → `git pull --ff-only` (skip when no
     `origin` remote); `git log --since` N-day diff; terminal summary + AI
     traversal checklist; missing `Muse/` dir → guidance, exit 0; all errors
     as `*apperror.AppError`. No `cli/cmd` import.
4. `cli/cmd/llm/llm_train_test.go`: stripping tests, `--days` default/value
   tests, missing-digest test (temp git repo fixture + resolver override).
5. `cli/cmdllm/llm_help_menu.go`: `train muse [--days N]` usage lines,
   Commands-section entry, `--days <N>` footer flag.
6. `cli/cmdllm/llmdocsgroups.go`: Utilities-group `llm train muse` matrix
   entry.
7. `cli/helpdoc/llm.md`: `train muse` section (invocations + sample output)
   + cheat-sheet row.
8. `cli/helpdoc/examples_dev.go`: extend `reg["llm-train"]` with
   `gitmap llm train muse --days 7` Invocation + ExpectedOutput (3-part).
9. `changelog.md`: `## [Unreleased]` entries only.
10. Verify: `python 03-ai-scripts/05-guideline-autofixer.py cli/cmd/llm
    --check-only --ext .go` (must scan > 0 files, exit 0); `gofmt -l` on
    touched files if gofmt exists.

## Acceptance criteria

- [x] `TrainTarget` / `DaysN` in `TrainOptions`; `DefaultTrainMuseDays = 7`.
- [x] `RunTrain` pre-scans and strips `muse` before flag parsing; routes to
      `runTrainMuse`; existing `llm train` behavior unchanged.
- [x] `--days` parsed with default 7 and `<= 0` clamp.
- [x] `runTrainMuse` implements clone-or-pull, N-day diff, terminal summary,
      AI traversal checklist, missing-digest guidance (exit 0), all
      `*apperror.AppError` errors, no `cli/cmd` import.
- [x] New tests added in existing style (stripping, `--days`, missing
      digest).
- [x] Help menu, usage text, docs matrix, helpdoc, examples, changelog
      updated as specced.
- [x] Autofixer check passes on `cli/cmd/llm` with files scanned > 0.

## Out of scope

- Digest authoring (another worker).
- Version bumps, releases, top-level command plumbing.
- Touching the other-run dirty-file list or shared indexes.
- Running builds or tests locally (R1).
