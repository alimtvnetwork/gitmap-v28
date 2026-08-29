package cmd

import (
	"fmt"
	"runtime"

	"github.com/atotto/clipboard"
	"github.com/pterm/pterm"
)

//nolint:unused
func runInstallScript() error {
	var script string

	if runtime.GOOS == "windows" {
		script = "irm https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.ps1 | iex"
	} else {
		script = "curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/gitmap-v28/main/install.sh | sh"
	}

	err := clipboard.WriteAll(script)
	if err != nil {
		pterm.Error.Printf("Failed to copy to clipboard: %v\n", err)
		fmt.Println("Here is the script instead:")
		fmt.Println(script)
		return nil
	}

	pterm.Success.Println("Install script copied to !clipboard")
	fmt.Println()
	fmt.Println(script)
	return nil
}
