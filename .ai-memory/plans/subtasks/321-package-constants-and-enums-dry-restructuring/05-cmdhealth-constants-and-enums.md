# Subtask 321-05: cmdhealth Constants & Enums DRY Restructuring

## Target Scope
- Target Package: `cli/cmdhealth` (2 files)
- Leaf Packages to Create:
  - `cli/cmdhealth/enums/`
    - `health_enums.go`
    - `enums_test.go`
  - `cli/cmdhealth/constants/`
    - `constants.go`
    - `constants_test.go`

## Specific Enums Required
1. `HealthCategoryType`: Check domains (`git`, `pipeline`, `fleet`, `db`, `system`, `network`).
2. `HealthSeverityType`: Severity levels (`ok`, `warn`, `error`, `fatal`).
3. `HealthOutputFormatType`: Output format (`terminal`, `json`, `summary`, `markdown`).
4. `HealthCheckStatusType`: State indicators (`green`, `yellow`, `red`, `skipped`).

Each enum must follow:
- Type suffix: `...Type`
- Methods: `IsValid() bool`, `String() string`, `Parse*(s string) (Type, error)`, `All*() []Type`.
- Positive booleans only (`is*`, `has*`).
- 8-15 line function limits with strict blank line formatting.

## Specific Constants Required
- Command names & aliases: `CmdHealth = "health"`, `CmdHealthAlias = "hlth"`, environment variable `EnvGitmapWorkdir = "GITMAP_WORKDIR"`.
- CLI flags: `--workdir`, `--help`, `-h`.
- Thresholds: `MaxDirtyReposDisplay = 5`, `MaxRedPipelinesDisplay = 5`, `MaxPipelineProbedRepos = 50`, `StaleDbDaysCutoff = 30`.
- UI strings, delimiters, sentinel strings (`NoteNodeRegistryUnreachable`), follow-up triage recommendations.

## Consuming Callers to Refactor
- `cli/cmdhealth/health.go`
- `cli/cmdhealth/health_checks.go`

## Verification
- Unit tests in `cli/cmdhealth/enums/enums_test.go` and `cli/cmdhealth/constants/constants_test.go`.
- `go vet ./cli/cmdhealth/...`
- `go test ./cli/cmdhealth/...`
