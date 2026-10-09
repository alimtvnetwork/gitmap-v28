// Package cmdnodes coordinates fleet-wide operations across unified cluster and SSH nodes.
package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/secrets"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// DeployRepoOptions stores configuration flags for fleet repository deployment.
type DeployRepoOptions struct {
	NodeFilterOptions
	Repo              string `json:"repo,omitempty"`
	Dest              string `json:"dest,omitempty"`
	WithIDEs          bool   `json:"with_ides"`
	WithPinned        bool   `json:"with_pinned"`
	WithConversations bool   `json:"with_conversations"`
	FromLocal         string `json:"from_local,omitempty"`
	Clone             bool   `json:"clone"`
	Clean             bool   `json:"clean"`
	ExcludeGitObjects bool   `json:"exclude_git_objects"`
	DryRun            bool   `json:"dry_run"`
	IsJSON            bool   `json:"is_json"`
}

// DeployRepoResult captures deployment telemetry to a fleet node.
type DeployRepoResult struct {
	NodeAlias     string        `json:"node_alias"`
	Host          string        `json:"host"`
	OS            string        `json:"os"`
	DestDir       string        `json:"dest_dir"`
	Status        string        `json:"status"`
	Latency       time.Duration `json:"latency"`
	LatencyMs     int64         `json:"latency_ms"`
	ArchiveSize   int64         `json:"archive_size,omitempty"`
	IDEs          []string      `json:"ides,omitempty"`
	Conversations int           `json:"conversations,omitempty"`
	Error         string        `json:"error,omitempty"`
}

// ParseDeployRepoOptions parses CLI arguments into DeployRepoOptions.
func ParseDeployRepoOptions(args []string) (DeployRepoOptions, error) {
	var opts DeployRepoOptions
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if isDeployRepoHelpFlag(arg) {
			return opts, nil
		}
		consumed, err := parseSingleDeployRepoFlag(&opts, args, i)
		if err != nil {
			return opts, err
		}
		if consumed > 0 {
			i += (consumed - 1)
			continue
		}
		parsePositionalDeployRepoArg(&opts, arg)
	}
	return opts, nil
}

func parseDeployRepoOptions(args []string) (DeployRepoOptions, error) {
	return ParseDeployRepoOptions(args)
}

func isDeployRepoHelpFlag(arg string) bool {
	return arg == "--help" || arg == "-h" || arg == "help"
}

func parsePositionalDeployRepoArg(opts *DeployRepoOptions, arg string) {
	if strings.HasPrefix(arg, "-") {
		return
	}
	if opts.Repo == "" {
		opts.Repo = arg
		return
	}
	if opts.Target == "" {
		opts.Target = arg
	}
}

func parseSingleDeployRepoFlag(opts *DeployRepoOptions, args []string, i int) (int, error) {
	arg := args[i]
	if isBooleanDeployRepoFlag(opts, arg) {
		return 1, nil
	}
	return parseParamDeployRepoFlag(opts, args, i)
}

func isBooleanDeployRepoFlag(opts *DeployRepoOptions, arg string) bool {
	switch arg {
	case "--include-main":
		opts.IncludeMain = true
		return true
	case "--open-only":
		opts.OpenOnly = true
		return true
	case "--with-ides":
		opts.WithIDEs = true
		return true
	case "--with-pinned":
		opts.WithPinned = true
		opts.WithIDEs = true
		return true
	case "--with-conversations":
		opts.WithConversations = true
		return true
	case "--clone":
		opts.Clone = true
		return true
	case "--clean":
		opts.Clean = true
		return true
	case "--exclude-git-objects", "--no-git-objects":
		opts.ExcludeGitObjects = true
		return true
	case "--dry-run", "-n":
		opts.DryRun = true
		return true
	case "--json", "-j":
		opts.IsJSON = true
		return true
	default:
		return false
	}
}

