package cmdos

import (
	"context"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/config"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func collectSSHFluidIdentities() []MachineIdentity {
	items := []MachineIdentity{resolveLocalMachineIdentity()}
	rootDB, err := openDB()
	if err != nil || rootDB == nil {
		return items
	}
	defer rootDB.Close()
	res := db.GetSSHConnections(context.Background(), rootDB.Conn())
	if res.IsFailure() {
		return items
	}
	for idx, conn := range res.Data {
		items = append(items, convertSSHConnToIdentity(idx+2, conn))
	}
	return items
}

func convertSSHConnToIdentity(seq int, conn db.SSHConnection) MachineIdentity {
	alias := strings.TrimSpace(conn.Alias)
	if len(alias) == 0 {
		alias = conn.IPAddress
	}
	nodeID := fmt.Sprintf("ssh-%02d", seq-1)
	prev, _ := config.GetVariable("ssh_prev", conn.IPAddress)
	osPlatform := formatOSPlatform(conn.OS)
	if conn.OSGroup != "" && conn.OSGroup != conn.OS {
		osPlatform = fmt.Sprintf("%s (%s)", conn.OSGroup, conn.OS)
	}
	return MachineIdentity{
		Sequence: seq, NodeID: nodeID, IPAddress: conn.IPAddress, Alias: alias,
		MachineName: alias, OSHostname: alias, PreviousAlias: prev,
		OSPlatform: osPlatform, OSVersion: conn.OSVersion, CurrentUser: conn.Username,
		GitMapVersion: conn.BuildVersion, Scope: "ssh",
	}
}

func executeSSHMachineAliasSet(mode string, args []string, isJSON bool) error {
	target, newVal := resolveSSHTargetAndValue(args)
	rootDB, err := openDB()
	if err != nil || rootDB == nil {
		return apperror.NewSimple("unable to open SSH database", "E_SSH_DB")
	}
	defer rootDB.Close()
	updateSSHConnRecord(rootDB, target, newVal)
	fmt.Printf("%s✔ Updated SSH fleet %s for %q -> %q%s\n", constants.ColorGreen, mode, target, newVal, constants.ColorReset)
	return renderFleetIfJSON(isJSON)
}

func executeSSHMachineAliasRevert(mode string, args []string, isJSON bool) error {
	target := "local-01"
	if len(args) > 0 {
		target = args[0]
	}
	prev, hasPrev := config.GetVariable("ssh_prev", target)
	if !hasPrev || len(prev) == 0 {
		prev = target
	}
	return executeSSHMachineAliasSet(mode, []string{target, prev}, isJSON)
}

func resolveSSHTargetAndValue(args []string) (string, string) {
	if len(args) >= 2 {
		return args[0], sanitizeMachineToken(args[1])
	}
	return "local-01", sanitizeMachineToken(args[0])
}

func updateSSHConnRecord(rootDB *store.DB, target, newVal string) {
	ctx := context.Background()
	_ = config.SetVariable("ssh_prev", target, target)
	_, _ = rootDB.Conn().ExecContext(ctx, `UPDATE SSHConnection SET Alias = ? WHERE Alias = ? OR IPAddress = ?`, newVal, target, target)
	_, _ = rootDB.Conn().ExecContext(ctx, `UPDATE ssh_hosts SET alias = ? WHERE alias = ? OR ip = ?`, newVal, target, target)
}
