// Package cmdssh — ssh_deploy_folder.go implements parallel folder deployment and worker pooling.
package cmdssh

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
)

type conflictPromptSession struct {
	mu             sync.Mutex
	isAllOverwrite bool
	isAllSkip      bool
}

func (s *conflictPromptSession) prompt(relPath string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.isAllOverwrite {
		return true
	}
	if s.isAllSkip {
		return false
	}
	return s.askUser(relPath)
}

func (s *conflictPromptSession) askUser(relPath string) bool {
	fmt.Printf("  File '%s' exists on remote. Overwrite? [y/N/a(ll)/s(kip all)]: ", relPath)
	var answer string
	if _, err := fmt.Scanln(&answer); err != nil {
		return false
	}
	clean := strings.TrimSpace(answer)
	return s.applyAnswer(clean)
}

func (s *conflictPromptSession) applyAnswer(clean string) bool {
	if strings.EqualFold(clean, "a") || strings.EqualFold(clean, "all") {
		s.isAllOverwrite = true
		return true
	}
	if strings.EqualFold(clean, "s") || strings.EqualFold(clean, "skip all") || strings.EqualFold(clean, "skip") {
		s.isAllSkip = true
		return false
	}
	return strings.EqualFold(clean, "y") || strings.EqualFold(clean, "yes")
}

// scanLocalFolder recursively scans the rootPath and returns all files as LocalFileInfo items.
func scanLocalFolder(rootPath string) ([]LocalFileInfo, error) {
	cleanRoot := filepath.Clean(rootPath)
	info, err := os.Stat(cleanRoot)
	if err != nil {
		return nil, apperror.WrapSimple(err, "scanLocalFolder.stat")
	}
	if !info.IsDir() {
		return nil, apperror.NewValidationError("path is not a directory: " + cleanRoot)
	}
	return walkDirectoryFiles(cleanRoot)
}

func walkDirectoryFiles(root string) ([]LocalFileInfo, error) {
	var items []LocalFileInfo
	walker := makeFileWalker(root, &items)
	if err := filepath.Walk(root, walker); err != nil {
		return nil, apperror.WrapSimple(err, "walkDirectoryFiles")
	}
	return items, nil
}

func makeFileWalker(root string, items *[]LocalFileInfo) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		*items = append(*items, buildLocalFileInfo(path, rel, info))
		return nil
	}
}

func buildLocalFileInfo(absPath, relPath string, info os.FileInfo) LocalFileInfo {
	normRel := filepath.ToSlash(relPath)
	return LocalFileInfo{
		RelPath: normRel,
		AbsPath: absPath,
		IsDir:   false,
		Size:    info.Size(),
		ModTime: info.ModTime(),
	}
}

// executeParallelFolderDeploy deploys items across a parallel worker pool.
func executeParallelFolderDeploy(client *ssh.Client, node db.SSHConnection, items []LocalFileInfo, destRoot string, opts DeployOptions) (DeployResult, error) {
	workerCount := resolveWorkerCount(opts.Parallel, len(items))
	jobs := make(chan LocalFileInfo, len(items))
	results := make(chan FileDeployResult, len(items))
	var wg sync.WaitGroup
	session := &conflictPromptSession{}

	spawnDeployWorkers(client, node, destRoot, opts, session, jobs, results, workerCount, &wg)
	feedJobsAndClose(items, jobs)
	go awaitWorkersAndClose(&wg, results)

	mode := resolveDeploySyncMode(opts)
	return collectFolderDeployResults(node, opts.SourcePath, destRoot, mode, results), nil
}

func resolveWorkerCount(requested, itemCount int) int {
	if itemCount <= 0 {
		return 1
	}
	workers := requested
	if workers <= 0 {
		workers = 4
	}
	if workers > 16 {
		workers = 16
	}
	if workers > itemCount {
		return itemCount
	}
	return workers
}

