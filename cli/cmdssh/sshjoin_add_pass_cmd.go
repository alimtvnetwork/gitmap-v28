package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/spf13/cobra"
)

var SJAddWithPassCmd = &cobra.Command{
	Use:     "add-with-pass <user@ip|ip> [password] [alias] [flags]",
	Aliases: []string{"add-pass", "add-password"},
	Short:   "Enroll an SSH machine with RSA-encrypted password storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		return executeEnrollWithPassCLI(cmd.Context(), args)
	},
}

const msgMissingAddPassTarget = `missing target host

Usage:
  gitmap ssh-join add-with-pass <user@ip|ip> [password] [alias] [flags]
  gitmap sj add-with-pass <user@ip|ip> [password] [alias] [flags]
  gitmap ssh join add-with-pass <user@ip|ip> [password] [alias] [flags]

Examples:
  gitmap ssh-join add-with-pass alim@192.168.1.14 secret123 devbox
  gitmap ssh-join add-with-pass root@192.168.1.14 P@ssw0rd! prod-server
  gitmap ssh-join add-with-pass alim@192.168.1.14 devbox
  gitmap sj add-pass 192.168.1.50 mypass staging`

type addPassEnrollParams struct {
	targetRaw string
	password  string
	alias     string
}

func parseAddPassParams(args []string) (addPassEnrollParams, error) {
	if len(args) == 0 {
		return addPassEnrollParams{}, apperror.NewValidationError(msgMissingAddPassTarget)
	}
	p := addPassEnrollParams{targetRaw: args[0]}
	if len(args) > 1 {
		p.password = args[1]
	}
	if len(args) > 2 {
		p.alias = args[2]
	}
	return p, nil
}

func promptPasswordIfMissing(ctx context.Context, pass string) (string, error) {
	if pass != "" {
		return pass, nil
	}
	return PromptSSHPassword(ctx, "Enter SSH password: ", int(os.Stdin.Fd()))
}

func printAddPassSuccess(alias, target string) {
	fmt.Printf("✓ Machine '%s' (%s) joined successfully with password.\n", alias, target)
	fmt.Println("  Password encrypted and stored securely using SSH RSA key.")
	fmt.Printf("  Recall anytime: gitmap ssh %s\n", alias)
	fmt.Printf("  Or connect directly: gitmap ssh %s\n", target)
}

func printAddPassJSON(alias, target string) error {
	out := map[string]any{
		"status":       "success",
		"alias":        alias,
		"target":       target,
		"has_password": true,
		"encryption":   "rsa-oaep",
	}
	return json.NewEncoder(os.Stdout).Encode(out)
}

func completeAddPassEnrollment(alias, target string, isJSON bool) error {
	if isJSON {
		return printAddPassJSON(alias, target)
	}
	printAddPassSuccess(alias, target)
	return nil
}

func persistHostWithEncryptedPass(ctx context.Context, host store.SSHHost, hist store.SSHHistory) error {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("persistHostWithEncryptedPass", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	return store.EnrollSSHHost(ctx, host, hist, dbConn.SQL())
}

func executeEnrollWithPassCLI(ctx context.Context, args []string) error {
	params, err := parseAddPassParams(args)
	if err != nil {
		return err
	}
	password, err := promptPasswordIfMissing(ctx, params.password)
	if err != nil {
		return err
	}
	target, err := ParseSSHTarget(params.targetRaw, resolveDefaultUsername(), 22)
	if err != nil {
		return err
	}
	encPass, err := EncryptSSHPassword(password)
	if err != nil {
		return err
	}
	opts := &SSHJoinOptions{
		Target:   target,
		Alias:    resolveJoinAlias(params.alias, nil, target),
		Password: password,
	}
	host, hist := buildHostAndHistory(opts)
	host.EncryptedPassword = encPass
	host.Port = target.Port
	if err := persistHostWithEncryptedPass(ctx, host, hist); err != nil {
		return err
	}
	return completeAddPassEnrollment(host.Alias, target.String(), false)
}
