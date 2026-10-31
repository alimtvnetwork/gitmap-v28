package cmdssh

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcopy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"golang.org/x/crypto/ssh"
)

type sshCPOptions struct {
	SrcNode     string
	SrcPath     string
	DestNode    string
	DestPath    string
	HasSrcNode  bool
	HasDestNode bool
	IsAllNodes  bool
	IsRecursive bool
	IsHelp      bool
}

// RunSSHCPCLI handles the 8 copy syntax variants for cluster nodes and clipboard.
func RunSSHCPCLI(args []string) error {
	opts := parseSSHCPArgs(args)
	if opts.IsHelp {
		printSSHCPUsage()

		return nil
	}

	if isClipboardFallback(opts) {
		return cmdcopy.RunCopyCmd(args)
	}

	return dispatchSSHCP(opts)
}

func isClipboardFallback(opts sshCPOptions) bool {
	if opts.HasSrcNode || opts.HasDestNode || opts.IsAllNodes || opts.IsRecursive {
		return false
	}

	return true
}

func dispatchSSHCP(opts sshCPOptions) error {
	if opts.IsRecursive {
		return executeRecursiveCopy(opts)
	}

	if opts.HasSrcNode {
		return dispatchRemoteSourceCopy(opts)
	}

	return dispatchLocalSourceCopy(opts)
}

func dispatchRemoteSourceCopy(opts sshCPOptions) error {
	if opts.IsAllNodes {
		return copyRemoteToAllNodes(opts.SrcNode, opts.SrcPath, opts.DestPath)
	}

	if opts.HasDestNode {
		return copyRemoteToRemote(opts.SrcNode, opts.SrcPath, opts.DestNode, opts.DestPath)
	}

	return copyRemoteToLocal(opts.SrcNode, opts.SrcPath, opts.DestPath)
}

func dispatchLocalSourceCopy(opts sshCPOptions) error {
	if opts.IsAllNodes {
		return copyLocalToAllNodes(opts.SrcPath, opts.DestPath)
	}

	return copyLocalToRemote(opts.SrcPath, opts.DestNode, opts.DestPath)
}

func executeRecursiveCopy(opts sshCPOptions) error {
	if opts.HasSrcNode && !opts.HasDestNode {
		return copyRemoteDirToLocal(opts.SrcNode, opts.SrcPath, opts.DestPath)
	}

	if !opts.HasSrcNode && opts.IsAllNodes {
		return copyLocalDirToAllNodes(opts.SrcPath, opts.DestPath)
	}

	if !opts.HasSrcNode && opts.HasDestNode {
		return copyLocalDirToRemote(opts.SrcPath, opts.DestNode, opts.DestPath)
	}

	return apperror.NewValidationError("recursive copy between two remote nodes is unsupported")
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

	if writeErr := writeLocalFile(destPath, data); writeErr != nil {
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

	if err = writeRemoteFileBytes(conn, destPath, data); err != nil {
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

	if err = writeRemoteFileBytes(destConn, destPath, data); err != nil {
		return err
	}

	fmt.Printf("  %s✓ Copied %s:%s -> %s:%s (%d bytes)%s\n",
		constants.ColorGreen, srcNode, srcPath, destNode, destPath, len(data), constants.ColorReset)

	return nil
}

func copyLocalToAllNodes(srcPath, destPath string) error {
	data, err := os.ReadFile(srcPath)
	if err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("read local source file '%s'", srcPath))
	}

	online, errNodes := queryOnlineConnections()
	if errNodes != nil {
		return errNodes
	}

	return broadcastDataToNodes(online, destPath, data)
}

func copyRemoteToAllNodes(srcNode, srcPath, destPath string) error {
	srcConn, err := findNodeConnection(srcNode)
	if err != nil {
		return err
	}

	data, err := readRemoteFileBytes(srcConn, srcPath)
	if err != nil {
		return err
	}

	online, errNodes := queryOnlineConnections()
	if errNodes != nil {
		return errNodes
	}

	return broadcastDataToNodes(online, destPath, data)
}

