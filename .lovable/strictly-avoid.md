Strictly avoid using flash model for high concurrency batches that exceed rate limits.
- Strictly avoid using `*_linux.go` file naming for generic non-Windows (`//go:build !windows`) code. Use `*_posix.go` or `*_unix.go` so macOS (`darwin`) and other POSIX systems are not excluded by Go build tools.
- NEVER add new top-level \Cmd*\ constants to \constants_cli.go\ without updating the \	opLevelCmds()\ map in \cmd_constants_test.go\ (or marking them with \// gitmap:cmd skip\), otherwise the AST parity test will fail the CI.
- ALWAYS use fenced code blocks (\\\) in \helptext/*.md\ files. Do not use 4-space indentations, or the examples golden test will fail.
- NEVER commit changes to \.github/workflows\ without verifying valid YAML syntax. Be wary of Cloudflare email obfuscation replacing \@version\ with \[email protected]\.
- NEVER bump go.mod to a new minor Go release (like 1.25) without also confirming that pinned linters (like golangci-lint) support that Go release.

- NEVER invoke a local GitHub Action composite (uses: ./.github/actions/...) without running actions/checkout FIRST in the job.

- NEVER use PowerShell literal newlines (\`n) as text replacement values in non-PowerShell files.
