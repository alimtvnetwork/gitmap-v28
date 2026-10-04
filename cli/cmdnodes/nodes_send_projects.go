// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

type sendProjectsNodeResult struct {
	alias       string
	hostAlias   string
	osType      string
	isSuccess   bool
	projsSynced int
	duration    time.Duration
	errMessage  string
}

// RunNodesSendProjects synchronizes workspace project registries across fleet nodes.
func RunNodesSendProjects(args []string) error {
	entries, err := loadLocalProjectEntries()
	if err != nil || len(entries) == 0 {
		return apperror.NewSimple("no workspace project entries found to send", "E9053")
	}

	conns, errConns := cmdssh.FetchAllSSHConnections()
	if errConns != nil || len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found", "E9054")
	}

	targetNode := parseTargetNodeFilter(args)
	if targetNode != "" && !strings.EqualFold(targetNode, "all") {
		return sendProjectsToSingleNode(targetNode, conns, entries)
	}

	return sendProjectsToAllNodes(conns, entries)
}

func parseTargetNodeFilter(args []string) string {
	for _, a := range args {
		if !strings.HasPrefix(a, "-") {
			return a
		}
	}
	return ""
}

func loadLocalProjectEntries() ([]vscodepm.Entry, error) {
	path, err := vscodepm.ProjectsJSONPath()
	if err == nil {
		data, readErr := os.ReadFile(path)
		if readErr == nil && len(data) > 0 {
			var entries []vscodepm.Entry
			if jsonErr := json.Unmarshal(data, &entries); jsonErr == nil && len(entries) > 0 {
				return entries, nil
			}
		}
	}

	entries, listErr := vscodepm.ListEntries()
	if listErr == nil && len(entries) > 0 {
		return entries, nil
	}

	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		return []vscodepm.Entry{
			{
				Name:     filepath.Base(cwd),
				RootPath: cwd,
				Paths:    []string{},
				Tags:     []string{"gitmap"},
				Enabled:  true,
			},
		}, nil
	}

	return nil, apperror.NewSimple("could not determine local project entries", "E9055")
}

func sendProjectsToSingleNode(alias string, conns []db.SSHConnection, entries []vscodepm.Entry) error {
	var targetConn *db.SSHConnection
	for _, c := range conns {
		if strings.EqualFold(c.Alias, alias) {
			targetConn = &c
			break
		}
	}

	if targetConn == nil {
		return apperror.NewNotFound("node alias", "E9056", fmt.Sprintf("node alias %q not found", alias))
	}

	start := time.Now()
	syncedCount, err := deployTranslatedProjects(*targetConn, entries)
	if err != nil {
		return err
	}

	printSendProjectsSuccess(targetConn.Alias, syncedCount, time.Since(start))
	return nil
}

func sendProjectsToAllNodes(conns []db.SSHConnection, entries []vscodepm.Entry) error {
	var results []sendProjectsNodeResult
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, c := range conns {
		if isLocalMachineConnection(c) {
			continue
		}
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			res := executeSendProjectsWorker(conn, entries)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(c)
	}

	wg.Wait()
	renderSendProjectsTable(results)
	return nil
}

func executeSendProjectsWorker(conn db.SSHConnection, entries []vscodepm.Entry) sendProjectsNodeResult {
	start := time.Now()
	hostLabel := conn.Alias
	if hostLabel == "" {
		hostLabel = "fleet-worker"
	}

	count, err := deployTranslatedProjects(conn, entries)
	if err != nil {
		return sendProjectsNodeResult{
			alias:       conn.Alias,
			hostAlias:   hostLabel,
			osType:      conn.OS,
			isSuccess:   false,
			projsSynced: 0,
			duration:    time.Since(start),
			errMessage:  err.Error(),
		}
	}

	return sendProjectsNodeResult{
		alias:       conn.Alias,
		hostAlias:   hostLabel,
		osType:      conn.OS,
		isSuccess:   true,
		projsSynced: count,
		duration:    time.Since(start),
	}
}