func parseParamDeployRepoFlag(opts *DeployRepoOptions, args []string, i int) (int, error) {
	if consumed := parseTargetDeployRepoFlag(opts, args, i); consumed > 0 {
		return consumed, nil
	}
	if consumed := parseDestDeployRepoFlag(opts, args, i); consumed > 0 {
		return consumed, nil
	}
	if consumed := parseFromLocalDeployRepoFlag(opts, args, i); consumed > 0 {
		return consumed, nil
	}
	if consumed := parseRepoTokenFlags(opts, args, i); consumed > 0 {
		return consumed, nil
	}
	return 0, nil
}

func parseTargetDeployRepoFlag(opts *DeployRepoOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--target=") {
		opts.Target = strings.TrimPrefix(arg, "--target=")
		return 1
	}
	if strings.HasPrefix(arg, "-t=") {
		opts.Target = strings.TrimPrefix(arg, "-t=")
		return 1
	}
	if (arg == "--target" || arg == "-t") && i+1 < len(args) {
		opts.Target = args[i+1]
		return 2
	}
	return 0
}

func parseDestDeployRepoFlag(opts *DeployRepoOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--dest=") {
		opts.Dest = strings.TrimPrefix(arg, "--dest=")
		return 1
	}
	if strings.HasPrefix(arg, "-d=") {
		opts.Dest = strings.TrimPrefix(arg, "-d=")
		return 1
	}
	if (arg == "--dest" || arg == "-d") && i+1 < len(args) {
		opts.Dest = args[i+1]
		return 2
	}
	return 0
}

func parseFromLocalDeployRepoFlag(opts *DeployRepoOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--from-local=") {
		opts.FromLocal = strings.TrimPrefix(arg, "--from-local=")
		return 1
	}
	if arg == "--from-local" && i+1 < len(args) {
		opts.FromLocal = args[i+1]
		return 2
	}
	return 0
}

func parseRepoTokenFlags(opts *DeployRepoOptions, args []string, i int) int {
	if consumed := parseExceptRepoFlag(opts, args, i); consumed > 0 {
		return consumed
	}
	return parseIncludeRepoFlag(opts, args, i)
}

func parseExceptRepoFlag(opts *DeployRepoOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--except=") {
		appendTokenList(&opts.Except, strings.TrimPrefix(arg, "--except="))
		return 1
	}
	if strings.HasPrefix(arg, "-e=") {
		appendTokenList(&opts.Except, strings.TrimPrefix(arg, "-e="))
		return 1
	}
	if strings.HasPrefix(arg, "--exclude=") {
		appendTokenList(&opts.Except, strings.TrimPrefix(arg, "--exclude="))
		return 1
	}
	if (arg == "--except" || arg == "-e" || arg == "--exclude") && i+1 < len(args) {
		appendTokenList(&opts.Except, args[i+1])
		return 2
	}
	return 0
}

func parseIncludeRepoFlag(opts *DeployRepoOptions, args []string, i int) int {
	arg := args[i]
	if strings.HasPrefix(arg, "--include=") {
		appendTokenList(&opts.Include, strings.TrimPrefix(arg, "--include="))
		return 1
	}
	if strings.HasPrefix(arg, "--accept=") {
		appendTokenList(&opts.Include, strings.TrimPrefix(arg, "--accept="))
		return 1
	}
	if (arg == "--include" || arg == "--accept") && i+1 < len(args) {
		appendTokenList(&opts.Include, args[i+1])
		return 2
	}
	return 0
}

func appendTokenList(target *[]string, raw string) {
	parts := strings.Split(raw, ",")
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			*target = append(*target, trimmed)
		}
	}
}

// FilterDeployRepoNodes filters fleet connections according to deploy repo options.
func FilterDeployRepoNodes(conns []db.SSHConnection, opts DeployRepoOptions) []db.SSHConnection {
	return FilterFleetNodes(conns, opts.NodeFilterOptions)
}

