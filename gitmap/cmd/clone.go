package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloneconcurrency"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/verbose"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/vscodepm"
)

// applySSHKey sets GIT_SSH_COMMAND if an SSH key name is provided.
func applySSHKey(name string) {
	if len(name) == 0 {
		return
	}

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHQuery, err)
		os.Exit(1)
	}
	defer db.Close()

	key, err := db.FindSSHKeyByName(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSSHNotFound, name)
		os.Exit(1)
	}

	sshCmd := fmt.Sprintf("ssh -i %s -o IdentitiesOnly=yes", key.PrivatePath)
	os.Setenv("GIT_SSH_COMMAND", sshCmd)
	fmt.Fprintf(os.Stdout, constants.MsgSSHCloneUsing, name, key.PrivatePath)
}

// runClone handles the "clone" subcommand.
func runClone(args []string) {
	checkHelp("clone", args)
	cf := parseCloneFlags(args)
	if len(cf.Source) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrSourceRequired)
		fmt.Fprintln(os.Stderr, constants.ErrCloneUsage)
		os.Exit(1)
	}
	initCloneVerbose(cf.Verbose)
	SetCloneDryRun(cf.DryRun)
	SetCloneAssumeYes(cf.IsAssumeYes)
	setCmdFaithfulVerify(cf.VerifyCmdFaithful)
	setCmdFaithfulExitOnMismatch(cf.VerifyCmdFaithfulExitOnMismatch)
	setCmdPrintArgv(cf.PrintCloneArgv)

	// Audit short-circuits all execution paths. It must run BEFORE
	// requireOnline / SSH key application so users can audit a manifest
	// while offline and without unlocking SSH agents.
	if cf.Audit {
		runCloneAudit(cf)
		maybeExitOnCmdFaithfulMismatch()

		return
	}

	requireOnline()
	applySSHKey(cf.SSHKeyName)
	applyCloneAssumeYesEnv(cf.IsAssumeYes)
	cf = applyURLSchemeFlags(cf)

	// Multi-URL form: any positional arg containing a comma, OR 2+ positional
	// args where the second one looks like a URL. This catches PowerShell's
	// silent comma-splitting of unquoted args (root cause of v3.78 regression).
	if shouldUseMultiClone(cf) {
		runCloneMulti(cf)
		maybeExitOnCmdFaithfulMismatch()

		return
	}

	if isDirectURL(cf.Source) {
		executeDirectClone(cf.Source, cf.FolderName, cf.GHDesktop, cf.NoReplace, cf.Output, cf.NoVSCodeSync)
		maybeExitOnCmdFaithfulMismatch()

		return
	}

	source := resolveCloneShorthand(cf.Source)
	executeClone(source, cf.TargetDir, cf.SafePull, cf.GHDesktop, cf.MaxConcurrency, cf.DefaultBranch, cf.NoVSCodeSync)
	maybeExitOnCmdFaithfulMismatch()
}

// shouldUseMultiClone returns true when the positional args describe a
// batch of URLs rather than a single source + optional folder name.
// Three triggers (any one is sufficient):
//  1. Any positional arg contains a list separator (`,` or `;`) — the
//     user explicitly listed URLs, even if PowerShell didn't pre-split.
//  2. 2+ positional args AND any arg beyond the first parses as a URL
//     — covers PowerShell's silent comma-split into separate argv slots
//     AND the `clone url1 url2 url3` space-only form.
//  3. The first arg flattens (after sanitisation) to 2+ valid URLs —
//     covers `clone "url1,url2"` where the whole list is one token.
func shouldUseMultiClone(cf CloneFlags) bool {
	for _, p := range cf.Positional {
		if strings.ContainsAny(p, urlListSeparators) {
			return true
		}
	}
	if len(cf.Positional) >= 2 {
		for _, p := range cf.Positional[1:] {
			if isDirectURL(sanitizeURLToken(p)) {
				return true
			}
		}
	}
	if len(cf.Positional) >= 1 {
		flat := flattenURLArgs(cf.Positional[:1])
		urlCount := 0
		for _, u := range flat {
			if isDirectURL(u) {
				urlCount++
			}
		}
		if urlCount >= 2 {
			return true
		}
	}

	return false
}

