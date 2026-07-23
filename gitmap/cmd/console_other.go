//go:build !windows

package cmd

// initConsole is a no-op on non-Windows platforms because Linux and macOS
// terminals already speak UTF-8 by default.
func initConsole() {}
