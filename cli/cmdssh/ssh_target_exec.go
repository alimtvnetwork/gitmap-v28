package cmdssh

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
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
	if isJSON {
		out := map[string]any{
			"status": "error",
			"host":   target,
			"error":  err.Error(),
		}
		_ = json.NewEncoder(os.Stdout).Encode(out)
	}
	return apperror.WrapSimple(err, "resolveConnectionForTarget")
}

func resolveConnectionForTarget(ctx context.Context, target string) (*db.SSHConnection, error) {
	conns, err := fetchAllSSHConnections()
	if err == nil {
		matched := filterConnectionsByTarget(conns, target)
		if len(matched) > 0 {
			return &matched[0], nil
		}
	}
	return queryFallbackConnection(ctx, target)
}

func queryFallbackConnection(ctx context.Context, target string) (*db.SSHConnection, error) {
	dbConn, err := openSSHDB()
	if err != nil {
		return nil, apperror.WrapSimple(err, "openSSHDB")
	}
	defer dbConn.Close()

	host, hostErr := store.GetHostByAlias(ctx, target, dbConn.Conn())
	if hostErr == nil {
		conn := convertSSHHostToConnection(host)
		return &conn, nil
	}
	hostIP, ipErr := store.GetHostByIP(ctx, target, dbConn.Conn())
	if ipErr == nil {
		conn := convertSSHHostToConnection(hostIP)
		return &conn, nil
	}
	return nil, apperror.NewNotFoundError(fmt.Sprintf("SSH host '%s' not found in registry", target))
}

func convertSSHHostToConnection(h store.SSHHost) db.SSHConnection {
	return db.SSHConnection{
		Alias:             h.Alias,
		IPAddress:         h.IP,
		Username:          h.Username,
		EncryptedPassword: h.EncryptedPassword,
		OS:                "linux",
	}
}

func executeRemoteTargetCommand(c db.SSHConnection, target string, args []string, isJSON bool) error {
	client, isConnected := connectSSHClient(c, fmt.Sprintf("[%s]", c.Alias))
	if !isConnected {
		err := fmt.Errorf("failed to connect or authenticate with host '%s'", target)
		return handleTargetNotFound(target, isJSON, err)
	}
	defer client.Close()

	c.OS = probeRemoteOSType(client)
	shellType, cmdStr, _ := resolveWorkerCommand(c, args)
	startTime := time.Now()
	out, err := crypto.RunCommand(client, cmdStr, shellType)
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
