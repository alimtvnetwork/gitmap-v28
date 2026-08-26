Strictly avoid using flash model for high concurrency batches that exceed rate limits.
- Strictly avoid using `*_linux.go` file naming for generic non-Windows (`//go:build !windows`) code. Use `*_posix.go` or `*_unix.go` so macOS (`darwin`) and other POSIX systems are not excluded by Go build tools.
- NEVER add new top-level \Cmd*\ constants to \constants_cli.go\ without updating the \	opLevelCmds()\ map in \cmd_constants_test.go\ (or marking them with \// gitmap:cmd skip\), otherwise the AST parity test will fail the CI.
