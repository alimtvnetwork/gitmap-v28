// Package cmdclone — clone_manifest.go handles manifest-based multi-repo cloning operations.
package cmdclone

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/cli/cloner"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdscan"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func initManifestCloneTask(cf CloneFlags) (int64, *store.DB) {
	absTarget, _ := filepath.Abs(cf.TargetDir)
	workDir, _ := os.Getwd()
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))

	return createPendingTask(constants.TaskTypeClone, absTarget, workDir, "clone", cmdArgs)
}

func handleManifestCloneFailure(taskDB *store.DB, taskID int64, source string, err error) error {
	failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneFailed, source, err))
	closeTaskDB(taskDB)
	fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, source, err)

	return err
}

// executeClone runs the clone operation and prints the summary.
func executeClone(cf CloneFlags) error {
	workers, ok := cloneconcurrency.Resolve(cf.MaxConcurrency)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrCloneMaxConcurrencyInvalid, cf.MaxConcurrency)

		return apperror.NewValidationError(fmt.Sprintf(constants.ErrCloneMaxConcurrencyInvalid, cf.MaxConcurrency))
	}

	cf.MaxConcurrency = workers
	taskID, taskDB := initManifestCloneTask(cf)

	summary, err := runCloneExecution(cf)
	if err != nil {
		return handleManifestCloneFailure(taskDB, taskID, cf.Source, err)
	}

	finalizeCloneExecution(summary, cf, taskDB, taskID)

	return nil
}

func finalizeCloneExecution(summary model.CloneSummary, cf CloneFlags, taskDB *store.DB, taskID int64) {
	fmt.Printf(constants.MsgCloneComplete, summary.Succeeded, summary.Failed)
	printCloneFailures(summary)
	registerCloned(summary, cf.TargetDir, cf.GHDesktop)
	syncClonedReposToRepoCache(summary, cf.TargetDir)
	syncManifestClonedReposToVSCodePM(summary, cf.TargetDir, cf.NoVSCodeSync)
	completePendingTask(taskDB, taskID)
	closeTaskDB(taskDB)

	fmt.Println("\nRunning gitmap status on target directory...")
	if err := os.Chdir(cf.TargetDir); err == nil {
		RunStatus([]string{})
	}
}

func syncClonedReposToRepoCache(summary model.CloneSummary, targetDir string) {
	if summary.Succeeded == 0 {
		return
	}
	records := make([]model.ScanRecord, 0, summary.Succeeded)
	for _, r := range summary.Cloned {
		rec := r.Record
		if len(rec.AbsolutePath) == 0 {
			rec.AbsolutePath = filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		}
		records = append(records, rec)
	}
	cmdscan.TryIncludeRepoInRepoCacheJSON(records)
}

func runCloneExecution(cf CloneFlags) (model.CloneSummary, error) {
	opts := cloner.CloneOptions{
		IsSafePull:     cf.SafePull,
		MaxConcurrency: cf.MaxConcurrency,
		DefaultBranch:  cf.DefaultBranch,
		IsClean:        cf.Clean,
		IsMissingOnly:  cf.MissingOnly,
		IsUseSSH:       cf.UseSSH,
	}
	records, err := cloner.LoadRecords(cf.Source)
	if err != nil {
		return model.CloneSummary{}, err
	}
	records = filterRecordsByOnly(records, cf.OnlyFilter)
	records = filterRecordsByExclude(records, cf.ExcludeFilter)
	if countMissingCloneRecords(records, cf.TargetDir) == 0 && !cf.Clean {
		fmt.Printf("✓ All %d repository(ies) in %s already exist on disk. Nothing to clone.\n", len(records), cf.Source)

		return model.CloneSummary{Succeeded: len(records)}, nil
	}
	records, err = prepareManifestRecordsAuth(records, cf)
	if err != nil {
		return model.CloneSummary{}, err
	}
	if len(records) == 0 {
		fmt.Println("Warning: no repositories matched clone filters")

		return model.CloneSummary{}, nil
	}

	return cloner.CloneRecords(records, cf.TargetDir, opts), nil
}

func prepareManifestRecordsAuth(records []model.ScanRecord, cf CloneFlags) ([]model.ScanRecord, error) {
	if cf.UseSSH {
		return convertRecordsToSSH(records), nil
	}

	return resolveAuthForRecords(records, cf.TargetDir, cf.MissingOnly || !cf.Clean)
}

func countMissingCloneRecords(records []model.ScanRecord, targetDir string) int {
	missing := 0
	for _, r := range records {
		dest := filepath.Join(targetDir, r.RelativePath)
		if r.RelativePath == "" {
			dest = filepath.Join(targetDir, r.RepoName)
		}
		if !isGitRepo(dest) {
			missing++
		}
	}

	return missing
}