func broadcastDataToNodes(nodes []db.SSHConnection, destPath string, data []byte) error {
	var count int
	for _, node := range nodes {
		if err := writeRemoteFileBytes(node, destPath, data); err == nil {
			count++
		}
	}

	fmt.Printf("  %s✓ Broadcast %s -> %d online nodes (%d bytes)%s\n",
		constants.ColorGreen, destPath, count, len(data), constants.ColorReset)

	return nil
}

func copyLocalDirToRemote(srcDir, destNode, destDir string) error {
	tarData, err := packLocalDirectoryTar(srcDir)
	if err != nil {
		return apperror.WrapSimple(err, "pack local directory")
	}

	conn, errConn := findNodeConnection(destNode)
	if errConn != nil {
		return errConn
	}

	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return apperror.NewExecutionError(fmt.Sprintf("connect to node %s", destNode))
	}
	defer client.Close()

	return streamDirTarToClient(client, destDir, tarData, isWindowsOS(conn.OS))
}

func copyLocalDirToAllNodes(srcDir, destDir string) error {
	tarData, err := packLocalDirectoryTar(srcDir)
	if err != nil {
		return apperror.WrapSimple(err, "pack local directory")
	}

	online, errNodes := queryOnlineConnections()
	if errNodes != nil {
		return errNodes
	}

	return broadcastDirTarToNodes(online, destDir, tarData)
}

func broadcastDirTarToNodes(nodes []db.SSHConnection, destDir string, tarData []byte) error {
	var count int
	for _, node := range nodes {
		client, isConnected := connectSSHClient(node)
		if isConnected {
			if err := streamDirTarToClient(client, destDir, tarData, isWindowsOS(node.OS)); err == nil {
				count++
			}
			client.Close()
		}
	}

	fmt.Printf("  %s✓ Broadcast dir %s -> %d online nodes (%d bytes)%s\n",
		constants.ColorGreen, destDir, count, len(tarData), constants.ColorReset)

	return nil
}

func streamDirTarToClient(client *ssh.Client, destDir string, tarData []byte, isWin bool) error {
	_ = ensureRemoteDir(client, destDir, isWin)
	session, err := client.NewSession()
	if err != nil {
		return apperror.WrapSimple(err, "NewSession")
	}
	defer session.Close()

	stdin, errPipe := session.StdinPipe()
	if errPipe != nil {
		return apperror.WrapSimple(errPipe, "StdinPipe")
	}

	return pipeTarToSession(session, stdin, destDir, tarData, isWin)
}

func pipeTarToSession(session *ssh.Session, stdin io.WriteCloser, destDir string, tarData []byte, isWin bool) error {
	var stderr bytes.Buffer
	session.Stderr = &stderr
	cmd := buildTarExtractCmd(destDir, isWin)
	if err := session.Start(cmd); err != nil {
		return apperror.WrapSimple(err, "Start tar extraction")
	}

	if _, err := stdin.Write(tarData); err != nil {
		_ = stdin.Close()
		return apperror.WrapSimple(err, "Write tar data")
	}
	_ = stdin.Close()

	if err := session.Wait(); err != nil {
		return apperror.WrapSimple(err, stderr.String())
	}

	return nil
}

func packLocalDirectoryTar(srcDir string) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	errWalk := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		return appendTarEntry(tw, srcDir, path, info)
	})
	if errWalk != nil {
		return nil, errWalk
	}
	if errClose := tw.Close(); errClose != nil {
		return nil, errClose
	}

	return buf.Bytes(), nil
}

func appendTarEntry(tw *tar.Writer, baseDir, fullPath string, info os.FileInfo) error {
	relPath, errRel := filepath.Rel(baseDir, fullPath)
	if errRel != nil || relPath == "." {
		return errRel
	}

	hdr, errHdr := tar.FileInfoHeader(info, "")
	if errHdr != nil {
		return errHdr
	}
	hdr.Name = filepath.ToSlash(relPath)
	if errWrite := tw.WriteHeader(hdr); errWrite != nil {
		return errWrite
	}
	if info.IsDir() {
		return nil
	}

	return writeTarFileData(tw, fullPath)
}

func writeTarFileData(tw *tar.Writer, filePath string) error {
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	_, errWrite := tw.Write(fileData)

	return errWrite
}

