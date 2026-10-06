package cmdpurge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/tempdir"
)

// PurgeExtendedOptions adds CLI-specific flags to PurgeOptions.
type PurgeExtendedOptions struct {
	PurgeOptions
	RepoDir   string `json:"repoDir"`
	IsRelease bool   `json:"isRelease"`
	IsJson    bool   `json:"isJson"`
	IsVerbose bool   `json:"isVerbose"`
	NoBackup  bool   `json:"noBackup"`
}

// PurgeExecutionReport describes the result of a history rewrite operation.
type PurgeExecutionReport struct {
	OperationId         int64      `json:"operationId"`
	RepoSlug            string     `json:"repoSlug"`
	Branch              string     `json:"branch"`
	TargetType          TargetType `json:"targetType"`
	TargetPath          string     `json:"targetPath"`
	OriginalHeadSha     string     `json:"originalHeadSha"`
	RewrittenHeadSha    string     `json:"rewrittenHeadSha"`
	BackupRef           string     `json:"backupRef"`
	TempBackupDir       string     `json:"tempBackupDir"`
	ScannedCommits      int        `json:"scannedCommits"`
	AffectedCommits     int        `json:"affectedCommits"`
	PurgedBlobsCount    int        `json:"purgedBlobsCount"`
	PurgedBlobsBytes    int64      `json:"purgedBlobsBytes"`
	ReleaseAssetsPurged int        `json:"releaseAssetsPurged"`
	ReleaseNotesUpdated int        `json:"releaseNotesUpdated"`
	NonWarrantyAdvisory string     `json:"nonWarrantyAdvisory"`
	IsDryRun            bool       `json:"isDryRun"`
	IsSuccess           bool       `json:"isSuccess"`
}

type rewriteContext struct {
	repoDir     string
	branch      string
	origHead    string
	opId        int64
	db          *store.PurgeHistoryDB
	commitRemap map[string]string
	commitMaps  []store.HistoryPurgeCommitMap
	fileRecords []store.HistoryPurgeFile
}

type rawCommitInfo struct {
	treeSha   string
	parents   []string
	authorEnv []string
	commitMsg string
}

// ExecutePurgeEngine runs the history purge pipeline using PurgeOptions.
func ExecutePurgeEngine(opts PurgeOptions) (*PurgeExecutionReport, error) {
	extOpts := PurgeExtendedOptions{
		PurgeOptions: opts,
		RepoDir:      ".",
	}

	return ExecutePurgeEngineExtended(extOpts)
}

// ExecutePurgeEngineExtended executes the full pre-flight, backup, tree rewrite, and GC pipeline.
func ExecutePurgeEngineExtended(opts PurgeExtendedOptions) (*PurgeExecutionReport, error) {
	repoDir := resolveRepoDir(opts.RepoDir)
	report, err := runPreflightGate(repoDir, opts)
	if err != nil || report == nil {
		return nil, err
	}
	if opts.IsDryRun || len(report.AffectedCommits) == 0 {
		return buildDryRunReport(opts, report), nil
	}

	db, err := store.OpenPurgeHistoryDB("")
	if err != nil {
		return nil, apperror.WrapSimple(err, "open splitdb for purge")
	}
	defer db.Close()

	return runRewritePipeline(db, repoDir, opts, report)
}

func runPreflightGate(repoDir string, opts PurgeExtendedOptions) (*PreflightReport, error) {
	report, err := RunPurgePreflight(repoDir, opts.PurgeOptions)
	if err != nil {
		return nil, err
	}
	if len(report.AffectedCommits) == 0 {
		if !opts.IsJson {
			fmt.Println("No matching target found in repository history. Nothing to purge.")
		}
		return report, nil
	}
	if !opts.IsAutoConfirm && !opts.IsJson && !opts.IsDryRun {
		if isConfirmed := confirmPurgeInteractively(report); !isConfirmed {
			return nil, nil
		}
	}

	return report, nil
}

