# Architecture Spec — GitMap and Antigravity CI/CD and Code Refactoring

Slug: gitmap-and-antigravity-cicd-and-code-refactoring
Sequence: 325
Owner repos (only): gitmap-v28, Antigravity-Manager, coding-guidelines-v24 (prompt update only)
Status: active — overnight run, user asleep, no questions; every ambiguity decided by the coordinator and recorded here.

## User Request (Verbatim)

The user ordered, as a high-priority non-negotiable overnight task: run `gitmap pa` after updating gitmap; get CI/CD green on Antigravity tools (Antigravity-Manager) and on gitmap; release each with a minor bump once green (bump via the repo bump script, tag via gitmap, release via `gitmap release`); and complete the enum/constants refactor across the whole gitmap repository, group by group and package by package, following the coding guidelines. Named defect exemplars: string `+` concatenation in cli_cmd_cargo/cargo_exe.go; a regular error returned by applyDefaultPragmas in the constants section; `==` flag-comparison chains in cli/cmdpool/pool_cli_input.go around line 79; non-lazy regex; sqlite_config.go naming (Database/Db) and the wrapped-return problem (two-value returns must become one class type, generic and also defined without generics so it looks flat). The enum CG execute prompt in coding-guidelines-v24 must strictly state the enum-per-package + type convention rule, with the cargo_exe.go example called out. The user checks in the morning; the run must proceed with no questions and no waiting.

## Decisions Made on Ambiguities (coordinator, no user input)

1. `gitmap space update` does not exist as an update command (`gitmap space` is common/backup-branch/send-email). Decoded per the user's shorthand rule ("space X" = "gitmap X") to `gitmap update`.
2. `gitmap update` itself fails at binary placement in this container (E1151; the swap step fails, downloads and SHA256 verify fine). Decision: manual placement of the verified v6.608.0 release asset (SHA256 matched checksums.txt), old binary kept as a backup. gitmap now reports v6.608.0.
3. Spec/plan numbering: existing 02-spec/21-app folders use slugs and .ai-memory/plans/pending uses numbers (latest 324). Decision: sequence 325, folder 02-spec/21-app/325-gitmap-and-antigravity-cicd-and-code-refactoring, plan .ai-memory/plans/pending/325-gitmap-and-antigravity-cicd-and-code-refactoring.md.
4. AGM version target in the original order (4.256.0 → 4.257.0) is stale: another session already shipped AGM up to 4.266.0. Decision: bump from whatever the current tip version is at bump time (4.266.0 → 4.267.0 unless the tip moved again).
5. gitmap version target (6.600.0 → 6.601.0) is likewise stale: the other session shipped up to v6.608.0. Decision: bump from the current tip version at bump time.
6. The other user session was actively pushing to gitmap-v28 and Antigravity-Manager during this run. Decision: never fight it — always fetch, work on the latest tip, rebase/ff-forward before each push, work on distinct file groups, and if a release lands from that session while a worker is mid-bump, the worker re-bases the bump onto the new tip instead of duplicating a release.
7. Skill conflict: execute-parent-task-with-n-steps-v6 prescribes one atomic commit at the end, zero builds/tests, zero CI verification, and zero releases unless asked. The parent task explicitly orders incremental commits, CI verification, and releases. Decision: parent instructions outrank the skill on all three points; the skill's swarm shape (3 workers, disjoint file boxes, GitMap-only search, relative paths, hyphen commit format, user identity) is still followed.
8. Skill conflict: the skill bans worker git commands. Workers in this run each own a different repo (no shared index), and the parent orders workers to execute including commits via the gitmap suite. Decision: workers may run gitmap cpf/cpb/release and read-only git commands inside their own repo only; the coordinator does all cross-repo verification and the final report.
9. If the whole-repo enum refactor cannot finish in one night, that is acceptable only if CI is green on both repos, both releases are done, and the remaining groups are listed precisely. Never claim completion that does not exist.

## Architecture

Three workstreams, disjoint by repo and by file group:

1. Antigravity-Manager CI-green + release. Scope is CI-green plus release only. No AGM-wide enum refactor. The repo's own bump script bumps the version; `gitmap release` creates the tag and the GitHub release; verification is tag on the remote via git ls-remote plus release published with assets and marked Latest. AGM main is protected by repo ruleset 23609023 (deletion + non_fast_forward, no bypass actors): force-push is rejected, so no history rewrites, ever.
2. gitmap-v28 CI-green + enum/constants refactor + release. CI-green first, then refactor waves group by group with incremental commits and a CI check per group or per batch, then the release. The other session already landed several "Refactor: CONSTANTS - dry restructuring of enums and constants" commits (6dd000dc, 656b1257, 6987f0e2 and earlier); the audit must measure what is already done before planning waves, so no group is refactored twice.
3. coding-guidelines-v24 prompt update. The enum CG execute prompt must strictly state: every enum gets its own package, following the type convention; the cargo_exe.go concatenation example is called out as the anti-pattern. Mirrored to .agents/skills and .cursor/skills mirrors as applicable. Committed via the gitmap suite as the user.

## Invariants (non-negotiable)

- Git identity MD ALIM UL KARIM <devorg.bd@gmail.com>; commits only via gitmap cpf/cpb/cpc/cpr/pcp, hyphen format, no colons in the message argument; never Muse <muse@local>; never raw git commit.
- Search only via GitMap tools (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train). Total ban on rg, ripgrep, grep, git grep, Select-String.
- Strictly relative Git paths everywhere, including release notes and release pages. Lowercase filenames.
- gitmap files stay under the 300-line cap; split when a refactor would exceed it.
- Verify via CI. Local go build + gofmt -l pre-flight is accepted; no local test suites.
- Track in-flight touches with gitmap rc heatmap status/tree/live/sync.
- No secrets in repos. No full-account cloning. No other repos.
