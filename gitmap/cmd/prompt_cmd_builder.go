// Package cmd — prompt_cmd_builder.go generates OS-aware command strings.
package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func BuildUnixPromptInstallCmd() string {
	return fmt.Sprintf("curl -sL %s | bash", constants.PromptArchitectBashURL)
}

func BuildWindowsPromptInstallCmd() string {
	return fmt.Sprintf(`Invoke-Expression "& { $(Invoke-RestMethod %s) }"`, constants.PromptArchitectPowerShellURL)
}
