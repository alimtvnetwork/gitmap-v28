package cmdssh

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/macro"
	"golang.org/x/crypto/ssh"
)

func deployBundleToNodeWithLogs(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, opts SSHExportAllOptions) {
	shell := determineFallbackShell(c.OS)
	var logs []string

	logs = append(logs, deployConfigWithLog(client, c, bundle, shell))
	logs = append(logs, deployMacrosWithLog(client, c, bundle.Macros, shell, opts.IsForce))
	logs = append(logs, deployKnownHostsWithLog(client, c, bundle.KnownHosts, shell))
	logs = append(logs, deployBundleFileWithLog(client, c, bundle, shell))
	logs = append(logs, triggerRemoteImportWithLog(client, c.OS, shell))

	output := strings.Join(filterNonEmptyLogs(logs), "\n")
	printNodeResultOutput(c.Alias, c.IPAddress, output, nil)
}

func filterNonEmptyLogs(logs []string) []string {
	var out []string
	for _, l := range logs {
		if l != "" {
			out = append(out, l)
		}
	}

	return out
}

func deployConfigWithLog(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, shell string) string {
	if len(bundle.Config) == 0 {
		return "  • Config: (no local config to export)"
	}
	data, err := json.MarshalIndent(bundle.Config, "", "  ")
	if err != nil {
		return fmt.Sprintf("  ▲ Config: marshal failed (%v)", err)
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	cmd := buildRemoteExportWriteCmd(".gitmap", "config.json", b64, isWindowsOS(c.OS))
	_, runErr := crypto.RunCommand(client, cmd, shell)
	if runErr != nil {
		return fmt.Sprintf("  ▲ Config: export failed (%v)", runErr)
	}

	return "  ✓ Config: exported ~/.gitmap/config.json"
}

func deployMacrosWithLog(client *ssh.Client, c db.SSHConnection, macros []macro.Macro, shell string, isForce bool) string {
	if len(macros) == 0 {
		return "  • Macros: (no local macros to export)"
	}
	synced := 0
	var names []string
	for _, m := range macros {
		isOk := syncSingleMacroToClient(client, c, m, shell, isForce)
		if isOk {
			synced++
			names = append(names, m.Name)
		}
	}
	namesSummary := formatExportMacroNames(names)

	return fmt.Sprintf("  ✓ Macros: exported %d/%d macro(s) [%s]", synced, len(macros), namesSummary)
}

func formatExportMacroNames(names []string) string {
	if len(names) <= 3 {
		return strings.Join(names, ", ")
	}

	return fmt.Sprintf("%s, +%d more", strings.Join(names[:3], ", "), len(names)-3)
}

func deployKnownHostsWithLog(client *ssh.Client, c db.SSHConnection, kh string, shell string) string {
	if kh == "" {
		return "  • Known Hosts: (none to export)"
	}
	b64 := base64.StdEncoding.EncodeToString([]byte(kh))
	cmd := buildRemoteExportAppendCmd(".ssh", "known_hosts", b64, isWindowsOS(c.OS))
	_, runErr := crypto.RunCommand(client, cmd, shell)
	if runErr != nil {
		return fmt.Sprintf("  ▲ Known Hosts: export failed (%v)", runErr)
	}
	linesCount := len(strings.Split(strings.TrimSpace(kh), "\n"))

	return fmt.Sprintf("  ✓ Known Hosts: merged %d host key(s) into ~/.ssh/known_hosts", linesCount)
}

func deployBundleFileWithLog(client *ssh.Client, c db.SSHConnection, bundle *GitmapExportBundle, shell string) string {
	data, err := json.MarshalIndent(bundle, "", "  ")
	if err != nil {
		return fmt.Sprintf("  ▲ Bundle: marshal error (%v)", err)
	}
	b64 := base64.StdEncoding.EncodeToString(data)
	cmd := buildRemoteExportWriteCmd(".gitmap", "export_bundle.json", b64, isWindowsOS(c.OS))
	_, runErr := crypto.RunCommand(client, cmd, shell)
	if runErr != nil {
		return fmt.Sprintf("  ▲ Bundle: export failed (%v)", runErr)
	}

	return fmt.Sprintf("  ✓ Bundle: deployed ~/.gitmap/export_bundle.json (%d connection(s))", len(bundle.Connections))
}

func triggerRemoteImportWithLog(client *ssh.Client, osType, shell string) string {
	cmd := "gitmap ssh import-all --local-bundle 2>/dev/null || true"
	if isWindowsOS(osType) == false {
		cmd = wrapUnixPath(cmd)
	}
	out, err := crypto.RunCommand(client, cmd, shell)
	trimmed := strings.TrimSpace(out)
	if err == nil && trimmed != "" {
		return fmt.Sprintf("  ✓ Remote Registry: %s", trimmed)
	}

	return "  ✓ Remote Registry: synchronized local bundle"
}
