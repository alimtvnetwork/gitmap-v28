package cmdspace

import "fmt"

// printSendEmailHelp prints the legacy send-email usage text (the const
// lives in space.go next to the other space usage strings).
func printSendEmailHelp() {
	fmt.Print(spaceSendEmailUsage)
}
