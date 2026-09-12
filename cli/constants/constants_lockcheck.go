package constants

// Lock-check messages.
const (
	MsgLockCheckScanning    = "Checking for processes locking %s...\n"
	MsgLockCheckFound       = "The following processes are using this folder:\n%s\n"
	MsgLockCheckKillPrompt  = "Terminate these processes to allow deletion? [y/N] "
	MsgLockCheckKilling     = "Terminating %s (PID %d)...\n"
	MsgLockCheckKilled      = "✓ Terminated %s\n"
	MsgLockCheckRetrying    = "Retrying folder removal...\n"
	WarnLockCheckKillFailed = "Warning: could not terminate %s (PID %d): %v\n"
	WarnLockCheckScanFailed = "Warning: could not scan for locking processes: %v\n"
	MsgLockCheckNoneFound   = "No locking processes detected.\n"
)
