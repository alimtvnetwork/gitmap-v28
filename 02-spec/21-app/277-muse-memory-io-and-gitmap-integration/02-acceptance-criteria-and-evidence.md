# Spec 277 — Muse Memory IO and GitMap Integration: acceptance criteria and evidence

Each criterion names the exact evidence the implementer must produce.
Evidence = file path + function name (+ sample output where marked).

- **AC-01 — Positional `muse` stripping:** `gitmap llm train muse --days 3`
  reaches the muse runner; `gitmap llm train` (no target) runs the existing
  curriculum unchanged. Evidence: `cli/cmd/llm/llm_train.go` `stripTrainTarget`
  + `RunTrain` routing on `opts.TrainTarget == "muse"`; test
  `TestStripTrainTargetMuse` in `cli/cmd/llm/llm_train_test.go`.
- **AC-02 — Flag-value safety:** `--skill-path muse` and
  `--heal-apply=muse` are NOT stripped as the target. Evidence:
  `cli/cmd/llm/llm_train.go` `stripTrainTarget`; test
  `TestStripTrainTargetPreservesFlagValue`.
- **AC-03 — `--days` default:** no `--days` → `opts.DaysN == 7`
  (`DefaultTrainMuseDays` in `cli/cmd/llm/llm_types.go`); `--days 0` or
  negative clamps to 7. Evidence: `cli/cmd/llm/llm_train.go`
  `parseTrainFlags`; tests `TestParseTrainFlagsDaysDefault`,
  `TestParseTrainFlagsDaysValue`.
- **AC-04 — repo-secrets resolution via store:** local path resolved through
  `cli/store` `OpenSpecialReposSplitDB().GetSpecialRepo("rs")` when the record
  has a `LocalPath` directory; package `llm` does NOT import `cli/cmd`.
  Evidence: `cli/cmd/llm/llm_train_muse.go` `defaultResolveRepoSecretsPath`;
  no `cli/cmd` import in `cli/cmd/llm`.
- **AC-05 — Clone-or-pull:** absent checkout → `gh repo clone
  https://github.com/alimtvnetwork/repo-secrets.git` with non-interactive
  stdin; present checkout → `git pull --ff-only` (skipped with a note when no
  `origin` remote). Missing `gh` or missing `gh auth` → immediate failure
  with actionable guidance, no prompt hang. Evidence:
  `cli/cmd/llm/llm_train_muse.go` `cloneRepoSecrets`, `checkGhReady`,
  `pullRepoSecrets`, `runGitCmd` (os/exec pattern per
  `cli/cmd/special_repos_ops.go:335`).
- **AC-06 — N-day diff:** changed files from
  `git log --since="<N> days ago" --name-only --pretty=format:` (deduplicated,
  sorted); commit count from `git rev-list --count --since=... HEAD`.
  Evidence: `cli/cmd/llm/llm_train_muse.go` `gitLogChangedFiles`,
  `uniqueNonEmptyLines`, `parseCount`.
- **AC-07 — Terminal summary:** prints digest root, commits, files changed,
  sequenced-digest count, sequenced `Muse/*.md` list, changed-file list.
  Evidence: `cli/cmd/llm/llm_train_muse.go` `printDigestSummary`,
  `listSequencedDigestFiles`; sample output recorded in
  `cli/helpdoc/llm.md`.
- **AC-08 — AI traversal checklist:** file-by-file paths + suggested
  `gitmap cat` / `gitmap aum search` / `gitmap lf` commands; paths rendered
  relative to CWD. Evidence: `cli/cmd/llm/llm_train_muse.go`
  `printTraversalChecklist`, `relForDisplay`.
- **AC-09 — Missing digest tolerated:** `repo-secrets` present but no `Muse/`
  dir → clear guidance message, exit 0 (no error). Evidence:
  `cli/cmd/llm/llm_train_muse.go` `printDigestMissingGuidance`; test
  `TestRunTrainMuseMissingDigest` (temp git repo fixture, resolver override).
- **AC-10 — Error discipline:** every failure path returns
  `*apperror.AppError` (wrapped, never swallowed). Evidence: all functions in
  `cli/cmd/llm/llm_train_muse.go` return `*apperror.AppError`; no bare
  `fmt.Errorf` escapes.
- **AC-11 — Style gates:** positive booleans only (`is`/`has` prefixes),
  functions ≤ 15 lines (prefer ≤ 8), wrapped errors, lowercase file names,
  relative paths only, no absolute paths in file content. Evidence:
  `python 03-ai-scripts/05-guideline-autofixer.py cli/cmd/llm --check-only
  --ext .go` exit 0 with files scanned > 0; `gofmt -l` clean on touched files.
- **AC-12 — Help surfaces:** `gitmap llm train muse --help` prints
  muse-specific usage (`printTrainMuseUsage` in `cli/cmd/llm/llm_train.go`);
  `buildLlmHelpMenu()` in `cli/cmdllm/llm_help_menu.go` carries
  `train muse [--days N]` usage lines + a `--days <N>` footer flag entry.
- **AC-13 — Docs + changelog:** `cli/helpdoc/llm.md` has a `train muse`
  section with concrete invocations + sample output; `reg["llm-train"]` in
  `cli/helpdoc/examples_dev.go` has a `gitmap llm train muse --days 7`
  Invocation + ExpectedOutput (3-part format); `changelog.md` entries live
  under `## [Unreleased]` only — no new version section, no version bump.
- **AC-14 — Boundary hygiene:** none of the dirty-file list from another run
  is touched; no git commands executed; no builds/tests run (R1); shared
  indexes (`.ai-memory/plans/readme.md`, `02-spec/21-app/readme.md`) untouched.
