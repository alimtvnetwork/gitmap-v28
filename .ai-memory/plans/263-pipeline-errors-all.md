# Plan 263 — `pipeline errors all` + `te` shortcut + minor release

## User request (verbatim)
Owner (2026-10-09): "Forget map, can you please include or add a new command that's called git map pipeline space errors space all. That is going to scan all the repositories that we have and all the repositories errors will be combined and show it to the terminal um and also if we do hyphen hyphen file then that would be coming from a single file which AI models can read and fix all these issues. Okay? So make sure you implement that command and make a minor bump and make a release. Is it clear?" Follow-up: "Git map uh pipeline error space all it can also have a shortcut like te space all okay write that and go ahead with this do not wait for my confirmation so start writing throughout parallelly."

## Scope
- Repo: `~/workspace/repos/gitmap-v28`, main at 8c9259e (pulled 2026-10-09; tree has 23 pre-existing modified files from parallel streams — never touched).
- Spec number 263 issued via `gitmap spec next` (dogfooded). Backup branch `backup/263-pipeline-errors-all` pushed to origin (git-direct; `space backup-branch` refused on dirty tree, per its D2 guard).
- Existing per-repo `pipeline errors`/`pe` behavior stays untouched.

## Research synthesis (R01 LEAD_FALLBACK + R02)
- Catalog: `store.OpenDefault()` + `db.ListRepos()` → `[]model.ScanRecord` (Slug + AbsolutePath). `cli/store/repo.go:90`.
- Per-repo pipeline DB: `pipelinedb.ResolvePipelineDbPath(slug)` → `<bindir>/data/pipeline/<slug>/sql.db`; opener `pipelinedb.OpenPipelineSplitDb(slug)`. Slug via `pipelinedb.SanitizeRepoSlug(r.Slug)`.
- Existing all-repos path: `executeAllPipelineErrorLogs` at `cli/cmdpipeline/pipeline_all_errors.go:36`, types `AllPipelineSummary`/`AllPipelineFailedItem`. **Gap 1:** discovers repos by filesystem (`discoverPipelineRepoSlugs`, pipeline_all_errors.go:77), not catalog. **Gap 2:** `--file` parsed (`PipelineErrorFlags.FilePath`, pipeline_flags.go:62) but NOT wired into all-repos path.
- `te` is FREE (verified: "Command 'te' is not recognized"). Registration pattern: `cli/cmd/rootutility.go:453` — `{[]string{"pe", "pipeline-errors", "pipeline_errors", "ee"}, func() error { return cmdpipeline.RunPipelineErrors(argsTail()) }}`. New entry mirrors it.
- `--file` convention: `os.WriteFile(path, content, 0644)`; `pe history-ai` does dual md+json. For 263: single AI-readable markdown file (+ terminal summary still shown).
- Byte-faithfulness: `--file` content must NOT pass through glyph filtering (the `cat` lesson). Write raw bytes.

## Decisions
- D1: "ALL" = all catalog repos via `ListRepos()`; repos without pipeline data appear as "no pipeline data" rows (not silently skipped).
- D2: `--file <path>` writes ONE markdown file (AI-readable: per-error blocks with repo path, workflow, step, timestamp, full error text, log excerpt). Terminal still prints the summary.
- D3: `te` routes to `cmdpipeline.RunPipelineErrors(argsTail())` — so `te all` ≡ `pe all` ≡ `pipeline errors all`, flags included.
- D4: Keep changes inside `cli/cmdpipeline/` (+ 1 dispatch line); new formatter in its own file ≤300 lines.

## Wave plan
- Wave 1 (done): research + spec 263 + subtasks.
- Wave 2: implement — (a) catalog-based discovery in all-repos path, (b) `--file` wiring with AI-readable markdown formatter (new file), (c) `te` dispatch entry + help.
- Wave 3: lead verification — `go build ./...` (lead-run, owner-authorized), live `pipeline errors all` + `te all --file /tmp/errors.md` against local catalog, file well-formedness check.
- Wave 4: atomic `gitmap cpf` + push (targeted `git add` of 263 files ONLY — other streams' dirty files must not be swept in), minor bump via `03-ai-scripts/37-bump-version.py`, tag + push, `gh` release, rebuild `~/.local/bin/gitmap`, verify version.

## Checkboxes
- [x] Step 0: git pull (8c9259e), spec 263 via `spec next`, backup branch pushed.
- [x] Research 01 (LEAD_FALLBACK: `te` free, dispatch pattern at rootutility.go:453) + Research 02 (catalog/DB/specs) DONE.
- [x] Spec 263: 01-overview.md + subtask 01-aggregator-implementation.md DONE (Spec Agent 1). Key decisions: D4 (--file+--json → JSON to file), D6 (os.Stat before OpenPipelineSplitDb — opener creates empty DBs as side effect), D7 (SanitizeRepoSlug on catalog slugs). Awaiting Spec Agent 2 (02-ai-report-format.md + subtask 02).
- [x] Wave 2: Worker 02 DONE — `te` dispatch entry (rootutility.go:454) + help usage line/tip (pipeline_help_menu.go); gofmt clean. DEFERRED to lead: one-word cobra alias addition in `cli/cmd/root_cobra_completion.go` (`makeTopLevelPECmd()` Aliases += `"te"`). Worker 01 (aggregator) still running.
- [x] Wave 2: Worker 01 DONE (aggregator: catalog discovery, --file wiring, 292-line formatter) + Worker 02 DONE (`te` dispatch + help). Lead fix: duplicate-slug repos dropped from markdown table (repoDataStatus map keyed by slug) — fixed via ordered repoList, verified 10/10 rows.
- [x] Wave 3: lead verification DONE (`go build ./...` exit 0 personally; live `pipeline errors all` + `te all --file` against real catalog; `te` ≡ `pe` byte-identical modulo timestamp; cobra `te` alias added by lead).
- [x] Wave 4: committed 19d1bf2 (pushed; NOTE: `gitmap cpf` swept 23 pre-existing other-stream files despite targeted add — content safe, lesson re-logged), minor bump to 6.522.0 via 37-bump-version.py, `go generate` re-run, tag v6.522.0 pushed, GitHub release published, ~/.local/bin/gitmap rebuilt → v6.522.0 verified.

## Assumptions
- Pipeline DBs exist for repos with pipeline runs; repos without get "no pipeline data" rows.
- `AllPipelineFailedItem` fields (repoSlug, workflowName, runId, branch, sha, errorSummary, updatedAt) are sufficient for the AI-readable blocks; workers extend if the record has more.

## Conflicts
- Parent brief (owner: builds authorized, per-wave commits+push, release requested) supersedes skill R1/R8/R10 per Top-Instruction Priority Mandate.

## Follow-ups (out of scope)
- Program 250 Wave 4 (stdout writer refactor) — separate program.
- The 23 pre-existing dirty files (other streams' work).
