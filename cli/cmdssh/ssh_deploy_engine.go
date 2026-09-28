// Package cmdssh — ssh_deploy_engine.go coordinates single-file and folder deployment pipelines.
package cmdssh

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

// DeployResult aggregates execution statistics and outcomes.
type DeployResult struct {
	Target           string             `json:"target"`
	IP               string             `json:"ip"`
	Source           string             `json:"source"`
	Destination      string             `json:"destination"`
	Mode             string             `json:"mode"`
	TotalFiles       int                `json:"totalFiles"`
	TransferredFiles int                `json:"transferredFiles"`
	SkippedFiles     int                `json:"skippedFiles"`
	TransferredBytes int64              `json:"transferredBytes"`
	Conflicts        int                `json:"conflicts"`
	DurationMs       int64              `json:"durationMs"`
	Success          bool               `json:"success"`
	ErrorMessage     string             `json:"errorMessage,omitempty"`
	FileResults      []FileDeployResult `json:"fileResults,omitempty"`
}

// FileDeployResult records the action and outcome of an individual file.
type FileDeployResult struct {
	RelPath  string                `json:"relPath"`
	Action   DeployAction          `json:"action"`
	Bytes    int64                 `json:"bytes"`
	Success  bool                  `json:"success"`
	Error    string                `json:"error,omitempty"`
	Conflict *DeployConflictRecord `json:"conflict,omitempty"`
}

// ExecuteDeploy orchestrates verification, connectivity, path resolution, and file/folder deployment.
func ExecuteDeploy(opts DeployOptions) error {
	srcInfo, err := verifySourcePath(opts.SourcePath)
	if err != nil {
		return err
	}
	node, err := resolveDeployTargetNode(opts)
	if err != nil {
		return err
	}
	client, err := dialTargetNodeClient(node)
	if err != nil {
		return err
	}
	defer client.Close()
	return runDeployPipeline(client, node, srcInfo, opts)
}

func verifySourcePath(srcPath string) (os.FileInfo, error) {
	cleanSrc := filepath.Clean(srcPath)
	info, err := os.Stat(cleanSrc)
	if err != nil {
		return nil, apperror.NewNotFound("source_path", "E404", fmt.Sprintf("source path not found: %s", srcPath))
	}
	return info, nil
}

func resolveDeployTargetNode(opts DeployOptions) (db.SSHConnection, error) {
	if opts.Conn.Alias != "" || opts.Conn.IPAddress != "" {
		return opts.Conn, nil
	}
	nodes, err := resolveDeployNodes(opts.TargetToken)
	if err != nil {
		return db.SSHConnection{}, err
	}
	if len(nodes) == 0 {
		return db.SSHConnection{}, apperror.NewNotFound("target", "E404", "no target node matched: "+opts.TargetToken)
	}
	return nodes[0], nil
}

func dialTargetNodeClient(node db.SSHConnection) (*ssh.Client, error) {
	isOnline, reason := CheckConnLiveness(context.Background(), node.IPAddress, 22, 2*time.Second)
	if !isOnline {
		return nil, apperror.NewExecutionError(fmt.Sprintf("node [%s|%s] is offline: %s", node.Alias, node.IPAddress, reason))
	}
	client, isConnected := connectSSHClient(node)
	if !isConnected {
		return nil, apperror.NewExecutionError(fmt.Sprintf("ssh connection failed for [%s|%s]", node.Alias, node.IPAddress))
	}
	return client, nil
}

func runDeployPipeline(client *ssh.Client, node db.SSHConnection, srcInfo os.FileInfo, opts DeployOptions) error {
	resolvedDest := resolveRemoteDestination(client, node, opts.DestPath, filepath.Base(opts.SourcePath), srcInfo.IsDir())
	t0 := time.Now()
	res, err := dispatchDeployTarget(client, node, resolvedDest, srcInfo, opts)
	if err != nil {
		return err
	}
	res.DurationMs = time.Since(t0).Milliseconds()
	return outputDeployResult(res, opts)
}

func dispatchDeployTarget(client *ssh.Client, node db.SSHConnection, destPath string, srcInfo os.FileInfo, opts DeployOptions) (DeployResult, error) {
	if !srcInfo.IsDir() {
		return executeSingleFileDeploy(client, node, opts.SourcePath, destPath, srcInfo, opts)
	}
	return dispatchFolderDeploy(client, node, destPath, opts)
}

