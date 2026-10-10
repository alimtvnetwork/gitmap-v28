// Package cmdsummary implements workspace-wide activity heatmap TreeView and dirty state remediation.
package cmdsummary

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/pterm/pterm"
)

// RunFullSummary executes workspace-wide activity heatmap analysis and TreeView rendering.
func RunFullSummary(args []string) error {
	opts := parseFullSummaryOptions(args)

	cwd, errCwd := os.Getwd()
	if errCwd != nil {
		cwd = "."
	}

	records := resolveWorkspaceRepositories(cwd)
	if len(records) == 0 {
		fmt.Println("[GitMap] No repositories discovered in current workspace.")
		return nil
	}

	cutoffTime := time.Now().Add(-time.Duration(opts.ActivityHours) * time.Hour)
	var activeRepos []RepoSummaryRecord
	dirtyReposCount := 0

	for _, rec := range records {
		repoSummary := evaluateRepoActivity(rec, cutoffTime, opts)
		if repoSummary == nil {
			continue
		}

		if repoSummary.IsDirty {
			dirtyReposCount++
		}
		activeRepos = append(activeRepos, *repoSummary)
	}

	payload := assembleFullSummaryPayload(len(records), activeRepos, dirtyReposCount, opts)

	if opts.IsJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	renderFullSummaryTreeView(payload, opts)
	return nil
}

func parseFullSummaryOptions(args []string) FullSummaryOptions {
	opts := FullSummaryOptions{
		ReleasesCount: 3,
		ActivityHours: 48,
	}

	for _, arg := range args {
		lower := strings.ToLower(arg)
		if lower == "--json" || lower == "-j" {
			opts.IsJSON = true
			continue
		}
		if lower == "--pe" || lower == "+pe" || strings.Contains(lower, "fspe") {
			opts.WithPE = true
			continue
		}
		if lower == "--force-all" || lower == "--all" {
			opts.ForceAll = true
			continue
		}
		if lower == "--verbose" || lower == "-v" {
			opts.IsVerbose = true
			continue
		}
		if n, err := strconv.Atoi(arg); err == nil && n > 0 {
			opts.ReleasesCount = n
			continue
		}
	}

	return opts
}

func evaluateRepoActivity(rec model.ScanRecord, cutoff time.Time, opts FullSummaryOptions) *RepoSummaryRecord {
	dir := rec.AbsolutePath
	if dir == "" {
		dir = rec.RelativePath
	}

	untracked, modified, staged, pendingFiles := queryRepoStatus(dir)
	dirtyFilesCount := untracked + modified + staged
	isDirty := dirtyFilesCount > 0

	lastCommitEpoch, _ := gitExec(dir, "log", "-1", "--format=%ct")
	var lastActivity time.Time
	if epochSec, err := strconv.ParseInt(lastCommitEpoch, 10, 64); err == nil && epochSec > 0 {
		lastActivity = time.Unix(epochSec, 0)
	} else {
		lastActivity = time.Now()
	}

	isRecentlyActive := lastActivity.After(cutoff)

	if !opts.ForceAll && !isDirty && !isRecentlyActive {
		return nil
	}

	remoteURL, _ := gitExec(dir, "remote", "get-url", "origin")
	if remoteURL == "" {
		remoteURL = rec.RepoName
	}
	branch, _ := gitExec(dir, "rev-parse", "--abbrev-ref", "HEAD")
	headHash, _ := gitExec(dir, "rev-parse", "HEAD")

	releases, _ := extractRepoReleaseSummaries(dir, remoteURL, opts.ReleasesCount)

	suggestedCmd := ""
	if isDirty {
		suggestedCmd = fmt.Sprintf("gitmap -C %s cpf \"wip: save active progress\"", dir)
	}

	record := &RepoSummaryRecord{
		RepoName:        rec.RepoName,
		CanonicalSlug:   rec.Slug,
		LocalPath:       dir,
		RemoteURL:       remoteURL,
		CurrentBranch:   branch,
		HeadCommitHash:  headHash,
		IsDirty:         isDirty,
		DirtyFilesCount: dirtyFilesCount,
		PendingFiles:    pendingFiles,
		LastActivityAt:  lastActivity,
		Releases:        releases,
		SuggestedCommit: suggestedCmd,
	}

	if opts.WithPE {
		attachPipelineErrorTelemetry(record)
	}

	_ = UpdateRepoHeadState(remoteURL, rec.Slug, dir, headHash, isDirty, dirtyFilesCount, lastActivity)
	return record
}

func queryRepoStatus(dir string) (int, int, int, []string) {
	porcelain, _ := gitExec(dir, "status", "--porcelain")
	return parsePorcelainStatusLines(porcelain)
}

