# Acceptance Criteria and Evidence — task 244-e2e-verify-program-243-changes

Spec owner: Spec Agent 2. Sibling spec: `01-test-scope-and-plan.md` (Spec Agent 1) holds the user request verbatim and the test plan; this file holds the acceptance IDs, expected outcomes, evidence format, and verdict rules for program 243 verification.

- Program under test: 243 (waves A–D, consolidated commit f127f7a on main)
- E2E binary: built from pristine HEAD f127f7a, lives OUTSIDE the repo at `~/workspace/gitmap-e2e-244/gitmap` (scratch, never committed)
- Out of scope: the in-flight 243b uncommitted consolidation (another coordinator is actively working in the tree)
- All paths below are relative to repo root (`~/workspace/repos/gitmap-v28`)

## User Request (Verbatim)

See `02-spec/21-app/244-e2e-verify-program-243-changes/01-test-scope-and-plan.md` (Spec Agent 1). Not duplicated here by design.

## Traceability

Each acceptance ID maps to a test-group section in `01-test-scope-and-plan.md`:

| Acceptance IDs | Test group in 01-test-scope-and-plan.md |
|---|---|
| A01–A10 | backup-branch contract |
| B01–B07 | help displayer invocations |
| C01–C04 | unknown-command handling |
| D01–D05 | regression and build |
| E01 | environment and binary identity |

## Acceptance Table

Conventions: matcher column is either `exact:` (whole-output or whole-line match) or `re:` (regex/substring must appear). "PASS if" is the exact condition the runner evaluates per ID.

### Group A — backup-branch contract

Form: `gitmap space backup-branch "<task string>" [--no-push] [--force|-f] [--help|-h]`. No `--dry-run` flag exists. Slug rules: lowercase, `[a-z0-9-]`, max 60 chars.

| ID | Command | Exit | Stdout/stderr matcher | PASS if |
|---|---|---|---|---|
| A01 | `gitmap space backup-branch --help` | 0 | re: stdout contains `Usage` and `--no-push` | exit 0, plain-text usage documents the flags |
| A02 | `gitmap space --help backup-branch` | 0 | re: contains `About`, `Behavior`, `Flags`, `Examples` group headers | exit 0, Displayer groups rendered |
| A03 | `gitmap space backup-branch` (no task string) | 1 | re: `missing task string` | exit 1, no branch created |
| A04 | `gitmap space backup-branch "!!!" --no-push` (empty slug) | 1 | re: `empty backup slug` | exit 1, no branch created |
| A05 | `gitmap space backup-branch "dirty tree test" --no-push` (scratch repo with uncommitted changes) | 1 | re: `uncommitted changes` | exit 1, no branch created |
| A06 | clean scratch repo (no origin), `gitmap space backup-branch "E2E Smoke Test_42" --no-push` | 0 | re: `^backup/e2e-smoke-test-42 @ [0-9a-f]+$` | exit 0, stdout matches, `git rev-parse --verify refs/heads/backup/e2e-smoke-test-42` succeeds |
| A07 | same as A06 without `--no-push` | 1 | re (case-insensitive): `push`, `origin`, or `fatal` | exit 1, clean push error, no hang |
| A08 | pre-create `backup/dup-case`, rerun task `dup case` without `--force` → exit 1; then with `--force` → exit 0 | 1 then 0 | re: `already exists`; after `--force` branch points at current HEAD | refusal then force-recreate both behave |
| A09 | 80-char task string with `--no-push` | 0 | slug part after `backup/` is ≤ 60 chars | exit 0, 60-char cap honored |
| A10 | informational: `backup/2026-10-08-pre-help-displayer` visible via `git branch -r` in the scratch clone | n/a | branch present | WARN-only, never FAIL |

### Group B — help displayer (7 invocations, all exit 0)

| ID | Command | Exit | Stdout matcher | PASS if |
|---|---|---|---|---|
| B01 | `gitmap scan --help` | 0 | re: header line and group headers | exit 0, Displayer rendered |
| B02 | `gitmap s --help` | 0 | same shape as B01 (alias) | exit 0, alias resolves |
| B03 | `gitmap space --help` | 0 | re: header, subcommand list incl. `backup-branch`, trailing `* ` suggestion line | exit 0, all present |
| B04 | `gitmap space --help common` | 0 | re: header and group headers | exit 0; suggestion-line check is WARN-only (subcommand help carries no `* ` line) |
| B05 | `gitmap space --help backup-branch` | 0 | re: `About`, `Behavior`, `Flags`, `Examples` | exit 0, groups rendered |
| B06 | `gitmap pull-all --help` | 0 | re: header, groups, trailing `* ` suggestion line | exit 0, all present |
| B07 | `gitmap pa --help` | 0 | same shape as B06 (alias) | exit 0, alias resolves |

