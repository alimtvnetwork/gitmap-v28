package cmd

import (
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
)

func safeRemoveWithRetry(path string) error {
	var err error
	for attempt := 0; attempt < 3; attempt++ {
		err = fsutil.SafeRemoveAll(path)
		if err == nil || !isProcessLockError(err) {
			return err
		}

		time.Sleep(150 * time.Millisecond)
	}

	return err
}

func isProcessLockError(err error) bool {
	if err == nil {
		return false
	}

	msg := strings.ToLower(err.Error())

	return strings.Contains(msg, "being used by another process") ||
		strings.Contains(msg, "access is denied") ||
		strings.Contains(msg, "unlinkat") ||
		strings.Contains(msg, "permission denied")
}
