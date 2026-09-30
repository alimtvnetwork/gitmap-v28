// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// DetectCloneFile inspects arguments or current directory for a manifest file.
func DetectCloneFile(args []string) (string, bool) {
	for _, arg := range args {
		if isManifestPath(arg) {
			return arg, true
		}
	}
	return inspectCurrentDirManifest()
}

func isManifestPath(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	if strings.HasSuffix(strings.ToLower(arg), ".json") {
		return isExistingFile(arg)
	}
	return false
}

func inspectCurrentDirManifest() (string, bool) {
	defaultManifest := "gitmap.json"
	if isExistingFile(defaultManifest) {
		return defaultManifest, true
	}
	return "", false
}

func isExistingFile(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

// ReadCloneFileBytes loads file bytes and returns the file basename.
func ReadCloneFileBytes(filePath string) ([]byte, string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, "", err
	}
	return data, filepath.Base(filePath), nil
}

// ResolveRemoteWorkDir returns the default fleet work directory based on node OS.
func ResolveRemoteWorkDir(osType string) string {
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	if isWin {
		return "D:/work"
	}
	return "~/work"
}

// ResolveRemoteDestPath computes the remote file path on a target node.
func ResolveRemoteDestPath(osType, fileName string) string {
	workDir := ResolveRemoteWorkDir(osType)
	isWin := strings.EqualFold(osType, "windows") || strings.EqualFold(osType, "win")
	if isWin {
		cleanDir := strings.ReplaceAll(workDir, "/", "\\")
		return cleanDir + "\\" + fileName
	}
	return strings.TrimRight(workDir, "/") + "/" + fileName
}

// StageFileToRemoteNode copies the manifest file to the remote node work directory.
func StageFileToRemoteNode(client *ssh.Client, conn db.SSHConnection, fileName string, data []byte) (string, error) {
	destPath := ResolveRemoteDestPath(conn.OS, fileName)
	err := cmdssh.StreamFileToRemote(client, destPath, data, conn.OS)
	return destPath, err
}
