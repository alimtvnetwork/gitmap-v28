// Package cmdvmware — vmware_crontab.go: crontab persistence management for VMware shared folders.
package cmdvmware

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

var (
	CrontabCommandFunc = exec.Command
	ReadCrontabFunc    = readCurrentCrontab
	WriteCrontabFunc   = writeCrontab
)

func readCurrentCrontab() string {
	cmd := CrontabCommandFunc("crontab", "-l")
	outBytes, err := cmd.Output()
	if err != nil {
		return ""
	}

	return cleanCrontabOutput(string(outBytes))
}

func cleanCrontabOutput(raw string) string {
	trimmed := strings.TrimSpace(raw)
	if strings.Contains(trimmed, "no crontab for") {
		return ""
	}

	return trimmed
}

func isCrontabPersisted() bool {
	current := ReadCrontabFunc()

	return strings.Contains(current, "vmhgfs-fuse") && strings.Contains(current, defaultMountPoint)
}

func buildUpdatedCrontab(current, rebootLine string) string {
	if current == "" {
		return rebootLine + "\n"
	}

	return current + "\n" + rebootLine + "\n"
}

func writeCrontab(content string) error {
	cmd := CrontabCommandFunc("crontab", "-")
	cmd.Stdin = strings.NewReader(content)
	out, err := cmd.CombinedOutput()
	if err != nil {
		errMsg := fmt.Sprintf("failed updating crontab: %s (%v)", string(out), err)

		return apperror.NewWithDetails("cmd.vmware.crontab", "E4004", errMsg, "cmd.vmware", apperror.ErrorTypeExecution, apperror.SeverityError, nil)
	}

	return nil
}

// EnsureCrontabPersistence guarantees VMware shared folder mounts persist across reboots.
func EnsureCrontabPersistence() error {
	if isCrontabPersisted() {
		return nil
	}

	updated := buildUpdatedCrontab(ReadCrontabFunc(), crontabRebootLine)

	return WriteCrontabFunc(updated)
}