func copyRemoteDirToLocal(srcNode, remoteDir, localDestDir string) error {
	conn, errConn := findNodeConnection(srcNode)
	if errConn != nil {
		return errConn
	}

	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return apperror.NewExecutionError(fmt.Sprintf("connect to node %s", srcNode))
	}
	defer client.Close()

	tarData, errFetch := fetchRemoteDirTar(client, remoteDir, isWindowsOS(conn.OS))
	if errFetch != nil {
		return errFetch
	}

	return unpackTarToLocalDir(localDestDir, tarData)
}

func fetchRemoteDirTar(client *ssh.Client, remoteDir string, isWin bool) ([]byte, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, apperror.WrapSimple(err, "NewSession")
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr
	cmd := buildRemoteTarCreateCmd(remoteDir, isWin)
	if errRun := session.Run(cmd); errRun != nil {
		return nil, apperror.WrapSimple(errRun, stderr.String())
	}

	return stdout.Bytes(), nil
}

func buildRemoteTarCreateCmd(dir string, isWin bool) string {
	if isWin {
		return fmt.Sprintf("cmd.exe /c tar.exe -cf - -C \"%s\" .", dir)
	}

	return fmt.Sprintf("sh -c \"tar -cf - -C '%s' .\"", dir)
}

func unpackTarToLocalDir(destDir string, tarData []byte) error {
	tr := tar.NewReader(bytes.NewReader(tarData))
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if errExtract := extractSingleTarEntry(destDir, hdr, tr); errExtract != nil {
			return errExtract
		}
	}

	return nil
}

func extractSingleTarEntry(destDir string, hdr *tar.Header, tr io.Reader) error {
	targetPath := filepath.Join(destDir, filepath.FromSlash(hdr.Name))
	if hdr.FileInfo().IsDir() {
		return os.MkdirAll(targetPath, 0755)
	}

	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return err
	}

	return writeExtractedFile(targetPath, tr, hdr.FileInfo().Mode())
}

func writeExtractedFile(targetPath string, tr io.Reader, mode os.FileMode) error {
	data, err := io.ReadAll(tr)
	if err != nil {
		return err
	}

	return os.WriteFile(targetPath, data, mode)
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

func parseSSHCPArgs(args []string) sshCPOptions {
	var opts sshCPOptions
	var posArgs []string
	var targetVal string

	for i := 0; i < len(args); i++ {
		i, posArgs, targetVal = processCPArg(args, i, &opts, posArgs, targetVal)
	}

	return populatePositionalOpts(opts, posArgs, targetVal)
}

func processCPArg(args []string, i int, opts *sshCPOptions, pos []string, target string) (int, []string, string) {
	arg := args[i]
	if isHelpToken(arg) {
		opts.IsHelp = true
		return i, pos, target
	}
	if isRecursiveToken(arg) {
		opts.IsRecursive = true
		return i, pos, target
	}
	if arg == "--all" {
		opts.IsAllNodes = true
		return i, pos, target
	}

	return processTargetOrPos(args, i, pos, target)
}

func processTargetOrPos(args []string, i int, pos []string, target string) (int, []string, string) {
	arg := args[i]
	if isTargetOptionFlag(arg) {
		return extractTargetNext(args, i, pos)
	}
	if hasTargetPrefix(arg) {
		return i, pos, extractTargetValue(arg)
	}
	if !strings.HasPrefix(arg, "-") {
		return i, append(pos, arg), target
	}

	return i, pos, target
}

func isTargetOptionFlag(arg string) bool {
	return arg == "-t" || arg == "--target"
}

func hasTargetPrefix(arg string) bool {
	return strings.HasPrefix(arg, "-t=") || strings.HasPrefix(arg, "--target=")
}

func extractTargetNext(args []string, i int, pos []string) (int, []string, string) {
	if i+1 < len(args) {
		return i + 1, pos, args[i+1]
	}

	return i, pos, ""
}

func extractTargetValue(arg string) string {
	if strings.HasPrefix(arg, "-t=") {
		return strings.TrimPrefix(arg, "-t=")
	}

	return strings.TrimPrefix(arg, "--target=")
}

func isHelpToken(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func isRecursiveToken(arg string) bool {
	return arg == "-r" || arg == "--recursive"
}

func populatePositionalOpts(opts sshCPOptions, pos []string, target string) sshCPOptions {
	if len(pos) < 2 {
		return opts
	}

	srcNode, srcPath, hasSrcNode := parseNodePath(pos[0])
	destNode, destPath, hasDestNode := parseNodePath(pos[1])
	opts.SrcNode = srcNode
	opts.SrcPath = srcPath
	opts.HasSrcNode = hasSrcNode
	opts.DestPath = destPath

	return resolveDestNode(opts, destNode, hasDestNode, target)
}

func resolveDestNode(opts sshCPOptions, destNode string, hasDestNode bool, target string) sshCPOptions {
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

func queryAllConnections() ([]db.SSHConnection, error) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "open connection database")
	}
	defer dbConn.Close()

	res := db.GetSSHConnections(dbConn.Context(), dbConn.SQL())
	if res.IsFailure() {
		return nil, apperror.NewExecutionError(fmt.Sprintf("load ssh nodes: %v", res.Err))
	}

	return res.Data, nil
}