func confirmPurgeInteractively(report *PreflightReport) bool {
	RenderPreflightBox(report, os.Stdout)
	fmt.Print("Proceed with history purge? [y/N]: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	clean := strings.ToLower(strings.TrimSpace(text))
	if clean != "y" && clean != "yes" {
		fmt.Println("[!] Purge cancelled by user. Repository commit graph remains unaltered.")
		return false
	}

	return true
}

func buildDryRunReport(opts PurgeExtendedOptions, report *PreflightReport) *PurgeExecutionReport {
	if opts.IsDryRun && !opts.IsJson {
		RenderPreflightBox(report, os.Stdout)
		fmt.Println("[Dry Run] History purge simulated. No repository objects modified.")
	}

	return &PurgeExecutionReport{
		RepoSlug:            report.RepoSlug,
		TargetType:          report.TargetType,
		TargetPath:          report.TargetPath,
		ScannedCommits:      report.TotalCommitsScanned,
		AffectedCommits:     len(report.AffectedCommits),
		PurgedBlobsBytes:    report.TotalEstimatedBytes,
		NonWarrantyAdvisory: makeAdvisoryString(0),
		IsDryRun:            opts.IsDryRun,
		IsSuccess:           true,
	}
}

func runRewritePipeline(db *store.PurgeHistoryDB, repoDir string, opts PurgeExtendedOptions, preReport *PreflightReport) (*PurgeExecutionReport, error) {
	targets := extractBlobTargets(repoDir, preReport.AffectedCommits)
	vaultPath, err := stageBackupVault(repoDir, preReport.RepoSlug, opts.NoBackup, targets)
	if err != nil {
		return nil, err
	}

	opId, err := recordPurgeStart(db, preReport, opts, vaultPath, len(targets))
	if err != nil {
		return nil, err
	}

	rCtx, err := initRewriteContext(db, repoDir, opId)
	if err != nil {
		return nil, err
	}

	newHead, err := rewriteCommitGraph(rCtx, opts.PurgeOptions)
	if err != nil {
		return nil, err
	}

	return finalizePurgeExecution(rCtx, opts, preReport, newHead, vaultPath)
}

func extractBlobTargets(repoDir string, nodes []AffectedCommitNode) []tempdir.PurgeBlobTarget {
	var targets []tempdir.PurgeBlobTarget
	for _, node := range nodes {
		for _, f := range node.MatchedFiles {
			data, _ := readBlobData(repoDir, node.CommitSha, f)
			targets = append(targets, tempdir.PurgeBlobTarget{
				CommitSha:    node.CommitSha,
				RelativePath: f,
				Data:         data,
				FileMode:     0644,
			})
		}
	}

	return targets
}

func readBlobData(repoDir, commitSha, relPath string) ([]byte, error) {
	cmd := exec.Command("git", "show", fmt.Sprintf("%s:%s", commitSha, relPath))
	cmd.Dir = repoDir

	return cmd.Output()
}

func stageBackupVault(repoDir, slug string, isNoBackup bool, targets []tempdir.PurgeBlobTarget) (string, error) {
	if isNoBackup || len(targets) == 0 {
		return "", nil
	}
	opSubdir := fmt.Sprintf("%d", time.Now().UnixNano())

	return tempdir.StageBlobsToTempBackup(repoDir, slug, opSubdir, targets)
}

func recordPurgeStart(db *store.PurgeHistoryDB, rep *PreflightReport, opts PurgeExtendedOptions, vaultPath string, count int) (int64, error) {
	commitsJson, _ := json.Marshal(opts.TargetCommits)
	op := &store.HistoryPurgeOperation{
		RepoSlug:             rep.RepoSlug,
		TargetType:           string(opts.TargetType),
		TargetPath:           opts.TargetPath,
		CommitListJson:       string(commitsJson),
		TotalCommitsScanned:  rep.TotalCommitsScanned,
		AffectedCommitsCount: len(rep.AffectedCommits),
		BackedUpFilesCount:   count,
		BackupVaultPath:      vaultPath,
		IsDryRun:             opts.IsDryRun,
		HasPushed:            rep.HasPushedCommits,
		CreatedAt:            time.Now().Unix(),
	}

	return db.InsertHistoryPurgeOperation(op)
}

func initRewriteContext(db *store.PurgeHistoryDB, repoDir string, opId int64) (*rewriteContext, error) {
	branch, err := queryCurrentBranch(repoDir)
	if err != nil {
		return nil, err
	}
	origHead, hasHead := verifyGitRef(repoDir, "HEAD")
	if !hasHead {
		return nil, apperror.NewSimple("ERR_PURGE_NO_HEAD", "failed to resolve HEAD commit")
	}

	backupRef := fmt.Sprintf("refs/gitmap-backup/%d/%s", opId, branch)
	cmd := exec.Command("git", "update-ref", backupRef, origHead)
	cmd.Dir = repoDir
	if err := cmd.Run(); err != nil {
		return nil, apperror.WrapSimple(err, "create safety backup reference")
	}

	return &rewriteContext{
		repoDir:     repoDir,
		branch:      branch,
		origHead:    origHead,
		opId:        opId,
		db:          db,
		commitRemap: make(map[string]string),
	}, nil
}

func rewriteCommitGraph(ctx *rewriteContext, opts PurgeOptions) (string, error) {
	commits, err := fetchTopologicalCommits(ctx.repoDir)
	if err != nil {
		return "", err
	}

	for _, sha := range commits {
		newSha, err := rewriteSingleCommit(ctx, sha, opts)
		if err != nil {
			return "", err
		}
		ctx.commitRemap[sha] = newSha
	}

	newHead := ctx.commitRemap[ctx.origHead]
	if err := updateBranchAndReset(ctx.repoDir, ctx.branch, newHead); err != nil {
		return "", err
	}

	return newHead, nil
}

func fetchTopologicalCommits(repoDir string) ([]string, error) {
	cmd := exec.Command("git", "rev-list", "--reverse", "--topo-order", "HEAD")
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "fetch topological commits")
	}

	return strings.Fields(string(out)), nil
}

