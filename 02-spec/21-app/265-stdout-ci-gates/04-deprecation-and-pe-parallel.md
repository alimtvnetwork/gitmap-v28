# 265 — Command Deprecation Policy + `pipeline errors all` Parallelism + Credential Audit Reference

Status: spec (draft) · Date: 2026-10-10 · Owner task: `265-stdout-ci-gates`

## Part A — Command deprecation policy (mirror the flag-alias pattern)

### Current state (verified by research — do not re-derive)

- A deprecation pattern EXISTS for FLAG aliases: `Msg*` constants in `cli/constants/constants_cli.go` (e.g. `--concurrency is deprecated; use --workers`), a stderr one-liner emitted at flag-resolution time (`cli/cmdscan/flags.go:152`), and tests asserting the contract.
- NO command-deprecation mechanism exists.

### Policy

1. **Message constant:** add a `Msg*` constant in `cli/constants/constants_cli.go` per deprecated command, e.g. `MsgFixSubcommandDeprecated = "fix %q is deprecated and will be removed in v%s; use %s"`.
2. **Dispatch-time stderr one-liner:** emit the warning at command dispatch time — the command equivalent of `cli/cmdscan/flags.go:152`. One line on stderr, then continue executing the old path unchanged.
3. **Two-release warning window:** a deprecated command warns for two minor releases, then is removed (or becomes a permanent alias).
4. **Aliases never break:** an alias is a permanent promise. Deprecation applies only to commands/subcommands that have a named successor — never to pure aliases.
5. **Tests:** every deprecation warning has a test asserting (a) the stderr line is emitted exactly once, and (b) the command still executes its old behavior while deprecated.

### Retroactive application: `fix` unknown-subcommand becomes a deprecation redirect

Current behavior: `cli/cmdautofix/fix.go:43-47` — an unknown `fix` subcommand produces a usage error mentioning `stash`/`wip`/`discard`; routed at `cli/cmd/rootutility.go:477`.

Change: convert that error path into the deprecation-warning pattern — emit the `Msg*` warning naming the real subcommands (`stash`, `wip`, `discard`) as the intended targets, then route to the suggested subcommand (or show its usage). The old bare-error behavior is the bug being fixed: users typing a displaced subcommand name get a deprecation-style redirect instead of a dead-end usage error.

## Part B — `pipeline errors all` parallelism (pe-all worker pool)

### Current state (verified by research — do not re-derive)

- `collectAllPipelineSummary()` in `cli/cmdpipeline/pipeline_all_errors.go` has a SEQUENTIAL per-repo loop.
- Worker-pool precedent exists: `executeParallelFetchWorkers` in `cli/cmdpipeline/pipeline_logs_fetch.go` (semaphore channel, currently cap 8).
- `--workers` flag pattern: `cli/cmdpipeline/pipeline_flags.go` (parsed like `Limit`).
- GitHub last-commit pattern: `gh run list --commit <sha> --json` at `cli/cmdpipeline/pipeline_logs_target.go:152`; alternatively direct REST with a Bearer token from `secrets.Resolve()`.

### Design

1. **Worker pool** in `collectAllPipelineSummary()`: semaphore-channel pattern mirroring `executeParallelFetchWorkers`.
2. **Worker count:** `workers = max(runtime.NumCPU(), 3)`; a `--workers` flag overrides it, parsed like `Limit` in `pipeline_flags.go`.
3. **No nested parallelism inside workers:** each worker handles exactly one repo, strictly sequentially. (One level only — the 243b consolidation incident is the cautionary tale.)
4. **Last-commit-only fetch per repo:** use the `gh run list --commit <sha> --json` pattern (`pipeline_logs_target.go:152`) or direct REST with Bearer from `secrets.Resolve()` — fetch only the latest commit's run, never full history.
5. **SQLite hash cache:** per repo, compare the recorded commit hash in the pipeline DB against the current HEAD hash — match means skip the fetch and mark the repo done; mismatch means fetch + record the new hash.
6. **Deterministic output:** re-sort results by repo slug after the parallel fetch; guard shared summary counters with a mutex.

### Acceptance

- `gitmap pipeline errors all` on a multi-repo catalog completes faster than the sequential baseline (record both timings in the subtask file).
- Output order is identical across runs (slug-sorted).
- `--workers 1` reproduces sequential behavior.
- Hash cache: a second run with no new commits performs zero fetches.

## Part C — Credential audit findings (reference only)

Findings from the credential audit, embedded here as reference for the follow-up program. **Findings only — fixes are a separate follow-up program; no remediation code in this task.**

| Severity | Finding | Location | Why it matters |
|---|---|---|---|
| HIGH | GitHub PAT stored in plaintext in `~/.gitconfig` (0644) | `cli/cmdlogin/login_token.go:85` | any local process/user can read the token |
| HIGH | Token persisted in audit DB `Args`/`Flags` plus raw command line in `history/commands.db` | `cli/cmdaudit/audit_db.go:77` | secrets leak into queryable history |
| HIGH | SMTP password + Telegram bot token in plaintext in `~/.gitmap/variables.json` (0644) | `cli/cmdagy/agy_telegram_email_settings.go:256` | world-readable credential store |
| MEDIUM | Fleet token deployment via shell interpolation | see coordinator's ledger | token visible in process list / shell history |
| MEDIUM | Hardcoded AES fallback keys for SSH passwords | see coordinator's ledger | fallback keys defeat the RSA vault |

The follow-up program (not this task) owns remediation: encrypted-at-rest storage, redaction in audit args, and permission tightening.
