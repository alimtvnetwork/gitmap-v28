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
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func extractTargetFlags(args []string) (bool, []string) {
	isJSON := false
	var cleanArgs []string
	for _, a := range args {
		if a == "--json" {
			isJSON = true
			continue
		}
		cleanArgs = append(cleanArgs, a)
	}
	return isJSON, cleanArgs
}

func handleDirectSSHTarget(ctx context.Context, parent *cobra.Command, target string, args []string) error {
	if ctx == nil {
		ctx = context.Background()
	}
	isJSON, cleanArgs := extractTargetFlags(args)
	hasCommand := len(cleanArgs) > 0
	if !hasCommand && !isJSON {
		return runSSHLogin(parent, []string{target}, ctx)
	}
	conn, err := resolveConnectionForTarget(ctx, target)
	if err != nil {
		return handleTargetNotFound(target, isJSON, err)
	}
	if !hasCommand && isJSON {
		return renderTargetNodeStatusJSON(ctx, target, *conn)
	}
	return executeRemoteTargetCommand(*conn, target, cleanArgs, isJSON)
}

func handleTargetNotFound(target string, isJSON bool, err error) error {
	sugg := collectSSHSuggestions(target)
	store.LogFailedCommand("ssh "+target, strings.Join(os.Args[1:], " "), "ssh", "E1001", err.Error(), sugg)
	if isJSON {
		out := map[string]any{
			"status":      "error",
			"host":        target,
			"error":       err.Error(),
			"suggestions": sugg,
		}
		_ = json.NewEncoder(os.Stdout).Encode(out)

		return apperror.WrapSimple(err, "resolveConnectionForTarget")
	}

	PrintSevenStepTroubleshootGuide(target, err)

	return apperror.WrapSimple(err, "resolveConnectionForTarget")
}

func resolveConnectionForTarget(ctx context.Context, target string) (*db.SSHConnection, error) {
	if conn := findMatchingTargetConnection(target); conn != nil {
		return conn, nil
	}
	return queryFallbackConnection(ctx, target)
}

func findMatchingTargetConnection(target string) *db.SSHConnection {
	conns, err := fetchAllSSHConnections()
	if err != nil {
		return nil
	}
	matched := filterConnectionsByTarget(conns, target)
	if len(matched) == 0 {
		return nil
	}
	return &matched[0]
}

func resolveTargetHostAndUser(target string) (string, string) {
	if !strings.Contains(target, "@") {
		return target, ""
	}
	parts := strings.Split(target, "@")
	return parts[len(parts)-1], parts[0]
}

func queryFallbackConnection(ctx context.Context, target string) (*db.SSHConnection, error) {
	dbConn, err := openSSHDB()
	if err != nil {
		return nil, apperror.WrapSimple(err, "openSSHDB")
	}
	defer dbConn.Close()

	cleanHost, explicitUser := resolveTargetHostAndUser(target)

	host, hostErr := store.GetHostByAlias(ctx, cleanHost, dbConn.Conn())
	if hostErr == nil {
		conn := convertSSHHostToConnection(host, explicitUser)
		return &conn, nil
	}
	hostIP, ipErr := store.GetHostByIP(ctx, cleanHost, dbConn.Conn())
	if ipErr == nil {
		conn := convertSSHHostToConnection(hostIP, explicitUser)
		return &conn, nil
	}
	return nil, apperror.NewNotFoundError(fmt.Sprintf("SSH host '%s' not found in registry", target))
}

func convertSSHHostToConnection(h store.SSHHost, explicitUser string) db.SSHConnection {
	user := h.Username
	if explicitUser != "" {
		user = explicitUser
	}
	return db.SSHConnection{
		Alias:             h.Alias,
		IPAddress:         h.IP,
		Username:          user,
		EncryptedPassword: h.EncryptedPassword,
		OS:                "linux",
	}
}

func executeRemoteTargetCommand(c db.SSHConnection, target string, args []string, isJSON bool) error {
	client, err := connectSSHClientWithErr(c, fmt.Sprintf("[%s]", c.Alias))
	if err != nil {
		return handleTargetNotFound(target, isJSON, err)
	}
	defer client.Close()

	c.OS = probeRemoteOSType(client)
	shellType, cmdStr, isDelegate := resolveWorkerCommand(c, args)
	if isDelegate && ensureDelegateInstalled(client, c) != nil {
		return apperror.NewExecutionError(fmt.Sprintf("failed to install gitmap delegate on host '%s'", target))
	}
	startTime := time.Now()
	out, err := secrets.RunCommand(client, cmdStr, shellType)
	durMs := time.Since(startTime).Milliseconds()
	exitCode := resolveProcessExitCode(err)

	if isJSON {
		return outputTargetJSONResult(target, c, cmdStr, out, exitCode, durMs)
	}
	fmt.Print(out)
	if exitCode != 0 {
		return apperror.NewExecutionError(fmt.Sprintf("remote command exited with code %d", exitCode))
	}
	return nil
}

func outputTargetJSONResult(target string, c db.SSHConnection, cmd, out string, exitCode int, durMs int64) error {
	res := map[string]any{
		"host":       target,
		"alias":      c.Alias,
		"ip":         c.IPAddress,
		"command":    cmd,
		"exitCode":   exitCode,
		"stdout":     out,
		"stderr":     "",
		"durationMs": durMs,
	}
	return json.NewEncoder(os.Stdout).Encode(res)
}

func renderTargetNodeStatusJSON(ctx context.Context, target string, c db.SSHConnection) error {
	isOnline, _ := CheckConnLiveness(ctx, c.IPAddress, 22, 500*time.Millisecond)
	status := "offline"
	if isOnline {
		status = "ready"
	}
	res := map[string]any{
		"host":   target,
		"alias":  c.Alias,
		"ip":     c.IPAddress,
		"port":   22,
		"user":   c.Username,
		"os":     c.OS,
		"online": isOnline,
		"status": status,
	}
	return json.NewEncoder(os.Stdout).Encode(res)
}
