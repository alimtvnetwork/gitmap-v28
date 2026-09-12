// Package cmdchrome — helpers.go provides shared utilities for Chrome commands.
package cmdchrome

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/helptext"
)

func confirmYesNo(prompt string) bool {
	fmt.Printf("%s (y/N): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

func printJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")

	return enc.Encode(v)
}

func checkHelp(command string, args []string) {
	if CheckHelpFn != nil {
		CheckHelpFn(command, args)
		return
	}

	for _, a := range args {
		if a == "--help" || a == "-h" || a == "help" {
			helptext.Print(command)
			cliexit.Exit(0)
		}
	}
}
