package cmdssh

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"golang.org/x/crypto/ssh"
)

type sshMvOptions struct {
	SrcNode     string
	SrcPath     string
	DestNode    string
	DestPath    string
	HasSrcNode  bool
	HasDestNode bool
	IsAllNodes  bool
	IsHelp      bool
}

// RunSSHMvExtended executes multi-node and local-fleet move operations with atomic remote cleanup.
func RunSSHMvExtended(args []string) error {
	opts := parseSSHMvArgs(args)
	if opts.IsHelp || opts.SrcPath == "" || opts.DestPath == "" {
		printSSHMvUsage()

		return nil
	}

	return dispatchSSHMv(opts)
}

func dispatchSSHMv(opts sshMvOptions) error {
	if opts.HasSrcNode {
		return dispatchRemoteSourceMv(opts)
	}

	return dispatchLocalSourceMv(opts)
}

func dispatchRemoteSourceMv(opts sshMvOptions) error {
	if opts.HasDestNode {
		return moveRemoteToRemote(opts.SrcNode, opts.SrcPath, opts.DestNode, opts.DestPath)
	}

	return moveRemoteToLocal(opts.SrcNode, opts.SrcPath, opts.DestPath)
}

func dispatchLocalSourceMv(opts sshMvOptions) error {
	if opts.IsAllNodes {
		return moveLocalToAllNodes(opts.SrcPath, opts.DestPath)
	}

	if opts.HasDestNode {
		return moveLocalToRemote(opts.SrcPath, opts.DestNode, opts.DestPath)
	}

	return apperror.NewValidationError("local move requires target node or fleet syntax")
}

func moveLocalToRemote(srcPath, destNode, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read local file '%s'", srcPath))
	}

	conn, errConn := findNodeConnection(destNode)
	if errConn != nil {
		return errConn
	}

	if errUpload := uploadFileToNode(conn, destPath, data); errUpload != nil {
		return errUpload
	}

	_ = os.Remove(srcPath)
	fmt.Printf("  %s✓ Moved %s -> %s:%s (%d bytes)%s\n",
		constants.ColorGreen, srcPath, destNode, destPath, len(data), constants.ColorReset)

	return nil
}

func uploadFileToNode(conn db.SSHConnection, destPath string, data []byte) error {
	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return apperror.NewExecutionError(fmt.Sprintf("connect to node %s", conn.Alias))
	}
	defer client.Close()

	return StreamFileToRemote(client, destPath, data, conn.OS)
}

func moveRemoteToLocal(srcNode, srcPath, destPath string) error {
	conn, errConn := findNodeConnection(srcNode)
	if errConn != nil {
		return errConn
	}

	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return apperror.NewExecutionError(fmt.Sprintf("connect to node %s", srcNode))
	}
	defer client.Close()

	data, errPull := StreamFileFromRemote(client, srcPath, conn.OS)
	if errPull != nil {
		return errPull
	}

	if errWrite := writeLocalFile(destPath, data); errWrite != nil {
		return apperror.WrapSimple(errWrite, "write local file")
	}

	if errDel := deleteRemoteFile(client, srcPath, isWindowsOS(conn.OS)); errDel != nil {
		return errDel
	}

	fmt.Printf("  %s✓ Moved %s:%s -> %s (%d bytes)%s\n",
		constants.ColorGreen, srcNode, srcPath, destPath, len(data), constants.ColorReset)

	return nil
}

func moveRemoteToRemote(srcNode, srcPath, destNode, destPath string) error {
	srcConn, errSrc := findNodeConnection(srcNode)
	if errSrc != nil {
		return errSrc
	}

	destConn, errDest := findNodeConnection(destNode)
	if errDest != nil {
		return errDest
	}

	srcClient, isSrcConnected := connectSSHClient(srcConn)
	if !isSrcConnected {
		return apperror.NewExecutionError(fmt.Sprintf("connect to source node %s", srcNode))
	}
	defer srcClient.Close()

	data, errPull := StreamFileFromRemote(srcClient, srcPath, srcConn.OS)
	if errPull != nil {
		return errPull
	}

	if errUpload := uploadFileToNode(destConn, destPath, data); errUpload != nil {
		return errUpload
	}

	if errDel := deleteRemoteFile(srcClient, srcPath, isWindowsOS(srcConn.OS)); errDel != nil {
		return errDel
	}

	fmt.Printf("  %s✓ Moved %s:%s -> %s:%s (%d bytes)%s\n",
		constants.ColorGreen, srcNode, srcPath, destNode, destPath, len(data), constants.ColorReset)

	return nil
}

