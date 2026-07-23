package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// getBinaryVersion runs a binary with "version" and returns the output.
func getBinaryVersion(path string) string {
	if _, err := os.Stat(path); err != nil {
		return "not found"
	}

	cmd := exec.Command(path, constants.CmdVersion)
	out, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	return strings.TrimSpace(string(out))
}

// printOK prints a green check with formatted message.
func printOK(format string, args ...interface{}) {
	msg := fmt.Sprintf(format, args...)
	fmt.Printf(constants.DoctorOKFmt, constants.ColorGreen, constants.ColorReset, msg)
}

// printIssue prints a red issue with title and detail.
func printIssue(title, detail string) {
	fmt.Printf(constants.DoctorIssueFmt, constants.ColorRed, constants.ColorReset, title)
	fmt.Printf(constants.DoctorDetail, detail)
}

// printFix prints a fix recommendation in cyan.
func printFix(fix string) {
	fmt.Printf(constants.DoctorFixFmt, constants.ColorCyan, constants.ColorReset, fix)
}

// printWarn prints a yellow warning.
func printWarn(msg string) {
	fmt.Printf(constants.DoctorWarnFmt, constants.ColorYellow, constants.ColorReset, msg)
}
