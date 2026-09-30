# Subtask 04: Minor Version Bump & Release Ceremony

## Scope
1. Version bump to `v6.426.0` across:
   - `package.json`
   - `version.json`
   - `cli/constants/constants.go`
2. Update `changelog.md` with release notes for `v6.426.0`.
3. Verify all 5 repository linters and formatting pass locally:
   - `python .github/scripts/go-format-check.py --check-only`
   - `python .github/scripts/file-size-check.py`
   - `python .github/scripts/check-constants-collisions.py`
   - `python .github/scripts/check-changelog-version-sync.py`
4. Atomic commit, create tag `v6.426.0`, push to `origin main`, `origin release/v6.426.0`, and `v6.426.0`.
5. Monitor GitHub Actions release and CI workflows until 100% green.
6. Run `gitmap update` on host to upgrade local binary and verify live commands.
