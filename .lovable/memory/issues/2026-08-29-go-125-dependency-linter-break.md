# Issue: Go 1.25 Dependency Version Bump Causing CI Linter Failure

- date: 2026-08-29
- status: resolved
- scope: build/ci

## 1. Why it happened
A sub-package upgrade or dependency addition pulled in bleeding-edge versions of `golang.org/x/*` packages (`v0.55.0`, `v0.47.0`, `v0.41.0`) which have minimum version requirements set to `go 1.25.0`, causing `gitmap/go.mod` to declare `go 1.25.0`.

## 2. How it happened
When CI runs `bash ../.github/scripts/full-suite-lint.sh`, it invokes `golangci-lint v1.64.8` (which is compiled against Go 1.24). The linter's `typecheck` analyzer refuses to parse packages declaring Go 1.25 or depending on modules requiring Go 1.25, emitting `package requires newer Go version go1.25 (application built with go1.24)` on every file.

## 3. Root Cause
- `gitmap/go.mod` line 3 declared `go 1.25.0` with `golang.org/x/crypto v0.55.0`, `golang.org/x/sys v0.47.0`, and `golang.org/x/text v0.41.0`.
- Incompatibility between Go 1.25 dependencies and the project-pinned `golangci-lint v1.64.8` Go 1.24 toolchain.

## 4. Code Fix
- Reverted `gitmap/go.mod` to `go 1.24.13` and downgraded `golang.org/x/*` dependencies to stable Go 1.24 releases (`crypto v0.36.0`, `sys v0.41.0`, `term v0.40.0`, `text v0.23.0`, `tools v0.31.0`).
- Regenerated `go.sum` via `go build -mod=mod ./...`.
- Renamed test fixture paths in `vscodepm/update_remove_test.go` from `other-app` to `second-app` to eliminate `misspell` scanner false positives on Windows backslashes.
- Verified all packages pass `go test ./...` and `.lovable/ai-fix-scripts/02-cicd-local-runner.py` with 0 failures.
