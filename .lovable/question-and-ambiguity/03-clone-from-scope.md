# 03 — `gitmap clone` reads JSON/CSV: scope vs existing commands

## Original task

> Add a `gitmap clone` command that reads JSON/CSV from disk, shows a dry-run plan by default, and then executes the git clones with a clear summary.

## Ambiguity

The literal phrase "Add a `gitmap clone` command" conflicts with the project state — `gitmap clone` already exists (`gitmap/cmd/clone.go`, accepts direct HTTPS/SSH URLs, see `mem://features/clone-direct-url`). And `gitmap cn --csv <path>` (`clonenextbatch.go`) ALREADY reads CSV from disk and processes batches.

Three reasonable interpretations:

### A — Extend existing `gitmap clone` to ALSO accept `--from <file>` flag
- **Pros**: One command, one mental model. URLs and files coexist.
- **Cons**: Mixes two semantically different inputs in one entry point. Behaviour change to a published command. The "dry-run by default" rule would either break the existing URL form or be limited to file mode (confusing).

### B — Add a new sibling subcommand `gitmap clone-from <file>` (RECOMMENDED)
- **Pros**: Zero behaviour change to existing `gitmap clone`. Clear semantics: file in, plan out, execute on confirm. Can have its own help, flags, and tests. Matches the existing `gitmap clone` / `gitmap cn` / `gitmap clone-next-batch` family naming style.
- **Cons**: One more verb on the CLI surface (mitigated by adding to the help index and completion generator).

### C — Tell the user `gitmap cn --csv` already exists and don't build anything
- **Pros**: No new code; honest about the prior art.
- **Cons**: `cn --csv` does VERSION BUMPING (clones `vN+1` of existing repos), NOT plain git clone of new repos from a list. Different semantics. JSON input not supported. No "dry-run by default" — it executes immediately. So `cn --csv` does NOT cover this request.

## Decision taken

**Option B — new `gitmap clone-from <file>` subcommand.** Implementing:

1. **CLI verb**: `gitmap clone-from <file>` (alias `cf` to match the project's two-letter alias style: `cn`, `cf`).
2. **Input formats**: auto-detect from file extension. `.json` → array of `{url, dest?, branch?, depth?}` objects; `.csv` → header row `url,dest,branch,depth` (only `url` required).
3. **Dry-run by default**: prints the planned `git clone` invocations as a table; exits 0 without cloning.
4. **`--execute` flag**: actually runs the clones. Sequential by default; `--parallel <N>` for fan-out (mirrors `cn --csv` pattern).
5. **Summary**: per-repo `ok` / `skipped` / `failed` counts + a CSV report at `.gitmap/clone-from-report-<unixts>.csv` (mirrors `cn --csv` report convention).
6. **Skip rule**: if `dest` already exists and is a non-empty directory → mark `skipped` (idempotent re-runs).
7. **No DB writes in dry-run**; on `--execute`, upserts each cloned repo into the standard scan DB (matches `gitmap clone <url>` behavior per `mem://features/clone-direct-url`).
8. **No shell handoff** (batch operation; cwd-changing makes no sense for N repos).
9. **Help text** at `gitmap/helptext/clone-from.md` per the help-system convention.
10. **Tests**: parser tests for both formats + a dry-run snapshot test + an end-to-end test using a local bare repo (mirrors `clonenextbatchconcurrent_e2e_csv_test.go` pattern).
11. **Constants**: all messages / errors / flag names in a new `constants_clonefrom.go`.
12. **Version bump**: 3.159.0 → 3.160.0 (minor, per project rule).
13. **Completion generator**: add `// gitmap:cmd top-level` marker so tab-completion picks it up.

User can override by saying e.g. "do option A" and I'll fold the file-reading behavior into the existing `gitmap clone`.
