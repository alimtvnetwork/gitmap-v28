// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

func filterRemoteConnections(conns []db.SSHConnection, opts NodesCloneOptions) []db.SSHConnection {
	var remote []db.SSHConnection
	for _, c := range conns {
		if isLocalMachineConnection(c) {
			continue
		}
		if isConnectionExcluded(c, opts) {
			continue
		}
		remote = append(remote, c)
	}
	return remote
}

func isConnectionExcluded(c db.SSHConnection, opts NodesCloneOptions) bool {
	if opts.TargetFilter != "" && !strings.EqualFold(c.Alias, opts.TargetFilter) && c.IPAddress != opts.TargetFilter {
		return true
	}
	if opts.ExcludeFilter != "" && (strings.EqualFold(c.Alias, opts.ExcludeFilter) || c.IPAddress == opts.ExcludeFilter) {
		return true
	}
	return false
}

func isLocalMachineConnection(c db.SSHConnection) bool {
	if c.IPAddress == "127.0.0.1" || strings.EqualFold(c.IPAddress, "localhost") {
		return true
	}
	if strings.EqualFold(c.Alias, "local") || strings.EqualFold(c.Alias, "localhost") {
		return true
	}
	return isLocalHostIP(c.IPAddress)
}

func isLocalHostIP(ip string) bool {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return false
	}
	for _, a := range addrs {
		ipNet, ok := a.(*net.IPNet)
		if ok && ipNet.IP.String() == ip {
			return true
		}
	}
	return false
}

func buildRemoteWorkDirExecString(kindStr, cmdArgs string, isWin bool) string {
	if isWin {
		return fmt.Sprintf("Set-Location D:\\work; gitmap %s %s", kindStr, cmdArgs)
	}
	return fmt.Sprintf("cd ~/work && gitmap %s %s", kindStr, cmdArgs)
}

func buildRemoteExecString(opts NodesCloneOptions, fileName string, isWin bool) string {
	cmdArgs := strings.Join(opts.PassArgs, " ")
	kindStr := string(opts.Kind)
	if !opts.HasFile {
		return fmt.Sprintf("gitmap %s %s", kindStr, cmdArgs)
	}
	return buildRemoteWorkDirExecString(kindStr, cmdArgs, isWin)
}

func runRemoteNodeWorker(conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) RemoteCloneNodeResult {
	start := time.Now()
	res := RemoteCloneNodeResult{Alias: conn.Alias, Host: conn.IPAddress, Role: "worker"}
	client, isConnected := cmdssh.ConnectSSHClient(conn)
	if !isConnected {
		res.Status = "auth_failed"
		res.Error = "authentication failed or node offline"
		return res
	}
	defer client.Close()
	return executeRemoteCloneSession(client, conn, opts, fileBytes, fileName, res, start)
}

func stageRemoteFileIfNeeded(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) error {
	if !opts.HasFile || len(fileBytes) == 0 {
		return nil
	}
	_, err := StageFileToRemoteNode(client, conn, fileName, fileBytes)
	return err
}

func executeRemoteCloneSession(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string, res RemoteCloneNodeResult, start time.Time) RemoteCloneNodeResult {
	if errStage := stageRemoteFileIfNeeded(client, conn, opts, fileBytes, fileName); errStage != nil {
		res.Status = "failed"
		res.Error = "manifest staging failed: " + errStage.Error()
		return res
	}
	return runRemoteExecOverSSH(client, conn, opts, fileName, res, start)
}

func runRemoteExecOverSSH(client *ssh.Client, conn db.SSHConnection, opts NodesCloneOptions, fileName string, res RemoteCloneNodeResult, start time.Time) RemoteCloneNodeResult {
	isWin := strings.EqualFold(conn.OS, "windows") || strings.EqualFold(conn.OS, "win")
	shell := "bash"
	if isWin {
		shell = "ps"
	}
	cmdStr := buildRemoteExecString(opts, fileName, isWin)
	out, errRun := crypto.RunCommand(client, cmdStr, shell)
	res.Duration = time.Since(start)
	res.DurationMs = res.Duration.Milliseconds()
	res.Stdout = strings.TrimSpace(out)
	if errRun != nil {
		res.Status = "failed"
		res.Error = errRun.Error()
		return res
	}
	res.Status = "success"
	return res
}

func executeFleetNodesParallel(conns []db.SSHConnection, opts NodesCloneOptions, fileBytes []byte, fileName string) []RemoteCloneNodeResult {
	var wg sync.WaitGroup
	var mu sync.Mutex
	results := make([]RemoteCloneNodeResult, 0, len(conns))
	for _, c := range conns {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			r := runRemoteNodeWorker(conn, opts, fileBytes, fileName)
			mu.Lock()
			results = append(results, r)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return results
}