func resolveLocalRepo(repoIdentifier, fromLocal string) (*model.ScanRecord, string, error) {
	if fromLocal != "" {
		return resolveRepoFromPath(fromLocal)
	}
	if repoIdentifier != "" {
		return resolveRepoByIdentifier(repoIdentifier)
	}
	return resolveRepoFromCwd()
}

func resolveRepoFromPath(path string) (*model.ScanRecord, string, error) {
	stat, err := os.Stat(path)
	if err != nil || !stat.IsDir() {
		return nil, "", apperror.NewNotFound("path", "E9071", fmt.Sprintf("local path %q not found or not a directory", path))
	}
	absPath, errAbs := filepath.Abs(path)
	if errAbs != nil {
		absPath = path
	}
	rec := findRepoInStore(absPath)
	if rec != nil {
		return rec, absPath, nil
	}
	base := filepath.Base(absPath)
	return &model.ScanRecord{
		RepoName:     base,
		Slug:         base,
		AbsolutePath: absPath,
	}, absPath, nil
}

func resolveRepoByIdentifier(query string) (*model.ScanRecord, string, error) {
	rec := findRepoInStore(query)
	if isStoredRepoValid(rec) {
		return rec, rec.AbsolutePath, nil
	}
	stat, err := os.Stat(query)
	if err == nil && stat.IsDir() {
		return resolveRepoFromPath(query)
	}
	if rec != nil {
		return rec, rec.AbsolutePath, nil
	}
	return nil, "", apperror.NewNotFound("repository", "E9072", fmt.Sprintf("repository %q not found in store or filesystem", query))
}

func isStoredRepoValid(rec *model.ScanRecord) bool {
	if rec == nil || rec.AbsolutePath == "" {
		return false
	}
	stat, err := os.Stat(rec.AbsolutePath)
	return err == nil && stat.IsDir()
}

func resolveRepoFromCwd() (*model.ScanRecord, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, "", apperror.WrapSimple(err, "resolve current working directory")
	}
	return resolveRepoFromPath(cwd)
}

func findRepoInStore(query string) *model.ScanRecord {
	st, err := store.OpenDefault()
	if err != nil {
		return nil
	}
	defer st.Close()
	return queryRepoFromStore(st, query)
}

func queryRepoFromStore(st *store.DB, query string) *model.ScanRecord {
	if byPath, err := st.FindByPath(query); err == nil && len(byPath) > 0 {
		return &byPath[0]
	}
	if bySlug, err := st.FindBySlug(query); err == nil && len(bySlug) > 0 {
		return &bySlug[0]
	}
	return scanStoreReposForMatch(st, query)
}

func scanStoreReposForMatch(st *store.DB, query string) *model.ScanRecord {
	all, err := st.ListRepos()
	if err != nil {
		return nil
	}
	for _, r := range all {
		if strings.EqualFold(r.RepoName, query) || strings.EqualFold(r.Slug, query) {
			return &r
		}
	}
	return nil
}

// PackageRepoArchive packages a local repository directory into a .tar.gz archive.
func PackageRepoArchive(repoPath string, opts DeployRepoOptions) ([]byte, int, error) {
	stat, err := os.Stat(repoPath)
	if err != nil || !stat.IsDir() {
		return nil, 0, apperror.NewNotFound("directory", "E9073", "repository directory "+repoPath+" not found")
	}
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	count, walkErr := walkAndArchiveRepo(tw, repoPath, opts)
	_ = tw.Close()
	_ = gw.Close()
	if walkErr != nil {
		return nil, 0, walkErr
	}
	return buf.Bytes(), count, nil
}

func walkAndArchiveRepo(tw *tar.Writer, repoPath string, opts DeployRepoOptions) (int, error) {
	fileCount := 0
	err := filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil || path == repoPath {
			return nil
		}
		rel, relErr := filepath.Rel(repoPath, path)
		if relErr != nil {
			return nil
		}
		if isWalkEntrySkipped(info, rel, opts) {
			return handleSkipDir(info)
		}
		if info.IsDir() {
			return nil
		}
		if writeTarFileEntry(tw, filepath.ToSlash(rel), path) {
			fileCount++
		}
		return nil
	})
	return fileCount, err
}