func dispatchFolderDeploy(client *ssh.Client, node db.SSHConnection, destPath string, opts DeployOptions) (DeployResult, error) {
	items, err := scanLocalFolder(opts.SourcePath)
	if err != nil {
		return DeployResult{}, err
	}
	return executeParallelFolderDeploy(client, node, items, destPath, opts)
}

func executeSingleFileDeploy(client *ssh.Client, node db.SSHConnection, srcPath, destPath string, srcInfo os.FileInfo, opts DeployOptions) (DeployResult, error) {
	isWin := isWindowsOS(node.OS)
	localItem := buildLocalFileInfo(srcPath, filepath.Base(srcPath), srcInfo)
	fileRes, mode, err := processSingleFileTransfer(client, node, localItem, destPath, opts, isWin)
	if err != nil {
		return DeployResult{}, err
	}
	res := buildInitialSingleDeployResult(node, srcPath, destPath, mode)
	applySingleDeployResult(&res, fileRes)
	return res, nil
}

func processSingleFileTransfer(client *ssh.Client, node db.SSHConnection, local LocalFileInfo, destPath string, opts DeployOptions, isWin bool) (FileDeployResult, DeploySyncMode, error) {
	remoteInfo, err := probeRemoteFileInfo(client, destPath, isWin)
	if err != nil {
		return FileDeployResult{}, SyncModeNone, err
	}
	mode := resolveDeploySyncMode(opts)
	action, err := evaluateFileConflict(local, remoteInfo, mode)
	if err != nil {
		return FileDeployResult{}, mode, err
	}
	session := &conflictPromptSession{}
	fileRes := executeDeployAction(client, node, local, destPath, remoteInfo, action, opts, session, isWin)
	return fileRes, mode, nil
}

func buildInitialSingleDeployResult(node db.SSHConnection, src, dest string, mode DeploySyncMode) DeployResult {
	return DeployResult{
		Target:      node.Alias,
		IP:          node.IPAddress,
		Source:      src,
		Destination: dest,
		Mode:        string(mode),
		TotalFiles:  1,
		Success:     true,
	}
}

func applySingleDeployResult(res *DeployResult, fileRes FileDeployResult) {
	res.FileResults = []FileDeployResult{fileRes}
	if !fileRes.Success {
		res.Success = false
		res.ErrorMessage = fileRes.Error
		return
	}
	applySingleDeploySuccess(res, fileRes)
}

func applySingleDeploySuccess(res *DeployResult, fileRes FileDeployResult) {
	if fileRes.Action == ActionSkip {
		res.SkippedFiles = 1
		return
	}
	if fileRes.Action == ActionConflictPrompt {
		res.Conflicts = 1
		return
	}
	res.TransferredFiles = 1
	res.TransferredBytes = fileRes.Bytes
}

func resolveRemoteDestination(client *ssh.Client, node db.SSHConnection, rawDest, srcBase string, isSrcDir bool) string {
	isWin := isWindowsOS(node.OS)
	baseDest := rawDest
	if !isAbsolutePath(rawDest, isWin) {
		baseDest = prefixRemoteWorkDir(rawDest, isWin)
	}
	if isBaseAppendNeeded(rawDest, srcBase, isSrcDir) {
		return joinRemotePath(baseDest, srcBase, isWin)
	}
	return baseDest
}

func isBaseAppendNeeded(rawDest, srcBase string, isSrcDir bool) bool {
	if isSrcDir {
		return strings.HasSuffix(rawDest, "/") || strings.HasSuffix(rawDest, "\\")
	}
	return isDirDestination(rawDest)
}

func prefixRemoteWorkDir(relPath string, isWin bool) string {
	clean := strings.TrimPrefix(strings.TrimPrefix(relPath, "./"), ".\\")
	workDir := "~"
	if isWin {
		workDir = "D:/work"
	}
	if clean == "" || clean == "." {
		return workDir
	}
	return workDir + "/" + clean
}

func isAbsolutePath(path string, isWin bool) bool {
	if strings.HasPrefix(path, "/") || strings.HasPrefix(path, "~") {
		return true
	}
	if isWin && len(path) >= 2 && path[1] == ':' {
		return true
	}
	if isWin && strings.HasPrefix(path, `\\`) {
		return true
	}
	return false
}

