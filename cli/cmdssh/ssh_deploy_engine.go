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
	RelPath string       `json:"relPath"`
	Action  DeployAction `json:"action"`
	Bytes   int64        `json:"bytes"`
	Success bool         `json:"success"`
	Error   string       `json:"error,omitempty"`
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
	resolvedDest := resolveRemoteDestination(client, node, opts.DestPath, filepath.Base(opts.SourcePath))
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
	remoteInfo, err := probeRemoteFileInfo(client, destPath, isWin)
	if err != nil {
		return DeployResult{}, err
	}
	mode := resolveDeploySyncMode(opts)
	action, err := evaluateFileConflict(localItem, remoteInfo, mode)
	if err != nil {
		return DeployResult{}, err
	}
	res := buildInitialSingleDeployResult(node, srcPath, destPath, mode)
	fileRes := executeDeployAction(client, node, localItem, destPath, action, opts, isWin)
	applySingleDeployResult(&res, fileRes)
	return res, nil
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

func resolveRemoteDestination(client *ssh.Client, node db.SSHConnection, rawDest, srcBase string) string {
	isWin := isWindowsOS(node.OS)
	baseDest := rawDest
	if !isAbsolutePath(rawDest, isWin) {
		baseDest = prefixRemoteWorkDir(rawDest, isWin)
	}
	if isDirDestination(rawDest) {
		return joinRemotePath(baseDest, srcBase, isWin)
	}
	return baseDest
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
	if opts.IsJSON {
		return RenderDeployJSONSummary(buildDeploySummaryJSON(res, opts))
	}
	renderDeployBanner(res)
	return nil
}

func buildDeploySummaryJSON(res DeployResult, opts DeployOptions) DeploySummaryJSON {
	exitCode := 0
	status := "success"
	if !res.Success {
		exitCode = 1
		status = "failed"
	}
	return DeploySummaryJSON{
		Status:    status,
		Command:   opts.SubCmd,
		Direction: res.Mode,
		TargetNode: DeployTargetNodeJSON{
			Alias:   res.Target,
			IP:      res.IP,
			OS:      opts.Conn.OS,
			WorkDir: opts.RemoteWorkDir,
		},
		Metrics: DeployMetricsJSON{
			FilesProcessed:   res.TotalFiles,
			FilesTransferred: res.TransferredFiles,
			FilesSkipped:     res.SkippedFiles,
			BytesTransferred: res.TransferredBytes,
			DurationMs:       res.DurationMs,
			ParallelWorkers:  opts.Parallel,
		},
		ExitCode: exitCode,
	}
}

func renderDeployBanner(res DeployResult) {
	fmt.Println()
	statusColor := constants.ColorGreen
	statusIcon := "✓"
	if !res.Success {
		statusColor = constants.ColorRed
		statusIcon = "✖"
	}
	fmt.Printf("  %s%s Deployment to [%s|%s] complete (%d ms)%s\n",
		statusColor, statusIcon, res.Target, res.IP, res.DurationMs, constants.ColorReset)
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
