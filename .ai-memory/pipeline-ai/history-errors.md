# Pipeline Historical Errors & AI Anti-Mistake Training Dossier

- **Repository:** `alimtvnetwork/gitmap-v28`
- **Generated At:** 2026-10-02T03:23:21Z
- **Failing Commits Analyzed:** 5
- **Total Failure Diagnostics:** 24

> **Purpose:** This dossier exposes all historical pipeline failures across recent commits
> to train AI coding agents so they NEVER repeat these architectural, linting, or testing mistakes.

## Executive Failure Taxonomy

| Commit | Workflow | Failing Step(s) | Primary Issue Category |
|--------|----------|-----------------|------------------------|
| `db17b55` | CI | Failed Step | Pipeline Execution |
| `f05646f` | CI | Failed Step | Pipeline Execution |
| `a98597d` | CI | Failed Step | Pipeline Execution |
| `9da3535` | CI | Failed Step | Pipeline Execution |
| `641a754` | CI | Failed Step | Pipeline Execution |

### [1] Commit `db17b55` (Branch: `main`)

- **Workflow:** CI
- **Run ID:** [36956759998](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36956759998)
- **Recorded At:** 2026-10-02T02:41:05Z (Duration: 431s)

#### Job: `CI` | Step: `Failed Step`
- **Summary:** [staticcheck] SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (NEW in /tmp/lint-current/report.json)
    ❌ store/migrations.go:270: [staticcheck] SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value
[staticcheck] SA4023: this comparison is always true (NEW in /tmp/lint-current/report.json)
    ❌ store/migrations.go:302: [staticcheck] SA4023: this comparison is always true
Process completed with exit code 1.

│ File:     cmdos/os_dev_clean.go:122:2
Run bash .github/scripts/check-single-linter-diff.sh cli
bash .github/scripts/check-single-linter-diff.sh cli
shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
env:
FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
GOTOOLCHAIN: local
LINTER: gocritic
BASELINE: /tmp/lint-gocritic-baseline/report.json
CURRENT_OUT: /tmp/lint-gocritic-current/report.json
========================================================================
FAIL: 5 new gocritic finding(s) detected!
Script:   .github/scripts/check-single-linter-diff.py
Target:   cli
To reproduce locally:
golangci-lint run --no-config --disable-all --enable=gocritic ./cli/...
Inspect the exact files and lines listed above to resolve.
GOCRITIC DIFF (baseline-diff, full-path only)
script   : .github/scripts/check-single-linter-diff.py (check-single-linter-diff.py)
analyzer : gocritic
target   : cli
command  : python .github/scripts/check-single-linter-diff.py cli --linter=gocritic
current  : /tmp/lint-gocritic-current/report.json
baseline : /tmp/lint-gocritic-baseline/report.json
+ NEW    : 5
[CI/CD ERROR REPORT] New gocritic Violations Detected (5 finding(s)):
[gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)
┌─ Error [1/5]: [gocritic]
│ Location: cmdos/os_dev_clean.go at line 122, col 2
│ Linter:   gocritic
│ Message:  ifElseChain: rewrite if-else to switch statement
│ Context:  NEW finding vs baseline
│ Script:   .github/scripts/check-single-linter-diff.py
└──────────────────────────────────────────────────────────
❌ cmdos/os_dev_clean.go:122:2: [gocritic] ifElseChain: rewrite if-else to switch statement
[gocritic] singleCaseSwitch: should rewrite switch statement to if statement (NEW vs baseline)
┌─ Error [2/5]: [gocritic]
│ File:     cmdos/os_tui_model.go:33:2
│ Location: cmdos/os_tui_model.go at line 33, col 2
│ Message:  singleCaseSwitch: should rewrite switch statement to if statement
❌ cmdos/os_tui_model.go:33:2: [gocritic] singleCaseSwitch: should rewrite switch statement to if statement
┌─ Error [3/5]: [gocritic]
│ File:     cmdpull/pull_efficient_render.go:51:3
│ Location: cmdpull/pull_efficient_render.go at line 51, col 3
❌ cmdpull/pull_efficient_render.go:51:3: [gocritic] ifElseChain: rewrite if-else to switch statement
[gocritic] appendAssign: append result not assigned to the same slice (NEW vs baseline)
┌─ Error [4/5]: [gocritic]
│ File:     gitignoreagm/sanitizer.go:78:12
│ Location: gitignoreagm/sanitizer.go at line 78, col 12
│ Message:  appendAssign: append result not assigned to the same slice
❌ gitignoreagm/sanitizer.go:78:12: [gocritic] appendAssign: append result not assigned to the same slice
┌─ Error [5/5]: [gocritic]
│ File:     gitutil/git_env.go:28:9
│ Location: gitutil/git_env.go at line 28, col 9
❌ gitutil/git_env.go:28:9: [gocritic] appendAssign: append result not assigned to the same slice
Process completed with exit code 1.