// runCloneMulti clones every URL in the flattened positional list, continuing
// on per-URL failure. Folder name is ignored in this mode (each repo lands in
// its own auto-derived folder). Exit codes follow mem://features/clone-multi.
func runCloneMulti(cf CloneFlags) {
	flat := flattenURLArgs(cf.Positional)
	urls, invalid := classifyURLs(flat)

	if len(urls) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrCloneAllInvalid)
		os.Exit(constants.ExitCloneMultiAllInvalid)
	}

	fmt.Printf(constants.MsgCloneMultiBegin, len(urls))

	succeeded := 0
	failed := 0
	pmPairs := make([]vscodepm.Pair, 0, len(urls))

	for idx, url := range urls {
		fmt.Printf(constants.MsgCloneMultiItem, idx+1, len(urls), url)

		// `--output terminal`: stream the standardized RepoTermBlock
		// BEFORE shelling out so the user sees branch/from/to/command
		// for THIS repo before its clone progress, then the next
		// repo's block after this one finishes. Mirrors the streamed
		// emission contract locked in chat.
		printCloneTermBlockForURL(cf.Output, idx+1, url, "")

		if err := executeDirectCloneOne(url, "", cf.GHDesktop, cf.NoReplace); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneMultiFailedFmt, idx+1, len(urls), url, err)
			failed++

			continue
		}
		succeeded++
		// Build a PM pair for the URL we just cloned. Mirrors the
		// folder resolution executeDirectCloneOne uses internally so
		// the absPath we hand to projects.json matches what landed
		// on disk.
		repoName := repoNameFromURL(url)
		folder := resolveCloneFolder(repoName, "")
		if abs, absErr := filepath.Abs(folder); absErr == nil {
			pmPairs = append(pmPairs, buildClonePMPair(abs, repoName))
		}
	}

	failed += len(invalid)

	fmt.Printf(constants.MsgCloneSummaryMultiFmt, succeeded, failed, len(urls)+len(invalid))

	// Single projects.json transaction for the whole batch. Soft-
	// fails when VS Code or the extension is missing. Skipped when
	// every per-URL clone failed (pmPairs is empty).
	syncClonedReposToVSCodePM(pmPairs, cf.NoVSCodeSync)

	if failed > 0 {
		os.Exit(constants.ExitCloneMultiPartialFail)
	}
}

// isDirectURL returns true when source is a git URL (not a file path).
// Accepts HTTPS, HTTP, SSH (`ssh://`), and SSH-shorthand (`git@host:owner/repo`).
// Kept in lockstep with isLikelyURL in rootflags.go so folder-name
// disambiguation and URL classification never disagree.
func isDirectURL(source string) bool {
	lower := strings.ToLower(source)
	if strings.HasPrefix(lower, constants.PrefixHTTPS) ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, constants.PrefixSSH) {
		return true
	}
	// SSH shorthand: git@host:owner/repo(.git)?  — must contain `:` after `@`.
	if strings.HasPrefix(lower, "git@") {
		at := strings.Index(lower, "@")
		colon := strings.Index(lower[at:], ":")

		return colon > 0
	}

	return false
}

// repoNameFromURL derives the repository name from a clone URL.
//
// Trailing slashes (and backslashes) are stripped first so URLs like
// `https://github.com/owner/repo/` or `git@host:owner/repo.git/` still
// resolve to `repo`. Without this, the basename would collapse to ""
// and the clone target would be the caller's CWD — which then trips
// the "target exists" replace flow and clobbers unrelated work.
func repoNameFromURL(url string) string {
	name := strings.TrimRight(url, "/\\")
	name = strings.TrimSuffix(name, ".git")
	name = strings.TrimRight(name, "/\\")
	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		name = name[idx+1:]
	}
	if idx := strings.LastIndex(name, ":"); idx >= 0 {
		name = name[idx+1:]
	}

	return name
}