func outputDeployResult(res DeployResult, opts DeployOptions) error {
	if !opts.IsJSON {
		renderDeployBanner(res)
		return nil
	}
	if res.Conflicts > 0 {
		return RenderConflictJSONPrompt(extractDeployConflicts(res))
	}
	return RenderDeployJSONSummary(buildDeploySummaryJSON(res, opts))
}

func extractDeployConflicts(res DeployResult) []DeployConflictRecord {
	var records []DeployConflictRecord
	for _, f := range res.FileResults {
		if f.Action == ActionConflictPrompt && f.Conflict != nil {
			records = append(records, *f.Conflict)
		}
	}
	return records
}

func buildDeploySummaryJSON(res DeployResult, opts DeployOptions) DeploySummaryJSON {
	targetNode := DeployTargetNodeJSON{ID: opts.Host.ID, Alias: res.Target, IP: res.IP, OS: opts.Conn.OS, WorkDir: opts.RemoteWorkDir}
	metrics := DeployMetricsJSON{FilesProcessed: res.TotalFiles, FilesTransferred: res.TransferredFiles, FilesSkipped: res.SkippedFiles, BytesTransferred: res.TransferredBytes, DurationMs: res.DurationMs, ParallelWorkers: opts.Parallel}
	return DeploySummaryJSON{
		Status:     resolveDeployStatus(res.Success),
		Command:    opts.SubCmd,
		Direction:  res.Mode,
		TargetNode: targetNode,
		Metrics:    metrics,
		Transfers:  buildDeployFileRecords(res.FileResults, opts),
		Conflicts:  extractDeployConflicts(res),
		ExitCode:   resolveDeployExitCode(res),
	}
}

func resolveDeployExitCode(res DeployResult) int {
	if !res.Success {
		return 1
	}
	if res.Conflicts > 0 {
		return 3
	}
	return 0
}

func resolveDeployStatus(isSuccess bool) string {
	if isSuccess {
		return "success"
	}
	return "failed"
}

func buildDeployFileRecords(results []FileDeployResult, opts DeployOptions) []DeployFileRecord {
	var records []DeployFileRecord
	for _, r := range results {
		records = append(records, buildSingleDeployFileRecord(r, opts))
	}
	return records
}

func buildSingleDeployFileRecord(r FileDeployResult, opts DeployOptions) DeployFileRecord {
	status := "transferred"
	if r.Action == ActionSkip {
		status = "skipped"
	}
	if r.Action == ActionConflictPrompt {
		status = "conflict"
	}
	return DeployFileRecord{
		Source:      r.RelPath,
		Destination: r.RelPath,
		Bytes:       r.Bytes,
		Status:      status,
		Direction:   opts.SubCmd,
	}
}

func renderDeployBanner(res DeployResult) {
	fmt.Println()
	color, icon := resolveBannerStatusStyle(res.Success)
	fmt.Printf("  %s%s Deployment to [%s|%s] complete (%d ms)%s\n",
		color, icon, res.Target, res.IP, res.DurationMs, constants.ColorReset)
	fmt.Printf("  • Source:      %s\n", res.Source)
	fmt.Printf("  • Destination: %s\n", res.Destination)
	fmt.Printf("  • Mode:        %s\n", res.Mode)
	fmt.Printf("  • Files:       %d total, %s%d transferred%s, %s%d skipped%s\n",
		res.TotalFiles,
		constants.ColorGreen, res.TransferredFiles, constants.ColorReset,
		constants.ColorYellow, res.SkippedFiles, constants.ColorReset)
	printDeployBytesAndConflicts(res)
	fmt.Println()
}

func resolveBannerStatusStyle(isSuccess bool) (string, string) {
	if isSuccess {
		return constants.ColorGreen, "✓"
	}
	return constants.ColorRed, "✖"
}

func printDeployBytesAndConflicts(res DeployResult) {
	if res.TransferredBytes > 0 {
		fmt.Printf("  • Transferred: %s\n", formatByteSize(res.TransferredBytes))
	}
	if res.Conflicts > 0 {
		fmt.Printf("  • Conflicts:   %d\n", res.Conflicts)
	}
}

func formatByteSize(bytes int64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1f KB", float64(bytes)/1024)
	}
	return fmt.Sprintf("%.2f MB", float64(bytes)/(1024*1024))
}
