package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

// SSHPassCmd represents the 'pass' command under ssh.
var SSHPassCmd = &cobra.Command{
	Use:     "pass [show|ls] <alias|ip> [flags]",
	Aliases: []string{"password"},
	Short:   "Inspect, review, or verify saved SSH encrypted node passwords",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHPassCLI(args)
	},
}

// RunSSHPassCLI routes pass subcommands (show, ls, or direct target).
func RunSSHPassCLI(args []string) error {
	hasHelp := len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help")
	if hasHelp || len(args) == 0 {
		return printPassHelp()
	}

	sub := strings.ToLower(args[0])
	if sub == "ls" || sub == "list" {
		return runPassList()
	}

	return dispatchPassTarget(sub, args)
}

func dispatchPassTarget(sub string, args []string) error {
	target := args[0]
	if sub == "show" || sub == "view" || sub == "get" {
		if len(args) < 2 {
			return apperror.NewValidationError("missing node alias or IP. Usage: gitmap ssh pass show <alias|ip>")
		}
		target = args[1]
	}

	return runPassShow(target)
}

func printPassHelp() error {
	fmt.Println("Usage: gitmap ssh pass <show|ls> [alias|ip] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  show <alias|ip>  Reveal the encrypted saved password for a node")
	fmt.Println("  ls, list         List all nodes with encrypted password status")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh pass show u2")
	fmt.Println("  gitmap ssh pass show 192.168.1.5")
	fmt.Println("  gitmap ssh pass ls")

	return nil
}

func runPassShow(target string) error {
	encPass, hostDesc, err := resolveNodeEncryptedPassword(target)
	if err != nil {
		return err
	}

	hasNoPass := encPass == ""
	if hasNoPass {
		fmt.Printf("No saved password for node: %s\n", target)

		return nil
	}

	plain, decErr := DecryptSSHPassword(encPass)
	if decErr != nil {
		return apperror.WrapSimple(decErr, "decrypt password")
	}

	renderPasswordDetails(target, hostDesc, plain)

	return nil
}

func renderPasswordDetails(target, hostDesc, plain string) {
	fmt.Printf("\n  %s● Saved SSH Password%s for '%s' (%s)\n",
		constants.ColorGreen, constants.ColorReset, target, hostDesc)
	fmt.Printf("    Password: %s%s%s\n", constants.ColorYellow, plain, constants.ColorReset)
	fmt.Println("    Storage:  RSA-OAEP / AES Encrypted Vault")
	fmt.Println("    Security: Password is encrypted at rest and only decrypted in memory.")
	fmt.Println()
}

func resolveNodeEncryptedPassword(target string) (string, string, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return "", "", apperror.WrapSimple(err, "open db")
	}
	defer dbConn.Close()

	ctx := context.Background()
	host, hostErr := store.GetSSHHostByAlias(ctx, target, dbConn.SQL())
	if hostErr == nil && host != nil {
		return host.EncryptedPassword, fmt.Sprintf("%s@%s", host.Username, host.IP), nil
	}

	hostByIP, ipErr := store.GetSSHHostByIP(ctx, target, dbConn.SQL())
	if ipErr == nil && hostByIP != nil {
		return hostByIP.EncryptedPassword, fmt.Sprintf("%s@%s", hostByIP.Username, hostByIP.IP), nil
	}

	return resolveFallbackConnPassword(ctx, dbConn.SQL(), target)
}

func resolveFallbackConnPassword(ctx context.Context, sqlDB *sql.DB, target string) (string, string, error) {
	connRes := db.GetSSHConnectionByAlias(ctx, sqlDB, target)
	if connRes.IsSuccess() {
		return connRes.Data.EncryptedPassword, fmt.Sprintf("%s@%s", connRes.Data.Username, connRes.Data.IPAddress), nil
	}

	connIPRes := db.GetSSHConnectionByIP(ctx, sqlDB, target)
	if connIPRes.IsSuccess() {
		return connIPRes.Data.EncryptedPassword, fmt.Sprintf("%s@%s", connIPRes.Data.Username, connIPRes.Data.IPAddress), nil
	}

	return "", "", apperror.NewNotFoundError("node not found: " + target)
}

func runPassList() error {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "open db")
	}
	defer dbConn.Close()

	hosts, err := store.ListSSHHosts(context.Background(), dbConn.SQL())
	if err != nil {
		return apperror.WrapSimple(err, "list hosts")
	}

	renderPassListTable(hosts)

	return nil
}

func renderPassListTable(hosts []store.SSHHost) {
	fmt.Println("\n● Registered Node Password Vault Status:")
	fmt.Printf("  %-16s %-16s %-12s %s\n", "ALIAS", "IP", "USER", "PASSWORD STATUS")
	fmt.Println("  --------------------------------------------------------------")
	for _, h := range hosts {
		status := "No password"
		hasPass := h.EncryptedPassword != ""
		if hasPass {
			status = "Encrypted (RSA/AES)"
		}
		fmt.Printf("  %-16s %-16s %-12s %s\n", h.Alias, h.IP, h.Username, status)
	}
	fmt.Println()
}
