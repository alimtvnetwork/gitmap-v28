package cmdssh

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunSSHCPCLI handles bidirectional file copy: local <-> remote or remote <-> remote.
func RunSSHCPCLI(args []string) error {
	if len(args) < 2 || hasHelpFlag(args) {
		printSSHCPUsage()
		return nil
	}

	srcNode, srcPath, hasSrcNode := parseNodePath(args[0])
	destNode, destPath, hasDestNode := parseNodePath(args[1])

	if hasSrcNode == false && hasDestNode == false {
		return runSSHCopyCLI(args)
	}

	if hasSrcNode && hasDestNode == false {
		return copyRemoteToLocal(srcNode, srcPath, destPath)
	}

	if hasSrcNode == false && hasDestNode {
		return copyLocalToRemote(srcPath, destNode, destPath)
	}

	return copyRemoteToRemote(srcNode, srcPath, destNode, destPath)
}

func parseNodePath(arg string) (string, string, bool) {
	if isWindowsDrivePath(arg) {
		return "", arg, false
	}
	idx := strings.Index(arg, ":")
	if idx <= 0 {
		return "", arg, false
	}
	node := arg[:idx]
	remotePath := arg[idx+1:]
	return node, remotePath, true
}

func isWindowsDrivePath(arg string) bool {
	if len(arg) < 2 {
		return false
	}
	if arg[1] != ':' {
		return false
	}
	return unicode.IsLetter(rune(arg[0]))
}

func copyRemoteToLocal(srcNode, srcPath, destPath string) error {
	conn, err := findNodeConnection(srcNode)
	if err != nil {
		return err
	}

	data, err := readRemoteFileBytes(conn, srcPath)
	if err != nil {
		return err
	}

	writeErr := writeLocalFile(destPath, data)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "write local file")
	}

	fmt.Printf("  %s✓ Copied %s:%s -> %s (%d bytes)%s\n",
		constants.ColorGreen, srcNode, srcPath, destPath, len(data), constants.ColorReset)
	return nil
}

func copyLocalToRemote(srcPath, destNode, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read local source file '%s'", srcPath))
	}

	conn, err := findNodeConnection(destNode)
	if err != nil {
		return err
	}

	err = writeRemoteFileBytes(conn, destPath, data)
	if err != nil {
		return err
	}

	fmt.Printf("  %s✓ Copied %s -> %s:%s (%d bytes)%s\n",
		constants.ColorGreen, srcPath, destNode, destPath, len(data), constants.ColorReset)
	return nil
}

func copyRemoteToRemote(srcNode, srcPath, destNode, destPath string) error {
	srcConn, err := findNodeConnection(srcNode)
	if err != nil {
		return err
	}

	destConn, err := findNodeConnection(destNode)
	if err != nil {
		return err
	}

	data, err := readRemoteFileBytes(srcConn, srcPath)
	if err != nil {
		return err
	}

	err = writeRemoteFileBytes(destConn, destPath, data)
	if err != nil {
		return err
	}

	fmt.Printf("  %s✓ Copied %s:%s -> %s:%s (%d bytes)%s\n",
		constants.ColorGreen, srcNode, srcPath, destNode, destPath, len(data), constants.ColorReset)
	return nil
}

func findNodeConnection(target string) (db.SSHConnection, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return db.SSHConnection{}, apperror.WrapSimple(err, "open connection database")
	}
	defer dbConn.Close()

	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return db.SSHConnection{}, apperror.NewExecutionError(fmt.Sprintf("load ssh nodes: %v", res.Err))
	}

	matched := filterConnectionsByTarget(res.Data, target)
	if len(matched) == 0 {
		return db.SSHConnection{}, apperror.NewValidationError(fmt.Sprintf("target node '%s' not found in cluster inventory", target))
	}

	return matched[0], nil
}

func readRemoteFileBytes(conn db.SSHConnection, remotePath string) ([]byte, error) {
	client, isConnected := connectSSHClient(conn)
	if isConnected == false {
		return nil, apperror.NewExecutionError(fmt.Sprintf("failed to authenticate with remote node '%s'", conn.Alias))
	}
	defer client.Close()

	readCmd := buildRemoteReadCmd(remotePath, isWindowsOS(conn.OS))
	shell := determineFallbackShell(conn.OS)
	out, err := crypto.RunCommand(client, readCmd, shell)
	if err != nil {
		return nil, apperror.WrapSimple(err, fmt.Sprintf("read remote file '%s' on %s", remotePath, conn.Alias))
	}

	cleanB64 := strings.ReplaceAll(strings.TrimSpace(out), "\r\n", "")
	cleanB64 = strings.ReplaceAll(cleanB64, "\n", "")
	data, decodeErr := base64.StdEncoding.DecodeString(cleanB64)
	if decodeErr != nil {
		return nil, apperror.WrapSimple(decodeErr, "decode remote file base64 payload")
	}

	return data, nil
}

func buildRemoteReadCmd(remotePath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf(`powershell -NoProfile -Command "[Convert]::ToBase64String([IO.File]::ReadAllBytes('%s'))"`, remotePath)
	}
	return fmt.Sprintf(`base64 < '%s'`, remotePath)
}

func writeRemoteFileBytes(conn db.SSHConnection, remotePath string, data []byte) error {
	client, isConnected := connectSSHClient(conn)
	if isConnected == false {
		return apperror.NewExecutionError(fmt.Sprintf("failed to authenticate with remote node '%s'", conn.Alias))
	}
	defer client.Close()

	writeCmd := buildRemoteWriteCmd(remotePath, data, isWindowsOS(conn.OS))
	shell := determineFallbackShell(conn.OS)
	_, err := crypto.RunCommand(client, writeCmd, shell)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("write remote file '%s' on %s", remotePath, conn.Alias))
	}

	return nil
}

func writeLocalFile(destPath string, data []byte) error {
	dir := filepath.Dir(destPath)
	if dir != "" && dir != "." {
		mkdirErr := os.MkdirAll(dir, 0755)
		if mkdirErr != nil {
			return mkdirErr
		}
	}
	return os.WriteFile(destPath, data, 0644)
}

func printSSHCPUsage() {
	fmt.Printf("\n%sUsage:%s gitmap ssh cp <src> <dest>\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Bidirectional file transfer across local and remote cluster nodes.")
	fmt.Println()
	fmt.Println("Formats:")
	fmt.Println("  gitmap ssh cp <local-path> <node>:<remote-path>     Push local file to remote node")
	fmt.Println("  gitmap ssh cp <node>:<remote-path> <local-path>     Pull remote file to local host")
	fmt.Println("  gitmap ssh cp <node1>:<path> <node2>:<path>         Transfer between remote nodes")
	fmt.Println("  gitmap ssh cp <local-path> <remote-path>            Broadcast to all online nodes")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  gitmap ssh cp ./app.tar.gz worker-1:/tmp/app.tar.gz")
	fmt.Println("  gitmap ssh cp worker-1:/etc/config.json ./local-config.json")
	fmt.Println("  gitmap ssh cp worker-1:/var/log/app.log worker-2:/tmp/worker1-app.log")
	fmt.Println()
}