// executeDirectClone clones a single repo from a direct URL.
// When no folder name is given, versioned URLs are auto-flattened
// (e.g., wp-onboarding-v13 clones into wp-onboarding/).
// By default, an existing target folder is replaced via the two-strategy
// flow in spec/01-app/96-clone-replace-existing-folder.md. Pass noReplace=true
// to restore the strict abort-on-exists behavior.
//
// noVSCodeSync, when true, skips the post-clone update of the
// alefragnani.project-manager projects.json file. See
// spec/01-vscode-project-manager-sync/02-clone-sync.md.
func executeDirectClone(url, folderName string, ghDesktopFlag, noReplace bool, output string, noVSCodeSync bool) {
	escapeNestedGitRepo()
	repoName := repoNameFromURL(url)

	if len(folderName) == 0 {
		// Local folder mirrors the URL verbatim (including any `-vN`
		// suffix). The previous behavior auto-flattened versioned
		// URLs (codex-june-6-v1 → codex-june-6) which surprised users
		// running `gitmap clone <url>` — they expect the folder to
		// match what they typed. Version-bump flattening lives in
		// `gitmap clone-next`, not here.
		folderName = repoName
	}

	// Defensive guard: if the resolved folder name itself looks like a URL,
	// the caller dispatched the wrong path — almost always because the user
	// is running a stale binary that pre-dates v3.80.0's multi-URL routing.
	// Refuse to build `D:\...\https:\github.com\...` paths that git can't
	// possibly create, and tell the user exactly why.
	if isDirectURL(folderName) {
		fmt.Fprintf(os.Stderr, constants.ErrCloneStaleBinaryFolderURL, folderName, constants.Version)
		os.Exit(1)
	}

	absPath, err := filepath.Abs(folderName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", folderName, err)
		absPath = folderName
	}

	// Strict mode: keep the original abort-on-exists behavior.
	if noReplace {
		if _, statErr := os.Stat(absPath); statErr == nil {
			fmt.Fprintf(os.Stderr, constants.ErrCloneURLExists, absPath)
			os.Exit(1)
		}
	}

	// Enqueue pending task.
	workDir, _ := os.Getwd()
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeClone, absPath, workDir, "clone", cmdArgs)

	// `--output terminal`: emit the standardized per-repo block to
	// stdout BEFORE the legacy "Cloning ..." line so the user sees
	// branch/from/to/command up-front. No-op when output is empty,
	// which preserves byte-identical legacy output.
	printCloneTermBlockForURL(output, 1, url, absPath)

	// Clone (default: replace; with --no-replace: clone into a guaranteed-empty target).
	fmt.Printf(constants.MsgCloneURLCloning, repoName, folderName)

	if noReplace {
		if cloneErr := runCloneCommand(url, absPath); cloneErr != nil {
			failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneURLFailed, url, cloneErr))
			closeTaskDB(taskDB)
			fmt.Fprintf(os.Stderr, constants.ErrCloneURLFailed, url, cloneErr)
			os.Exit(1)
		}
	} else {
		if _, replaceErr := cloneReplacing(url, absPath); replaceErr != nil {
			failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneURLFailed, url, replaceErr))
			closeTaskDB(taskDB)
			fmt.Fprintf(os.Stderr, constants.ErrCloneURLFailed, url, replaceErr)
			os.Exit(1)
		}
	}

	fmt.Printf(constants.MsgCloneURLDone, repoName)

	// Upsert to database.
	upsertDirectClone(url, repoName, folderName, absPath)

	// GitHub Desktop registration (auto-register by default for direct URL).
	registerSingleDesktop(repoName, absPath)

	// Shell handoff: cd the parent shell into the freshly cloned folder
	// when invoked via the wrapper function (mirrors `cn` and `cd`).
	// Only fires for the single-repo direct-URL path — runCloneMulti
	// deliberately skips handoff because the destination is ambiguous.
	WriteShellHandoff(absPath)

	// Open in VS Code if available.
	openInVSCode(absPath)

	// VS Code Project Manager: register the freshly-cloned repo so
	// it appears in the sidebar without a separate `gitmap code`
	// step. Soft-fails when VS Code or the extension is missing.
	syncSingleClonedRepoToVSCodePM(absPath, repoName, noVSCodeSync)

	completePendingTask(taskDB, taskID)
	closeTaskDB(taskDB)
}

