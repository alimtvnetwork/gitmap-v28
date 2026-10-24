# Migration and Verification Spec — GitMap and Antigravity CI/CD and Code Refactoring

## Migration notes

- No data migrations. No schema changes except where a refactor group explicitly requires one (none planned).
- The other session's in-flight work is a moving target: both repos may gain commits during this run. Migration rule for every worker: fetch before starting, fetch before pushing, ff-forward or re-base onto the new tip, never force-push, never overwrite another session's commits.
- Enum moves are behavior-preserving by construction: same values, same String() output, same call sites (updated imports only). Any behavior change found during a wave stops that wave; the worker reports it instead of pushing it.
- The v6.600.0 tag stays broken and untouched. The gitmap recovery is a fresh minor release, not a repair of v6.600.0.
- The AGM v4.256.0 tag/history rewrite from the previous session is done and closed; this task performs no rewrites.

## Verification contract (evidence or it did not happen)

For each claim in the morning report, cite one of:

- CI: `gh run list --repo <owner>/<repo> --branch main` output showing every workflow success on the release commit, plus run URLs.
- Release: `git ls-remote origin "vX.Y.Z^{}"` output; `gh release view vX.Y.Z` showing published, not draft, asset count, Latest status, release URL.
- Refactor group: commit SHA (via gitmap suite), file list changed, gofmt + go build pre-flight exit 0, CI status on that commit.
- Prompt update: commit SHA, files changed including mirrors, quoted rule text.
- Audit: counts of enum/constant sites refactored vs remaining, grouped by package, with file lists for the remaining groups.

## Acceptance criteria

1. Antigravity-Manager: CI green on the release commit; minor-bump release published with tag verified on the remote and marked Latest.
2. gitmap-v28: CI green on the release commit; minor-bump release published with tag verified on the remote.
3. Enum/constants refactor: every group either completed with commit SHAs or listed precisely as remaining with file lists. No group claimed complete without a commit.
4. Prompt update committed with mirrors.
5. Morning report delivered to the parent with CI status per repo with run links, releases (tags, URLs, asset counts), refactor groups completed vs remaining with file lists, every commit SHA, and every ambiguity decision.

## Known risks

- The other session may push conflicting refactors to the same packages. Mitigation: disjoint file groups, RC heatmap tracking, fetch-before-push; on collision, the worker reports and the coordinator re-plans that group.
- CI runs take 10–30 minutes per round; refactor waves must batch sensibly so the night is not spent waiting on single-file rounds.
- The wrapped-return pattern for sqlite_config.go depends on a prompt convention the worker must read first; if the convention is ambiguous, the worker picks the flat generic form, records the choice, and proceeds.
