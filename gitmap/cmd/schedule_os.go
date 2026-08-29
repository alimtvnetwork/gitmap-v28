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
	switch runtime.GOOS {
	case "windows":
		executeOSCommand("restart", exec.Command("shutdown", "/r", "/t", "0"))
	case "darwin":
		executeOSCommand("restart", exec.Command("sudo", "shutdown", "-r", "now"))
	default:
		executeOSCommand("restart", exec.Command("sudo", "reboot"))
	}
	return nil
}

// runScheduleShutdown runs native OS shutdown.
func runScheduleShutdown() error {
	switch runtime.GOOS {
	case "windows":
		executeOSCommand("shutdown", exec.Command("shutdown", "/s", "/t", "0"))
	case "darwin":
		executeOSCommand("shutdown", exec.Command("sudo", "shutdown", "-h", "now"))
	default:
		executeOSCommand("shutdown", exec.Command("sudo", "poweroff"))
	}
	return nil
}