func handleSkipDir(info os.FileInfo) error {
	if info.IsDir() {
		return filepath.SkipDir
	}
	return nil
}

func isWalkEntrySkipped(info os.FileInfo, relPath string, opts DeployRepoOptions) bool {
	if opts.Clean && isBuildArtifactDir(info.Name()) {
		return true
	}
	if opts.Clean || opts.ExcludeGitObjects {
		return isGitObjectsPath(relPath)
	}
	return false
}

func isBuildArtifactDir(name string) bool {
	switch strings.ToLower(name) {
	case "node_modules", "target", ".venv", "dist":
		return true
	default:
		return false
	}
}

func isGitObjectsPath(relPath string) bool {
	norm := filepath.ToSlash(relPath)
	return norm == ".git/objects" || strings.HasPrefix(norm, ".git/objects/")
}

// RunNodesDeployRepo deploys a repository across candidate fleet nodes.
func RunNodesDeployRepo(args []string) error {
	opts, err := ParseDeployRepoOptions(args)
	if err != nil {
		return err
	}
	repoRecord, localRepoPath, errRepo := resolveLocalRepo(opts.Repo, opts.FromLocal)
	if errRepo != nil {
		return errRepo
	}
	conns, errConns := cmdssh.FetchAllSSHConnections()
	if errConns != nil || len(conns) == 0 {
		return apperror.NewSimple("no registered fleet SSH nodes found", "E9074")
	}
	targets := FilterFleetNodes(conns, opts.NodeFilterOptions)
	if len(targets) == 0 {
		return apperror.NewSimple("no candidate fleet nodes found matching filter criteria", "E9075")
	}
	if opts.DryRun {
		return executeDryRunDeployRepo(targets, opts, repoRecord, localRepoPath)
	}
	tarData, fileCount, errPkg := PackageRepoArchive(localRepoPath, opts)
	if errPkg != nil {
		return errPkg
	}
	return executeFleetDeployRepo(targets, opts, repoRecord, localRepoPath, tarData, fileCount)
}

// RunNodesDeployRepos orchestrates repository deployments across candidate fleet nodes.
func RunNodesDeployRepos(args []string) error {
	return RunNodesDeployRepo(args)
}

func executeDryRunDeployRepo(targets []db.SSHConnection, opts DeployRepoOptions, repoRecord *model.ScanRecord, localRepoPath string) error {
	var results []DeployRepoResult
	for _, c := range targets {
		destDir := resolveRemoteDestDir(c, repoRecord.Slug, opts.Dest)
		ides := []string{"vscode", "cursor", "antigravity"}
		results = append(results, DeployRepoResult{
			NodeAlias: c.Alias,
			Host:      c.IPAddress,
			OS:        c.OS,
			DestDir:   destDir,
			Status:    "DRY-RUN",
			IDEs:      ides,
		})
	}
	return outputRepoResults(results, opts, repoRecord.Slug, 0, len(results))
}

func executeFleetDeployRepo(targets []db.SSHConnection, opts DeployRepoOptions, repoRecord *model.ScanRecord, localRepoPath string, tarData []byte, fileCount int) error {
	results := deployRepoToTargets(targets, opts, repoRecord, localRepoPath, tarData)
	return outputRepoResults(results, opts, repoRecord.Slug, fileCount, len(tarData))
}

func deployRepoToTargets(targets []db.SSHConnection, opts DeployRepoOptions, repoRecord *model.ScanRecord, localRepoPath string, tarData []byte) []DeployRepoResult {
	var results []DeployRepoResult
	var wg sync.WaitGroup
	var mu sync.Mutex
	for _, c := range targets {
		wg.Add(1)
		go func(conn db.SSHConnection) {
			defer wg.Done()
			res := executeDeployRepoWorker(conn, opts, repoRecord, localRepoPath, tarData)
			mu.Lock()
			results = append(results, res)
			mu.Unlock()
		}(c)
	}
	wg.Wait()
	return filterOpenOnlyRepoResults(results, opts.OpenOnly)
}

