package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

// appendSSHHelp appends SSH documentation to the provided writer.
func appendSSHHelp(cmd *cobra.Command, args []string, buf io.Writer) error {
	_, err := fmt.Fprint(buf, `
## SSH Commands

Examples:
- gitmap ssh m1
- gitmap ssh 192.168.1.5

### Alias Resolution
When you use a short alias like 'm1' or an IP, gitmap checks the local SQLite 'ssh_hosts' table. 
If an alias is registered (e.g. via 'gitmap ssh <ip> as <alias>'), it will resolve to the stored username and IP, automatically connecting without requiring you to type the full credentials each time.
`)
	if err != nil {
		return apperror.NewSimple("appendSSHHelp", "E_INTERNAL_ERROR")
	}

	return nil
}
