# Spec 277 — Muse Memory IO and GitMap Integration: scope and plan

**Version:** 1.0.0
**Updated:** 2026-10-17
**Spec:** 277 (auto-discovered next free number; highest existing is 276)
**Parent program:** muse-memory-io (side chat "Muse memory IO")
**Owner:** Worker 01 (spec files, command implementation, docs, changelog)

---

## 1. Scope

IN SCOPE:

1. `gitmap llm train muse [--days N]` — positional sub-target of `llm train`
   (NOT a top-level command; the new-command checklist for top-level commands
   does not apply):
   - Pre-scan `llm train` args for the positional target `muse` and strip it
     BEFORE `parseTrainFlags()` runs (today every arg is parsed as a flag, so a
     bare `muse` would fail flag parsing).
   - Add `--days N` flag (default 7, the digest window in days).
   - Route `TrainTarget == "muse"` to the new runner; all existing `llm train`
     behavior is untouched.
2. Clone-or-pull of the repo-secrets vault:
   - Resolve the local checkout via the store special-repo entry:
     `cli/store` `OpenSpecialReposSplitDB().GetSpecialRepo("rs")`
     (short key `rs` verified in `cli/store/special_repos_split_db.go`).
   - If a configured `LocalPath` directory exists, use it.
   - If absent: clone `https://github.com/alimtvnetwork/repo-secrets.git` via
     `gh repo clone` into the work-folder root (current directory +
     configured name). Fail fast with clear guidance when `gh` is missing or
     `gh auth` is not authenticated — NEVER hang on a credential prompt
     (stdin non-interactive).
   - If present: `git pull --ff-only` (skipped with a note when no `origin`
     remote exists).
3. N-day diff: `git log --since="<N> days ago"` to enumerate changed files and
   `git rev-list --count --since=... HEAD` for the commit count.
4. Terminal summary + raw AI traversal checklist:
   - Summary: digest root, commits, files changed, sequenced digest count.
   - Checklist: file-by-file paths + suggested `gitmap cat` /
     `gitmap aum search` / `gitmap lf` commands so any LLM can walk the digest.
5. Digest contract consumed (authored by the memory-IO program, another
   worker): digest root at repo-secrets `Muse/README.md`, sequenced
   `Muse/01-*.md` … files, `Muse/chat-seeds/`.
   The command MUST tolerate the digest's absence: clear guidance message,
   exit 0.
6. Error discipline: ALL errors returned as real `*apperror.AppError` so the
   exit code is trustworthy (never swallow errors).
7. Go style (coding-guidelines): positive booleans only (`is`/`has` prefixes),
   functions <= 15 lines (prefer <= 8), wrapped errors, relative paths only,
   lowercase file names, no absolute paths in file content.
8. Package `llm` CANNOT import `cli/cmd` (circular). Clone-or-pull is
   reimplemented inside package `llm` using `cli/store` + `os/exec` git,
   following the `runGitInSpecialRepo` os/exec pattern
   (`cli/cmd/special_repos_ops.go:335`).
9. Colocated tests (source only, not run): positional `muse` stripping +
   `--days` parsing (default 7), missing-digest graceful message.
10. Docs + changelog: `cli/helpdoc/llm.md` `train muse` section with concrete
    invocations + sample output; `cli/helpdoc/examples_dev.go`
    `reg["llm-train"]` 3-part example (`gitmap llm train muse --days 7`);
    `changelog.md` entries under `## [Unreleased]` only; optional
    `cli/cmdllm/llmdocsgroups.go` matrix entry.

OUT OF SCOPE:

1. Digest authoring itself (`Muse/README.md`, sequenced `Muse/01-*.md`,
   `Muse/chat-seeds/` content) — authored by the memory-IO program on another
   box/worker. This spec only consumes the digest contract.
2. Version bumps, `version.json` edits, releases of any kind.
3. New top-level command plumbing — `train muse` is a positional sub-target of
   the existing `llm train` command only.
4. Touching any other command's behavior, flags, or help text.
5. Workflow/CI changes (`go build`, `go test` runs are CI's job — R1).

---

## 2. Command design

`gitmap llm train [muse] [--days N] [existing flags…]`

| Invocation | Behavior |
|---|---|
| `gitmap llm train muse` | Sync digest, 7-day window |
| `gitmap llm train muse --days 14` | Sync digest, 14-day window |
| `gitmap llm train muse --help` | Muse-target usage |
| `gitmap llm train` (no target) | Existing curriculum, unchanged |

Arg pre-scan strips the FIRST bare `muse` token that is not a flag and not a
flag value (`--skill-path muse`, `--heal-apply=muse` are preserved).

---

## 3. File plan

| File | Change |
|---|---|
| `cli/cmd/llm/llm_types.go` | Add `TrainTarget string`, `DaysN int` to `TrainOptions`; add `DefaultTrainMuseDays = 7` const |
| `cli/cmd/llm/llm_train.go` | Pre-scan/strip in `RunTrain`; add `--days` flag + clamp in `parseTrainFlags`; route `TrainTarget == "muse"`; update `printTrainUsage`; add `printTrainMuseUsage` |
| `cli/cmd/llm/llm_train_muse.go` | NEW — `runTrainMuse(opts)` + clone/pull/diff/summary/checklist helpers |
| `cli/cmd/llm/llm_train_test.go` | New tests (stripping, `--days`, missing-digest) |
| `cli/cmdllm/llm_help_menu.go` | `train muse [--days N]` usage lines + `--days` flag entry |
| `cli/cmdllm/llmdocsgroups.go` | Optional Utilities-group matrix entry |
| `cli/helpdoc/llm.md` | `train muse` section + cheat-sheet row |
| `cli/helpdoc/examples_dev.go` | `reg["llm-train"]` second example |
| `changelog.md` | `## [Unreleased]` entries |
