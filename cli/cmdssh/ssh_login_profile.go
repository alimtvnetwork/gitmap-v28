package cmdssh

import (
	"context"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	dbpkg "github.com/alimtvnetwork/gitmap-v28/cli/db"
	"golang.org/x/crypto/ssh"
)

func isNodeProfiled(ctx context.Context, target, ip string) bool {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return false
	}
	defer dbConn.Close()

	conn, err := dbpkg.GetSSHConnectionByAlias(ctx, dbConn.SQL(), target)
	if err == nil && conn != nil && conn.OSVersion != "" && !conn.FirstRunAt.IsZero() {
		return true
	}
	connIP, errIP := dbpkg.GetSSHConnectionByIP(ctx, dbConn.SQL(), ip)
	return errIP == nil && connIP != nil && connIP.OSVersion != "" && !connIP.FirstRunAt.IsZero()
}

func connectProbeClient(sshTarget *SSHTarget, password string) *ssh.Client {
	if password == "" {
		return tryConnectDefaultKey(sshTarget)
	}
	c, err := crypto.ConnectWithPassword(sshTarget.IP, sshTarget.Username, password)
	if err == nil && c != nil {
		return c
	}
	return tryConnectDefaultKey(sshTarget)
}

func ensureRemoteGitmapForLogin(client *ssh.Client, alias string) {
	osType := probeRemoteOSType(client)
	if !isRemoteGitmapInstalled(client, osType) {
		fmt.Printf("ℹ [%s] GitMap missing on remote machine, auto-installing...\n", alias)
		_ = bootstrapRemoteGitmap(client, alias, osType)
	}
}

func resolveOSGroupFromType(osType string) string {
	if isWindowsOS(osType) {
		return constants.OSGroupWindows
	}
	if osType == constants.OSTargetMac {
		return constants.OSGroupMac
	}
	return constants.OSGroupUnix
}

func queryOrFallbackReport(client *ssh.Client) *cmdos.OSInfoReport {
	rep, hasReport := probeRemoteGitmapWhichOS(client)
	if hasReport && rep != nil {
		return rep
	}
	osType, fullVer, arch := probeTargetOSDetails(client)
	return &cmdos.OSInfoReport{
		OSType:       osType,
		OSGroup:      resolveOSGroupFromType(osType),
		OSVersion:    fullVer,
		Architecture: arch,
	}
}

func buildProfileConnection(
	target string,
	sshTarget *SSHTarget,
	rep *cmdos.OSInfoReport,
	now time.Time,
) dbpkg.SSHConnection {
	return dbpkg.SSHConnection{
		Alias:             target,
		IPAddress:         sshTarget.IP,
		Username:          sshTarget.Username,
		EncryptedPassword: sshTarget.EncryptedPassword,
		KeyPath:           findDefaultUserSSHKey(),
		OS:                rep.OSType,
		OSGroup:           rep.OSGroup,
		OSVersion:         rep.OSVersion,
		BuildVersion:      rep.BuildVersion,
		FirstRunAt:        now,
		CreatedAt:         now,
	}
}

func saveProfileConnection(
	ctx context.Context,
	target string,
	sshTarget *SSHTarget,
	rep *cmdos.OSInfoReport,
) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		return
	}
	defer dbConn.Close()

	conn := buildProfileConnection(target, sshTarget, rep, time.Now().UTC())
	_ = dbpkg.InsertOrUpdateSSHConnection(ctx, dbConn.SQL(), conn)
}

func executeNodeLoginProfiling(
	ctx context.Context,
	target string,
	sshTarget *SSHTarget,
	client *ssh.Client,
) {
	fmt.Printf("ℹ First login to node %s (%s): running GitMap which-os profiling...\n", target, sshTarget.IP)
	ensureRemoteGitmapForLogin(client, target)
	rep := queryOrFallbackReport(client)
	fmt.Printf("✔ Node %s: Identified via GitMap which-os: %s (%s, %s, %s)\n",
		target, rep.OSType, rep.OSGroup, rep.OSVersion, rep.Architecture)
	saveProfileConnection(ctx, target, sshTarget, rep)
}

func probeAndEnsureNodeProfile(
	ctx context.Context,
	target string,
	sshTarget *SSHTarget,
	password string,
) {
	if isNodeProfiled(ctx, target, sshTarget.IP) {
		return
	}
	client := connectProbeClient(sshTarget, password)
	if client == nil {
		return
	}
	defer client.Close()

	executeNodeLoginProfiling(ctx, target, sshTarget, client)
}
