# 10 — Go filename GOOS-suffix trap (`lang_js.go` excluded on linux)

## Original task
Step 1 of the 12-step commit-in implementation: audit + repair Phases
1-7 before wiring orchestration.

## Ambiguity
Phase 7 created per-language detector files named after the language
(`lang_go.go`, `lang_js.go`, `lang_ts.go`, …) following the spec's §6.5
naming hint. On `go build` this surfaced an "undefined: matchJsFunc"
error from `lang_ts.go` even though `matchJsFunc` was defined in
`lang_js.go` in the same package. Root cause: Go's build system treats
filenames matching `*_<GOOS>.go` as build-constrained; `js` IS a valid
GOOS (js/wasm), so `lang_js.go` was silently excluded on linux/darwin/
windows builds. The spec did not anticipate this.

## Options considered
1. Rename only `lang_js.go` → `lang_javascript.go` (matches enum
   token `CommitInLanguageJavaScript`).
2. Rename every per-language file to its full enum-token form
   (`lang_go.go` → `lang_go.go` is fine, but rename `lang_ts.go` →
   `lang_typescript.go`, `lang_php.go` → `lang_php.go` fine, etc.) for
   uniformity.
3. Move all detectors into one file (`detectors.go`) sidestepping the
   issue but violating the "one detector per file" spec hint.

## Recommendation
Option 1 — minimal blast radius. `js` is the only short token that
collides with a real GOOS; `ts`, `php`, `rs`, `py`, `java`, `cs` are
safe (none are GOOS values). Renaming all files would churn diffs and
memory for no compile-correctness benefit.

## Decision taken
Option 1: renamed `gitmap/cmd/commitin/funcintel/lang_js.go` →
`lang_javascript.go`. Verified `go build ./...` and
`go test ./cmd/commitin/...` are green.

## Future-proofing
Added rule to internal memory: when naming per-language Go source
files, AVOID GOOS-suffix collisions (`*_js.go`, `*_linux.go`,
`*_windows.go`, `*_darwin.go`, `*_android.go`, `*_ios.go`, `*_plan9.go`,
`*_aix.go`, `*_freebsd.go`, `*_netbsd.go`, `*_openbsd.go`, `*_solaris.go`,
`*_dragonfly.go`, `*_illumos.go`, `*_zos.go`, `*_wasip1.go`). Same
applies to GOARCH suffixes (`*_amd64.go`, `*_arm64.go`, …).