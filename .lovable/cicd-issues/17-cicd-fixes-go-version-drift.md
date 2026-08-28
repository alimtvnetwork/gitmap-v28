# 17. CI/CD Fixes: Go Version Drift & Sync

## Error Summary
Three CI/CD checks failed:
1. golangci-lint: `the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.0)`
2. Generates drift: `Generated files are out of sync with constants.`
3. CHANGELOG version-sync: `changelog.md not found`
4. Test failure: `TestResolveProfileTree_Aliases` failed because `vscode` resolved to `vscode-settings` instead of `ubuntu+vscode`.

## Root Cause Analysis
1. `gitmap/go.mod` was updated to `go 1.25.0`, but the CI/CD pipeline pins `golangci-lint` to `v1.64.8`, which was compiled against Go 1.24. This caused a runtime failure when analyzing Go 1.25.0 code.
2. Changes to command files or ASTs weren't followed by a local `go generate ./...`, causing `completion/allcommands_generated.go` to desync.
3. The bash script `.github/scripts/check-changelog-version-sync.sh` explicitly searched for `changelog.md` in all caps, while the repository file is named `changelog.md` (lowercase). Additionally, the version string was `6.103.0` in `gitmap/constants/constants.go` instead of syncing with `version.json`.
4. A recent feature added `vscode-settings` installer profile and mapped the alias `vscode` to it, conflicting with the older test assertion that expected `vscode` to resolve to `ubuntu+vscode`.

## Solution Applied
1. Downgraded `gitmap/go.mod` back to `go 1.24.0` (and ran `go mod tidy`) to safely pass the pinned `golangci-lint` check without violating the `mem://core` pin policy.
2. Ran `go generate ./...` in the `gitmap` directory to regenerate `allcommands_generated.go`.
3. Updated `.github/scripts/check-changelog-version-sync.sh` to target `changelog.md` instead of `changelog.md`. 
4. Synced `gitmap/constants/constants.go` to match `version.json`'s version (`6.27.0`) and ensured the changelog heading uses `v` prefix (`## [v6.27.0]`).
5. Updated `gitmap/cmd/install_profile_tree_test.go` to expect `vscode-settings` when the `vscode` slug is used.

## What NOT to Repeat
- NEVER bump `go.mod` to a new minor Go release (like 1.25) without also confirming that pinned linters (like golangci-lint) support that Go release.
- ALWAYS run `go generate ./...` locally before pushing changes that modify CLI commands or structures.
- ALWAYS match the exact filename casing (`changelog.md`) in CI scripts, since Linux runners are case-sensitive.
- ALWAYS ensure `gitmap/constants/constants.go` is synchronized with `version.json`.
- ALWAYS verify tests pass (`go test ./...`) after adding aliases or mapping rules to core engines like `ResolveProfileTree`.