func moveLocalToAllNodes(srcPath, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read local source file '%s'", srcPath))
	}

	online, errNodes := queryOnlineConnections()
	if errNodes != nil {
		return errNodes
	}

	if errBroad := broadcastDataToNodes(online, destPath, data); errBroad != nil {
		return errBroad
	}

	_ = os.Remove(srcPath)
	fmt.Printf("  %s✓ Local source '%s' removed after fleet move%s\n",
		constants.ColorGreen, srcPath, constants.ColorReset)

	return nil
}

func deleteRemoteFile(client *ssh.Client, remotePath string, isWin bool) error {
	cmd := buildDeleteRemoteCmd(remotePath, isWin)
	_, err := secrets.RunCommand(client, cmd, "")
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("delete remote file %s", remotePath))
	}

	return nil
}

func buildDeleteRemoteCmd(remotePath string, isWin bool) string {
	if isWin {
		return fmt.Sprintf("powershell.exe -NoProfile -Command \"Remove-Item -Force -Path '%s'\"", remotePath)
	}

	return fmt.Sprintf("rm -f '%s'", remotePath)
}

func parseSSHMvArgs(args []string) sshMvOptions {
	var opts sshMvOptions
	var posArgs []string
	var targetVal string

	for i := 0; i < len(args); i++ {
		i, posArgs, targetVal = processMvArg(args, i, &opts, posArgs, targetVal)
	}

	return populateMvPositionalOpts(opts, posArgs, targetVal)
}

func processMvArg(args []string, i int, opts *sshMvOptions, pos []string, target string) (int, []string, string) {
	arg := args[i]
	if isHelpToken(arg) {
		opts.IsHelp = true
		return i, pos, target
	}
	if arg == "--all" {
		opts.IsAllNodes = true
		return i, pos, target
	}

	return processTargetOrPos(args, i, pos, target)
}

func populateMvPositionalOpts(opts sshMvOptions, pos []string, target string) sshMvOptions {
	if len(pos) < 2 {
		return opts
	}

	srcNode, srcPath, hasSrcNode := parseNodePath(pos[0])
	destNode, destPath, hasDestNode := parseNodePath(pos[1])
	opts.SrcNode = srcNode
	opts.SrcPath = srcPath
	opts.HasSrcNode = hasSrcNode
	opts.DestPath = destPath

	return resolveMvDestNode(opts, destNode, hasDestNode, target)
}

func resolveMvDestNode(opts sshMvOptions, destNode string, hasDestNode bool, target string) sshMvOptions {
	if hasDestNode {
		opts.DestNode = destNode
		opts.HasDestNode = true
		if destNode == "all" {
			opts.IsAllNodes = true
		}
		return opts
	}

	if target != "" {
		opts.DestNode = target
		opts.HasDestNode = true
		if target == "all" {
			opts.IsAllNodes = true
		}
	}

	return opts
}

func printSSHMvUsage() {
	fmt.Printf("\n%sUsage:%s gitmap mv <src> <dest> [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Multi-node file move with atomic remote cleanup.")
	fmt.Println()
	fmt.Println("Move variants:")
	fmt.Println("  gitmap mv <local-path> <node>:<remote-path>     Local to Node Move")
	fmt.Println("  gitmap mv <local-path> <dest> -t <node>         Local to Node with flag")
	fmt.Println("  gitmap mv <local-path> all:<dest>               Local to Fleet Move")
	fmt.Println("  gitmap mv <node>:<remote-path> <local-path>     Node to Local Move")
	fmt.Println("  gitmap mv <node1>:<path> <node2>:<path>         Node to Node Move")
	fmt.Println()
}
