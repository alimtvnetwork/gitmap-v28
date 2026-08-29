# CI/CD Issues Index

Tracks every CI/CD pipeline failure or hardening decision encountered, its root cause, and resolution. New entries go in `.lovable/cicd-issues/XX-short-name.md` with sequential numeric prefixes.

## Conventions
- File naming: `XX-kebab-case-name.md` (XX = zero-padded sequence starting at `01`).
- One file per distinct issue. Do **not** duplicate — if the same root cause recurs, append a "Recurrence" section to the existing file.
- Status values: `✅ Resolved`, `🔄 In Progress`, `⏳ Pending`, `🚫 Blocked`.

## Issues

| # | Title | Tool / Stage | Status | File |
|---|-------|--------------|--------|------|
| 01 | misspell: `labeled` → `labeled` | golangci-lint (misspell) | ✅ Resolved | [01-misspell-labelled.md](cicd-issues/01-misspell-labelled.md) |
| 02 | `lint-regression-guard` → `lint-hard-floor` → `lint-baseline-guard` (now uniformly baseline-diff for all 5 linters) | golangci-lint (baseline-guard job) | ✅ Resolved | [02-lint-regression-guard-semantics.md](cicd-issues/02-lint-regression-guard-semantics.md) |
| 03 | pterm SpinnerPrinter DATA RACE | go test -race | ✅ Resolved | [03-pterm-spinner-data-race.md](cicd-issues/03-pterm-spinner-data-race.md) |
| 04 | Legacy SQLite Migration Failures & `cmdFaithful` Concurrency Race | go test -race / Test Matrix | ✅ Resolved | [04-zipgroupitems-migration-and-cmdfaithful-race.md](cicd-issues/04-zipgroupitems-migration-and-cmdfaithful-race.md) |
| 05 | Cluster TLS Dial Timeout, Relative ModuleRoot in Subprocess Tests & Test Env Race | go test -race / Test Matrix | ✅ Resolved | [05-cluster-tls-dial-timeout-and-test-env-race.md](cicd-issues/05-cluster-tls-dial-timeout-and-test-env-race.md) |
| 06 | LookPath Injection in Coding Guidelines & Cross-Platform Fallback Assertion | go test -race / Test Matrix | ✅ Resolved | [06-codingguidelines-lookpath-injection.md](cicd-issues/06-codingguidelines-lookpath-injection.md) |
| 07 | installctx `ctxExplainEnabled` Concurrency Race Guard | go test -race / Test Matrix | ✅ Resolved | [07-installctx-explain-concurrency-race.md](cicd-issues/07-installctx-explain-concurrency-race.md) |
| 08 | clonepretty Global Flag Atomic Synchronization | go test -race / Test Matrix | ✅ Resolved | [08-clonepretty-atomic-state-synchronization.md](cicd-issues/08-clonepretty-atomic-state-synchronization.md) |
| 09 | Test $HOME Mutation Isolation in Tilde Expansion Tests | go test -race / Test Matrix | ✅ Resolved | [09-home-env-race-in-tilde-tests.md](cicd-issues/09-home-env-race-in-tilde-tests.md) |
| 10 | VS Code Project Manager Cross-Platform Slash Normalization | go test / Test Matrix | ✅ Resolved | [10-vscodepm-cross-platform-path-normalization.md](cicd-issues/10-vscodepm-cross-platform-path-normalization.md) |
| 11 | cmdFaithfulExiter Mutex Isolation & maybeExitOnCmdFaithfulMismatch Race | go test -race / Test Matrix | ✅ Resolved | [11-clonetermverifyexit-process-exit-race.md](cicd-issues/11-clonetermverifyexit-process-exit-race.md) |
| 12 | Query Wrapper Property Mismatch & E2E Typos | Audit / Test Matrix | ✅ Resolved | [12-query-wrapper-isfail-typo.md](cicd-issues/12-query-wrapper-isfail-typo.md) |
| 32 | Panic on Zero-Args & Legacy Exit Refactor | CLI Runtime / Unit Tests | ✅ Resolved | [32-root-cli-panic-on-zero-args.md](cicd-issues/32-root-cli-panic-on-zero-args.md) |
| 33 | Gocritic appendAssign Diff Failure & Unused Helpers | CI Diff Gate / golangci-lint | ✅ Resolved | [33-gocritic-appendassign-diff.md](cicd-issues/33-gocritic-appendassign-diff.md) |
| 34 | Gocritic ifElseChain Diff Failure & Switch Refactor | CI Diff Gate / golangci-lint | ✅ Resolved | [34-gocritic-ifelsechain-diff.md](cicd-issues/34-gocritic-ifelsechain-diff.md) |
| 35 | Exhaustive Switch Diff Failure & Complete Enum Handlers | CI Diff Gate / golangci-lint | ✅ Resolved | [35-exhaustive-switch-diff.md](cicd-issues/35-exhaustive-switch-diff.md) |
| 36 | Misspell Changed Files & US English Standardization | Spell Check / misspell | ✅ Resolved | [36-misspell-changed-diff.md](cicd-issues/36-misspell-changed-diff.md) |
| 37 | Installer Smoke Python Runner & Strict Relative Paths | Installer Smoke / Policy Gates | ✅ Resolved | [37-installer-smoke-release-diff.md](cicd-issues/37-installer-smoke-release-diff.md) |

