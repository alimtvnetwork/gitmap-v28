// Package cmdssh — ssh_oneliner.go generates a single-line command for SSH nodes import with clipboard copy.
package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/atotto/clipboard"
)

// RunSSHExportOnelinerCLI generates a single-line command containing all SSH nodes encoded in Base64 and copies it to clipboard.
func RunSSHExportOnelinerCLI(args []string) error {
	envelope, err := BuildCompactSSHNodesExportEnvelope()
	if err != nil {
		return err
	}
	compactBytes, err := json.Marshal(envelope)
	if err != nil {
		return err
	}
	b64 := base64.StdEncoding.EncodeToString(compactBytes)
	cmd := fmt.Sprintf("gitmap ssh nodes import-json --base64 \"%s\"", b64)

	fmt.Println(cmd)
	copyOnelinerToClipboard(cmd)

	return nil
}

func copyOnelinerToClipboard(cmd string) {
	err := clipboard.WriteAll(cmd)
	if err != nil {
		fmt.Printf("\n  %s[clip] Note: clipboard unavailable in session; copy the command above manually.%s\n", constants.ColorDim, constants.ColorReset)
		return
	}
	fmt.Printf("\n  %s📋 Copied single-line import command to clipboard!%s\n", constants.ColorGreen, constants.ColorReset)
}
