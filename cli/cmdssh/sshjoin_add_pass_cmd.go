package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var SJAddWithPassCmd = &cobra.Command{
	Use:     "add-with-pass <user@ip|ip> [password] [alias] [flags]",
	Aliases: []string{"add-pass", "add-password"},
	Short:   "Enroll an SSH machine with RSA-encrypted password storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		if ctx == nil {
			ctx = context.Background()
		}
		return executeEnrollWithPassCLI(ctx, args)
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
	osType    string
	isJSON    bool
}

func parseAddPassParams(args []string) (addPassEnrollParams, error) {
	if len(args) == 0 {
		return addPassEnrollParams{}, apperror.NewValidationError(msgMissingAddPassTarget)
	}
	p := addPassEnrollParams{}
	var positionals []string
	idx := 0
	for idx < len(args) {
		a := args[idx]
		if a == "--json" {
			p.isJSON = true
			idx++
			continue
		}
		if (a == "--os" || a == "-o") && idx+1 < len(args) {
			p.osType = args[idx+1]
			idx += 2
			continue
		}
		if (a == "--password" || a == "--pass") && idx+1 < len(args) {
			p.password = args[idx+1]
			idx += 2
			continue
		}
		if (a == "--alias" || a == "-n") && idx+1 < len(args) {
			p.alias = args[idx+1]
			idx += 2
			continue
		}
		if !strings.HasPrefix(a, "-") {
			positionals = append(positionals, a)
		}
		idx++
	}
	if len(positionals) == 0 {
		return addPassEnrollParams{}, apperror.NewValidationError(msgMissingAddPassTarget)
	}
	p.targetRaw = positionals[0]
	if len(positionals) > 1 && p.password == "" {
		p.password = positionals[1]
	}
	if len(positionals) > 2 && p.alias == "" {
		p.alias = positionals[2]
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

func persistHostWithEncryptedPass(ctx context.Context, host store.SSHHost, hist store.SSHHistory, osType string) error {
	return persistHostWithEncryptedPassAndVersion(ctx, host, hist, osType, "")
}

func persistHostWithEncryptedPassAndVersion(ctx context.Context, host store.SSHHost, hist store.SSHHistory, osType, osVersion string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return apperror.New("persistHostWithEncryptedPass", "E_INTERNAL_ERROR", map[string]any{"cause": err.Error()})
	}
	defer dbConn.Close()
	pair := hostHistoryPair{host: host, hist: hist, osType: osType, osVersion: osVersion}
	if appErr := persistDualTables(ctx, dbConn.SQL(), pair); appErr != nil {
		return appErr
	}
	return nil
}

func executeEnrollWithPassCLI(ctx context.Context, args []string) error {
	if ctx == nil {
		ctx = context.Background()
	}
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
	osType := resolveDefaultOS(params.osType)
	if params.osType == "" {
		osType = probeHostOSOrFallback(target.IP, target.Port, target.Username, password, osType)
	}
	if err := persistHostWithEncryptedPass(ctx, host, hist, osType); err != nil {
		return err
	}
	return completeAddPassEnrollment(host.Alias, target.String(), params.isJSON)
}

func probeHostOSOrFallback(ip string, port int, user, pass, fallback string) string {
	if !probeTCPQuick(ip, port, 200*time.Millisecond) {
		return fallback
	}
	client, err := crypto.ConnectWithPassword(ip, user, pass)
	if err != nil || client == nil {
		return fallback
	}
	defer client.Close()
	return probeRemoteOSType(client)
}
