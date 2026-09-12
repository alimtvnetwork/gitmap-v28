package cmddb

// RunDB executes the main database management CLI command.
func RunDB(args []string) error {
	return runDB(args)
}

// RunStartFresh executes the fresh database reconstruction CLI command.
func RunStartFresh(args []string) error {
	return runStartFresh(args)
}

// FormatBytes formats byte counts into human-readable strings.
func FormatBytes(bytes int64) string {
	return formatBytes(bytes)
}

// ConfirmOrSkip checks interactive confirmation or CLI bypass flags.
func ConfirmOrSkip(msg string, args []string) bool {
	return confirmOrSkip(msg, args)
}

// IsInteractiveStdin checks whether stdin is an interactive character device.
func IsInteractiveStdin() bool {
	return isInteractiveStdin()
}

// HasConfirmFlag checks if args contain confirmation/force flags.
func HasConfirmFlag(args []string) bool {
	return hasConfirmFlag(args)
}

// ParseConfirmFlag parses the confirm flag.
func ParseConfirmFlag(cmdName string, args []string) bool {
	return parseConfirmFlag(cmdName, args)
}

// TruncateStr truncates a string with ellipsis if exceeding maxLen.
func TruncateStr(s string, maxLen int) string {
	return truncateStr(s, maxLen)
}