// upsertDirectClone persists the cloned repo in the database.
func upsertDirectClone(url, repoName, folderName, absPath string) {
	rec := model.ScanRecord{
		Slug:         strings.ToLower(repoName),
		RepoName:     repoName,
		RelativePath: folderName,
		AbsolutePath: absPath,
	}
	populateDirectCloneURLs(&rec, url)

	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not open database: %v\n", err)

		return
	}
	defer db.Close()

	if upsertErr := db.UpsertRepos([]model.ScanRecord{rec}); upsertErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not save repo to database: %v\n", upsertErr)

		return
	}
	if markErr := db.MarkCloned(absPath); markErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not stamp clone time: %v\n", markErr)
	}
}

// registerSingleDesktop registers a single repo with GitHub Desktop.
func registerSingleDesktop(name, absPath string) {
	records := []model.ScanRecord{{
		RepoName:     name,
		AbsolutePath: absPath,
	}}
	result := desktop.AddRepos(records)
	if result.Added > 0 {
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}

// initCloneVerbose sets up verbose logging if enabled.
func initCloneVerbose(enabled bool) {
	if enabled {
		log, err := verbose.Init()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.WarnVerboseLogFailed, err)

			return
		}
		defer log.Close()
	}
}

// resolveCloneShorthand maps "json", "csv", and "text" to default output paths.
func resolveCloneShorthand(source string) string {
	shorthandMap := map[string]string{
		constants.ShorthandJSON: filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile),
		constants.ShorthandCSV:  filepath.Join(constants.DefaultOutputFolder, constants.DefaultCSVFile),
		constants.ShorthandText: filepath.Join(constants.DefaultOutputFolder, constants.DefaultTextFile),
	}
	resolved, ok := shorthandMap[strings.ToLower(source)]
	if ok {
		return validateShorthandPath(resolved)
	}

	return source
}

// validateShorthandPath checks that the resolved shorthand file exists.
func validateShorthandPath(resolved string) string {
	_, err := os.Stat(resolved)
	if err == nil {
		return resolved
	}
	fmt.Fprintf(os.Stderr, constants.ErrShorthandNotFound, resolved)
	os.Exit(1)

	return ""
}