store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
    if err := insertRemappedGroupRepos(dbConn); err != nil {
    ^
store/migrations.go:270:6: SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (staticcheck)
    func insertRemappedGroupRepos(dbConn *sql.DB) error {
Process completed with exit code 1.

❌ FAIL: Found 24 absolute path / URI violation(s):
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:49: Hardcoded absolute repo path (D:\work\gitmap): SQLite text columns default to `BINARY` collation. Under `BINARY` collation, string comparisons eval
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:418: Hardcoded absolute repo path (D:\work\gitmap): - Given a database containing duplicate rows (`RepoId=1, AbsolutePath="D:\work\gitmap"` and `RepoId=
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Absolute file:/// URI with drive letter: - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Hardcoded absolute repo path (D:\work\gitmap): - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Absolute file:/// URI with drive letter: - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Hardcoded absolute repo path (D:\work\gitmap): - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:81: Hardcoded absolute repo path (D:\work\gitmap): | `D:\work\gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:82: Hardcoded absolute repo path (D:\work\gitmap): | `d:\work\gitmap\` | `d:\work\gitmap` | `d:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:83: Hardcoded absolute repo path (D:\work\gitmap): | `D:/work/./gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:84: Hardcoded absolute repo path (D:\work\gitmap): | `d:/WORK/gitmap` | `d:\WORK\gitmap` | `d:/WORK/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Absolute file:/// URI with drive letter: - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Hardcoded absolute repo path (D:\work\gitmap): - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Absolute file:/// URI with drive letter: - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Hardcoded absolute repo path (D:\work\gitmap): - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Absolute file:/// URI with drive letter: - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Hardcoded absolute repo path (D:\work\gitmap): - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Absolute file:/// URI with drive letter: - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Hardcoded absolute repo path (D:\work\gitmap): - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Absolute file:/// URI with drive letter: - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Hardcoded absolute repo path (D:\work\gitmap): - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Absolute file:/// URI with drive letter: - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Hardcoded absolute repo path (D:\work\gitmap): - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:443: Hardcoded absolute repo path (D:\work\gitmap): - Verify `D:\work\gitmap`, `d:\work\gitmap`, `D:/work/gitmap`, and `d:\work\gitmap\` all resolve to
    02-spec/22-app-issues/62-pull-all-duplicate-repositories-and-case-sensitivity-rca.md:51: Hardcoded absolute repo path (D:\work\gitmap): When parallel workers consumed the job queue, Worker A pulled `D:\work\gitmap` and Worker B pulled `
Process completed with exit code 1.

❌ FAIL: Found 3 nested-if / anti-compression violation(s) across 3 file(s):
    cli/cmd/nodes_cmd_test.go:93: Nested if statement found (depth 2 inside conditional block): if strings.Contains(between, ) {
    cli/cmdpull/pull_efficient_render.go:438: Nested if statement found (depth 2 inside conditional block): if relErr == nil && !strings.HasPrefix(rel, ) {
    cli/cmdpull/pull_remediation_hint.go:186: Nested if statement found (depth 2 inside conditional block): if hasOnlyUntracked {
Process completed with exit code 1.

# Tally PASS/FAIL package lines for the PR summary comment.
    passed=$(grep -cE '^ok[[:space:]]' /tmp/full-suite/test-output.txt || true)
    failed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)
    echo "packages_passed=$passed" >> "$GITHUB_OUTPUT"
    echo "packages_failed=$failed" >> "$GITHUB_OUTPUT"
    exit "$rc"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
--- FAIL: TestEveryHelpFileHasExamples (0.00s)
    examples_golden_test.go:42: help files missing `## Examples` section (2):
    - devtool.md
    - pulle.md
    FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/helptext	0.007s
Process completed with exit code 1.

store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
    if err := insertRemappedGroupRepos(dbConn); err != nil {
    ^
store/migrations.go:270:6: SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (staticcheck)
    func insertRemappedGroupRepos(dbConn *sql.DB) error {
Process completed with exit code 1.

❌ FAILED: Found 2 violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_efficient_render.go:438: Nested 'if' detected (depth 2): 'if relErr == nil && !strings.HasPrefix(rel, ) {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_remediation_hint.go:186: Nested 'if' detected (depth 2): 'if hasOnlyUntracked {'
Process completed with exit code 1.

❌ FAILED: Found 3 error management violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/store/migrations.go:242: Swallowed database call error detected: '_, _ = dbConn.Exec(sqlPrune)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:86: Swallowed database call error detected: '_, _ = conn.Exec(`DELETE FROM SearchQueryLog`)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:87: Swallowed database call error detected: '_, _ = conn.Exec(`VACUUM`)'
Process completed with exit code 1.
```text
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:44:09.9667870Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9668322Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9668738Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9669692Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9670345Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9671128Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9671857Z ^[[36;1m  BASELINE=""^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9672280Z ^[[36;1melif [ ! -f "$BASELINE" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9672932Z ^[[36;1m  echo "::notice::No baseline found — this run will seed the cache."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:44:09.9673576Z ^[[36;1mfi^[[0m
```

### [2] Commit `f05646f` (Branch: `main`)

- **Workflow:** CI
- **Run ID:** [36955632940](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36955632940)
- **Recorded At:** 2026-10-02T02:26:27Z (Duration: 339s)

#### Job: `CI` | Step: `Failed Step`
- **Summary:** store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
    if err := insertRemappedGroupRepos(dbConn); err != nil {
    ^
store/migrations.go:270:6: SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (staticcheck)
    func insertRemappedGroupRepos(dbConn *sql.DB) error {
Process completed with exit code 1.

❌ FAILED: Found 3 error management violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/store/migrations.go:242: Swallowed database call error detected: '_, _ = dbConn.Exec(sqlPrune)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:86: Swallowed database call error detected: '_, _ = conn.Exec(`DELETE FROM SearchQueryLog`)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:87: Swallowed database call error detected: '_, _ = conn.Exec(`VACUUM`)'
Process completed with exit code 1.

[staticcheck] SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (NEW in /tmp/lint-current/report.json)
    ❌ store/migrations.go:270: [staticcheck] SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value
[staticcheck] SA4023: this comparison is always true (NEW in /tmp/lint-current/report.json)
    ❌ store/migrations.go:302: [staticcheck] SA4023: this comparison is always true
Process completed with exit code 1.

│ File:     cmdos/os_dev_clean.go:122:2
Run bash .github/scripts/check-single-linter-diff.sh cli
bash .github/scripts/check-single-linter-diff.sh cli
shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
env:
FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
GOTOOLCHAIN: local
LINTER: gocritic
BASELINE: /tmp/lint-gocritic-baseline/report.json
CURRENT_OUT: /tmp/lint-gocritic-current/report.json
========================================================================
FAIL: 5 new gocritic finding(s) detected!
Script:   .github/scripts/check-single-linter-diff.py
Target:   cli
To reproduce locally:
golangci-lint run --no-config --disable-all --enable=gocritic ./cli/...
Inspect the exact files and lines listed above to resolve.
GOCRITIC DIFF (baseline-diff, full-path only)
script   : .github/scripts/check-single-linter-diff.py (check-single-linter-diff.py)
analyzer : gocritic
target   : cli
command  : python .github/scripts/check-single-linter-diff.py cli --linter=gocritic
current  : /tmp/lint-gocritic-current/report.json
baseline : /tmp/lint-gocritic-baseline/report.json
+ NEW    : 5
[CI/CD ERROR REPORT] New gocritic Violations Detected (5 finding(s)):
[gocritic] ifElseChain: rewrite if-else to switch statement (NEW vs baseline)
┌─ Error [1/5]: [gocritic]
│ Location: cmdos/os_dev_clean.go at line 122, col 2
│ Linter:   gocritic
│ Message:  ifElseChain: rewrite if-else to switch statement
│ Context:  NEW finding vs baseline
│ Script:   .github/scripts/check-single-linter-diff.py
└──────────────────────────────────────────────────────────
❌ cmdos/os_dev_clean.go:122:2: [gocritic] ifElseChain: rewrite if-else to switch statement
[gocritic] singleCaseSwitch: should rewrite switch statement to if statement (NEW vs baseline)
┌─ Error [2/5]: [gocritic]
│ File:     cmdos/os_tui_model.go:33:2
│ Location: cmdos/os_tui_model.go at line 33, col 2
│ Message:  singleCaseSwitch: should rewrite switch statement to if statement
❌ cmdos/os_tui_model.go:33:2: [gocritic] singleCaseSwitch: should rewrite switch statement to if statement
┌─ Error [3/5]: [gocritic]
│ File:     cmdpull/pull_efficient_render.go:51:3
│ Location: cmdpull/pull_efficient_render.go at line 51, col 3
❌ cmdpull/pull_efficient_render.go:51:3: [gocritic] ifElseChain: rewrite if-else to switch statement
[gocritic] appendAssign: append result not assigned to the same slice (NEW vs baseline)
┌─ Error [4/5]: [gocritic]
│ File:     gitignoreagm/sanitizer.go:78:12
│ Location: gitignoreagm/sanitizer.go at line 78, col 12
│ Message:  appendAssign: append result not assigned to the same slice
❌ gitignoreagm/sanitizer.go:78:12: [gocritic] appendAssign: append result not assigned to the same slice
┌─ Error [5/5]: [gocritic]
│ File:     gitutil/git_env.go:28:9
│ Location: gitutil/git_env.go at line 28, col 9
❌ gitutil/git_env.go:28:9: [gocritic] appendAssign: append result not assigned to the same slice
Process completed with exit code 1.

❌ FAILED: Found 2 violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_remediation_hint.go:186: Nested 'if' detected (depth 2): 'if hasOnlyUntracked {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_efficient_render.go:438: Nested 'if' detected (depth 2): 'if relErr == nil && !strings.HasPrefix(rel, ) {'
Process completed with exit code 1.

❌ FAIL: Found 3 nested-if / anti-compression violation(s) across 3 file(s):
    cli/cmd/nodes_cmd_test.go:93: Nested if statement found (depth 2 inside conditional block): if strings.Contains(between, ) {
    cli/cmdpull/pull_efficient_render.go:438: Nested if statement found (depth 2 inside conditional block): if relErr == nil && !strings.HasPrefix(rel, ) {
    cli/cmdpull/pull_remediation_hint.go:186: Nested if statement found (depth 2 inside conditional block): if hasOnlyUntracked {
Process completed with exit code 1.

❌ FAIL: Found 24 absolute path / URI violation(s):
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Absolute file:/// URI with drive letter: - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Hardcoded absolute repo path (D:\work\gitmap): - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Absolute file:/// URI with drive letter: - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Hardcoded absolute repo path (D:\work\gitmap): - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:81: Hardcoded absolute repo path (D:\work\gitmap): | `D:\work\gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:82: Hardcoded absolute repo path (D:\work\gitmap): | `d:\work\gitmap\` | `d:\work\gitmap` | `d:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:83: Hardcoded absolute repo path (D:\work\gitmap): | `D:/work/./gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:84: Hardcoded absolute repo path (D:\work\gitmap): | `d:/WORK/gitmap` | `d:\WORK\gitmap` | `d:/WORK/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Absolute file:/// URI with drive letter: - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Hardcoded absolute repo path (D:\work\gitmap): - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Absolute file:/// URI with drive letter: - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Hardcoded absolute repo path (D:\work\gitmap): - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Absolute file:/// URI with drive letter: - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Hardcoded absolute repo path (D:\work\gitmap): - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Absolute file:/// URI with drive letter: - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Hardcoded absolute repo path (D:\work\gitmap): - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Absolute file:/// URI with drive letter: - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Hardcoded absolute repo path (D:\work\gitmap): - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Absolute file:/// URI with drive letter: - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Hardcoded absolute repo path (D:\work\gitmap): - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:443: Hardcoded absolute repo path (D:\work\gitmap): - Verify `D:\work\gitmap`, `d:\work\gitmap`, `D:/work/gitmap`, and `d:\work\gitmap\` all resolve to
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:49: Hardcoded absolute repo path (D:\work\gitmap): SQLite text columns default to `BINARY` collation. Under `BINARY` collation, string comparisons eval
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:418: Hardcoded absolute repo path (D:\work\gitmap): - Given a database containing duplicate rows (`RepoId=1, AbsolutePath="D:\work\gitmap"` and `RepoId=
    02-spec/22-app-issues/62-pull-all-duplicate-repositories-and-case-sensitivity-rca.md:51: Hardcoded absolute repo path (D:\work\gitmap): When parallel workers consumed the job queue, Worker A pulled `D:\work\gitmap` and Worker B pulled `
Process completed with exit code 1.

# Tally PASS/FAIL package lines for the PR summary comment.
    passed=$(grep -cE '^ok[[:space:]]' /tmp/full-suite/test-output.txt || true)
    failed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)
    echo "packages_passed=$passed" >> "$GITHUB_OUTPUT"
    echo "packages_failed=$failed" >> "$GITHUB_OUTPUT"
    exit "$rc"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
--- FAIL: TestEveryHelpFileHasExamples (0.01s)
    examples_golden_test.go:42: help files missing `## Examples` section (2):
    - devtool.md
    - pulle.md
    FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/helptext	0.036s
Process completed with exit code 1.

store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
    if err := insertRemappedGroupRepos(dbConn); err != nil {
    ^
store/migrations.go:270:6: SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (staticcheck)
    func insertRemappedGroupRepos(dbConn *sql.DB) error {
Process completed with exit code 1.
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `golangci-lint (strict, fail on any error)`
- **Summary:** store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `Run ./.github/actions/policy-check`
- **Summary:** ❌ FAIL: Found 24 absolute path / URI violation(s):
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `Diff vs baseline (fail only on NEW findings)`
- **Summary:** [staticcheck] SA4023(related information): github.com/alimtvnetwork/gitmap-v28/cli/store.insertRemappedGroupRepos never returns a nil interface value (NEW in /tmp/lint-current/report.json)
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `Gocritic diff (baseline-diff, full-path only)`
- **Summary:** │ File:     cmdos/os_dev_clean.go:122:2
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `Run full test suite`
- **Summary:** # Tally PASS/FAIL package lines for the PR summary comment.
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

#### Job: `CI` | Step: `golangci-lint (strict, full suite)`
- **Summary:** store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
```text
Lint	golangci-lint (strict, fail on any error)	﻿2026-10-02T02:28:07.7010278Z ##[group]Run golangci-lint run --timeout=5m --issues-exit-code=1
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7010577Z ^[[36;1mgolangci-lint run --timeout=5m --issues-exit-code=1^[[0m
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7063837Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064040Z env:
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064175Z   FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064337Z   GOTOOLCHAIN: local
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:07.7064472Z ##[endgroup]
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7228010Z ##[error]store/migrations.go:302:46: SA4023: this comparison is always true (staticcheck)
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7234776Z 	if err := insertRemappedGroupRepos(dbConn); err != nil {
Lint	golangci-lint (strict, fail on any error)	2026-10-02T02:28:28.7235352Z 	                                            ^
```

### [3] Commit `a98597d` (Branch: `main`)

- **Workflow:** CI
- **Run ID:** [36955109204](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36955109204)
- **Recorded At:** 2026-10-02T02:19:43Z (Duration: 303s)

#### Job: `CI` | Step: `Failed Step`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
    # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
    # [github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate]
vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath))
    cmdupdate/exports.go
    + L1 [typecheck] : # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate [github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate.test]
cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
    ========================================================================
[typecheck] could not import github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate (-: # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)) (NEW in /tmp/lint-current/report.json)
    ❌ cmd/clihelpers.go:35: [typecheck] could not import github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate (-: # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
[typecheck] : # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate [github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate.test]
cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath) (NEW in /tmp/lint-current/report.json)
    ❌ cmdupdate/exports.go:1: [typecheck] : # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate [github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate.test]
Process completed with exit code 1.

❌ FAIL: Found 8 nested-if / anti-compression violation(s) across 4 file(s):
    cli/cmd/nodes_cmd_test.go:93: Nested if statement found (depth 2 inside conditional block): if strings.Contains(between, ) {
    cli/cmdpull/pull_efficient_render.go:438: Nested if statement found (depth 2 inside conditional block): if relErr == nil && !strings.HasPrefix(rel, ) {
    cli/cmdpull/pull_remediation_hint.go:186: Nested if statement found (depth 2 inside conditional block): if hasOnlyUntracked {
    cli/cmdupdate/update_fleet.go:878: Nested if statement found (depth 2 inside conditional block): if agmData, agmErr := locatePackageBinary(); agmErr == nil && len(agmData) > 0 {
    cli/cmdupdate/update_fleet.go:918: Nested if statement found (depth 2 inside conditional block): if data, readErr := os.ReadFile(execPath); readErr == nil && len(data) > 0 {
    cli/cmdupdate/update_fleet.go:923: Nested if statement found (depth 2 inside conditional block): if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {
    cli/cmdupdate/update_fleet.go:932: Nested if statement found (depth 2 inside conditional block): if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {
    cli/cmdupdate/update_fleet.go:937: Nested if statement found (depth 2 inside conditional block): if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {
Process completed with exit code 1.

❌ FAILED: Found 3 error management violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/store/migrations.go:242: Swallowed database call error detected: '_, _ = dbConn.Exec(sqlPrune)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:86: Swallowed database call error detected: '_, _ = conn.Exec(`DELETE FROM SearchQueryLog`)'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/searcher/search_history_db.go:87: Swallowed database call error detected: '_, _ = conn.Exec(`VACUUM`)'
Process completed with exit code 1.

cmdupdate\update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
    Gitmap binary not found. Searched locations:
    - D:\a\gitmap-v28\gitmap-v28\bin\gitmap.exe
    - D:\a\gitmap-v28\gitmap-v28\gitmap.exe
Process completed with exit code 1.

❌ FAILED: Found 7 violation(s):
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_remediation_hint.go:186: Nested 'if' detected (depth 2): 'if hasOnlyUntracked {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdpull/pull_efficient_render.go:438: Nested 'if' detected (depth 2): 'if relErr == nil && !strings.HasPrefix(rel, ) {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdupdate/update_fleet.go:878: Nested 'if' detected (depth 2): 'if agmData, agmErr := locatePackageBinary(); agmErr == nil && len(agmData) > 0 {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdupdate/update_fleet.go:918: Nested 'if' detected (depth 2): 'if data, readErr := os.ReadFile(execPath); readErr == nil && len(data) > 0 {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdupdate/update_fleet.go:923: Nested 'if' detected (depth 2): 'if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdupdate/update_fleet.go:932: Nested 'if' detected (depth 2): 'if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {'
    - /home/runner/work/gitmap-v28/gitmap-v28/cli/cmdupdate/update_fleet.go:937: Nested 'if' detected (depth 2): 'if data, readErr := os.ReadFile(lp); readErr == nil && len(data) > 0 {'
Process completed with exit code 1.

❌ FAIL: Found 26 absolute path / URI violation(s):
    02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md:188: Hardcoded absolute repo path (D:\work\gitmap): - `d:\work\gitmap` vs `D:\work\gitmap`
    02-spec/21-app/65-gitmap-update-all-zip-and-fixes/01-architecture-spec.md:189: Hardcoded absolute repo path (D:\work\gitmap): - `d:\work\gitmap\` vs `d:\work\gitmap`
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:49: Hardcoded absolute repo path (D:\work\gitmap): SQLite text columns default to `BINARY` collation. Under `BINARY` collation, string comparisons eval
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/01-architecture-spec.md:418: Hardcoded absolute repo path (D:\work\gitmap): - Given a database containing duplicate rows (`RepoId=1, AbsolutePath="D:\work\gitmap"` and `RepoId=
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Absolute file:/// URI with drive letter: - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:58: Hardcoded absolute repo path (D:\work\gitmap): - **Primary Source:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Absolute file:/// URI with drive letter: - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:59: Hardcoded absolute repo path (D:\work\gitmap): - **Helpers:** [helpers.go](file:///d:/work/gitmap/cli/cmdpull/helpers.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:81: Hardcoded absolute repo path (D:\work\gitmap): | `D:\work\gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:82: Hardcoded absolute repo path (D:\work\gitmap): | `d:\work\gitmap\` | `d:\work\gitmap` | `d:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:83: Hardcoded absolute repo path (D:\work\gitmap): | `D:/work/./gitmap` | `D:\work\gitmap` | `D:/work/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:84: Hardcoded absolute repo path (D:\work\gitmap): | `d:/WORK/gitmap` | `d:\WORK\gitmap` | `d:/WORK/gitmap` | `d:/work/gitmap` |
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Absolute file:/// URI with drive letter: - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:180: Hardcoded absolute repo path (D:\work\gitmap): - **Concurrency Orchestration:** [pull_concurrency.go](file:///d:/work/gitmap/cli/cmdpull/pull_concu
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Absolute file:/// URI with drive letter: - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:181: Hardcoded absolute repo path (D:\work\gitmap): - **Issue Aggregation & Remediation:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Absolute file:/// URI with drive letter: - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:271: Hardcoded absolute repo path (D:\work\gitmap): - **Summary Render Pipeline:** [pull_efficient_render.go](file:///d:/work/gitmap/cli/cmdpull/pull_ef
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Absolute file:/// URI with drive letter: - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:272: Hardcoded absolute repo path (D:\work\gitmap): - **Batch Output Coordinator:** [pull.go](file:///d:/work/gitmap/cli/cmdpull/pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Absolute file:/// URI with drive letter: - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:364: Hardcoded absolute repo path (D:\work\gitmap): - **Safe Pull Execution:** [safe_pull.go](file:///d:/work/gitmap/cli/cloner/safe_pull.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Absolute file:/// URI with drive letter: - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:365: Hardcoded absolute repo path (D:\work\gitmap): - **Pull Diagnostics:** [pulldiag.go](file:///d:/work/gitmap/cli/cloner/pulldiag.go)
    02-spec/21-app/66-fix-gitmap-pa-duplicate-repos/02-component-spec.md:443: Hardcoded absolute repo path (D:\work\gitmap): - Verify `D:\work\gitmap`, `d:\work\gitmap`, `D:/work/gitmap`, and `d:\work\gitmap\` all resolve to
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
grep -E -A 5 '^--- FAIL:' /tmp/json-snapshot.log || true
    echo "========================================="
    fi
    exit "$exit_code"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
    # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmd [build failed]
    FAIL
    =========================================
    JSON SNAPSHOT / SCHEMA REGRESSION
    A startup-list / find-next / latest-branch JSON
    contract test failed. Likely causes:
    1. Encoder shape changed (added/renamed/reordered field)
    2. Key ordering drifted (map iteration leaked in)
    3. Schema registry vN.json out of date
    To accept an intentional schema change:
    GITMAP_UPDATE_SCHEMA=<name> go test ./cmd/...
    To regenerate golden fixtures:
    GITMAP_UPDATE_GOLDEN=1 go test ./cmd/...
Process completed with exit code 1.

❌ FAIL: Found 1 boolean guideline violation(s) across 1 file(s):
    cli/cmdupdate/update_fleet.go:904: Banned function prefix returning bool (shouldIncludeAgmInZip): func shouldIncludeAgmInZip(pkg string) bool {
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
# Tally PASS/FAIL package lines for the PR summary comment.
    passed=$(grep -cE '^ok[[:space:]]' /tmp/full-suite/test-output.txt || true)
    failed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)
    echo "packages_passed=$passed" >> "$GITHUB_OUTPUT"
    echo "packages_failed=$failed" >> "$GITHUB_OUTPUT"
    exit "$rc"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
    # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
FAIL	github.com/alimtvnetwork/gitmap-v28/cli [build failed]
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmd [build failed]
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate [build failed]
--- FAIL: TestEveryHelpFileHasExamples (0.01s)
    examples_golden_test.go:42: help files missing `## Examples` section (2):
    - devtool.md
    - pulle.md
    FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/helptext	0.017s
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/tests/cmd_test [build failed]
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/tests/heavy_test [build failed]
Process completed with exit code 1.

cmd/clihelpers.go:35:2: could not import github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate (-: # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)) (typecheck)
    "github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate"
    ^
    cmdupdate/exports.go:1: : # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate [github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate.test]
cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath) (typecheck)
    package cmdupdate
Process completed with exit code 1.

FAIL: test_smoke_installer_source_mode (__main__.TestSmokeInstaller.test_smoke_installer_source_mode)
    ----------------------------------------------------------------------
Traceback (most recent call last):
    File "/home/runner/work/gitmap-v28/gitmap-v28/.github/scripts/tests/test_ci_scripts.py", line 157, in test_smoke_installer_source_mode
    self.assertEqual(res.returncode, 0)
AssertionError: 3 != 0
    Ran 18 tests in 56.713s
FAILED (failures=1)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.

cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Process completed with exit code 1.
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Go vet`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Diff vs baseline (fail only on NEW findings)`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath))
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run ./.github/actions/policy-check`
- **Summary:** ❌ FAIL: Found 1 boolean guideline violation(s) across 1 file(s):
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run E2E CLI Smoke Tests Windows`
- **Summary:** cmdupdate\update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run JSON snapshot fast tests`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run full test suite`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `golangci-lint (strict, full suite)`
- **Summary:** cmd/clihelpers.go:35:2: could not import github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate (-: # github.com/alimtvnetwork/gitmap-v28/cli/cmdupdate
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run lint-script unit tests`
- **Summary:** FAIL: test_smoke_installer_source_mode (__main__.TestSmokeInstaller.test_smoke_installer_source_mode)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Run E2E CLI Smoke Tests`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

#### Job: `CI` | Step: `Build binary`
- **Summary:** cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
```text
Lint	Go vet	2026-10-02T02:20:45.6291565Z ##[error]cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.0376458Z ##[error]vet: cmdupdate/update_fleet.go:438:16: h.KeyPath undefined (type store.SSHHost has no field or method KeyPath)
Lint	Go vet	2026-10-02T02:20:46.1061262Z ##[error]Process completed with exit code 1.
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	﻿2026-10-02T02:22:32.2451415Z ##[group]Run BASELINE=/tmp/lint-baseline/report.json
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2451886Z ^[[36;1mBASELINE=/tmp/lint-baseline/report.json^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452298Z ^[[36;1m# When disabled via dispatch, force seeding mode by pointing at^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2452770Z ^[[36;1m# an empty path so lint-diff.py treats it as "no baseline".^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453149Z ^[[36;1mif [ "$BASELINE_DISABLED" = "true" ]; then^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2453629Z ^[[36;1m  echo "::notice::Baseline cache disabled by workflow_dispatch input — seeding mode."^[[0m
Lint Baseline Diff	Diff vs baseline (fail only on NEW findings)	2026-10-02T02:22:32.2454072Z ^[[36;1m  BASELINE=""^[[0m
```

### [4] Commit `9da3535` (Branch: `main`)

- **Workflow:** CI
- **Run ID:** [36953655930](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36953655930)
- **Recorded At:** 2026-10-02T02:01:08Z (Duration: 397s)

#### Job: `CI` | Step: `Failed Step`
- **Summary:** Process completed with exit code 2.

# Tally PASS/FAIL package lines for the PR summary comment.
    passed=$(grep -cE '^ok[[:space:]]' /tmp/full-suite/test-output.txt || true)
    failed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)
    echo "packages_passed=$passed" >> "$GITHUB_OUTPUT"
    echo "packages_failed=$failed" >> "$GITHUB_OUTPUT"
    exit "$rc"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
--- FAIL: TestEveryHelpFileHasExamples (0.01s)
    examples_golden_test.go:42: help files missing `## Examples` section (2):
    - devtool.md
    - pulle.md
    FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/helptext	0.014s
Process completed with exit code 1.
```text
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.3661978Z   fail-on-cache-miss: false
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.6113035Z (node:2400) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:46.7388115Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:47.5999108Z ##[error]Process completed with exit code 2.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.4909843Z   fail-on-cache-miss: false
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.6690694Z (node:2392) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164193Z ##[group]Run set -euo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164598Z ^[[36;1mset -euo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2508660Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2703880Z ##[group]Run set -uo pipefail
```

#### Job: `CI` | Step: `Run misspell on changed files`
- **Summary:** Process completed with exit code 2.
```text
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.3661978Z   fail-on-cache-miss: false
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.6113035Z (node:2400) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:46.7388115Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:47.5999108Z ##[error]Process completed with exit code 2.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.4909843Z   fail-on-cache-miss: false
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.6690694Z (node:2392) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164193Z ##[group]Run set -euo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164598Z ^[[36;1mset -euo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2508660Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2703880Z ##[group]Run set -uo pipefail
```

#### Job: `CI` | Step: `Run full test suite`
- **Summary:** # Tally PASS/FAIL package lines for the PR summary comment.
```text
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.3661978Z   fail-on-cache-miss: false
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:45.6113035Z (node:2400) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:46.7388115Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Spell Check (misspell, US locale)	UNKNOWN STEP	2026-10-02T02:01:47.5999108Z ##[error]Process completed with exit code 2.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.4909843Z   fail-on-cache-miss: false
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:01:55.6690694Z (node:2392) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164193Z ##[group]Run set -euo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2164598Z ^[[36;1mset -euo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2508660Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Full Suite Guard	UNKNOWN STEP	2026-10-02T02:02:00.2703880Z ##[group]Run set -uo pipefail
```

### [5] Commit `641a754` (Branch: `main`)

- **Workflow:** CI
- **Run ID:** [36920857215](https://github.com/alimtvnetwork/gitmap-v28/actions/runs/36920857215)
- **Recorded At:** 2026-10-01T20:20:31Z (Duration: 375s)

#### Job: `CI` | Step: `Failed Step`
- **Summary:** # Tally PASS/FAIL package lines for the PR summary comment.
    passed=$(grep -cE '^ok[[:space:]]' /tmp/full-suite/test-output.txt || true)
    failed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)
    echo "packages_passed=$passed" >> "$GITHUB_OUTPUT"
    echo "packages_failed=$failed" >> "$GITHUB_OUTPUT"
    exit "$rc"
    shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
    env:
    FORCE_JAVASCRIPT_ACTIONS_TO_NODE24: true
    GOTOOLCHAIN: local
--- FAIL: TestEveryHelpFileHasExamples (0.01s)
    examples_golden_test.go:42: help files missing `## Examples` section (2):
    - devtool.md
    - pulle.md
    FAIL
FAIL	github.com/alimtvnetwork/gitmap-v28/cli/helptext	0.011s
Process completed with exit code 1.
```text
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:11.5458597Z   fail-on-cache-miss: false
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:11.8675607Z (node:2389) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:18.9210370Z ##[group]Run set -euo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:18.9210769Z ^[[36;1mset -euo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1684080Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1850891Z ##[group]Run set -uo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1851235Z ^[[36;1mset -uo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1853430Z ^[[36;1m# Tally PASS/FAIL package lines for the PR summary comment.^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1854738Z ^[[36;1mfailed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1855654Z ^[[36;1mecho "packages_failed=$failed" >> "$GITHUB_OUTPUT"^[[0m
```

#### Job: `CI` | Step: `Run full test suite`
- **Summary:** # Tally PASS/FAIL package lines for the PR summary comment.
```text
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:11.5458597Z   fail-on-cache-miss: false
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:11.8675607Z (node:2389) [DEP0169] DeprecationWarning: `url.parse()` behavior is not standardized and prone to errors that have security implications. Use the WHATWG URL API instead. CVEs are not issued for `url.parse()` vulnerabilities.
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:18.9210370Z ##[group]Run set -euo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:18.9210769Z ^[[36;1mset -euo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1684080Z shell: /usr/bin/bash --noprofile --norc -e -o pipefail {0}
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1850891Z ##[group]Run set -uo pipefail
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1851235Z ^[[36;1mset -uo pipefail^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1853430Z ^[[36;1m# Tally PASS/FAIL package lines for the PR summary comment.^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1854738Z ^[[36;1mfailed=$(grep -cE '^FAIL[[:space:]]' /tmp/full-suite/test-output.txt || true)^[[0m
Full Suite Guard	UNKNOWN STEP	2026-10-01T20:21:19.1855654Z ^[[36;1mecho "packages_failed=$failed" >> "$GITHUB_OUTPUT"^[[0m
```

## AI Agent Directives to Prevent Mistakes

1. **Pre-Commit Verification:** Run `python 03-ai-scripts/05-guideline-autofixer.py <dir>` before staging code.
2. **Strict Sizing Rules:** Keep Go functions $\le 15$ lines and files $\le 100$ lines (Rule R14).
3. **Positive Boolean Hygiene:** Use affirmative identifiers (`isSuccess`, `hasErrors`, `isReady`); reject double negatives.
4. **Universal Error Handling:** Always wrap and return errors using `*appfault.AppError` / `apperror.WrapSimple`.
5. **Cross-Platform Safety:** Ensure all path operations use `filepath.ToSlash` and tests do NOT rely on unmocked OS commands.

