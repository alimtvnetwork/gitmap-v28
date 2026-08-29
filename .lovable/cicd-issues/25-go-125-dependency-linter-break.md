# 25-go-125-dependency-linter-break

## Error Summary
The CI/CD pipeline failed during `bash ../.github/scripts/full-suite-lint.sh` with repeated typecheck errors:
`Error: package requires newer Go version go1.25 (application built with go1.24) (typecheck)` across `apperror`, `result`, `ghtoken`, `verbose`, `downloaderconfig`, `cmd/*`, `detector`, `worker`, `indexer`, and external dependencies `golang.org/x/text`, `golang.org/x/crypto`.

## Root Cause Analysis
- `gitmap/go.mod` was pinned to `go 1.25.0` with bleeding-edge `golang.org/x/*` dependencies (`golang.org/x/crypto v0.55.0`, `golang.org/x/sys v0.47.0`, `golang.org/x/term v0.45.0`, `golang.org/x/text v0.41.0`, `golang.org/x/tools v0.49.0`).
- The GitHub Actions workflow (`ci.yml`) installs pinned `golangci-lint v1.64.8`, which was compiled with Go 1.24.
- When `golangci-lint v1.64.8` invoked its internal `go/packages` typechecker, it refused to parse packages declaring `go 1.25.0` or importing dependencies requiring `go 1.25.0`, triggering build-wide typecheck failure.

## Solution Applied
1. Downgraded `gitmap/go.mod` to `go 1.24.13` and pinned `golang.org/x/*` dependencies to their stable Go 1.24-compatible releases (`golang.org/x/crypto v0.36.0`, `golang.org/x/sys v0.41.0`, `golang.org/x/term v0.40.0`, `golang.org/x/text v0.23.0`, `golang.org/x/tools v0.31.0`, `golang.org/x/mod v0.24.0`).
2. Ran `go build -mod=mod ./...` to refresh `go.sum` hashes without allowing local Go toolchains to auto-upgrade to 1.25.
3. Fixed `misspell` false positive on `D:\work\other-app` in `vscodepm/update_remove_test.go` by renaming test fixtures to `second-app`.
4. Verified that all tests (`go test ./... -count=1`) and `.lovable/ai-fix-scripts/02-cicd-local-runner.py` pass 100% locally with 0 lint errors and 0 vet warnings.

## What NOT to Repeat
- NEVER bump `golang.org/x/*` dependencies to bleeding-edge versions that require Go 1.25 when the CI linter is pinned to Go 1.24.
- ALWAYS test with `.lovable/ai-fix-scripts/02-cicd-local-runner.py` before pushing.
