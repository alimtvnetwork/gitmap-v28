// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type syncSettingsResult struct {
	alias      string
	hostAlias  string
	osType     string
	isSuccess  bool
	errMessage string
	duration   time.Duration
}

// RunNodesPushSettings exports local Antigravity settings and pushes them to a single remote node.
func RunNodesPushSettings(args []string) error {
	if len(args) == 0 {
		return apperror.NewSimple("target node alias required: gitmap nodes push-settings <node-alias>", "E9050")
	}

	targetAlias := args[0]
	conn, err := resolveSSHConnectionByAlias(targetAlias)
	if err != nil {
		return err
	}

	exportPath, data, errExport := exportLocalAgySettings()
	if errExport != nil {
		return errExport
	}

	start := time.Now()
	destPath, pushErr := pushSettingsToNode(conn, data)
	if pushErr != nil {
		return pushErr
	}

	printPushSettingsSuccess(conn.Alias, destPath, exportPath, time.Since(start))
	return nil
}

// RunNodesSyncSettings broadcasts local Antigravity settings across all reachable fleet nodes.
func RunNodesSyncSettings(_ []string) error {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil || len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found to synchronize settings", "E9051")
	}

	_, data, errExport := exportLocalAgySettings()
	if errExport != nil {
		return errExport
	}

	var results []syncSettingsResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, c := range conns {
		if isLocalMachineConnection(c) {
			continue
		}
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			res := executeSyncSettingsWorker(conn, data)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	renderSyncSettingsTable(results)
	return nil
}

func resolveSSHConnectionByAlias(alias string) (db.SSHConnection, error) {
	conns, err := cmdssh.FetchAllSSHConnections()
	if err != nil {
		return db.SSHConnection{}, apperror.WrapSimple(err, "fetch ssh connections")
	}

	for _, c := range conns {
		if strings.EqualFold(c.Alias, alias) {
			return c, nil
		}
	}

	return db.SSHConnection{}, apperror.NewNotFound("node alias", "E9052", fmt.Sprintf("node alias %q not found in fleet registry", alias))
}

func exportLocalAgySettings() (string, []byte, error) {
	tempDir := filepath.Join(".ai-memory", "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", nil, apperror.WrapSimple(err, "ensure temp directory")
	}

	exportPath := filepath.Join(tempDir, "agy-settings-export.json")
	if err := cmdagy.ExecuteAgySettingsExport(exportPath); err != nil {
		return "", nil, apperror.WrapSimple(err, "export antigravity settings")
	}

	data, err := os.ReadFile(exportPath)
	if err != nil {
		return "", nil, apperror.WrapSimple(err, "read exported settings file")
	}

	return exportPath, data, nil
}

func resolveRemoteStagingPath(conn db.SSHConnection) string {
	if isWindowsNode(conn) {
		return `C:\Windows\Temp\agy_settings_sync.json`
	}
	return "/tmp/agy_settings_sync.json"
}

func pushSettingsToNode(conn db.SSHConnection, data []byte) (string, error) {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return "", apperror.WrapSimple(errConnect, fmt.Sprintf("dial node %q", conn.Alias))
	}
	defer client.Close()

	destPath := resolveRemoteStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, destPath, data, conn.OS); errStream != nil {
		return "", apperror.WrapSimple(errStream, fmt.Sprintf("stream settings to node %q", conn.Alias))
	}

	importCmd := fmt.Sprintf("gitmap agy settings import %s", destPath)
	shell := "sh"
	if isWindowsNode(conn) {
		shell = "ps"
	}

	out, errExec := crypto.RunCommand(client, importCmd, shell)
	if errExec != nil {
		return "", apperror.WrapSimple(errExec, fmt.Sprintf("remote import failed on %q: %s", conn.Alias, out))
	}

	return destPath, nil
}

func executeSyncSettingsWorker(conn db.SSHConnection, data []byte) syncSettingsResult {
	start := time.Now()
	hostLabel := conn.Alias
	if hostLabel == "" {
		hostLabel = "fleet-worker"
	}

	_, err := pushSettingsToNode(conn, data)
	if err != nil {
		return syncSettingsResult{
			alias:      conn.Alias,
			hostAlias:  hostLabel,
			osType:     conn.OS,
			isSuccess:  false,
			errMessage: err.Error(),
			duration:   time.Since(start),
		}
	}

	return syncSettingsResult{
		alias:     conn.Alias,
		hostAlias: hostLabel,
		osType:    conn.OS,
		isSuccess: true,
		duration:  time.Since(start),
	}
}

func printPushSettingsSuccess(alias, remotePath, stagingPath string, dur time.Duration) {
	fmt.Printf("\n  %s Successfully pushed Antigravity settings to node %s\n", constants.ColorGreen+"✓"+constants.ColorReset, alias)
	fmt.Printf("    • Target Node:  %s\n", alias)
	fmt.Printf("    • Staging File: %s\n", stagingPath)
	fmt.Printf("    • Remote Dest:  %s\n", remotePath)
	fmt.Printf("    • Duration:     %v\n\n", dur.Round(time.Millisecond))
}

func renderSyncSettingsTable(results []syncSettingsResult) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-12s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST ALIAS", "OS", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 70))

	for _, r := range results {
		statusStr := constants.ColorGreen + "● SUCCESS" + constants.ColorReset
		if !r.isSuccess {
			statusStr = constants.ColorRed + "▲ FAILED" + constants.ColorReset
		}
		fmt.Printf("  %-18s %-16s %-10s %-12s %-10v\n",
			r.alias, r.hostAlias, r.osType, statusStr, r.duration.Round(time.Millisecond))
	}
	fmt.Println()
}