func resolveAuthForRecords(records []model.ScanRecord, targetDir string, skipExisting bool) ([]model.ScanRecord, error) {
	out := make([]model.ScanRecord, 0, len(records))
	for _, r := range records {
		resolved, hasResolved, err := resolveSingleRecordAuth(r, targetDir, skipExisting)
		if err != nil {
			return nil, err
		}
		if hasResolved {
			out = append(out, resolved)
		}
	}

	return out, nil
}

func resolveSingleRecordAuth(r model.ScanRecord, targetDir string, skipExisting bool) (model.ScanRecord, bool, error) {
	dest := filepath.Join(targetDir, r.RelativePath)
	if r.RelativePath == "" {
		dest = filepath.Join(targetDir, r.RepoName)
	}

	if skipExisting && isGitRepo(dest) {
		return r, true, nil
	}

	resolved, err := ResolveRepoAuth(r)
	if err == ErrAuthSkipped {
		fmt.Printf("  Skipping %s per user request.\n", getRecordDisplayName(r))

		return model.ScanRecord{}, false, nil
	}

	if errors.Is(err, ErrAuthNonInteractive) {
		fmt.Fprintf(os.Stderr, "\n%s\n", err.Error())

		return model.ScanRecord{}, false, err
	}

	return resolved, true, nil
}

func pickResolvedURL(rec model.ScanRecord, fallback string) string {
	if rec.Transport == "ssh" && len(rec.SSHUrl) > 0 {
		return rec.SSHUrl
	}
	if len(rec.HTTPSUrl) > 0 {
		return rec.HTTPSUrl
	}

	return fallback
}

func convertRecordsToSSH(records []model.ScanRecord) []model.ScanRecord {
	for i, r := range records {
		records[i].Transport = "ssh"
		if sshURL, ok := ConvertURLToSSH(resolveRecordRemoteURL(r)); ok {
			records[i].SSHUrl = sshURL
		}
	}
	return records
}

func syncManifestClonedReposToVSCodePM(summary model.CloneSummary, targetDir string, skip bool) {
	if summary.Succeeded == 0 {
		return
	}

	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		absTarget = targetDir
	}

	pairs := make([]vscodepm.Pair, 0, summary.Succeeded)
	for _, r := range summary.Cloned {
		abs := filepath.Join(absTarget, model.CleanRelativePath(r.Record.RelativePath))
		pairs = append(pairs, buildClonePMPair(abs, r.Record.RepoName))
	}

	syncClonedReposToVSCodePM(pairs, skip)
}

func printCloneFailures(s model.CloneSummary) {
	if s.Failed == 0 {
		return
	}

	fmt.Println(constants.MsgFailedClones)
	for _, e := range s.Errors {
		fmt.Printf(constants.MsgFailedEntry,
			e.Record.RepoName, e.Record.RelativePath, e.Error)
	}
}

func registerCloned(s model.CloneSummary, targetDir string, enabled bool) {
	if !enabled {
		return
	}

	absTarget, absErr := filepath.Abs(targetDir)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", targetDir, absErr)
		absTarget = targetDir
	}

	if s.Succeeded == 0 {
		return
	}

	records := make([]model.ScanRecord, 0, s.Succeeded)
	for _, r := range s.Cloned {
		r.Record.AbsolutePath = filepath.Join(absTarget, model.CleanRelativePath(r.Record.RelativePath))
		records = append(records, r.Record)
	}

	result := desktop.AddRepos(records)
	fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
}

func applyURLSchemeFlags(cf CloneFlags) CloneFlags {
	if !cf.UseSSH && !cf.UseHTTPS {
		return cf
	}

	toSSH := cf.UseSSH
	if cf.UseSSH && cf.UseHTTPS {
		fmt.Fprintln(os.Stderr,
			"  Warning: --ssh and --https are mutually exclusive; --ssh wins.")
	}

	rewrite := func(in string) string {
		isNonDirectURL := !isDirectURL(in)
		if isNonDirectURL {
			return in
		}

		out, ok := "", false
		if toSSH {
			out, ok = ConvertURLToSSH(in)
		} else {
			out, ok = ConvertURLToHTTPS(in)
		}

		if ok {
			return out
		}

		return in
	}

	before := cf.Source
	cf.Source = rewrite(cf.Source)
	for i, p := range cf.Positional {
		cf.Positional[i] = rewrite(p)
	}

	if before != cf.Source {
		fmt.Printf("  ↪ %s rewrite: %s → %s\n", schemeLabel(toSSH), before, cf.Source)
	}

	return cf
}

func schemeLabel(toSSH bool) string {
	if toSSH {
		return "--ssh"
	}

	return "--https"
}
