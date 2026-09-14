package cmdssh

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/spf13/cobra"
)

const leafHelpAdd = `Enroll remote SSH machine into GitMap registry with an alias and optional key auth.

Usage:
  gitmap ssh-join <user@ip|ip> [alias] [flags]
  gitmap ssh-join add <user@ip|ip> [alias] [flags]
  gitmap ssh join <user@ip|ip> [alias] [flags]
  gitmap sj <user@ip|ip> [alias] [flags]

Aliases:
  add, enroll, join, new

Flags:
  -u, --user string   Remote SSH username
  -n, --name string   Memorable alias name
  -p, --port int      Target SSH port (default 22)
      --auth          Push local public key to authorized_keys
  -f, --force         Overwrite existing alias mapping

Examples:
  gitmap ssh-join user@192.168.1.14
  gitmap ssh-join root@192.168.1.14 prod-server
  gitmap ssh-join alim@192.168.1.14 devbox --auth
  gitmap ssh-join add ubuntu@192.168.1.50 cloud-vm
  gitmap ssh join admin@10.0.0.12:2222 gateway
  gitmap sj 192.168.1.50 devbox`

const leafHelpScan = `Scan local subnet or CIDR network for active SSH machines on port 22.

Usage:
  gitmap ssh-join scan [subnet] [flags]
  gitmap ssh join scan [subnet] [flags]
  gitmap sj scan [subnet] [flags]

Aliases:
  find, discover, probe

Flags:
  -p, --port int          Target SSH port to probe (default 22)
  -t, --timeout duration  Per-host probe timeout (default 800ms)
  -w, --workers int       Concurrent worker pool size (default 40)

Examples:
  gitmap ssh-join scan
  gitmap ssh-join scan 192.168.1.0/24
  gitmap sj find 10.0.0.0/24 --port 2222
  gitmap sj discover --timeout 1s --workers 50`

const leafHelpStatus = `Check machine connectivity, latency, and health of registered SSH machines.

Usage:
  gitmap ssh-join status [alias|ip] [flags]
  gitmap ssh join status [alias|ip] [flags]
  gitmap sj status [alias|ip] [flags]

Aliases:
  ping, health, check

Flags:
  -p, --port int          Target SSH port to probe (default 22)
  -t, --timeout duration  Per-host probe timeout (default 1500ms)

Examples:
  gitmap ssh-join status
  gitmap ssh-join status devbox
  gitmap ssh-join ping 192.168.1.50
  gitmap sj health --port 2222 --timeout 2s`

const leafHelpLs = `List all enrolled SSH machines, aliases, users, and creation timestamps.

Usage:
  gitmap ssh-join ls [flags]
  gitmap ssh join ls [flags]
  gitmap sj ls [flags]

Aliases:
  list

Examples:
  gitmap ssh-join ls
  gitmap sj ls
  gitmap ssh join ls`

const leafHelpRm = `Remove an enrolled SSH machine from registry by alias or IP.

Usage:
  gitmap ssh-join rm <alias|ip>
  gitmap ssh join rm <alias|ip>
  gitmap sj rm <alias|ip>

Aliases:
  remove, delete

Examples:
  gitmap ssh-join rm devbox
  gitmap ssh-join rm 192.168.1.50
  gitmap ssh join rm staging`

const leafHelpAuth = `Push local SSH public key to remote machine authorized_keys.

Usage:
  gitmap ssh-join add-auth <user@ip|alias|ip>
  gitmap ssh join add-auth <user@ip|alias|ip>
  gitmap sj add-auth <user@ip|alias|ip>

Aliases:
  auth

Examples:
  gitmap ssh-join add-auth user@192.168.1.14
  gitmap ssh-join add-auth root@192.168.1.14
  gitmap ssh-join add-auth devbox
  gitmap sj add-auth 192.168.1.50`

const leafHelpHistory = `Display audit trail of joined machines with timestamps and users.

Usage:
  gitmap ssh-join history [filter]
  gitmap ssh join history [filter]
  gitmap sj history [filter]

Aliases:
  hist

Examples:
  gitmap ssh-join history
  gitmap ssh-join history 192.168.1.50
  gitmap sj hist ubuntu`

var leafHelpCatalog = map[string]string{
	"add":      leafHelpAdd,
	"enroll":   leafHelpAdd,
	"join":     leafHelpAdd,
	"scan":     leafHelpScan,
	"find":     leafHelpScan,
	"discover": leafHelpScan,
	"probe":    leafHelpScan,
	"status":   leafHelpStatus,
	"ping":     leafHelpStatus,
	"health":   leafHelpStatus,
	"check":    leafHelpStatus,
	"ls":       leafHelpLs,
	"list":     leafHelpLs,
	"rm":       leafHelpRm,
	"remove":   leafHelpRm,
	"delete":   leafHelpRm,
	"add-auth": leafHelpAuth,
	"auth":     leafHelpAuth,
	"history":  leafHelpHistory,
	"hist":     leafHelpHistory,
}

func lookupLeafHelp(cmdName string) (string, bool) {
	norm := strings.ToLower(strings.TrimSpace(cmdName))
	text, hasHelp := leafHelpCatalog[norm]
	return text, hasHelp
}

// GetLeafHelp returns dedicated formatted help text for a leaf subcommand.
func GetLeafHelp(cmdName string) (string, bool) {
	return lookupLeafHelp(cmdName)
}

// PrintLeafHelp renders dedicated help text to the provided writer.
func PrintLeafHelp(out io.Writer, cmdName string) error {
	text, hasHelp := lookupLeafHelp(cmdName)
	if !hasHelp {
		return apperror.NewValidationError(fmt.Sprintf("unknown leaf command: %s", cmdName))
	}
	fmt.Fprintln(out, text)
	return nil
}

// AttachLeafHelp configures custom help rendering on a Cobra command.
func AttachLeafHelp(cmd *cobra.Command, cmdName string) {
	if cmd == nil {
		return
	}
	text, hasHelp := lookupLeafHelp(cmdName)
	if !hasHelp {
		return
	}
	cmd.SetHelpFunc(func(c *cobra.Command, args []string) {
		fmt.Fprintln(c.OutOrStdout(), text)
	})
}

func attachAllLeafCommandsHelp() {
	AttachLeafHelp(SJAddCmd, "add")
	AttachLeafHelp(SJScanCmd, "scan")
	AttachLeafHelp(SJStatusCmd, "status")
	AttachLeafHelp(SJLsCmd, "ls")
	AttachLeafHelp(SJRmCmd, "rm")
	AttachLeafHelp(SJAddAuthCmd, "add-auth")
	AttachLeafHelp(SJHistCmd, "history")
}

func init() {
	attachAllLeafCommandsHelp()
}
