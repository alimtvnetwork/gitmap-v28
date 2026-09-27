// Package cmdssh — ssh_pull_inventory.go fetches and aggregates remote machine inventories into repo-secrets.
package cmdssh

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

var defaultMachineFolders = map[string]string{
	"w1": "04-w1-machine",
	"w2": "05-w2-machine",
	"w3": "06-w3-machine",
}

// RunSSHPullInventoryCLI orchestrates remote scans and inventory fetching into repo-secrets.
func RunSSHPullInventoryCLI(args []string) error {
	target := "all"
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		target = args[0]
	}

	nodes, err := resolveDeployNodes(target)
	if err != nil {
		return err
	}

	repoSecretsRoot := resolveRepoSecretsRoot()
	fmt.Printf("\n%s📡 GitMap Fleet Remote Inventory Aggregation%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  • Destination Root: %s\n", repoSecretsRoot)
	fmt.Printf("  • Target Nodes:     %d node(s)\n\n", len(nodes))

	for _, node := range nodes {
		pullSingleNodeInventory(node, repoSecretsRoot)
	}

	fmt.Printf("\n  %s✓ Fleet inventory aggregation complete.%s\n\n", constants.ColorGreen, constants.ColorReset)

	return nil
}

func resolveRepoSecretsRoot() string {
	candidates := []string{
		`D:\work\repo-secrets`,
		filepath.Join("..", "repo-secrets"),
		"repo-secrets",
	}

	for _, c := range candidates {
		if info, err := os.Stat(c); err == nil && info.IsDir() {
			return c
		}
	}

	return `D:\work\repo-secrets`
}

func pullSingleNodeInventory(node db.SSHConnection, secretsRoot string) {
	aliasLower := strings.ToLower(node.Alias)
	folderName, exists := defaultMachineFolders[aliasLower]
	if !exists {
		folderName = fmt.Sprintf("machine-%s", aliasLower)
	}

	destDir := filepath.Join(secretsRoot, folderName)
	_ = os.MkdirAll(destDir, 0755)

	fmt.Printf("  ▸ [%s|%s] Scanning and fetching inventory -> %s\n", node.Alias, node.IPAddress, folderName)
	client, isConnected := connectSSHClient(node)
	if !isConnected {
		fmt.Printf("    %s✖ Failed to connect via SSH%s\n", constants.ColorRed, constants.ColorReset)

		return
	}
	defer client.Close()

	runRemoteScanAndFetch(client, node, destDir)
}

func runRemoteScanAndFetch(client *ssh.Client, node db.SSHConnection, destDir string) {
	t0 := time.Now()
	scanCmd := "gitmap scan D:/work --output json"
	_, _ = crypto.RunCommand(client, scanCmd, "")

	fetchRemoteFile(client, node, "D:/work/.gitmap/output/gitmap.json", filepath.Join(destDir, "gitmap.json"))
	fetchRemoteFile(client, node, "D:/work/ooshutup10.cfg", filepath.Join(destDir, "ooshutup10.cfg"))
	copyFleetNodeConfigs(destDir)

	dur := time.Since(t0).Round(time.Millisecond)
	fmt.Printf("    %s✓ Inventory collected in %s%s\n", constants.ColorGreen, dur, constants.ColorReset)
}

func fetchRemoteFile(client *ssh.Client, node db.SSHConnection, remotePath, localPath string) bool {
	session, err := client.NewSession()
	if err != nil {
		return false
	}
	defer session.Close()

	localFile, err := os.Create(localPath)
	if err != nil {
		return false
	}
	defer localFile.Close()

	session.Stdout = localFile
	cmd := fmt.Sprintf(`powershell -NoProfile -Command "if (Test-Path '%s') { [System.IO.File]::OpenRead('%s').CopyTo([System.Console]::OpenStandardOutput()) }"`, remotePath, remotePath)
	if !isWindowsOS(node.OS) {
		cmd = fmt.Sprintf(`cat '%s' 2>/dev/null`, remotePath)
	}

	runErr := session.Run(cmd)
	stat, _ := localFile.Stat()
	if stat != nil && stat.Size() > 0 && runErr == nil {
		fmt.Printf("    %s✓ Fetched %s (%d bytes)%s\n", constants.ColorGreen, filepath.Base(localPath), stat.Size(), constants.ColorReset)

		return true
	}

	_ = os.Remove(localPath)

	return false
}

func copyFleetNodeConfigs(destDir string) {
	srcNodes := `D:\work\repo-secrets\gitmap-ssh-nodes.json`
	if data, err := os.ReadFile(srcNodes); err == nil {
		_ = os.WriteFile(filepath.Join(destDir, "gitmap-ssh-nodes.json"), data, 0644)
		_ = os.WriteFile(filepath.Join(destDir, "gitmap-ssh.json"), data, 0644)
	}
}