### Group C — unknown-command handling

| ID | Command | Exit | Stderr matcher | PASS if |
|---|---|---|---|---|
| C01 | `gitmap scna` | 1 | re: stderr contains `Did you mean` and `scan` | exit 1, suggestion on stderr |
| C02 | `gitmap docker` | 1 | re: stderr contains `Did you mean` | exit 1, suggestion on stderr |
| C03 | `gitmap releas` | 1 | re: stderr contains `Did you mean` and `release` | exit 1, suggestion on stderr |
| C04 | `gitmap somethingrandom` | 1 | re (case-insensitive): stderr contains `unrecognized command` and NOT `Did you mean` | exit 1, plain unknown-command error (E1001) on stderr |

### Group D — regression and build

| ID | Command / check | Exit | Matcher | PASS if |
|---|---|---|---|---|
| D01 | `gitmap --version` | 0 | exact: stdout line is `gitmap v6.515.0` | exit 0 and exact string match |
| D02 | file-size gate default: inspect `03-ai-scripts/52-file-size-check.py` (pristine f127f7a) | n/a | re: script defines `DEFAULT_LINE_LIMIT = 500` | default limit is 500 — the committed reality at f127f7a (see Findings for the 300-doc gap) |
| D03 | `cli/enums/` exists in the checked-out repo | n/a | dir listing contains `cli/enums/` with Go files | directory exists with enum sources |
| D04 | `go build ./...` run from `cli/` | 0 | n/a (exit code only) | exit 0, no build errors |
| D05 | pre-commit hook wires the file-size check | n/a | re: `50-fastgate.py` references the gate with limit 500 (docstring "500-line", message "(limit 500)") | hook docstring/message consistent with the 500 default |

### Group E — environment and binary identity

| ID | Check | PASS if |
|---|---|---|
| E01 | scratch binary at `~/workspace/gitmap-e2e-244/gitmap` exists, is executable, and its sha256 is recorded at run time | binary present and runnable; sha256 recorded in the report (path stays outside the repo, never committed) |

## Evidence Format for verification-report.md

The report is written to `~/workspace/gitmap-e2e-244/verification-report.md` (outside the repo, never committed). Required structure:

```markdown
# Verification Report — program 243 (task 244-e2e-verify-program-243-changes)

## Environment
- binary: ~/workspace/gitmap-e2e-244/gitmap
- binary sha256: <recorded at run time>
- HEAD: f127f7a
- go version: <go version output>
- run date: <YYYY-MM-DD, Asia/Kuala_Lumpur>

## Group A — backup-branch contract
| ID | command | exit | output excerpt (first 5 lines or matched line) | result | note |
... one row per ID (A01..A10, B01..B07, C01..C04, D01..D05, E01)

## Group B — help displayer
...

## Group C — unknown-command handling
...

## Group D — regression and build
...

## Group E — environment and binary identity
...

## Findings (informational, never verdict-failing)
- <finding text>

## Verdict
<VERIFIED | NOT VERIFIED — see RCA notes>
```

Per entry, the runner records: the exact command, exit code, output excerpt (first 5 lines of stdout/stderr or the single matched line for matcher IDs), PASS/FAIL (WARN allowed only for informational assertions, see verdict rule), and a one-line note.

## Verdict Rule

- ALL assertions (A01–A10, B01–B07, C01–C04, D01–D05, E01) PASS → verdict line: `program 243 committed changes VERIFIED`
- Any FAIL → verdict line: `NOT VERIFIED — see RCA notes`, followed by the list of failing IDs with observed-vs-expected (command, expected, observed) under the Verdict section
- Informational assertions (e.g. backup branch visible in remote listing after a push) may be recorded as WARN without failing the verdict; WARN rows must still carry a note explaining why

## Discrepancy Policy

Known discrepancy (NOT a 243 failure): at committed HEAD f127f7a the file-size gate enforces 500 (`52-file-size-check.py` default and `50-fastgate.py` docstring/message agree) while the 300-line cap exists only in docs (mindset file, spec 243, commit f127f7a message). The in-flight 243b working tree has already flipped the code default to 300 — that alignment is 243b's scope. Logged under Findings, NEVER counted as a failure against program 243.

## Notes for the Runner

- The e2e binary is built from pristine HEAD f127f7a; do not rebuild from the working tree (243b consolidation is in flight and out of scope)
- Scripts and report live outside the repo under `~/workspace/gitmap-e2e-244/`; nothing from the e2e run is committed
- Never run git commands inside `~/workspace/repos/gitmap-v28` (read-only git is allowed only inside the scratch dirs under `~/workspace/gitmap-e2e-244/`)
- All paths recorded in the report must be relative to repo root (never absolute, no file:/// URIs)
