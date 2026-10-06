package cmdpurge

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// TargetType defines the resolution mode for history purging.
type TargetType string

const (
	TargetTypeFolder TargetType = "folder"
	TargetTypeFile   TargetType = "file"
	TargetTypeCommit TargetType = "commit"
)

// PreflightScanOptions encapsulates inputs for scanning the commit graph.
type PreflightScanOptions struct {
	RepoDir       string     `json:"repoDir"`
	TargetType    TargetType `json:"targetType"`
	TargetPath    string     `json:"targetPath"` // Relative file or folder path
	TargetCommits []string   `json:"targetCommits"`
	BranchFilter  string     `json:"branchFilter"`
	IsAllowDirty  bool       `json:"isAllowDirty"`
}

// PurgeOptions encapsulates options for the history purge preflight check.
type PurgeOptions struct {
	TargetType    TargetType `json:"targetType"`
	TargetPath    string     `json:"targetPath"`
	TargetCommits []string   `json:"targetCommits"`
	BranchFilter  string     `json:"branchFilter"`
	IsAutoConfirm bool       `json:"isAutoConfirm"`
	IsDryRun      bool       `json:"isDryRun"`
	IsAllowDirty  bool       `json:"isAllowDirty"`
}

// AffectedCommitNode describes a single commit touched by the target pattern.
type AffectedCommitNode struct {
	CommitSha       string   `json:"commitSha"`
	ShortSha        string   `json:"shortSha"`
	AuthorName      string   `json:"authorName"`
	CommitTimestamp int64    `json:"commitTimestamp"`
	CommitMessage   string   `json:"commitMessage"`
	ParentShas      []string `json:"parentShas"`
	IsAffected      bool     `json:"isAffected"`
	MatchedFiles    []string `json:"matchedFiles"`
}

// PreflightReport summarizes the impact and blast radius of the proposed purge.
type PreflightReport struct {
	RepoSlug            string               `json:"repoSlug"`
	TargetType          TargetType           `json:"targetType"`
	TargetPath          string               `json:"targetPath"`
	TotalCommitsScanned int                  `json:"totalCommitsScanned"`
	AffectedCommits     []AffectedCommitNode `json:"affectedCommits"`
	AffectedBranches    []string             `json:"affectedBranches"`
	AffectedTags        []string             `json:"affectedTags"`
	DistinctFilesCount  int                  `json:"distinctFilesCount"`
	TotalEstimatedBytes int64                `json:"totalEstimatedBytes"`
	HasPushedCommits    bool                 `json:"hasPushedCommits"`
	IsWorkingTreeDirty  bool                 `json:"isWorkingTreeDirty"`
}

// RunPurgePreflight scans git commits and compiles a blast-radius report.
func RunPurgePreflight(repoDir string, opts PurgeOptions) (*PreflightReport, error) {
	scanOpts := PreflightScanOptions{
		RepoDir:       repoDir,
		TargetType:    opts.TargetType,
		TargetPath:    opts.TargetPath,
		TargetCommits: opts.TargetCommits,
		BranchFilter:  opts.BranchFilter,
		IsAllowDirty:  opts.IsAllowDirty,
	}

	return RunPreflightScan(scanOpts)
}

// RunPreflightScan executes the commit graph traversal and impact calculation.
func RunPreflightScan(opts PreflightScanOptions) (*PreflightReport, error) {
	isDirty, err := CheckWorkingTreeDirty(opts.RepoDir)
	if err != nil {
		return nil, err
	}
	if isDirty && !opts.IsAllowDirty {
		return nil, apperror.NewValidation("RunPreflightScan", "ERR_PURGE_DIRTY_TREE", "working tree is dirty")
	}

	normOpts := normalizeScanOptions(opts)

	return executePreflightInspection(normOpts, isDirty)
}

func normalizeScanOptions(opts PreflightScanOptions) PreflightScanOptions {
	res := opts
	res.TargetPath = filepath.ToSlash(filepath.Clean(opts.TargetPath))
	if res.TargetType == "" {
		res.TargetType = detectTargetType(res.TargetPath, res.TargetCommits)
	}

	return res
}

func detectTargetType(path string, commits []string) TargetType {
	if len(commits) > 0 {
		return TargetTypeCommit
	}
	if strings.HasSuffix(path, "/") {
		return TargetTypeFolder
	}

	return TargetTypeFile
}