func rewriteSingleCommit(ctx *rewriteContext, sha string, opts PurgeOptions) (string, error) {
	raw, err := inspectCommitRaw(ctx.repoDir, sha)
	if err != nil {
		return "", err
	}

	newTreeSha, err := filterAndRebuildTree(ctx.repoDir, raw.treeSha, opts)
	if err != nil {
		return "", err
	}

	remappedParents := remapParentShas(raw.parents, ctx.commitRemap)
	newCommitSha, err := createRewrittenCommit(ctx.repoDir, newTreeSha, remappedParents, raw)
	if err != nil {
		return "", err
	}

	recordCommitMapEntry(ctx, sha, newCommitSha, raw)

	return newCommitSha, nil
}

func inspectCommitRaw(repoDir, sha string) (*rawCommitInfo, error) {
	cmd := exec.Command("git", "log", "-1", "--format=raw", sha)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "read raw commit info")
	}

	return parseRawCommitOutput(string(out)), nil
}

func parseRawCommitOutput(raw string) *rawCommitInfo {
	info := &rawCommitInfo{}
	scanner := bufio.NewScanner(strings.NewReader(raw))
	inMsg := false
	var msgLines []string

	for scanner.Scan() {
		line := scanner.Text()
		if inMsg {
			msgLines = append(msgLines, strings.TrimPrefix(line, "    "))
			continue
		}
		if line == "" {
			inMsg = true
			continue
		}
		parseCommitHeaderLine(line, info)
	}
	info.commitMsg = strings.Join(msgLines, "\n")

	return info
}

func parseCommitHeaderLine(line string, info *rawCommitInfo) {
	parts := strings.SplitN(line, " ", 2)
	if len(parts) < 2 {
		return
	}
	switch parts[0] {
	case "tree":
		info.treeSha = parts[1]
	case "parent":
		info.parents = append(info.parents, parts[1])
	case "author":
		info.authorEnv = append(info.authorEnv, parseIdentityEnv("GIT_AUTHOR", parts[1])...)
	case "committer":
		info.authorEnv = append(info.authorEnv, parseIdentityEnv("GIT_COMMITTER", parts[1])...)
	}
}

func parseIdentityEnv(prefix, identity string) []string {
	parts := strings.Fields(identity)
	if len(parts) < 3 {
		return nil
	}
	tz := parts[len(parts)-1]
	timestamp := parts[len(parts)-2]
	nameAndEmail := strings.Join(parts[:len(parts)-2], " ")
	email := ""
	name := nameAndEmail
	if start := strings.Index(nameAndEmail, "<"); start >= 0 {
		if end := strings.Index(nameAndEmail, ">"); end > start {
			email = nameAndEmail[start+1 : end]
			name = strings.TrimSpace(nameAndEmail[:start])
		}
	}

	return []string{
		fmt.Sprintf("%s_NAME=%s", prefix, name),
		fmt.Sprintf("%s_EMAIL=%s", prefix, email),
		fmt.Sprintf("%s_DATE=%s %s", prefix, timestamp, tz),
	}
}

