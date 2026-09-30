# Subtask 04: Minor Version Bump & Release Ceremony

## Scope
1. Version bump to `v6.426.0` across:
   - `package.json`
   - `version.json`
   - `cli/constants/constants.go`
2. Update `changelog.md` with release notes for `v6.426.0`.
3. Verify all repository linters and formatting pass:
   - `python linter-scripts/check-nested-ifs.py` (0 violations)
   - `python linter-scripts/check-enum-and-boolean.py` (0 violations)
   - `gofmt -w` on all modified files
4. Atomic commit, create tag `v6.426.0`, push to `origin main`, `origin release/v6.426.0`, and `v6.426.0`.
5. Monitor GitHub Actions release and CI workflows until 100% green.
6. Run `gitmap update` on host to upgrade local binary and verify live commands.

## Status: COMPLETED
- `v6.426.0` released and verified.
- `v6.427.0` release underway for clone only-missing subcommand routing, live progress, and pending task deduplication.