func deployTranslatedProjects(conn db.SSHConnection, entries []vscodepm.Entry) (int, error) {
	translated := translateProjectEntries(entries, conn)
	data, errMarshal := json.MarshalIndent(translated, "", "  ")
	if errMarshal != nil {
		return 0, apperror.WrapSimple(errMarshal, "marshal translated projects.json")
	}

	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return 0, apperror.WrapSimple(errConnect, fmt.Sprintf("dial node %q", conn.Alias))
	}
	defer client.Close()

	destPath := resolveRemoteProjectsStagingPath(conn)
	if errStream := cmdssh.StreamFileToRemote(client, destPath, data, conn.OS); errStream != nil {
		return 0, apperror.WrapSimple(errStream, fmt.Sprintf("stream projects to %q", conn.Alias))
	}

	if errExec := executeRemotePlacement(client, destPath, conn); errExec != nil {
		return 0, errExec
	}

	return len(translated), nil
}

func resolveRemoteProjectsStagingPath(conn db.SSHConnection) string {
	if isWindowsNode(conn) {
		return `C:\Windows\Temp\projects_sync.json`
	}
	return "/tmp/projects_sync.json"
}

func translateProjectEntries(entries []vscodepm.Entry, conn db.SSHConnection) []vscodepm.Entry {
	isWin := isWindowsNode(conn)
	translated := make([]vscodepm.Entry, len(entries))
	for i, e := range entries {
		t := e
		t.RootPath = translatePathForNode(e.RootPath, isWin)
		if len(e.Paths) > 0 {
			t.Paths = make([]string, len(e.Paths))
			for j, p := range e.Paths {
				t.Paths[j] = translatePathForNode(p, isWin)
			}
		}
		translated[i] = t
	}
	return translated
}

func translatePathForNode(path string, isWin bool) string {
	if path == "" {
		return path
	}
	base := filepath.Base(path)
	if base == "." || base == "/" || base == "\\" {
		base = "workspace"
	}
	if isWin {
		clean := strings.ReplaceAll(path, "/", "\\")
		if strings.HasPrefix(clean, `D:\work\`) || strings.HasPrefix(clean, `C:\`) {
			return clean
		}
		return `D:\work\` + base
	}
	clean := strings.ReplaceAll(path, "\\", "/")
	if strings.HasPrefix(clean, "~/") || strings.HasPrefix(clean, "/home/") || strings.HasPrefix(clean, "/root/") {
		return clean
	}
	return "~/work/" + base
}

func executeRemotePlacement(client *ssh.Client, stagingPath string, conn db.SSHConnection) error {
	if isWindowsNode(conn) {
		psCmd := fmt.Sprintf(`powershell.exe -NoProfile -Command "$appdata = [Environment]::GetFolderPath('ApplicationData'); $dir1 = Join-Path $appdata 'Code\User\globalStorage\alefragnani.project-manager'; $dir2 = Join-Path $appdata 'Cursor\User\globalStorage\alefragnani.project-manager'; foreach ($d in @($dir1, $dir2)) { if (-not (Test-Path $d)) { New-Item -ItemType Directory -Force -Path $d | Out-Null }; Copy-Item '%s' -Destination (Join-Path $d 'projects.json') -Force }"`, stagingPath)
		_, err := crypto.RunCommand(client, psCmd, "ps")
		return err
	}

	unixCmd := fmt.Sprintf(`sh -c "mkdir -p \"$HOME/.config/Code/User/globalStorage/alefragnani.project-manager\" \"$HOME/.config/Cursor/User/globalStorage/alefragnani.project-manager\" && cp '%s' \"$HOME/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json\" && cp '%s' \"$HOME/.config/Cursor/User/globalStorage/alefragnani.project-manager/projects.json\""`, stagingPath, stagingPath)
	_, err := crypto.RunCommand(client, unixCmd, "sh")
	return err
}

func printSendProjectsSuccess(alias string, count int, dur time.Duration) {
	fmt.Printf("\n  %s Synchronized %d workspace project(s) to node %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, count, alias)
	fmt.Printf("    • Target Node:     %s\n", alias)
	fmt.Printf("    • Projects Synced: %d\n", count)
	fmt.Printf("    • Duration:        %v\n\n", dur.Round(time.Millisecond))
}

func renderSendProjectsTable(results []sendProjectsNodeResult) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-10s %-12s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST ALIAS", "OS", "PROJECTS", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 78))

	for _, r := range results {
		statusStr := constants.ColorGreen + "● SUCCESS" + constants.ColorReset
		if !r.isSuccess {
			statusStr = constants.ColorRed + "▲ FAILED" + constants.ColorReset
		}
		fmt.Printf("  %-18s %-16s %-10s %-10d %-12s %-10v\n",
			r.alias, r.hostAlias, r.osType, r.projsSynced, statusStr, r.duration.Round(time.Millisecond))
	}
	fmt.Println()
}
