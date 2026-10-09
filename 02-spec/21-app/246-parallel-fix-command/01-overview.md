# Spec 246 — Parallel Fix Command (native port of the AI autofix scripts)

## Mindset
Same staging-repo mindset as spec 243 (`.ai-memory/what-to-read.md` → "The Mindset"):
velocity over ceremony, backup-first, small files, aliases are a kept feature,
DRY, enforcement beats documentation.

## User Request (verbatim)
```text
there was a AI script for fixing the file and coding right you can check want you to have that feature inside GitMap and also add this feature in the LLM train and also skills the reason I'm saying this because GitMap should be able to fix this Unicode encoding or any other issues parallelly yeah applying the grouping something like this so it it should have this process dealing with and fixing this parallelly and we can change the worker threads so this is how it should be changeable from the CLI
```

## Program goal
Port the 8 AI autofix scripts in `03-ai-scripts/` (encoding normalizer, newline
fixer, guideline autofixer, relative-path fixer, naming autofixer, Go code
formatter, misspell auditor, Markdown gap fixer — ~1096 lines total, stdlib-only,
sharing `02-shared-engine.py`) into a native GitMap command that:

1. fixes Unicode encoding and other file-hygiene issues **in parallel** with a
   CLI-controllable worker thread pool (`--workers N` / `-w N`);
2. applies the owner's requested **grouping**: one `--category` per concern
   (`encoding`, `newlines`, `naming`, `paths`, `gofmt`, `misspell`, `markdown`,
   `guidelines`), selectable as a subset, defaulting to all;
3. is **teachable**: the feature is documented inside `gitmap llm train` output
   and the Antigravity skill (`.agents/skills/gitmap/SKILL.md`), so AI agents
   discover and use it the skill way.

## Non-goals
- Do NOT delete, move, or modify anything under `03-ai-scripts/` — the Python
  scripts remain the reference implementation and the CI fallback.
- Do NOT touch the existing `gitmap fix` (git-state remediation) or the
  `cli/cmdfix/` package — see D2.
- Do NOT touch versioning, release flow, or the version cadence.
- Do NOT prune or rename existing aliases — aliases are a kept feature.
- `--staged` mode (script 26 has it), editor/IDE integration, network calls:
  none of these in v1.
- No new dependencies outside the Go stdlib plus existing `cli/*` packages.
- No big-bang help migration (program 243 rule) — only the new command's topic.

## Decisions (2026-10-09, recorded)
- D1: Pure-Go reimplementation, not subprocess calls into the Python scripts.
  The scripts stay untouched as reference; Go becomes the fast parallel path.
- D2: NAME COLLISION — `fix` is TAKEN. `cli/cmd/rootcore.go:79` already
  dispatches `{[]string{"fix"}, ...}` → `cmdfix.RunFix` (git-state remediation:
  stash/wip/discard flows), and the `cli/cmdfix/` package exists with that
  command's implementation. The new command is `gitmap autofix`
  (short alias `afx`), implemented in a NEW package `cli/cmdautofix/`.
  Both names were verified free of dispatch entries and package collisions.
- D3: Dry-run by default (audit only, exit 1 on violations). `--apply` writes.
  This unifies the scripts' inconsistent CLI shapes (`--fix` vs `--check-only`,
  `--path/-p` vs positional) and drops script 05's inverted default-fix.
- D4: Categories are the 8 groups above; `--category a,b` selects a subset,
  default = all. Unknown category name = usage error (exit 2).
- D5: The worker pool runs over FILES (not categories): each worker applies all
  selected categories to one file, then takes the next. `--workers N` / `-w N`,
  default `runtime.NumCPU()`. This mirrors script 26's thread-pool-over-files.
- D6: Exit codes: 0 clean · 1 violations found (dry-run) or unfixable findings
  remain (`--apply`) · 2 tool error (bad flag, unreadable path, `gofmt` missing).
  CI-compatible, same contract as the scripts.
- D7: Byte-safety: null-byte probe (first 8 KiB) skips binaries; files that are
  not valid UTF-8 are violations in dry-run and SKIPPED under `--apply` (never
  write lossy bytes); rewrite a file only when the transformed bytes differ
  from the original (no-op writes avoided, mtime preserved).
- D8: `gofmt` category execs the `gofmt` binary (`gofmt -w`) exactly like script
  26 does — NOT `go/format` stdlib — so output matches the repo's toolchain.
  If `gofmt` is absent from PATH, the category reports a tool error → exit 2.
- D9: `misspell` does CASE-PRESERVING replacement, fixing script 27's known bug
  (its case-insensitive replace lowercased capitalized words: `Behaviour` →
  `behavior`). Patterns: ALL-UPPER, Title-Case, lowercase; anything else falls
  back to the dictionary's lowercase form. Script 27's "delegates to misspell
  binary" docstring was never implemented — do NOT port it.
- D10: `paths` category: the hard-coded `coding-guidelines` URI pattern from
  script 07 becomes a package constant (`defaultRepoURI`) overridable via
  `--uri-pattern`. Default behavior is byte-identical to the script; it is now
  parameterizable, not removed.
- D11: `naming` (script 08) stays REPORT-ONLY even under `--apply`:
  `== true`/`=== true`/`== True` style needs human judgment; the script has no
  fix mode and we do not invent one.
- D12: `guidelines` composite = `newlines` fix + `naming` check (mirrors script
  05), always obeying the command's unified dry-run/`--apply` rule.
- D13: LLM-train + skill integration: `cli/cmdautofix` exposes
  `SkillSnippetMD() string` (static markdown block, spec 02 §7). The lead wires
  it into `gitmap llm train` output (`cli/cmd/llm/llm_train.go`) and
  `.agents/skills/gitmap/SKILL.md`. The 246 subtasks do NOT touch those
  lead-owned shared surfaces.
- D14: Help is DRY per program 243: the command's help content lives in Go
  structs (source of truth); `helptext/*.md` topics are GENERATED, never
  hand-edited.

## Scope boundaries
- New code lives ONLY in `cli/cmdautofix/` (+ one dispatch entry and one import
  in `cli/cmd/rootutility.go`, owned by subtask 02).
- Reads: anything, via `gitmap aum search` / `gitmap find` / `gitmap cat` only
  (TOTAL BAN on grep/rg/git grep/Select-String). Relative paths everywhere.
- Every Go file ≤ ~300 lines, grouped by concern (standing rule).
- Implementers verify with `go build ./...` from `cli/`. NEVER run `go test`
  without the owner's explicit command (standing rule).
