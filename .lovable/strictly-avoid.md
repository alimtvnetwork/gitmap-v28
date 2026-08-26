Strictly avoid using flash model for high concurrency batches that exceed rate limits.
- Strictly avoid using `*_linux.go` file naming for generic non-Windows (`//go:build !windows`) code. Use `*_posix.go` or `*_unix.go` so macOS (`darwin`) and other POSIX systems are not excluded by Go build tools.