func filterOpenOnlyRepoResults(results []DeployRepoResult, openOnly bool) []DeployRepoResult {
	if !openOnly {
		return results
	}
	var out []DeployRepoResult
	for _, r := range results {
		if r.Status != "OFFLINE" {
			out = append(out, r)
		}
	}
	return out
}

func executeDeployRepoWorker(conn db.SSHConnection, opts DeployRepoOptions, repoRecord *model.ScanRecord, localRepoPath string, tarData []byte) DeployRepoResult {
	start := time.Now()
	destDir := resolveRemoteDestDir(conn, repoRecord.Slug, opts.Dest)
	res := DeployRepoResult{
		NodeAlias:   conn.Alias,
		Host:        conn.IPAddress,
		OS:          conn.OS,
		DestDir:     destDir,
		ArchiveSize: int64(len(tarData)),
	}
	if !probeAGMNodeLiveness(context.Background(), conn, 1500*time.Millisecond) {
		return buildOfflineRepoResult(res, start)
	}
	return runRemoteRepoDeploy(conn, opts, repoRecord, localRepoPath, destDir, tarData, res, start)
}

func buildOfflineRepoResult(res DeployRepoResult, start time.Time) DeployRepoResult {
	res.Status = "OFFLINE"
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	res.Error = "connection timed out or unreachable"
	return res
}

func runRemoteRepoDeploy(conn db.SSHConnection, opts DeployRepoOptions, repoRecord *model.ScanRecord, localRepoPath, destDir string, tarData []byte, res DeployRepoResult, start time.Time) DeployRepoResult {
	client, errConnect := cmdssh.ConnectSSHClientWithErr(conn)
	if errConnect != nil {
		return handleRepoConnectFailure(res, errConnect, start)
	}
	defer client.Close()

	stagingPath := resolveRepoStagingPath(conn, repoRecord.Slug)
	if errStream := cmdssh.StreamFileToRemote(client, stagingPath, tarData, conn.OS); errStream != nil {
		return buildFailedRepoResult(res, "stream failed: "+errStream.Error(), start)
	}

	extractCmd, shell := resolveRepoExtractCommand(conn, stagingPath, destDir)
	if _, errExtract := secrets.RunCommand(client, extractCmd, shell); errExtract != nil {
		return buildFailedRepoResult(res, "extract failed: "+errExtract.Error(), start)
	}

	runRemoteRescan(client, conn)
	res.IDEs = runRemoteIDEsHook(client, conn, repoRecord.RepoName, destDir, opts)
	res.Conversations = runRemoteConvsHook(client, conn, localRepoPath, opts)
	res.Status = "SUCCESS"
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	return res
}

func buildFailedRepoResult(res DeployRepoResult, errMsg string, start time.Time) DeployRepoResult {
	res.Status = "FAILED"
	res.Error = errMsg
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	return res
}

func handleRepoConnectFailure(res DeployRepoResult, errConnect error, start time.Time) DeployRepoResult {
	res.Latency = time.Since(start)
	res.LatencyMs = res.Latency.Milliseconds()
	if isOfflineError(errConnect.Error()) {
		res.Status = "OFFLINE"
		res.Error = errConnect.Error()
		return res
	}
	res.Status = "FAILED"
	res.Error = errConnect.Error()
	return res
}

func runRemoteRescan(client *ssh.Client, conn db.SSHConnection) {
	shell := "sh"
	if isWindowsNode(conn) {
		shell = "cmd"
	}
	_, _ = secrets.RunCommand(client, "gitmap rescan", shell)
}

func runRemoteIDEsHook(client *ssh.Client, conn db.SSHConnection, repoName, destDir string, opts DeployRepoOptions) []string {
	if !opts.WithIDEs && !opts.WithPinned {
		return nil
	}
	ides, err := RegisterRemoteIDEs(client, conn, repoName, destDir, opts.WithPinned)
	if err != nil {
		return nil
	}
	return ides
}