func executePreflightInspection(opts PreflightScanOptions, isDirty bool) (*PreflightReport, error) {
	rawLog, err := fetchGitLogRaw(opts.RepoDir, opts.BranchFilter)
	if err != nil {
		return nil, err
	}

	nodes, total := parseCommitNodes(rawLog, opts)
	affNodes, affShas := filterAffectedNodes(nodes)
	branches := collectBranchesForCommits(opts.RepoDir, affShas)
	tags := collectTagsForCommits(opts.RepoDir, affShas)
	hasPushed, _ := CheckCommitsPushed(opts.RepoDir, affShas)
	distinctFiles, totalBytes := calculateBlastRadius(opts.RepoDir, affNodes)

	return &PreflightReport{
		RepoSlug:            resolveRepoSlug(opts.RepoDir),
		TargetType:          opts.TargetType,
		TargetPath:          opts.TargetPath,
		TotalCommitsScanned: total,
		AffectedCommits:     affNodes,
		AffectedBranches:    branches,
		AffectedTags:        tags,
		DistinctFilesCount:  len(distinctFiles),
		TotalEstimatedBytes: totalBytes,
		HasPushedCommits:    hasPushed,
		IsWorkingTreeDirty:  isDirty,
	}, nil
}

// CheckWorkingTreeDirty returns true if uncommitted changes exist.
func CheckWorkingTreeDirty(repoDir string) (bool, error) {
	cmd := exec.Command("git", "-C", repoDir, "status", "--porcelain")
	out, err := cmd.Output()
	if err != nil {
		return false, apperror.WrapSimple(err, "git status check")
	}

	return strings.TrimSpace(string(out)) != "", nil
}

// CheckCommitsPushed returns true if any affected commit exists on remote tracking branches.
func CheckCommitsPushed(repoDir string, commits []string) (bool, error) {
	if len(commits) == 0 {
		return false, nil
	}
	for _, sha := range commits {
		cmd := exec.Command("git", "-C", repoDir, "branch", "-r", "--contains", sha)
		out, err := cmd.Output()
		if err == nil && strings.TrimSpace(string(out)) != "" {
			return true, nil
		}
	}

	return false, nil
}

// ResolveAffectedBranches returns local branches containing the commit.
func ResolveAffectedBranches(repoDir string, commitSha string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoDir, "branch", "--contains", commitSha)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "resolve affected branches")
	}

	return parseBranchNames(string(out)), nil
}

func parseBranchNames(raw string) []string {
	var list []string
	for _, line := range strings.Split(raw, "\n") {
		clean := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), "*"))
		if clean != "" && !strings.Contains(clean, "->") {
			list = append(list, clean)
		}
	}

	return list
}

// ResolveAffectedTags returns tags pointing to or descending from the commit.
func ResolveAffectedTags(repoDir string, commitSha string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoDir, "tag", "--contains", commitSha)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "resolve affected tags")
	}

	return parseTagNames(string(out)), nil
}

func parseTagNames(raw string) []string {
	var list []string
	for _, line := range strings.Split(raw, "\n") {
		clean := strings.TrimSpace(line)
		if clean != "" {
			list = append(list, clean)
		}
	}

	return list
}

// PromptConfirmation asks the user for explicit confirmation unless bypassed.
func PromptConfirmation(isAutoYes bool, in io.Reader, out io.Writer) (bool, error) {
	if isAutoYes {
		_, _ = fmt.Fprintln(out, "[Auto-confirmed via -y flag]")
		return true, nil
	}

	_, _ = fmt.Fprint(out, "Proceed with history purge? [y/N]: ")
	scanner := bufio.NewScanner(in)
	if !scanner.Scan() {
		return false, nil
	}

	resp := strings.ToLower(strings.TrimSpace(scanner.Text()))

	return resp == "y" || resp == "yes", nil
}

func fetchGitLogRaw(repoDir, branchFilter string) (string, error) {
	args := []string{"-C", repoDir, "log", "--name-only", "--format=COMMIT_DIV%x00%H%x00%h%x00%an%x00%ct%x00%s%x00%P"}
	if branchFilter != "" {
		args = append(args, branchFilter)
	} else {
		args = append(args, "--all")
	}

	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "fetch git log")
	}

	return string(out), nil
}

func parseCommitNodes(rawLog string, opts PreflightScanOptions) ([]AffectedCommitNode, int) {
	chunks := strings.Split(rawLog, "COMMIT_DIV\x00")
	var nodes []AffectedCommitNode
	total := 0
	for _, chunk := range chunks {
		trimmed := strings.TrimSpace(chunk)
		if trimmed == "" {
			continue
		}
		total++
		node := buildCommitNode(trimmed, opts)
		nodes = append(nodes, node)
	}

	return nodes, total
}