## Patterns Learned
- **US-English everywhere in Go**: `misspell` flags British spellings in comments and identifiers. Avoid `labeled`, `canceled`, `behavior`, `color`, `occurred`, `receive`, `separate`.
- **Pinned linter versions**: golangci-lint is pinned to `v1.64.8`; `goimports` to `v0.24.0`; `govulncheck` to `v1.1.4`. Do not assume newer rules.
- **ARIA attributes are exempt**: `aria-labelledby` is a standard HTML/ARIA token and must never be "corrected".
- **Decide hard-floor vs baseline-diff per linter** before adding to CI; document choice in the script header. Don't conflate them under one job.
- **Compute `-local` prefix from `go.mod`** — never hardcode the module path in CI scripts.
- **Cache `GOMODCACHE`/`GOCACHE`** keyed on `go.sum` hash for compile-gate + matrix speed.
- **Conditional UI Component Allocation**: Always wrap background UI routines (e.g. `pterm.SpinnerPrinter`) in `isMultiActive` checks. Unconditional start leads to `-race` failures when tests disable output.
- **Guard Legacy Schema Migrations with Table Existence**: In SQLite database migrations, check `db.tableExists("Table")` before issuing `ALTER TABLE` or `UPDATE` on legacy tables, and execute benign-tolerant DDL statements via `db.conn.Exec` rather than `ExecWrapper` to prevent stderr error spam.
- **Synchronize Package Globals and Mocks**: Use `atomic.Bool` and `sync.RWMutex` on package-level state and test hooks to guarantee race-free execution during `go test -race` runs.
- **Bound Network Dialing with Timeouts**: Always use `tls.DialWithDialer` or `net.DialTimeout` with bounded connection deadlines (e.g. 500ms - 2s) to prevent tests and commands from hanging on OS TCP SYN retransmission loops when connecting to unreachable nodes.
- **Use `t.Setenv` Exclusively in Tests**: Never use raw `os.Setenv` in unit/integration test helpers; `t.Setenv` ensures test-isolated and race-safe environment restoration.
- **Derive Absolute Module Roots in Test Subprocess Builders**: Never use relative paths (like `cmd.Dir = ".."`) when building binaries inside `go test` helpers; concurrent `os.Chdir()` calls will corrupt relative paths. Use `runtime.Caller(0)` to obtain an immutable absolute path.
- **Inject Executable Lookups Instead of Mutating `$PATH`**: Never zero out `$PATH` with `t.Setenv("PATH", "")` when packages run parallel tests (`t.Parallel()`); provide injectable `LookPath func(string) (string, error)` seams to test missing binary paths safely.
- **Protect Feature Flag Globals with RWMutex / atomic.Bool**: Wrap feature-flag globals like `ctxExplainEnabled` with `sync.RWMutex` and dry-run/spinner flags with `atomic.Bool` to eliminate data races during concurrent test suites.
- **Avoid Stale Environment Snapshots in Assertions**: Validate structural invariants (prefix/suffix removal) rather than comparing against static `os.UserHomeDir()` snapshots that may be mutated by concurrent `t.Setenv` test fixtures.
- **Normalize Slashes with `filepath.ToSlash` for Cross-Platform Path Keys**: Always normalize backslashes before comparing filesystem keys across operating systems, since POSIX treats `\` as a regular character.
- **Lock Out Exit-Triggering Hooks Exclusively During Tests**: When overriding `os.Exit` with a test stub, hold an exclusive mutex lock across the entire test to prevent concurrent tests from evaluating the exit trigger while the global handler transitions.
- **Never use `*_linux.go` for multi-OS POSIX fallbacks**: In Go, file suffix `*_linux.go` forces `GOOS=linux` even with `//go:build !windows`. Always use `*_posix.go` or `*_unix.go` with `//go:build !windows` to include macOS (`darwin`), FreeBSD, etc.
