package cmd

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// executeOSCommand executes an OS command and prints error.
func executeOSCommand(op string, cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		appErr := apperror.WrapSimple(err, op)
		fmt.Println("Error:", appErr)
	} else {
		fmt.Println(op, "executed successfully")
	}
}

// runScheduleRestart runs native OS restart.
func runScheduleRestart() error {
	isWindows := runtime.GOOS == "windows"
	isMac := runtime.GOOS == "darwin"
	if isWindows {
		executeOSCommand("restart", exec.Command("shutdown", "/r", "/t", "0"))
	} else if isMac {
		executeOSCommand("restart", exec.Command("sudo", "shutdown", "-r", "now"))
	} else {
		executeOSCommand("restart", exec.Command("sudo", "reboot"))
	}
	return nil
}

// runScheduleShutdown runs native OS shutdown.
func runScheduleShutdown() error {
	isWindows := runtime.GOOS == "windows"
	isMac := runtime.GOOS == "darwin"
	if isWindows {
		executeOSCommand("shutdown", exec.Command("shutdown", "/s", "/t", "0"))
	} else if isMac {
		executeOSCommand("shutdown", exec.Command("sudo", "shutdown", "-h", "now"))
	} else {
		executeOSCommand("shutdown", exec.Command("sudo", "poweroff"))
	}
	return nil
}
