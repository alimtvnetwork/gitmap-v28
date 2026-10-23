# Subtask 311-02: CI/CD Fixes

## Objective

Fix every root cause found in subtask 311-01 at the source so the full CI
suite passes. No check may be silenced, skipped, or baselined away unless the
baseline/entry itself is proven to be the defect.

## Known Fix Queue (from baseline annotations; confirm against 311-01)

- Delete verified-dead functions (callers checked via gitmap aum search
  first): recordVacuumTelemetry, handleVacuumError, executeVacuumOnDB,
  vacuumSingleTarget (cli/cmddb/db_vacuum.go — note cli/cmddb/ is in the
  spec-310 dirty set; coordinate before editing), requireOnlineCheck
  (cli/cmd/helpcheck.go), hasCredentialPersistenceFailure
  (cli/cloner/safe_pull.go), formatFolderFileURI, hasSkipIDEsFlag,
  parseGitHubSlugFromURL (locate first).
- Misspell: DOUB -> DOUBT at the source occurrence.
- Legacy Refs: reword 02-spec/01-spec-authoring-guide/13-root-readme-conventions.md
  lines 20 and 257 to drop /gitmap-v[567]\b/ references.
- Repair the failing unit/race tests at their underlying defect.
- Repair the failing linter-script unit test at its source.
- Relative Path failures: fix only files outside the spec-310 dirty set;
  leave its files to that stream.

## Rules

- Respect the 300-line file-size gate on every touched .go file.
- Respect the nested-if and boolean/enum linters on every touched file
  (they scan whole touched files; a touched file's pre-existing violations
  become this commit's problem).
- gofmt-clean on every touched .go file.
- Verify imports with bare-pattern aum searches before removing any import
  (escaped patterns can silently miss usages).

## Exit Criteria

- A push whose GitHub check-runs show the previously failing checks green,
  with no new failures introduced.