func spawnDeployWorkers(client *ssh.Client, node db.SSHConnection, destRoot string, opts DeployOptions, session *conflictPromptSession, jobs <-chan LocalFileInfo, results chan<- FileDeployResult, count int, wg *sync.WaitGroup) {
	for w := 0; w < count; w++ {
		wg.Add(1)
		go runFolderDeployWorker(client, node, destRoot, opts, session, jobs, results, wg)
	}
}

func feedJobsAndClose(items []LocalFileInfo, jobs chan<- LocalFileInfo) {
	for _, item := range items {
		jobs <- item
	}
	close(jobs)
}

func awaitWorkersAndClose(wg *sync.WaitGroup, results chan<- FileDeployResult) {
	wg.Wait()
	close(results)
}

func runFolderDeployWorker(client *ssh.Client, node db.SSHConnection, destRoot string, opts DeployOptions, session *conflictPromptSession, jobs <-chan LocalFileInfo, results chan<- FileDeployResult, wg *sync.WaitGroup) {
	defer wg.Done()
	isWin := isWindowsOS(node.OS)
	for item := range jobs {
		res := deploySingleFolderItem(client, node, destRoot, item, opts, session, isWin)
		results <- res
	}
}

func deploySingleFolderItem(client *ssh.Client, node db.SSHConnection, destRoot string, item LocalFileInfo, opts DeployOptions, session *conflictPromptSession, isWin bool) FileDeployResult {
	destPath := joinRemotePath(destRoot, item.RelPath, isWin)
	remoteInfo, err := probeRemoteFileInfo(client, destPath, isWin)
	if err != nil {
		return FileDeployResult{RelPath: item.RelPath, Success: false, Error: err.Error()}
	}
	mode := resolveDeploySyncMode(opts)
	action, err := evaluateFileConflict(item, remoteInfo, mode)
	if err != nil {
		return FileDeployResult{RelPath: item.RelPath, Success: false, Error: err.Error()}
	}
	return executeDeployAction(client, node, item, destPath, remoteInfo, action, opts, session, isWin)
}

func executeDeployAction(client *ssh.Client, node db.SSHConnection, item LocalFileInfo, destPath string, remote RemoteFileInfo, action DeployAction, opts DeployOptions, session *conflictPromptSession, isWin bool) FileDeployResult {
	if action == ActionSkip {
		return FileDeployResult{RelPath: item.RelPath, Action: ActionSkip, Success: true}
	}
	if action == ActionConflictPrompt {
		return handleConflictPromptAction(client, node, item, destPath, remote, opts, session)
	}
	if action == ActionTransferToLocal {
		return transferRemoteToLocal(client, destPath, item.AbsPath, isWin)
	}
	return transferLocalToRemote(client, node, item, destPath, opts)
}

func handleConflictPromptAction(client *ssh.Client, node db.SSHConnection, item LocalFileInfo, destPath string, remote RemoteFileInfo, opts DeployOptions, session *conflictPromptSession) FileDeployResult {
	if !opts.IsInteractive {
		return buildConflictResult(item, destPath, remote)
	}
	if session != nil && session.prompt(item.RelPath) {
		return transferLocalToRemote(client, node, item, destPath, opts)
	}
	return FileDeployResult{RelPath: item.RelPath, Action: ActionSkip, Success: true}
}

func buildConflictResult(item LocalFileInfo, destPath string, remote RemoteFileInfo) FileDeployResult {
	record := DeployConflictRecord{
		Source:      item.RelPath,
		Destination: destPath,
		LocalMtime:  item.ModTime.Unix(),
		RemoteMtime: remote.ModTime.Unix(),
		LocalSize:   item.Size,
		RemoteSize:  remote.Size,
		Reason:      "remote file exists and no sync mode specified",
	}
	return FileDeployResult{RelPath: item.RelPath, Action: ActionConflictPrompt, Success: true, Conflict: &record}
}