func findNodeConnection(target string) (db.SSHConnection, error) {
	conns, err := queryAllConnections()
	if err != nil {
		return db.SSHConnection{}, err
	}

	matched := filterConnectionsByTarget(conns, target)
	if len(matched) == 0 {
		return db.SSHConnection{}, apperror.NewValidationError(fmt.Sprintf("target node '%s' not found in cluster inventory", target))
	}

	return matched[0], nil
}

func queryOnlineConnections() ([]db.SSHConnection, error) {
	conns, err := queryAllConnections()
	if err != nil {
		return nil, err
	}

	online, _ := partitionOnlineOffline(conns)
	if len(online) == 0 {
		return nil, apperror.NewValidationError("no online cluster nodes found")
	}

	return online, nil
}

func readRemoteFileBytes(conn db.SSHConnection, remotePath string) ([]byte, error) {
	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return nil, apperror.NewExecutionError(fmt.Sprintf("failed to authenticate with remote node '%s'", conn.Alias))
	}
	defer client.Close()

	return StreamFileFromRemote(client, remotePath, conn.OS)
}

func writeRemoteFileBytes(conn db.SSHConnection, remotePath string, data []byte) error {
	client, isConnected := connectSSHClient(conn)
	if !isConnected {
		return apperror.NewExecutionError(fmt.Sprintf("failed to authenticate with remote node '%s'", conn.Alias))
	}
	defer client.Close()

	if err := StreamFileToRemote(client, remotePath, data, conn.OS); err != nil {
		return apperror.WrapSimple(err, fmt.Sprintf("write remote file '%s' on %s", remotePath, conn.Alias))
	}

	return nil
}

func ensureLocalDir(destPath string) error {
	dir := filepath.Dir(destPath)
	if dir == "" || dir == "." {
		return nil
	}

	return os.MkdirAll(dir, 0755)
}

func writeLocalFile(destPath string, data []byte) error {
	if err := ensureLocalDir(destPath); err != nil {
		return err
	}

	return os.WriteFile(destPath, data, 0644)
}

func printSSHCPUsage() {
	fmt.Printf("\n%sUsage:%s gitmap cp <src> <dest> [flags]\n\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("Bidirectional file and directory transfer across local and remote cluster nodes.")
	fmt.Println()
	fmt.Println("8 Copy Syntax Variants:")
	fmt.Println("  1. gitmap cp <local-src> <node>:<dest>             Push local file to node")
	fmt.Println("  2. gitmap cp <local-src> <dest> -t <node>          Push local file with target flag")
	fmt.Println("  3. gitmap cp <local-src> all:<dest>                Broadcast to all online nodes")
	fmt.Println("  4. gitmap cp <node>:<remote-src> <local-dest>      Pull remote file to local")
	fmt.Println("  5. gitmap cp <node1>:<src> <node2>:<dest>          Inter-node transfer")
	fmt.Println("  6. gitmap cp <node1>:<src> all:<dest>              Node to fleet broadcast")
	fmt.Println("  7. gitmap cp -r <local-dir> <node>:<dest-dir>      Recursive directory copy")
	fmt.Println("  8. gitmap cp <text> / gitmap copy                  Memory & clipboard fallback")
	fmt.Println()
}