func runRemoteConvsHook(client *ssh.Client, conn db.SSHConnection, localRepoPath string, opts DeployRepoOptions) int {
	if !opts.WithConversations {
		return 0
	}
	count, err := DeployConversationsForRepo(client, conn, localRepoPath)
	if err != nil {
		return 0
	}
	return count
}

func resolveRemoteDestDir(conn db.SSHConnection, slug, customDest string) string {
	if customDest != "" {
		return customDest
	}
	if isWindowsNode(conn) {
		return `D:\work\` + slug
	}
	return "~/work/" + slug
}

func resolveRepoStagingPath(conn db.SSHConnection, slug string) string {
	if isWindowsNode(conn) {
		return fmt.Sprintf(`C:\Windows\Temp\repo_deploy_%s.tar.gz`, slug)
	}
	return fmt.Sprintf(`/tmp/repo_deploy_%s.tar.gz`, slug)
}

func resolveRepoExtractCommand(conn db.SSHConnection, stagingPath, destDir string) (string, string) {
	if isWindowsNode(conn) {
		cmd := fmt.Sprintf(`powershell.exe -NoProfile -Command "$dest = '%s'; if (-not (Test-Path $dest)) { New-Item -ItemType Directory -Force -Path $dest | Out-Null }; tar.exe -xzf '%s' -C $dest; Remove-Item '%s' -Force -ErrorAction SilentlyContinue"`, destDir, stagingPath, stagingPath)
		return cmd, "ps"
	}
	cmd := fmt.Sprintf(`sh -c "mkdir -p \"%s\" && tar -xzf \"%s\" -C \"%s\" && rm -f \"%s\""`, destDir, stagingPath, destDir, stagingPath)
	return cmd, "sh"
}

func outputRepoResults(results []DeployRepoResult, opts DeployRepoOptions, slug string, fileCount, sizeBytes int) error {
	if opts.IsJSON {
		return emitDeployRepoJSON(results)
	}
	renderDeployRepoTable(results)
	printDeployRepoSummary(results, slug, opts.DryRun, fileCount, sizeBytes)
	return nil
}

func emitDeployRepoJSON(results []DeployRepoResult) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(results)
}

func renderDeployRepoTable(results []DeployRepoResult) {
	fmt.Println()
	fmt.Printf("  %s%-18s %-16s %-10s %-24s %-12s %-10s%s\n",
		constants.ColorCyan, "NODE ALIAS", "HOST", "OS", "DEST DIRECTORY", "STATUS", "LATENCY", constants.ColorReset)
	fmt.Println("  " + strings.Repeat("-", 96))
	for _, r := range results {
		statusStr := formatStatusBadge(r.Status)
		osStr := r.OS
		if osStr == "" {
			osStr = "linux"
		}
		fmt.Printf("  %-18s %-16s %-10s %-24s %-12s %-10v\n",
			r.NodeAlias, r.Host, osStr, truncateMiddle(r.DestDir, 24), statusStr, r.Latency.Round(time.Millisecond))
	}
	fmt.Println()
}

func truncateMiddle(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	half := (maxLen - 3) / 2
	return s[:half] + "..." + s[len(s)-half:]
}

func printDeployRepoSummary(results []DeployRepoResult, slug string, isDryRun bool, fileCount, sizeBytes int) {
	successCount := 0
	offlineCount := 0
	for _, r := range results {
		if r.Status == "SUCCESS" || r.Status == "DRY-RUN" {
			successCount++
		}
		if r.Status == "OFFLINE" {
			offlineCount++
		}
	}
	prefix := constants.ColorGreen + "✓" + constants.ColorReset
	mode := "Deployed"
	if isDryRun {
		mode = "Simulated deployment of"
	}
	fmt.Printf("  %s %s repository %q across %d/%d node(s) (%d offline)\n\n",
		prefix, mode, slug, successCount, len(results), offlineCount)
}
