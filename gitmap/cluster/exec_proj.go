package cluster

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/db"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

const (
	fileRunPs1         = "run.ps1"
	fileRunSh          = "run.sh"
	cmdPwsh            = "pwsh"
	flagNonInteractive = "-NonInteractive"
	flagFile           = "-File"
	cmdSh              = "sh"
	msgNoRunScript     = "neither run.ps1 nor run.sh found"
	msgProjectNotFound = "project %q not found"
	msgCreateCICDStub  = "create-cicd: reserved for future spec"
	maxErrorLines      = 20
	newline            = "\n"
	emptyString        = ""
)

type ProjRunResult struct {
	ProjectName string
	Succeeded   bool
	Stderr      string
}

// ExecProjRun scans registered repo paths on the node, finds run.ps1 or run.sh, executes, captures last 20 lines on failure.

func ExecProjRun(ctx context.Context, node ClusterNode, projectNames []string) ([]ProjRunResult, error) {
	database, err := store.OpenDefault()
	isDbError := err != nil
	if isDbError == true {
		return nil, fmt.Errorf("failed to open store: %w", err)
	}
	defer database.Close()

	records, err := database.ListRepos()
	isListError := err != nil
	if isListError == true {
		return nil, fmt.Errorf("failed to list repos: %w", err)
	}

	var results []ProjRunResult
	for _, projName := range projectNames {
		var foundPath string
		for _, rec := range records {
			isMatch := strings.EqualFold(rec.RepoName, projName) || strings.EqualFold(rec.Slug, projName)
			if isMatch == true {
				foundPath = rec.AbsolutePath
				break
			}
		}

		isMissing := foundPath == emptyString
		if isMissing {
			results = append(results, ProjRunResult{
				ProjectName: projName,
				Succeeded:   false,
				Stderr:      fmt.Sprintf(msgProjectNotFound, projName),
			})
			continue
		}

		res := executeProjectRun(ctx, projName, foundPath)
		results = append(results, res)
	}

	return results, nil
}

func executeProjectRun(ctx context.Context, projName, absPath string) ProjRunResult {
	ps1Path := filepath.Join(absPath, fileRunPs1)
	shPath := filepath.Join(absPath, fileRunSh)

	var cmd *exec.Cmd

	_, errPs1 := os.Stat(ps1Path)
	hasPs1 := errPs1 == nil
	if hasPs1 {
		cmd = exec.CommandContext(ctx, cmdPwsh, flagNonInteractive, flagFile, ps1Path)
	}

	_, errSh := os.Stat(shPath)
	hasSh := errSh == nil
	if !hasPs1 && hasSh {
		cmd = exec.CommandContext(ctx, cmdSh, shPath)
	}

	if !hasPs1 && !hasSh {
		return ProjRunResult{
			ProjectName: projName,
			Succeeded:   false,
			Stderr:      msgNoRunScript,
		}
	}

	cmd.Dir = absPath
	output, err := cmd.CombinedOutput()
	isSuccess := err == nil
	if isSuccess == true {
		return ProjRunResult{
			ProjectName: projName,
			Succeeded:   true,
			Stderr:      emptyString,
		}
	}

	lines := strings.Split(strings.TrimSpace(string(output)), newline)
	totalLines := len(lines)
	isOverMax := totalLines > maxErrorLines
	if isOverMax == true {
		lines = lines[totalLines-maxErrorLines:]
	}

	errMsg := err.Error()
	combinedStderr := fmt.Sprintf("%s%s%s", errMsg, newline, strings.Join(lines, newline))

	return ProjRunResult{
		ProjectName: projName,
		Succeeded:   false,
		Stderr:      strings.TrimSpace(combinedStderr),
	}
}

// ExecProjCreateCICD is a stub returning db.ResultStatusDeferred with a future spec message.

func ExecProjCreateCICD(ctx context.Context, node ClusterNode, projectNames []string) (db.ResultStatusType, string) {
	return db.ResultStatusDeferred, msgCreateCICDStub
}