func buildCommitNode(chunk string, opts PreflightScanOptions) AffectedCommitNode {
	lines := strings.Split(chunk, "\n")
	headerParts := strings.Split(lines[0], "\x00")
	files := extractChangedFiles(lines[1:])
	matched := matchTargetFiles(files, headerParts[0], opts)
	ts, _ := strconv.ParseInt(getHeaderPart(headerParts, 3), 10, 64)
	parents := strings.Fields(getHeaderPart(headerParts, 5))

	return AffectedCommitNode{
		CommitSha:       getHeaderPart(headerParts, 0),
		ShortSha:        getHeaderPart(headerParts, 1),
		AuthorName:      getHeaderPart(headerParts, 2),
		CommitTimestamp: ts,
		CommitMessage:   getHeaderPart(headerParts, 4),
		ParentShas:      parents,
		IsAffected:      len(matched) > 0,
		MatchedFiles:    matched,
	}
}

func getHeaderPart(parts []string, idx int) string {
	if idx < len(parts) {
		return parts[idx]
	}

	return ""
}

func extractChangedFiles(lines []string) []string {
	var files []string
	for _, l := range lines {
		clean := strings.TrimSpace(l)
		if clean != "" {
			files = append(files, clean)
		}
	}

	return files
}

func matchTargetFiles(files []string, commitSha string, opts PreflightScanOptions) []string {
	if opts.TargetType == TargetTypeCommit {
		return matchCommitShaTarget(commitSha, opts.TargetCommits)
	}

	var matched []string
	for _, f := range files {
		cleanFile := filepath.ToSlash(filepath.Clean(f))
		if isPathMatch(cleanFile, opts.TargetPath, opts.TargetType) {
			matched = append(matched, cleanFile)
		}
	}

	return matched
}

func matchCommitShaTarget(sha string, targetShas []string) []string {
	for _, target := range targetShas {
		if strings.HasPrefix(sha, target) || strings.HasPrefix(target, sha) {
			return []string{"[commit-targeted]"}
		}
	}

	return nil
}

func isPathMatch(filePath, targetPath string, targetType TargetType) bool {
	if targetType == TargetTypeFolder {
		folderPrefix := strings.TrimSuffix(targetPath, "/") + "/"

		return strings.HasPrefix(filePath, folderPrefix) || filePath == targetPath
	}

	return filePath == targetPath
}

func filterAffectedNodes(nodes []AffectedCommitNode) ([]AffectedCommitNode, []string) {
	var affNodes []AffectedCommitNode
	var affShas []string
	for _, n := range nodes {
		if n.IsAffected {
			affNodes = append(affNodes, n)
			affShas = append(affShas, n.CommitSha)
		}
	}

	return affNodes, affShas
}

func collectBranchesForCommits(repoDir string, shas []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, sha := range shas {
		branches, _ := ResolveAffectedBranches(repoDir, sha)
		for _, b := range branches {
			if !seen[b] {
				seen[b] = true
				result = append(result, b)
			}
		}
	}

	return result
}

func collectTagsForCommits(repoDir string, shas []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, sha := range shas {
		tags, _ := ResolveAffectedTags(repoDir, sha)
		for _, t := range tags {
			if !seen[t] {
				seen[t] = true
				result = append(result, t)
			}
		}
	}

	return result
}

func calculateBlastRadius(repoDir string, nodes []AffectedCommitNode) ([]string, int64) {
	distinctMap := make(map[string]bool)
	var totalBytes int64
	for _, n := range nodes {
		for _, f := range n.MatchedFiles {
			if f == "[commit-targeted]" {
				continue
			}
			if !distinctMap[f] {
				distinctMap[f] = true
				totalBytes += estimateBlobSize(repoDir, n.CommitSha, f)
			}
		}
	}

	var files []string
	for f := range distinctMap {
		files = append(files, f)
	}

	return files, totalBytes
}

func estimateBlobSize(repoDir, commitSha, relPath string) int64 {
	cmd := exec.Command("git", "-C", repoDir, "cat-file", "-s", commitSha+":"+relPath)
	out, err := cmd.Output()
	if err != nil {
		return 0
	}
	bytes, _ := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)

	return bytes
}

func resolveRepoSlug(repoDir string) string {
	cmd := exec.Command("git", "-C", repoDir, "remote", "get-url", "origin")
	out, err := cmd.Output()
	if err == nil && strings.TrimSpace(string(out)) != "" {
		return gitutil.CanonicalRepoID(strings.TrimSpace(string(out)))
	}

	return filepath.Base(repoDir)
}