func transferLocalToRemote(client *ssh.Client, node db.SSHConnection, item LocalFileInfo, destPath string, opts DeployOptions) FileDeployResult {
	if opts.IsDryRun {
		return FileDeployResult{RelPath: item.RelPath, Action: ActionTransferToRemote, Bytes: item.Size, Success: true}
	}
	data, err := os.ReadFile(item.AbsPath)
	if err != nil {
		return FileDeployResult{RelPath: item.RelPath, Action: ActionTransferToRemote, Success: false, Error: err.Error()}
	}
	if err := StreamFileToRemote(client, destPath, data, node.OS); err != nil {
		return FileDeployResult{RelPath: item.RelPath, Action: ActionTransferToRemote, Success: false, Error: err.Error()}
	}
	return FileDeployResult{RelPath: item.RelPath, Action: ActionTransferToRemote, Bytes: int64(len(data)), Success: true}
}

func transferRemoteToLocal(client *ssh.Client, remotePath, localPath string, isWin bool) FileDeployResult {
	data, err := pullRemoteFileBytes(client, remotePath, isWin)
	if err != nil {
		return FileDeployResult{RelPath: remotePath, Action: ActionTransferToLocal, Success: false, Error: err.Error()}
	}
	if err := saveLocalFileBytes(localPath, data); err != nil {
		return FileDeployResult{RelPath: remotePath, Action: ActionTransferToLocal, Success: false, Error: err.Error()}
	}
	return FileDeployResult{RelPath: remotePath, Action: ActionTransferToLocal, Bytes: int64(len(data)), Success: true}
}

func saveLocalFileBytes(localPath string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return err
	}
	return os.WriteFile(localPath, data, 0644)
}

func pullRemoteFileBytes(client *ssh.Client, remotePath string, isWin bool) ([]byte, error) {
	cmd := buildPullFileCmd(remotePath, isWin)
	b64Out, err := secrets.RunCommand(client, cmd, "")
	if err != nil {
		return nil, apperror.WrapSimple(err, "pullRemoteFileBytes")
	}
	data, err := base64.StdEncoding.DecodeString(strings.TrimSpace(b64Out))
	if err != nil {
		return nil, apperror.WrapSimple(err, "pullRemoteFileBytes.decode")
	}
	return data, nil
}

func buildPullFileCmd(remotePath string, isWin bool) string {
	if isWin {
		escaped := strings.ReplaceAll(remotePath, "'", "''")
		return fmt.Sprintf("powershell.exe -NoProfile -Command \"[Convert]::ToBase64String([System.IO.File]::ReadAllBytes('%s'))\"", escaped)
	}
	escaped := strings.ReplaceAll(remotePath, "'", "'\\''")
	return fmt.Sprintf("base64 < '%s'", escaped)
}

func collectFolderDeployResults(node db.SSHConnection, src, dest string, mode DeploySyncMode, results <-chan FileDeployResult) DeployResult {
	res := makeInitialFolderDeployResult(node, src, dest, mode)
	for r := range results {
		res.TotalFiles++
		res.FileResults = append(res.FileResults, r)
		updateFolderDeployMetrics(&res, r)
	}
	return res
}

func makeInitialFolderDeployResult(node db.SSHConnection, src, dest string, mode DeploySyncMode) DeployResult {
	return DeployResult{
		Target:      node.Alias,
		IP:          node.IPAddress,
		Source:      src,
		Destination: dest,
		Mode:        string(mode),
		Success:     true,
	}
}

func updateFolderDeployMetrics(res *DeployResult, r FileDeployResult) {
	if !r.Success {
		res.Success = false
		return
	}
	applySuccessfulFileMetric(res, r)
}

func applySuccessfulFileMetric(res *DeployResult, r FileDeployResult) {
	if r.Action == ActionSkip {
		res.SkippedFiles++
		return
	}
	if r.Action == ActionConflictPrompt {
		res.Conflicts++
		return
	}
	res.TransferredFiles++
	res.TransferredBytes += r.Bytes
}