// executeClone runs the clone operation and prints the summary.
//
// maxConcurrency is the worker count plumbed in from --max-concurrency.
// Values <= 1 keep the legacy sequential runner; > 1 enables the
// bounded worker pool in gitmap/cloner/concurrent.go. The on-disk
// nested folder hierarchy is preserved at any N because each repo
// still lands at filepath.Join(targetDir, rec.RelativePath).
//
// defaultBranch is the optional `--default-branch` fallback. Empty
// keeps the legacy "remote default HEAD" behavior for rows with an
// untrustworthy Branch / BranchSource. Non-empty rewrites those rows
// in cloner.applyDefaultBranchFallback so they go through the
// trusted `git clone -b <fallback>` path.
func executeClone(source, targetDir string, safePull, ghDesktop bool, maxConcurrency int, defaultBranch string, noVSCodeSync bool) {
	workers, ok := cloneconcurrency.Resolve(maxConcurrency)
	if !ok {
		fmt.Fprintf(os.Stderr, constants.ErrCloneMaxConcurrencyInvalid, maxConcurrency)
		os.Exit(1)
	}
	maxConcurrency = workers

	// Enqueue clone as a pending task before execution.
	absTarget, absErr := filepath.Abs(targetDir)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", targetDir, absErr)
		absTarget = targetDir
	}
	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  Warning: could not determine working directory: %v\n", wdErr)
	}
	cmdArgs := buildCommandArgs(append([]string{"clone"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeClone, absTarget, workDir, "clone", cmdArgs)

	summary, err := cloner.CloneFromFileWithOptions(source, targetDir, cloner.CloneOptions{
		SafePull:       safePull,
		MaxConcurrency: maxConcurrency,
		DefaultBranch:  defaultBranch,
	})
	if err != nil {
		failPendingTask(taskDB, taskID, fmt.Sprintf(constants.ErrCloneFailed, source, err))
		closeTaskDB(taskDB)
		fmt.Fprintf(os.Stderr, constants.ErrCloneFailed, source, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgCloneComplete, summary.Succeeded, summary.Failed)
	printCloneFailures(summary)
	registerCloned(summary, targetDir, ghDesktop)

	// VS Code Project Manager: build one pair per successfully
	// cloned repo and run a single Sync. Mirrors registerCloned's
	// abs-path resolution so projects.json points at the same path
	// the GitHub Desktop registration uses.
	syncManifestClonedReposToVSCodePM(summary, targetDir, noVSCodeSync)

	// Mark clone task as completed after all steps succeed.
	completePendingTask(taskDB, taskID)
	closeTaskDB(taskDB)
}

// syncManifestClonedReposToVSCodePM converts a manifest-style
// CloneSummary into VS Code Project Manager pairs and syncs them in
// one shot. Skipped when summary has zero successes (avoids a noisy
// "0 added, 0 updated" line).
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

// printCloneFailures lists any repos that failed to clone.
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

// registerCloned adds successfully cloned repos to GitHub Desktop.
func registerCloned(s model.CloneSummary, targetDir string, enabled bool) {
	if enabled {
		absTarget, absErr := filepath.Abs(targetDir)
		if absErr != nil {
			fmt.Fprintf(os.Stderr, "  Warning: could not resolve absolute path for %s: %v\n", targetDir, absErr)
			absTarget = targetDir
		}
		records := make([]model.ScanRecord, 0, s.Succeeded)
		for _, r := range s.Cloned {
			r.Record.AbsolutePath = filepath.Join(absTarget, model.CleanRelativePath(r.Record.RelativePath))
			records = append(records, r.Record)
		}
		result := desktop.AddRepos(records)
		fmt.Printf(constants.MsgDesktopSummary, result.Added, result.Failed)
	}
}

// applyURLSchemeFlags rewrites cf.Source and every positional URL via
// ConvertURLToSSH / ConvertURLToHTTPS when the user passes `--ssh` or
// `--https`. Non-URL positionals (folder names, manifest shorthands
// like `json` / `csv`) are passed through unchanged so a stray flag
// can't corrupt a manifest-style invocation.
//
// `--ssh` and `--https` are mutually exclusive — when both are set,
// `--ssh` wins and a one-line stderr warning is printed so the user
// can spot the conflict.
//
// Spec: spec/01-app/110-clone-ssh-flag.md
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
		if !isDirectURL(in) {
			return in
		}
		if toSSH {
			if out, ok := ConvertURLToSSH(in); ok {
				return out
			}
		} else {
			if out, ok := ConvertURLToHTTPS(in); ok {
				return out
			}
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

// schemeLabel returns the human-readable label for the active rewrite.
func schemeLabel(toSSH bool) string {
	if toSSH {
		return "--ssh"
	}

	return "--https"
}
