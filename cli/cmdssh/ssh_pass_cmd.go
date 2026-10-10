package cmdssh

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

// SSHPassCmd represents the 'pass' command under ssh.
var SSHPassCmd = &cobra.Command{
	Use:     "pass [show|ls|encrypt|decrypt] <alias|ip|password> [flags]",
	Aliases: []string{"password"},
	Short:   "Inspect, review, encrypt, or decrypt SSH node passwords",
	RunE: func(cmd *cobra.Command, args []string) error {
		return RunSSHPassCLI(args)
	},
}

// RunSSHPassCLI routes pass subcommands (show, ls, encrypt, decrypt, or direct target).
func RunSSHPassCLI(args []string) error {
	hasHelp := len(args) > 0 && (args[0] == "--help" || args[0] == "-h" || args[0] == "help")
	if hasHelp || len(args) == 0 {
		return printPassHelp()
	}

	sub := strings.ToLower(args[0])
	if sub == "ls" || sub == "list" {
		return runPassList()
	}
	if sub == "encrypt" || sub == "enc" {
		return runPassEncrypt(args[1:])
	}
	if sub == "decrypt" || sub == "dec" {
		return runPassDecrypt(args[1:])
	}

	return dispatchPassTarget(sub, args)
}

func resolvePassTarget(sub string, args []string) (string, error) {
	isShowAction := sub == "show" || sub == "view" || sub == "get"
	if isShowAction == false {
		return args[0], nil
	}
	if len(args) < 2 {
		return "", apperror.NewValidationError("missing node alias or IP. Usage: gitmap ssh pass show <alias|ip>")
	}

	return args[1], nil
}

func dispatchPassTarget(sub string, args []string) error {
	target, err := resolvePassTarget(sub, args)
	if err != nil {
		return err
	}

	return runPassShow(target)
}

func printPassHelp() error {
	fmt.Println("Usage: gitmap ssh pass <show|ls|encrypt|decrypt> [alias|ip|password] [flags]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  show <alias|ip>                   Reveal the encrypted saved password for a node")
	fmt.Println("  ls, list                          List all nodes with encrypted password status")
	fmt.Println("  encrypt, enc <password> [flags]   Encrypt a password using salted rotation cipher or Caesar")
	fmt.Println("  decrypt, dec <ciphertext>         Decrypt a cipher string (salt:, caesar:, or aes:)")
	fmt.Println()
	fmt.Println("Encryption Flags:")
	fmt.Println("  --salt <salt>     Provide custom salt (default: auto-generated 8-byte hex)")
	fmt.Println("  --caesar, --shift <N> Use Caesar rotation cipher with shift N")
	fmt.Println("  --aes             Use AES-GCM encryption with machine-bound key")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh pass encrypt rtyrty123@")
	fmt.Println("  gitmap ssh pass encrypt rtyrty123@ --salt 9f8b2c4e")
	fmt.Println("  gitmap ssh pass encrypt rtyrty123@ --caesar 13")
	fmt.Println("  gitmap ssh pass decrypt salt:9f8b2c4e:UypnMmpIN3BMLg==")
	fmt.Println("  gitmap ssh pass show u2")
	fmt.Println("  gitmap ssh pass ls")

	return nil
}

func parseCaesarShiftArg(val string) int {
	s, err := strconv.Atoi(val)
	if err == nil {
		return s
	}
	return 0
}

func parsePassEncryptFlags(args []string) (string, string, int, bool, bool, error) {
	if len(args) == 0 {
		return "", "", 0, false, false, apperror.NewValidationError("missing plaintext password. Usage: gitmap ssh pass encrypt <password> [--salt <salt>] [--caesar <shift>]")
	}
	plain := args[0]
	salt := ""
	caesarShift := 0
	useCaesar := false
	useAES := false

	for i := 1; i < len(args); i++ {
		arg := args[i]
		if arg == "--salt" && i+1 < len(args) {
			salt = args[i+1]
			i++
			continue
		}
		if (arg == "--caesar" || arg == "--shift") && i+1 < len(args) {
			useCaesar = true
			caesarShift = parseCaesarShiftArg(args[i+1])
			i++
			continue
		}
		if arg == "--aes" {
			useAES = true
			continue
		}
	}
	return plain, salt, caesarShift, useCaesar, useAES, nil
}