func filterAndRebuildTree(repoDir, treeSha string, opts PurgeOptions) (string, error) {
	entries, err := fetchTreeEntries(repoDir, treeSha)
	if err != nil {
		return "", err
	}

	filtered := filterTreeEntries(entries, opts)
	if len(filtered) == 0 {
		return makeEmptyTree(repoDir)
	}

	return buildMktree(repoDir, filtered)
}

func fetchTreeEntries(repoDir, treeSha string) ([]string, error) {
	cmd := exec.Command("git", "ls-tree", "-r", "-z", treeSha)
	cmd.Dir = repoDir
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "ls-tree failed")
	}

	return strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00"), nil
}

func filterTreeEntries(entries []string, opts PurgeOptions) []string {
	var kept []string
	cleanTarget := filepath.ToSlash(filepath.Clean(opts.TargetPath))

	for _, entry := range entries {
		if len(entry) == 0 {
			continue
		}
		path := extractPathFromTreeEntry(entry)
		if isTreePathExcluded(path, cleanTarget, opts.TargetType) {
			continue
		}
		kept = append(kept, entry)
	}

	return kept
}

func extractPathFromTreeEntry(entry string) string {
	parts := strings.SplitN(entry, "\t", 2)
	if len(parts) == 2 {
		return parts[1]
	}

	return entry
}

func isTreePathExcluded(path, target string, t TargetType) bool {
	norm := filepath.ToSlash(filepath.Clean(path))
	if t == TargetTypeFolder {
		return norm == target || strings.HasPrefix(norm, target+"/")
	}

	return norm == target
}

func buildMktree(repoDir string, entries []string) (string, error) {
	cmd := exec.Command("git", "mktree", "-z")
	cmd.Dir = repoDir
	input := strings.Join(entries, "\x00") + "\x00"
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "mktree failed")
	}

	return strings.TrimSpace(string(out)), nil
}

func makeEmptyTree(repoDir string) (string, error) {
	cmd := exec.Command("git", "mktree", "-z")
	cmd.Dir = repoDir
	cmd.Stdin = strings.NewReader("")
	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "mktree empty failed")
	}

	return strings.TrimSpace(string(out)), nil
}

func remapParentShas(parents []string, remap map[string]string) []string {
	var remapped []string
	seen := make(map[string]bool)
	for _, p := range parents {
		target := p
		if mapped, ok := remap[p]; ok {
			target = mapped
		}
		if !seen[target] && len(target) > 0 {
			seen[target] = true
			remapped = append(remapped, target)
		}
	}

	return remapped
}

func createRewrittenCommit(repoDir, treeSha string, parents []string, raw *rawCommitInfo) (string, error) {
	args := []string{"commit-tree", treeSha}
	for _, p := range parents {
		args = append(args, "-p", p)
	}
	cmd := exec.Command("git", args...)
	cmd.Dir = repoDir
	cmd.Env = append(os.Environ(), raw.authorEnv...)
	cmd.Stdin = strings.NewReader(raw.commitMsg)

	out, err := cmd.Output()
	if err != nil {
		return "", apperror.WrapSimple(err, "commit-tree failed")
	}

	return strings.TrimSpace(string(out)), nil
}

func recordCommitMapEntry(ctx *rewriteContext, origSha, newSha string, raw *rawCommitInfo) {
	parentOrig := ""
	if len(raw.parents) > 0 {
		parentOrig = raw.parents[0]
	}
	m := store.HistoryPurgeCommitMap{
		OperationId:        ctx.opId,
		OriginalCommitSha:  origSha,
		RewrittenCommitSha: newSha,
		ParentOriginalSha:  parentOrig,
		CommitMessage:      raw.commitMsg,
		CreatedAt:          time.Now().Unix(),
	}
	ctx.commitMaps = append(ctx.commitMaps, m)
}

