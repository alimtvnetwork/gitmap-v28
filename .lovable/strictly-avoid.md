Strictly avoid using flash model for high concurrency batches that exceed rate limits.
- Strictly avoid using `*_linux.go` file naming for generic non-Windows (`//go:build !windows`) code. Use `*_posix.go` or `*_unix.go` so macOS (`darwin`) and other POSIX systems are not excluded by Go build tools.
- NEVER add new top-level \Cmd*\ constants to \constants_cli.go\ without updating the \	opLevelCmds()\ map in \cmd_constants_test.go\ (or marking them with \// gitmap:cmd skip\), otherwise the AST parity test will fail the CI.
- ALWAYS use fenced code blocks (\\\) in \helptext/*.md\ files. Do not use 4-space indentations, or the examples golden test will fail.
- NEVER commit changes to \.github/workflows\ without verifying valid YAML syntax. Be wary of Cloudflare email obfuscation replacing \@version\ with \[email protected]\.
- NEVER bump go.mod to a new minor Go release (like 1.25) without also confirming that pinned linters (like golangci-lint) support that Go release.

- NEVER invoke a local GitHub Action composite (uses: ./.github/actions/...) without running actions/checkout FIRST in the job.

- NEVER use PowerShell literal newlines (\`n) as text replacement values in non-PowerShell files.

- NEVER open files in Python scripts without explicitly specifying encoding="utf-8".

- NEVER use global or greedy regex replacements for version bumping in documentation; always target a specific anchor.

- NEVER run text-replacements on .go files without subsequently running gofmt -w ..
- NEVER change a return type inside a deeply nested utility without cascading it all the way to the top-level invoker.
- NEVER rename API JSON boolean keys without verifying downstream Python/Bash scripts that consume them.

- NEVER add or modify CLI constants, shell commands, or help text without subsequently running go generate ./... in the gitmap directory to prevent generated file drift.
- NEVER blindly commit a dirty working tree without verifying git status/diff, to avoid reverting previous linter fixes.
- NEVER use British English spelling (e.g., `behavior`, `recognize`) in the codebase; use US English to pass the misspell linter.
- NEVER write bare `fmt.Fprintln(os.Stderr, err)` in `gitmap/cmd/`; always use `cliexit.Reportf` or `cliexit.Fail`.
- NEVER access dictionary keys with direct index `res["is_failure"]` when consuming `query_wrapper` output; always use `res.get("is_fail")`.
- ALWAYS annotate historical mentions of legacy version names with `<!-- gitmap-legacy-ref-allow -->` to prevent policy check scan triggers.
- NEVER run `git diff --exit-code` in CI code generation drift checks without scoping to the module directory (`git diff --exit-code .`), to prevent Git LFS raw file conversions on root folders from causing false drift detection.
- NEVER bump versions across `version.json` or `readme.md` without synchronizing `var Version` in `gitmap/constants/constants.go` and adding a matching `## [vX.Y.Z]` heading in `changelog.md`.
- NEVER pass unsanitized variables to `jq --argjson` in bash scripts without validating numeric integer format (`[[ "$LINE" =~ ^[0-9]+$ ]]`).


### No Shell Wrappers for Python Scripts
Cross-platform Python scripts should be invoked natively via \python script.py\ in CI/CD workflows. **Strictly avoid** creating \.sh\ or \.ps1\ wrappers that simply forward arguments to a Python script, as these wrappers can fail unexpectedly in cross-platform GitHub Actions environments.