func encryptPassAES(plain string) (string, string, error) {
	enc, err := EncryptSSHPassword(plain)
	if err != nil {
		return "", "", apperror.WrapSimple(err, "encrypt AES")
	}
	return enc, "AES-GCM (random nonce)", nil
}

func encryptPassSalted(plain, salt string) (string, string) {
	effectiveSalt := salt
	if effectiveSalt == "" {
		effectiveSalt = secrets.GenerateRandomSalt(8)
	}
	cipherText := EncryptSaltedPassword(plain, effectiveSalt)
	modeDesc := fmt.Sprintf("Salted rotation cipher (salt=%s)", effectiveSalt)
	return cipherText, modeDesc
}

func handlePassEncryptAES(plain string) error {
	cipherText, modeDesc, err := encryptPassAES(plain)
	if err != nil {
		return err
	}
	renderEncryptedPasswordResult(plain, cipherText, modeDesc)
	return nil
}

func runPassEncrypt(args []string) error {
	plain, salt, caesarShift, useCaesar, useAES, err := parsePassEncryptFlags(args)
	if err != nil {
		return err
	}

	if useAES {
		return handlePassEncryptAES(plain)
	}
	if useCaesar {
		cipherText := EncryptCaesarPassword(plain, caesarShift, salt)
		modeDesc := fmt.Sprintf("Caesar rotation cipher (shift=%d, salt=%q)", caesarShift, salt)
		renderEncryptedPasswordResult(plain, cipherText, modeDesc)
		return nil
	}
	cipherText, modeDesc := encryptPassSalted(plain, salt)
	renderEncryptedPasswordResult(plain, cipherText, modeDesc)
	return nil
}

func renderEncryptedPasswordResult(plain, cipherText, modeDesc string) {
	fmt.Println()
	fmt.Printf("  %s● Encrypted Password%s (%s)\n", constants.ColorGreen, constants.ColorReset, modeDesc)
	fmt.Printf("    Plaintext:  %s%s%s\n", constants.ColorYellow, plain, constants.ColorReset)
	fmt.Printf("    Ciphertext: %s%s%s\n", constants.ColorCyan, cipherText, constants.ColorReset)
	fmt.Println()
	fmt.Println("    JSON configuration variable:")
	fmt.Printf("      \"winPass\": \"%s\",\n", cipherText)
	fmt.Println()
}

func runPassDecrypt(args []string) error {
	if len(args) == 0 {
		return apperror.NewValidationError("missing ciphertext to decrypt. Usage: gitmap ssh pass decrypt <ciphertext>")
	}
	cipherText := args[0]
	plain, err := DecryptSSHPassword(cipherText)
	if err != nil {
		return apperror.WrapSimple(err, "decrypt ciphertext")
	}

	fmt.Println()
	fmt.Printf("  %s● Decrypted Password%s\n", constants.ColorGreen, constants.ColorReset)
	fmt.Printf("    Ciphertext: %s%s%s\n", constants.ColorCyan, cipherText, constants.ColorReset)
	fmt.Printf("    Plaintext:  %s%s%s\n", constants.ColorYellow, plain, constants.ColorReset)
	fmt.Println()
	return nil
}

func runPassShow(target string) error {
	encPass, hostDesc, err := resolveNodeEncryptedPassword(target)
	if err != nil {
		return err
	}

	hasPass := encPass != ""
	if hasPass == false {
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
	dbConn, err := openSSHDBFunc()
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
	conn, connErr := db.GetSSHConnectionByAlias(ctx, sqlDB, target)
	if connErr == nil && conn != nil {
		return conn.EncryptedPassword, fmt.Sprintf("%s@%s", conn.Username, conn.IPAddress), nil
	}

	connIP, ipErr := db.GetSSHConnectionByIP(ctx, sqlDB, target)
	if ipErr == nil && connIP != nil {
		return connIP.EncryptedPassword, fmt.Sprintf("%s@%s", connIP.Username, connIP.IPAddress), nil
	}

	return "", "", apperror.NewNotFoundError("node not found: " + target)
}

func runPassList() error {
	dbConn, err := openSSHDBFunc()
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