func assembleFullSummaryPayload(total int, active []RepoSummaryRecord, dirtyCount int, opts FullSummaryOptions) FullSummaryPayload {
	var payload FullSummaryPayload
	payload.Attributes.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	payload.Attributes.GitMapVersion = "6.520.0"
	payload.Attributes.Command = "full-summary"
	if opts.WithPE {
		payload.Attributes.Command = "full-summary+pe"
	}
	payload.Attributes.ActivityHours = opts.ActivityHours

	payload.Data.TotalReposDiscovered = total
	payload.Data.ActiveReposCount = len(active)
	payload.Data.DirtyReposCount = dirtyCount
	if dirtyCount > 0 {
		payload.Data.MasterSanitizeCmd = "gitmap sanitize-all --message \"wip: save active progress across dirty repos\""
	}
	payload.Data.Repositories = active

	return payload
}

func renderFullSummaryTreeView(payload FullSummaryPayload, opts FullSummaryOptions) {
	fmt.Printf("[GitMap] Workspace Full Summary (Activity Window: Last %d Hours)\n", payload.Attributes.ActivityHours)
	fmt.Printf("Discovered %d active/dirty repositories out of %d total repositories.\n\n",
		payload.Data.ActiveReposCount, payload.Data.TotalReposDiscovered)

	if len(payload.Data.Repositories) == 0 {
		fmt.Println("No active or modified repositories found within the activity window.")
		return
	}

	for i, repo := range payload.Data.Repositories {
		isLastRepo := i == len(payload.Data.Repositories)-1
		prefix := "├──"
		if isLastRepo {
			prefix = "└──"
		}

		statusLabel := pterm.Green("[CLEAN]")
		if repo.IsDirty {
			statusLabel = pterm.Yellow(fmt.Sprintf("[DIRTY: %d uncommitted]", repo.DirtyFilesCount))
		}

		fmt.Printf("%s 📁 %s (%s) %s\n", prefix, pterm.Bold.Sprint(repo.RepoName), repo.RemoteURL, statusLabel)
		childIndent := "│   "
		if isLastRepo {
			childIndent = "    "
		}

		if repo.IsDirty && len(repo.PendingFiles) > 0 {
			fmt.Printf("%s├── Pending Changes:\n", childIndent)
			for _, file := range repo.PendingFiles {
				fmt.Printf("%s│   ├── %s\n", childIndent, file)
			}
		}

		if len(repo.Releases) > 0 {
			fmt.Printf("%s├── Releases:\n", childIndent)
			for j, rel := range repo.Releases {
				relPrefix := "├──"
				if j == len(repo.Releases)-1 {
					relPrefix = "└──"
				}
				checkMark := pterm.Green("✓")
				fmt.Printf("%s│   %s %s %s (%s): %s\n", childIndent, relPrefix, checkMark, rel.TagName, rel.ReleaseDate, rel.SummaryGist)
			}
		}

		if repo.PipelineError != nil {
			pe := repo.PipelineError
			fmt.Printf("%s├── ❌ CI/CD Pipeline FAILED (%s)\n", childIndent, pe.WorkflowName)
			fmt.Printf("%s│   Job: %s | Step: %s | Exit Code: %d\n", childIndent, pe.JobName, pe.StepName, pe.ExitCode)
			if pe.ErrorSummary != "" {
				fmt.Printf("%s│   Summary: %s\n", childIndent, pe.ErrorSummary)
			}
			if pe.StackTrace != "" {
				fmt.Printf("%s│   Stack Trace (last 25 lines):\n", childIndent)
				lines := strings.Split(pe.StackTrace, "\n")
				for _, line := range lines {
					if strings.TrimSpace(line) != "" {
						fmt.Printf("%s│     %s\n", childIndent, line)
					}
				}
			}
		} else if opts.WithPE && repo.IsPipelineClean {
			fmt.Printf("%s├── ✓ CI/CD Pipeline PASSING\n", childIndent)
		}

		if repo.IsDirty && repo.SuggestedCommit != "" {
			fmt.Printf("%s└── Suggested Fix: %s\n", childIndent, repo.SuggestedCommit)
		}

		fmt.Println()
	}

	if payload.Data.DirtyReposCount > 0 {
		fmt.Println("========================================================================================")
		fmt.Printf("[Pending Work Notice]: %d repositories have uncommitted changes.\n", payload.Data.DirtyReposCount)
		fmt.Println("To commit all pending changes in one shot, run:")
		fmt.Printf("  %s\n", payload.Data.MasterSanitizeCmd)
		fmt.Println("========================================================================================")
	}
}