func updateBranchAndReset(repoDir, branch, newHeadSha string) error {
	ref := fmt.Sprintf("refs/heads/%s", branch)
	upCmd := exec.Command("git", "update-ref", ref, newHeadSha)
	upCmd.Dir = repoDir
	if err := upCmd.Run(); err != nil {
		return apperror.WrapSimple(err, "atomic branch update")
	}

	resetCmd := exec.Command("git", "reset", "--hard", newHeadSha)
	resetCmd.Dir = repoDir

	return resetCmd.Run()
}

func finalizePurgeExecution(ctx *rewriteContext, opts PurgeExtendedOptions, rep *PreflightReport, newHead, vault string) (*PurgeExecutionReport, error) {
	_ = ctx.db.InsertHistoryPurgeCommitMapBatch(ctx.opId, ctx.commitMaps)
	expireReflogsAndGC(ctx.repoDir)

	relDeleted, relUpdated := 0, 0
	if opts.IsRelease {
		relSummary, _ := PruneReleaseAssets(ctx.repoDir, []string{opts.TargetPath})
		if relSummary != nil {
			relDeleted = relSummary.AssetsDeletedCount
			relUpdated = relSummary.NotesUpdatedCount
		}
	}

	_ = ctx.db.UpdateHistoryPurgeOperationStatus(store.UpdatePurgeStatusOptions{
		OperationId: ctx.opId,
		IsSuccess:   true,
		IsVerified:  true,
		Notes:       fmt.Sprintf("rewritten %d commits", len(ctx.commitMaps)),
	})

	adv := makeAdvisoryString(ctx.opId)
	if !opts.IsJson {
		printPostPurgeSummary(ctx.opId, newHead, adv)
	}

	return &PurgeExecutionReport{
		OperationId:         ctx.opId,
		RepoSlug:            rep.RepoSlug,
		Branch:              ctx.branch,
		TargetType:          opts.TargetType,
		TargetPath:          opts.TargetPath,
		OriginalHeadSha:     ctx.origHead,
		RewrittenHeadSha:    newHead,
		BackupRef:           fmt.Sprintf("refs/gitmap-backup/%d/%s", ctx.opId, ctx.branch),
		TempBackupDir:       vault,
		ScannedCommits:      rep.TotalCommitsScanned,
		AffectedCommits:     len(rep.AffectedCommits),
		ReleaseAssetsPurged: relDeleted,
		ReleaseNotesUpdated: relUpdated,
		NonWarrantyAdvisory: adv,
		IsSuccess:           true,
	}, nil
}

func expireReflogsAndGC(repoDir string) {
	c1 := exec.Command("git", "reflog", "expire", "--expire=now", "--all")
	c1.Dir = repoDir
	_ = c1.Run()

	c2 := exec.Command("git", "prune", "--expire=now")
	c2.Dir = repoDir
	_ = c2.Run()

	c3 := exec.Command("git", "gc", "--prune=now")
	c3.Dir = repoDir
	_ = c3.Run()
}

func makeAdvisoryString(opId int64) string {
	return fmt.Sprintf("You can undo this if you wanted to. We do not confirm this, but you can try: gitmap history undo %d", opId)
}

func printPostPurgeSummary(opId int64, newHead, advisory string) {
	fmt.Println("┌─ Operation Completed: Git History Purged ──────────────────────────┐")
	fmt.Printf("│ Operation ID      : %-46d │\n", opId)
	fmt.Printf("│ Rewritten HEAD    : %-46s │\n", truncateString(newHead, 46))
	fmt.Println("├────────────────────────────────────────────────────────────────────┤")
	fmt.Printf("│ %-66s │\n", advisory)
	fmt.Println("└────────────────────────────────────────────────────────────────────┘")
}

// doPurge is a backward-compatible wrapper for legacy purge invocations.
func doPurge(db *store.DB, repoPath, pattern string, isAutoConfirm bool) error {
	opts := PurgeExtendedOptions{
		PurgeOptions: PurgeOptions{
			TargetType:    TargetTypeFolder,
			TargetPath:    pattern,
			IsAutoConfirm: isAutoConfirm,
		},
		RepoDir: repoPath,
	}
	_, err := ExecutePurgeEngineExtended(opts)

	return err
}

func normalizePurgePattern(pattern string) string {
	return strings.ReplaceAll(pattern, "\\", "/")
}
